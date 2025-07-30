package ollama

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
