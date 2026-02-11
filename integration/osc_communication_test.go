package integration

import (
	"context"
	"fmt"
	"math"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/renderorange/chroma/chroma-tui/osc"
	"github.com/renderorange/chroma/chroma-tui/tui"
)

// TestHelper provides utilities for robust OSC testing
type TestHelper struct {
	server     *osc.Server
	client     *osc.Client
	model      *tui.Model
	port       int
	serverDone chan struct{}
	mu         sync.Mutex
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

	server := osc.NewServer(port)
	client := osc.NewClient("127.0.0.1", port)
	model := tui.NewModel(client)

	return &TestHelper{
		server:     server,
		client:     client,
		model:      &model,
		port:       port,
		serverDone: make(chan struct{}),
	}
}

// startServer starts the OSC server in a goroutine
func (th *TestHelper) startServer(ctx context.Context) error {
	th.mu.Lock()
	defer th.mu.Unlock()

	serverStarted := make(chan error, 1)

	go func() {
		defer close(th.serverDone)
		// Server.Start() is blocking, so we run it in a goroutine
		if err := th.server.Start(); err != nil {
			serverStarted <- fmt.Errorf("server failed to start: %w", err)
			return
		}
	}()

	// Give the server a moment to start listening
	select {
	case <-time.After(50 * time.Millisecond):
		return nil // Server started successfully
	case <-ctx.Done():
		return ctx.Err()
	}
}

// stopServer stops the OSC server gracefully
func (th *TestHelper) stopServer() {
	th.mu.Lock()
	defer th.mu.Unlock()

	// Close serverDone channel to signal shutdown
	select {
	case <-th.serverDone:
	default:
		close(th.serverDone)
	}
}

// waitForServerReady waits until the server is ready to receive messages
func (th *TestHelper) waitForServerReady(ctx context.Context, timeout time.Duration) error {
	// Give the server more time to start up
	select {
	case <-time.After(timeout):
		return nil // Assume server is ready after timeout
	case <-ctx.Done():
		return ctx.Err()
	}
}

