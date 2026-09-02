(() => {
  'use strict';

  const endpoint = window.FTN_UPSTREAM_MONITORING_URL || '/api/v1/upstream/monitoring';
  const grid = document.getElementById('upstreamMonitoringGrid');
  const summary = document.getElementById('upstreamMonitoringSummary');
  if (!grid || !summary) return;

  const esc = (value) => String(value ?? '—')
    .replaceAll('&', '&amp;').replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;').replaceAll('"', '&quot;').replaceAll("'", '&#39;');

  const stateClass = (state) => {
    const normalized = String(state || '').toLowerCase();
    if (normalized === 'healthy' || normalized === 'up') return 'healthy';
    if (normalized === 'degraded' || normalized === 'warning') return 'warning';
    if (normalized === 'down' || normalized === 'critical') return 'critical';
    return 'muted';
  };

  const normalize = (payload) => {
    const items = Array.isArray(payload) ? payload : (payload?.upstreams || payload?.items || []);
    return items.map((item) => ({
      name: item.name || item.id || item.provider || 'Unknown upstream',
      type: item.type || item.kind || 'upstream',
      state: item.state || (item.healthy ? 'healthy' : 'unknown'),
      latency: Number.isFinite(Number(item.latency_ms)) ? Number(item.latency_ms) : null,
      loss: Number.isFinite(Number(item.packet_loss_pct)) ? Number(item.packet_loss_pct) : null,
      prefixes: Number.isFinite(Number(item.prefix_count)) ? Number(item.prefix_count) : null,
      limit: Number.isFinite(Number(item.prefix_limit)) ? Number(item.prefix_limit) : null,
      checked: item.last_checked || item.last_check_timestamp || null,
    }));
  };

  const render = (items) => {
    const healthy = items.filter((x) => stateClass(x.state) === 'healthy').length;
    const degraded = items.filter((x) => stateClass(x.state) === 'warning').length;
    const critical = items.filter((x) => stateClass(x.state) === 'critical').length;
    summary.textContent = `${items.length} upstream${items.length === 1 ? '' : 's'} · ${healthy} healthy · ${degraded} degraded · ${critical} critical`;

    if (!items.length) {
      grid.innerHTML = '<article class="card muted"><h3>No upstream telemetry</h3><p>The control API has not supplied upstream monitoring evidence yet.</p></article>';
      return;
    }

    grid.innerHTML = items.map((item) => {
      const prefix = item.prefixes == null ? '—' : `${item.prefixes}${item.limit == null ? '' : ` / ${item.limit}`}`;
      return `<article class="card upstream-card">
        <div class="upstream-card-head"><div><h3>${esc(item.name)}</h3><small>${esc(item.type)}</small></div><span class="status ${stateClass(item.state)}">${esc(item.state)}</span></div>
        <div class="upstream-stats">
          <span><b>${item.latency == null ? '—' : `${esc(item.latency)} ms`}</b><small>latency</small></span>
          <span><b>${item.loss == null ? '—' : `${esc(item.loss)}%`}</b><small>packet loss</small></span>
          <span><b>${esc(prefix)}</b><small>prefixes / limit</small></span>
        </div>
        <small class="muted">Last check: ${esc(item.checked || '—')}</small>
      </article>`;
    }).join('');
  };

  async function refresh() {
    grid.innerHTML = '<article class="card muted">Refreshing upstream telemetry…</article>';
    try {
      const headers = {Accept: 'application/json'};
      const token = sessionStorage.getItem('ftn-control-token');
      if (token) headers.Authorization = `Bearer ${token}`;
      const response = await fetch(endpoint, {headers, cache: 'no-store'});
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      render(normalize(await response.json()));
    } catch (error) {
      summary.textContent = 'Telemetry unavailable';
      grid.innerHTML = `<article class="card muted"><h3>Monitoring API unavailable</h3><p>No fabricated values are shown. Endpoint: <code>${esc(endpoint)}</code></p><small>${esc(error.message)}</small></article>`;
    }
  }

  window.FTNUpstreamMonitoring = {refresh};
  refresh();
  setInterval(refresh, 15000);
})();
