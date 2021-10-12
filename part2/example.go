package main

import (
	"context"
	"fmt"
	"log"

	"entimport-example/ent"
	"entimport-example/ent/user"

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

	// Create Car
	vw := client.Car.
		Create().
		SetModel("volkswagen").
		SetColor("blue").
		SetEngineSize(1400).
		SaveX(ctx)
	fmt.Println("First car created:", vw)

	// Update the user - add the car relation.
	client.User.Update().Where(user.ID(zeev.ID)).AddCars(vw).SaveX(ctx)

	// Query all cars that belong to user.
	cars := zeev.QueryCars().AllX(ctx)
	fmt.Println("User cars:", cars)

	// Create a second Car.
	delorean := client.Car.
		Create().
		SetModel("delorean").
		SetColor("silver").
		SetEngineSize(9999).
		SaveX(ctx)
	fmt.Println("Second car created:", delorean)

	// Update the user - add another car relation.
	client.User.Update().Where(user.ID(zeev.ID)).AddCars(delorean).SaveX(ctx)

	// Traverse the sub-graph.
	cars = delorean.
		QueryUser().
		QueryCars().
		AllX(ctx)
	fmt.Println("User cars:", cars)
}
