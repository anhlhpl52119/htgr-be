package helpers

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintStruct(s any) {
	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
}
