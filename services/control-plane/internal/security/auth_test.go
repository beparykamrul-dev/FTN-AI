package security

import (
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
)

func TestMiddlewareOpenWhenTokenUnset(t *testing.T) {
    t.Setenv("FTN_API_AUTH_TOKEN", "")
    called := false
    h := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        called = true
        w.WriteHeader(http.StatusNoContent)
    }))
    rr := httptest.NewRecorder()
    h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/services", nil))
    if rr.Code != http.StatusNoContent || !called {
        t.Fatalf("expected open middleware, status=%d called=%v", rr.Code, called)
    }
}

func TestMiddlewareRejectsMissingBearer(t *testing.T) {
    t.Setenv("FTN_API_AUTH_TOKEN", "test-secret")
    h := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        t.Fatal("protected handler must not run")
    }))
    rr := httptest.NewRecorder()
    h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/services", nil))
    if rr.Code != http.StatusUnauthorized {
        t.Fatalf("expected 401, got %d", rr.Code)
    }
}

func TestMiddlewareRejectsEmptyBearer(t *testing.T) {
    t.Setenv("FTN_API_AUTH_TOKEN", "test-secret")
    h := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        t.Fatal("protected handler must not run")
    }))
    req := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
    req.Header.Set("Authorization", "Bearer ")
    rr := httptest.NewRecorder()
    h.ServeHTTP(rr, req)
    if rr.Code != http.StatusUnauthorized {
        t.Fatalf("expected 401 for empty bearer, got %d", rr.Code)
    }
}

func TestMiddlewareRejectsInvalidBearer(t *testing.T) {
    t.Setenv("FTN_API_AUTH_TOKEN", "test-secret")
    h := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        t.Fatal("protected handler must not run")
    }))
    req := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
    req.Header.Set("Authorization", "Bearer wrong")
    rr := httptest.NewRecorder()
    h.ServeHTTP(rr, req)
    if rr.Code != http.StatusUnauthorized {
        t.Fatalf("expected 401 for invalid bearer, got %d", rr.Code)
    }
}

func TestMiddlewareAcceptsValidBearer(t *testing.T) {
    t.Setenv("FTN_API_AUTH_TOKEN", "test-secret")
    called := false
    h := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        called = true
        w.WriteHeader(http.StatusNoContent)
    }))
    req := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
    req.Header.Set("Authorization", "Bearer test-secret")
    rr := httptest.NewRecorder()
    h.ServeHTTP(rr, req)
    if rr.Code != http.StatusNoContent || !called {
        t.Fatalf("expected authenticated request, status=%d called=%v", rr.Code, called)
    }
}

func TestMiddlewareKeepsHealthPublic(t *testing.T) {
    t.Setenv("FTN_API_AUTH_TOKEN", "test-secret")
    h := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNoContent)
    }))
    for _, path := range []string{"/healthz", "/readyz", "/metrics"} {
        rr := httptest.NewRecorder()
        h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
        if rr.Code != http.StatusNoContent {
            t.Fatalf("%s: expected public endpoint, got %d", path, rr.Code)
        }
    }
}

func TestMiddlewareDoesNotReadProcessEnvironmentDirectlyAfterConstruction(t *testing.T) {
    // This test documents that middleware construction snapshots the token.
    t.Setenv("FTN_API_AUTH_TOKEN", "first")
    h := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNoContent)
    }))
    _ = os.Setenv("FTN_API_AUTH_TOKEN", "second")
    defer os.Unsetenv("FTN_API_AUTH_TOKEN")
    req := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
    req.Header.Set("Authorization", "Bearer first")
    rr := httptest.NewRecorder()
    h.ServeHTTP(rr, req)
    if rr.Code != http.StatusNoContent {
        t.Fatalf("expected middleware to use construction-time token, got %d", rr.Code)
    }
}
