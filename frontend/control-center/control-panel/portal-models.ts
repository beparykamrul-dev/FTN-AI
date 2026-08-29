export type AccountRow = {
  id: string;
  name: string;
  status: "active" | "suspended" | "pending";
  serviceCount: number;
};

export type InvoiceRow = {
  id: string;
  accountId: string;
  amountMinor: number;
  currency: string;
  status: "draft" | "open" | "paid" | "overdue" | "void";
};

export type PaymentRow = {
  id: string;
  accountId: string;
  amountMinor: number;
  currency: string;
  status: "pending" | "succeeded" | "failed" | "refunded";
};

export type NocIncidentRow = {
  id: string;
  severity: "info" | "warning" | "critical";
  title: string;
  status: "open" | "acknowledged" | "resolved";
  occurredAt: string;
};
