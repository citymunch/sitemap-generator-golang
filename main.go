package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	templater "text/template"
	"time"
	"unicode/utf8"
)

type Config struct {
	ApiEndpoint string
	ApiKey      string
	WebAppUrl   string
	OutputFile  string
}

type FindAuthorisedRestaurantsResponse struct {
	Results []Restaurants
}

type Restaurants struct {
	Restaurant Restaurant
}

type Restaurant struct {
	Id string
}

func main() {
	config := createConfig()
	printConfig(config)

	// Limit to 50000 restaurants because that's the maximum URLs allowed in a sitemap file.
	url := config.ApiEndpoint + "/restaurants/search/authorised-restaurants?limit=50000"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Partner "+config.ApiKey)
	req.Header.Add("X-CM-MaxVersionSupported", "9")
	req.Header.Add("User-Agent", "CityMunch web app sitemap generator")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	data := new(FindAuthorisedRestaurantsResponse)
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		log.Fatal(err)
	}

	template, err := templater.ParseFiles("template.xml")
	if err != nil {
		log.Fatal(err)
	}

	target, err := os.Create(config.OutputFile)
	if err != nil {
		log.Fatal(err)
	}
	err = template.Execute(target, struct {
		WebAppUrl string
		Response  *FindAuthorisedRestaurantsResponse
	}{config.WebAppUrl, data})
	if err != nil {
		log.Fatal(err)
	}

	target.Close()

	fmt.Println("Generated " + config.OutputFile)
}

func createConfig() Config {
	if len(os.Args) < 2 {
		log.Fatal("No output file specified")
	}
	outputFile := os.Args[1]
	if outputFile == "" {
		log.Fatal("No output file specified")
	}

	return Config{
		ApiEndpoint: getConfigVarFromEnv("CM_API_ENDPOINT"),
		ApiKey:      getConfigVarFromEnv("CM_API_KEY"),
		WebAppUrl:   getConfigVarFromEnv("CM_WEB_APP_URL"),
		OutputFile:  outputFile,
	}
}

func getConfigVarFromEnv(envVar string) string {
	value := os.Getenv(envVar)

	if value == "" {
		log.Fatal("Environment variable not set: " + envVar)
	}

	return value
}

func printConfig(config Config) {
	fmt.Println("Running with configuration:")

	fmt.Println("API endpoint:", config.ApiEndpoint)

	// Redact the API key incase the output is redirected to a log file and, for example, sent
	// to a Kibana log server, where the users of the log server shouldn't know the private key.
	redactedApiKey := config.ApiKey[0:10] + strings.Repeat("*", utf8.RuneCountInString(config.ApiKey)-10)
	fmt.Println("API key:", redactedApiKey)

	fmt.Println("Web app URL:", config.WebAppUrl)
}
