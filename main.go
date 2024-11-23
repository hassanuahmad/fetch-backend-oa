package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/receipts/process", handleProcessReceipts)

	log.Println("Server is running on Port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
