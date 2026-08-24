package edge

import (
    "context"
    "fmt"
)

type ClientPlatform string
const (
    PlatformAndroid ClientPlatform = "android"
    PlatformWindows ClientPlatform = "windows"
    PlatformLinux ClientPlatform = "linux"
    PlatformMacOS ClientPlatform = "macos"
)

type ClientEnrollment struct {
    DeviceID string `json:"device_id"`
    UserID string `json:"user_id"`
    Platform ClientPlatform `json:"platform"`
    NodeID string `json:"node_id"`
    PolicyVersion string `json:"policy_version"`
    CertificateID string `json:"certificate_id"`
}

type ClientEnrollmentStore interface {
    Enroll(context.Context, ClientEnrollment) error
    Get(context.Context, string) (ClientEnrollment, error)
    Sync(context.Context, string) error
    Revoke(context.Context, string) error
}

func ValidateClientEnrollment(e ClientEnrollment) error {
    if e.DeviceID == "" || e.UserID == "" || e.NodeID == "" { return fmt.Errorf("device enrollment is incomplete") }
    switch e.Platform { case PlatformAndroid, PlatformWindows, PlatformLinux, PlatformMacOS: default: return fmt.Errorf("unsupported client platform") }
    if e.CertificateID == "" { return fmt.Errorf("certificate_id is required") }
    return nil
}
