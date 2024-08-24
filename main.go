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

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/mmcdole/gofeed"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

func init() {
	functions.HTTP("Chat", Chat)
}

// GCP Cloud Run Function
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

	docs := parseFeed("https://thomasvn.dev/feed/")

	llm, _ := openai.New()
	stuffQAChain := chains.LoadStuffQA(llm)
	answer, _ := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": docs,
		"question":        d.Message,
	})

	fmt.Fprintf(w, "Hello, %s!", html.EscapeString(answer["text"].(string)))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run main.go <question>")
		os.Exit(1)
	}

	docs := parseFeed("https://thomasvn.dev/feed/")

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

func parseFeed(url string) []schema.Document {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)

	results := []schema.Document{}
	for _, item := range feed.Items {
		content := "TITLE: " + item.Title + "\n\n" + html.UnescapeString(item.Content)
		d := schema.Document{
			PageContent: content,
			Metadata:    map[string]any{"title": item.Title, "link": item.Link, "updated": item.Updated, "published": item.Published},
		}
		results = append(results, d)
	}
	return results
}
