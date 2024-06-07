package main

import (
	"math/rand"
	"net/http"

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
		},
	}

	task2 := &boomer.Task{
		Weight: 1,
		Fn: func() {
			responseTime := rand.Int63n(1000)
			_ = rand.Int63n(1024 * 1024)
			boomer.RecordFailure(http.MethodGet, "bar", responseTime, errStatusCode)
		},
	}

	boomer.Run(task1, task2)
}
