package cidirectinvesting

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func TestClient_Dashboard(t *testing.T) {
    t.Run("failure case", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusBadRequest)
        }))
        defer server.Close()

        client, err := New(WithAPIEndpoint(server.URL))
        if err != nil {
            t.Fatalf("client creation failed: %s", err)
        }
        _, err = client.Dashboard(DashboardRequest{})
        if err == nil {
            t.Fatalf("dashboard request succeeded, but shouldn't")
        }
    })
    
    t.Run("successful case", func(t *testing.T) {
        expectedBalance := 2.2
        expectedAccountID := 10
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, fmt.Sprintf(`{"balance": %2.2f, "accounts": [{"id": %d}]}`, expectedBalance, expectedAccountID))
        }))
        defer server.Close()

        client, err := New(WithAPIEndpoint(server.URL))
        if err != nil {
            t.Fatalf("client creation failed: %s", err)
        }

        response, err := client.Dashboard(DashboardRequest{
            Currencies: map[Currency]struct{}{
                CurrencyCAD: {},
            },
            StartDate: time.Now().AddDate(0, -1, 0),
            EndDate:   time.Now(),
        })
        if err != nil {
            t.Fatalf("expected request to succeed, but didn't: %s", err)
        }
        if response.Balance != expectedBalance {
            t.Fatalf("wrong balance returned: expected %2.2f, got %2.2f", expectedBalance, response.Balance)
        }
        if len(response.Accounts) != 1 {
            t.Fatalf("expected 1 account but got %d", len(response.Accounts))
        }
        if response.Accounts[0].Id != expectedAccountID {
            t.Fatalf("response contained wrong account id: expected %d, got %d", expectedAccountID, response.Accounts[0].Id)
        }
    })
}
