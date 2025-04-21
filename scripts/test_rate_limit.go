package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func makeRequest(client *http.Client, url, token string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return
	}

	if token != "" {
		req.Header.Set("API_KEY", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer requisição:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d | Resposta: %s\n", resp.StatusCode, string(body))
}

func main() {
	client := &http.Client{Timeout: 2 * time.Second}
	url := "http://localhost:8080/"

	fmt.Println("Teste por IP (sem token):")
	for i := 1; i <= 15; i++ {
		fmt.Printf("Requisição #%d - ", i)
		makeRequest(client, url, "")
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nTeste com token personalizado:")
	for i := 1; i <= 15; i++ {
		fmt.Printf("Requisição #%d - ", i)
		makeRequest(client, url, "abc123")
		time.Sleep(100 * time.Millisecond)
	}
}
