package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		operation string
		a         *float64
		b         *float64
		want      float64
		wantErr   error
	}{
		{name: "add positives", operation: OpAdd, a: ptr(19), b: ptr(9), want: 28},
		{name: "add negatives", operation: OpAdd, a: ptr(-3), b: ptr(-2), want: -5},
		{name: "add mixed signs", operation: OpAdd, a: ptr(-1.5), b: ptr(2), want: 0.5},
		{name: "add zeros", operation: OpAdd, a: ptr(0), b: ptr(0), want: 0},
		{name: "add decimals", operation: OpAdd, a: ptr(0.1), b: ptr(0.2), want: 0.3},

		{name: "subtract positives", operation: OpSubtract, a: ptr(5), b: ptr(10), want: -5},
		{name: "subtract negative operand", operation: OpSubtract, a: ptr(5), b: ptr(-3), want: 8},
		{name: "subtract zeros", operation: OpSubtract, a: ptr(0), b: ptr(0), want: 0},
		{name: "subtract decimals", operation: OpSubtract, a: ptr(1.5), b: ptr(0.25), want: 1.25},

		{name: "multiply positives", operation: OpMultiply, a: ptr(6), b: ptr(7), want: 42},
		{name: "multiply negatives", operation: OpMultiply, a: ptr(-2.5), b: ptr(4), want: -10},
		{name: "multiply by zero", operation: OpMultiply, a: ptr(12), b: ptr(0), want: 0},
		{name: "multiply decimals", operation: OpMultiply, a: ptr(0.5), b: ptr(0.5), want: 0.25},

		{name: "divide positives", operation: OpDivide, a: ptr(9), b: ptr(3), want: 3},
		{name: "divide negatives", operation: OpDivide, a: ptr(-9), b: ptr(3), want: -3},
		{name: "divide zero by nonzero", operation: OpDivide, a: ptr(0), b: ptr(5), want: 0},
		{name: "divide decimals", operation: OpDivide, a: ptr(1), b: ptr(4), want: 0.25},
		{name: "divide by zero", operation: OpDivide, a: ptr(1), b: ptr(0), wantErr: ErrDivisionByZero},
		{name: "divide by negative zero", operation: OpDivide, a: ptr(-5), b: ptr(math.Copysign(0, -1)), wantErr: ErrDivisionByZero},

		{name: "power positive exponent", operation: OpPower, a: ptr(2), b: ptr(10), want: 1024},
		{name: "power negative exponent", operation: OpPower, a: ptr(2), b: ptr(-1), want: 0.5},
		{name: "power negative base integer exponent", operation: OpPower, a: ptr(-2), b: ptr(3), want: -8},
		{name: "power zero exponent", operation: OpPower, a: ptr(9), b: ptr(0), want: 1},
		{name: "power zero to the zero", operation: OpPower, a: ptr(0), b: ptr(0), wantErr: ErrUndefinedPower},
		{name: "power decimals", operation: OpPower, a: ptr(9), b: ptr(0.5), want: 3},
		{name: "power negative base fractional exponent", operation: OpPower, a: ptr(-2), b: ptr(0.5), wantErr: ErrNonFiniteResult},

		{name: "sqrt positive", operation: OpSqrt, a: ptr(9), want: 3},
		{name: "sqrt zero", operation: OpSqrt, a: ptr(0), want: 0},
		{name: "sqrt decimal", operation: OpSqrt, a: ptr(0.25), want: 0.5},
		{name: "sqrt ignores unused operand", operation: OpSqrt, a: ptr(16), b: ptr(99), want: 4},
		{name: "sqrt negative", operation: OpSqrt, a: ptr(-1), wantErr: ErrNegativeSqrt},
		{name: "sqrt negative decimal", operation: OpSqrt, a: ptr(-0.01), wantErr: ErrNegativeSqrt},

		{name: "percentage example", operation: OpPercentage, a: ptr(200), b: ptr(15), want: 30},
		{name: "percentage of zero", operation: OpPercentage, a: ptr(50), b: ptr(0), want: 0},
		{name: "percentage one hundred", operation: OpPercentage, a: ptr(42), b: ptr(100), want: 42},
		{name: "percentage one", operation: OpPercentage, a: ptr(250), b: ptr(1), want: 2.5},
		{name: "percentage negative rate", operation: OpPercentage, a: ptr(80), b: ptr(-25), want: -20},
		{name: "percentage decimal rate", operation: OpPercentage, a: ptr(80), b: ptr(12.5), want: 10},
		{name: "percentage large base", operation: OpPercentage, a: ptr(math.MaxFloat64), b: ptr(50), want: math.MaxFloat64 / 2},
		{name: "percentage large rate", operation: OpPercentage, a: ptr(50), b: ptr(math.MaxFloat64), want: math.MaxFloat64 / 2},
		{name: "percentage large negative base", operation: OpPercentage, a: ptr(-1e308), b: ptr(15), want: -1.5e307},
		{name: "percentage overflow", operation: OpPercentage, a: ptr(math.MaxFloat64), b: ptr(math.MaxFloat64), wantErr: ErrNonFiniteResult},
		{name: "percentage overflow large operands", operation: OpPercentage, a: ptr(1e308), b: ptr(1e308), wantErr: ErrNonFiniteResult},

		{name: "missing operation", operation: "", a: ptr(1), b: ptr(2), wantErr: ErrMissingOperation},
		{name: "unsupported operation", operation: "modulo", a: ptr(1), b: ptr(2), wantErr: ErrUnsupportedOperation},
		{name: "uppercase operation", operation: "ADD", a: ptr(1), b: ptr(2), wantErr: ErrUnsupportedOperation},
		{name: "add missing a", operation: OpAdd, b: ptr(1), wantErr: ErrMissingOperandA},
		{name: "add missing b", operation: OpAdd, a: ptr(1), wantErr: ErrMissingOperandB},
		{name: "sqrt missing a", operation: OpSqrt, wantErr: ErrMissingOperandA},

		{name: "add overflow", operation: OpAdd, a: ptr(math.MaxFloat64), b: ptr(math.MaxFloat64), wantErr: ErrNonFiniteResult},
		{name: "subtract overflow", operation: OpSubtract, a: ptr(math.MaxFloat64), b: ptr(-math.MaxFloat64), wantErr: ErrNonFiniteResult},
		{name: "multiply overflow", operation: OpMultiply, a: ptr(math.MaxFloat64), b: ptr(2), wantErr: ErrNonFiniteResult},
		{name: "power overflow", operation: OpPower, a: ptr(10), b: ptr(400), wantErr: ErrNonFiniteResult},
		{name: "infinite operand", operation: OpAdd, a: ptr(math.Inf(1)), b: ptr(1), wantErr: ErrNonFiniteOperand},
		{name: "nan operand", operation: OpAdd, a: ptr(math.NaN()), b: ptr(1), wantErr: ErrNonFiniteOperand},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Calculate(tt.operation, tt.a, tt.b)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Calculate() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Calculate() unexpected error = %v", err)
			}
			if !almostEqual(got, tt.want) {
				t.Fatalf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptr(v float64) *float64 {
	return &v
}

func almostEqual(got, want float64) bool {
	const epsilon = 1e-12
	diff := math.Abs(got - want)
	scale := math.Max(1, math.Max(math.Abs(got), math.Abs(want)))
	return diff <= epsilon*scale
}
