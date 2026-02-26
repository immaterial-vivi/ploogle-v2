package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryResult struct {
	Hits        []SearchHit
	TsQuery     string
	Query       string
	Performance QueryPerformance
	Page        Page
}

type Page struct {
	Results int
	Pages   int
	Page    int
	PerPage int
}

type QueryPerformance struct {
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
	Total_Count int    `json:"-"`
}

type Query struct {
	query  string
	limit  int
	offset int
}

func ImFeelingPlucky(query string, dbpool *pgxpool.Pool) (SearchHit, error) {

	queryObj := Query{query: query, limit: 100, offset: 0}

	queryResult, _ := Search(queryObj, dbpool)

	hit := queryResult.Hits[rand.Intn(len(queryResult.Hits))]

	return hit, nil
}

func Search(query Query, dbpool *pgxpool.Pool) (*QueryResult, error) {

	log.Println("searching", query.query)

	var result QueryResult

	var queryPerformance QueryPerformance
	queryPerformance.StartTime = int64(time.Now().UnixNano())

	rows, _ := dbpool.Query(
		context.Background(),
		`select book_url,
			chapter_url,
			title,
			chapter,
			author,
			ts_headline(
				title || ' ' || chapter_title || ' ' || chapter_text,
				query,
				'StartSel=**,StopSel=**, MaxFragments=1, MinWords=5, MaxWords=100'
			) as excerpt,
			rank,
			query,
			count
		from (
				select *,
					count(*) over () as count
				from (
						select distinct on (book_id, rank) *
						from (
								select *,
									books.url as book_url,
									ts_rank_cd('{0.014925373, 0.2, 0.4, 1.0}', textsearchable_index_col, query) as rank
								from (
										select *,
											chapters.url as chapter_url
										from chapters,
											websearch_to_tsquery($1) as query
										where textsearchable_index_col @@ query
									)
									join books on book_id = books.id
								order by rank desc
							) order by rank desc
					)
				limit $2 offset $3
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
	result.TsQuery = tsquery
	result.Query = query.query
	log.Println(tsquery)

	return &result, nil

}
