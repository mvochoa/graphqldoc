package graphqldoc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	b, err := ioutil.ReadFile("./schema.graphql")
	checkError(err)

	query = string(b)
}

// HTTP execute query to the GraphQL endpoint
func HTTP(endpoint string) {
	var response Response

	b, err := json.Marshal(query)
	checkError(err)
	query = `{"query":` + string(b) + `}`

	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(query))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(query)))

	client := &http.Client{}
	resp, err := client.Do(req)
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
