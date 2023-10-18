package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Handler() {
	url := "https://example.com"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalln("HTTP request failed with status code:", resp.StatusCode)
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error during reading body:", err)
		return
	}
	fmt.Println(string(bytes))
}
