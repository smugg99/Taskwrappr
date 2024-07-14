package main

import (
	"fmt"
	"log"

	"smuggr.xyz/taskwrappr"
)

func main() {
	memoryMap := taskwrappr.GetInternals()

	memoryMap.Variables["someVar"] = taskwrappr.NewVariable("dupa")
	memoryMap.Variables["someOtherVar"] = taskwrappr.NewVariable(true)

	memoryMap.Actions["navigate"] = taskwrappr.NewAction(func(s *taskwrappr.Script, args ...interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("navigate action requires at least 1 argument")
		}

		url := args[0].(string)
		fmt.Printf("Navigating to: %s\n", url)

		return url, nil
	})

	memoryMap.Actions["externalfunction"] = taskwrappr.NewAction(func(s *taskwrappr.Script, args ...interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("externalfunction action requires at least 1 argument")
		}

		return args[0], nil
	})

	script, err := taskwrappr.NewScript("../scripts/test.tw", memoryMap)
	if err != nil {
		log.Fatal(err)
	}

	success, err := script.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("script execution finished with status: %v\n", success)
}
