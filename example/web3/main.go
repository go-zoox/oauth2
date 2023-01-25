package main

import (
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-zoox/dotenv"
	"github.com/go-zoox/random"

	// "github.com/go-zoox/oauth2/web3"
	"github.com/go-zoox/crypto/hmac"
	"github.com/go-zoox/zoox"
	zd "github.com/go-zoox/zoox/default"
)

var secret = "secret"

type User struct {
	Nickname    string `json:"nickname"`
	Web3Address string `json:"address"`
	Nonce       string `json:"nonce"`
	// Nonce         uint64 `json:"nonce"`
}

func getUser(address string) *User {
	user := &User{
		Nickname:    "Zzz",
		Web3Address: address,
	}

	return user
}

func createUser(address string) *User {
	user := &User{
		Nickname:    "Zzz",
		Web3Address: address,
	}

	return user
}

func generateToken(address string) string {
	return hmac.Sha256(secret, address)
}

func main() {
	var cfg struct {
		ClientID     string `env:"CLIENT_ID"`
		ClientSecret string `env:"CLIENT_SECRET"`
	}
	if err := dotenv.Load(&cfg); err != nil {
		log.Fatal(err)
	}

	// var client, _ = web3.New(
	// 	cfg.ClientID,
	// 	cfg.ClientSecret,
	// 	"http://127.0.0.1:8080/login/web3/callback",
	// )

	r := zd.Default()

	r.Any("/", func(ctx *zoox.Context) {
		ctx.JSON(http.StatusOK, zoox.H{
			"message":   "Hello World",
			"user_id":   ctx.Session.Get("user_id"),
			"user_name": ctx.Session.Get("user_name"),
		})
	})

	login := r.Group("/login")

	login.Get("/web3", func(ctx *zoox.Context) {
		// client.Authorize("any", func(loginUrl string) {
		// 	ctx.Redirect(loginUrl)
		// })
		action := ctx.Query("action")
		switch action {
		case "authorize":
			address := ctx.Query("address")
			user := getUser(address)
			if user == nil {
				user = createUser(address)
			}

			// generate nonce
			user.Nonce = random.String(32)

			ctx.JSON(http.StatusOK, user)
			return
		case "token":
			address := ctx.Query("address")
			signature := ctx.Query("signature")

			user := getUser(address)

			signatureX := crypto.Keccak256(
				crypto.Keccak256([]byte("string challenge")),
				crypto.Keccak256([]byte(user.Nonce)),
			)
			pubkey, err := crypto.SigToPub(
				signatureX,
				[]byte(signature),
			)
			if err != nil {
				panic(err)
			}

			addressX := crypto.PubkeyToAddress(*pubkey)
			if string(addressX[:]) != address {
				ctx.JSON(http.StatusUnauthorized, zoox.H{
					"message": "Unauthorized",
				})
				return
			}

			token := generateToken(address)

			ctx.JSON(http.StatusOK, zoox.H{
				"message": "OK",
				"token":   token,
			})
			return
		// case "signup":
		// 	address := ctx.Bodies()["address"].(string)
		// 	user := createUser(address)
		// 	ctx.JSON(http.StatusOK, user)
		default:
			ctx.Render(200, "login.html", zoox.H{})
			return
		}
	})

	// login.Get("/web3/callback", func(ctx *zoox.Context) {
	// 	code := ctx.Query("code")
	// 	state := ctx.Query("state")

	// 	client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
	// 		if err != nil {
	// 			log.Println("[OAUTH2] Login Callback Error", err)
	// 			ctx.Redirect("/login/web3")
	// 			return
	// 		}

	// 		fmt.Println("user:", user)

	// 		// login success
	// 		ctx.Session.Set("user_id", user.ID)
	// 		ctx.Session.Set("user_name", user.Nickname)

	// 		ctx.Redirect("/")
	// 	})
	// })

	r.Run(":8080")
}
