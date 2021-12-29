package main

import (
	"github.com/nicolai86/cidirectinvesting-go"
	"log"
	"os"
	"time"
)

func main() {
	client, err := cidirectinvesting.New(cidirectinvesting.WithThirdPartyKey(os.Getenv("CDI_KEY_ID"), os.Getenv("CDI_SECRET_KEY")))
	if err != nil {
		log.Fatalf("Client creation failed: %s\n", err)
	}

	if _, err := client.Login(); err != nil {
		log.Fatalf("Login failed: %s\n", err)
	}

	startDate := time.Now().AddDate(0, -2, 0)
	endDate := time.Now()
	dashboard, err := client.Dashboard(cidirectinvesting.DashboardRequest{
		Currencies: map[cidirectinvesting.Currency]struct{}{
			cidirectinvesting.CurrencyCAD: {},
		},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		log.Fatalf("Failed to load dashboard: %s\n", err)
	}

	historyAcrossAccounts, err := client.History(cidirectinvesting.HistoryRequest{
		StartDate: startDate,
		EndDate:   endDate,
	})
	_ = historyAcrossAccounts
	if err != nil {
		log.Fatalf("Unable to load account history: %s", err)
	}

	log.Println("Account list:")
	for _, account := range dashboard.Accounts {

		details, err := client.AccountDetails(cidirectinvesting.AccountDetailRequest{
			AccountID: account.Id,
		})
		if err != nil {
			log.Fatalf("Unable to load account details: %s", err)
		}
		log.Printf("%s(%d, %s, since %s): %2.2f %s\n", account.AccountName, account.Id, details.AccountNumber, account.InceptionDate, account.Balance, account.Currency)
		history, err := client.History(cidirectinvesting.HistoryRequest{
			StartDate: startDate,
			EndDate:   endDate,
			AccountID: &account.Id,
		})
		_ = history
		if err != nil {
			log.Fatalf("Unable to load account history: %s", err)
		}

		transactions, err := client.Transactions(cidirectinvesting.TransactionRequest{
			Limit:     10,
			Page:      1,
			StartDate: startDate,
			EndDate:   endDate,
			AccountID: account.Id,
			Filter: &cidirectinvesting.Filter{
				Quantity:     &cidirectinvesting.Range{},
				Price:        &cidirectinvesting.Range{},
				Value:        &cidirectinvesting.Range{},
				Transactions: []string{"Deposit"},
			},
		})
		if err != nil {
			log.Fatalf("Unable to load account transactions: %s", err)
		}
		log.Printf("%d deposits between %s and %s\n",
			transactions.TotalCount,
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"),
		)
	}
}
