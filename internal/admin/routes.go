package admin

import (
	"net/http"

	"github.com/yohamta/dagu/internal/admin/handlers"
)

type route struct {
	method  string
	pattern string
	handler http.HandlerFunc
}

func defaultRoutes(cfg *Config) []*route {
	return []*route{
		{method: http.MethodGet, pattern: `^/?$`, handler: handlers.HandleGetList(
			&handlers.DAGListHandlerConfig{
				DAGsDir: cfg.DAGs,
			},
		)},
		{method: http.MethodGet, pattern: `^/dags/?$`, handler: handlers.HandleGetList(
			&handlers.DAGListHandlerConfig{
				DAGsDir: cfg.DAGs,
			},
		)},
		{method: http.MethodGet, pattern: `^/dags/([^/]+)$`, handler: handlers.HandleGetDAG(
			&handlers.DAGHandlerConfig{
				DAGsDir:            cfg.DAGs,
				LogEncodingCharset: cfg.LogEncodingCharset,
			},
		)},
		{method: http.MethodPost, pattern: `^/dags/([^/]+)$`, handler: handlers.HandlePostDAGAction(
			&handlers.PostDAGHandlerConfig{
				DAGsDir: cfg.DAGs,
				Bin:     cfg.Command,
				WkDir:   cfg.WorkDir,
			},
		)},
	}
}
