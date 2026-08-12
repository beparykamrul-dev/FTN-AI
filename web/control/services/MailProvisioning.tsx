import { useState } from "react";

type Props = { identityId: string };

export function MailProvisioning({ identityId }: Props) {
  const [localPart, setLocalPart] = useState("");
  const [busy, setBusy] = useState(false);
  const [message, setMessage] = useState("");

  async function provision() {
    setBusy(true);
    setMessage("");
    try {
      const response = await fetch("/api/v1/control/mail/provision", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ identity_id: identityId, local_part: localPart }),
      });
      if (!response.ok) throw new Error("provision failed");
      setMessage(`Mailbox ${localPart}@familytimenet.com provisioned.`);
      setLocalPart("");
    } catch {
      setMessage("Mailbox could not be provisioned. Check permissions and the mail service.");
    } finally {
      setBusy(false);
    }
  }

  return (
    <section aria-labelledby="mail-provision-title">
      <h2 id="mail-provision-title">FTN Mail</h2>
      <p>Create a private mailbox for an existing FTN identity.</p>
      <label>
        Mailbox name
        <input
          value={localPart}
          onChange={(e) => setLocalPart(e.target.value.toLowerCase())}
          placeholder="name"
          pattern="[a-z0-9._-]+"
          autoComplete="off"
        />
      </label>
      <span>@familytimenet.com</span>
      <button disabled={busy || !localPart} onClick={provision}>
        {busy ? "Provisioning…" : "Create mailbox"}
      </button>
      {message && <p role="status">{message}</p>}
    </section>
  );
}
