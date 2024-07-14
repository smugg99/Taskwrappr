// actioner.go
package taskwrappr

type Block struct {
	Actions []*Action
	Memory  *MemoryMap
}

type Action struct {
	ExecuteFunc func(s *Script, args ...interface{}) (interface{}, error)
	Arguments   []interface{}
	Block       *Block
}

func (a *Action) ProcessArgs(s *Script) ([]interface{}, error) {
	processedArgs := make([]interface{}, len(a.Arguments))

	for i, arg := range a.Arguments {
		switch v := arg.(type) {
		case *Action:
			processedArg, err := v.Execute(s)
			if err != nil {
				return nil, err
			}
			processedArgs[i] = processedArg
			// fmt.Printf("Type: %T, Value: %v, Index: %d\n", processedArg, processedArg, i)
		default:
			processedArgs[i] = v
			// fmt.Printf("Type: %T, Value: %v, Index: %d\n", v, v, i)
		}
	}

	return processedArgs, nil
}

func (a *Action) Execute(s *Script) (interface{}, error) {
	processedArgs, err := a.ProcessArgs(s)
	if err != nil {
		return nil, err
	}

	// fmt.Println("Executing action", a, "with args", processedArgs)

	return a.ExecuteFunc(s, processedArgs...)
}

func NewAction(executeFunc func(s *Script, args ...interface{}) (interface{}, error)) *Action {
	return &Action{
		ExecuteFunc: executeFunc,
	}
}