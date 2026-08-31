import type { ControlRole, PanelRoute } from "./rbac";
import type { DashboardViewModel } from "./dashboard-model";
import type { SidebarSection } from "./sidebar-model";

export type ControlPanelContext = {
  role: ControlRole;
  activePath: string;
  routes: readonly PanelRoute[];
  sidebar: readonly SidebarSection[];
  dashboard: DashboardViewModel;
};

export type NavigationAction = {
  path: string;
  replace?: boolean;
};

export function createNavigationAction(path: string, replace = false): NavigationAction {
  return { path, replace };
}
