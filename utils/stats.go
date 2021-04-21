// Figure out a way to run this file!!!
package utils

import (
	"fmt"
	"net/http"
)

func main() {
	// testing
	MakeRequest("rakaar/facebook-clone")
}

func IsAfterKWoC(timestamp string) {
	// checks if given timestamp is after coding period of KWoC has begun
	// CODING_PERIOD_BEGINS := "2019-04-01T00:00:00Z"
	fmt.Println("timestamp ", timestamp)
}

func MakeRequest(URL string) {
	// make HTTP request
	fmt.Println("url is ", URL)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "token 7920176cccc1c0b6bc2bb6c594873e36db871350")
	res, _ := client.Do(req)

	fmt.Println(res)		

}

func FetchLatestCommits(repo string) {
	fmt.Println("repo is ",repo)
}

func FetchLatestPRs(repo string) {
	fmt.Println("repo is ",repo)
}
