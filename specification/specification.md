# BacVM Bytecode Instruction Specification
This page contains all the specifications for BacVM (shortened to BVM here) bytecode instructions. Moreover, there are other pages dediccated to go into more depth for some instructions, as well as giving examples.

## BVM Versioning
BVM follows a simple yet effective version scheme. A BVM version consists of the following:
1. A major version
2. A minor version

These two versions are put together into one version string, separated by periods:
> major.minor

Therefore, major version 3, minor version will take the following format:
> 3.2

When the major is incremented, the minor will reset to 0.

There are several guidelines related to compatibility between versions.
* Differing major versions will never be compatible.
* Differing minor versions, with the same major versions, will be compatible.

Take a VM running version 1.2 as an example:
* Bytecode versioned 1.3 and 1.x will run.
* Bytecode versioned 2.0 and 2.x will not run.

## Formatting & Structure
BVM bytecode instructions are compiled to a UTF-8 encoded text file.

#### The First Line
The first line is a special line, because it is not like any other line. The first line simply consists of the version string.

#### The Instructions
Each line represents one instruction. An instruction contains the operation and can contain zero or one argument(s), no more. The operation will always be one word, two letters, and is case sensitive! BVM specifies the operations to be all lowercase. In the case that there are arguments, these are separated from the instructions via a space.

#### Example
```
1.0.0
ex
```
In this example, the version, `1.0.0`, is defined in the first line. Then, the first and only instruction, `ex` (exit), is defined, which takes no parameters.

## Operations
Operations are tasks that BVM will execute. An operation can have an argument. There are two types of operations, feed and non-feed operations. Generally, feed operations require multiple instructions, non-feed operations only require one. These will be explored more in depth below.

#### The Buffer
The buffer is an important part of the VM. It is a first-in-last-out collection, storing data. It exists to make operations much easier and doable. Values are loaded onto the buffer using specific operations, and can be popped using `pb`. It is recommended to be implemented as a linked list.

#### Non-Feed Operations
These are as simple as operations can get. They are only one instruction.

Examples (with and without argument):
```
ex
gt 0
```
In the examples, `ex` takes no argument, `gt` takes the argument `11`.

#### Feed Operations
Sometimes, only one argument will not suffice. For example, if you want to declare a variable, you need to provide at the very least the name and the value, which is not possible with only 1 argument (technically you could split the argument, but this gets very messy). It is because of this that feed operations exist. The feed instructions are as follows:
1. Specify that a feed has started, and provide the type of feed operation.
2. Feed data (repeat until complete).
3. Complete the feed operation.

Feeds can be nested inside of each other. You can execute other operations during a feed, including starting a new feed.

As an example, presuming the operation "example" exists and requires three values:
```
fi example
pu 1
pu 2
pb
ff
```
This will initialize a feed operation of type "example", push the values `1` and `2`, and pop the buffered value as the third argument. The `ff` indicates that the feed has been finalized.

#### Operation List
\- indicates that there is no argument required.
? indicates that the argument is optional

Operation | Description | Argument
-- | -- | --
bd | Drops the first buffered value | \-
ex | Exits the VM | \-
ff | Stops a feed operation | \-
fi | Starts a feed operation | The feed type
gc | Garbage cleans a variable | The variable identifier
gt | Goes to a specific instruction | The instruction (0 based)
pb | Pushes the first buffered value | \-
pu | Pushes a specified value | The value to push?
sf | Decrements the current scope | - 
si | Increments the current scope | -
rn | Inverts the reading | -
rt | Terminates the reading | -
vl | Loads a variabl value into the buffer | The variable identifier

#### Feed Reference
Type | Description | Values
-- | -- | --
c_eg | Compares equal or greater than | Conditional
c_eq | Compares equality | Conditional
c_es | Compares equal or smaller than | Conditional
c_g | Compares greather than | Conditional
c_s | Compares smaller than | Conditional
f | Native function call | Function
v | Variable declaration | Variable

#### Feed Values

##### Conditional
1) The left hand side of the comparison.
2) The right hand side of the comparison.

This will start a new reader, with whether or not it's reading being determined by the output of the condition.

As an example, comparing equality on two numbers
```
fi c_eq
pu 2
pu 2
ff
```

##### Function
1. The function name.

A function call can specify more values. These are treated as parameters. Per parameter, there are two values that need to be pushed: the variable type and variable value. As you may notice, exactly these two values are pushed onto the buffer when you load a variable using `vl`.

As an example, executing the function `foo` with no parameters:
```
fi f
pu foo
ff
```

As an example, executing the function `stdout` with the variable "Hello, world":
```
fi f
pu stdout
pu string
pu Hello, world
ff
```

As an example, executing the function `stdout` referencing the variable `test` as a parameter:
```
fi f
pu stdout
vl test
pb
pb
ff
```

##### Variable
1. The variable identifier.
3. The initial value (can be blank for 0).

Variable values must be real numbers. No other type is allowed.

As an example, creating the variable `constant` with the value `19`:
```
fi v
pu constant
pu 19
ff
```
As an example, creating the copy variable `copy` with the value of `constant`:
```
fi v
pu copy
vl constant
pb
ff
```
