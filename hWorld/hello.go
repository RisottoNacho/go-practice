package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}
type ArticleRecord struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

const (
	host     = "127.0.0.2"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "postgres"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func openConection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db

}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	enableCORS(myRouter)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles).Methods("GET")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET")
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	db := openConection()

	vars := mux.Vars(r)
	id := vars["id"]

	var article ArticleRecord
	err := db.QueryRow(`SELECT id, name FROM cliets WHERE id = ?`, id).Scan(
		&article.Id,
		&article.Title,
		&article.Desc,
		&article.Content,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	defer db.Close()
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(article)

}
func createNewArticle(w http.ResponseWriter, r *http.Request) {
	db := openConection()

	var newArticle Article
	err := json.NewDecoder(r.Body).Decode(&newArticle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sqlStatement := "INSERT INTO article (title,description,content_) VALUES ($1,$2,$3)"
	_, err = db.Exec(sqlStatement, newArticle.Title, newArticle.Desc, newArticle.Content)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newArticle)
	defer db.Close()

}
func updateArticle(w http.ResponseWriter, r *http.Request) {
	// reqBody, _ := ioutil.ReadAll(r.Body)
	// vars := mux.Vars(r)
	// key := vars["id"]
	// var updatedArticle Article

	// json.Unmarshal([]byte(reqBody), &updatedArticle)

	// for i, article := range Articles {
	// 	if article.Id == key {
	// 		updatedArticle.Id = article.Id
	// 		Articles[i] = updatedArticle
	// 	}
	// }
	// json.NewEncoder(w).Encode(Articles)
}
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// key := vars["id"]
	// for i, article := range Articles {
	// 	if article.Id == key {
	// 		Articles = append(Articles[:i], Articles[i+1:]...)
	// 		break
	// 	}
	// }
	// json.NewEncoder(w).Encode(Articles)
}

func main() {
	// Articles = []Article{
	// 	{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
	// 	{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	// }
	handleRequests()
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	db := openConection()

	var articles []ArticleRecord

	rows, err := db.Query(`SELECT * FROM article`)

	if err != nil {
		//some error handling
		return
	}

	for rows.Next() {
		article := ArticleRecord{}

		err = rows.Scan(
			&article.Title,
			&article.Desc,
			&article.Content,
		)

		if err != nil {
			panic(err)
		}

		articles = append(articles, article)
	}
	defer db.Close()
	defer rows.Close()

	json.NewEncoder(w).Encode(articles)
}

func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			// Just put some headers to allow CORS...
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			// and call next handler!
			next.ServeHTTP(w, req)
		})
}
