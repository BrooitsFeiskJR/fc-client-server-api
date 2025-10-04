package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ServerDolarResponse struct {
	Bid string `json:"bid"`
}

func Test() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Printf("Error ao criar request: %v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error ao fazer request: %v", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error ao ler body: %v", err)
	}

	var dolarResponse ServerDolarResponse
	if err = json.Unmarshal(body, &dolarResponse); err != nil {
		log.Printf("Error ao desserializar body: %v", err)
	}
	if err = saveQuotationInTextFile(dolarResponse); err != nil {
		log.Printf("Falha ao salvar cotação no arquivo: %v", err)
	}
	log.Printf("Bid: %v", dolarResponse.Bid)
}

func saveQuotationInTextFile(response ServerDolarResponse) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	dolarFileFormat := fmt.Sprintf("Dólar: %v\n", response.Bid)
	if _, err = f.WriteString(dolarFileFormat); err != nil {
		return err
	}
	return nil
}
