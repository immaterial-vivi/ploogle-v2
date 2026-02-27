package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChartResponse struct {
	StatCards []StatCard
	ChartData ChartData `json:"chart_data"`
}

type StatCard struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []DataSet `json:"datasets"`
}

type DataSet struct {
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	BackgroundColor string    `json:"chart_data"`
}

type BlacklistEntryDTO struct { //i know that that's not quite the right word but idk
	BlacklistEntryId int       `json:"blacklist_entry_id"`
	BookId           int       `json:"book_id"`
	BookTitle        string    `json:"book_title"`
	Reason           string    `json:"reason"`
	BlacklistedAt    time.Time `json:"blacklisted_at"`
}

func getBlacklist(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, _ := dbpool.Query(context.Background(), "select blacklist.id, book_id, title, reason, blacklist.created_at from blacklist join books on books.id = blacklist.book_id")

		results, err := pgx.CollectRows(rows, pgx.RowToStructByPos[BlacklistEntryDTO])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)

		}
		w.WriteHeader(http.StatusOK)
	}
}

func postAddToBlacklist(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		url := r.Form.Get("url")
		if url == "" {
			http.Error(w, "url is required", http.StatusBadRequest)
			return
		}

		reason := r.Form.Get("reason")
		if reason == "" {
			http.Error(w, "reason is required", http.StatusBadRequest)
			return
		}

		var bookId int
		err := dbpool.QueryRow(context.Background(), "select id from books where url like $1", url).Scan(&bookId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		_, err = dbpool.Exec(context.Background(), `insert into blacklist(book_id, reason)  values ($1,$2)`, bookId, reason)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		http.Redirect(w, r, "/dash/blacklist.html", http.StatusFound)
	}
}

func postRemoveFromBlacklist(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		blacklistEntryTarget := struct {
			BlacklistId int `json:"blacklist_id"`
		}{}

		d := json.NewDecoder(r.Body)

		err := d.Decode(&blacklistEntryTarget)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		log.Println(blacklistEntryTarget.BlacklistId)

		_, err = dbpool.Exec(context.Background(), `delete from blacklist where id = $1`, blacklistEntryTarget.BlacklistId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func getRequestStats(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func getPluckyStats(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func getLatencyStats(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func getRecentSearches(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func DashboardApi(dbpool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./dash"))

	mux.Handle("/", fileServer)
	mux.HandleFunc("GET /blacklist", getBlacklist(dbpool))
	mux.HandleFunc("POST /blacklist/add", postAddToBlacklist(dbpool))
	mux.HandleFunc("DELETE /blacklist/remove", postRemoveFromBlacklist(dbpool))

	mux.HandleFunc("/requestStats", getRequestStats(dbpool))
	mux.HandleFunc("/pluckyStats", getPluckyStats(dbpool))
	mux.HandleFunc("/latencyStats", getLatencyStats(dbpool))
	mux.HandleFunc("/recentSearches", getRecentSearches(dbpool))
	return mux
}
