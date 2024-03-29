package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/arraying/bacvm/native"

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

// main will run all the specified files in the appropriate order.
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
	vm.Natives(&bacvm.Native{
		Name: "stdout",
		Task: native.StdOut,
	}, &bacvm.Native{
		Name: "stdoutln",
		Task: native.StdOutLn,
	})
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
				argument := parts[1]
				if comment := strings.Index(argument, ";"); comment != -1 {
					argument = argument[:comment-1]
				}
				instruction.Argument = argument
			}
			instructions = append(instructions, instruction)
		}
	}
	err = vm.Run(func() []*bacvm.Instruction {
		return instructions
	})
	if err != nil {
		fmt.Println(vm.Dump())
	}
	return err
}
