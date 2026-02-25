package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookMeta struct {
	url        string
	lastUpdate time.Time
}

func findEpubDest(commandOutput string) (string, error) {
	lines := strings.Split(commandOutput, "\n")

	var filepath string

	for _, line := range lines {
		log.Println(line)
		if strings.Contains(line, "Successfully wrote") {

			r, _ := regexp.Compile(`\'(?P<filename>.*)\'`)

			matches := r.FindStringSubmatch(line)
			filepath = matches[1]
			log.Println(filepath)
		}
	}

	if filepath == "" {
		return "", errors.New("crawl failed!")
	}

	return filepath, nil
}

func fetchBook(url string) (string, error) {
	fanficfare := "fanficfare" // TODO docker!!!!
	// TODO ALSO REMEMBER TO BACKDATE package date on Ladies and Ladybugs-ao3_66691849.epub !!
	archiveUrl := "https://archiveofourown.org"

	fullUrl := archiveUrl + url

	log.Println("Updating book at", url)

	cmd := exec.Command(fanficfare, "-p", "-d", "--non-interactive", "--force", fullUrl)

	cmd.Dir = "../test books"
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
		return "", err
	}

	cmdOutput := string(out)

	fetchedPdfFilepath, err := findEpubDest(cmdOutput)

	if err != nil {
		return "", err
	}

	if fetchedPdfFilepath == "" {
		return "", errors.New("no file path emitted :<")
	}

	return cmd.Dir + "/" + fetchedPdfFilepath, nil
}

// could be inlined?
func updateBook(url string, dbpool *pgxpool.Pool) error {
	epubFilePath, err := fetchBook(url)
	log.Println("update book>", epubFilePath)
	if err != nil || epubFilePath == "" {
		saveCrawlFailed(url, dbpool)
		return err

	}
	return insertOrUpdateBook(epubFilePath, dbpool)

}

func parsePage(reader io.ReadCloser) ([]BookMeta, error) {

	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return nil, err
	}

	var links []BookMeta

	log.Println("Books on page:")

	doc.Find("ol > li > div > .datetime").Each(func(i int, s *goquery.Selection) {
		ao3DateLayout := "02 Jan 2006"
		date, err := time.Parse(ao3DateLayout, s.Text())

		if err != nil {
			log.Println(err)
			return
		}

		link, _ := s.Parent().Find("h4").Find("a").Attr("href")
		formattedDate := date.Format(time.DateOnly)
		log.Println("	", link, "last updated at", formattedDate)
		links = append(links, BookMeta{link, date})
	})

	return links, err

}

func fetchPage(pageNr int) (io.ReadCloser, error) {

	baseUrl := "https://archiveofourown.org/tags/Human%20Domestication%20Guide%20-%20GlitchyRobo/works?page="

	pageNrStr := strconv.Itoa(pageNr)

	pageUrl := baseUrl + pageNrStr

	res, err := http.Get(pageUrl)

	if err != nil {
		log.Println("Error fetching books:", err)
		return nil, err
	}

	if res.StatusCode >= 400 {
		log.Println("HTTP error:", res.StatusCode, res.Status)
		return nil, errors.New(res.Status)
	}

	return res.Body, nil

}

func fetchBooks(dbpool *pgxpool.Pool) error {

	// we need:
	// [x] an html parser lib, batteries included baby!
	// [x] the date of the last time a book was _updated_
	// [x] the url of the first page of results sorted by date desc

	var lastPackageDate time.Time
	dbpool.QueryRow(context.Background(), "select packaged_at from books order by packaged_at desc").Scan(&lastPackageDate)

	log.Println("Starting book update...")
	log.Println("Last was update from ", lastPackageDate.Format(time.DateOnly))

	needNextPage := true
	pageNr := 1
scrape:
	for needNextPage || pageNr > 200 { // wow i hope this is enough pages, change this if we double again to not fail full library rebuilds

		log.Println("Fetching page ", pageNr)

		body, err := fetchPage(pageNr)
		links, err := parsePage(body)

		if err != nil {
			log.Println("Error parsinng page: ", err)
			return err
		}

		for _, link := range links {
			if link.lastUpdate.After(lastPackageDate) {
				updateBook(link.url, dbpool)
			} else {
				needNextPage = false
				break scrape
			}
		}

		body.Close()

		pageNr = pageNr + 1
	}

	log.Println("Update successful!")
	cleanFailedCrawls(dbpool)

	return nil
}

type FailedCrawl struct {
	Id        int
	WorkUrl   string    `db:"work_url"`
	CreatedAt time.Time `db:"created_at"`
}

func cleanFailedCrawls(dbpool *pgxpool.Pool) {

	// get failed crawls
	rows, err := dbpool.Query(
		context.Background(),
		"select * from crawl_error",
	)

	if err != nil {
		log.Println(err)
		return
	}

	fails, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[FailedCrawl])

	if err != nil {
		log.Println(err)
		return
	}

	// abort if there's nothing to do
	if len(fails) == 0 {
		return

	}

	log.Println("Cleaning failed crawls!")

	for _, fail := range fails {
		url := fail.WorkUrl

		err := updateBook(url, dbpool)

		if err != nil {
			log.Println(err)
			continue
		}

		dbpool.Exec(context.Background(),
			"delete from crawl_error where id = $1",
			fail.Id,
		)
	}

}
