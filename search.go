package main

import (
	"context"
	"log"
	"math/rand"
	"slices"
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
	Book_Url             string
	Chapter_Url          string
	Title                string
	Chapter              int
	Author               string
	Summary              string
	Excerpt              string
	Rank                 float32
	Total_Count          int `json:"-"`
	Blacklisted          bool
	Blacklist_Reason     string
	Blacklist_Created_at time.Time
	Direct_Title_Match   bool
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

func mergeHits(results []SearchHit, additionalHits []SearchHit) ([]SearchHit, int) {
	removed := 0
	additionalHits = slices.DeleteFunc(additionalHits, func(e SearchHit) bool {
		for _, hit := range results {
			if hit.Book_Url == e.Book_Url {
				removed++
				return true
			}
		}
		return false
	})
	return append(results, additionalHits...), removed
}

func Search(query Query, dbpool *pgxpool.Pool) (*QueryResult, error) {

	log.Println("searching", query.query)

	var result QueryResult

	var queryPerformance QueryPerformance
	queryPerformance.StartTime = time.Now().UnixNano()

	// match exact titles
	rows, err := dbpool.Query(context.Background(),
		`select 
    			url as book_url,
    			'' as chapter_url,
    			title as title,
    			0 as chapter,
    			author as author,
    			summary as summary,
    			'' as excerpt,
    			100 as rank,
			   count(*) over () as total_count,
			   reason is not NULL    as blacklisted,
			   coalesce(reason, '')  as blacklist_reason,
			   coalesce(blacklist.created_at, to_timestamp(0))  as blacklist_created_at,
			   true as direct_title_match
    from books
        left join blacklist on books.id = blacklist.book_id
    where title ilike $1
     limit $2
    offset $3`,
		query.query, query.limit, query.offset)

	exactTitleHits, err := pgx.CollectRows(rows, pgx.RowToStructByPos[SearchHit])

	// match partial titles
	rows, err = dbpool.Query(context.Background(),
		`select 
    			url as book_url,
    			'' as chapter_url,
    			title as title,
    			0 as chapter,
    			author as author,
    			summary as summary,
    			'' as excerpt,
    			100 as rank,
			   count(*) over () as total_count,
			   reason is not NULL    as blacklisted,
			   coalesce(reason, '')  as blacklist_reason,
			   coalesce(blacklist.created_at, to_timestamp(0))  as blacklist_created_at,
			   true as direct_title_match
    from books
        left join blacklist on books.id = blacklist.book_id
    where title ilike '%' || $1 || '%' 
     limit $2
    offset $3`,
		query.query, query.limit, query.offset)

	titleHits, err := pgx.CollectRows(rows, pgx.RowToStructByPos[SearchHit])

	rows, err = dbpool.Query(
		context.Background(),
		`select books.url            as book_url,
				   chapter_ranks.url     as chapter_url,
				   title,
				   chapter_ranks.chapter as chapter,
				   author,
				   books.summary         as summary,
				   ts_headline(
						   title || ' ' || chapter_title || ' ' || chapter_text,
						   query,
						   'StartSel=<em>,StopSel=</em>, MaxFragments=1, MinWords=5, MaxWords=100'
				   )                     as excerpt,
				   rank                  as rank,
				   count(*) over () as total_count,
				   reason is not NULL    as blacklisted,
				   coalesce(reason, '')  as blacklist_reason,
				   coalesce(blacklist.created_at, to_timestamp(0))  as blacklist_created_at,
				   false as direct_title_match
			-- select chapter_ranks.book_id, title, chapter, rank
			from (select distinct on (book_id) book_id,
											   chapter_title,
											   chapter_text,
											   url,
											   chapter,
											   ts_rank_cd('{0.014925373, 0.2, 0.4, 1.0}', textsearchable_index_col, query) as rank,
											   query
	
				  from (select *,
							   chapters.url as chapter_url
						from chapters,
							 websearch_to_tsquery($1) as query
						where textsearchable_index_col @@ query)
				 ) as chapter_ranks
					 join books on chapter_ranks.book_id = books.id
					 left join blacklist on chapter_ranks.book_id = blacklist.book_id
			where blacklist.shadow is not true order by rank desc limit $2 offset $3`,
		query.query, query.limit, query.offset)

	if err != nil {
		log.Println("error getting rows from database:", err)
		return nil, err
	}
	semanticHits, err := pgx.CollectRows(rows, pgx.RowToStructByPos[SearchHit])
	if err != nil {
		log.Println("error getting rows from database:", err)
		return nil, err
	}

	result.Hits, _ = mergeHits(result.Hits, exactTitleHits)
	result.Hits, _ = mergeHits(result.Hits, titleHits)
	result.Hits, _ = mergeHits(result.Hits, semanticHits)

	queryPerformance.EndTime = time.Now().UnixNano()
	queryPerformance.DeltaTime = queryPerformance.EndTime - queryPerformance.StartTime

	var tsquery string
	err = dbpool.QueryRow(context.Background(), "select websearch_to_tsquery('english', $1);", query.query).Scan(&tsquery)

	result.Performance = queryPerformance
	result.TsQuery = tsquery
	result.Query = query.query

	log.Println("searched:", query.query)
	log.Println("num found:", len(result.Hits))

	return &result, nil

}
