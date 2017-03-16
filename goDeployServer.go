package main

import (//"log"
	"./PostInfo"
	"encoding/json"
	"os/exec"
	"os"
	"fmt"
	"net/http"
)

const botPath = "/TelegramBots"
const AllUsersBotPath = botPath + "/" + "AllUsersBot"

func getNewCommit(r *http.Request) string {
	r.ParseForm()
	var postInfo postInfo.GithubPostInfo
	err := json.Unmarshal([]byte(r.FormValue("payload")), &postInfo)
	if err != nil {
		fmt.Println("error:", err)
	}
	return postInfo.Before
}

func deployAllUsersBot(w http.ResponseWriter, r *http.Request){
	newCommit := getNewCommit(r)
	fmt.Fprintf(w, "Commit entrante: <%s>\n", newCommit)
	os.Chdir(AllUsersBotPath)
	localCommit, _ := exec.Command("getLastCommit").Output()
	fmt.Fprintf(w, "Commit local: <%s>\n", localCommit)
	if string(localCommit) != string(newCommit + "\n") {
		fmt.Fprintf(w, "Deteniendo bot...\n")
		out, _ := exec.Command("git", "pull").Output()
		fmt.Fprintf(w, "Haciendo pull...\n")
		out, _ = exec.Command("git", "pull").Output()
		fmt.Fprintf(w, string(out))
		out, _ = exec.Command("git", "pull").Output()
		fmt.Fprintf(w, "Iniciando bot...\n")
	} else {
		fmt.Fprintf(w, "No hay cambios que sincronizar. Ignorando...")
	}
	fmt.Fprintf(w, "No hay cambios que sincronizar. Ignorando...")
}


func main() {
	http.HandleFunc("/deploy/AllUsersBot", deployAllUsersBot)
	http.ListenAndServe(":29090", nil)
}
