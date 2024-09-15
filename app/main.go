// main.go
package main

import (
	"fmt"
	//"os"
	"time"

	"smuggr.xyz/taskwrappr"
)

func printTokenizeStats(tokens []taskwrappr.Token, startTime time.Time) {
	if len(tokens) == 0 {
		return
	}

	endTime := time.Since(startTime)
	tokenCount := len(tokens)
	lineCount := tokens[len(tokens)-1].Line()

	endTimeMs := float64(endTime) / float64(time.Millisecond)
	perTokenMs := endTimeMs / float64(tokenCount)
	perLineMs := endTimeMs / float64(lineCount)

	fmt.Printf("\nTokenize time: %.3fms\n", endTimeMs)
	fmt.Printf("Tokens: %d\n", tokenCount)
	fmt.Printf("Lines: %d\n", lineCount)
	fmt.Printf("Time per token: %.3fms\n", perTokenMs)
	fmt.Printf("Time per line: %.3fms\n", perLineMs)
}

func printParseStats(nodes []taskwrappr.Node, startTime time.Time) {
	if len(nodes) == 0 {
		return
	}

	endTime := time.Since(startTime)
	nodeCount := len(nodes)

	endTimeMs := float64(endTime) / float64(time.Millisecond)
	perNodeMs := endTimeMs / float64(nodeCount)

	fmt.Printf("\nParse time: %.3fms\n", endTimeMs)
	fmt.Printf("Nodes: %d\n", nodeCount)
	fmt.Printf("Time per node: %.3fms\n", perNodeMs)
}

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

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nTokens: %d\n", len(tokens))
	for _, token := range tokens {
		fmt.Printf("[%s:%d:%d] %s\n", scriptPath, token.Line(), token.IndexSinceLine(), token)
	}
	fmt.Printf("\n")

	parser := taskwrappr.NewParser(tokens, scriptPath)
	_parserStartTime := time.Now()
	nodes, err := parser.Parse()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nNodes: %d\n", len(nodes))
	for _, node := range nodes {
		fmt.Printf("[%s:%d:%d] %s\n", scriptPath, node.Line(), node.IndexSinceLine(), node)
	}
	fmt.Printf("\n")

	printTokenizeStats(tokens, _tokenizerStartTime)
	printParseStats(nodes, _parserStartTime)
}
