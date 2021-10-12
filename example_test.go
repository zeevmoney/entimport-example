package main

import (
	"context"
	"log"

	"entimport-example/ent"

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
	example(ctx, client)
	// Output:
	// User created: User(id=1, age=33, name=Zeev, last_name=Manilovich, phone=)
	// First car created: Car(id=1, model=volkswagen, color=blue, engine_size=1400, user_id=0)
	// User cars: [Car(id=1, model=volkswagen, color=blue, engine_size=1400, user_id=1)]
	// Second car created: Car(id=2, model=delorean, color=silver, engine_size=9999, user_id=0)
	// User cars: [Car(id=1, model=volkswagen, color=blue, engine_size=1400, user_id=1) Car(id=2, model=delorean, color=silver, engine_size=9999, user_id=1)]
}
