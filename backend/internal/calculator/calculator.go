// Package calculator evaluates traditional calculator operations.
// It has no dependency on HTTP or JSON.
package calculator

import (
	"errors"
	"fmt"
	"math"
)

const (
	OpAdd        = "add"
	OpSubtract   = "subtract"
	OpMultiply   = "multiply"
	OpDivide     = "divide"
	OpPower      = "power"
	OpSqrt       = "sqrt"
	OpPercentage = "percentage"
)

var (
	ErrMissingOperation     = errors.New("operation is required")
	ErrUnsupportedOperation = errors.New("operation is not supported")
	ErrMissingOperandA      = errors.New("operand a is required")
	ErrMissingOperandB      = errors.New("operand b is required")
	ErrDivisionByZero       = errors.New("division by zero")
	ErrNegativeSqrt         = errors.New("square root of a negative number is not defined")
	ErrUndefinedPower       = errors.New("zero raised to the power of zero is undefined")
	ErrNonFiniteOperand     = errors.New("operands must be finite numbers")
	ErrNonFiniteResult      = errors.New("calculation result is not a finite number")
)

// Calculate evaluates operation using operands a and b.
//
// A nil operand means the caller omitted it. Zero is a valid operand.
// Binary operations require both a and b. Square root requires only a
// and ignores b when it is present.
//
// Percentage is b percent of a, defined as (a * b) / 100.
// For example, a=200 and b=15 returns 30.
//
// Power uses math.Pow. Zero raised to the power of zero is rejected,
// including negative zero, because this API treats that input as undefined.
// A negative base raised to a non-integer power has no real result and is rejected.
// Operation names are case-sensitive.
func Calculate(operation string, a, b *float64) (float64, error) {
	switch operation {
	case OpAdd, OpSubtract, OpMultiply, OpDivide, OpPower, OpPercentage:
		if err := requireBinary(a, b); err != nil {
			return 0, err
		}
		if err := requireFinite(*a, *b); err != nil {
			return 0, err
		}
		return evalBinary(operation, *a, *b)
	case OpSqrt:
		if a == nil {
			return 0, ErrMissingOperandA
		}
		if err := requireFinite(*a); err != nil {
			return 0, err
		}
		return sqrt(*a)
	case "":
		return 0, ErrMissingOperation
	default:
		return 0, fmt.Errorf("%w: %q", ErrUnsupportedOperation, operation)
	}
}

func requireBinary(a, b *float64) error {
	if a == nil {
		return ErrMissingOperandA
	}
	if b == nil {
		return ErrMissingOperandB
	}
	return nil
}

func requireFinite(values ...float64) error {
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return ErrNonFiniteOperand
		}
	}
	return nil
}

func evalBinary(operation string, a, b float64) (float64, error) {
	// Calculations use float64, so results may contain
	// normal binary floating-point rounding differences.
	switch operation {
	case OpAdd:
		return finite(a + b)
	case OpSubtract:
		return finite(a - b)
	case OpMultiply:
		return finite(a * b)
	case OpDivide:
		return divide(a, b)
	case OpPower:
		return power(a, b)
	case OpPercentage:
		return percentage(a, b)
	default:
		return 0, fmt.Errorf("%w: %q", ErrUnsupportedOperation, operation)
	}
}

func power(a, b float64) (float64, error) {
	// math.Pow(0, 0) returns 1. This API rejects that input instead.
	if a == 0 && b == 0 {
		return 0, ErrUndefinedPower
	}
	return finite(math.Pow(a, b))
}

func divide(a, b float64) (float64, error) {
	// Reject both +0 and -0 before IEEE division can produce infinity.
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return finite(a / b)
}

func sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrt
	}
	return finite(math.Sqrt(a))
}

// percentage returns b percent of a, (a * b) / 100.
// Dividing the larger-magnitude operand by 100 first avoids some
// intermediate overflows from multiplying a and b directly. The final
// value is still checked for a non-finite result. This does not recover
// every result that is mathematically representable as a float64.
func percentage(a, b float64) (float64, error) {
	if math.Abs(a) >= math.Abs(b) {
		return finite((a / 100) * b)
	}
	return finite(a * (b / 100))
}

func finite(result float64) (float64, error) {
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrNonFiniteResult
	}
	return result, nil
}
