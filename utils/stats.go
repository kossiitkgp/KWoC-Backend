package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kwoc20-backend/models"
	"net/http"
	"os"
	"strings"
)

func Testing() string {
	// testing
	// fetch all projects PR
	// FetchLatestPulls("kossiitkgp/kwoc20-backend")
	db := GetDB()
	defer db.Close()

	var projects []models.Project

	err := db.Preload("Mentor").Preload("SecondaryMentor").Find(&projects).Error
	if err != nil {
		fmt.Println("Error in Fetching projects - TODO - Log this")
	}

	fmt.Println("We are testing Pulls ---------------------")
	for _, project := range projects {
		trimmed_repo_link := strings.Replace(project.RepoLink, "https://github.com/", "", 1)
		// FetchLatestPulls(trimmed_repo_link, project.LastPullDate, project.ID)
		FetchLatestCommits(trimmed_repo_link, project.Branch, project.LastCommitSHA, project.ID)
	}

	return "testing for now"
}

func GetExtension(filename string) string {
	split_arr := strings.Split(filename, ".")
	extension := "." + split_arr[len(split_arr)-1]
	return extension
}

func GetLanguagesFromFilenames(filenames []string) []string {
	var languages []string

	json_file, err := os.Open("languages.json")
	if err != nil {
		fmt.Println(err)
	}
	defer json_file.Close()

	var ext2Lang map[string]string
	ext2Lang_bytes, _ := ioutil.ReadAll(json_file)
	_ = json.Unmarshal(ext2Lang_bytes, &ext2Lang)

	// parse the file extensions
	exts := make(map[string]bool)
	for i := range filenames {
		exts[GetExtension(filenames[i])] = true
	}

	// Get extension
	for key := range exts {
		languages = append(languages, ext2Lang[key])
	}

	return languages
}

func IsBeforeKWoC(timestamp string) bool {
	// returns true if the timestamp is before KWoC
	fmt.Println("timestamp ", timestamp)
	KWOC_STARTING_DATE := "2020-05-05T18:30:01Z"
	return timestamp < KWOC_STARTING_DATE
}

