package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

const (
	REPO_URL = "https://github.com/thomasvn/content"
	DATA_DIR = "./data"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	question := parseArgs()
	fmt.Println("Question:", question)

	llm, err := openai.New()
	if err != nil {
		return err
	}

	err = downloadDocuments()
	if err != nil {
		fmt.Println(err)
	}
	docs, err := loadDocuments()
	if err != nil {
		fmt.Println(err)
	}

	// Stuffs all documents into the prompt of the llm, and returns an answer to
	// the question. Suitable for a small number of documents.
	stuffQAChain := chains.LoadStuffQA(llm)
	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": docs,
		"question":        question,
	})
	if err != nil {
		return err
	}
	fmt.Println(answer)
	return nil
}

func parseArgs() string {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <question>")
		os.Exit(1)
	}

	return args[1]
}

// downloadDocuments clones the REPO_URL to the DATA_DIR.
func downloadDocuments() error {
	cmd := exec.Command("git", "clone", REPO_URL, DATA_DIR)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v, output: %s", err, output)
	}
	return nil
}

// loadDocuments walks through the DATA_DIR and loads all Markdown files into
// the `Document` struct.
func loadDocuments() ([]schema.Document, error) {
	var documents []schema.Document

	err := filepath.Walk(DATA_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() && filepath.Ext(path) == ".md" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			document := schema.Document{PageContent: string(content)}
			documents = append(documents, document)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error loading documents: %v", err)
	}

	return documents, nil
}
