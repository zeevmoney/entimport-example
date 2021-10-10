## Background:

A few months ago we announced the Schema Import Initiative, its goal is to to help support many use cases for generating
Ent schemas from external resources. Today, we are happy to announce the release of **entimport** - an importent (sorry)
command line tool designed to create ent schemas from existing SQL databases. This is a feature the community has been
asking for a long time. It can help ent users or potential users to transition an existing setup in another language or
ORM to ent. It can also help with use cases where you would like to access the same data from different platforms (
automatically sync between them). The first version supports both `MySQL` and `PostgreSQL` databases.

## Getting Started:

We will do a quick example of end to end usage with a MySQL database, let’s get started.

### Create a Database

Prerequisites:

- Docker - https://docs.docker.com/get-docker/
- Golang - https://golang.org/doc/install

Start the project in a new directory called `entimport-exmaple`. Create a file named `docker-compose.yaml` and paste the
following content inside:

```yaml
version: "3.7"

services:

  mysql8:
    platform: linux/amd64
    image: mysql
    environment:
      MYSQL_DATABASE: entimport
      MYSQL_ROOT_PASSWORD: pass
    healthcheck:
      test: mysqladmin ping -ppass
    ports:
      - "3308:3306"
```

This file contains service configuration for MySQL.  
Run it with the following command:

```shell
docker-compose up
# to run in the background add -d 
```

Next, we will create a simple schema. For this example we will use a relation between two entities:

- User
- Car

Connect to your database using MySQL console or any other DB tool and execute the following statements:

```mysql
# Users Table
create table users
(
    id        bigint auto_increment
        primary key,
    age       bigint       not null,
    name      varchar(255) not null,
    last_name varchar(255) null comment 'surname'
)
    collate = utf8mb4_bin;

# Cars Table

create table cars
(
    id          bigint auto_increment
        primary key,
    model       varchar(255) not null,
    color       varchar(255) not null,
    engine_size mediumint    not null,
    user_id     bigint       null,
    constraint cars_owners
        foreign key (user_id) references users (id)
            on delete set null
)
    collate = utf8mb4_bin;
```

We have created the tables mentioned above, with a One To Many relation:

- One user can have many cars
- Each car can have only one owner.

### Initialize Ent Project

In order to import our schema to `ent` we need to init an ent project first.

First, we need to init a go module inside the project root dir:

```shell
go mod init entimport-example
```

Run `ent` Init:

```shell
go install entgo.io/ent/cmd/ent
ent init
# alternatively you can do:
go get entgo.io/ent/cmd/ent
go run entgo.io/ent/cmd/ent init 
```

Your project should look like this:

```
├── docker-compose.yaml
├── ent
│   ├── generate.go
│   └── schema
└── go.mod
```

### Install entimport

Download `entimport`:

```shell
go get ariga.io/entimport
```

Check if it's working:

```shell
go run ariga.io/entimport/cmd/entimport -h
```

This command will print:

```
Usage of ntimport:
  -dialect string
        database dialect (default "mysql")
  -dsn string
        data source name (connection information)
  -schema-path string
        output path for ent schema (default "./ent/schema")
  -tables value
        comma-separated list of tables to inspect (all if empty)
```

### Run entimport

We are ready to import our MySQL schema to `ent`!

We will do it with the following command:
> This command will import all tables in our schema, you can also limit to specific tables using -tables flag.

```shell
go run ariga.io/entimport/cmd/entimport -dialect mysql -dsn "root:pass@tcp(localhost:3306)/entimport"
```

Now the magic happened and `enimport` wrote our schema to the `ent/schmea` directory:

``` {5-6}
├── docker-compose.yaml
├── ent
│   ├── generate.go
│   └── schema
│       ├── car.go
│       └── user.go
├── go.mod
└── go.sum
```

Let's check the Schemas:

- User Schema

```go title="<project>/ent/schema/user.go"
type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{field.Int("id"), field.Int("age"), field.String("name"), field.String("last_name").Optional().Comment("surname")}
}
func (User) Edges() []ent.Edge {
	return []ent.Edge{edge.To("cars", Car.Type)}
}
func (User) Annotations() []schema.Annotation {
	return nil
}
```

- Car Schema:

```go title="<project>/ent/schema/car.go"
type Car struct {
	ent.Schema
}

func (Car) Fields() []ent.Field {
	return []ent.Field{field.Int("id"), field.String("model"), field.String("color"), field.Int32("engine_size"), field.Int("user_id").Optional()}
}
func (Car) Edges() []ent.Edge {
	return []ent.Edge{edge.From("user", User.Type).Ref("cars").Unique().Field("user_id")}
}
func (Car) Annotations() []schema.Annotation {
	return nil
}
```

> **`entimport` successfully created entities and their relation!**

Let's try them out. First we must generate the ent schema. In order to [generate](https://entgo.io/docs/code-gen) `ent`
files from the produced schemas, run:

```shell
go run -mod=mod entgo.io/ent/cmd/ent generate ./schema

# OR:

go generate ./ent
```

