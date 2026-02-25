package main

import (
	"log"
	"net/http"

	"github.com/go-co-op/gocron-ui/server"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func QueueJobs(dbpool *pgxpool.Pool) (http.Handler, error) {
	log.Println("cron starting")

	scheduler, err := gocron.NewScheduler()

	if err != nil {
		log.Println("Failed to start cron scheduler!", err)
		return nil, err
	}

	_, err = scheduler.NewJob(
		gocron.CronJob("32 22 * * *", false),
		gocron.NewTask(
			func(dbpool *pgxpool.Pool) error {
				return fetchBooks(dbpool)
			},
			dbpool,
		),
		gocron.WithName("book updates"),
	)

	if err != nil {
		log.Println("Failed to start cron job for bookUpdateTask", err)
		return nil, err
	}

	scheduler.Start()

	cronUiServer := server.NewServer(scheduler, 0)
	// log.Println("GoCron UI running on :9006")
	// log.Fatal(http.ListenAndServe(":9006", cronUiServer.Router))
	return cronUiServer.Router, nil
}
