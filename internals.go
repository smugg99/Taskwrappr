// internals.go

package taskwrappr

import (
    "fmt"
    "time"
)

func GetInternals() *MemoryMap {
    actions := make(map[string]*Action)
    variables := make(map[string]*Variable)

    actions["print"] = NewAction(func(s *Script, args ...interface{}) (interface{}, error) {
        fmt.Println(args...)
        return args, nil
    }, nil)

    actions["wait"] = NewAction(func(s *Script, args ...interface{}) (interface{}, error) {
        if len(args) < 1 {
            return nil, fmt.Errorf("wait action requires at least 1 argument")
        }

        durationStr := fmt.Sprintf("%v", args[0])
        duration, err := time.ParseDuration(durationStr + "ms")
        if err != nil {
            return nil, err
        }

        time.Sleep(duration)
        return nil, nil
    }, nil)

    actions["if"] = NewAction(func(s *Script, args ...interface{}) (interface{}, error) {
        if len(args) == 0 {
            return false, fmt.Errorf("if function requires exactly one argument")
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
            return false, fmt.Errorf("too many arguments for if function")
        }
    }, func(s *Script, a *Action) error {
        if a.Block == nil {
            return fmt.Errorf("if action must have a preceding code block")
        }
        return nil
    })

    actions["else"] = NewAction(func(s *Script, args ...interface{}) (interface{}, error) {
        if s.CurrentBlock == nil {
            return false, fmt.Errorf("else action must have a preceding code block")
        }
        if s.CurrentBlock != nil && s.CurrentBlock.LastResult != nil {
            if lastResult, ok := s.CurrentBlock.LastResult.(bool); ok {
                return !lastResult, nil
            }
        }

        return false, nil
    }, func(s *Script, a *Action) error {
        if a.Block == nil {
            return fmt.Errorf("else action must have a preceding code block")
        }
        return nil
    })

    actions["and"] = NewAction(func(s *Script, args ...interface{}) (interface{}, error) {
        if len(args) == 0 {
            return false, fmt.Errorf("and function requires at least one argument")
        }
        for _, arg := range args {
            if v, ok := arg.(bool); ok {
                if !v {
                    return false, nil
                }
            } else {
                return false, fmt.Errorf("and function only supports boolean arguments")
            }
        }
        return true, nil
    }, nil)

    actions["or"] = NewAction(func(s *Script, args ...interface{}) (interface{}, error) {
        if len(args) == 0 {
            return false, fmt.Errorf("or function requires at least one argument")
        }
        for _, arg := range args {
            if v, ok := arg.(bool); ok {
                if v {
                    return true, nil
                }
            } else {
                return false, fmt.Errorf("or function only supports boolean arguments")
            }
        }
        return false, nil
    }, nil)

    actions["not"] = NewAction(func(s *Script, args ...interface{}) (interface{}, error) {
        if len(args) != 1 {
            return false, fmt.Errorf("not function requires exactly one argument")
        }
        if v, ok := args[0].(bool); ok {
            return !v, nil
        }
        return false, fmt.Errorf("not function only supports a boolean argument")
    }, nil)

    actions["xor"] = NewAction(func(s *Script, args ...interface{}) (interface{}, error) {
        if len(args) < 2 {
            return false, fmt.Errorf("xor function requires at least two arguments")
        }
        countTrue := 0
        for _, arg := range args {
            if v, ok := arg.(bool); ok {
                if v {
                    countTrue++
                }
            } else {
                return false, fmt.Errorf("xor function only supports boolean arguments")
            }
        }
        return countTrue%2 == 1, nil
    }, nil)

    actions["nand"] = NewAction(func(s *Script, args ...interface{}) (interface{}, error) {
        if len(args) == 0 {
            return false, fmt.Errorf("nand function requires at least one argument")
        }
        for _, arg := range args {
            if v, ok := arg.(bool); ok {
                if !v {
                    return true, nil
                }
            } else {
                return false, fmt.Errorf("nand function only supports boolean arguments")
            }
        }
        return false, nil
    }, nil)

    return &MemoryMap{
        Actions:   actions,
        Variables: variables,
    }
}
