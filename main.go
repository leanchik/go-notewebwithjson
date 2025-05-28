package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type Note struct {
	Title string
	Body  string
}

var notes []Note

func loadNotes() {
	file, err := os.ReadFile("notes.json")
	if err != nil {
		fmt.Println("Файл с заметками не найден", err)
		return
	}
	err = json.Unmarshal(file, &notes)
	if err != nil {
		fmt.Println("Ошибка при чтении заметок: ")
	}
}
func saveNotes() {
	data, err := json.MarshalIndent(notes, "", " ")
	if err != nil {
		fmt.Println("Ошибка при сохранении заметок: ", err)
	}
	err = os.WriteFile("notes.json", data, 0644)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		body := r.FormValue("body")
		if body != "" {
			newNote := Note{
				Title: title,
				Body:  body,
			}
			notes = append(notes, newNote)
			saveNotes()
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	data := struct {
		Author string
		Notes  []Note
	}{
		Author: "leanchik",
		Notes:  notes,
	}

	tmpl.Execute(w, data)
}

func main() {
	loadNotes()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/delete", deleteHandler)
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		indexStr := r.FormValue("index")
		index, err := strconv.Atoi(indexStr)
		if err != nil || index < 0 || index >= len(notes) {
			http.Error(w, "Неверный индекс", http.StatusBadRequest)
			return
		}
		notes = append(notes[:index], notes[index+1:]...)
		saveNotes()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
