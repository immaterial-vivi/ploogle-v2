package main

import (
	"context"
	"log"
	"os"
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

// type SearchHit struct {
// 	ParsedQuery  string `db:"parsed_query"`
// 	Title        string
// 	Url          string
// 	Excerpt      string
// 	Author       string
// 	ChapterTitle string `db:"chapter_title"`
// 	Chapter      int
// 	Rank         float32
// 	TotalCount   int `db:"total_count" json:"-"`
// }

type SearchHit struct {
	Book_Url string
	Title    string
	Chapter  int
	Query    string `json:"-"`
	Author   string
	Excerpt  string
	Rank     float32
}

type Query struct {
	query  string
	limit  int
	offset int
}

// statement_s_ is a lie but i'm working on it :p
//
// statements prepared:
// - ploogle_websearch_query
func initSearchStatements(dbpool *pgxpool.Pool) error {

	b, err := os.ReadFile("./searchquery.sql")
	if err != nil {
		log.Fatalln(err)
	}

	searchquery := string(b)

	_, err = dbpool.Exec(context.Background(), searchquery)
	return err
}

func search(query Query, dbpool *pgxpool.Pool) (*QueryResult, error) {

	log.Println("searching", query.query)

	var result QueryResult

	var queryPerformance QueryPerformance
	queryPerformance.StartTime = int64(time.Now().UnixNano())

	// rows, _ := dbpool.Query(
	// 	context.Background(),
	// 	fmt.Sprintf(
	// 		"execute ploogle_websearch_query(websearch_to_tsquery('%s'), %d, %d)",
	// 		query.query, query.limit, query.offset),
	// )

	// conn, _ := dbpool.Acquire(context.Background())

	// rows, _ := dbpool.Query(
	// 	context.Background(),
	// 	fmt.Sprintf(`select * from (
	// 		select distinct on (book_id) query, books.title as title, chapters.url as url, ts_headline(title || ' ' || chapter_title || ' ' || chapter_text, query) as excerpt, books.author as author, chapters.chapter_title as chapter_title, chapters.chapter as chapter, ts_rank_cd(textsearchable_index_col, query) as rank, COUNT(*) OVER () AS total_count
	// 		from chapters join books on book_id=books.id, websearch_to_tsquery('%s') as query
	// 	where textsearchable_index_col @@ query)
	// 	order by rank DESC
	// 	limit %d
	// 	offset %d;`, query.query, query.limit, query.offset))

	// rows, _ := dbpool.Query(
	// 	context.Background(),
	// 	`select * from (
	// 		select distinct on (book_id) query, books.title as title, chapters.url as url, ts_headline(title || ' ' || chapter_title || ' ' || chapter_text, query) as excerpt, books.author as author, chapters.chapter_title as chapter_title, chapters.chapter as chapter, ts_rank_cd(textsearchable_index_col, query) as rank, COUNT(*) OVER () AS total_count
	// 		from chapters join books on book_id=books.id, websearch_to_tsquery($1) as query
	// 	where textsearchable_index_col @@ query)
	// 	order by rank DESC
	// 	limit $2
	// 	offset $3;`,
	// 	query.query, query.limit, query.offset)

	rows, _ := dbpool.Query(
		context.Background(),
		`select book_url,
			title,
			chapter,
			query,
			author,
			ts_headline(
				title || ' ' || chapter_title || ' ' || chapter_text,
				query
			) as excerpt, rank
		from (
				select *, books.url as book_url, ts_rank_cd(textsearchable_index_col, query) as rank from (
						select *
						from chapters,
							websearch_to_tsquery($1) as query
						where textsearchable_index_col @@ query
					)
					join books on book_id = books.id
				order by rank
				limit $2 offset $3
    )`,
		query.query, query.limit, query.offset)

	//

	// log.Println(rows.Values())
	queryPerformance.EndTime = int64(time.Now().UnixNano())
	queryPerformance.DeltaTime = queryPerformance.EndTime - queryPerformance.StartTime

	var err error
	result.Hits, err = pgx.CollectRows(rows, pgx.RowToStructByPos[SearchHit])

	if err != nil {
		return nil, err
	}

	// if len(result.Hits) > 0 {
	// 	queryPerformance.NumResults = result.Hits[0].TotalCount
	// }

	var tsquery string
	err = dbpool.QueryRow(context.Background(), "select websearch_to_tsquery('english', $1);", query.query).Scan(&tsquery)

	result.Performance = queryPerformance
	result.Query = tsquery

	log.Println(tsquery)

	return &result, nil

}
