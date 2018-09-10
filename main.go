package main

import (
    "time"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
	
)

var db *sql.DB
var tpl *template.Template

const (
        DB_USER     = "postgres"
        DB_PASSWORD = "postgres"
        DB_NAME     = "postgres"
    )

func init() {
	var err error
	/*dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
            DB_USER, DB_PASSWORD, DB_NAME)*/
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	/*db, err = sql.Open("postgres", dbinfo)*/
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

type Comment struct {
	Id  int
	Comment  string
	Ip  string
	Date string
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/comments", commentsIndex)
	http.HandleFunc("/comments/create", commentsCreate)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/comments", http.StatusSeeOther)
}

func commentsIndex(w http.ResponseWriter, r *http.Request) {

	/*fmt.Println(db)*/
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}


	rows, err := db.Query("SELECT * FROM comments")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()


	comments := make([]Comment, 0)
	for rows.Next() {
		comment := Comment{}
		err := rows.Scan(&comment.Id, &comment.Comment, &comment.Ip, &comment.Date) // order matters
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tpl.ExecuteTemplate(w, "comments.gohtml", comments)
}


func commentsCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("chao")
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	/*loc, _ := time.LoadLocation("America/Santiago")*/


	// get form values
	comment := Comment{}
	comment.Ip = r.RemoteAddr
	comment.Date = time.Now().Format("2006-01-02 15:04:05")
	comment.Comment = r.FormValue("comment")

	if comment.Comment == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}


	_, err := db.Exec("INSERT INTO comments (comment, ip, date) VALUES ($1, $2, $3)", 
		comment.Comment, comment.Ip, comment.Date)

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}


	tpl.ExecuteTemplate(w, "created.gohtml", comment)

}