func MakeRequest(URL string) (string, string) {
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

func FilterAndSaveCommits(API_URL string, LAST_COMMIT_SHA string, id uint) (bool, string) { // returns true if LATEST commit is found, else false
	db := GetDB()
	defer db.Close()

	res, link_in_headers := MakeRequest(API_URL)
	resBytes := []byte(res)

	var commits []map[string]interface{}
	err := json.Unmarshal(resBytes, &commits)
	if err != nil {
		fmt.Println("err in unmarshal commits ", err)
	}

	for i := range commits {
		// need to check if commit date is after KWoC coding period began or NOT
		commit_info_map := commits[i]["commit"].(map[string]interface{})
		commit_info_author_map := commit_info_map["author"].(map[string]interface{})
		commit_date := commit_info_author_map["date"].(string)
		fmt.Print(commit_date)

		fmt.Print((!IsBeforeKWoC(commit_date)), "hello")
		// For the first page Save the latest commit's SHA
		if (link_in_headers == "" || !strings.Contains(link_in_headers, "rel=\"prev\"")) && (!IsBeforeKWoC(commit_date)) && (i == 0) {
			latest_commit_sha := commits[i]["sha"].(string)
			projects := models.Project{}
			project := &models.Project{
				LastCommitSHA: latest_commit_sha,
			}
			db.First(&projects, id).Select("LastCommitSHA").Updates(project)
			fmt.Println("This is the latest SHA of the project ", latest_commit_sha)
		}
		fmt.Print(LAST_COMMIT_SHA, "1234", commits[i]["sha"])
		if IsBeforeKWoC(commit_date) || commits[i]["sha"] == LAST_COMMIT_SHA {
			return true, ""
		}

		commit_url := commits[i]["html_url"].(string)
		fmt.Println("needed info -> commit URL ", commit_url) // remove this print later
		fmt.Println("SHA is ", commits[i]["sha"])             // remove this print later

		author_data_map, _ := commits[i]["author"].(map[string]interface{})
		student_username := author_data_map["login"].(string)
		fmt.Println("Student username ", student_username)
		// Checking if commit_author is a registered student or not

		student := &models.Student{}
		// CAN BE WRONG
		db.Where("username=?", student_username).First(&student)
		fmt.Print(student)
		if student.ID == 0 {
			continue
		} else {
			// making another API request to get more info about the commit like stats
			api_url, _ := commits[i]["url"].(string)
			res, _ := MakeRequest(api_url)
			resBytes := []byte(res)
			var commit_info map[string]interface{}
			_ = json.Unmarshal(resBytes, &commit_info)

			commit_stats_map, _ := commit_info["stats"].(map[string]interface{})
			lines_added := commit_stats_map["additions"].(float64)
			lines_removed := commit_stats_map["deletions"].(float64)
			fmt.Println("needed_info -> lines-added ", lines_added)     // remove this print later
			fmt.Println("needed_info -> lines-removed ", lines_removed) // remove this print later

			commit_message := commit_info_map["message"].(string)
			fmt.Println("needed info -> message ", commit_message)

			// Fetches the tech on which student worked using file names
			files_arr, _ := commit_info["files"].([]interface{})
			var file_names []string
			for j := range files_arr {
				file_map := files_arr[j].(map[string]interface{})
				file_name := file_map["filename"].(string)
				file_names = append(file_names, file_name)
			}

			languages_worked := GetLanguagesFromFilenames(file_names)
			fmt.Println("languages worked is ", languages_worked)
			// TODO: Update the Languages Worked Field under Student row
			fmt.Print(student_username, "stuent username")

			languages, _ := json.Marshal(languages_worked)
			fmt.Print(languages, "Lnaguages worked on ")
			db.Exec("UPDATE students SET tech_worked = ? WHERE username=?  ", languages, student_username)

			// TODO: Save the commit message in the the DB, the commit model contains
			// URL  : commit_url
			// Message : commit_message
			// LinesAdded : lines_added
			// LinesRemoved: lines_removed
			// SHA : commits[i][sha]

			// project: that will be parameter passed or from the repo name, u can get the object
			// student : you can get the student object based on "student_username"

			Project := models.Project{}
			db.First(&Project, id)

			Student := models.Student{}
			db.Where("username=?", student_username).First(&Student)

			commits := &models.Commits{
				URL:          commit_url,
				Message:      commit_message,
				LinesAdded:   lines_added,
				LinesRemoved: lines_removed,
				SHA:          commits[i]["sha"].(string),

				Project: Project,
				Student: Student,
			}

			db.Create(&commits)

			// Addding the summary stats - increase commit count in Project, and Student
			db.Exec("UPDATE projects SET added_lines=added_lines+? , removed_lines=removed_lines+? WHERE id = ?", lines_added, lines_removed, id)
			db.Exec("UPDATE projects SET commit_count = commit_count + 1 WHERE id = ?", id)

			db.Exec("UPDATE students SET added_lines=added_lines+? , removed_lines=removed_lines+? WHERE username = ?", lines_added, lines_removed, student_username)
			db.Exec("UPDATE students SET commit_count = commit_count + 1 WHERE username = ?", student_username)
			// TODO:
			// Take the Student object and increase the commit_count by 1
			// Take the Project object and increase the commit_count by 1
		}
	}

	// TODO: Update the last commit SHA of the project with commits[0]'s SHA in the FIRST PAGE
	if link_in_headers == "" || !strings.Contains(link_in_headers, "rel=\"next\"") {
		return true, ""
	} else {
		untrimmed_next_url := strings.Split(link_in_headers, ">")[0]
		next_url := strings.TrimLeft(untrimmed_next_url, "<")
		return false, next_url
	}
}

func FetchLatestCommits(repo string, branch string, last_commit_sha string, id uint) { // TODO: Here mostly a project Object will be passed
	fmt.Println("repo is ", repo)

	LAST_COMMIT_SHA := last_commit_sha
	fmt.Print(LAST_COMMIT_SHA) // TODO: need to be fetched from Project object
	LATEST_COMMITS_FETCHED := false
	API_URL := "https://api.github.com/repos/" + repo + "/commits?sha=" + branch
	for !LATEST_COMMITS_FETCHED {
		LATEST_COMMITS_FETCHED, API_URL = FilterAndSaveCommits(API_URL, LAST_COMMIT_SHA, id)
		fmt.Println("API_URL IS -----------------", API_URL)
		fmt.Println("LAST_COMMITS_FETCHED IS -----------------------", LATEST_COMMITS_FETCHED)
	}
}

func FilterAndSavePulls(API_URL string, LAST_PULL_DATE string, project_id uint) (bool, string) {
	db := GetDB()
	defer db.Close()

	res, link_in_headers := MakeRequest(API_URL)
	resBytes := []byte(res)

	var pulls []map[string]interface{}
	err := json.Unmarshal(resBytes, &pulls)
	if err != nil {
		fmt.Println("err in unmarshal commits ", err)
	}

	for i := range pulls {
		pull_date := pulls[i]["created_at"].(string)

		// For the first page Save the latest pull's created date
		if (link_in_headers == "" || !strings.Contains(link_in_headers, "rel=\"prev\"")) && (!IsBeforeKWoC(pull_date)) && (i == 0) {
			latest_pull_date := pulls[i]["created_at"].(string)

			projects := models.Project{}
			project := &models.Project{
				LastPullDate: latest_pull_date,
			}
			db.Preload("Mentor").First(&projects, project_id).Select("LastPullDate").Updates(project)
			fmt.Println("This is the latest pull date of the project ", latest_pull_date)
		}

		if IsBeforeKWoC(pull_date) || pull_date == LAST_PULL_DATE {
			return true, ""
		}

		pull_url := pulls[i]["html_url"].(string)
		title := pulls[i]["title"].(string)
		fmt.Println("pul_url is ", pull_url) // remove this print later
		fmt.Println("Pull ttle is ", title)  // remove this later

		user_info, _ := pulls[i]["user"].(map[string]interface{})
		pr_author := user_info["login"]
		fmt.Println("Author of PR is ", pr_author)

		Project := models.Project{}
		db.First(&Project, project_id)

		fmt.Print(pr_author)
		Student := models.Student{}
		db.Where("username=?", pr_author).First(&Student)

		pull_request := &models.PullRequest{
			URL:     pull_url,
			Title:   title,
			Project: Project,
			Student: Student,
		}

		db.Create(&pull_request)

		fmt.Print(project_id)

		db.Exec("UPDATE projects SET pr_count = pr_count + 1 WHERE id = ?", project_id)

		db.Exec("UPDATE students SET pr_count = pr_count + 1 WHERE username = ?", pr_author)

	}

	if link_in_headers == "" || !strings.Contains(link_in_headers, "rel=\"next\"") {
		return true, ""
	} else {
		untrimmed_next_url := strings.Split(link_in_headers, ">")[0]
		next_url := strings.TrimLeft(untrimmed_next_url, "<")
		return false, next_url
	}
}

func FetchLatestPulls(repo string, last_pull_date string, project_id uint) {
	fmt.Println("repo is ", repo)
	LAST_PULL_DATE := last_pull_date
	LATEST_PULLS_FETCHED := false
	API_URL := "https://api.github.com/repos/" + repo + "/pulls?state=all"
	for !LATEST_PULLS_FETCHED {
		LATEST_PULLS_FETCHED, API_URL = FilterAndSavePulls(API_URL, LAST_PULL_DATE, project_id)
		fmt.Println("API_URL IS ----- ", API_URL)
		fmt.Println("LATEST_PULLS_FETCHED ", LATEST_PULLS_FETCHED)
	}
}
