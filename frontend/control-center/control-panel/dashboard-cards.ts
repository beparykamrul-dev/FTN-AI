import type { DashboardMetric } from "./dashboard-model";

export const DASHBOARD_CARDS: readonly DashboardMetric[] = [
  { id: "accounts", label: "Accounts", value: "—", status: "neutral" },
  { id: "active-services", label: "Active services", value: "—", status: "neutral" },
  { id: "open-incidents", label: "Open incidents", value: "—", status: "neutral" },
  { id: "devices", label: "Devices", value: "—", status: "neutral" },
  { id: "network-health", label: "Network health", value: "—", status: "neutral" },
  { id: "pending-recovery", label: "Pending recovery", value: "—", status: "neutral" },
];
