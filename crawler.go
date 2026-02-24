package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgxpool"
)

func fetchBooks(dbpool *pgxpool.Pool) error {

	// we need:
	// [x] an html parser lib, batteries included baby!
	// [x] the date of the last time a book was _updated_
	// [x] the url of the first page of results sorted by date desc

	baseUrl := "https://archiveofourown.org/tags/Human%20Domestication%20Guide%20-%20GlitchyRobo/works?page=" + "1"

	var err error

	var lastPackageDate time.Time
	dbpool.QueryRow(context.Background(), "select packaged_at from books order by packaged_at desc").Scan(&lastPackageDate)

	log.Println("Starting book update...")
	log.Println("Last was update from ", lastPackageDate.Format(time.DateOnly))

	res, err := http.Get(baseUrl)

	if err != nil {
		log.Println("Error fetching books:", err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		log.Println("HTTP error:", res.StatusCode, res.Status)
		return errors.New(res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return err
	}

	doc.Find("ol > li > div > .datetime").Each(func(i int, s *goquery.Selection) {
		date := s.Text()
		log.Println(date)
	})

	return err
}

const test = `FFF: DEBUG: 2026-02-24 23:10:20,736: fetcher_requests.py(114): 
---------- REQ (GET) RequestsFetcher
https://archiveofourown.org/works/66691849/chapters/172065550?view_adult=true
FFF: DEBUG: 2026-02-24 23:10:21,027: fetcher_requests.py(127): response code:200
FFF: DEBUG: 2026-02-24 23:10:21,027: decorators.py(118): fromcache:False
FFF: DEBUG: 2026-02-24 23:10:21,027: decorators.py(129): random sleep(2.00-6.00):2.26
.FFF: DEBUG: 2026-02-24 23:10:23,289: requestable.py(55): Encoding:utf8
FFF: INFO: 2026-02-24 23:10:23,400: writer_epub.py(395): Saving EPUB Version 2.0
FFF: DEBUG: 2026-02-24 23:10:23,458: cli.py(63): Successfully wrote 'Ladies and Ladybugs-ao3_66691849.epub'`

func findPdfDest(commandOutput string) (string, error) {

	lines := strings.Split(commandOutput, "\n")

	for _, line := range lines {
		if strings.Contains(line, "Successfully wrote") {

			r, _ := regexp.Compile(`\'(?P<filename>.*)\'`)

			matches := r.FindStringSubmatch(line)
			log.Println(matches[1])

		}

	}

	return "", nil

}

func fetchBook(url string) (string, error) {

	fanficfare := "fanficfare" // TODO docker!!!!
	// TODO ALSO REMEMBER TO BACKDATE package date on Ladies and Ladybugs-ao3_66691849.epub !!
	archiveUrl := "https://archiveofourown.org"

	fullUrl := archiveUrl + url

	cmd := exec.Command(fanficfare, "-p", "-d", "--non-interactive", "--force", fullUrl)

	cmd.Dir = "../test books"
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
		return "", err
	}

	cmdOutput := string(out)

	fetchedPdfFilepath, err := findPdfDest(cmdOutput)

	if err != nil {
		return "", err
	}

	return fetchedPdfFilepath, nil
}
