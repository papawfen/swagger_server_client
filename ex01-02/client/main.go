package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type CandyRequest struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

type CandyResponse struct {
	Change int    `json:"change"`
	Thanks string `json:"thanks"`
}

func readFlags(request CandyRequest) CandyRequest {
	name := flag.String("k", "", "key")
	count := flag.Int("c", 0, "count")
	money := flag.Int("m", 0, "money")
	flag.Parse()
	request.CandyCount = *count
	request.CandyType = *name
	request.Money = *money
	return request
}

func main() {
	var request CandyRequest
	request = readFlags(request)
	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}

	client := GetClient()

	resp, err := client.Post("https://localhost:8080/buy_candy", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	resp_struct := map[string]interface{}{
		"error":  "",
		"change": 0,
		"thanks": "",
	}
	err = json.Unmarshal(res, &resp_struct)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode == 201 {
		fmt.Printf("%s\n", resp_struct["thanks"])
	} else {
		fmt.Printf("%s\n", resp_struct["error"])
	}
}

func GetClient() *http.Client {
	data, err := os.ReadFile("../certs/minica.pem")
	if err != nil {
		log.Fatalln(err)
	}
	cp, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalln(err)
	}
	cp.AppendCertsFromPEM(data)

	cfg := &tls.Config{
		RootCAs: cp,
		GetClientCertificate: ClientCertReqFunc(
			"../certs/cert.pem",
			"../certs/key.pem",
		),
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: cfg,
		},
	}
}

func ClientCertReqFunc(certfile, keyfile string) func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	c, err := getCert(certfile, keyfile)
	return func(certReq *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		if err != nil || certfile == "" {
			log.Fatalln("no certificate provided " + err.Error())
		} else {
			if err != nil {
				log.Fatalf("%v\n", err)
			}
		}
		return &c, nil
	}
}

func getCert(certfile, keyfile string) (c tls.Certificate, err error) {
	if certfile != "" && keyfile != "" {
		c, err = tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			log.Fatalf("error loading key pair: %v\n", err)
		}
	} else {
		err = fmt.Errorf("no certificate provided")
	}
	return
}
