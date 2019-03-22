package bacvm

import (
	"fmt"
	"strconv"
)

type (
	// Variable represents a VM variable.
	Variable struct {
		Name  string
		Type  string
		Value interface{}
	}
)

const (
	// VariableTypeNumber represents a 64 bit floating point number.
	VariableTypeNumber = "number"
	// VariableTypeBoolean represents a true or false value.
	VariableTypeBoolean = "bool"
	// VariableTypeString represents a string.
	VariableTypeString = "string"
)

// Assign assigns the variable to the string value.
func (variable *Variable) Assign(value string) (err error) {
	switch variable.Type {
	case VariableTypeNumber:
		if value == "" {
			variable.Value = 0
			return
		}
		fp, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return ErrorVariableType
		}
		variable.Value = fp
	case VariableTypeBoolean:
		if value == "" {
			variable.Value = false
			return
		}
		boolean, err := strconv.ParseBool(value)
		if err != nil {
			return ErrorVariableType
		}
		variable.Value = boolean
	case VariableTypeString:
		variable.Value = value
	}
	return
}

// Compare compares two variables.
func (variable *Variable) Compare(comparison string, other *Variable) bool {
	switch comparison {
	case feederComparisonTypeA:
		return variable.Type == VariableTypeBoolean && variable.Value.(bool) && other.Type == VariableTypeBoolean && other.Value.(bool)
	case feederComparisonTypeEg:
		// TODO
	case feederComparisonTypeEq:
		if variable.Type == other.Type && variable.Value == other.Value {
			return true
		}
		return false
	case feederComparisonTypeEs:
		// TODO
	case feederComparisonTypeG:
		// TODO
	case feederComparisonTypeO:
		if variable.Type == VariableTypeBoolean && variable.Value.(bool) {
			return true
		}
		if other.Type == VariableTypeBoolean && variable.Value.(bool) {
			return true
		}
		return false
	case feederComparisonTypeS:
		// TODO
	default:
		return false
	}
	// TODO
	return false
}

// ValueString gets the value as a string.
func (variable *Variable) ValueString() string {
	return fmt.Sprintf("%v", variable.Value)
}

// Weight gets the weight of the variable type.
func (variable *Variable) Weight() int {
	switch variable.Type {
	case VariableTypeNumber:
		return 2
	case VariableTypeBoolean:
		return 3
	case VariableTypeString:
		return 1
	default:
		return 0
	}
}
