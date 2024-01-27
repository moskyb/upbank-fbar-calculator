package upapi

import (
	"fmt"
)

const (
	AccountTypeSaver         = "SAVER"
	AccountTypeTransactional = "TRANSACTIONAL"
	AccountTypeHomeLoan      = "HOME_LOAN"
)

const (
	OwnershipTypeIndividual = "INDIVIDUAL"
	OwnershipTypeJoint      = "JOINT"
)

type Money struct {
	CurrencyCode     string `json:"currencyCode"`
	Value            string `json:"value"`
	ValueInBaseUnits int    `json:"valueInBaseUnits"`
}

type Response[T any] struct {
	Data  T `json:"data"`
	Links struct {
		Prev *string `json:"prev"`
		Next *string `json:"next"`
	} `json:"links"`
}

type ErrorResponse struct {
	Errors []struct {
		Status string `json:"status"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Source *struct {
			Parameter string `json:"parameter"`
			Pointer   string `json:"pointer"`
		} `json:"source,omitempty"`
	} `json:"errors"`
}

func (e *ErrorResponse) Error() string {
	eb := "Up API returned error: "
	if len(e.Errors) > 1 {
		eb += "Up API Returned Errors:\n"
	}

	for _, err := range e.Errors {
		eb += fmt.Sprintf("	%s: %s - %s", err.Status, err.Title, err.Detail)
		if err.Source != nil {
			if err.Source.Parameter != "" {
				eb += fmt.Sprintf(" (Parameter: %s)", err.Source.Parameter)
			}
			if err.Source.Pointer != "" {
				eb += fmt.Sprintf(" (Pointer: %s)", err.Source.Pointer)
			}
		}
		eb += "\n"
	}

	return eb
}
