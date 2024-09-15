// main.go
package main

import (
	"fmt"
	//"os"
	"time"

	"smuggr.xyz/taskwrappr"
)

func main() {
	//var scriptPath string
	// if len(os.Args) >= 2 {
	// 	scriptPath = os.Args[1]

	// 	if scriptPath == "" {
	// 		fmt.Println("please provide a non empty .tw script path")
	// 		return
	// 	}
	// } else {
	// 	fmt.Println("please provide a valid .tw script path")
	// 	return
	// }

	scriptPath := "../scripts/tests.tw"
	tokenizer := taskwrappr.NewTokenizer(scriptPath)
	startTime := time.Now()
	tokens, err := tokenizer.Tokenize()
	defer func() {
		if err != nil || len(tokens) == 0 {
			return
		}

		endTime := time.Since(startTime)
		tokenCount := len(tokens)
		lineCount := tokenizer.Line

		endTimeMs := float64(endTime) / float64(time.Millisecond)
		perTokenMs := endTimeMs / float64(tokenCount)
		perLineMs := endTimeMs / float64(lineCount)

		fmt.Printf("Tokenize time: %.3fms\n", endTimeMs)
		fmt.Printf("Tokens: %d\n", tokenCount)
		fmt.Printf("Lines: %d\n", lineCount)
		fmt.Printf("Time per token: %.3fms\n", perTokenMs)
		fmt.Printf("Time per line: %.3fms\n", perLineMs)
	}()

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
}
