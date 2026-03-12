package main

import (
	"fmt"
	"permgraph/internal/db"
)

func main() {

	conn := "postgres://postgres:password@localhost:5432/permgraph"

	_, err := db.NewPool(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println("DB ready")
}