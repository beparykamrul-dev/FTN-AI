export type DashboardMetric = {
  id: string;
  label: string;
  value: string;
  trend?: string;
  status: "healthy" | "warning" | "critical" | "neutral";
};

export type DashboardViewModel = {
  title: string;
  subtitle: string;
  metrics: readonly DashboardMetric[];
};

export const EMPTY_DASHBOARD: DashboardViewModel = {
  title: "FTN Control Center",
  subtitle: "Unified ISP operations",
  metrics: [
    { id: "accounts", label: "Accounts", value: "—", status: "neutral" },
    { id: "active-services", label: "Active services", value: "—", status: "neutral" },
    { id: "open-incidents", label: "Open incidents", value: "—", status: "neutral" },
    { id: "devices", label: "Devices", value: "—", status: "neutral" },
  ],
};
