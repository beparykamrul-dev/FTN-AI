package controlplane

import("strings";"sync";"time")
type SessionState string
const(SessionConnecting SessionState="connecting";SessionAuthenticated SessionState="authenticated";SessionClosed SessionState="closed")
type WSSession struct{ID string `json:"id"`;Principal string `json:"principal"`;State SessionState `json:"state"`;ConnectedAt time.Time `json:"connected_at"`;LastSeen time.Time `json:"last_seen"`}
type SessionRegistry struct{mu sync.RWMutex;sessions map[string]WSSession}
func NewSessionRegistry()*SessionRegistry{return &SessionRegistry{sessions:make(map[string]WSSession)}}
func(r *SessionRegistry)Upsert(s WSSession){if r==nil{return};s.ID=strings.TrimSpace(s.ID);s.Principal=strings.TrimSpace(s.Principal);if s.ID==""{return};if s.ConnectedAt.IsZero(){s.ConnectedAt=time.Now().UTC()}else{s.ConnectedAt=s.ConnectedAt.UTC()};if s.LastSeen.IsZero(){s.LastSeen=s.ConnectedAt}else{s.LastSeen=s.LastSeen.UTC()};if s.State==""{s.State=SessionConnecting};r.mu.Lock();if r.sessions==nil{r.sessions=make(map[string]WSSession)};r.sessions[s.ID]=s;r.mu.Unlock()}
func(r *SessionRegistry)Get(id string)(WSSession,bool){if r==nil{return WSSession{},false};id=strings.TrimSpace(id);if id==""{return WSSession{},false};r.mu.RLock();defer r.mu.RUnlock();s,ok:=r.sessions[id];return s,ok}
func(r *SessionRegistry)Remove(id string){if r==nil{return};id=strings.TrimSpace(id);if id==""{return};r.mu.Lock();delete(r.sessions,id);r.mu.Unlock()}
