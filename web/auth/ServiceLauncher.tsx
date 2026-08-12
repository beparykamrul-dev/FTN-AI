import { useEffect, useState } from "react";
import { getSession, signIn, signOut, type FTNSession } from "./client";

export function ServiceLauncher() {
  const [session, setSession] = useState<FTNSession | null>(null);
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getSession().then(setSession).catch(() => {}).finally(() => setLoading(false));
  }, []);

  async function submit(event: React.FormEvent) {
    event.preventDefault();
    setError("");
    try {
      await signIn(login, password);
      setPassword("");
      setSession(await getSession());
    } catch {
      setError("Login failed. Check your credentials or service access.");
    }
  }

  if (loading) return <main className="ftn-auth"><p>Loading FTN…</p></main>;

  if (!session) {
    return (
      <main className="ftn-auth">
        <section className="ftn-auth-card">
          <header><strong>FTN</strong><span>Private Services</span></header>
          <h1>Welcome back</h1>
          <p>Sign in with your FTN Identity.</p>
          <form onSubmit={submit}>
            <label>Username or email<input autoComplete="username" value={login} onChange={e => setLogin(e.target.value)} required /></label>
            <label>Password<input type="password" autoComplete="current-password" value={password} onChange={e => setPassword(e.target.value)} required /></label>
            {error && <p role="alert">{error}</p>}
            <button type="submit">Sign in</button>
          </form>
        </section>
      </main>
    );
  }

  return (
    <main className="ftn-auth">
      <section className="ftn-launcher">
        <header><strong>FTN</strong><span>Your Services</span></header>
        <div className="ftn-service-grid">
          {session.services.filter(s => s.enabled).map(service => (
            <a className="ftn-service-card" href={`/services/${service.id}`} key={service.id}>
              <strong>{service.name}</strong>
              <span>Open service →</span>
            </a>
          ))}
        </div>
        <button className="ftn-signout" onClick={async () => { await signOut(); setSession(null); }}>Sign out</button>
      </section>
    </main>
  );
}
