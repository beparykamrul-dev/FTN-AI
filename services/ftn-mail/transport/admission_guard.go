package transport

import (
	"net"
	"time"
)

type AdmissionGuard struct {
	Admission *IPAdmission
	Throttle  *AuthThrottle
}

func (g *AdmissionGuard) Allow(ip net.IP, now time.Time) bool {
	if g == nil || g.Admission == nil { return false }
	return g.Admission.Allow(ip, now)
}

func (g *AdmissionGuard) AuthFailure(ip net.IP, now time.Time) bool {
	if g == nil || g.Throttle == nil { return false }
	return g.Throttle.Failed(ip, now)
}

func (g *AdmissionGuard) AuthSuccess(ip net.IP) {
	if g == nil || g.Throttle == nil { return }
	g.Throttle.Success(ip)
}

func (g *AdmissionGuard) Cleanup(now time.Time) {
	if g == nil { return }
	if g.Admission != nil { g.Admission.Cleanup(now) }
	if g.Throttle != nil { g.Throttle.Cleanup(now) }
}
