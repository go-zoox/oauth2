package doreamon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-zoox/cookie"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/oauth2/doreamon"
)

type CreateOAuth2DoreamonHandlerConfig struct {
	ApplicationName string
	ClientID        string
	ClientSecret    string
	RedirectURI     string
}

type VerifyUserConfig struct {
	CookieKey string
	Token     *Token
}

type SaveUserConfig struct {
	CookieKey string
	Token     *Token
}

type Token struct {
	CookieKey string
	Cookie    func() cookie.Cookie
}

func (t *Token) Get() string {
	return t.Cookie().Get(t.CookieKey)
}

func (t *Token) Set(token string) error {
	t.Cookie().Set(t.CookieKey, token, &cookie.Config{
		MaxAge: 7 * 24 * time.Hour,
	})
	return nil
}

func CreateOAuth2DoreamonHandler(cfg *CreateOAuth2DoreamonHandlerConfig) func(
	w http.ResponseWriter,
	r *http.Request,
	VerifyUser func(cfg *VerifyUserConfig, token string, r *http.Request) error,
	SaveUser func(cfg *SaveUserConfig, user *oauth2.User, token *oauth2.Token) (tokenString string, err error),
	Next func() error,
) error {
	originPathCookieKey := "login_from"

	client, err := doreamon.New(&doreamon.DoreamonConfig{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURI:  cfg.RedirectURI,
		Scope:        "user,email",
		Version:      "2",
	})
	if err != nil {
		panic(err)
	}

	CookieKey := "go-zoox_oauth2_token"
	VerifyUserCfg := &VerifyUserConfig{
		CookieKey: CookieKey,
	}
	SaveUserCfg := &SaveUserConfig{
		CookieKey: CookieKey,
	}

	return func(
		w http.ResponseWriter,
		r *http.Request,
		VerifyUser func(cfg *VerifyUserConfig, token string, r *http.Request) error,
		SaveUser func(cfg *SaveUserConfig, user *oauth2.User, token *oauth2.Token) (tokenString string, err error),
		Next func() error,
	) error {
		if r.Method != "GET" {
			tokeString := VerifyUserCfg.Token.Get()
			if tokeString == "" {
				logger.Info("[oauth2] failed to verify user(1): %#v", fmt.Errorf("[oauth2][VerifyUser] failed to get cookie by key(%s), value: empty string", CookieKey))
				time.Sleep(1 * time.Second)

				// http.Redirect(w, r, "/login", http.StatusFound)

				w.WriteHeader(401)

				accept := r.Header.Get("Accept")
				acceptJSON := accept == "*/*" || strings.Contains(accept, "application/json")
				if acceptJSON {
					data, _ := json.Marshal(map[string]any{
						"code":    401000,
						"message": "Unauthorized",
					})
					w.Write(data)
					return nil
				}

				w.Write([]byte("Unauthorized"))
				return nil
			}

			return Next()
		}

		path := r.URL.Path
		if matched, _ := regexp.MatchString("\\.(js|css|json|txt|map|webmanifest|manifest|png|jpg|jpeg|webp|gif)$", path); err == nil && matched {
			// logger.Infof("[oauth2] ignore visit files \\.(js|css|json|txt|map|webmanifest|manifest)$ ...")
			return Next()
		}

		logger.Infof("[oauth2] comming (from: %s)...", r.URL.Path)

		if path == "/login" {
			logger.Infof("[oauth2] go login (from: %s %s)...", r.Method, r.URL.Path)
			client.Authorize(cfg.ApplicationName, func(loginUrl string) {
				http.Redirect(w, r, loginUrl, http.StatusFound)
			})
			return nil
		}

		if path == "/logout" {
			logger.Infof("[oauth2] go logout ...")

			client.Logout(cfg.ApplicationName, func(logoutUrl string) {
				http.Redirect(w, r, logoutUrl, http.StatusFound)
			})
			return nil
		}

		if path == "/login/doreamon/callback" {
			code := r.FormValue("code")
			state := r.FormValue("state")

			logger.Infof("[oauth2] login callback ...")
			client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
				if err != nil {
					log.Println("[OAUTH2] Login Callback Error", err)
					time.Sleep(3 * time.Second)
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}

				SaveUserCfg.Token = &Token{
					CookieKey: CookieKey,
					Cookie: func() cookie.Cookie {
						return cookie.New(w, r)
					},
				}
				tokenString, err := SaveUser(SaveUserCfg, user, token)
				if err != nil {
					logger.Info("failed to save user: %#v", err)
					time.Sleep(1 * time.Second)

					w.WriteHeader(500)
					w.Write([]byte("Failed to create user: " + user.Email))
					return
				}

				SaveUserCfg.Token.Set(tokenString)

				logger.Infof("[oauth2] login successed ...")
				http.Redirect(w, r, "/", http.StatusFound)
			})

			return nil
		}

		VerifyUserCfg.Token = &Token{
			CookieKey: CookieKey,
			Cookie: func() cookie.Cookie {
				return cookie.New(w, r)
			},
		}

		tokeString := VerifyUserCfg.Token.Get()
		if tokeString == "" {
			logger.Info("[oauth2] failed to verify user(1): %#v", fmt.Errorf("[oauth2][VerifyUser] failed to get cookie by key(%s), value: empty string", CookieKey))
			time.Sleep(1 * time.Second)
			http.SetCookie(w, &http.Cookie{
				Name:  "OriginPath",
				Value: path,
			})

			http.Redirect(w, r, "/login", http.StatusFound)
			return nil
		}

		logger.Infof("[oauth2] verify user ...")
		if err := VerifyUser(VerifyUserCfg, tokeString, r); err != nil {
			logger.Info("[oauth2] failed to verify user(2): %#v", err)
			time.Sleep(1 * time.Second)
			http.SetCookie(w, &http.Cookie{
				Name:  "OriginPath",
				Value: path,
			})

			http.Redirect(w, r, "/login", http.StatusFound)
			return nil
		}

		// success
		if OriginPath, err := r.Cookie(originPathCookieKey); err == nil && OriginPath.Value != "" {
			time.Sleep(1 * time.Second)

			http.SetCookie(w, &http.Cookie{
				Name:    originPathCookieKey,
				Value:   "",
				Expires: time.Unix(0, 0),
			})

			logger.Infof("[oauth2] save origin path for redirect back ...")
			http.Redirect(w, r, OriginPath.Value, http.StatusFound)
			return nil
		}

		return Next()
	}
}

