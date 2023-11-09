package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ex, err := fetchApiExchange()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = saveExchangeOnFile(ex["rate"])
	if err != nil {
		log.Fatal(err)
		return
	}
}

func fetchApiExchange() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	var bid map[string]string

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
		return bid, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return bid, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return bid, err
	}
	err = json.Unmarshal(body, &bid)
	if err != nil {
		log.Fatal(err)
		return bid, err
	}
	return bid, nil
}

func saveExchangeOnFile(rate string) error {
	file := createOrOpenFile()
	defer file.Close()
	_, err := file.WriteString("Dólar: " + rate + "\n")
	if err != nil {
		return err
	}
	return nil
}

func createOrOpenFile() *os.File {
	file, err := os.OpenFile("cotacao.txt", os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		file, _ = os.Create("cotacao.txt")
		return file
	}
	return file
}
