// main.go
package main

import (
	"fmt"
	"os"
	"time"

	"smuggr.xyz/taskwrappr"
)

func main() {
	script, err := os.ReadFile("../scripts/tests.tw")
	if err != nil {
		fmt.Println("error reading script file:", err)
		return
	}

	tokenizer := taskwrappr.NewTokenizer(string(script))
	startTime := time.Now()
	tokens, err := tokenizer.Tokenize()
	defer func() {
		if err != nil || len(tokens) == 0 {
			return
		}
		endTime := time.Since(startTime)
		fmt.Println("Tokenize time:", endTime, "per token:", endTime / time.Duration(len(tokens)), "tokens:", len(tokens), "per line:", endTime / time.Duration(tokenizer.Line), "lines:", tokenizer.Line)
	}()

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
}
