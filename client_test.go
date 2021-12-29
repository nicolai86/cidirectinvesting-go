package cidirectinvesting

import (
    "testing"
)

func TestWithThirdPartyKey(t *testing.T) {
    expectedKey, expectedSecret := "key", "secret"
    c, err := New(WithThirdPartyKey(expectedKey, expectedSecret))
    if err != nil {
        t.Fatalf("unexpected Client creation error: %s", err)
    }
    if c.apiKey != expectedKey {
        t.Fatalf("wrong API key: expected %q, got %q", expectedKey, c.apiKey)
    }
    if c.apiSecret != expectedSecret {
        t.Fatalf("wrong API key: expected %q, got %q", expectedSecret, c.apiSecret)
    }
}
