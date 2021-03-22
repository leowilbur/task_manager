package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/urfave/cli"

	"github.com/leowilbur/task_manager/pkg/models"
	"github.com/leowilbur/task_manager/pkg/rest"
)

func main() {
	app := cli.NewApp()
	app.Name = "examination"
	app.Usage = "examination service for Autonomous"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "port",
			Usage:  "Port the server listens to",
			EnvVar: "PORT",
			Value:  "3000",
		},
	}
	app.Action = startAPIRestService()

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("server failed with error %v", err)
	}

}

func startAPIRestService() cli.ActionFunc {
	return func(c *cli.Context) error {
		port := c.String("port")

		workerContext := &models.Context{
			Cron:    cron.New(),
			JobInfo: make(map[cron.EntryID]models.Task),
		}
		go func() {
			workerContext.Cron.Run()
		}()

		// New rest api
		router, err := rest.New(workerContext)
		if err != nil {
			return err
		}

		_, ctxCancel := context.WithCancel(context.Background())
		go func() {
			osSignal := make(chan os.Signal, 1)
			signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
			<-osSignal
			ctxCancel()
			// Stop cron and graceful shutdown
			go func() {
				router.Context.Cron.Stop()
				timer := time.NewTimer(3 * time.Second)
				<-timer.C

				log.Fatal(".........Graceful Shutdown.........")
			}()
		}()

		if err := http.ListenAndServe(":"+port, router); err != nil {
			return err
		}

		return nil
	}
}
