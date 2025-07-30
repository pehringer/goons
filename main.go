package main

import (
	"fmt"

	"github.com/pehringer/goons/ollama"
)

func main() {
        model := "llama3.2:1b"
	message := "Hello, my name is Jacob."
        fmt.Printf("Prompting Ollama: %s - %s\n", model, message)
	reply, err := ollama.SendPrompt(model, message)
	if err != nil {
		fmt.Printf("Ollama error: %v", err)
	} else {
		fmt.Println(reply)
	}
}
