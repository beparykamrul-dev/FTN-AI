import { canAccess, type ControlRole, type PanelRoute } from "./rbac";

export type RouteDecision =
  | { allowed: true; route: PanelRoute }
  | { allowed: false; reason: "forbidden" | "not-found" };

export function authorizePanelRoute(role: ControlRole, path: string): RouteDecision {
  const route = findRoute(path);
  if (!route) return { allowed: false, reason: "not-found" };
  if (!canAccess(role, route)) return { allowed: false, reason: "forbidden" };
  return { allowed: true, route };
}

function findRoute(path: string): PanelRoute | undefined {
  const normalized = path.replace(/\/$/, "") || "/";
  return (require("./rbac") as typeof import("./rbac")).PANEL_ROUTES.find((route) => route.path === normalized);
}
