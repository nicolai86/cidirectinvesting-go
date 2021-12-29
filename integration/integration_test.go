package integration

import (
	"github.com/nicolai86/cidirectinvesting-go"
	"os"
	"testing"
	"time"
)

func requireConfiguredEnvironment(t *testing.T) {
	if os.Getenv("CDI_KEY_ID") == "" || os.Getenv("CDI_SECRET_KEY") == "" {
		t.Skip("environment not configured: missing CDI_KEY_ID or CDI_SECRET_KEY variable")
	}
}

func mustNewClient(t *testing.T) *cidirectinvesting.Client {
	client, err := cidirectinvesting.New(cidirectinvesting.WithThirdPartyKey(os.Getenv("CDI_KEY_ID"), os.Getenv("CDI_SECRET_KEY")))
	if err != nil {
		t.Fatalf("failed to create new client: %s", err)
	}
	return client
}

func mustNewClientWithLogin(t *testing.T) *cidirectinvesting.Client {
	client, err := cidirectinvesting.New(cidirectinvesting.WithThirdPartyKey(os.Getenv("CDI_KEY_ID"), os.Getenv("CDI_SECRET_KEY")))
	if err != nil {
		t.Fatalf("failed to create new client: %s", err)
	}
	if _, err := client.Login(); err != nil {
		t.Fatalf("failed to login: %s", err)
	}
	return client
}

func TestClient_Login(t *testing.T) {
	t.Parallel()
	requireConfiguredEnvironment(t)

	client := mustNewClient(t)
	session, err := client.Login()
	if err != nil {
		t.Fatalf("Login failed: %s", err)
	}

	if session.Email == "" {
		t.Fatalf("expected email to be set, but wasn't")
	}

	// nothing else is tested because it depends on the customer data
}

func TestClient_Dashboard(t *testing.T) {
	t.Parallel()
	requireConfiguredEnvironment(t)

	client := mustNewClientWithLogin(t)
	dashboard, err := client.Dashboard(cidirectinvesting.DashboardRequest{
		Currencies: map[cidirectinvesting.Currency]struct{}{
			cidirectinvesting.CurrencyCAD: {},
		},
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	})
	if err != nil {
		t.Fatalf("expected request to succeed, but didn't: %v", err)
	}

	if dashboard.InceptionDate == "" {
		t.Fatalf("expected inception date")
	}
	// no other attributes are tested because this depends on the actual client data
}

func TestClient_AccountDetails(t *testing.T) {
	t.Parallel()
	requireConfiguredEnvironment(t)

	client := mustNewClientWithLogin(t)
	dashboard, err := client.Dashboard(cidirectinvesting.DashboardRequest{
		Currencies: map[cidirectinvesting.Currency]struct{}{
			cidirectinvesting.CurrencyCAD: {},
		},
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	})
	if err != nil {
		t.Fatalf("expected dashboard request to succeed, but didn't: %v", err)
	}
	if len(dashboard.Accounts) == 0 {
		t.Skipf("Can't test without active accounts")
	}
	account := dashboard.Accounts[0]

	details, err := client.AccountDetails(cidirectinvesting.AccountDetailRequest{
		AccountID: account.Id,
	})
	if err != nil {
		t.Fatalf("expected account details request to succeed, but didn't: %v", err)
	}
	if details.Id != account.Id {
		t.Fatalf("inconsistent response: expected account %d, got %d", account.Id, details.Id)
	}
}

func TestClient_History(t *testing.T) {
	t.Parallel()
	requireConfiguredEnvironment(t)

	client := mustNewClientWithLogin(t)
	dashboard, err := client.Dashboard(cidirectinvesting.DashboardRequest{
		Currencies: map[cidirectinvesting.Currency]struct{}{
			cidirectinvesting.CurrencyCAD: {},
		},
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	})
	if err != nil {
		t.Fatalf("expected dashboard request to succeed, but didn't: %v", err)
	}

	if len(dashboard.Accounts) == 0 {
		t.Skipf("Can't test without active accounts")
	}
	account := dashboard.Accounts[0]
	history, err := client.History(cidirectinvesting.HistoryRequest{
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
		AccountID: &account.Id,
	})
	if err != nil {
		t.Fatalf("expected history request to succeed, but didn't: %v", err)
	}
	// no assertion on the history because it depends on the account
	_ = history
}

func TestClient_Transactions(t *testing.T) {
	t.Parallel()
	requireConfiguredEnvironment(t)

	client := mustNewClientWithLogin(t)
	dashboard, err := client.Dashboard(cidirectinvesting.DashboardRequest{
		Currencies: map[cidirectinvesting.Currency]struct{}{
			cidirectinvesting.CurrencyCAD: {},
		},
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	})
	if err != nil {
		t.Fatalf("expected dashboard request to succeed, but didn't: %v", err)
	}

	if len(dashboard.Accounts) == 0 {
		t.Skipf("Can't test without active accounts")
	}
	account := dashboard.Accounts[0]
	transactions, err := client.Transactions(cidirectinvesting.TransactionRequest{
		AccountID: account.Id,
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
		Page:      1,
		Limit:     5,
		Filter: &cidirectinvesting.Filter{
			Transactions: []string{
				"Deposit",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected history request to succeed, but didn't: %v", err)
	}
	if len(transactions.Data) < 1 {
		t.Fatalf("expected some transactions, but got none. This might be a false positive if you never transacted.")
	}
	// no assertion on the history because it depends on the account
	_ = transactions
}
