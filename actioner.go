// actioner.go
package taskwrappr

type Action struct {
	Execute func(s *ScriptRunner, args ...interface{}) (interface{}, error)
}

func NewAction(execute func(s *ScriptRunner, args ...interface{}) (interface{}, error)) *Action {
	return &Action{
		Execute: execute,
	}
}