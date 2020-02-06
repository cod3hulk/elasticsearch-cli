package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Request struct {
}

func main() {
	newClient := &http.Client{
		Timeout: time.Second * 10,
	}

	host := os.Args[1]
	query := os.Args[2]

	now := time.Now()
	start := now.Add(-4 * time.Hour).Format("2006-01-02 15:04:05")
	end := now.Format("2006-01-02 15:04:05")

	search := fmt.Sprintf(`
		{
			"query": {
				"bool": {
					"must": [{
						"query_string": {
							"query": "%s"
						}
					}, {
						"range": {
							"@timestamp": {
								"gte": "%s",
								"lte": "%s",
								"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
							}
						}
					}]
				}
			}
		}
	`, query, start, end)

	req, _ := http.NewRequest("GET", host, strings.NewReader(search))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := newClient.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))

}
