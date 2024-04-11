package main

import (
	"log"
	"operations"
	"restapi"
	"github.com/go-openapi/loads"
	"buy_candy"
)

func main() {
	spec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewCandyServerAPI(spec)
	server := restapi.NewServer(api)
	server.EnabledListeners = []string{"https"}
	server.TLSCertificate = "../../certs/cert.pem"
	server.TLSCertificateKey = "../../certs/key.pem"
	server.TLSPort = 8080
	defer server.Shutdown()

	api.BuyCandyHandler = operations.BuyCandyHandlerFunc(
		buy_candy.BuyCandyHandler,
	)

	if err = server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
