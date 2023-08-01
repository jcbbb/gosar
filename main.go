package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jcbbb/gosar/auth"
	"github.com/jcbbb/gosar/common"
	"github.com/jcbbb/gosar/db"
	"github.com/jcbbb/gosar/user"
	"github.com/joho/godotenv"
)

func main() {
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

	mux.HandleFunc("/user/auth", common.MakeHandlerFuncMap(map[string]common.ApiFunc{
		"POST": auth.HandleLogin,
	}))

	mux.HandleFunc("/user/", common.MakeHandlerFuncMap(map[string]common.ApiFunc{
		"GET": auth.EnsureAuth(user.HandleGetUser),
	}))

	mux.HandleFunc("/user/phone/", common.MakeHandlerFuncMap(map[string]common.ApiFunc{
		"POST":   auth.EnsureAuth(user.HandleAddPhone),
		"GET":    auth.EnsureAuth(user.HandleGetPhones),
		"PUT":    auth.EnsureAuth(user.HandleUpdatePhone),
		"DELETE": auth.EnsureAuth(user.HandleDeletePhone),
	}))

	log.Fatal(http.ListenAndServe(":3000", mux))
	defer db.Pool.Close()
}
