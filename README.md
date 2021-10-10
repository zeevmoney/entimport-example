# Background:

A few months ago we announced the Schema Import Initiative, its goal is to to help support many use cases for generating
Ent schemas from external resources. Today, we are happy to announce the release of entimport - a command line tool
designed to create ent schemas from existing SQL databases. This is a feature the community has been asking for a long
time. It can help ent users or potential users to transition an existing setup in another language or ORM to ent. It can
also help with use cases where you would like to access the same data from different platforms (automatically sync
between them). The first version supports both `MySQL` and `PostgreSQL` databases.

# Getting Started:

We will do a quick example of end to end usage with a MySQL database, let’s get started.

## Create a Database

### Prerequisites:

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
    user_cars   bigint       null,
    constraint cars_users_cars
        foreign key (user_cars) references users (id)
            on delete set null
)
    collate = utf8mb4_bin;
```

## Initialize Ent Project

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

## Install entimport

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

## Running entimport

Now we are ready to import our MySQL schema to `ent`!
We will do it with the following command:

```shell
go run ariga.io/entimport/cmd/entimport
```



