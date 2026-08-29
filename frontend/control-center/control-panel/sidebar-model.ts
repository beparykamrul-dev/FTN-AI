import { routesForRole, type ControlRole, type PanelRoute } from "./rbac";

export type SidebarSection = {
  title: string;
  items: readonly PanelRoute[];
};

export function sidebarForRole(role: ControlRole): readonly SidebarSection[] {
  const routes = routesForRole(role);
  const groups: Record<string, PanelRoute[]> = {
    Operations: [],
    Network: [],
    Finance: [],
    Administration: [],
  };

  for (const route of routes) {
    if (["devices", "discovery", "topology", "monitoring", "incidents", "ai-operations", "recovery"].includes(route.id)) groups.Network.push(route);
    else if (["billing", "payments", "accounts"].includes(route.id)) groups.Finance.push(route);
    else if (["audit", "system"].includes(route.id)) groups.Administration.push(route);
    else groups.Operations.push(route);
  }

  return Object.entries(groups)
    .filter(([, items]) => items.length > 0)
    .map(([title, items]) => ({ title, items }));
}
