package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	db "github.com/petersid2022/todo-go/db"
)

func main() {
	http.HandleFunc("/", Handler)
	http.HandleFunc("/todos/", DeleteHandler)
	http.HandleFunc("/submit", NewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	filePath := "./db/data.json"
	dat, err := os.ReadFile("./index.html")
	if err != nil {
		log.Fatal("ERROR trying to read index.html\n", err)
	}

	tableHTML := ""

	database, err := db.GetDatabase(filePath)
	if err != nil {
		log.Fatal("ERROR when fetching db from db.go\n", err)
	}

	for _, item := range database {
		rowHTML := fmt.Sprintf(
			`<tr>
                <td>%s</td>
                <td>%t</td>
                <td>%s</td>
                <td>
                    <button class="btn btn-danger" hx-delete="/todos/%d">
                        Delete
                    </button>
                </td>
            </tr>`,
			item.Name, item.Status, item.Date, item.Id,
		)
		tableHTML += rowHTML
	}

	out := fmt.Sprintf(string(dat), tableHTML)
	io.WriteString(w, out)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		log.Fatal("ERROR Method not allowed\n", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/todos/"):]
	fmt.Printf("Received a DELETE request for todo ID: %s\n", id)

	filePath := "./db/data.json"
	out, _ := strconv.Atoi(id)
	db.DeleteEntry(out, filePath)
}

func NewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Fatal("ERROR Method not allowed\n", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Received a POST request for a new todo")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	todo_name := r.Form.Get("data")

    if todo_name == "" {
        return
    }

	todo_status := false
	todo_date := time.Now().Format("02/01/2006")
	filePath := "./db/data.json"

	db.NewEntry(todo_name, todo_status, todo_date, filePath)
}
