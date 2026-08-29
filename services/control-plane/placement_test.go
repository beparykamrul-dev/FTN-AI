package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPlacementChoosesHealthyLowLatencyNode(t *testing.T) {
	old := nodes
	nodes = []Node{
		{ID:"slow",Provider:"p1",Region:"BD",Services:[]string{"media"},CPUPercent:20,RAMPercent:20,LatencyMs:80,Healthy:true},
		{ID:"fast",Provider:"p2",Region:"BD",Services:[]string{"media"},CPUPercent:30,RAMPercent:30,LatencyMs:10,Healthy:true},
		{ID:"down",Provider:"p3",Region:"BD",Services:[]string{"media"},LatencyMs:1,Healthy:false},
	}
	defer func(){nodes=old}()
	r:=httptest.NewRequest(http.MethodPost,"/api/v1/placement/preview",strings.NewReader(`{"service_id":"media","region":"BD"}`))
	w:=httptest.NewRecorder()
	(&App{}).placement(w,r)
	if w.Code!=http.StatusOK {t.Fatalf("status=%d body=%s",w.Code,w.Body.String())}
	if !strings.Contains(w.Body.String(),`"node_id":"fast"`){t.Fatalf("fast node not selected: %s",w.Body.String())}
	if !strings.Contains(w.Body.String(),`"execution":"approval_required"`){t.Fatal("placement must remain approval-gated")}
}

func TestPlacementRejectsUnsupportedService(t *testing.T) {
	old:=nodes
	nodes=[]Node{{ID:"n1",Provider:"p1",Services:[]string{"dns"},Healthy:true}}
	defer func(){nodes=old}()
	r:=httptest.NewRequest(http.MethodPost,"/api/v1/placement/preview",strings.NewReader(`{"service_id":"media"}`))
	w:=httptest.NewRecorder()
	(&App{}).placement(w,r)
	if w.Code!=http.StatusServiceUnavailable {t.Fatalf("status=%d",w.Code)}
}
