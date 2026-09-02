(() => {
  'use strict';
  const tokenKey = 'ftn-control-token';
  const $ = id => document.getElementById(id);
  const esc = v => String(v ?? '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
  const token = () => sessionStorage.getItem(tokenKey) || '';
  async function api(path) {
    const headers = {Accept:'application/json'};
    if (token()) headers.Authorization = `Bearer ${token()}`;
    const r = await fetch(path, {headers, cache:'no-store'});
    if (!r.ok) throw new Error(`${r.status} ${path}`);
    return r.json();
  }
  function setConnection(ok, text) { const e=$('connection'); e.textContent=text; e.className=`status ${ok?'online':'offline'}`; }
  function renderServices(data, ent) {
    const active = new Set((ent?.entitlements || []).filter(x=>x.active).map(x=>x.service_id));
    const services = data.services || [];
    $('serviceCount').textContent = services.length;
    $('servicesGrid').innerHTML = services.map(s => `<article class="service ${active.has(s.id)?'active':''}"><div class="service-top"><h3>${esc(s.name)}</h3><span>${esc(s.status)}</span></div><p>${(s.platforms||[]).map(p=>`<span class="tag">${esc(p)}</span>`).join('')}</p><small>${active.has(s.id)?'Entitled for this principal':'Catalogued service'}</small></article>`).join('') || '<article class="card muted">No services returned.</article>';
  }
  async function load() {
    $('updated').textContent = `Refreshing ${new Date().toLocaleTimeString()}`;
    try { const h=await fetch('/healthz',{cache:'no-store'}); if(!h.ok) throw new Error(); const d=await h.json(); setConnection(true,'API: online'); $('apiStatus').textContent='Healthy'; $('apiDetail').textContent=new Date(d.time).toLocaleTimeString(); }
    catch { setConnection(false,'API: offline'); $('apiStatus').textContent='Offline'; $('apiDetail').textContent='Health check failed'; }
    try { const d=await api('/api/v1/services'); let e=null; try { e=await api('/api/v1/entitlements'); } catch {} renderServices(d,e); } catch { $('servicesGrid').innerHTML='<article class="card">Services require an authenticated session.</article>'; $('serviceCount').textContent='—'; }
    try { const d=await api('/api/v1/network/health'); $('networkGrid').innerHTML=`<article class="card"><h3>Network health</h3><pre>${esc(JSON.stringify(d,null,2))}</pre></article>`; } catch { $('networkGrid').innerHTML='<article class="card muted">Network health is protected or unavailable.</article>'; }
    try { const d=await api('/api/v1/dns-guard/summary'); $('dnsStatus').textContent='Active'; $('dnsDetail').textContent=`${d.profiles ?? 0} profiles · ${d.active_bindings ?? 0} bindings · ${d.events_24h ?? 0} events/24h`; $('dnsGrid').innerHTML=`<article class="card"><h3>Policy summary</h3><div class="big">${esc(d.profiles ?? 0)}</div><small>profiles</small></article><article class="card"><h3>Active bindings</h3><div class="big">${esc(d.active_bindings ?? 0)}</div><small>bindings</small></article><article class="card"><h3>Recent events</h3><div class="big">${esc(d.events_24h ?? 0)}</div><small>last 24h</small></article>`; } catch { $('dnsStatus').textContent='Protected'; $('dnsDetail').textContent='Authenticated API required'; $('dnsGrid').innerHTML='<article class="card muted">DNS Guard summary requires an authenticated session.</article>'; }
    try { const d=await api('/api/v1/data-governor/assets'); const n=Array.isArray(d.assets)?d.assets.length:0; $('governanceGrid').innerHTML=`<article class="card"><h3>Governed assets</h3><div class="big">${n}</div><small>returned by policy API</small></article><article class="card"><h3>Raw secrets</h3><div class="big">Never</div><small>not rendered by this UI</small></article>`; } catch { $('governanceGrid').innerHTML='<article class="card muted">Governance state requires an authenticated session.</article>'; }
    $('updated').textContent = `Updated ${new Date().toLocaleString()}`;
  }
  $('refresh').addEventListener('click', load);
  $('saveToken').addEventListener('click', () => { const v=$('token').value.trim(); if(v) sessionStorage.setItem(tokenKey,v); else sessionStorage.removeItem(tokenKey); $('token').value=''; load(); });
  $('clearToken').addEventListener('click', () => { sessionStorage.removeItem(tokenKey); $('token').value=''; load(); });
  load();
})();
