// internals.go

package taskwrappr

import (
    "fmt"
    "time"
)

func GetInternals() *MemoryMap {
    actions := make(map[string]*Action)
    variables := make(map[string]*Variable)

    actions["if"] = NewAction(IfAction, IfActionValidator)
    actions["else"] = NewAction(ElseAction, ElseActionValidator)
    actions["for"] = NewAction(ForAction, ForActionValidator)
    actions["print"] = NewAction(PrintAction, nil)
    actions["wait"] = NewAction(WaitAction, nil)
    actions["pass"] = NewAction(PassAction, nil)
    actions["and"] = NewAction(AndAction, nil)
    actions["or"] = NewAction(OrAction, nil)
    actions["not"] = NewAction(NotAction, nil)
    actions["xor"] = NewAction(XorAction, nil)
    actions["nand"] = NewAction(NandAction, nil)

    return &MemoryMap{
        Actions:   actions,
        Variables: variables,
    }
}

func IfAction(s *Script, args ...interface{}) (interface{}, error) {
    if len(args) == 0 {
        return false, fmt.Errorf("'if' action requires exactly one argument")
    }

    if len(args) == 1 {
        switch v := args[0].(type) {
        case bool:
            return v, nil
        case string:
            return v != "", nil
        default:
            return false, fmt.Errorf("unsupported argument type: %T", args[0])
        }
    } else {
        return false, fmt.Errorf("too many arguments for 'if' action")
    }
}

func IfActionValidator(s *Script, a *Action) error {
    if a.Block == nil {
        return fmt.Errorf("'if' action must have a preceding code block")
    }
    return nil
}

func ElseAction(s *Script, args ...interface{}) (interface{}, error) {
    if s.CurrentBlock.LastResult != nil {
        if lastResult, ok := s.CurrentBlock.LastResult.(bool); ok {
            return !lastResult, nil
        }
    }

    return false, nil
}

func ElseActionValidator(s *Script, a *Action) error {
    if a.Block == nil {
        return fmt.Errorf("'else' action must have a code block")
    }
    return nil
}

func ForAction(s *Script, args ...interface{}) (interface{}, error) {
    if len(args) == 0 {
        return false, fmt.Errorf("'for' action requires at least one argument")
    }

    if len(args) == 1 {
        switch v := args[0].(type) {
        case bool:
            return v, nil
        case string:
            return v != "", nil
        default:
            return false, fmt.Errorf("unsupported argument type: %T", args[0])
        }
    } else {
        return false, fmt.Errorf("too many arguments for 'for' action")
    }
}

func ForActionValidator(s *Script, a *Action) error {
    if a.Block == nil {
        return fmt.Errorf("'for' action must have a code block")
    }
    return nil
}

func PrintAction(s *Script, args ...interface{}) (interface{}, error) {
    fmt.Println(args...)
    return args, nil
}

func WaitAction(s *Script, args ...interface{}) (interface{}, error) {
    if len(args) < 1 {
        return nil, fmt.Errorf("'wait' action requires at least 1 argument")
    }

    durationStr := fmt.Sprintf("%v", args[0])
    duration, err := time.ParseDuration(durationStr + "ms")
    if err != nil {
        return nil, err
    }

    time.Sleep(duration)
    return nil, nil
}

func PassAction(s *Script, args ...interface{}) (interface{}, error) {
    if len(args) == 0 {
        return false, fmt.Errorf("'pass' action requires at least one argument")
    }

    return args, nil
}

func AndAction(s *Script, args ...interface{}) (interface{}, error) {
    if len(args) == 0 {
        return false, fmt.Errorf("'and' action requires at least one argument")
    }
    for _, arg := range args {
        if v, ok := arg.(bool); ok {
            if !v {
                return false, nil
            }
        } else {
            return false, fmt.Errorf("'and' action only supports boolean arguments")
        }
    }
    return true, nil
}

func OrAction(s *Script, args ...interface{}) (interface{}, error) {
    if len(args) == 0 {
        return false, fmt.Errorf("'or' action requires at least one argument")
    }
    for _, arg := range args {
        if v, ok := arg.(bool); ok {
            if v {
                return true, nil
            }
        } else {
            return false, fmt.Errorf("'or' action only supports boolean arguments")
        }
    }
    return false, nil
}

func NotAction(s *Script, args ...interface{}) (interface{}, error) {
    if len(args) != 1 {
        return false, fmt.Errorf("'not' action requires exactly one argument")
    }
    if v, ok := args[0].(bool); ok {
        return !v, nil
    }
    return false, fmt.Errorf("'not' action only supports a boolean argument")
}

func XorAction(s *Script, args ...interface{}) (interface{}, error) {
    if len(args) < 2 {
        return false, fmt.Errorf("'xor' action requires at least two arguments")
    }
    countTrue := 0
    for _, arg := range args {
        if v, ok := arg.(bool); ok {
            if v {
                countTrue++
            }
        } else {
            return false, fmt.Errorf("'xor' action only supports boolean arguments")
        }
    }
    return countTrue%2 == 1, nil
}

func NandAction(s *Script, args ...interface{}) (interface{}, error) {
    if len(args) == 0 {
        return false, fmt.Errorf("'nand' action requires at least one argument")
    }
    for _, arg := range args {
        if v, ok := arg.(bool); ok {
            if !v {
                return true, nil
            }
        } else {
            return false, fmt.Errorf("'nand' action only supports boolean arguments")
        }
    }
    return false, nil
}
