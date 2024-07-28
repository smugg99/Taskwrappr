// internals.go

package taskwrappr

import (
    "fmt"
    "time"
)

func printVariable(v *Variable) {
    switch v.Type {
    case StringType, IntegerType, FloatType, BooleanType:
        fmt.Print(v.Value)
    case ArrayType:
        array := v.Value.([]*Variable)
        fmt.Printf("%c", BracketOpenSymbol)
        for i, elem := range array {
            printVariable(elem)
            if i != len(array)-1 {
                fmt.Printf("%c", SpaceSymbol)
            }
        }
        fmt.Printf("%c", BracketCloseSymbol)
    case NilType:
        fmt.Print("nil")
    default:
        fmt.Printf("unsupported argument type: %v\n", v.Type)
    }
}

func GetInternals() *MemoryMap {
    actions := make(map[string]*Action)
    variables := make(map[string]*Variable)

    actions["if"]     = NewAction(IfAction, IfActionValidator)
    actions["elseIf"] = NewAction(ElseIfAction, ElseIfActionValidator)
    actions["else"]   = NewAction(ElseAction, ElseActionValidator)
    actions["for"]    = NewAction(ForAction, ForActionValidator)
    actions["print"]  = NewAction(PrintAction, nil)
    actions["wait"]   = NewAction(WaitAction, nil)
    actions["pass"]   = NewAction(PassAction, nil)
    actions["type"]   = NewAction(TypeAction, nil)
    actions["bool"]   = NewAction(BoolAction, nil)
    actions["int"]    = NewAction(IntAction, nil)
    actions["float"]  = NewAction(FloatAction, nil)
    actions["string"] = NewAction(StringAction, nil)

    variables[TrueString] = NewVariable(true, BooleanType)
    variables[FalseString] = NewVariable(false, BooleanType)
    variables[NilString] = NewVariable(nil, NilType)

    return &MemoryMap{
        Actions:   actions,
        Variables: variables,
    }
}

func IfAction(s *Script, args ...*Variable) ([]*Variable, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("'if' action requires exactly one argument")
    }

    arg := args[0]
    if arg.Type != BooleanType {
        return nil, fmt.Errorf("'if' action requires a boolean argument")
    }

    if arg.Value.(bool) {
        s.CurrentBlock.LastResult = arg
        return []*Variable{NewVariable(true, BooleanType)}, nil
    }
    s.CurrentBlock.LastResult = NewVariable(false, BooleanType)

    return []*Variable{NewVariable(false, BooleanType)}, nil
}

func IfActionValidator(s *Script, a *Action) error {
    if a.Block == nil {
        return fmt.Errorf("'if' action must have a sequent code block")
    }
    return nil
}

func ElseIfAction(s *Script, args ...*Variable) ([]*Variable, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("'elseif' action requires exactly one argument")
    }
    arg := args[0]
    if arg.Type != BooleanType {
        return nil, fmt.Errorf("'elseif' action requires a boolean argument")
    }
    if s.CurrentBlock.LastResult != nil {
        if lastResult, ok := s.CurrentBlock.LastResult.Value.(bool); ok && !lastResult {
            if arg.Value.(bool) {
                s.CurrentBlock.LastResult = arg
                return []*Variable{NewVariable(true, BooleanType)}, nil
            }
        }
    }
    return []*Variable{NewVariable(false, BooleanType)}, nil
}

func ElseIfActionValidator(s *Script, a *Action) error {
    if a.Block == nil {
        return fmt.Errorf("'else if' action must have a sequent code block")
    }
    return nil
}

func ElseAction(s *Script, args ...*Variable) ([]*Variable, error) {
    if len(args) > 1 {
        return nil, fmt.Errorf("'else' action requires at most one argument")
    }
    
    if s.CurrentBlock.LastResult != nil {
        if lastResult, ok := s.CurrentBlock.LastResult.Value.(bool); ok && !lastResult {
            s.CurrentBlock.LastResult = NewVariable(true, BooleanType)
            return []*Variable{NewVariable(true, BooleanType)}, nil
        }
    }
    s.CurrentBlock.LastResult = NewVariable(false, BooleanType)

    return []*Variable{NewVariable(false, BooleanType)}, nil
}

func ElseActionValidator(s *Script, a *Action) error {
    if a.Block == nil {
        return fmt.Errorf("'else' action must have a code block")
    }
    return nil
}

func ForAction(s *Script, args ...*Variable) ([]*Variable, error) {
    if len(args) == 0 {
        return nil, fmt.Errorf("'for' action requires at least one argument")
    }
    if len(args) == 1 {
        arg := args[0].Value
        switch v := arg.(type) {
        case bool:
            return []*Variable{NewVariable(v, BooleanType)}, nil
        case string:
            return []*Variable{NewVariable(v != "", BooleanType)}, nil
        default:
            return nil, fmt.Errorf("unsupported argument type: %T", arg)
        }
    }
    return nil, fmt.Errorf("too many arguments for 'for' action")
}

func ForActionValidator(s *Script, a *Action) error {
    if a.Block == nil {
        return fmt.Errorf("'for' action must have a code block")
    }
    return nil
}

func PrintAction(s *Script, args ...*Variable) ([]*Variable, error) {
    for i, arg := range args {
        printVariable(arg)
        if i != len(args)-1 {
            fmt.Print(" ")
        }
    }
    fmt.Println()

    return nil, nil
}

func WaitAction(s *Script, args ...*Variable) ([]*Variable, error) {
    if len(args) < 1 {
        return nil, fmt.Errorf("'wait' action requires at least 1 argument")
    }
    durationStr := fmt.Sprintf("%v", args[0].Value)
    duration, err := time.ParseDuration(durationStr + "ms")
    if err != nil {
        return nil, err
    }
    time.Sleep(duration)
    return nil, nil
}

func PassAction(s *Script, args ...*Variable) ([]*Variable, error) { 
    return args, nil
}

func TypeAction(s *Script, args ...*Variable) ([]*Variable, error) { 
    if len(args) < 1 {
        return nil, fmt.Errorf("'type' action requires exactly 1 argument")
    }

    return []*Variable{NewVariable(args[0].Type.String(), StringType)}, nil
}

func BoolAction(s *Script, args ...*Variable) ([]*Variable, error) { 
    if len(args) < 1 {
        return nil, fmt.Errorf("'bool' action requires exactly 1 argument")
    }

    arg := args[0]
    value, err := arg.toBool() 
    if err != nil {
        return nil, err
    }

    return []*Variable{NewVariable(value, BooleanType)}, nil
}

func IntAction(s *Script, args ...*Variable) ([]*Variable, error) { 
    if len(args) < 1 {
        return nil, fmt.Errorf("'int' action requires exactly 1 argument")
    }

    arg := args[0]
    value, err := arg.toInt() 
    if err != nil {
        return nil, err
    }

    return []*Variable{NewVariable(value, IntegerType)}, nil
}

func FloatAction(s *Script, args ...*Variable) ([]*Variable, error) { 
    if len(args) < 1 {
        return nil, fmt.Errorf("'float' action requires exactly 1 argument")
    }

    arg := args[0]
    value, err := arg.toFloat() 
    if err != nil {
        return nil, err
    }

    return []*Variable{NewVariable(value, FloatType)}, nil
}

func StringAction(s *Script, args ...*Variable) ([]*Variable, error) { 
    if len(args) < 1 {
        return nil, fmt.Errorf("'string' action requires exactly 1 argument")
    }

    arg := args[0]
    value, err := arg.toString()
    if err != nil {
        return nil, err
    }

    return []*Variable{NewVariable(value, StringType)}, nil
}