package ollama

import (
        "bytes"
        "encoding/json"
        "fmt"
	"io"
	"net/http"
)

type functionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type toolCall struct {
	Function functionCall `json:"function"`
}

type chatMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []toolCall `json:"tool_calls"`
	ToolName  string     `json:"tool_name"`
}

type property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type parameters struct {
	Type       string              `json:"type"`
	Properties map[string]property `json:"properties"`
	Required   []string            `json:"required"`
}

type function struct {
	Name        string     `json:"name"`
        Description string     `json:"description"`
        Parameters  parameters `json:"parameters"`
}

type tool struct {
	Type     string   `json:"type"`
	Function function `json:"function"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Tools	 []tool        `json:"tools"`
	Stream   bool          `json:"stream"`
}

type chatResponse struct {
	Model     string      `json:"model"`
	CreatedAt string      `json:"created_at"`
	Message   chatMessage `json:"message"`
	Done      bool        `json:"done"`
}

func (c *chatRequest) apiChat(message chatMessage) (chatResponse, error) {
	c.Messages = append(c.Messages, message)
	body, err := json.Marshal(c)
	if err != nil {
		return chatResponse{}, fmt.Errorf("error marshaling ollama chat request: %w", err)
	}
	buffer := bytes.NewBuffer(body)
	request, err := http.NewRequest("POST", "http://localhost:11434/api/chat", buffer)
	if err != nil {
		return chatResponse{}, fmt.Errorf("error creating ollama chat request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return chatResponse{}, fmt.Errorf("error sending ollama chat request: %w", err)
	}
	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return chatResponse{}, fmt.Errorf("error reading ollama chat response body: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return chatResponse{}, fmt.Errorf("ollama chat error: %d - %s", response.StatusCode, string(body))
	}
    	result := chatResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
        	return chatResponse{}, fmt.Errorf("error unmarshaling ollama chat response: %w", err)
	}
	c.Messages = append(c.Messages, result.Message)
	return result, nil
}
