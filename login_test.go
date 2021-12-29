package cidirectinvesting

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestClient_Login(t *testing.T) {
    t.Run("failure case", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusForbidden)
        }))
        defer server.Close()

        client, err := New(WithAPIEndpoint(server.URL))
        if err != nil {
            t.Fatalf("client creation failed: %s", err)
        }
        _, err = client.Login()
        if err == nil {
            t.Fatalf("login succeeded, but shouldn't")
        }
    })

    t.Run("successful case", func(t *testing.T) {
        expectedEmail := "test@example.com"
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, fmt.Sprintf(`{"email": %q}`, expectedEmail))
        }))
        defer server.Close()

        client, err := New(WithAPIEndpoint(server.URL))
        if err != nil {
            t.Fatalf("client creation failed: %s", err)
        }
        session, err := client.Login()
        if err != nil {
            t.Fatalf("login failed: %s", err)
        }
        if session.Email != expectedEmail {
            t.Fatalf("wrong email address: expected %q, got %q", expectedEmail, session.Email)
        }
    })

}
