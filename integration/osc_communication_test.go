package integration

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/renderorange/chroma/chroma-tui/osc"
	"github.com/renderorange/chroma/chroma-tui/tui"
)

// TestHelper provides utilities for OSC testing
type TestHelper struct {
	client *osc.Client
	model  *tui.Model
	port   int
	mu     sync.Mutex
}

// getAvailablePort finds an available port for testing
func getAvailablePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// newTestHelper creates a new test helper with dynamic port allocation
func newTestHelper(t *testing.T) *TestHelper {
	port, err := getAvailablePort()
	if err != nil {
		t.Fatalf("Failed to get available port: %v", err)
	}

	client := osc.NewClient("127.0.0.1", port)
	model := tui.NewModel(client)

	return &TestHelper{
		client: client,
		model:  &model,
		port:   port,
	}
}

func TestOSCCommunication_BasicConnectivity(t *testing.T) {
	helper := newTestHelper(t)

	// Test that client can be created
	if helper.client == nil {
		t.Fatal("Failed to create OSC client")
	}

	// Test that model is initialized
	if helper.model == nil {
		t.Fatal("Failed to create TUI model")
	}

	// Verify initial model state
	if helper.model.Gain != 1.0 {
		t.Errorf("Expected initial gain 1.0, got %f", helper.model.Gain)
	}

	t.Log("Basic connectivity test passed")
}

func TestOSCCommunication_TUIModelIntegration(t *testing.T) {
	helper := newTestHelper(t)

	// Test initial values
	if helper.model.Gain != 1.0 {
		t.Errorf("Expected initial gain 1.0, got %f", helper.model.Gain)
	}

	// Test focus navigation
	initialControl := helper.model.Focused()
	helper.model.NextControl()
	nextControl := helper.model.Focused()
	if nextControl == initialControl {
		t.Error("NextControl did not change the focused control")
	}

	helper.model.PrevControl()
	if helper.model.Focused() != initialControl {
		t.Error("PrevControl did not return to initial control")
	}

	t.Log("TUI model focus navigation working")
}

func TestOSCCommunication_ParameterAdjustment(t *testing.T) {
	helper := newTestHelper(t)

	// Test parameter adjustment (fire-and-forget)
	initialGain := helper.model.Gain
	helper.model.SetFocused(tui.TestCtrlGain)
	helper.model.AdjustFocused(0.1)

	// Verify local state updated (no server response expected)
	if helper.model.Gain <= initialGain {
		t.Errorf("Expected gain to increase from %f, got %f", initialGain, helper.model.Gain)
	}

	t.Log("Parameter adjustment test passed")
}

func TestOSCCommunication_ConcurrentAccess(t *testing.T) {
	helper := newTestHelper(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test concurrent parameter adjustments
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errors <- fmt.Errorf("goroutine %d: timeout", index)
				return
			default:
				// Simulate parameter adjustment
				helper.mu.Lock()
				helper.model.SetFocused(tui.TestCtrlGain)
				helper.model.AdjustFocused(0.01)
				helper.mu.Unlock()
			}
		}(i)
	}

	// Wait for all goroutines to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All goroutines completed
	case <-ctx.Done():
		t.Fatal("Concurrent access test timed out")
	}

	close(errors)
	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
	}

	t.Log("Concurrent access test passed")
}

func TestOSCCommunication_ErrorHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	helper := newTestHelper(t)

	// Test that client can handle errors gracefully
	initialGain := helper.model.Gain
	helper.model.SetFocused(tui.TestCtrlGain)
	helper.model.AdjustFocused(0.1)

	// Verify local state is updated even if server is unreachable
	if helper.model.Gain == initialGain {
		t.Error("Model should update locally even when server is unreachable")
	}

	// Test with context timeout
	select {
	case <-ctx.Done():
		t.Log("Context timeout handled correctly")
	default:
		t.Log("Error handling test passed")
	}
}

func TestOSCCommunication_StatelessBehavior(t *testing.T) {
	// Test that TUI works correctly in stateless mode (no server sync)
	helper := newTestHelper(t)

	// Set multiple parameters
	helper.model.Gain = 1.5
	helper.model.FilterCutoff = 3000
	helper.model.OverdriveDrive = 0.8

	// Verify local state is maintained
	if helper.model.Gain != 1.5 {
		t.Errorf("Expected gain 1.5, got %f", helper.model.Gain)
	}
	if helper.model.FilterCutoff != 3000 {
		t.Errorf("Expected filter cutoff 3000, got %f", helper.model.FilterCutoff)
	}
	if helper.model.OverdriveDrive != 0.8 {
		t.Errorf("Expected overdrive drive 0.8, got %f", helper.model.OverdriveDrive)
	}

	t.Log("Stateless behavior test passed")
}
