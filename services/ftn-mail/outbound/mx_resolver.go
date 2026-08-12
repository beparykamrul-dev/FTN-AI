package outbound

import (
	"context"
	"errors"
	"net"
	"sort"
)

type MXResolver interface { LookupMX(context.Context, string) ([]*net.MX, error) }

type NetMXResolver struct { Resolver *net.Resolver }

func (r NetMXResolver) LookupMX(ctx context.Context, domain string) ([]*net.MX, error) {
	if domain == "" { return nil, errors.New("empty mail domain") }
	resolver := r.Resolver
	if resolver == nil { resolver = net.DefaultResolver }
	records, err := resolver.LookupMX(ctx, domain)
	if err != nil { return nil, err }
	sort.SliceStable(records, func(i,j int) bool { return records[i].Pref < records[j].Pref })
	return records, nil
}
