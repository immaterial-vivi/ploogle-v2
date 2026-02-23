package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"flag"

	"github.com/ArcadiaLin/go-epub"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

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

func readBooks() {
	book, err := epub.ReadBook("../books/Human Domestication Guide-ao3_45190954.epub")
	if err != nil {
		log.Println("could not read book:", err)
	}

	// Dublin Core metadata helpers
	if title, err := book.Title(); err == nil {
		fmt.Println("Title:", title)
	}

	// Generic metadata access
	values, err := book.MetadataByKey("language")

	fmt.Println("language:", values)

	// Iterate over the full metadata map (includes <meta> extensions)
	metadata := book.AllMetadata()
	fmt.Println("metadata:", metadata)

	// Work with chapters
	fmt.Println("Total chapters:", book.ChapterCount())
	//firstChapter, _ := book.ChapterByIndex(0)
	//fmt.Println(firstChapter.Text())

	// Concatenate the whole book into a single text blob
	//fmt.Println(book.AllChaptersText())

}

func insertOrUpdateBook(path string, dbpool *pgxpool.Pool) error {
	book, err := epub.ReadBook(path)
	if err != nil {
		log.Println("could not read book:", err)
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
	ebookData.author = author
	ebookData.title = title
	ebookData.url = bookUrl
	ebookData.summary = summary

	ebookData.chapters = book.Chapters

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

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, all config coming from system env")
	}

	databaseUser := os.Getenv("POSTGRES_USER")
	databasePassword := os.Getenv("POSTGRES_PASSWORD")
	databaseHost := os.Getenv("POSTGRES_HOST")
	databaseDBPath := os.Getenv("POSTGRES_DB")

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s", databaseUser, databasePassword, databaseHost, databaseDBPath)

	dbconfig, err := pgxpool.ParseConfig(databaseUrl)

	if err != nil {
		log.Println("dbconfig invalid")
	}

	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbconfig)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// var greeting string
	// err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	// if err != nil {
	// 	log.Fatalf("QueryRow failed: %v\n", err)
	// }

	// fmt.Println(greeting)

	shouldMigrateDbPtr := flag.Bool("m", false, "migrate the db schema")
	shouldLoadDbPtr := flag.Bool("l", false, "load books from bookdir into the db")

	queryPtr := flag.String("q", "", "search for this")
	flag.Parse()

	if *shouldMigrateDbPtr {
		migrate(dbpool)
	}
	// database setup completed!

	if *shouldLoadDbPtr {
		fillDB(dbpool)
	}
	//fillDB(dbpool)

	// err = insertOrUpdateBook("../books/Human Domestication Guide-ao3_45190954.epub", dbpool)
	// if err != nil {
	// 	log.Fatalf("insert failed: %v\n", err)
	// }

	var query Query

	query.query = *queryPtr

	res, err := search(query, dbpool)

	fmt.Println(res, err)

}

func _search(queryPtr *string, dbpool *pgxpool.Pool) error {

	log.Println("searching for query:", *queryPtr)
	rows, _ := dbpool.Query(context.Background(), "select title, author, url from books where title like $1", *queryPtr)

	for rows.Next() {
		var title string
		var author string
		var url string
		err := rows.Scan(&title, &author, &url)
		if err != nil {
			return err
		}
		result := fmt.Sprintf("%s by %s, at %s", title, author, url)
		log.Println(result)
	}

	return rows.Err()

	// panic("unimplemented")
}
