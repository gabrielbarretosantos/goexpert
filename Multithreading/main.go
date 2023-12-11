package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Result struct {
	Data   map[string]interface{}
	Source string
	Error  bool
}

func main() {
	fmt.Print("Digite o CEP: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	cep := scanner.Text()
	brasilapiChan := make(chan Result)
	viacepChan := make(chan Result)
	go fetchCEP("https://brasilapi.com.br/api/cep/v1/"+cep, brasilapiChan)
	go fetchCEP("http://viacep.com.br/ws/"+cep+"/json/", viacepChan)
	var receveid = false
	for !receveid {
		select {
		case message := <-brasilapiChan:
			receveid = output(message)
		case message := <-viacepChan:
			receveid = output(message)
		case <-time.After(time.Second):
			message := Result{
				Error: true,
			}
			receveid = output(message)
		}

	}

}

func output(message Result) bool {
	if message.Error {
		fmt.Println("ERROR-TIMEOUT")
		return true
	}
	fmt.Printf("API: %s\n", message.Source)
	for key, value := range message.Data {
		fmt.Printf("%s: %v\n", key, value)
	}
	return true
}

func fetchCEP(url string, resultChan chan<- Result) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("err on request on url: %s\n", err)
		resultChan <- Result{}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err on request on body: %s\n", err)
		resultChan <- Result{}
		return
	}

	var addressData map[string]interface{}
	err = json.Unmarshal(body, &addressData)
	if err != nil {
		fmt.Printf("err on unmarshal: %s\n", err)
		resultChan <- Result{}
		return
	}

	resultChan <- Result{Data: addressData, Source: url}
}
