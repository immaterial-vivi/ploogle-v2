package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryResult struct {
	hits        []Chapter
	query       Query
	performance QueryPerformance
}

type QueryPerformance struct {
	numResults int
	// all time in nano seconds
	startTime int64
	endTime   int64
	deltaTime int64
}

type BookResult struct {
	book                Book
	hightlightedChapter Chapter
}

type Query struct {
	query string
}

var searchSqlQuery = `select * from chapters where textsearchable_index_col @@ to_tsquery('english', $1);`

func search(query Query, dbpool *pgxpool.Pool) (*QueryResult, error) {

	log.Println("searching", query.query)

	var result QueryResult

	var queryPerformance QueryPerformance
	result.performance = queryPerformance
	queryPerformance.startTime = int64(time.Now().Nanosecond())

	rows, _ := dbpool.Query(context.Background(),
		"select url from chapters")

	// chapterResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[Chapter])

	// if err != nil {
	// 	return nil, err
	// }

	// for _, chapterResult := range chapterResults {
	// 	fmt.Println(chapterResult)
	// }

	for rows.Next() {
		// var chapterResult Chapter
		var url string

		err := rows.Scan(&url)
		fmt.Println(url)
		// result.hits = append(result.hits, chapterResult)
		// fmt.Println(chapterResult)
		if err != nil {
			return nil, err
		}
	}

	queryPerformance.endTime = int64(time.Now().Nanosecond())
	queryPerformance.deltaTime = queryPerformance.endTime - queryPerformance.startTime

	return &result, nil

}
