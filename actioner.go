package taskwrappr

type Block struct {
    Actions    []*Action
	Executed   bool
	LastResult interface{}
}

type Action struct {
    ExecuteFunc func(s *Script, args ...interface{}) (interface{}, error)
    Arguments   []interface{}
    Block       *Block
}

func NewBlock() *Block {
    return &Block{
        Actions: []*Action{},
    }
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
        default:
            processedArgs[i] = v
        }
    }

    return processedArgs, nil
}

func (a *Action) Execute(s *Script) (interface{}, error) {
	processedArgs, err := a.ProcessArgs(s)
	if err != nil {
		return nil, err
	}

	return a.ExecuteFunc(s, processedArgs...)
}

func NewAction(executeFunc func(s *Script, args ...interface{}) (interface{}, error)) *Action {
    return &Action{
        ExecuteFunc: executeFunc,
    }
}