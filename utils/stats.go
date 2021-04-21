// TODO - FIGURE OUT A WAY TO RUN THIS FILE
package utils

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"os"
)


func main() {
	// testing
	FetchLatestCommits("lttkgp/metadata-extractor", "master")
}

func IsBeforeKWoC(timestamp string) bool{
	// returns true if the timestamp is before KWoC
	fmt.Println("timestamp ", timestamp)
	KWOC_STARTING_DATE := "2016-11T11:23:26Z"
	return timestamp < KWOC_STARTING_DATE
}

func MakeRequest(URL string) (string, string){
	// make HTTP request
	fmt.Println("url is ", URL)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_STATS_TOKEN"))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Err is", err)
	}
	defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	response := string(resBody)

	link_in_headers := res.Header.Get("Link")
	return response, link_in_headers 
}

func FilterAndSaveCommits(API_URL string, LAST_COMMIT_SHA string) (bool, string){ // returns true if LATEST commit is found, else false
	res, link_in_headers := MakeRequest(API_URL)
	resBytes := []byte(res)

	var commits []map[string]interface{}
	err := json.Unmarshal(resBytes, &commits)
	if err != nil {
		fmt.Println("err in unmarshal commits ",err)
	}

	for i := range commits {
		// need to check if commit date is after KWoC coding period began or NOT
		commit_info_map := commits[i]["commit"].(map[string]interface{})
		commit_info_author_map := commit_info_map["author"].(map[string]interface{})
		commit_date := commit_info_author_map["date"].(string)
		// if(IsBeforeKWoC(commit_date)){
		// 	continue
		// }
		if(IsBeforeKWoC(commit_date) || commits[i]["sha"] == LAST_COMMIT_SHA) {
			// TODO: Update the LAST COMMIT SHA of the project
			return true, ""
		}
		
		commit_url := commits[i]["html_url"]
		fmt.Println("needed info -> commit URL ", commit_url) // remove this print later
		fmt.Println("SHA is ",commits[i]["sha"]) // remove this print later

		author_data_map, _ := commits[i]["author"].(map[string]interface{}) 
		student_username := author_data_map["login"]
		fmt.Println("Student username ", student_username)
		// Checking if commit_author is a registered student or not
		// TODO: Need to check if student_username is in database or not
		// If in DB, proceed to check more info about commit 
		// if NOT in DB,  "continue" the loop i.e check the next commit
		
		// making another API request to get more info about the commit like stats
		api_url, _ := commits[i]["url"].(string)
		res, _ := MakeRequest(api_url)
		resBytes := []byte(res)
		var commit_info map[string]interface{}
		_ = json.Unmarshal(resBytes, &commit_info)

		commit_stats_map, _ := commit_info["stats"].(map[string]interface{})
		lines_added := commit_stats_map["additions"]
		lines_removed := commit_stats_map["deletions"]
		fmt.Println("needed_info -> lines-added ",lines_added) // remove this print later
		fmt.Println("needed_info -> lines-removed ",lines_removed) // remove this print later
		
		commit_message := commit_info_map["message"]
		fmt.Println("needed info -> message ", commit_message)

		// TODO: Save the commit message in the the DB, the commit model contains 
		// URL  : commit_url
		// Message : commit_message
		// LinesAdded : lines_added
		// LinesRemoved: lines_removed
		
		// project: that will be parameter passed or from the repo name, u can get the object
		// student : you can get the student object based on "student_username"

		//Addding the summary stats - increase commit count in Project, and Student
		// TODO:
		// Take the Student object and increase the commit_count by 1
		// Take the Project object and increase the commit_count by 1
	}

	if(link_in_headers == "" || strings.Contains(link_in_headers, "rel=\"next\"") == false) {
		// TODO: Update the last commit SHA of the project with commits[0]'s SHA
		return true,""
	} else {
		untrimmed_next_url := strings.Split(link_in_headers, ">")[0]
		next_url := strings.TrimLeft(untrimmed_next_url, "<")
		return false, next_url
	}

}

func FetchLatestCommits(repo string, branch string) { // TODO: Here mostly a project Object will be passed
	fmt.Println("repo is ",repo)
	LAST_COMMIT_SHA := "" // TOOD: need to be fetched from Project object 
	LATEST_COMMITS_FETCHED := false
	API_URL := "https://api.github.com/repos/" + repo + "/commits?sha=" + branch
	for LATEST_COMMITS_FETCHED == false {
		LATEST_COMMITS_FETCHED, API_URL = FilterAndSaveCommits(API_URL, LAST_COMMIT_SHA)
	}
}

func FetchLatestPRs(repo string) {
	fmt.Println("repo is ",repo)
}
