package ollama

func SendPrompt(model, message string) (string, error) {
	c := &chatRequest{
		Model: model,
		Stream: false,
	}
	m := chatMessage{
		Role:    "user",
		Content: message,
	}
	r, err := c.apiChat(m)
	return r.Message.Content, err
}
