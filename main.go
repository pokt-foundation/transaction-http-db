package main

import (
	"github.com/pokt-foundation/utils-go/environment"
	"github.com/sirupsen/logrus"
)

const (
	// TODO - declare needed environment variable names here
	todo1 = "TO_DO_1"
	todo2 = "TO_DO_2"

	todoDefault = 12345
)

type options struct {
	// TODO - declare needed environment variable struct fields here
	todo1 string
	todo2 int64
}

func gatherOptions() options {
	return options{
		// TODO - access environnment variables using methods here

		// MustGet<TYPE> funcs will panic if env var not found
		todo1: environment.MustGetString(todo1),

		// Get<TYPE> funcs will default to the second argument if not found
		todo2: environment.GetInt64(todo2, todoDefault),
	}
}

func main() {
	log := logrus.New()
	// log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	// TODO - use environment variables from gatherOptions
	_ = gatherOptions()
}
