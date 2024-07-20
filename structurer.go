// structurer.go
package taskwrappr

import "fmt"

type Action struct {
    ExecuteFunc  func(s *Script, args ...interface{}) (interface{}, error)
    Arguments    []interface{}
    Block        *Block
    ValidateFunc func(s *Script, a *Action) error
}

type Variable struct {
	Value interface{}
    Type  VariableType
}

type MemoryMap struct {
	Actions   map[string]*Action
	Variables map[string]*Variable
}

type Block struct {
    Actions    []*Action
	Executed   bool
    Memory     *MemoryMap
	LastResult interface{}
}

type VariableType int

const (
	StringType VariableType = iota
	IntegerType
	FloatType
	BooleanType
	InvalidType
)

func (v VariableType) String() string {
    switch v {
    case StringType:
        return "string"
    case IntegerType:
        return "integer"
    case FloatType:
        return "float"
    case BooleanType:
        return "boolean"
    default:
        return "invalid"
    }
}

func NewVariable(value interface{}, variableType VariableType) *Variable {
	return &Variable{
		Value: value,
        Type:  variableType,
	}
}

func NewMemoryMap() *MemoryMap {
	return &MemoryMap{
		Actions:   make(map[string]*Action),
		Variables: make(map[string]*Variable),
	}
}

func (m *MemoryMap) GetAction(name string) *Action {
	action, ok := m.Actions[name]
	if !ok {
		return nil
	}
	return action
}

func (m *MemoryMap) GetVariable(name string) *Variable {
	variable, ok := m.Variables[name]
	if !ok {
		return nil
	}
	return variable
}

func NewBlock() *Block {
    return &Block{
        Actions: []*Action{},
    }
}

func NewAction(executeFunc func(s *Script, args ...interface{}) (interface{}, error), validateFunc func(s *Script, a *Action) error) *Action {
    return &Action{
        ExecuteFunc:  executeFunc,
        ValidateFunc: validateFunc,
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
        case *Variable:
            processedArgs[i] = v.Value
        default:
            return nil, fmt.Errorf("unsupported argument type: %T", arg)
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

func (a *Action) Validate(s *Script) (error) {
    if a.ValidateFunc == nil {
        return nil
    }

    return a.ValidateFunc(s, a)
}