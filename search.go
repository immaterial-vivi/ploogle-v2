package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryResult struct {
	Hits        []SearchHit
	Query       string
	Performance QueryPerformance
}

type QueryPerformance struct {
	NumResults int
	// all time in nano seconds
	StartTime int64
	EndTime   int64
	DeltaTime int64
}

type SearchHit struct {
	Book_Url    string
	Chapter_Url string
	Title       string
	Chapter     int
	Author      string
	Excerpt     string
	Rank        float32
	Query       string `json:"-"`
}

type Query struct {
	query  string
	limit  int
	offset int
}

func search(query Query, dbpool *pgxpool.Pool) (*QueryResult, error) {

	log.Println("searching", query.query)

	var result QueryResult

	var queryPerformance QueryPerformance
	queryPerformance.StartTime = int64(time.Now().UnixNano())

	rows, _ := dbpool.Query(
		context.Background(),
		`select
			book_url,
			chapter_url,
			title,
			chapter,
			author,
			ts_headline(
				title || ' ' || chapter_title || ' ' || chapter_text,
				query,
				'StartSel=**,StopSel=**, MaxFragments=3, MinWords=5, MaxWords=10'
			) as excerpt,
			rank,
			query
		from (
				select *, books.url as book_url, ts_rank_cd(textsearchable_index_col, query) as rank from (
						select *, chapters.url as chapter_url
						from chapters,
							websearch_to_tsquery($1) as query
						where textsearchable_index_col @@ query
					)
					join books on book_id = books.id
				order by rank desc
				limit $2
				offset $3
   		);`,
		query.query, query.limit, query.offset)

	queryPerformance.EndTime = int64(time.Now().UnixNano())
	queryPerformance.DeltaTime = queryPerformance.EndTime - queryPerformance.StartTime

	var err error
	result.Hits, err = pgx.CollectRows(rows, pgx.RowToStructByPos[SearchHit])

	if err != nil {
		return nil, err
	}

	var tsquery string
	err = dbpool.QueryRow(context.Background(), "select websearch_to_tsquery('english', $1);", query.query).Scan(&tsquery)

	result.Performance = queryPerformance
	result.Query = tsquery

	log.Println(tsquery)

	return &result, nil

}
