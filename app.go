package main

import (
	"fmt"
	"github.com/onionltd/go-omg"
	"github.com/oniontree-org/go-oniontree"
	"github.com/urfave/cli/v2"
	pgperrors "golang.org/x/crypto/openpgp/errors"
	"net/http"
	"net/url"
	"strings"
)

const Version = "0.1"

type Application struct {
	ot  *oniontree.OnionTree
	app *cli.App
}

func (a *Application) handleOnionTreeOpen() cli.BeforeFunc {
	return func(c *cli.Context) error {
		ot, err := oniontree.Open(c.String("C"))
		if err != nil {
			return fmt.Errorf("failed to open OnionTree repository: %s", err)
		}
		a.ot = ot
		return nil
	}
}

func (a *Application) handleSyncCommand() cli.ActionFunc {
	normalizeURL := func(u string) (string, error) {
		nu, err := url.ParseRequestURI(u)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s://%s", nu.Scheme, nu.Host), nil
	}
	return func(c *cli.Context) error {
		id := c.Args().First()
		if id == "" {
			return fmt.Errorf("Missing a service ID")
		}

		client := goomg.NewClient(&http.Client{
			Timeout: c.Duration("timeout"),
		})

		service, err := a.ot.GetService(id)
		if err != nil {
			return fmt.Errorf("failed to get existing service data: %s", err)
		}

		if c.Bool("no-verify-signature") {
			fmt.Println("WARNING: signature verification is disabled!")
		}

		for _, u := range service.URLs {
			mirrors, err := client.GetMirrorsMessage(u)
			if err != nil {
				continue
			}

			if !c.Bool("no-verify-signature") {
				if _, err := mirrors.VerifySignature(service.PublicKeys); err != nil {
					if err == pgperrors.ErrUnknownIssuer {
						return fmt.Errorf("failed to verify signed message: unknown issuer")
					}
					continue
				}
			}

			urls, err := mirrors.List()
			if err != nil {
				continue
			}

			validUrls := make([]string, 0, len(urls))
			for i := range urls {
				u, err := normalizeURL(urls[i])
				if err != nil {
					return fmt.Errorf("failed to normalize URL: %s", err)
				}
				if !strings.HasSuffix(u, ".onion") {
					continue
				}
				validUrls = append(validUrls, u)
			}
			urls = validUrls

			if len(urls) == 0 {
				break
			}

			if c.Bool("replace") {
				service.SetURLs(urls)
			} else {
				service.AddURLs(urls)
			}

			if err := a.ot.UpdateService(service); err != nil {
				return fmt.Errorf("failed to update service: %s", err)
			}

			break
		}

		return nil
	}
}

func (a *Application) Run(args []string) error {
	return a.app.Run(args)
}
