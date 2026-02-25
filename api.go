package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ResponseData struct {
	Status  string `json:"status"`
	Message any    `json:"message"`
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func getSearchRequest(dbpool *pgxpool.Pool, pageSize int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get query and build query object

		values := r.URL.Query()

		log.Print(values)

		var qString string

		if values.Has("q") {
			qString = values["q"][0]
		} else {
			log.Println("query required")
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		var pString string
		var pageNr int

		if values.Has("p") {
			pString = values["p"][0]
			pageNr, _ = strconv.Atoi(pString)
		}

		var query Query
		query.query = qString
		query.limit = pageSize
		query.offset = pageNr * pageSize

		// send it to db
		result, err := search(query, dbpool)

		// respond
		// jsonMessage, err := json.Marshal(result)
		if err != nil {
			log.Println("error", err)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		response := ResponseData{Status: "success", Message: result}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

func ploogleV2Api(dbpool *pgxpool.Pool, pageSize int) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/", getSearchRequest(dbpool, pageSize))
	return mux, nil
}
