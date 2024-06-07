package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/myzhan/boomer"
)

var errStatusCode = "unexpected response status code"

func main() {
	task1 := &boomer.Task{
		Weight: 1,
		Fn: func() {
			responseTime := rand.Int63n(1000)
			contentLength := rand.Int63n(1024 * 1024)
			boomer.RecordSuccess(http.MethodGet, "foo", responseTime, contentLength)
			time.Sleep(100 * time.Millisecond)
		},
	}

	task2 := &boomer.Task{
		Weight: 1,
		Fn: func() {
			responseTime := rand.Int63n(1000)
			_ = rand.Int63n(1024 * 1024)
			boomer.RecordFailure(http.MethodGet, "bar", responseTime, errStatusCode)
			time.Sleep(100 * time.Millisecond)
		},
	}

	boomer.Run(task1, task2)
}
