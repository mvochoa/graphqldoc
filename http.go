package graphqldoc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Data struct {
	Schema Schema `json:"__schema"`
}

// Response estructura de una respuesta HTTP
type Response struct {
	Data   Data                     `json:"data"`
	Errors []map[string]interface{} `json:"errors"`
}

var query string

func init() {
	b, err := Asset("template/schema.graphql")
	checkError(err)

	query = string(b)
}

// HTTP execute query to the GraphQL endpoint
func HTTP(endpoint string) {
	var response Response

	resp, err := http.PostForm(endpoint,
		url.Values{"query": {query}})

	checkError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	json.Unmarshal(body, &response)
	if len(response.Errors) != 0 {
		log.Fatal(response.Errors)
	}

	generateDocs(response.Data.Schema)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
