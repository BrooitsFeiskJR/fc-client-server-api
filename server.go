package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ApiResponse struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type DolarResponse struct {
	Bid string `json:"bid"`
}

// Consumir a API: https://economia.awesomeapi.com.br/json/last/USD-BRL [ ]
// Retornar o JSON apenas do valor atual do dolar (campo Bid) para o client [ ]
// Registrar no banco de dados SQLite cada cotação recebida [ ]

// Timeouts:
// - chamar a API de cotação do dólar deverá ser de 200ms
// - persistir os dados no banco deverá ser de 10ms

func saveDolarQuotationInDatabase(response ApiResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	db, err := sql.Open("sqlite3", "quotations.db")
	if err != nil {
		return fmt.Errorf("erro ao abrir banco de dados: %w", err)
	}
	defer db.Close()

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS quotations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bid TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	if _, err = db.Exec(createTableQuery); err != nil {
		return fmt.Errorf("erro ao criar tabela: %w", err)
	}

	insertQuery := "INSERT INTO quotations (bid) VALUES (?)"
	if _, err = db.ExecContext(ctx, insertQuery, response.USDBRL.Bid); err != nil {
		return fmt.Errorf("erro ao inserir cotação: %w", err)
	}
	return nil
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo não permitido", http.StatusMethodNotAllowed)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Printf("Error ao criar request: %v", err)
		http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error ao fazer request: %v", err)
		http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error ao ler body: %v", err)
		http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
		return
	}
	var Apiresponse ApiResponse
	if err = json.Unmarshal(body, &Apiresponse); err != nil {
		log.Printf("Error ao desserializar body: %v", err)
		http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
		return
	}
	if err = saveDolarQuotationInDatabase(Apiresponse); err != nil {
		log.Printf("Error ao salvar cotação no banco de dados: %v", err)
		http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
		return
	}
	response := DolarResponse{
		Bid: Apiresponse.USDBRL.Bid,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error ao serializar response: %v", err)
		http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
		return
	}
}

func StartServer() {
	http.HandleFunc("/cotacao", cotacaoHandler)

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
