# BacVM Bytecode Variable Specification
BacVM has three types of variables:
1. number (64 bit FP).
2. bool.
3. string.

## Typing

The type cannot be inferred from the initial value (BVM has no way to tell whether a value is `true` or `"true"`, for example), so the initial type must be specified. If no initial value is passed in, the zero value for that type will be used:
Type | Zero Value
- | -
number | 0
bool | false
string | *empty string*

When a variable is re-assigned to a new value, the variable's type will change to the new value.

## Comparisons
Variables of different types will never be equal to each other. The order of magnitude is as follows:
> string > number > bool

Presuming three variables, `s` (string), `n` (number), `b` (bool), the following will be the case:
* s > n evaluates to true
* b > n evaluates to false
* b = s evaluates to false
* n >= b evaluates to true
* And so on...

All variables that are not of type boolean will have the value `false` in logical operations.