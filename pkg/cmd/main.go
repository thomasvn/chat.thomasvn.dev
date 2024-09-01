package main

import (
	"context"
	"fmt"
	"os"

	chat "thomasvn.dev/chat"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
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

	docs := chat.ParseFeed(chat.FeedURL)

	llm, _ := openai.New()

	qaPromptSelector := chains.ConditionalPromptSelector{
		DefaultPrompt: prompts.NewPromptTemplate(
			chat.MyStuffQAPromptTemplate,
			[]string{"context", "question"},
		),
	}
	prompt := qaPromptSelector.GetPrompt(llm)
	llmChain := chains.NewLLMChain(llm, prompt)
	stuffQAChain := chains.NewStuffDocuments(llmChain)

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
