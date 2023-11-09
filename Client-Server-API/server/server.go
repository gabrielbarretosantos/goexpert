package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type USDBRL struct {
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
}
type Exchange struct {
	USDBRL USDBRL
}

func main() {
	log.Println("Server it's ON :)")
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handleCotation)
	http.ListenAndServe(":8080", mux)
}

func handleCotation(w http.ResponseWriter, r *http.Request) {
	ex, err := fetchExchange()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = saveDB(ex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]string{"rate": ex.USDBRL.Bid}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func connectionDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../exchange.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}

func saveDB(ex *Exchange) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	db, err := connectionDatabase()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS exchanges (id INTEGER PRIMARY KEY, rate varchar(10), timestamp DATETIME DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		log.Fatal(err)
		return err
	}

	stmt, err := db.PrepareContext(ctx, "INSERT INTO exchanges (rate) VALUES (?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, &ex.USDBRL.Bid)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func fetchExchange() (*Exchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var ex Exchange
	if err = json.Unmarshal(body, &ex); err != nil {
		return nil, err
	}
	return &ex, nil
}
