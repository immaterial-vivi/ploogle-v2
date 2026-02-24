package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

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
	ebookData.Author = author
	ebookData.Title = title
	ebookData.Url = bookUrl
	ebookData.Summary = summary

	ebookData.Chapters = book.Chapters

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

	// parse commandline arguments
	// -----
	shouldMigrateDbPtr := flag.Bool("m", false, "migrate the db schema")
	shouldLoadDbPtr := flag.Bool("l", false, "load books from bookdir into the db")

	queryPtr := flag.String("q", "", "search for this")
	pagePtr := flag.Int("p", 0, "page")
	pageSize, err := strconv.Atoi(os.Getenv("PLOOGLE_PAGE_SIZE"))

	fmt.Println(pageSize, err)
	flag.Parse()
	// -----

	// db setup
	// ------
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

	// prepare prepared statements
	// clean out any leftovers, probably not needed
	//dbpool.Exec(context.Background(), "deallocate all;")

	// load search related statments
	//initSearchStatements(dbpool)
	// -----

	// apply database schema
	if *shouldMigrateDbPtr {
		migrate(dbpool)
	}

	// reload books from disk
	if *shouldLoadDbPtr {
		fillDB(dbpool)
	}

	var query Query

	query.query = *queryPtr
	query.limit = pageSize
	query.offset = *pagePtr * pageSize

	res, err := search(query, dbpool)

	if err != nil {
		log.Fatalln(err)
	}
	resJson, err := json.Marshal(res)

	fmt.Println(string(resJson), err)

}
