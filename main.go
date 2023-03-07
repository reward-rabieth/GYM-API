package main

import (
	"log"

	"github.com/reward-rabieth/gym/api"
	"github.com/reward-rabieth/gym/storage"
	"github.com/reward-rabieth/gym/utils"
)

func main() {
	utils.InitLogger()
	cfg := utils.Loadconfig()

	store, err := storage.NewPostgresStorage(cfg.Dbcfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewApiServer(cfg.Servercfg.Port, store)
	server.Run()

}
