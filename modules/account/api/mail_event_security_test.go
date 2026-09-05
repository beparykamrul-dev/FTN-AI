package api

import("net/http";"net/http/httptest";"strings";"testing")
type testMailProcessor struct{}
func(testMailProcessor)Process(MailEvent)error{return nil}
func TestMailEventRejectsTrailingJSON(t *testing.T){api:=MailEventAPI{Processor:testMailProcessor{}};req:=httptest.NewRequest(http.MethodPost,"/",strings.NewReader(`{"id":"1","source":"smtp","event_type":"accepted","confidence":1}{}`));rr:=httptest.NewRecorder();api.Process(rr,req);if rr.Code!=http.StatusBadRequest{t.Fatalf("status=%d",rr.Code)}}
