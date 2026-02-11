package osc

import (
	"testing"
)

func TestOSCClient_CreationWithValidHostAndPort(t *testing.T) {
	c := NewClient("127.0.0.1", 57120)
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.client == nil {
		t.Fatal("expected non-nil internal client")
	}
}

func TestOSCClient_BoolToIntConversion(t *testing.T) {
	if boolToInt(true) != 1 {
		t.Error("expected true to be 1")
	}
	if boolToInt(false) != 0 {
		t.Error("expected false to be 0")
	}
}

func TestClientEffectsReordering(t *testing.T) {
	client := NewClient("127.0.0.1", 57120)

	order := []string{"filter", "granular", "delay"}
	err := client.SetEffectsOrder(order)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// GetEffectsOrder sends the request but response is handled asynchronously
	// through the server component. For this test, we verify the request
	// can be sent without error.
	receivedOrder, err := client.GetEffectsOrder()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify that GetEffectsOrder returns a valid order (even if it's default for now)
	if len(receivedOrder) == 0 {
		t.Error("expected non-empty effects order")
	}
}
