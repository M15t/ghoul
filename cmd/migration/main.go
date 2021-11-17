package main

import (
	"fmt"

	"ghoul/internal/functions/migration"
)

func main() {
	if err := migration.Run(); err != nil {
		fmt.Printf("ERROR: %+v\n", err)
	}
}
