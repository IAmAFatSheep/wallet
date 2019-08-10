package main

import (
	"fmt"

	"github.com/blockcypher/gobcy"
)

func main() {
	btc := gobcy.API{"a32f202761be4948affc85fd1f0a7f93", "btc", "main"}
	hook, err := btc.CreateHook(gobcy.Hook{Event: "tx-confirmation", Address: "15qx9ug952GWGTNn7Uiv6vode4RcGrRemh", URL: "http://localhost:9000/callbacks/forwards"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", hook)
}