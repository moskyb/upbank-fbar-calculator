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

type Entry struct {
	ID string `json:"id" csv:"id"`

	CreatedAt time.Time  `json:"created_at" csv:"created_at"`
	SettledAt *time.Time `json:"settled_at" csv:"settled_at"`

	Description string  `json:"description" csv:"description"`
	Message     *string `json:"message" csv:"message"`

	Amount       int `json:"amount" csv:"amount"`
	BalanceAfter int `json:"balance_after" csv:"balance_after"`
}

func FromTransactions(accountName string, xacts []upapi.Transaction) *Ledger {
	ledger := &Ledger{AccountName: accountName}

	slices.Reverse(xacts)
	for _, xact := range xacts {
		if xact.Attributes.Cashback != nil {
			fmt.Println("cashback", xact.ID, xact.Attributes.Cashback.Amount.ValueInBaseUnits)
		}

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

			Amount:       amount,
			BalanceAfter: ledger.CurrentBalance + amount,
		})

		ledger.CurrentBalance += amount
	}

	return ledger
}

func (l *Ledger) HighWaterMark() int {
	hwm := 0
	for _, entry := range l.Entries {
		if entry.BalanceAfter > hwm {
			hwm = entry.BalanceAfter
		}
	}

	return hwm
}

func (l *Ledger) DumpCSV() error {
	name := "./" + strings.ReplaceAll(l.AccountName, " ", "-") + ".csv"
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	err = gocsv.MarshalFile(l.Entries, f)
	if err != nil {
		return fmt.Errorf("failed to marshal CSV: %w", err)
	}

	return nil
}
