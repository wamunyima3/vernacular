// Command admin runs one-off account tasks that need no web UI, such as
// making the first admin: admin -email you@example.org -role admin
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "admin:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	email := flag.String("email", "", "the user's email (they must have signed up)")
	role := flag.String("role", string(gen.RoleAdmin), "developer, contributor or admin")
	flag.Parse()
	if *email == "" || !gen.Role(*role).Valid() {
		flag.Usage()
		return errors.New("a valid -email and -role are required")
	}

	pool, err := database.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

	return database.InTx(ctx, pool, "cli:admin", func(tx pgx.Tx) error {
		n, err := dbgen.New(tx).SetRole(ctx, dbgen.SetRoleParams{Email: *email, Role: *role})
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("no user with email %s; sign up first", *email)
		}
		fmt.Printf("%s is now %s\n", *email, *role)
		return nil
	})
}
