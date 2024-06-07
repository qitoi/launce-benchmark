package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"math/rand"
	"net/http"

	"os/signal"
	"syscall"
	"time"

	"github.com/qitoi/launce"
	"github.com/qitoi/launce/taskset"
)

var (
	_ taskset.User = (*User)(nil)
)

type User struct {
	taskset.UserImpl
}

var errStatusCode = errors.New("unexpected response status code")

func (u *User) WaitTime() launce.WaitTimeFunc {
	return nil
}

func (u *User) TaskSet() taskset.TaskSet {
	return taskset.NewRandom(
		taskset.TaskFunc(func(ctx context.Context, u launce.User, s taskset.Scheduler) error {
			responseTime := time.Duration(rand.Int63n(int64(1000 * time.Millisecond)))
			contentLength := rand.Int63n(1024 * 1024)
			u.Report(http.MethodGet, "/foo", responseTime, contentLength, nil)
			return nil
		}),
		taskset.TaskFunc(func(ctx context.Context, u launce.User, s taskset.Scheduler) error {
			responseTime := time.Duration(rand.Int63n(int64(1000 * time.Millisecond)))
			contentLength := rand.Int63n(1024 * 1024)
			u.Report(http.MethodGet, "/bar", responseTime, contentLength, errStatusCode)
			return nil
		}),
	)
}

func main() {
	host := flag.String("master-host", "localhost", "master host")
	flag.Parse()

	transport := launce.NewZmqTransport(*host, 5557)
	worker, err := launce.NewWorker(transport)
	if err != nil {
		log.Fatal(err)
	}

	// 統計情報の集約プロセスを1ユーザーごとに割り当てる (初期値: 100)
	worker.StatsAggregationUsers = 1

	worker.RegisterUser("MyUser", func() launce.User {
		return &User{}
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		worker.Quit()
	}()

	if err := worker.Join(); err != nil {
		log.Fatal(err)
	}
}
