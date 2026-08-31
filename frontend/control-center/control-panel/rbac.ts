export type ControlRole =
  | "super-admin"
  | "admin"
  | "billing"
  | "noc"
  | "engineer"
  | "support";

export type PanelRoute = {
  id: string;
  label: string;
  path: string;
  roles: readonly ControlRole[];
};

export const PANEL_ROUTES: readonly PanelRoute[] = [
  { id: "overview", label: "Overview", path: "/control/overview", roles: ["super-admin", "admin", "billing", "noc", "engineer", "support"] },
  { id: "accounts", label: "Accounts", path: "/control/accounts", roles: ["super-admin", "admin", "billing"] },
  { id: "billing", label: "Billing", path: "/control/billing", roles: ["super-admin", "admin", "billing"] },
  { id: "payments", label: "Payments", path: "/control/payments", roles: ["super-admin", "admin", "billing"] },
  { id: "devices", label: "Devices", path: "/control/devices", roles: ["super-admin", "admin", "noc", "engineer"] },
  { id: "discovery", label: "Discovery", path: "/control/discovery", roles: ["super-admin", "noc", "engineer"] },
  { id: "topology", label: "Topology", path: "/control/topology", roles: ["super-admin", "noc", "engineer"] },
  { id: "monitoring", label: "Monitoring", path: "/control/monitoring", roles: ["super-admin", "noc", "engineer"] },
  { id: "incidents", label: "Incidents", path: "/control/incidents", roles: ["super-admin", "admin", "noc", "engineer"] },
  { id: "ai-operations", label: "AI Operations", path: "/control/ai-operations", roles: ["super-admin", "noc"] },
  { id: "recovery", label: "Recovery approvals", path: "/control/recovery", roles: ["super-admin", "noc"] },
  { id: "audit", label: "Audit", path: "/control/audit", roles: ["super-admin", "admin", "billing", "noc"] },
  { id: "support", label: "Support", path: "/control/support", roles: ["super-admin", "admin", "support"] },
  { id: "system", label: "System", path: "/control/system", roles: ["super-admin"] },
];

export function canAccess(role: ControlRole, route: PanelRoute): boolean {
  return route.roles.includes(role);
}

export function routesForRole(role: ControlRole): readonly PanelRoute[] {
  return PANEL_ROUTES.filter((route) => canAccess(role, route));
}
