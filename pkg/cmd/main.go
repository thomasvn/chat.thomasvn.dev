package main

import (
	"context"
	"fmt"
	"os"

	"main/pkg/chat"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	localrun()
	localserver()
}

func localrun() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run main.go <question>")
		os.Exit(1)
	}

	docs := chat.ParseFeed("https://thomasvn.dev/feed/")

	// Suitable for a small number of documents.
	llm, _ := openai.New()
	stuffQAChain := chains.LoadStuffQA(llm)
	answer, _ := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": docs,
		"question":        os.Args[1],
	})

	fmt.Println("Question: ", os.Args[1])
	fmt.Println("Answer: ", answer["text"].(string))
}

func localserver() {
	port := "8080"
	funcframework.Start(port)
}
