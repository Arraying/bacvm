package bacvm

import (
	"fmt"
	"strconv"
)

// Variable represents a number. It is stored as a 64 bit floating point number, the maximum value of integers is 53 bit.
type Variable float64

// Compare compares two variables.
func (variable Variable) Compare(operation string, other Variable) bool {
	switch operation {
	case feederComparisonTypeEg:
		return variable >= other
	case feederComparisonTypeEq:
		return variable == other
	case feederComparisonTypeEs:
		return variable <= other
	case feederComparisonTypeG:
		return variable > other
	case feederComparisonTypeS:
		return variable < other
	}
	return false
}

// Value returns the value of the variable as a string.
func (variable Variable) Value() string {
	return fmt.Sprintf("%f", variable)
}

// VariableBounded determines whether the integer can be represented within 53 bits, the highest integer value.
func VariableBounded(input int) bool {
	if input < 0 {
		input *= -1
	}
	count := 0
	for input != 0 {
		count++
		input >>= 1
	}
	return count <= 53
}

// variableCreate creates a variable from a piece of string.
// If the string cannot be parsed into a 64 bit floating point integer then it will fall back to 0.
func variableCreate(data string) Variable {
	val, _ := strconv.ParseFloat(data, 64)
	return Variable(val)
}
