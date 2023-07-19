package helloplugin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cs3org/reva"
	"github.com/cs3org/reva/pkg/appctx"
	"github.com/cs3org/reva/pkg/rhttp/global"
	"github.com/cs3org/reva/pkg/utils/cfg"
)

func init() {
	reva.RegisterPlugin(HelloWorld{})
}

type HelloWorld struct {
	c *config
}

type config struct {
	Name string `mapstructure:"name"`
}

func (c *config) ApplyDefaults() {
	if c.Name == "" {
		c.Name = "world"
	}
}

func (HelloWorld) RevaPlugin() reva.PluginInfo {
	return reva.PluginInfo{
		ID:  "http.services.helloplugin",
		New: New,
	}
}

func New(ctx context.Context, m map[string]any) (global.Service, error) {
	var c config
	if err := cfg.Decode(m, &c); err != nil {
		return nil, err
	}
	return &HelloWorld{c: &c}, nil
}

func (s *HelloWorld) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.GetLogger(r.Context())
		msg := fmt.Sprintf("Hello %s from a reva plugin!", s.c.Name)
		if _, err := w.Write([]byte(msg)); err != nil {
			log.Err(err).Msg("error writing response")
			return
		}
		log.Debug().Msgf("replied %s", msg)
	})
}

func (s *HelloWorld) Prefix() string {
	return "/helloplugin"
}

func (s *HelloWorld) Close() error { return nil }

func (s *HelloWorld) Unprotected() []string {
	return []string{"/"}
}

// guards to ensure HelloWorld implements a reva http service
var _ global.Service = (*HelloWorld)(nil)
var _ global.NewService = New
