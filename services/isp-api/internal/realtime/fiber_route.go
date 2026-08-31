package realtime

import (
	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes defines the WebSocket boundary. Authentication and origin
// policy are intentionally injected by the application so deployments can
// choose their trusted identity provider without weakening this package.
func RegisterRoutes(app *fiber.App, authenticate fiber.Handler, upgrade fiber.Handler) {
	app.Use("/ws/v1", authenticate)
	app.Get("/ws/v1", upgrade)
}
