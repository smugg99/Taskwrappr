// structurer.go
package taskwrappr

import (
	"fmt"
	"strconv"
)

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

func (v *Variable) CastTo(targetType VariableType) (interface{}, error) {
	switch targetType {
	case StringType:
		return v.toString()
	case IntegerType:
		return v.toInt()
	case FloatType:
		return v.toFloat()
	case BooleanType:
		return v.toBool()
	default:
		return nil, fmt.Errorf("invalid variable type: %v", targetType)
	}
}

func (v *Variable) toString() (string, error) {
	switch v.Type {
	case StringType:
		return v.Value.(string), nil
	case IntegerType:
		return strconv.Itoa(v.Value.(int)), nil
	case FloatType:
		return fmt.Sprintf("%g", v.Value.(float64)), nil
	case BooleanType:
		return strconv.FormatBool(v.Value.(bool)), nil
	default:
		return "", fmt.Errorf("cannot convert %v to string", v.Type)
	}
}

func (v *Variable) toInt() (int, error) {
	switch v.Type {
	case StringType:
		if i, err := strconv.Atoi(v.Value.(string)); err == nil {
			return i, nil
		} else if f, err := strconv.ParseFloat(v.Value.(string), 64); err == nil {
			return int(f), nil
		}
	case IntegerType:
		return v.Value.(int), nil
	case FloatType:
		return int(v.Value.(float64)), nil
	case BooleanType:
		if v.Value.(bool) {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("cannot convert %v to integer", v.Type)
}

func (v *Variable) toFloat() (float64, error) {
	switch v.Type {
	case StringType:
		if f, err := strconv.ParseFloat(v.Value.(string), 64); err == nil {
			return f, nil
		}
	case IntegerType:
		return float64(v.Value.(int)), nil
	case FloatType:
		return v.Value.(float64), nil
	case BooleanType:
		if v.Value.(bool) {
			return 1.0, nil
		}
		return 0.0, nil
	}
	return 0.0, fmt.Errorf("cannot convert %v to float", v.Type)
}

func (v *Variable) toBool() (interface{}, error) {
	switch v.Type {
	case StringType:
		if b, err := strconv.ParseBool(v.Value.(string)); err == nil {
			return b, nil
		}
	case IntegerType:
	case FloatType:
		var value float64
		if v.Type == IntegerType {
			value = float64(v.Value.(int))
		} else {
			value = v.Value.(float64)
		}
		return value != 0, nil
	case BooleanType:
		return v.Value, nil
	}
	return nil, fmt.Errorf("cannot convert %v to boolean", v.Type)
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