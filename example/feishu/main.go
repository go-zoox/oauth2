package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-zoox/dotenv"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/oauth2/feishu"
	"github.com/go-zoox/zoox"
	zd "github.com/go-zoox/zoox/default"
)

func main() {
	var cfg struct {
		ClientID     string `env:"CLIENT_ID"`
		ClientSecret string `env:"CLIENT_SECRET"`
	}
	if err := dotenv.Load(&cfg); err != nil {
		log.Fatal(err)
	}

	var client, _ = feishu.New(
		cfg.ClientID,
		cfg.ClientSecret,
		"http://127.0.0.1:8080/login/feishu/callback",
	)

	r := zd.Default()

	r.Get("/", func(ctx *zoox.Context) {
		ctx.JSON(http.StatusOK, zoox.H{
			"message":   "Hello World",
			"user_id":   ctx.Session.Get("user_id"),
			"user_name": ctx.Session.Get("user_name"),
		})
	})

	login := r.Group("/login")

	login.Get("/feishu", func(ctx *zoox.Context) {
		client.Authorize("any", func(loginUrl string) {
			ctx.Redirect(loginUrl)
		})
	})

	login.Get("/feishu/callback", func(ctx *zoox.Context) {
		code := ctx.Query("code")
		state := ctx.Query("state")

		client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
			if err != nil {
				log.Println("[OAUTH2] Login Callback Error", err)
				ctx.Redirect("/login/feishu")
				return
			}

			fmt.Println("user:", user)

			// login success
			ctx.Session.Set("user_id", user.ID)
			ctx.Session.Set("user_name", user.Nickname)

			ctx.Redirect("/")
		})
	})

	r.Run(":8080")
}
