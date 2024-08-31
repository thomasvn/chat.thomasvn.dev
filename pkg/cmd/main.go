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
	// localserver()
}

func localrun() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run main.go <question>")
		os.Exit(1)
	}

	docs := chat.ParseFeed("https://thomasvn.dev/feed/")

	llm, _ := openai.New()

	// PREVIOUS
	// stuffQAChain := chains.LoadStuffQA(llm)

	// TEST
	const myStuffQATemplate = `Use the following pieces of context to answer the question at the end. If you're unsure about the answer, provide your best guess based on the available information. Always return a response, even if you're not completely certain. If the context doesn't contain relevant information, use your general knowledge to provide a plausible answer. Clearly state when you're making an educated guess.

	{{.context}}

	Question: {{.question}}
	Helpful Answer:`
	qaPromptSelector := chains.ConditionalPromptSelector{
		DefaultPrompt: prompts.NewPromptTemplate(
			myStuffQATemplate,
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
