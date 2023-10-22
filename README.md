# Langchain Testing

## Usage

Local

```sh
source .env
go run main.go "What kinds of topics regarding containers has Thomas written about"
```

```txt
Question:  What kinds of topics regarding containers has Thomas written about
Answer:  Thomas has written about topics such as "squashing containers" and "using Homebrew as a package manager for macOS."
```

GCloud Function

```sh
source .env
gcloud functions deploy HelloWorld \
  --runtime go113 \
  --trigger-http \
  --allow-unauthenticated \
  --update-env-vars MY_VARIABLE_NAME=new_value
```

<!-- 
IDEAS
- RAG (retrieval augmented API). Pull contents of all my blog posts. Make it a chat interface.
  - Expose it as an API. Serverless.
  - Make it a chat interface, where you can follow up on questions
  - https://github.com/tmc/langchaingo/blob/main/examples/document-qa-example/document_qa.go
  - https://github.com/tmc/langchaingo/blob/main/examples/chroma-vectorstore-example/chroma_vectorstore_example.go
- Pull contents of all Kubecost codebases & docs
-->

<!-- 
DONE
- Questions are parameterized and passed as CLI Args
- Graceful failure when cloning the repo

-->

<!-- 
// API returned unexpected status code: 400: This model's maximum context length
// is 4097 tokens. However, your messages resulted in 11030 tokens. Please
// reduce the length of the messages.

-->
