package main

import (
	"context"
	"fmt"
	"log"

	"entimport-example/ent"

	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	client, err := ent.Open(dialect.MySQL, "root:pass@tcp(localhost:3306)/entimport?parseTime=True")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	example(ctx, client)
}

func example(ctx context.Context, client *ent.Client) {
	// Create a User.
	zeev := client.User.
		Create().
		SetAge(33).
		SetName("Zeev").
		SetLastName("Manilovich").
		SaveX(ctx)
	fmt.Println("User created:", zeev)
}
