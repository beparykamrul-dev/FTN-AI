package control
import("errors";"strings";"sync";"time")
type ApprovalStatus string
const(Pending ApprovalStatus="pending";Approved ApprovalStatus="approved";Rejected ApprovalStatus="rejected")
type Request struct{ID string `json:"id"`;ServerID string `json:"server_id"`;Operation Operation `json:"operation"`;Reason string `json:"reason"`;Status ApprovalStatus `json:"status"`;CreatedAt time.Time `json:"created_at"`;UpdatedAt time.Time `json:"updated_at"`}
type ApprovalStore struct{mu sync.RWMutex;m map[string]Request}
func NewApprovalStore()*ApprovalStore{return &ApprovalStore{m:make(map[string]Request)}}
func(s *ApprovalStore)Create(r Request)error{if s==nil{return errors.New("approval store is required")};r.ID=strings.TrimSpace(r.ID);r.ServerID=strings.TrimSpace(r.ServerID);r.Reason=strings.TrimSpace(r.Reason);if r.ID==""||r.ServerID==""||r.Operation==""{return errors.New("invalid approval request")};if err:=Validate(r.Operation);err!=nil{return err};s.mu.Lock();defer s.mu.Unlock();if s.m==nil{s.m=make(map[string]Request)};if _,exists:=s.m[r.ID];exists{return errors.New("approval request already exists")};now:=time.Now().UTC();r.Status=Pending;r.CreatedAt=now;r.UpdatedAt=now;s.m[r.ID]=r;return nil}
func(s *ApprovalStore)SetStatus(id string,status ApprovalStatus)error{if s==nil{return errors.New("approval store is required")};id=strings.TrimSpace(id);status=ApprovalStatus(strings.ToLower(strings.TrimSpace(string(status))));if id==""||(status!=Approved&&status!=Rejected){return errors.New("invalid approval status")};s.mu.Lock();defer s.mu.Unlock();r,ok:=s.m[id];if !ok{return errors.New("approval request not found")};if r.Status!=Pending{return errors.New("approval request is already decided")};r.Status=status;r.UpdatedAt=time.Now().UTC();s.m[id]=r;return nil}
func(s *ApprovalStore)Get(id string)(Request,bool){if s==nil{return Request{},false};id=strings.TrimSpace(id);if id==""{return Request{},false};s.mu.RLock();defer s.mu.RUnlock();r,ok:=s.m[id];return r,ok}
