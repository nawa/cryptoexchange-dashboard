package main

import (
	"github.com/nawa/cryptoexchange-wallet-info/cmd"
)

func main() {
	cmd.Execute()

	//TODO remove me
	// server := http.NewServer(8080)
	// err := server.Start()
	// if err != nil {
	// 	panic(err)
	// }
}
