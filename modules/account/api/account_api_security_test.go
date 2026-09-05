package api

import("net/http";"net/http/httptest";"strings";"testing")
type testAccountService struct{}
func(testAccountService)CreateAccount(AccountInput)(any,error){return map[string]string{"ok":"1"},nil}
func(testAccountService)ListAccounts()(any,error){return []string{},nil}
func TestAccountCreateRejectsOversizedFields(t *testing.T){api:=AccountAPI{Service:testAccountService{}};req:=httptest.NewRequest(http.MethodPost,"/",strings.NewReader(`{"name":"`+strings.Repeat("x",257)+`","type":"x","currency":"USD"}`));rr:=httptest.NewRecorder();api.Create(rr,req);if rr.Code!=http.StatusBadRequest{t.Fatalf("unexpected status %d",rr.Code)}}
func TestAccountCreateAcceptsValidInput(t *testing.T){api:=AccountAPI{Service:testAccountService{}};req:=httptest.NewRequest(http.MethodPost,"/",strings.NewReader(`{"name":"ok","type":"customer","currency":"USD"}`));rr:=httptest.NewRecorder();api.Create(rr,req);if rr.Code!=http.StatusCreated{t.Fatalf("unexpected status %d",rr.Code)}}
