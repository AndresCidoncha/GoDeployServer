package main

import ("./PostInfo"
	"encoding/json"
	"os/exec"
	"os"
	"fmt"
	"net/http"
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

func deploy(w http.ResponseWriter, r *http.Request, botName string){
	newCommit := getNewCommit(r)
	fmt.Fprintf(w, "Commit entrante: <%s>\n", newCommit)
	os.Chdir(botsPath + "/" + botName)
	localCommit, _ := exec.Command("getLastCommit").Output()
	fmt.Fprintf(w, "Commit local: <%s>\n", localCommit)
	if string(localCommit) != string(newCommit + "\n") {
		fmt.Fprintf(w, "Redesplegando %s...\n", botName)
		out, err := exec.Command("python", botName+".py", "stop").Output()
		fmt.Fprintf(w, "%s\n", out)
		if err != nil {
			fmt.Fprintf(w, "ERROR STOPPING: %s\n", err)
		}
		out, err = exec.Command("git", "pull").Output()
		fmt.Fprintf(w, "%s\n", out)
		if err != nil {
			fmt.Fprintf(w, "ERROR PULLING: %s\n", err)
		}
		out, err = exec.Command("python", botName+".py", "start").Output()
		fmt.Fprintf(w, "%s\n", out)
		if err != nil {
			fmt.Fprintf(w, "ERROR STARTING: %s\n", err)
		}
	} else {
		fmt.Fprintf(w, "No hay cambios que sincronizar. Ignorando...")
	}
}

func deployAllUsersBot(w http.ResponseWriter, r *http.Request){
	deploy(w, r, "AllUsersBot")
	fmt.Fprintf(w, "Done.")
}


func main() {
	http.HandleFunc("/deploy/AllUsersBot", deployAllUsersBot)
	http.ListenAndServe(":29090", nil)
}
