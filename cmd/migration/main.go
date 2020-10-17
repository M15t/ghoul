package main

import (
	"fmt"

	"github.com/M15t/ghoul/internal/functions/migration"
)

func main() {
	if err := migration.Run(); err != nil {
		fmt.Printf("ERROR: %+v\n", err)
	}
}
