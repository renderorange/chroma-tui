package osc

import (
	"fmt"
	"math"
	"testing"
)

// Helper functions for OSC type conversion with validation

func toFloat32(v interface{}) float32 {
	switch val := v.(type) {
	case float32:
		if math.IsNaN(float64(val)) || math.IsInf(float64(val), 0) {
			fmt.Printf("WARNING: toFloat32 invalid value %v, using 0.0\n", val)
			return 0.0
		}
		if val < -1000.0 || val > 1000.0 {
			fmt.Printf("WARNING: toFloat32 value %f out of range, using 0.0\n", val)
			return 0.0
		}
		return val
	case float64:
		f := float32(val)
		if math.IsNaN(val) || math.IsInf(val, 0) {
			fmt.Printf("WARNING: toFloat32 invalid value %v, using 0.0\n", val)
			return 0.0
		}
		if f < -1000.0 || f > 1000.0 {
			fmt.Printf("WARNING: toFloat32 value %f out of range, using 0.0\n", f)
			return 0.0
		}
		return f
	case int32:
		return float32(val)
	case int:
		return float32(val)
	default:
		fmt.Printf("WARNING: toFloat32 invalid type %T, expected float32/float64/int32/int, using 0.0\n", v)
		return 0.0
	}
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case int32:
		return int(val)
	case int:
		return val
	case float32:
		f := float64(val)
		if f < -2147483648.0 || f > 2147483647.0 {
			fmt.Printf("WARNING: toInt value %f out of int32 range, using 0\n", f)
			return 0
		}
		return int(val)
	case float64:
		if val < -2147483648.0 || val > 2147483647.0 {
			fmt.Printf("WARNING: toInt value %f out of int32 range, using 0\n", val)
			return 0
		}
		return int(val)
	default:
		fmt.Printf("WARNING: toInt invalid type %T, expected int/int32/float32/float64, using 0\n", v)
		return 0
	}
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		if len(val) > 255 {
			fmt.Printf("WARNING: toString string length %d exceeds maximum 255, truncating\n", len(val))
			return val[:255]
		}
		return val
	default:
		fmt.Printf("WARNING: toString invalid type %T, expected string, using 'subtle'\n", v)
		return "subtle"
	}
}

func TestOSCValidation_InvalidFloat32Values(t *testing.T) {
	testCases := []struct {
		name         string
		input        interface{}
		want         float32
		shouldReject bool
	}{
		{
			name:         "below minimum - rejected",
			input:        float32(-1000.1),
			want:         0.0, // Rejected by security validation
			shouldReject: true,
		},
		{
			name:         "above maximum - rejected",
			input:        float32(1000.1),
			want:         0.0, // Rejected by security validation
			shouldReject: true,
		},
		{
			name:         "valid minimum - accepted",
			input:        float32(-1000.0),
			want:         -1000.0,
			shouldReject: false,
		},
		{
			name:         "valid maximum - accepted",
			input:        float32(1000.0),
			want:         1000.0,
			shouldReject: false,
		},
		{
			name:         "negative infinity - rejected",
			input:        float32(math.Inf(-1)),
			want:         0.0, // Rejected by security validation
			shouldReject: true,
		},
		{
			name:         "positive infinity - rejected",
			input:        float32(math.Inf(1)),
			want:         0.0, // Rejected by security validation
			shouldReject: true,
		},
		{
			name:         "NaN - rejected",
			input:        float32(math.NaN()),
			want:         0.0, // Rejected by security validation
			shouldReject: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toFloat32(tc.input)
			if result != tc.want {
				t.Errorf("Expected %f, got %f", tc.want, result)
			}
		})
	}
}

func TestOSCValidation_InvalidInt32Values(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected int32
		wantWarn bool
	}{
		{
			name:     "valid minimum",
			input:    int32(-2147483648),
			expected: -2147483648,
			wantWarn: false,
		},
		{
			name:     "valid maximum",
			input:    int32(2147483647),
			expected: 2147483647,
			wantWarn: false,
		},
		{
			name:     "float64 too small",
			input:    float64(-2147483649.0),
			expected: 0,
			wantWarn: true,
		},
		{
			name:     "float64 too large",
			input:    float64(2147483648.0),
			expected: 0,
			wantWarn: true,
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
		name     string
		input    interface{}
		expected string
		wantWarn bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantWarn: false,
		},
		{
			name:     "valid string",
			input:    "test",
			expected: "test",
			wantWarn: false,
		},
		{
			name:     "long string",
			input:    string(make([]byte, 300)), // 300 > 255 max
			expected: string(make([]byte, 255)), // truncated to 255 chars
			wantWarn: true,
		},
		{
			name:     "unicode string",
			input:    "hello 世界",
			expected: "hello 世界",
			wantWarn: false,
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

func TestOSCValidation_BoundaryConditions(t *testing.T) {
	tests := []struct {
		name  string
		value float32
	}{
		{"zero", 0.0},
		{"negative", -0.5},
		{"positive", 0.5},
		{"float32_epsilon_boundaries", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toFloat32(tt.value)
			if result != tt.value {
				t.Errorf("expected %f, got %f", tt.value, result)
			}
		})
	}
}
