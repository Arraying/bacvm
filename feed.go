package bacvm

type (
	feeder interface {
		feed(*VM, string) error
		finalize(*VM) error
	}
	feederComparison struct {
		a     Variable
		b     Variable
		index int
		kind  string
	}
	feederFunction struct {
		arguments  []Variable
		name       string
		parameters bool
	}
	feederVariable struct {
		index    int
		name     string
		variable Variable
	}
	feeding struct {
		current  feeder
		previous *feeding
	}
)

const (
	feederComparisonTypeEg = "c_eg"
	feederComparisonTypeEq = "c_eq"
	feederComparisonTypeEs = "c_es"
	feederComparisonTypeG  = "c_g"
	feederComparisonTypeS  = "c_s"
	feederFunc             = "f"
	feederVar              = "v"
)

func (feeder *feederComparison) feed(vm *VM, argument string) error {
	switch feeder.index {
	case 0:
		feeder.a = variableCreate(argument)
	case 1:
		feeder.b = variableCreate(argument)
	default:
		return ErrorFeedSize
	}
	feeder.index++
	return nil
}

func (feeder *feederComparison) finalize(vm *VM) error {
	if feeder.index != 2 {
		return ErrorFeedQuantity
	}
	res := feeder.a.Compare(feeder.kind, feeder.b)
	reading := &reading{
		current:  res,
		previous: vm.reading,
	}
	vm.reading = reading
	return nil
}

func (feeder *feederFunction) feed(vm *VM, argument string) error {
	if feeder.parameters {
		feeder.arguments = append(feeder.arguments, variableCreate(argument))
	} else {
		feeder.name = argument
		feeder.parameters = true
	}
	return nil
}

func (feeder *feederFunction) finalize(vm *VM) error {
	fn := vm.natives[feeder.name]
	if fn == nil {
		return ErrorFunctionReference
	}
	result := fn(vm, feeder.arguments)
	vm.bufferPush(result.Value())
	return nil
}

func (feeder *feederVariable) feed(vm *VM, argument string) error {
	switch feeder.index {
	case 0:
		feeder.name = argument
	case 1:
		feeder.variable = variableCreate(argument)
	default:
		return ErrorFeedSize
	}
	feeder.index++
	return nil
}

func (feeder *feederVariable) finalize(vm *VM) error {
	if feeder.index != 2 {
		return ErrorFeedSize
	}
	vm.Scope.Put(feeder.name, feeder.variable)
	return nil
}
