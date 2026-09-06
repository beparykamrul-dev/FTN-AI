package realtime

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
)

const MaxBroadcastPayload = 1 << 20
const MaxRoomPeers = 100000

type Hub struct { mu sync.RWMutex; rooms map[string]map[string]Peer }
func NewHub() *Hub { return &Hub{rooms: make(map[string]map[string]Peer)} }
func (h *Hub) Join(roomID string, peer Peer) error { if h==nil||peer==nil{return fmt.Errorf("hub and peer are required")};roomID=strings.TrimSpace(roomID);peerID:=strings.TrimSpace(peer.ID());if roomID==""||peerID==""||len(roomID)>256||len(peerID)>256{return fmt.Errorf("room and peer are required")};h.mu.Lock();defer h.mu.Unlock();if h.rooms==nil{h.rooms=make(map[string]map[string]Peer)};room:=h.rooms[roomID];if room==nil{room=make(map[string]Peer);h.rooms[roomID]=room};if _,exists:=room[peerID];!exists&&len(room)>=MaxRoomPeers{return fmt.Errorf("room peer limit exceeded")};room[peerID]=peer;return nil }
func (h *Hub) Leave(roomID,peerID string){if h==nil{return};roomID=strings.TrimSpace(roomID);peerID=strings.TrimSpace(peerID);if roomID==""||peerID==""{return};h.mu.Lock();defer h.mu.Unlock();room:=h.rooms[roomID];if room==nil{return};delete(room,peerID);if len(room)==0{delete(h.rooms,roomID)}}
func (h *Hub) Broadcast(ctx context.Context,roomID,senderID string,payload []byte)map[string]error{result:=make(map[string]error);if h==nil{result["_"]=fmt.Errorf("hub is nil");return result};if ctx==nil{result["_"]=fmt.Errorf("context is nil");return result};if err:=ctx.Err();err!=nil{result["_"]=err;return result};roomID=strings.TrimSpace(roomID);senderID=strings.TrimSpace(senderID);if roomID==""||len(roomID)>256{result["_"]=fmt.Errorf("room is required");return result};if len(payload)>MaxBroadcastPayload{result["_"]=fmt.Errorf("payload exceeds %d bytes",MaxBroadcastPayload);return result};body:=append([]byte(nil),payload...);h.mu.RLock();room:=h.rooms[roomID];ids:=make([]string,0,len(room));for id:=range room{if id!=senderID{ids=append(ids,id)}};sort.Strings(ids);peers:=make([]Peer,0,len(ids));for _,id:=range ids{if p:=room[id];p!=nil{peers=append(peers,p)}};h.mu.RUnlock();for _,p:=range peers{if err:=ctx.Err();err!=nil{result[p.ID()]=err;continue};result[strings.TrimSpace(p.ID())]=p.Send(ctx,body)};return result}
