package main

import (
	"context"
	"fmt"
	"permgraph/internal/db"
	"permgraph/internal/tenant"
)

func main() {

	conn := "postgres://postgres:pass@localhost:5432/permgraph"

	pool, err := db.NewPool(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println("DB ready")

	tenantService := tenant.NewService(pool)

	id, err := tenantService.CreateTenant(context.Background(), "acme")
	if err != nil {
		panic(err)
	}

	fmt.Println("Created tenant:", id)

}