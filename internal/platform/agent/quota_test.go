package agent

import (
	"context"
	"testing"
	"time"
)

type testQuotaStore struct{usage Usage}
func(s *testQuotaStore)Get(context.Context,Scope)(Usage,error){return s.usage,nil}
func(s *testQuotaStore)Put(context.Context,Scope,Usage)error{return nil}

func TestQuotaGateRejectsNegativeStoredUsage(t *testing.T){store:=&testQuotaStore{usage:Usage{Requests:-1}};q:=&QuotaGate{Store:store};plan:=Plans["free"];if err:=q.CheckAndConsume(context.Background(),Scope{},plan,1);err==nil{t.Fatal("expected invalid stored usage to fail")}}
func TestQuotaGateResetsExpiredUsage(t *testing.T){store:=&testQuotaStore{usage:Usage{Requests:20,Tokens:20000,ResetAt:time.Now().Add(-time.Hour)}};q:=&QuotaGate{Store:store};plan:=Plans["free"];if err:=q.CheckAndConsume(context.Background(),Scope{},plan,1);err!=nil{t.Fatalf("expected expired quota to reset: %v",err)}}
