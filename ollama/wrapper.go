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
