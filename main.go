package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	cep            = "01153000"
	brasilApiURL   = "https://brasilapi.com.br/api/cep/v1/" + cep
	viaCepURL      = "http://viacep.com.br/ws/" + cep + "/json/"
	requestTimeout = 1 * time.Second
)

type BrasilAPIResponse struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
}

func fetchFromBrasilAPI(ctx context.Context, ch chan<- string) {
	req, err := http.NewRequestWithContext(ctx, "GET", brasilApiURL, nil)
	if err != nil {
		ch <- fmt.Sprintf("BrasilAPI error: %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- fmt.Sprintf("BrasilAPI error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- fmt.Sprintf("BrasilAPI returned non-200 status: %d", resp.StatusCode)
		return
	}

	var result BrasilAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		ch <- fmt.Sprintf("BrasilAPI error: %v", err)
		return
	}

	ch <- fmt.Sprintf("BrasilAPI response: %+v", result)
}

func fetchFromViaCEP(ctx context.Context, ch chan<- string) {
	req, err := http.NewRequestWithContext(ctx, "GET", viaCepURL, nil)
	if err != nil {
		ch <- fmt.Sprintf("ViaCEP error: %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- fmt.Sprintf("ViaCEP error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- fmt.Sprintf("ViaCEP returned non-200 status: %d", resp.StatusCode)
		return
	}

	var result ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		ch <- fmt.Sprintf("ViaCEP error: %v", err)
		return
	}

	ch <- fmt.Sprintf("ViaCEP response: %+v", result)
}

func main() {
	ch := make(chan string, 2)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	go fetchFromBrasilAPI(ctx, ch)
	go fetchFromViaCEP(ctx, ch)

	select {
	case res := <-ch:
		fmt.Println(res)
	case <-ctx.Done():
		fmt.Println("Timeout: no response received within 1 second")
	}
}
