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
	_tokenizerStartTime := time.Now()
	tokens, err := tokenizer.Tokenize()
	defer func() {
		if err != nil || len(tokens) == 0 {
			return
		}

		endTime := time.Since(_tokenizerStartTime)
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
		fmt.Printf("[%s:%d:%d] %s\n", scriptPath, token.Line(), token.IndexSinceLine(), token)
	}

	parser := taskwrappr.NewParser(tokens, scriptPath)
	_parserStartTime := time.Now()
	nodes, err := parser.Parse()
	defer func() {
		if err != nil || len(nodes) == 0 {
			return
		}

		endTime := time.Since(_parserStartTime)
		nodeCount := len(nodes)

		endTimeMs := float64(endTime) / float64(time.Millisecond)
		perNodeMs := endTimeMs / float64(nodeCount)

		fmt.Printf("Parse time: %.3fms\n", endTimeMs)
		fmt.Printf("Nodes: %d\n", nodeCount)
		fmt.Printf("Time per node: %.3fms\n", perNodeMs)
	}()

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, node := range nodes {
		fmt.Printf("[%s:%d:%d] %s\n", scriptPath, node.Line(), node.IndexSinceLine(), node)
	}
}
