package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", handleProcessReceipts)
	r.HandleFunc("/receipts/{id}/points", handleGetPoints)

	log.Println("Server is running on Port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
