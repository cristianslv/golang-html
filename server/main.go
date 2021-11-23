package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

const (
	Port = ":5000"
)

type Todo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../public/index.html"))

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer sqliteDatabase.Close()

	todos := getTodos(sqliteDatabase)

	items := struct {
		Todos []Todo
	}{
		Todos: todos,
	}

	tmpl.Execute(w, items)
}

type BodyRequest struct {
	Description string `json:"description"`
}

type BodyRequestDelete struct {
	Id int `json:"id"`
}

func todo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}

		var bodyRequest BodyRequest

		json.Unmarshal([]byte(string(body)), &bodyRequest)

		fmt.Printf("Description: %s\n", bodyRequest.Description)

		sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
		defer sqliteDatabase.Close()

		insertTodo(sqliteDatabase, bodyRequest.Description)

		fmt.Fprintln(w, "teste")
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}

		var bodyRequestDelete BodyRequestDelete

		json.Unmarshal([]byte(string(body)), &bodyRequestDelete)

		sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
		defer sqliteDatabase.Close()

		deleteTodo(sqliteDatabase, bodyRequestDelete.Id)
	}
}

func main() {
	startDatabase()

	http.HandleFunc("/", index)
	http.HandleFunc("/todo", todo)
	http.HandleFunc("/delete", delete)

	http.ListenAndServe(Port, nil)
}

func startDatabase() {
	if _, err := os.Stat("./sqlite-database.db"); err == nil {
		log.Println("Database already exists. Starting Server!")
	} else if errors.Is(err, os.ErrNotExist) {
		os.Remove("sqlite-database.db")

		log.Println("Creating sqlite-database.db...")

		file, err := os.Create("sqlite-database.db")

		if err != nil {
			log.Fatal(err.Error())
		}

		file.Close()

		log.Println("sqlite-database.db created")

		sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")

		defer sqliteDatabase.Close()

		createTable(sqliteDatabase)

		insertTodo(sqliteDatabase, "My first todo.")

		todos := getTodos(sqliteDatabase)

		log.Println(todos)
	}
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE todo (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"description" TEXT
	  );`

	statement, err := db.Prepare(createStudentTableSQL)

	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec()
}

func insertTodo(db *sql.DB, description string) {
	insertStudentSQL := `INSERT INTO todo(description) VALUES (?)`
	statement, err := db.Prepare(insertStudentSQL)

	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(description)

	if err != nil {
		log.Fatalln(err.Error())
	}
}

func deleteTodo(db *sql.DB, id int) {
	insertStudentSQL := `DELETE FROM todo WHERE id = ?`
	statement, err := db.Prepare(insertStudentSQL)

	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(id)

	if err != nil {
		log.Fatalln(err.Error())
	}
}

func getTodos(db *sql.DB) []Todo {
	todos := make([]Todo, 0)

	row, err := db.Query("SELECT * FROM todo ORDER BY id")

	if err != nil {
		log.Fatal(err)
	}

	defer row.Close()

	for row.Next() {
		var value1 int
		var value2 string

		row.Scan(&value1, &value2)

		todo := Todo{
			Id:          value1,
			Description: value2,
		}

		todos = append(todos, todo)

	}

	return todos
}
