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
	"github.com/jcbbb/gosar/user"
	_ "github.com/joho/godotenv/autoload"
)

var (
	POSTGRES_URI = os.Getenv("POSTGRES_URI")
	ADDR         = common.GetEnvStr("ADDR", ":3000")
)

func main() {
	var err error
	db.Pool, err = pgxpool.New(context.Background(), POSTGRES_URI)

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

	fmt.Printf("Server started at %v\n", ADDR)
	log.Fatal(http.ListenAndServe(ADDR, mux))

	defer db.Pool.Close()
}
