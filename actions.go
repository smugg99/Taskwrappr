// actions.go
package taskwrappr

type Action struct {
    Block        *Block
	arguments    []*Action
	executeFunc  func(s *Script, args ...*Variable) ([]*Variable, error)
    validateFunc func(s *Script, a *Action) error
}

func NewAction(executeFunc func(s *Script, args ...*Variable) ([]*Variable, error), validateFunc func(s *Script, a *Action) error) *Action {
    return &Action{
        executeFunc:  executeFunc,
        validateFunc: validateFunc,
    }
}

func CloneAction(a *Action) *Action {
	return &Action{
		Block:        a.Block,
		arguments:    a.arguments,
		executeFunc:  a.executeFunc,
		validateFunc: a.validateFunc,
	}
}

func (a *Action) ProcessArgs(s *Script) ([]*Variable, error) {
	args := a.GetArguments()
    processedArgs := make([]*Variable, len(args))

    for i, arg := range args {
        processedArg, err := arg.Execute(s)
		if err != nil {
			return nil, err
		}

		if len(processedArg) == 1 {
			processedArgs[i] = processedArg[0]
		} else {
			processedArgs[i] = NewVariable(processedArg, ArrayType)
		}
    }

    return processedArgs, nil
}

func (a *Action) SetArguments(args []*Action) {
	a.arguments = args
}

func (a *Action) GetArguments() ([]*Action) {
	return a.arguments
}

func (a *Action) Execute(s *Script) ([]*Variable, error) {
    processedArgs, err := a.ProcessArgs(s)
    if err != nil {
        return nil, err
    }
    return a.executeFunc(s, processedArgs...)
}

func (a *Action) Validate(s *Script) (error) {
    if a.validateFunc == nil {
        return nil
    }

    return a.validateFunc(s, a)
}