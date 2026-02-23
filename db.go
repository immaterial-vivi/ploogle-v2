package main

import (
	"context"
	"log"
	"time"

	"github.com/ArcadiaLin/go-epub"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EBookData struct {
	url      string
	title    string
	author   string
	summary  string
	chapters []epub.Chapter
}

type Book struct {
	id        int       `db:"id"`
	url       string    `db:"url"`
	title     string    `db:"title"`
	author    string    `db:"author"`
	summary   string    `db:"summary"`
	createdAt time.Time `db:"created_at"`

	// author's notes are stored as chapter 0
}

type Chapter struct {
	id            int    `db:"id"`
	chapterNumber int    `db:"chapter"`
	title         string `db:"chapter_title"`
	chapterUrl    string `db:"url"`
	bookId        int    `db:"book_id"`
	// text         string    `db:"chapter_text"`
	createdAt time.Time `db:"created_at"`
	//the chapter notes are stored in the text directly
}

type BlacklistEntry struct {
	id        int       `db:"id"`
	bookId    int       `db:"book_id"`
	reason    string    `db:"reason"`
	createdAt time.Time `db:"created_at"`
}

// we only ever update a whole book at once because we fetch entire ebooks through fanficfare
// this is overfetching, but ao3 doesn't have an api and the alternative is web scraping. This feels more stable
func upsertBook(book EBookData, dbpool *pgxpool.Pool) error {

	_, err := dbpool.Exec(context.Background(),
		"insert into books(title, author, url, summary) values ($1, $2, $3, $4) on conflict (url) do update set title=$1, author=$2, summary=$4",
		book.title, book.author, book.url, book.summary)

	var bookId int
	err = dbpool.QueryRow(context.Background(), "select id from books where url like $1", book.url).Scan(&bookId)

	for index, chapter := range book.chapters {
		text := chapter.Text()
		title := chapter.Title
		chapter_url := chapter.Url
		if chapter_url == "" {
			chapter_url = book.url
		}

		_, err = dbpool.Exec(context.Background(),
			"insert into chapters(chapter, chapter_title, url, book_id, chapter_text) values ($1, $2, $3, $4, $5) on conflict (url) do update set chapter=$1, chapter_title=$2, chapter_text=$5",
			index, title, chapter_url, bookId, text)

		if err != nil {
			log.Println(err)
		}
	}

	return err
}

func blacklistBook(bookId int, reason string, dbpool *pgxpool.Pool) error {
	_, err := dbpool.Exec(context.Background(),
		"insert into blacklist(book_id, reason) values($1, $2)",
		bookId, reason)

	return err
}

func deleteBlacklistEntry(id int, dbpool *pgxpool.Pool) error {

	_, err := dbpool.Exec(context.Background(), "delete from blacklist where id=$1", id)
	return err
}

func listBlacklist(dbpool *pgxpool.Pool) ([]BlacklistEntry, error) {

	rows, _ := dbpool.Query(context.Background(), "select * from blacklist")

	var blacklist []BlacklistEntry

	for rows.Next() {
		var entry BlacklistEntry
		err := rows.Scan(&entry)

		if err != nil {
			return nil, err
		}

		blacklist = append(blacklist, entry)
	}

	return nil, rows.Err()
}
