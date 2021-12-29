package cidirectinvesting

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_History(t *testing.T) {
	t.Run("failure case", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client, err := New(WithAPIEndpoint(server.URL))
		if err != nil {
			t.Fatalf("client creation failed: %s", err)
		}
		_, err = client.History(HistoryRequest{})
		if err == nil {
			t.Fatalf("history request succeeded, but shouldn't")
		}
	})

	t.Run("successful case specific account", func(t *testing.T) {
		accountID := 7
		valueAOne, valueATwo := 100.25, 250.10
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, fmt.Sprintf(`[
    {"date":"2021-12-26","%d":%.2f,"total":350.35},
    {"date":"2021-12-27","%d":%.2f,"total":350.60}
]`, accountID, valueAOne, accountID, valueATwo))
		}))
		defer server.Close()
		client, err := New(WithAPIEndpoint(server.URL))
		if err != nil {
			t.Fatalf("client creation failed: %s", err)
		}

		response, err := client.History(HistoryRequest{})
		if err != nil {
			t.Fatalf("expected request to succeed, but didn't: %s", err)
		}
		if len(*response) != 2 {
			t.Fatalf("expected 2 data points, got %d", len(*response))
		}
		first := (*response)[0]
		if first.Accounts[accountID] != valueAOne {
			t.Fatalf("expected %.2f, got %2.2f", valueAOne, first.Accounts[accountID])
		}
		second := (*response)[1]
		if second.Accounts[accountID] != valueATwo {
			t.Fatalf("expected %.2f, got %2.2f", valueATwo, second.Accounts[accountID])
		}
	})

	t.Run("successful case all accounts", func(t *testing.T) {
		accountIDA, accountIDB := 7, 42
		valueAOne, valueBOne := 100.25, 250.10
		valueATwo, valueBTwo := 100.50, 250.10
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, fmt.Sprintf(`[
    {"date":"2021-12-26","%d":%.2f,"%d":%.2f,"total":350.35},
    {"date":"2021-12-27","%d":%.2f,"%d":%.2f,"total":350.60}
]`, accountIDA, valueAOne, accountIDB, valueBOne, accountIDA, valueATwo, accountIDB, valueBTwo))
		}))
		defer server.Close()

		client, err := New(WithAPIEndpoint(server.URL))
		if err != nil {
			t.Fatalf("client creation failed: %s", err)
		}

		response, err := client.History(HistoryRequest{})
		if err != nil {
			t.Fatalf("expected request to succeed, but didn't: %s", err)
		}
		if len(*response) != 2 {
			t.Fatalf("expected 2 data points, got %d", len(*response))
		}
		first := (*response)[0]
		if first.Accounts[accountIDA] != valueAOne {
			t.Fatalf("expected %.2f, got %2.2f", valueAOne, first.Accounts[accountIDA])
		}
		if first.Accounts[accountIDB] != valueBOne {
			t.Fatalf("expected %.2f, got %2.2f", valueBOne, first.Accounts[accountIDB])
		}
		second := (*response)[1]
		if second.Accounts[accountIDA] != valueATwo {
			t.Fatalf("expected %.2f, got %2.2f", valueATwo, second.Accounts[accountIDA])
		}
		if second.Accounts[accountIDB] != valueBTwo {
			t.Fatalf("expected %.2f, got %2.2f", valueBTwo, second.Accounts[accountIDB])
		}
	})
}
