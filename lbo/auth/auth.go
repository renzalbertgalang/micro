package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/micro/cli"
	"github.com/micro/go-micro/client"
	"github.com/micro/micro/plugin"
	proto "github.com/orangeapps/lbo-endgame/c/auth/api"
	"io/ioutil"
	"log"
	"net/http"
)

type lBOAuth struct {}

func NewLBOAuth() plugin.Plugin {
	return &lBOAuth{}
}

func (a *lBOAuth) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "lbo_auth",
			Usage: "Enable authentication for lbo-endgame platform.",
		},
	}
}

func (g *lBOAuth) Commands() []cli.Command {
	return nil
}

func (a *lBOAuth) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			c := ioutil.NopCloser(bytes.NewBuffer(b))
			decoder := json.NewDecoder(c)
			var t map[string]interface{}
			err = decoder.Decode(&t)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			method, ok := t["method"]
			if !ok {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if method != "AuthService.Authenticate" && method != "AuthService.Validate"  {
				token := r.Header.Get("Token")
				if token == "" {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}
				log.Print(token)
				protoCli := proto.NewAuthService("c.auth", client.NewClient())
				_, err := protoCli.Validate(context.Background(), &proto.ValidateReq{
					AccessToken: token,
				})
				if err != nil {
					log.Print(err)
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			h.ServeHTTP(w, r)
		})
	}
}

func (p *lBOAuth) Init(ctx *cli.Context) error {
	return nil
}

func (p *lBOAuth) String() string {
	return "lbo_auth"
}
