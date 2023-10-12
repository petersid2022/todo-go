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

var filePath string = "./db/data.json"

func main() {
	http.HandleFunc("/", Handler)
	http.HandleFunc("/todos/", DeleteHandler)
	http.HandleFunc("/mark-done/", ToggleHandler)
	http.HandleFunc("/submit", NewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Handler(w http.ResponseWriter, _ *http.Request) {
	dat, err := os.ReadFile("./index.html")
	if err != nil {
		log.Fatal("ERROR trying to read index.html\n", err)
	}

	var tableHTML string

	database, err := db.GetDatabase(filePath)
	if err != nil {
		log.Fatal("ERROR when fetching db from db.go\n", err)
	}

	for _, item := range database {
		status := "done"
		if item.Status == true {
			status = "not done"
		}
		rowHTML := fmt.Sprintf(
			`<tr>
            <td>%s</td>
            <td>%t</td>
            <td>%s</td>
            <td>
                <button hx-delete="/todos/%d" hx-swap="outerHTML">
                    Delete
                </button>
            </td>
            <td>
            <button hx-put="/mark-done/%d" hx-swap="outerHTML">
                    %s
                </button>
            </td>
        </tr>`,
			item.Name, item.Status, item.Date, item.Id, item.Id, status,
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

	out, _ := strconv.Atoi(id)
	db.DeleteEntry(out, filePath)
}

func ToggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Fatal("ERROR Method not allowed\n", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/mark-done/"):]
	fmt.Printf("Received a PUT request saying that todo ID: %s is DONE\n", id)

	out, _ := strconv.Atoi(id)
	db.ToggleEntry(out, filePath)
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
	db.NewEntry(todo_name, todo_status, todo_date, filePath)
}
