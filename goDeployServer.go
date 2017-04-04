package main

import ("./PostInfo"
	"encoding/json"
	"os/exec"
	"os"
	"time"
	"fmt"
	"net/http"
	"strings"
)

const botsPath = "/TelegramBots"

func getNewCommit(r *http.Request) string {
	r.ParseForm()
	var postInfo postInfo.GitHubPostInfo
	err := json.Unmarshal([]byte(r.FormValue("payload")), &postInfo)
	if err != nil {
		fmt.Println("error:", err)
	}
	return postInfo.Before
}


func compareCommits(w http.ResponseWriter, r *http.Request) bool {
	newCommit := strings.TrimSpace(getNewCommit(r))
	fmt.Fprintf(w, "New commit: <%s>\n", newCommit)
	out, err := exec.Command("getLastCommit").Output()
	localCommit := strings.TrimSpace(string(out))
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Fprintf(w, "Local commit: <%s>\n========================================================\n", localCommit)
	return localCommit != newCommit
}


func handlePythonBot(w http.ResponseWriter, command string){
	fmt.Fprintf(w, "========================================================\n-> %s BOT\n", strings.ToUpper(command))
	_, err := exec.Command("/usr/bin/python", "bot.py", command).Output()
	if err != nil {
		fmt.Fprintf(w, "********************************************************\nError on %s bot: %s\n********************************************************\n", command, err)
	}
}


func pullChanges(w http.ResponseWriter){
	out, err := exec.Command("git", "pull").Output()
	fmt.Fprintf(w, "========================================================\n-> PULL:\n%s", out)
	if err != nil {
		fmt.Fprintf(w, "********************************************************\nError on pull: %s\n********************************************************\n", err)
	}
}


func deployPythonBot(w http.ResponseWriter, r *http.Request, botName string){
	os.Chdir(botsPath)
	os.Chdir(botsPath + "/" + botName)
	fmt.Fprintf(w, "PUSH TO DEPLOY: %s\n========================================================\n", botName)
	if compareCommits(w, r) {
		fmt.Fprintf(w, "-> REDEPLOY\n")
		handlePythonBot(w, "stop")
		pullChanges(w)
		time.Sleep(time.Second * 2)
		handlePythonBot(w, "start")
	} else {
		fmt.Fprintf(w, "No changes to sync. Skipping...")
	}
	fmt.Fprintf(w, "========================================================\nDONE")
}


func deployAllUsersBot(w http.ResponseWriter, r *http.Request){
	deployPythonBot(w, r, "AllUsersBot")
}


func main() {
	http.HandleFunc("/deploy/AllUsersBot", deployAllUsersBot)
	http.ListenAndServe(":29090", nil)
}
