package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/reward-rabieth/gym/api"
	"github.com/reward-rabieth/gym/storage"
	"github.com/reward-rabieth/gym/types"
	"github.com/reward-rabieth/gym/utils"
)

func SeedMember(storer storage.Storage, members types.GymParams, pwd string) *types.Gymmember {

	member, err := types.NewGymMember(members, pwd)

	if err != nil {
		log.Fatal(err)

	}

	if err := storer.CreateMember(member); err != nil {
		log.Fatal(err)

	}
	fmt.Println("new member=> ", member.Number)
	return member
}

func SeedMembers(s storage.Storage) {
	SeedMember(s, types.GymParams{
		Name:            "reward",
		Age:             26,
		Gender:          "male",
		Height:          5.2,
		Weight:          56,
		Membership:      "gold",
		PersonalTrainer: "goldie",
	}, "cypher99")

}
func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()
	utils.InitLogger()
	cfg := utils.Loadconfig()

	store, err := storage.NewPostgresStorage(cfg.Dbcfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	if *seed {
		fmt.Println("seeding the database")
		SeedMembers(store)
	}

	server := api.NewApiServer(cfg.Servercfg.Port, store)
	server.Run()

}
