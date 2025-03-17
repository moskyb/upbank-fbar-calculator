package ledger

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/moskyb/upbank-fbar-calculator/upapi"
)

type Ledger struct {
	CurrentBalance int
	AccountName    string
	Entries        []Entry
}

type Money int

func (m Money) MarshalCSV() (string, error) {
	return fmt.Sprintf("%.2f", float64(m)/100), nil
}

type Entry struct {
	ID string `json:"id" csv:"id"`

	CreatedAt time.Time  `json:"created_at" csv:"created_at"`
	SettledAt *time.Time `json:"settled_at" csv:"settled_at"`

	Description string  `json:"description" csv:"description"`
	Message     *string `json:"message" csv:"message"`

	Amount       Money `json:"amount" csv:"amount"`
	BalanceAfter Money `json:"balance_after" csv:"balance_after"`
}

func FromTransactions(accountName string, xacts []upapi.Transaction) *Ledger {
	ledger := &Ledger{AccountName: accountName}

	slices.Reverse(xacts)
	for _, xact := range xacts {
		amount := xact.Attributes.Amount.ValueInBaseUnits

		if xact.Attributes.RoundUp != nil {
			amount += xact.Attributes.RoundUp.Amount.ValueInBaseUnits
		}

		if xact.Attributes.Cashback != nil {
			amount += xact.Attributes.Cashback.Amount.ValueInBaseUnits
		}

		ledger.Entries = append(ledger.Entries, Entry{
			ID: xact.ID,

			CreatedAt: xact.Attributes.CreatedAt,
			SettledAt: xact.Attributes.SettledAt,

			Description: xact.Attributes.Description,
			Message:     xact.Attributes.Message,

			Amount:       Money(amount),
			BalanceAfter: Money(ledger.CurrentBalance + amount),
		})

		ledger.CurrentBalance += amount
	}

	return ledger
}

func (l *Ledger) HighWaterMark(year int) int {
	hwm := Money(0)

	if len(l.Entries) != 0 && len(l.TransactionsForYear(year)) == 0 {
		panic(fmt.Sprintf("MISSING FUNCTIONALITY: Attempting to calculate high watermark for account %s in year %d. %s has had transactions, but none this year. Please fix this, ben", l.AccountName, year, l.AccountName))
	}

	for _, entry := range l.Entries {
		if entry.CreatedAt.Year() == year && entry.BalanceAfter > hwm {
			hwm = entry.BalanceAfter
		}
	}

	return int(hwm)
}

func (l *Ledger) TransactionsForYear(year int) []Entry {
	var xacts []Entry
	for _, entry := range l.Entries {
		if entry.CreatedAt.Year() == year {
			xacts = append(xacts, entry)
		}
	}

	return xacts
}

func (l *Ledger) DumpCSV(year int) error {
	name := "./" + strings.ReplaceAll(l.AccountName, " ", "-") + ".csv"
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	err = gocsv.MarshalFile(l.TransactionsForYear(year), f)
	if err != nil {
		return fmt.Errorf("failed to marshal CSV: %w", err)
	}

	return nil
}
