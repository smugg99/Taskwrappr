package main

import (
	"fmt"
	"log"
	"time"

	"smuggr.xyz/taskwrappr"
)

func main() {
	memoryMap := taskwrappr.NewMemoryMap()

	memoryMap.Variables["someVar"] = taskwrappr.NewVariable("dupa")
	memoryMap.Variables["someOtherVar"] = taskwrappr.NewVariable(true)

	memoryMap.Actions["navigate"] = taskwrappr.NewAction(func(args ...interface{}) (interface{}, error) {
		url := args[0].(string)
		log.Printf("Navigating to: %s\n", url)

		return true, nil
	})

	memoryMap.Actions["print"] = taskwrappr.NewAction(func(args ...interface{}) (interface{}, error) {
		log.Println(args...)

		return args[0], nil
	})

	memoryMap.Actions["wait"] = taskwrappr.NewAction(func(args ...interface{}) (interface{}, error) {
		log.Println(args...)
		for _, arg := range args {
			log.Printf("Type of arg: %T\n", arg)
		}

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

	memoryMap.Actions["externalfunction"] = taskwrappr.NewAction(func(args ...interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("externalfunction action requires at least 1 argument")
		}

		return args[0], nil
	})

	script, err := taskwrappr.NewScript("../scripts/test.tw")
	if err != nil {
		log.Fatal(err)
	}

	runner := taskwrappr.NewScriptRunner(script, memoryMap)
	success, err := runner.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Script execution finished with status: %v\n", success)
}
