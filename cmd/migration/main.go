package main

import "github.com/M15t/ghoul/internal/functions/migration"

func main() {
	checkErr(migration.Run())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
