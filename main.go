package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/techschool/samplebank/api"
	db "github.com/techschool/samplebank/db/sqlc"
	"github.com/techschool/samplebank/util"
)

func main() {

	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Println("Error in Loadin config")
	}

	conn, err := sql.Open(conf.DbDriver, conf.DbSource)
	if err != nil {
		log.Println(err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(conf, store)
	if err != nil {
		log.Println("Error in creating server", err)
	}

	err = server.Start(conf.ServerAddress)
	if err != nil {
		log.Println(err)
	}

}
