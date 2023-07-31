package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jcbbb/gosar/auth"
	"github.com/jcbbb/gosar/common"
	"github.com/jcbbb/gosar/db"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello world")

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	postgresUri := os.Getenv("POSTGRES_URI")

	db.Pool, err = pgxpool.New(context.Background(), postgresUri)

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/user/register", common.MakeHandlerFuncMap(map[string]common.ApiFunc{
		"POST": auth.HandleSignup,
	}))

	log.Fatal(http.ListenAndServe(":3000", mux))
	defer db.Pool.Close()
}
