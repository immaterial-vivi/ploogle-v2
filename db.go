package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ArcadiaLin/go-epub"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EBookData struct {
	Url         string
	Title       string
	Author      string
	Summary     string
	Chapters    []epub.Chapter
	PublishedAt time.Time
	UpdatedAt   time.Time
	PackagedAt  time.Time
}

type Book struct {
	Id        int       `db:"id"`
	Url       string    `db:"url"`
	Title     string    `db:"title"`
	Author    string    `db:"author"`
	Summary   string    `db:"summary"`
	CreatedAt time.Time `db:"created_at"`

	// author's notes are stored as chapter 0
}

type Chapter struct {
	Id            int    `db:"id"`
	ChapterNumber int    `db:"chapter"`
	Title         string `db:"chapter_title"`
	ChapterUrl    string `db:"url"`
	BookId        int    `db:"book_id"`
	// text         string    `db:"chapter_text"`
	CreatedAt time.Time `db:"created_at"`
	//the chapter notes are stored in the text directly
}

type BlacklistEntry struct {
	Id        int       `db:"id"`
	BookId    int       `db:"book_id"`
	Reason    string    `db:"reason"`
	CreatedAt time.Time `db:"created_at"`
}

func migrate(dbpool *pgxpool.Pool) {

	log.Print("Running db migration...")

	// get the schema defs and init the db
	b, err := os.ReadFile("schema.sql") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	schema := string(b)

	_, err = dbpool.Exec(context.Background(), schema)

	if err != nil {
		log.Fatalln("Unable to load schema:", err)
	}

	log.Print("Done! \n")
}

func saveCrawlFailed(url string, dbpool *pgxpool.Pool) error {
	_, err := dbpool.Exec(context.Background(), "insert into crawl_error (url) values ($1)", url)
	return err
}

func insertOrUpdateBook(path string, dbpool *pgxpool.Pool) error {
	book, err := epub.ReadBook(path)
	if err != nil {
		log.Println("could not read book at", path, err)
		return err
	}

	title, err := book.Title()
	author, err := book.Creator()
	url, err := book.MetadataByKey("source")
	bookUrl := url[0]

	metadata := book.AllMetadata()
	// fmt.Println("metadata:", metadata)\
	description, ok := metadata["description"]
	summary := ""
	if ok {
		summary = description[0]
	}

	if err != nil {
		return err
	}

	var ebookData EBookData
	ebookData.Author = author
	ebookData.Title = title
	ebookData.Url = bookUrl
	ebookData.Summary = summary
	ebookData.Chapters = book.Chapters

	// get the dates
	publishedAtDate, err := time.Parse(time.DateOnly, metadata["date"][0])
	ebookData.PublishedAt = publishedAtDate
	packagedAtDate, err := time.Parse(time.DateOnly, metadata["date"][1])
	ebookData.PackagedAt = packagedAtDate
	updatedAtDate, err := time.Parse(time.DateOnly, metadata["date"][2])
	ebookData.UpdatedAt = updatedAtDate
	if err != nil {
		log.Println(err)
	}
	// log.Println(ebookData)
	upsertBook(ebookData, dbpool)
	return err

}

func fillDB(dbpool *pgxpool.Pool) {
	booksDir := os.Getenv("BOOKS_DIR")

	entries, err := os.ReadDir(booksDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		log.Println("adding", e.Name())
		bookPath := fmt.Sprintf("%s/%s", booksDir, e.Name())
		insertOrUpdateBook(bookPath, dbpool)
	}

}

// we only ever update a whole book at once because we fetch entire ebooks through fanficfare
// this is overfetching, but ao3 doesn't have an api and the alternative is web scraping. This feels more stable
func upsertBook(book EBookData, dbpool *pgxpool.Pool) error {

	_, err := dbpool.Exec(context.Background(),
		"insert into books(title, author, url, summary, published_at, updated_at, packaged_at) values ($1, $2, $3, $4, $5, $6, $7) on conflict (url) do update set title=$1, author=$2, summary=$4, updated_at=$6, packaged_at=$7",
		book.Title, book.Author, book.Url, book.Summary, book.PublishedAt, book.UpdatedAt, book.PackagedAt)

	if err != nil {
		log.Println(err)
	}

	var bookId int
	err = dbpool.QueryRow(context.Background(), "select id from books where url like $1", book.Url).Scan(&bookId)

	if err != nil {
		log.Println(err)
	}

	for index, chapter := range book.Chapters {
		text := chapter.Text()
		title := chapter.Title
		chapter_url := chapter.Url
		if chapter_url == "" {
			chapter_url = book.Url
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
