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
	Complete_Title_Match bool
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

func mergeHits(results []SearchHit, additionalHits []SearchHit, acc int) ([]SearchHit, int) {
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

	resultsLen := 0
	if len(results) > 0 && acc == 0 {
		resultsLen = results[0].Total_Count
	}

	additionalHitsLen := 0
	if len(additionalHits) > 0 {
		additionalHitsLen = additionalHits[0].Total_Count
	}

	return append(results, additionalHits...), resultsLen + additionalHitsLen - removed
}

func Search(query Query, dbpool *pgxpool.Pool) (*QueryResult, error) {

	log.Println("searching", query.query)

	var result QueryResult

	var queryPerformance QueryPerformance
	queryPerformance.StartTime = time.Now().UnixNano()

	var countExactTitleHits, countTitleHits, totalHits int
	limit := query.limit
	offset := query.offset

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
			   true as direct_title_match,
				true as complete_title_match
    from books
        left join blacklist on books.id = blacklist.book_id
    where title ilike $1
     limit $2
    offset $3`,
		query.query, limit, offset)

	exactTitleHits, err := pgx.CollectRows(rows, pgx.RowToStructByPos[SearchHit])
	result.Hits, countExactTitleHits = mergeHits(result.Hits, exactTitleHits, 0)

	limit = max(0, limit-len(result.Hits))
	offset = max(0, offset-len(result.Hits))

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
			   true as direct_title_match,
				false as complete_title_match

    from books
        left join blacklist on books.id = blacklist.book_id
    where title ilike '%' || $1 || '%' 
     limit $2
    offset $3`,
		query.query, limit, offset)

	titleHits, err := pgx.CollectRows(rows, pgx.RowToStructByPos[SearchHit])
	result.Hits, countTitleHits = mergeHits(result.Hits, titleHits, countExactTitleHits)

	limit = max(0, limit-len(result.Hits))
	offset = max(0, offset-len(result.Hits))

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
				   false as direct_title_match,
				false as complete_title_match

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
		query.query, limit, offset)

	if err != nil {
		log.Println("error getting rows from database:", err)
		return nil, err
	}
	semanticHits, err := pgx.CollectRows(rows, pgx.RowToStructByPos[SearchHit])
	if err != nil {
		log.Println("error getting rows from database:", err)
		return nil, err
	}

	result.Hits, totalHits = mergeHits(result.Hits, semanticHits, countTitleHits)

	if limit == 0 {

		row := dbpool.QueryRow(
			context.Background(),
			`select
				   count(*) over () ::integer
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
			where blacklist.shadow is not true order by rank desc limit 1`,
			query.query)

		var count int
		err := row.Scan(&count)

		if err != nil {
			log.Println("error getting rows from database:", err)
			return nil, err
		}
		totalHits = count + countTitleHits
	}

	queryPerformance.EndTime = time.Now().UnixNano()
	queryPerformance.DeltaTime = queryPerformance.EndTime - queryPerformance.StartTime

	var tsquery string
	err = dbpool.QueryRow(context.Background(), "select websearch_to_tsquery('english', $1);", query.query).Scan(&tsquery)

	result.Performance = queryPerformance
	result.TsQuery = tsquery
	result.Query = query.query
	result.Page.Results = totalHits
	log.Println(totalHits, countExactTitleHits, countTitleHits)
	log.Println("searched:", query.query)
	log.Println("num found:", totalHits)

	return &result, nil

}
