package native

import (
	"fmt"

	"github.com/arraying/bacvm"
)

func stdout(vm *bacvm.VM, vars []*bacvm.Variable) bacvm.Variable {
	args := make([]string, 0)
	for _, variable := range vars {
		args = append(args, variable.ValueString())
	}
	i, _ := fmt.Println(args)
	return bacvm.Variable{
		Type:  bacvm.VariableTypeNumber,
		Value: float64(i),
	}
}
