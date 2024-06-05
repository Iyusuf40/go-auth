package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	postgresUrl := "postgres://yusuf:0@localhost:5432/test"
	conn, err := pgx.Connect(context.Background(), postgresUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	createTableStmt := `
		CREATE TABLE IF NOT EXISTS users (
			email 			varchar(256),
			first_name 		varchar(128),
			last_name 		varchar(128)
		)
	`
	cmdTag, err := conn.Exec(context.Background(), createTableStmt)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
		os.Exit(1)
	}

	insertStmt := `
		INSERT INTO users VALUES ($1, $2, $3);
	`

	cmdTag, err = conn.Exec(context.Background(), insertStmt, "mail", "f_name", "l_name")

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to insert: %v\n", err)
		os.Exit(1)
	}

	_, err = conn.Exec(context.Background(), "drop table users")

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to drop table: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(cmdTag)

	defer conn.Close(context.Background())
}
