package helpers

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintStruct(s any) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", data)
}