func CreateHTTPHandler(
	ApplicationName string,
	VerifyUser func(cfg *VerifyUserConfig, token string, r *http.Request, w http.ResponseWriter) error,
	SaveUser func(cfg *SaveUserConfig, user *oauth2.User, token *oauth2.Token, r *http.Request, w http.ResponseWriter) (tokenString string, err error),
	Next func(w http.ResponseWriter, r *http.Request) error,
) http.Handler {
	if os.Getenv("DOREAMON_CLIENT_ID") == "" {
		panic("env DOREAMON_CLIENT_ID is required")
	}
	if os.Getenv("DOREAMON_CLIENT_SECRET") == "" {
		panic("env DOREAMON_CLIENT_SECRET is required")
	}
	if os.Getenv("DOREAMON_REDIRECT_URI") == "" {
		panic("env DOREAMON_REDIRECT_URI is required")
	}

	handler := CreateOAuth2DoreamonHandler(&CreateOAuth2DoreamonHandlerConfig{
		ApplicationName: ApplicationName,
		ClientID:        os.Getenv("DOREAMON_CLIENT_ID"),
		ClientSecret:    os.Getenv("DOREAMON_CLIENT_SECRET"),
		RedirectURI:     os.Getenv("DOREAMON_REDIRECT_URI"),
	})

	hfn := func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r, func(cfg *VerifyUserConfig, token string, r *http.Request) error {
			return VerifyUser(cfg, token, r, w)
		}, func(cfg *SaveUserConfig, user *oauth2.User, token *oauth2.Token) (tokenString string, err error) {
			return SaveUser(cfg, user, token, r, w)
		}, func() error {
			return Next(w, r)
		})

		if err != nil {
			fmt.Println("oauth 2error:", err)
			return
		}
	}

	return http.HandlerFunc(hfn)
}
