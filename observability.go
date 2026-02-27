package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestInfo struct {
	Status    int       `db:"status" json:"status"`
	Method    string    `db:"method" json:"method"`
	Path      string    `db:"path" json:"path"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type QueryInfo struct {
	QueryString string    `db:"query_string" json:"query_string"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Latency     int64     `db:"latency" json:"latency"`
	ResultCount int       `db:"result_count" json:"result_count"`
	TsQuery     string    `db:"ts_query" json:"ts_query"`
	Page        int       `db:"page" json:"page"`
}

type PluckyInfo struct {
	QueryString string    `db:"query_string" json:"query_string"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Latency     int64     `db:"latency" json:"latency"`
	BookUrl     string    `db:"book_url" json:"book_url"`
}

func LogQuery(dbpool *pgxpool.Pool, queryInfo QueryInfo) {
	_, err := dbpool.Exec(context.Background(), `insert into query_log(query_string, latency, result_count, ts_query,page) values ($1,$2,$3,$4,$5)`, queryInfo.QueryString, queryInfo.Latency, queryInfo.ResultCount, queryInfo.TsQuery, queryInfo.Page)
	if err != nil {
		log.Println(err)
	}
}

func LogPlucky(dbpool *pgxpool.Pool, pluckyInfo PluckyInfo) {
	_, err := dbpool.Exec(context.Background(), `insert into plucky_log(query, latency, book_url) values ($1,$2,$3)`, pluckyInfo.QueryString, pluckyInfo.Latency, pluckyInfo.BookUrl)
	if err != nil {
		log.Println(err)
	}
}

func LogRequest(dbpool *pgxpool.Pool, requestInfo RequestInfo) {
	_, err := dbpool.Exec(context.Background(),
		`insert into request_log(status, method, path, user_agent) values ($1, $2, $3, $4)`,
		requestInfo.Status, requestInfo.Method, requestInfo.Path, requestInfo.UserAgent)

	if err != nil {
		log.Println(err)
	}
}
