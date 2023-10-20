package handler

import (
	"os"

	"github.com/pektezol/gobrr/src/pkg/html"
)

func Handler() {
	// url := "https://example.com"
	// resp, err := http.Get(url)
	// if err != nil {
	// 	log.Fatalln("Error making HTTP request:", err)
	// 	return
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	log.Fatalln("HTTP request failed with status code:", resp.StatusCode)
	// 	return
	// }
	file, _ := os.Open("test/html/test_cbracco.html")
	tokenizer := html.NewLexer(file)
	tokenizer.Read()
}
