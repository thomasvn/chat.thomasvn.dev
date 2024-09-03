package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/mmcdole/gofeed"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
)

const FeedURL = "https://thomasvn.dev/feed/"
const MyStuffQAPromptTemplate = `You are Thomas, a software engineer and the author of the blog at thomasvn.dev. Respond to the user's question in your authentic voice, using the context provided. Key points about your communication style:

1. Be concise and to the point.
2. Use a friendly, conversational tone.
3. Include technical details when relevant, but explain them clearly.
4. Reference your blog posts or personal experiences when applicable.
5. If you're unsure, provide your best guess based on the available information.
6. Always return a response, even if you're not completely certain.
7. If the context doesn't contain relevant information, use your general knowledge to provide a plausible answer.
8. Clearly state when you're making an educated guess.

Use the following pieces of context to answer the question at the end:
{{.context}}

Question: {{.question}}
Thomas' Response:`

func init() {
	functions.HTTP("Chat", Chat)
}

// GCP Cloud Run Function
//
// Example:
//
//	curl -L 'https://us-west1-thomasvn0.cloudfunctions.net/thomasvn-chat' \
//	-H 'Content-Type: application/json' \
//	-d '{
//	    "message": "THE QUESTION GOES HERE"
//	}'
func Chat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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

	docs := ParseFeed(FeedURL)

	llm, _ := openai.New()
	qaPromptSelector := chains.ConditionalPromptSelector{
		DefaultPrompt: prompts.NewPromptTemplate(
			MyStuffQAPromptTemplate,
			[]string{"context", "question"},
		),
	}
	prompt := qaPromptSelector.GetPrompt(llm)
	llmChain := chains.NewLLMChain(llm, prompt)
	stuffQAChain := chains.NewStuffDocuments(llmChain)
	answer, _ := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": docs,
		"question":        d.Message,
	})

	fmt.Fprintf(w, "%s", answer["text"].(string))
}

func ParseFeed(url string) []schema.Document {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)

	results := []schema.Document{}
	for _, item := range feed.Items {
		converter := md.NewConverter("", true, nil)
		markdown, _ := converter.ConvertString(item.Content)
		metadata := "TITLE: " + item.Title + " LINK: " + item.Link + " UPDATED: " + item.Updated + " PUBLISHED: " + item.Published + "\n\n"
		content := metadata + markdown
		d := schema.Document{
			PageContent: content,
			Metadata:    map[string]any{"title": item.Title, "link": item.Link, "updated": item.Updated, "published": item.Published},
		}
		results = append(results, d)
	}
	return results
}