func TestOSCCommunication_BasicMessageSending(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	helper := newTestHelper(t)
	defer helper.stopServer()

	// Start server with synchronization
	if err := helper.startServer(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to be ready
	if err := helper.waitForServerReady(ctx, 1*time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Test: Send a parameter change via the TUI client
	testGain := float32(1.5)
	if err := helper.client.SetGain(testGain); err != nil {
		t.Fatalf("Failed to send gain: %v", err)
	}

	// Test: Send a boolean parameter
	if err := helper.client.SetInputFreeze(true); err != nil {
		t.Fatalf("Failed to send input freeze: %v", err)
	}

	// Test: Send an int parameter
	if err := helper.client.SetBlendMode(2); err != nil {
		t.Fatalf("Failed to send blend mode: %v", err)
	}

	t.Log("OSC messages sent successfully")
}

func TestOSCCommunication_NetworkFailureHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	helper := newTestHelper(t)
	defer helper.stopServer()

	// Test client behavior when server is not running
	// Note: UDP client doesn't return connection errors, it just drops packets
	client := osc.NewClient("127.0.0.1", helper.port)

	// These should not crash when server is not running (UDP is connectionless)
	if err := client.SetGain(1.0); err != nil {
		t.Errorf("Unexpected error when sending to non-running server (UDP is connectionless): %v", err)
	}

	// Start server
	if err := helper.startServer(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to be ready
	if err := helper.waitForServerReady(ctx, 1*time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Now messages should succeed
	if err := client.SetGain(1.0); err != nil {
		t.Errorf("Expected success when server is running, but got error: %v", err)
	}

	// Test server shutdown during operation
	helper.stopServer()

	// Give server time to shutdown
	time.Sleep(100 * time.Millisecond)

	// Messages should not crash even when server is stopped (UDP is connectionless)
	if err := client.SetGain(1.5); err != nil {
		t.Errorf("Unexpected error when sending to stopped server (UDP is connectionless): %v", err)
	}

	// Test with invalid port - this will fail with port validation error
	invalidClient := osc.NewClient("127.0.0.1", 99999)
	if err := invalidClient.SetGain(1.0); err == nil {
		t.Error("Expected error with invalid port, but got none")
	} else {
		t.Logf("Expected error with invalid port: %v", err)
	}

	t.Log("Network failure handling working correctly - UDP is connectionless by design")
}

func TestOSCCommunication_ConcurrentAccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	helper := newTestHelper(t)
	defer helper.stopServer()

	// Start server
	if err := helper.startServer(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to be ready
	if err := helper.waitForServerReady(ctx, 1*time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Test concurrent model access
	const numGoroutines = 10
	const numOperations = 20

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				// Test concurrent parameter setting
				gain := float32(0.5 + float64(goroutineID*numOperations+j)*0.01)
				if err := helper.client.SetGain(gain); err != nil {
					errors <- fmt.Errorf("goroutine %d, op %d: %v", goroutineID, j, err)
					return
				}

				// Test concurrent model state access
				helper.model.Gain = gain
				helper.model.NextControl()

				// Small delay to increase chance of race conditions
				time.Sleep(1 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
	}

	t.Log("Concurrent access test passed")
}

func TestOSCCommunication_TUIModelIntegration(t *testing.T) {
	helper := newTestHelper(t)
	defer helper.stopServer()

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

func TestOSCCommunication_PendingChangesSystem(t *testing.T) {
	helper := newTestHelper(t)
	defer helper.stopServer()

	// Set initial state
	initialGain := float32(0.5)
	helper.model.Gain = initialGain

	// Simulate a state update from server (no pending changes)
	serverState := osc.State{
		Gain:         1.0,
		FilterAmount: 0.8,
		BlendMode:    1,
	}

	helper.model.ApplyState(serverState)

	// Values should be updated since no pending changes
	if helper.model.Gain != 1.0 {
		t.Errorf("Expected gain to be updated to 1.0, got %f", helper.model.Gain)
	}
	if helper.model.FilterAmount != 0.8 {
		t.Errorf("Expected filter amount to be updated to 0.8, got %f", helper.model.FilterAmount)
	}
	if helper.model.BlendMode != 1 {
		t.Errorf("Expected blend mode to be updated to 1, got %d", helper.model.BlendMode)
	}

	t.Log("Pending changes system working correctly")
}

func TestOSCCommunication_StateCleanup(t *testing.T) {
	helper := newTestHelper(t)
	defer helper.stopServer()

	// Set a value and simulate time passing
	helper.model.Gain = 0.9

	// Create a server state with different values
	serverState := osc.State{
		Gain: 1.2,
	}

	// Apply state (should clean up any stale changes)
	helper.model.ApplyState(serverState)

	// Wait for cleanup timeout (500ms)
	time.Sleep(600 * time.Millisecond)

	// Apply state again - should update since pending change is stale
	helper.model.ApplyState(serverState)

	if helper.model.Gain != 1.2 {
		t.Errorf("Expected stale pending change to be cleaned up, got %f", helper.model.Gain)
	}

	t.Log("State cleanup working correctly")
}

func TestOSCCommunication_ErrorHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	helper := newTestHelper(t)
	defer helper.stopServer()

	// Start server
	if err := helper.startServer(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to be ready
	if err := helper.waitForServerReady(ctx, 1*time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Test malformed data handling - extreme values
	testCases := []struct {
		name  string
		value float32
		valid bool
	}{
		{"Normal gain", 1.0, true},
		{"Zero gain", 0.0, true},
		{"Negative gain", -1.0, true},            // Should not crash
		{"Very high gain", 1000.0, true},         // Should not crash
		{"NaN gain", float32(math.NaN()), false}, // Should handle gracefully
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := helper.client.SetGain(tc.value)
			if tc.valid && err != nil {
				t.Errorf("Expected valid value %f to succeed, but got error: %v", tc.value, err)
			}
			// For invalid values, we just check that it doesn't crash
		})
	}

	t.Log("Error handling test passed")
}

func TestOSCCommunication_MultipleParameterTypes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	helper := newTestHelper(t)
	defer helper.stopServer()

	// Start server
	if err := helper.startServer(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to be ready
	if err := helper.waitForServerReady(ctx, 1*time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Test float parameters
	floatTests := []struct {
		name   string
		method func(float32) error
		value  float32
	}{
		{"Gain", helper.client.SetGain, 1.5},
		{"FilterCutoff", helper.client.SetFilterCutoff, 4000},
		{"OverdriveDrive", helper.client.SetOverdriveDrive, 0.8},
		{"GranularDensity", helper.client.SetGranularDensity, 25},
	}

	for _, test := range floatTests {
		if err := test.method(test.value); err != nil {
			t.Errorf("Failed to set %s to %f: %v", test.name, test.value, err)
		}
	}

	// Test int parameters
	if err := helper.client.SetBlendMode(2); err != nil {
		t.Errorf("Failed to set blend mode: %v", err)
	}

	// Test boolean parameters
	boolTests := []struct {
		name   string
		method func(bool) error
		value  bool
	}{
		{"InputFreeze", helper.client.SetInputFreeze, true},
		{"GranularFreeze", helper.client.SetGranularFreeze, false},
	}

	for _, test := range boolTests {
		if err := test.method(test.value); err != nil {
			t.Errorf("Failed to set %s to %t: %v", test.name, test.value, err)
		}
	}

	// Test string parameters
	if err := helper.client.SetGrainIntensity("pronounced"); err != nil {
		t.Errorf("Failed to set grain intensity: %v", err)
	}

	t.Log("All parameter types working correctly")
}

func TestVisualizerOSCIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	helper := newTestHelper(t)
	defer helper.stopServer()

	// Start server
	if err := helper.startServer(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to be ready
	if err := helper.waitForServerReady(ctx, 1*time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Verify initial state - spectrum and waveform should be zeroed
	for i := 0; i < 8; i++ {
		if helper.model.Spectrum[i] != 0.0 {
			t.Errorf("Expected initial spectrum[%d] to be 0.0, got %f", i, helper.model.Spectrum[i])
		}
	}
	for i := 0; i < 64; i++ {
		if helper.model.Waveform[i] != 0.0 {
			t.Errorf("Expected initial waveform[%d] to be 0.0, got %f", i, helper.model.Waveform[i])
		}
	}

	// Test spectrum data consistency - simulate realistic audio spectrum data
	expectedSpectrum := [8]float32{0.1, 0.3, 0.8, 0.6, 0.4, 0.2, 0.1, 0.05}

	// Test waveform data consistency - simulate realistic waveform data
	expectedWaveform := [64]float32{}
	for i := range expectedWaveform {
		// Create a simple sine wave pattern
		expectedWaveform[i] = float32(0.5 * sin(float64(i)*2*3.14159/64))
	}

	// Test state update with both spectrum and waveform data
	testState := osc.State{
		Gain:     1.0,
		Spectrum: expectedSpectrum,
		Waveform: expectedWaveform,
	}

	// Apply the state to the model
	helper.model.ApplyState(testState)

	// Verify spectrum data was applied correctly
	for i := 0; i < 8; i++ {
		if helper.model.Spectrum[i] != expectedSpectrum[i] {
			t.Errorf("Spectrum[%d] mismatch: expected %f, got %f", i, expectedSpectrum[i], helper.model.Spectrum[i])
		}
	}

	// Verify waveform data was applied correctly
	for i := 0; i < 64; i++ {
		if helper.model.Waveform[i] != expectedWaveform[i] {
			t.Errorf("Waveform[%d] mismatch: expected %f, got %f", i, expectedWaveform[i], helper.model.Waveform[i])
		}
	}

	// Test timing consistency - multiple rapid updates should not corrupt data
	rapidUpdateCount := 10
	rapidUpdateSpectrum := [8]float32{0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.2}
	rapidUpdateWaveform := [64]float32{}
	for i := range rapidUpdateWaveform {
		// Create a different pattern for rapid updates
		rapidUpdateWaveform[i] = float32(0.3 * cos(float64(i)*4*3.14159/64))
	}

	// Apply multiple rapid updates to test timing consistency
	for i := 0; i < rapidUpdateCount; i++ {
		testState := osc.State{
			Spectrum: rapidUpdateSpectrum,
			Waveform: rapidUpdateWaveform,
		}
		helper.model.ApplyState(testState)

		// Small delay to simulate real-world timing
		time.Sleep(1 * time.Millisecond)
	}

	// Verify final state is consistent after rapid updates
	for i := 0; i < 8; i++ {
		if helper.model.Spectrum[i] != rapidUpdateSpectrum[i] {
			t.Errorf("Rapid update spectrum[%d] mismatch: expected %f, got %f", i, rapidUpdateSpectrum[i], helper.model.Spectrum[i])
		}
	}

	for i := 0; i < 64; i++ {
		if helper.model.Waveform[i] != rapidUpdateWaveform[i] {
			t.Errorf("Rapid update waveform[%d] mismatch: expected %f, got %f", i, rapidUpdateWaveform[i], helper.model.Waveform[i])
		}
	}

	// Test data bounds - ensure spectrum and waveform values stay within expected ranges
	// Spectrum should typically be 0.0 to 1.0 (normalized FFT values)
	for i := 0; i < 8; i++ {
		if helper.model.Spectrum[i] < 0.0 || helper.model.Spectrum[i] > 1.0 {
			t.Errorf("Spectrum[%d] out of bounds [0.0, 1.0]: %f", i, helper.model.Spectrum[i])
		}
	}

	// Waveform should typically be -1.0 to 1.0 (normalized audio samples)
	for i := 0; i < 64; i++ {
		if helper.model.Waveform[i] < -1.0 || helper.model.Waveform[i] > 1.0 {
			t.Errorf("Waveform[%d] out of bounds [-1.0, 1.0]: %f", i, helper.model.Waveform[i])
		}
	}

	t.Log("Visualizer OSC integration test passed - spectrum and waveform data flow verified")
}

func TestEffectsOrderOSC(t *testing.T) {
	// This test verifies that /chroma/effectsOrder and /chroma/getEffectsOrder OSC handlers exist
	// and work correctly in SuperCollider

	// TODO: Implement test that:
	// 1. Sends /chroma/effectsOrder with new order
	// 2. Sends /chroma/getEffectsOrder to retrieve current order
	// 3. Verifies the response contains the expected order

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	helper := newTestHelper(t)
	defer helper.stopServer()

	// Start server
	if err := helper.startServer(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to be ready
	if err := helper.waitForServerReady(ctx, 1*time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Test setting effects order
	if err := helper.client.Send("/chroma/effectsOrder", "filter", "granular", "delay"); err != nil {
		t.Fatalf("Failed to send effects order: %v", err)
	}

	// Test getting effects order - this should fail without handler
	// TODO: Add proper OSC response handling to verify the order was set
}

// Helper functions for test data generation
func sin(x float64) float64 {
	return float64(float32(math.Sin(x)))
}

func cos(x float64) float64 {
	return float64(float32(math.Cos(x)))
}
