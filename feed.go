package bacvm

import "fmt"

type (
	feeder interface {
		feed(*VM, string) error
		finalize(*VM) error
	}
	feederComparison struct {
		a     interface{}
		b     interface{}
		index int
		kind  string
	}
	feederFunction struct {
		arguments []*Variable
		index     int
		name      string
	}
	feederVariable struct {
		index    int
		variable *Variable
	}
	feeding struct {
		current  feeder
		previous *feeding
	}
)

const (
	feederVarKindMutable  = "m"
	feederVarKindConstant = "c"
)

const (
	feederComparisonTypeA  = "l_a"
	feederComparisonTypeEg = "c_eg"
	feederComparisonTypeEq = "c_eq"
	feederComparisonTypeEs = "c_es"
	feederComparisonTypeG  = "c_g"
	feederComparisonTypeO  = "l_o"
	feederComparisonTypeS  = "c_s"
	feederVar              = "v"
)

func (feeder *feederComparison) feed(vm *VM, argument string) error {
	// TODO
	return nil
}

func (feeder *feederComparison) finalize(vm *VM) error {
	// TODO
	return nil
}

func (feeder *feederFunction) feed(vm *VM, argument string) error {
	// TODO
	return nil
}

func (feeder *feederVariable) feed(vm *VM, argument string) error {
	var log string
	switch feeder.index {
	case 0:
		feeder.variable = &Variable{
			Name: argument,
		}
		log = "Creating variable " + argument
	case 1:
		switch argument {
		case VariableTypeNumber, VariableTypeBoolean, VariableTypeString:
			feeder.variable.Type = argument
		default:
			return ErrorOperationArgument
		}
		log = "Set type " + argument
	case 2:
		if err := feeder.variable.Assign(argument); err != nil {
			return err
		}
		log = "Assigned value"
	default:
		return ErrorFeedSize
	}
	feeder.index++
	if vm.Verbose {
		fmt.Println(log)
	}
	return nil
}

func (feeder *feederVariable) finalize(vm *VM) error {
	if feeder.index != 3 {
		return ErrorFeedSize
	}
	vm.Scope.Put(feeder.variable)
	return nil
}
