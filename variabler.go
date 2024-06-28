// variabler.go
package taskwrappr

type Variable struct {
	Value interface{}
}

func NewVariable(value interface{}) *Variable {
	return &Variable{
		Value: value,
	}
}
