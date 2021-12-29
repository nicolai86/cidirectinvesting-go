package cidirectinvesting

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_AccountDetails(t *testing.T) {
	t.Run("failure case", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client, err := New(WithAPIEndpoint(server.URL))
		if err != nil {
			t.Fatalf("client creation failed: %s", err)
		}
		_, err = client.AccountDetails(AccountDetailRequest{})
		if err == nil {
			t.Fatalf("account details request succeeded, but shouldn't")
		}
	})

	t.Run("successful case", func(t *testing.T) {
		expectedAccountName := "Test Account"
		expectedAccountID := 10
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, fmt.Sprintf(`{"id": %d, "account_name": %q}`, expectedAccountID, expectedAccountName))
		}))
		defer server.Close()

		client, err := New(WithAPIEndpoint(server.URL))
		if err != nil {
			t.Fatalf("client creation failed: %s", err)
		}

		response, err := client.AccountDetails(AccountDetailRequest{
			AccountID: expectedAccountID,
		})
		if err != nil {
			t.Fatalf("expected request to succeed, but didn't: %s", err)
		}
		if response.AccountName != expectedAccountName {
			t.Fatalf("wrong balance returned: expected %q, got %q", expectedAccountName, response.AccountName)
		}
		if response.Id != expectedAccountID {
			t.Fatalf("response contained wrong account id: expected %d, got %d", expectedAccountID, response.Id)
		}
	})
}
