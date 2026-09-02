package main

import "testing"

func TestPrivilegedJobType(t *testing.T) {
    cases := []struct { job, action string; want bool }{
        {"network.route", "network.route", true},
        {"router.config", "router.config", true},
        {"dns.change", "dns.change", true},
        {"security.policy", "security.policy", true},
        {"config.reload", "config.reload", true},
        {"media.transcode", "", false},
        {"billing.invoice", "billing.invoice", false},
    }
    for _, tc := range cases {
        if got := privilegedJobType(tc.job, tc.action); got != tc.want { t.Errorf("privilegedJobType(%q,%q)=%v, want %v", tc.job,tc.action,got,tc.want) }
    }
}
