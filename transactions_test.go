package cidirectinvesting

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Transactions(t *testing.T) {
	t.Run("failure case", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client, err := New(WithAPIEndpoint(server.URL))
		if err != nil {
			t.Fatalf("client creation failed: %s", err)
		}
		_, err = client.Transactions(TransactionRequest{
			AccountID: 10,
			StartDate: time.Now().AddDate(0, -1, 0),
			EndDate:   time.Now(),
		})
		if err == nil {
			t.Fatalf("transactions request succeeded, but shouldn't")
		}
	})

	t.Run("successful case", func(t *testing.T) {
		expectedTransactionID := 10
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, fmt.Sprintf(`{"data": [{"id": %d}]}`, expectedTransactionID))
		}))
		defer server.Close()

		client, err := New(WithAPIEndpoint(server.URL))
		if err != nil {
			t.Fatalf("client creation failed: %s", err)
		}

		response, err := client.Transactions(TransactionRequest{
			StartDate: time.Now().AddDate(0, -1, 0),
			EndDate:   time.Now(),
			AccountID: 10,
		})
		if err != nil {
			t.Fatalf("expected request to succeed, but didn't: %s", err)
		}
		if len(response.Data) != 1 {
			t.Fatalf("expected 1 transactions but got %d", len(response.Data))
		}
		if response.Data[0].Id != expectedTransactionID {
			t.Fatalf("response contained wrong account id: expected %d, got %d", expectedTransactionID, response.Data[0].Id)
		}
	})
}
