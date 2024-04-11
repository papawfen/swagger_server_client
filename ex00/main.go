package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type CandyRequest struct {
	Money int 			`json:"money"`
	CandyType string 	`json:"candyType"`
	CandyCount int 		`json:"candyCount"`
}

type CandyResponse struct {
	Change int 		`json:"change"`
	Thanks string 	`json:"thanks"`
}

type BadResponse struct {
	Error string		`json:"error"`
}

var prices = map[string]int {
	"CE": 10,
	"AA": 15,
	"NT": 17,
	"DE": 21,
	"YR": 23,
}

func checkrequest(request CandyRequest) (BadResponse, int) {
	var badResponse BadResponse
	if _, found := prices[request.CandyType]; found {
		change := request.Money - prices[request.CandyType] * request.CandyCount
		if change < 0 {
			badResponse.Error = "You need " + strconv.Itoa(change * -1) + " more money"
			return badResponse, 402
		}
	} else {
		badResponse.Error = "Candy " + request.CandyType + " not found"
		return badResponse, 400
	}
	if request.Money < 0 {
		badResponse.Error = "Error in input data. Money is negative"
		return badResponse, 400
	} else if request.CandyCount <= 0 {
		badResponse.Error = "Error in input data. Candy count is negative or zero"
		return badResponse, 400
	}
	return badResponse, 201
}

func buyCandy(request CandyRequest) CandyResponse {
	var response CandyResponse
	response.Change = request.Money - prices[request.CandyType] * request.CandyCount
	response.Thanks = "Thank You!"
	return response
}

func buyRequestHandler (w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/buy_candy" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	var request CandyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	badResponse, code := checkrequest(request)
	badResp, _ := json.Marshal(badResponse)
	if code == 400 {
		http.Error(w, "400 wrong input data", http.StatusBadRequest)
		w.Write(badResp)
		return
	} else if code == 402 {
		http.Error(w, "402 not enough money", http.StatusBadRequest)
		w.Write(badResp)
		return
	} else if code == 201 {
		response := buyCandy(request)
		resp, _ := json.Marshal(response)
		w.Write(resp)
	}
}

func main() {
	fmt.Printf("Starting server at port 3333\n")
	http.HandleFunc("/buy_candy", buyRequestHandler)
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Fatal(err)
	}
}
