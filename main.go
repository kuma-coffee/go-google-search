package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

var (
	apiKey string
	cx     string
	query  string
)

func customSearch(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	cx := os.Getenv("CX")
	query := os.Getenv("QUERY")

	client := &http.Client{Transport: &transport.APIKey{Key: apiKey}}

	svc, err := customsearch.New(client)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := svc.Cse.List().Cx(cx).Q(query).Do()
	if err != nil {
		log.Fatal(err)
	}

	var searchResult []string

	for i, result := range resp.Items {
		searchResult = append(searchResult, fmt.Sprintf("#%d: %s\n", i+1, result.Title), fmt.Sprintf("\t%s\n", result.Snippet), fmt.Sprintf("\t%s\n", result.Link))
	}

	c.JSON(200, searchResult)
}

func main() {
	r := gin.Default()

	r.GET("/", customSearch)

	r.Run()
}
