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
    switch t := arg.Type; t {
    case BooleanType:
        return []*Variable{NewVariable(arg.Value, BooleanType)}, nil
    case StringType:
        return []*Variable{NewVariable(arg.Value != "", BooleanType)}, nil
    default:
        return nil, fmt.Errorf("unsupported argument type: %s", arg.Type)
    }
}

func IfActionValidator(s *Script, a *Action) error {
    if a.Block == nil {
        return fmt.Errorf("'if' action must have a preceding code block")
    }
    return nil
}

func ElseIfAction(s *Script, args ...*Variable) ([]*Variable, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("'else if' action requires exactly one argument")
    }

    if s.CurrentBlock.LastResult != nil {
        if lastResult, ok := s.CurrentBlock.LastResult.Value.(bool); ok {
            if lastResult {
                return nil, nil
            } else {
                switch v := args[0].Value.(type) {
                case bool:
                    return []*Variable{NewVariable(v, BooleanType)}, nil
                case string:
                    return []*Variable{NewVariable(v != "", BooleanType)}, nil
                default:
                    return nil, fmt.Errorf("unsupported argument type: %T", v)
                }
            }
        }
    }

    return nil, nil
}

func ElseIfActionValidator(s *Script, a *Action) error {
    if a.Block == nil {
        return fmt.Errorf("'else if' action must have a preceding code block")
    }
    return nil
}

func ElseAction(s *Script, args ...*Variable) ([]*Variable, error) {
    if len(args) == 0 {
        if s.CurrentBlock.LastResult != nil {
            if lastResult, ok := s.CurrentBlock.LastResult.Value.(bool); ok {
                return []*Variable{NewVariable(!lastResult, BooleanType)}, nil
            }
        } else {
            return nil, nil
        }
    } else if len(args) == 1 {
        arg := args[0]
        switch t := arg.Type; t {
        case BooleanType:
            return []*Variable{NewVariable(!arg.Value.(bool), BooleanType)}, nil
        case StringType:
            return []*Variable{NewVariable(arg.Value == "", BooleanType)}, nil
        default:
            return nil, fmt.Errorf("unsupported argument type: %s", arg)
        }
    }

    return nil, fmt.Errorf("too many arguments for 'else' action")
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
