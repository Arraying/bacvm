package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/arraying/bacvm"
)

var (
	verbose bool
	files   []string
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Whether to be verbose.")
	flag.Parse()
	files = flag.Args()
}

func main() {
	for _, arg := range files {
		if err := run(arg); err != nil {
			panic(err)
		}
	}
}

func run(arg string) error {
	vm := bacvm.VM{
		Verbose: verbose,
	}
	vm.Natives()
	file, err := os.Open(arg)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	instructions := make([]*bacvm.Instruction, 0)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if i == 0 {
			if !bacvm.Offer(line) {
				return errors.New("version incompatibility: VM running " + bacvm.Version() + " bytecode running " + line)
			}
		} else {
			parts := strings.SplitN(line, " ", 2)
			instruction := &bacvm.Instruction{
				Operation: parts[0],
			}
			if len(parts) > 1 {
				instruction.Argument = parts[1]
			}
			instructions = append(instructions, instruction)
		}
	}
	err = vm.Run(func() []*bacvm.Instruction {
		return instructions
	})
	if err != nil {
		return err
	}
	fmt.Println(vm.Scope.Get("name"))
	return nil
}
