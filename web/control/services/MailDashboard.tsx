import { useEffect, useState } from "react";

type MailItem = { id: string; subject?: string; from?: string; to?: string[]; status?: string; created_at?: string };

type Props = { identityId: string };

export function MailDashboard({ identityId }: Props) {
  const [items, setItems] = useState<MailItem[]>([]);
  const [busy, setBusy] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const r = await fetch(`/api/v1/control/mail/messages?identity_id=${encodeURIComponent(identityId)}`, { credentials: "include" });
        if (!r.ok) throw new Error("mail load failed");
        const data = await r.json() as { items?: MailItem[] };
        if (active) setItems(data.items ?? []);
      } catch { if (active) setError("Mail data could not be loaded."); }
      finally { if (active) setBusy(false); }
    })();
    return () => { active = false; };
  }, [identityId]);

  if (busy) return <section aria-busy="true"><h2>FTN Mail</h2><p>Loading mail…</p></section>;
  return (
    <section aria-labelledby="ftn-mail-dashboard-title">
      <header>
        <h2 id="ftn-mail-dashboard-title">FTN Mail</h2>
        <p>Private mailbox and delivery activity.</p>
      </header>
      {error && <p role="alert">{error}</p>}
      {!error && items.length === 0 && <p>No mail activity yet.</p>}
      {items.length > 0 && <ul>
        {items.map((item) => <li key={item.id}>
          <strong>{item.subject || "(No subject)"}</strong>
          <span>{item.from || "Unknown sender"}</span>
          <small>{item.status || "unknown"}{item.created_at ? ` · ${new Date(item.created_at).toLocaleString()}` : ""}</small>
        </li>)}
      </ul>}
    </section>
  );
}
