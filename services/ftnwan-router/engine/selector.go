package engine

import (
	"context"
	"errors"
	"sort"

	router "github.com/beparykamrul-dev/FTN-AI/services/ftnwan-router/contracts"
)

var ErrNoRoute = errors.New("no usable route")

// SelectRoute deterministically selects the best currently-installed route.
// Lower metric wins; protocol and next-hop provide stable tie-breakers.
func SelectRoute(ctx context.Context, routes []router.Route) (router.Route, error) {
	if err := ctx.Err(); err != nil {
		return router.Route{}, err
	}
	usable := make([]router.Route, 0, len(routes))
	for _, r := range routes {
		if r.Installed && r.Prefix != "" && r.Interface != "" {
			usable = append(usable, r)
		}
	}
	if len(usable) == 0 {
		return router.Route{}, ErrNoRoute
	}
	sort.SliceStable(usable, func(i, j int) bool {
		if usable[i].Metric != usable[j].Metric {
			return usable[i].Metric < usable[j].Metric
		}
		if usable[i].Protocol != usable[j].Protocol {
			return usable[i].Protocol < usable[j].Protocol
		}
		if usable[i].NextHop != usable[j].NextHop {
			return usable[i].NextHop < usable[j].NextHop
		}
		return usable[i].Interface < usable[j].Interface
	})
	return usable[0], nil
}
