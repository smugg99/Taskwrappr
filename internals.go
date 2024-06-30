// internals.go
package taskwrappr

import (
	"fmt"
	"time"
)

func GetInternals() (map[string]*Action, map[string]*Variable) {
	actions := make(map[string]*Action)
	variables := make(map[string]*Variable)

	actions["print"] = NewAction(func(s *ScriptRunner, args ...interface{}) (interface{}, error) {
		fmt.Println(args...)

		return args[0], nil
	})

	actions["wait"] = NewAction(func(s *ScriptRunner, args ...interface{}) (interface{}, error) {
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
	})

	actions["if"] = NewAction(func(s *ScriptRunner, args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return false, fmt.Errorf("if function requires at least one argument")
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
		}

		for _, arg := range args {
			switch v := arg.(type) {
			case bool:
				if !v {
					return false, nil
				}
			case string:
				if v == "" {
					return false, nil
				}
			default:
				return false, fmt.Errorf("unsupported argument type: %T", arg)
			}
		}

		return true, nil
	})

	actions["else"] = NewAction(func(s *ScriptRunner, args ...interface{}) (interface{}, error) {
		return !s.flags.LastActionSuccess, nil
	})

	return actions, variables
}