Let's see our `ent` directory:

```
...
├── ent
│   ├── car
│   │   ├── car.go
│   │   └── where.go
│   ├── car.go
│   ├── car_create.go
│   ├── car_delete.go
│   ├── car_query.go
│   ├── car_update.go
│   ├── client.go
│   ├── config.go
│   ├── context.go
│   ├── ent.go
│   ├── enttest
│   │   └── enttest.go
│   ├── generate.go
│   ├── hook
│   │   └── hook.go
│   ├── migrate
│   │   ├── migrate.go
│   │   └── schema.go
│   ├── mutation.go
│   ├── predicate
│   │   └── predicate.go
│   ├── runtime
│   │   └── runtime.go
│   ├── runtime.go
│   ├── schema
│   │   ├── car.go
│   │   └── user.go
│   ├── tx.go
│   ├── user
│   │   ├── user.go
│   │   └── where.go
│   ├── user.go
│   ├── user_create.go
│   ├── user_delete.go
│   ├── user_query.go
│   └── user_update.go
...
```

### Ent Example

Let's run a quick example to prove that our schema works:

create a file named `example.go` in the root of the project, with the following content:

```go title="<project>/example.go"
import (
	"context"
	"log"

	"<project>/ent"

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
}
```

Let's try to add a user, write the following code at the end of the file:

```go title="<project>/example.go"
	// Create User
	zeev := client.User.
		Create().
		SetAge(33).
		SetName("Zeev").
		SetLastName("Manilovich").
		SaveX(ctx)
	fmt.Println("User created:", zeev)
```

Then run:

```shell
go run example.go
```

This will output: `# User created: User(id=2, age=33, name=Zeev, last_name=Manilovich)`

Let's check with the database if the user was really added:

```mysql
SELECT *
FROM users
WHERE name = 'Zeev';

+--+---+----+----------+
|id|age|name|last_name |
+--+---+----+----------+
|1 |33 |Zeev|Manilovich|
+--+---+----+----------+
```

Great! now let's play a little more with `ent` and add some relations:

```go title="<project>/example.go"
    import "entimport-tutorial/ent/user"

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
```

After Running the code above, your DB should hold a user with 2 cars in a O2M relation.

```mysql
+--+---+----+----------+
|id|age|name|last_name |
+--+---+----+----------+
|1 |33 |Zeev|Manilovich|
+--+---+----+----------+

+--+----------+------+-----------+-------+
|id|model     |color |engine_size|user_id|
+--+----------+------+-----------+-------+
|1 |volkswagen|blue  |1400       |1      |
|2 |delorean  |silver|9999       |1      |
+--+----------+------+-----------+-------+
```

### Syncing DB changes

Since we want to keep the DB in sync, we want `entimport` to be able to change the schema after the db was changed.
Let's see how it works. Run the following SQL code:

```mysql
alter table users
    add phone varchar(255) null;

create unique index users_phone_uindex
    on users (phone);
```

Now let's run `entimport` again:

```shell
go run ariga.io/entimport/cmd/entimport -dialect mysql -dsn "root:pass@tcp(localhost:3306)/entimport"
```

We can see that the `user.go` file was changed:

```go title="<project>/ent/schema/user.go"
func (User) Fields() []ent.Field {
	return []ent.Field{field.Int("id"), ..., field.String("phone").Optional().Unique()}
}
```

Now we can run `go generate ./ent` and use the new schema do a phone to the `User`.

## Future Plans and Caveats

As mentioned above this initial version supports MySQL and PostgreSQL databases.  
It also supports all types of SQL relations. We have plans to further upgrade the tool and add more features such as:

- Index support (currently Unique index is supported).
- Support for all data types (for example `uuid` in Postgres).
- Support for Default value in columns.
- Support for editing schema both manually and automatically (real upsert and not only overwrite)
- Postgres special types: postgres.NetworkType, postgres.BitType, *schema.SpatialType, postgres.CurrencyType,
  postgres.XMLType, postgres.ArrayType, postgres.UserDefinedType.

### Known Caveats:

- Schema files are overwritten by new calls to `entimport`.
- There is no difference in DB schema between `O2O Bidirectional` and `O2O Same Type` - both will result in the same
  `ent` schema.
- There is no difference in DB schema between `M2M Bidirectional` and `M2M Same Type` - both will result in the same
- `ent` schema.
- In recursive relations the `edge` names will be prefixed with `child_` & `parent_`.
- For example: `users` with M2M relation to itself will result in:

```go
func (User) Edges() []ent.Edge {
return []ent.Edge{edge.To("child_users", User.Type), edge.From("parent_users", User.Type)}
}
```

## Wrapping Up

In this post, we presented the Upsert API, a long-anticipated capability, that is available by feature-flag in Ent
v0.9.0. We discussed where upserts are commonly used in applications and the way they are implemented using common
relational databases. Finally, we showed a simple example of how to get started with the Upsert API using Ent.

Have questions? Need help with getting started? Feel free to [join our Slack channel](https://entgo.io/docs/slack/).


