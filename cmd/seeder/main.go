package main

import (
	"bitbucket.org/microret/oxus/internal/config"
	"bitbucket.org/microret/oxus/internal/db/connect"
	postgres "bitbucket.org/microret/oxus/internal/db/init"
	"bitbucket.org/microret/oxus/seeds"
	"flag"
	"fmt"
	"strings"
)

type flags []string

func (f *flags) String() string {
	return strings.Join(*f, ", ")
}

func (f *flags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {
	var sources flags

	flag.Var(&sources, "s", "Specify the source for seeding")
	flag.Parse()

	conf, err := config.Load()
	if err != nil {
		fmt.Printf("err config.Load() %s\n", err)
		return
	}

	dbase, err := postgres.Connect(conf.Db)
	if err != nil {
		fmt.Printf("err db.Connect() %s\n", err)
		return
	}
	connect.DB = dbase

	seeds.RunSeeds(sources)
}
