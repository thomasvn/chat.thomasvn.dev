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
- Google Cloud Function v2. Deploy via API. Restructure code. https://cloud.google.com/functions/docs/create-deploy-http-go
- RAG (retrieval augmented API). Pull contents of all my blog posts. Make it a chat interface.
  - Serverless API can't be cloning the Repo every time. Should I put all my data onto a GCP bucket?
  - Make it a chat interface, where you can follow up on questions
  - https://github.com/tmc/langchaingo/blob/main/examples/document-qa-example/document_qa.go
  - https://github.com/tmc/langchaingo/blob/main/examples/chroma-vectorstore-example/chroma_vectorstore_example.go
- Pull contents of all Kubecost codebases & docs
-->

<!-- 
DONE
- Questions are parameterized and passed as CLI Args
- Graceful failure when cloning the repo
- Expose it as an API via GCP Cloud Functions
-->

<!-- 
Memory limit of 128 MiB exceeded with 131 MiB used. Consider increasing the memory limit
Function execution took 40880 ms, finished with status: 'connection error'
-->

<!-- 
// API returned unexpected status code: 400: This model's maximum context length
// is 4097 tokens. However, your messages resulted in 11030 tokens. Please
// reduce the length of the messages.
-->
