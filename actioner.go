// actioner.go
package taskwrappr

type Action struct {
	Execute func(args ...interface{}) (interface{}, error)
}

func NewAction(execute func(args ...interface{}) (interface{}, error)) *Action {
	return &Action{
		Execute: execute,
	}
}