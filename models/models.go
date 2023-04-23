package models

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	log.Printf("Get DB env %s\n", env)
	_ = pop.LoadConfigFile()
	log.Printf("Connections %s\n", pop.Connections)
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"
}
