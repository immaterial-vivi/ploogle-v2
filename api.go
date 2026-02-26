package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
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

func getPluckyRequest(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		var qString string
		var url string

		if values.Has("q") {
			qString = values["q"][0]
		}

		if qString == "" {
			url, _ = GetRandomBookUrl(dbpool)
		} else {
			hit, _ := ImFeelingPlucky(qString, dbpool)
			url = hit.Book_Url
		}

		response := ResponseData{Status: "success", Message: url}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

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
		pageNr := 1

		if values.Has("p") {
			pString = values["p"][0]
			pInt, _ := strconv.Atoi(pString)
			pageNr = max(pInt, 1)
		}

		var query Query
		query.query = qString
		query.limit = pageSize
		query.offset = (pageNr - 1) * pageSize

		// send it to db
		result, err := Search(query, dbpool)

		// insert
		if len(result.Hits) > 0 {
			result.Page.Results = result.Hits[0].Total_Count
		}

		result.Page.Page = pageNr
		result.Page.Pages = int(math.Ceil(float64(result.Page.Results) / float64(pageSize)))
		result.Page.PerPage = pageSize

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

func PloogleV2Api(dbpool *pgxpool.Pool, pageSize int) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/search", getSearchRequest(dbpool, pageSize))
	mux.HandleFunc("/plucky", getPluckyRequest(dbpool))

	return mux, nil
}
