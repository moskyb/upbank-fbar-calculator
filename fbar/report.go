package fbar

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/moskyb/upbank-fbar-calculator/ledger"
	"github.com/moskyb/upbank-fbar-calculator/upapi"
)

type Report struct {
	FinancialYear int
	Entries       map[string]ReportEntry
}

type AccountRecord struct {
	DisplayName    string `json:"display_name"`
	AccountType    string `json:"account_type"`
	Ownership      string `json:"ownership"`
	ClosingBalance int    `json:"closing_balance"`
	HighWaterMark  int    `json:"high_water_mark"`
}

func GenerateReport(upAPIToken string, year int) (*Report, error) {
	zone, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %w", err)
	}

	client := upapi.NewClient(upAPIToken, upapi.WithQuiet())
	accounts, err := client.PaginateAllAccounts(context.Background(), upapi.ListAccountsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	errsMtx := sync.Mutex{}
	errs := make([]error, 0, len(accounts))

	r := &Report{FinancialYear: year}
	r.Entries = make(map[string]ReportEntry, len(accounts))

	wg := sync.WaitGroup{}
	for _, acc := range accounts {
		wg.Add(1)
		go func(acc upapi.Account) {
			defer wg.Done()

			if acc.Attributes.CreatedAt.After(time.Date(year+1, time.January, 1, 0, 0, 0, 0, zone)) {
				// Account created after the end of the calendar year we're looking at, so it doesn't need to be reported on
				return
			}

			xacts, err := client.PaginateAllTransactionsForAccount(context.Background(), acc.ID, upapi.ListTransactionsParams{
				Until: time.Date(year+1, time.January, 1, 0, 0, 0, 0, zone),
			})
			if err != nil {
				errsMtx.Lock()
				errs = append(errs, fmt.Errorf("failed to list transactions for account %s: %w", acc.ID, err))
				errsMtx.Unlock()

				return
			}

			ledger := ledger.FromTransactions(acc.Attributes.DisplayName, xacts)
			r.Entries[acc.Attributes.DisplayName] = ReportEntry{
				AccountName:      acc.Attributes.DisplayName,
				HighWaterMark:    ledger.HighWaterMark(year),
				ClosingBalance:   ledger.CurrentBalance,
				TransactionCount: len(ledger.TransactionsForYear(year)),
			}

			err = ledger.DumpCSV(year)
			if err != nil {
				errsMtx.Lock()
				errs = append(errs, fmt.Errorf("failed to dump CSV for account %s: %w", acc.ID, err))
				errsMtx.Unlock()

				return
			}

		}(acc)
	}
	wg.Wait()

	return r, errors.Join(errs...)
}

func (r *Report) PrettyString() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("FBAR Report for Upbank, CY%d\n\n", r.FinancialYear))
	sb.WriteString(fmt.Sprintf("%d accounts held in %d:\n", len(r.Entries), r.FinancialYear))

	sortedEntries := slices.SortedFunc(maps.Values(r.Entries), func(i, j ReportEntry) int {
		return strings.Compare(stripEmoji(i.AccountName), stripEmoji(j.AccountName))
	})

	for _, entry := range sortedEntries {
		sb.WriteString(fmt.Sprintf("\t%s\n", entry.AccountName))
	}
	sb.WriteString("\n")

	for _, entry := range sortedEntries {
		sb.WriteString(fmt.Sprintf("Account: %s\n", entry.AccountName))
		sb.WriteString(fmt.Sprintf("\tTransaction count: %d\n", entry.TransactionCount))
		sb.WriteString(fmt.Sprintf("\tHigh water mark: %s\n", PrettyMoney(entry.HighWaterMark)))
		sb.WriteString(fmt.Sprintf("\tClosing balance: %s\n", PrettyMoney(entry.ClosingBalance)))
		sb.WriteString("\n")
	}

	return sb.String()
}

var emojiRE = regexp.MustCompile(`[[:^ascii:]]`)

func stripEmoji(s string) string {
	return strings.TrimSpace(emojiRE.ReplaceAllString(s, ""))
}

func PrettyMoney(amount int) string {
	return fmt.Sprintf("AUD $%d.%02d", amount/100.0, amount%100.0)
}

type ReportEntry struct {
	AccountName      string
	TransactionCount int
	HighWaterMark    int
	OpeningBalance   int
	ClosingBalance   int
}
