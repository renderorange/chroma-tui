package osc

import (
	"math"
	"testing"
)

func TestOSCValidation_InvalidFloat32Values(t *testing.T) {
	// Test boundary conditions
	testCases := []struct {
		name       string
		input      interface{}
		expected   float32
		shouldWarn bool
	}{
		{
			name:       "below minimum",
			input:      float32(-1000.1),
			expected:   0.0,
			shouldWarn: true,
		},
		{
			name:       "above maximum",
			input:      float32(1000.1),
			expected:   0.0,
			shouldWarn: true,
		},
		{
			name:       "valid minimum",
			input:      float32(-1000.0),
			expected:   -1000.0,
			shouldWarn: false,
		},
		{
			name:       "valid maximum",
			input:      float32(1000.0),
			expected:   1000.0,
			shouldWarn: false,
		},
		{
			name:       "negative infinity",
			input:      float32(math.Inf(-1)),
			expected:   0.0,
			shouldWarn: true,
		},
		{
			name:       "positive infinity",
			input:      float32(math.Inf(1)),
			expected:   0.0,
			shouldWarn: true,
		},
		{
			name:       "NaN",
			input:      float32(math.NaN()),
			expected:   0.0,
			shouldWarn: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toFloat32(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}

func TestOSCValidation_InvalidInt32Values(t *testing.T) {
	testCases := []struct {
		name       string
		input      interface{}
		expected   int32
		shouldWarn bool
	}{
		{
			name:       "below minimum",
			input:      int32(-2147483648),
			expected:   0,
			shouldWarn: true,
		},
		{
			name:       "above maximum",
			input:      int32(2147483647),
			expected:   0,
			shouldWarn: true,
		},
		{
			name:       "above maximum",
			input:      int32(2147483647),
			expected:   0,
			shouldWarn: true,
		},
		{
			name:       "above maximum",
			input:      int32(2147483647),
			expected:   0,
			shouldWarn: true,
		},
		{
			name:       "valid minimum",
			input:      int32(-2147483648),
			expected:   -2147483648,
			shouldWarn: false,
		},
		{
			name:       "valid maximum",
			input:      int32(2147483647),
			expected:   2147483647,
			shouldWarn: false,
		},
		{
			name:       "float64 too small",
			input:      float64(-2147483648.0),
			expected:   0,
			shouldWarn: true,
		},
		{
			name:       "float64 too large",
			input:      float64(2147483648.0),
			expected:   0,
			shouldWarn: true,
		},
		{
			name:       "float64 too large",
			input:      float64(2147483648.0),
			expected:   0,
			shouldWarn: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toInt(tc.input)
			if int32(result) != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestOSCValidation_InvalidStringValues(t *testing.T) {
	testCases := []struct {
		name       string
		input      interface{}
		expected   string
		shouldWarn bool
	}{
		{
			name:       "empty string",
			input:      "",
			expected:   "",
			shouldWarn: false,
		},
		{
			name:       "valid string",
			input:      "test",
			expected:   "test",
			shouldWarn: false,
		},
		{
			name:       "long string",
			input:      string(make([]byte, 300)), // 300 > 255 max
			expected:   "test",                    // truncated to 255 chars
			shouldWarn: true,
		},
		{
			name:       "unicode string",
			input:      "hello 世界",
			expected:   "hello 世界",
			shouldWarn: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toString(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

// TestStruct is used to test struct input validation
type TestStruct struct {
	Field int
}

func TestOSCValidation_InvalidTypes(t *testing.T) {
	testCases := []struct {
		name       string
		input      interface{}
		shouldWarn bool
	}{
		{
			name:       "nil input",
			input:      nil,
			shouldWarn: true,
		},
		{
			name:       "slice input",
			input:      []int{1, 2, 3},
			shouldWarn: true,
		},
		{
			name:       "map input",
			input:      map[string]int{"key": 1},
			shouldWarn: true,
		},
		{
			name:       "bool input",
			input:      true,
			shouldWarn: true,
		},
		{
			name:       "byte input",
			input:      byte(255),
			shouldWarn: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// These should trigger warning messages and return safe defaults
			result := toFloat32(tc.input)
			if result != 0.0 {
				t.Errorf("Expected safe default 0.0 for invalid type %T", tc.input)
			}
			resultInt := toInt(tc.input)
			if resultInt != 0 {
				t.Errorf("Expected safe default 0 for invalid type %T", tc.input)
			}
			resultStr := toString(tc.input)
			if resultStr != "subtle" {
				t.Errorf("Expected 'subtle' for invalid type %T", tc.input)
			}
		})
	}
}

func TestOSCValidation_BoundaryConditions(t *testing.T) {
	// Test edge cases around boundaries
	t.Run("float32 epsilon boundaries", func(t *testing.T) {
		// Test values just outside bounds
		below := toFloat32(float32(-1000.0 - 0.0001))
		above := toFloat32(float32(1000.0 + 0.0001))

		if below != 0.0 || above != 0.0 {
			t.Errorf("Expected clamping to 0.0 for out-of-bounds values")
		}

		// Test values just inside bounds
		validBelow := toFloat32(float32(-999.9))
		validAbove := toFloat32(float32(999.9))

		if validBelow != -999.9 || validAbove != 999.9 {
			t.Errorf("Expected valid values to pass through: got %f, %f", validBelow, validAbove)
		}
	})
}
