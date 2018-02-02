package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/abourget/eosapi"
)

func main() {
	api := eosapi.New("http://testnet1.eos.io")

	out, err := api.GetInfo()
	if err != nil {
		log.Fatalln("error:", err)
	}
	fmt.Println("GetInfo", out)

	cnt, err := json.MarshalIndent(out, "", "  ")
	fmt.Println("JSON:", string(cnt))

	// out, err := api.GetAccount("initn")
	// if err != nil {
	// 	log.Fatalln("error:", err)
	// }
	// fmt.Println("Account initm", out)

	// out2, err := api.GetCode("currency")
	// if err != nil {
	// 	log.Fatalln("error:", err)
	// }
	// fmt.Printf("Contract `currency: %+v\n", out2)

}
