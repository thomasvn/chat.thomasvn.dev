package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

func main() {
	// if err := run(); err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }

	downloadDocuments()
	docs, err := loadDocuments()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", docs)
}

func run() error {
	llm, err := openai.New()
	if err != nil {
		return err
	}

	// We can use LoadStuffQA to create a chain that takes input documents and a question,
	// stuffs all the documents into the prompt of the llm and returns an answer to the
	// question. It is suitable for a small number of documents.
	stuffQAChain := chains.LoadStuffQA(llm)
	docs := []schema.Document{
		{PageContent: "Harrison went to Harvard."},
		{PageContent: "Ankush went to Princeton."},
	}

	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": docs,
		"question":        "Where did Harrison go to collage?",
	})
	if err != nil {
		return err
	}
	fmt.Println(answer)

	// Another option is to use the refine documents chain for question answering. This
	// chain iterates over the input documents one by one, updating an intermediate answer
	// with each iteration. It uses the previous version of the answer and the next document
	// as context. The downside of this type of chain is that it uses multiple llm calls that
	// cant be done in parallel.
	refineQAChain := chains.LoadRefineQA(llm)
	answer, err = chains.Call(context.Background(), refineQAChain, map[string]any{
		"input_documents": docs,
		"question":        "Where did Ankush go to collage?",
	})
	fmt.Println(answer)

	return nil
}

func downloadDocuments() {
	repoURL := "https://github.com/thomasvn/content"
	destination := "./data/"
	cmd := exec.Command("git", "clone", repoURL, destination)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to clone repository: %v, output: %s", err, output)
	}
}

func loadDocuments() ([]schema.Document, error) {
	var documents []schema.Document

	// Path to the directory
	dirPath := "./data"

	// Walk through each file in the directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only consider Markdown regular files (not directories)
		if info.Mode().IsRegular() && filepath.Ext(path) == ".md" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			document := schema.Document{PageContent: string(content)}
			documents = append(documents, document)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error loading documents: %v\n", err)
	}

	return documents, nil
}
