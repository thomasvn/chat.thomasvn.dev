package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

const (
	REPO_URL = "github.com/thomasvn/content.git"
	DATA_DIR = "/tmp/data"
)

func main() {
	question, err := parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	answer, err := run(question)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Question: ", question)
	fmt.Println("Answer: ", answer)
}

func run(question string) (string, error) {
	llm, err := openai.New()
	if err != nil {
		return "", err
	}

	// Clones docs from the `REPO_URL`, then loads them into `[]schema.Document`
	// for the llm.
	err = downloadDocuments()
	if err != nil {
		return "", err
	}
	docs, err := loadDocuments()
	if err != nil {
		return "", err
	}

	// Stuffs all documents into the prompt of the llm, and returns an answer to
	// the question. Suitable for a small number of documents.
	stuffQAChain := chains.LoadStuffQA(llm)
	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": docs,
		"question":        question,
	})
	if err != nil {
		return "", err
	}

	// `answer["text"]` of type `any` needs to be converted to a `string`.
	answerString, ok := answer["text"].(string)
	if !ok {
		return "", fmt.Errorf("failed to convert answer to string")
	}
	return answerString, nil
}

// parseArgs parses the command line arguments. It expects the question to the
// llm as the sole argument.
func parseArgs() (string, error) {
	args := os.Args
	if len(args) != 2 {
		return "", fmt.Errorf("usage: go run main.go <question>")
	}
	return args[1], nil
}

// downloadDocuments clones the REPO_URL to the DATA_DIR.
func downloadDocuments() error {
	if _, err := os.Stat(DATA_DIR); !os.IsNotExist(err) {
		fmt.Println("Directory already exists. Skipping clone.")
		return nil
	}

	if os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN") == "" {
		return fmt.Errorf("GITHUB_PERSONAL_ACCESS_TOKEN environment variable not set")
	}

	var repoUrlWithToken = "https://x-access-token:" + os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN") + "@" + REPO_URL
	cmd := exec.Command("git", "clone", repoUrlWithToken, DATA_DIR)
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

// GCP Cloud Function
//
// Example:
//
//	curl -L 'https://us-west1-thomasvn0.cloudfunctions.net/CLOUD_FUNCTION_NAME' \
//	-H 'Content-Type: application/json' \
//	-d '{
//	    "message": "THE QUESTION GOES HERE"
//	}'
func Chat(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		switch err {
		case io.EOF:
			fmt.Fprint(w, "Example usage: curl -L 'https://us-west1-thomasvn0.cloudfunctions.net/CLOUD_FUNCTION_NAME' -H 'Content-Type: application/json' -d '{\"message\": \"THE QUESTION GOES HERE\"}'")
			return
		default:
			log.Printf("json.NewDecoder: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	result, err := run(d.Message)
	if err != nil {
		log.Printf("run: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, html.EscapeString(result))
}
