package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

type IdResponse struct {
	Id string `json:"id"`
}

func handleGetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiptId := vars["id"]

	receipt, exists := database[receiptId]
	if !exists {
		http.Error(w, "receipt not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]int{"points": receipt.Points}); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func handleProcessReceipts(w http.ResponseWriter, r *http.Request) {
	var payloadData Receipt
	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	points := 0

	// 1
	retailerNameLen := 0
	for _, letter := range payloadData.Retailer {
		if unicode.IsLetter(letter) || unicode.IsDigit(letter) {
			retailerNameLen++
		}
	}
	points += retailerNameLen

	// 2
	evenOddTotal := payloadData.Total
	totalSplit := strings.Split(evenOddTotal, ".")
	// Assuming all total prices have a decimal point
	if totalSplit[1] == "00" {
		points += 50
	}

	// 3
	multiple25, err := strconv.ParseFloat(evenOddTotal, 64)
	if err != nil {
		fmt.Printf("could not convery multiple25 to an int: %s\n", err)
	}

	z := math.Mod(multiple25, 0.25)
	if z == 0 {
		points += 25
	}

	// 4
	points += (len(payloadData.Items) / 2) * 5

	// 5
	for _, item := range payloadData.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				fmt.Printf("error parsing price: %s\n", err)
				continue
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	// 6
	purchaseDate := payloadData.PurchaseDate
	purchaseDay, err := strconv.Atoi(purchaseDate[len(purchaseDate)-2:])
	if err != nil {
		fmt.Printf("Error converting purchaseDay to integer: %s\n", err)
		return
	}
	dayMod := math.Mod(float64(purchaseDay), 2)
	if int(dayMod) == 1 {
		points += 6
	}

	// 7
	purchaseTime := payloadData.PurchaseTime
	purchaseHour, err := strconv.Atoi(purchaseTime[0:2])
	if err != nil {
		fmt.Printf("Error converting purchaseHour to integer: %s\n", err)
		return
	}
	purchaseMin, err := strconv.Atoi(purchaseTime[len(purchaseTime)-2:])
	if err != nil {
		fmt.Printf("Error converting purchaseDay to integer: %s\n", err)
		return
	}

	if (purchaseHour == 14 && purchaseMin != 00) || purchaseHour == 15 {
		points += 10
	}

	receiptId := "receipt" + strconv.Itoa(len(database)+1)
	payloadData.Points = points
	database[receiptId] = payloadData

	data := &IdResponse{Id: receiptId}
	receiptIdJson, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Could not marshal json: %s\n", err)
		return
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(receiptIdJson)
}
