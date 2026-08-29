export const CONTROL_ENDPOINTS = {
  dashboard: "/api/v1/control/dashboard",
  accounts: "/api/v1/accounts",
  invoices: "/api/v1/billing/invoices",
  payments: "/api/v1/payments",
  devices: "/api/v1/network/devices",
  discovery: "/api/v1/network/discovery",
  topology: "/api/v1/network/topology",
  incidents: "/api/v1/monitoring/incidents",
  recovery: "/api/v1/ai/recovery",
  audit: "/api/v1/audit/events",
} as const;
