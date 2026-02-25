package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"flag"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, all config coming from system env")
	}

	// parse commandline arguments
	// -----
	shouldMigrateDbPtr := flag.Bool("m", false, "migrate the db schema")
	shouldLoadDbPtr := flag.Bool("l", false, "load books from bookdir into the db")
	shouldRecrawlPtr := flag.Bool("c", false, "crawl ao3")

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
	// -----

	// apply database schema
	if *shouldMigrateDbPtr {
		migrate(dbpool)
	}

	// reload books from disk
	if *shouldLoadDbPtr {
		fillDB(dbpool)
	}

	// crawl web site
	if *shouldRecrawlPtr {
		fetchBooks(dbpool)
	}

	if *queryPtr != "" {
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

}
