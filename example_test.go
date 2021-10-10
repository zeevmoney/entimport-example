package main

import (
	"context"
	"fmt"
	"log"

	"entimport-tutorial/ent"
	"entimport-tutorial/ent/user"

	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
)

func Example_EntImport() {
	client, err := ent.Open(dialect.MySQL, "root:pass@tcp(localhost:3306)/entimport?parseTime=True")

	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	if err := Do(ctx, client); err != nil {
		log.Fatal(err)
	}
	// Output:
	// User created: User(id=1, age=33, name=Zeev, last_name=Manilovich)
	// First car created: Car(id=1, model=volkswagen, color=blue, engine_size=1400, user_id=0)
	// User cars: [Car(id=1, model=volkswagen, color=blue, engine_size=1400, user_id=1)]
	// Second car created: Car(id=2, model=delorean, color=silver, engine_size=9999, user_id=0)
	// User cars: [Car(id=1, model=volkswagen, color=blue, engine_size=1400, user_id=1) Car(id=2, model=delorean, color=silver, engine_size=9999, user_id=1)]

}

func Do(ctx context.Context, client *ent.Client) error {
	// Create User
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

	// Update the user - add the car relation
	client.User.Update().Where(user.ID(zeev.ID)).AddCars(vw).SaveX(ctx)

	// Query all cars that belong to user
	cars := zeev.QueryCars().AllX(ctx)
	fmt.Println("User cars:", cars)

	// Create a second Car
	delorean := client.Car.
		Create().
		SetModel("delorean").
		SetColor("silver").
		SetEngineSize(9999).
		SaveX(ctx)
	fmt.Println("Second car created:", delorean)

	// Update the user - add another the car relation
	client.User.Update().Where(user.ID(zeev.ID)).AddCars(delorean).SaveX(ctx)

	// Traverse the sub-graph.
	cars = delorean.
		QueryUser().
		QueryCars().
		AllX(ctx)
	fmt.Println("User cars:", cars)
	return nil
}
