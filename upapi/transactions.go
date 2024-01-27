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

type Transaction struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		Status          string  `json:"status"`
		RawText         *string `json:"rawText,omitempty"`
		Description     string  `json:"description"`
		Message         *string `json:"message,omitempty"`
		IsCategorizable bool    `json:"isCategorizable"`

		HoldInfo *struct {
			Amount        Money  `json:"amount"`
			ForeignAmount *Money `json:"foreignAmount,omitempty"`
		} `json:"holdInfo,omitempty"`

		RoundUp *struct {
			Amount       Money  `json:"amount"`
			BoostPortion *Money `json:"boostPortion,omitempty"`
		} `json:"roundUp,omitempty"`

		Cashback *struct {
			Description string `json:"description"`
			Amount      Money  `json:"amount"`
		} `json:"cashback,omitempty"`

		Amount             Money  `json:"amount"`
		ForeignAmount      *Money `json:"foreignAmount,omitempty"`
		CardPurchaseMethod *struct {
			Method           string  `json:"method"`
			CardNumberSuffix *string `json:"cardNumberSuffix,omitempty"`
		} `json:"cardPurchaseMethod,omitempty"`

		SettledAt *time.Time `json:"settledAt,omitempty"`
		CreatedAt time.Time  `json:"createdAt"`
	} `json:"attributes"`
}

type ListTransactionsParams struct {
	Status   string    `qparam:"filter[status]"`
	Since    time.Time `qparam:"filter[since]"`
	Until    time.Time `qparam:"filter[until]"`
	Category string    `qparam:"filter[category]"`
	Tag      string    `qparam:"filter[tag]"`

	Before string `qparam:"page[before]"`
	After  string `qparam:"page[after]"`
}

func (c *Client) PaginateAllTransactions(ctx context.Context, params ListTransactionsParams) ([]Transaction, error) {
	var xacts []Transaction
	for {
		resp, err := c.ListTransactions(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list transactions: %w", err)
		}

		xacts = append(xacts, resp.Data...)

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

	return xacts, nil
}

func (c *Client) ListTransactions(ctx context.Context, params ListTransactionsParams) (*Response[[]Transaction], error) {
	path, err := c.buildURL("transactions")
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q, err := qparam.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	req.URL.RawQuery = q.Encode()

	body, err := c.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	var resp Response[[]Transaction]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

func (c *Client) PaginateAllTransactionsForAccount(ctx context.Context, accountID string, params ListTransactionsParams) ([]Transaction, error) {
	var xacts []Transaction
	for {
		resp, err := c.ListTransactionsForAccount(ctx, accountID, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list transactions: %w", err)
		}

		xacts = append(xacts, resp.Data...)

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

	return xacts, nil
}

func (c *Client) ListTransactionsForAccount(ctx context.Context, accoundID string, params ListTransactionsParams) (*Response[[]Transaction], error) {
	path, err := c.buildURL(fmt.Sprintf("accounts/%s/transactions", accoundID))
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q, err := qparam.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	req.URL.RawQuery = q.Encode()

	body, err := c.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	var resp Response[[]Transaction]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}
