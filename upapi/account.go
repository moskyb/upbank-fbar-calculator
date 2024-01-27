package upapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/moskyb/upbank-fbar-calculator/qparam"
)

type Account struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		DisplayName   string    `json:"displayName"`
		AccountType   string    `json:"accountType"`
		OwnershipType string    `json:"ownershipType"`
		Balance       Money     `json:"balance"`
		CreatedAt     time.Time `json:"createdAt"`
	} `json:"attributes"`
}

type ListAccountsParams struct {
	AccountType string `qparam:"filter[accountType]"`
	Ownership   string `qparam:"filter[ownershipType]"`

	Before string `qparam:"page[before]"`
	After  string `qparam:"page[after]"`
}

func (c *Client) PaginateAllAccounts(ctx context.Context, params ListAccountsParams) ([]Account, error) {
	var accounts []Account

	for {
		resp, err := c.ListAccounts(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list accounts: %w", err)
		}

		accounts = append(accounts, resp.Data...)

		if resp.Links.Next == nil {
			break
		}

		nextURL, err := url.Parse(*resp.Links.Next)
		if err != nil {
			return nil, fmt.Errorf("failed to parse next URL: %w", err)
		}

		after := nextURL.Query().Get("page[after]")
		params.After = after
	}

	return accounts, nil
}

func (c *Client) ListAccounts(ctx context.Context, params ListAccountsParams) (*Response[[]Account], error) {
	url, err := c.buildURL("accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	query, err := qparam.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode query params: %w", err)
	}

	req.URL.RawQuery = qparam.Merge(req.URL.Query(), query).Encode()

	body, err := c.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	resp := Response[[]Account]{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Response[Account], error) {
	url, err := c.buildURL(fmt.Sprintf("accounts/%s", id))
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	body, err := c.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	resp := Response[Account]{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}
