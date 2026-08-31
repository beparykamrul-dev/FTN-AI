import type { ControlRole, PanelRoute } from "./rbac";

export type PanelUser = {
  id: string;
  displayName: string;
  role: ControlRole;
};

export type PanelShellState = {
  user: PanelUser;
  activePath: string;
  sidebarOpen: boolean;
  notificationsUnread: number;
  routes: readonly PanelRoute[];
};

export function createShellState(
  user: PanelUser,
  routes: readonly PanelRoute[],
  activePath = "/control/overview",
): PanelShellState {
  return {
    user,
    activePath,
    sidebarOpen: true,
    notificationsUnread: 0,
    routes,
  };
}
