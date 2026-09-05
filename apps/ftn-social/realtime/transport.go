package realtime

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Peer interface { ID() string; Send(context.Context, []byte) error; Close() error }
type Transport struct { mu sync.RWMutex; peers map[string]Peer }
const MaxTransportPayload = 1 << 20
func NewTransport()*Transport{return &Transport{peers:make(map[string]Peer)}}
func(t *Transport)Register(peer Peer)error{if t==nil||peer==nil{return fmt.Errorf("invalid peer")};id:=strings.TrimSpace(peer.ID());if id==""{return fmt.Errorf("peer id is required")};t.mu.Lock();defer t.mu.Unlock();if t.peers==nil{t.peers=make(map[string]Peer)};if _,ok:=t.peers[id];ok{return fmt.Errorf("peer already registered: %s",id)};t.peers[id]=peer;return nil}
func(t *Transport)Unregister(id string)error{if t==nil{return nil};id=strings.TrimSpace(id);if id==""{return nil};t.mu.Lock();peer,ok:=t.peers[id];if ok{delete(t.peers,id)};t.mu.Unlock();if !ok||peer==nil{return nil};return peer.Close()}
func(t *Transport)Send(ctx context.Context,id string,payload []byte)error{if t==nil{return fmt.Errorf("transport is nil")};if ctx==nil{return fmt.Errorf("context is nil")};id=strings.TrimSpace(id);if id==""{return fmt.Errorf("peer id is required")};if len(payload)==0||len(payload)>MaxTransportPayload{return fmt.Errorf("payload size out of bounds")};select{case<-ctx.Done():return ctx.Err();default:};t.mu.RLock();peer,ok:=t.peers[id];t.mu.RUnlock();if !ok||peer==nil{return fmt.Errorf("peer not connected: %s",id)};return peer.Send(ctx,payload)}
