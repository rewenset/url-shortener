package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"time"
)

const csvFileName = "storage.csv"
const idLength = 5
const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	if _, err := os.Stat(csvFileName); os.IsNotExist(err) {
		f, err := os.Create(csvFileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.Printf("initialized file storage '%s'", csvFileName)
	}

	rand.Seed(time.Now().UnixNano())
}

func generateID() string {
	b := make([]byte, idLength)
	for i := range b {
		b[i] = symbols[rand.Int63()%int64(len(symbols))]
	}
	return string(b)
}

func save(url string) string {
	f, err := os.OpenFile(csvFileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	urlID := generateID()
	err = writer.Write([]string{urlID, url})
	if err != nil {
		panic(err)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		panic(err)
	}
	log.Printf("saved url %s and its id %s to storage", url, urlID)
	return urlID
}

func load(urlID string) string {
	log.Println("loading url from csv")
	f, err := os.Open(csvFileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	for _, r := range records {
		if urlID == r[0] {
			return r[1]
		}
	}

	return ""
}
