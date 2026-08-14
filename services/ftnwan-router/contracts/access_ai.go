package router

import "context"

// AccessAI provides a read/diagnose/plan boundary for router-facing assistants.
// Implementations must not bypass router policy or execute arbitrary commands.
type AccessAI interface {
	Explain(context.Context, string) (string, error)
	Diagnose(context.Context, string) (string, error)
	Recommend(context.Context, string) (string, error)
}

// AccessSurface identifies a supported FTN access surface.
type AccessSurface string

const (
	SurfaceLocalRouter AccessSurface = "local-router"
	SurfaceAndroid     AccessSurface = "android"
	SurfacePC          AccessSurface = "pc"
	SurfaceVPN         AccessSurface = "vpn"
	SurfaceMesh        AccessSurface = "mesh"
)
