package main

import (
	"os"
	"fmt"
	"bufio"
	"net/http"

	"github.com/pehringer/goons/llm"
	"github.com/pehringer/goons/ollama"
)

func main() {
	session, err := llm.Server(&ollama.Server{
		URL:    "http://localhost:11434",
		Client: http.Client{},
	}).Chat("llama3.2:1b")
	if err != nil {
		fmt.Printf("Ollama error: %v", err)
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		reply, err := session.Chat(scanner.Text())
		if err != nil {
			fmt.Printf("Ollama error: %v", err)
			return
		}
		fmt.Println(reply)
	}
}
