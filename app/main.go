// main.go
package main

import (
	"log"

	"smuggr.xyz/taskwrappr"
)

func main() {
	memoryMap := taskwrappr.GetInternals()

	memoryMap.Variables["someStringVar"] = taskwrappr.NewVariable("dupa", taskwrappr.StringType)
	memoryMap.Variables["someCastableVar"] = taskwrappr.NewVariable("-6.9", taskwrappr.StringType)
	memoryMap.Variables["someNotCastableVar"] = taskwrappr.NewVariable("duppa", taskwrappr.StringType)
	memoryMap.Variables["someBoolVar"] = taskwrappr.NewVariable(true, taskwrappr.BooleanType)
	memoryMap.Variables["someIntVar"] = taskwrappr.NewVariable(42, taskwrappr.IntegerType)
	memoryMap.Variables["someFloatVar"] = taskwrappr.NewVariable(3.14, taskwrappr.FloatType)
	memoryMap.Variables["someNegativeVar"] = taskwrappr.NewVariable(-21.37, taskwrappr.FloatType)

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
