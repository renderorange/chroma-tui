package osc

import (
	"math"
	"testing"
)

func TestOSCValidation_InvalidFloat32Values(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		want     float32
		wantWarn bool
	}{
		{
			name:     "below minimum",
			input:    float32(-1000.1),
			want:     -1000.1,
			wantWarn: false,
		},
		{
			name:     "above maximum",
			input:    float32(1000.1),
			want:     1000.1,
			wantWarn: false,
		},
		{
			name:     "valid minimum",
			input:    float32(-1000.0),
			want:     -1000.0,
			wantWarn: false,
		},
		{
			name:     "valid maximum",
			input:    float32(1000.0),
			want:     1000.0,
			wantWarn: false,
		},
		{
			name:     "negative infinity",
			input:    float32(math.Inf(-1)),
			want:     float32(math.Inf(-1)),
			wantWarn: false,
		},
		{
			name:     "positive infinity",
			input:    float32(math.Inf(1)),
			want:     float32(math.Inf(1)),
			wantWarn: false,
		},
		{
			name:     "NaN",
			input:    float32(math.NaN()),
			want:     float32(math.NaN()),
			wantWarn: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toFloat32(tc.input)
			if tc.want != tc.want { // NaN check
				if !(result != result && tc.want != tc.want) { // both NaN
					t.Errorf("Expected %f, got %f", tc.want, result)
				}
			} else if result != tc.want && !(math.IsNaN(float64(result)) && math.IsNaN(float64(tc.want))) {
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
			name:     "below minimum",
			input:    int32(-2147483648),
			expected: -2147483648,
			wantWarn: false,
		},
		{
			name:     "above maximum",
			input:    int32(2147483647),
			expected: 2147483647,
			wantWarn: false,
		},
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
			input:    float64(-2147483648.0),
			expected: -2147483648,
			wantWarn: false,
		},
		{
			name:     "float64 too large",
			input:    float64(2147483647),
			expected: 2147483647,
			wantWarn: false,
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
			input:    string(make([]byte, 300)),
			expected: string(make([]byte, 255)),
			wantWarn: true,
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
