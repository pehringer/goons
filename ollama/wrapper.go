package ollama

import (
        "io"
	"fmt"
	"bytes"
        "net/http"
	"encoding/json"

	"github.com/pehringer/goons/llm"
)

type Session struct {
	url     string
	client  *http.Client
	request *chatRequest
}

type Server struct {
	URL    string
	Client http.Client
}

func (s *Server) Chat(model string) (llm.Session, error) {
	return &Session{
		url:     s.URL,
		client:  &s.Client,
		request: &chatRequest{
			Model:  model,
			Stream: false,
                        Messages: []chatMessage{
                                {
                                        Role:    "system",
                                        Content: `
You are a tool-using assistant.

You can invoke tools by emitting only a single JSON object:
  {"toolname": [input1, input2, ...]}

The user will reply with a single JSON object:
  {"result": [output1, output2, ...]}

Only use a tool if it helps answer the user's question.
Do not ask whether a tool can be used, just use it.
Only use tools that are explicitly listed.
Never simulate or guess a tool result.
Do not describe, explain, or narrate tool behavior.
Do not provide a final answer until the tool result has been returned.
Use natural language when no tool is needed.

Available tools:
  {"addition": [a, b]}  ->  {"result": [a + b]}
  {"subtraction": [a, b]}  ->  {"result": [a - b]}

Example tool usage:

User: What is 5 + 3?
Assistant: {"addition": [5, 3]}
User: {"result": [8]}
Assistant: The answer is 8.

User: What is 10 minus 4?
Assistant: {"subtraction": [10, 4]}
User: {"result": [6]}
Assistant: The result is 6.

User: Hello!
Assistant: Hi there!

User: 7 - 20
Assistant: {"subtraction": [20, 7]}
User: {"result": [13]}
Assistant: 13

User: -2 plus 5?
Assistant: {"addition": [-2, 5]}
User: {"result": [3]}
Assistant: Equals 3.
`,
                                },
                        },

		},
	}, nil
}

func (s *Session) Chat(message string) (string, error) {
	s.request.Messages = append(s.request.Messages, chatMessage{
		Role:    "user",
		Content: message,
	})
	body, err := json.Marshal(s.request)
	if err != nil {
		return "", fmt.Errorf("error marshaling ollama chat request: %w", err)
	}
	request, err := http.NewRequest("POST", s.url+"/api/chat", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error creating ollama chat request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := s.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error sending ollama chat request: %w", err)
	}
	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading ollama chat response body: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama chat error: %d - %s", response.StatusCode, string(body))
	}
	result := &chatResponse{}
	if err := json.Unmarshal(body, result); err != nil {
		return "", fmt.Errorf("error unmarshaling ollama chat response: %w", err)
	}
	s.request.Messages = append(s.request.Messages, result.Message)
	return result.Message.Content, nil
}
