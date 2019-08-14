package monzo

import (
	"fmt"
	"testing"
)

func TestPing(t *testing.T) {
	a1 := NewClient("")
	err := a1.Ping()
	if err == nil {
		fmt.Println("expected an error when using empty string for token. Didn't get one")
		t.FailNow()
	}

	a2 := NewClient("invalid_token")
	err = a2.Ping()
	if err == nil {
		fmt.Print("expected an error when using an invalid token. Didn't get one")
		t.FailNow()
	}
}
