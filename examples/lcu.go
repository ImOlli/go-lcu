package main

import (
	"github.com/ImOlli/go-lcu/lcu"
	"log"
)

func main() {
	info, err := lcu.FindLCUConnectInfo()

	if err != nil {

		if lcu.IsProcessNotFoundError(err) {
			log.Println("No LeagueClient process found!")
			return
		}

		panic(err)
	}

	log.Printf("LeagueClient is running on port %s and you can authenticate with following token: %s", info.Port, info.AuthToken)
}
