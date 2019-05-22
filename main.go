package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func SeperateSentences(sentences string) []string {
	sentences = strings.Replace(sentences, " ", "", -1) //remove space
	sentences = strings.Replace(sentences, "​", "", -1) //remove zero space
	return strings.SplitAfter(sentences, "។")           // separate sentences into sentence
}

func GetContent(file *os.File, index int) {
	// Make HTTP request
	response, err := http.Get("https://www.khmerload.com/news/" + strconv.Itoa(index))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		return
	}

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Get the article content
	document.Find(".article-body .article-content div").Each(func(i int, element *goquery.Selection) {
		text := element.Find("p").Text()
		if text != "" {
			for _, sentence := range SeperateSentences(text) {
				// write to file
				fmt.Fprintln(file, sentence)
			}
		}
	})
}

func main() {
	file, err := os.Create("khmer_compus.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	// 101511 is number of article
	for i := 1; i <= 10; i++ {
		fmt.Println(i)
		GetContent(file, i)
	}
}
