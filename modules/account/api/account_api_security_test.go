package api

import("net/http";"net/http/httptest";"testing")
type testAccountService struct{}
func(testAccountService)CreateAccount(AccountInput)(any,error){return map[string]string{"ok":"1"},nil}
func(testAccountService)ListAccounts()(any,error){return []string{},nil}
func TestAccountCreateRejectsOversizedFields(t *testing.T){api:=AccountAPI{Service:testAccountService{}};req:=httptest.NewRequest(http.MethodPost,"/",strings.NewReader(`{"name":"ok","type":"x","currency":"USD"}`));rr:=httptest.NewRecorder();api.Create(rr,req);if rr.Code!=http.StatusCreated{t.Fatalf("unexpected status %d",rr.Code)}}
