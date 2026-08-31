package main

import (
    "testing"
    "time"
)

func TestValidateRouteSafety(t *testing.T) {
    base := RouteSafetyInput{Prefix:"203.0.113.0/24", RPKIStatus:"valid", ASPathLength:2, MaxASPathLength:5, BFDUp:true, LastUpdate:time.Now(), MaxStale:time.Minute}
    cases := []struct{name, mutate, want string}{
        {"valid", func(){}, "ok"},
        {"rpki", func(){base.RPKIStatus="invalid"}, "rpki_not_valid"},
        {"aspath", func(){base.RPKIStatus="valid"; base.ASPathLength=9}, "as_path_too_long"},
        {"bfd", func(){base.ASPathLength=2; base.BFDUp=false}, "bfd_down"},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T){
            x := base; tc.mutate = tc.mutate
            switch tc.name { case "rpki": x.RPKIStatus="invalid"; case "aspath": x.ASPathLength=9; case "bfd": x.BFDUp=false }
            ok, reason := validateRouteSafety(x)
            if tc.want == "ok" && !ok { t.Fatalf("want valid, got %s", reason) }
            if tc.want != "ok" && (ok || reason != tc.want) { t.Fatalf("got ok=%v reason=%s want=%s", ok, reason, tc.want) }
        })
    }
}
