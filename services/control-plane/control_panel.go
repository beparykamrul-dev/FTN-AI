package main

import (
	"net/http"
)

func (a *App) controlPanel(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(controlPanelHTML))
}

const controlPanelHTML = `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>FTN Control Panel</title>
<style>
:root{color-scheme:dark}*{box-sizing:border-box}body{margin:0;font-family:system-ui,-apple-system,sans-serif;background:#07111f;color:#e9f1ff}main{max-width:1280px;margin:auto;padding:24px}.top{display:flex;justify-content:space-between;gap:16px;align-items:center}.muted{color:#91a5bf}.grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(210px,1fr));gap:14px}.card{background:#0e1b2d;border:1px solid #223750;border-radius:14px;padding:18px}.value{font-size:28px;font-weight:700;margin-top:6px}.ok{color:#55d99b}.warn{color:#ffd166}.critical{color:#ff7b7b}.bar{height:7px;background:#1c2d43;border-radius:8px;overflow:hidden;margin-top:12px}.bar>i{display:block;height:100%;width:0;background:#55d99b}.row{display:flex;gap:8px;flex-wrap:wrap}.pill{padding:5px 9px;border-radius:999px;background:#172a43;font-size:12px}.finding{border-left:3px solid #ffd166;padding-left:10px;margin:12px 0}.finding.critical{border-color:#ff7b7b}.btn{background:#18314f;color:#fff;border:1px solid #315274;border-radius:9px;padding:9px 12px;cursor:pointer}.btn:hover{background:#224367}pre{white-space:pre-wrap;color:#b9c9dd}
</style></head>
<body><main>
<div class="top"><div><h1>FTN Control Panel</h1><p class="muted">Unified API + AI decision-support control plane</p></div><button class="btn" onclick="loadAll()">Refresh</button></div>
<section class="grid" id="summary"><article class="card"><span class="muted">API</span><div class="value">Loading</div></article></section>
<h2>Infrastructure Services</h2><section class="grid" id="services"></section>
<h2>FTN AI Core</h2><section class="card"><div class="row"><span class="pill" id="aiStatus">checking</span><span class="pill">decision-support only</span><span class="pill">approval-first</span></div><p class="muted">Analysis uses supplied telemetry signals. Recommendations never execute device changes directly.</p>
<div class="row"><button class="btn" onclick="analyze({packet_loss_pct:6.2,latency_ms:118,cpu_pct:61,memory_pct:72,bgp_established_pct:100,dns_error_rate_pct:1.1})">Analyze sample telemetry</button></div><div id="findings"></div></section>
<h2>Identity</h2><section class="card" id="identity">Loading…</section>
</main>
<script>
const get=async p=>{const r=await fetch(p,{headers:{Accept:'application/json'}});if(!r.ok)throw new Error(r.status);return r.json()};
async function loadAll(){
 try{const [root,me,dash,services,ai]=await Promise.all([get('/api/v1'),get('/api/v1/me'),get('/api/v1/dashboard'),get('/api/v1/services'),get('/api/v1/ai/status')]);
 document.getElementById('summary').innerHTML='<article class="card"><span class="muted">API</span><div class="value ok">'+root.status+'</div><span class="pill">'+root.version+'</span></article><article class="card"><span class="muted">Services</span><div class="value">'+dash.services.total+'</div><span class="pill">available '+dash.services.available+'</span></article><article class="card"><span class="muted">AI</span><div class="value ok">'+ai.status+'</div><span class="pill">'+ai.mode+'</span></article><article class="card"><span class="muted">Controls</span><div class="value">'+(dash.control_plane.approval_first?'Approval-first':'Review')+'</div><span class="pill">audit enabled</span></article>';
 document.getElementById('services').innerHTML=services.services.map(s=>'<article class="card"><strong>'+s.name+'</strong><div class="ok">'+s.status+'</div><div class="row">'+s.platforms.map(p=>'<span class="pill">'+p+'</span>').join('')+'</div></article>').join('');
 document.getElementById('aiStatus').textContent=ai.status+' • '+ai.mode;document.getElementById('identity').innerHTML='<pre>'+JSON.stringify(me,null,2)+'</pre>';
 }catch(e){document.getElementById('summary').innerHTML='<article class="card"><div class="critical">Control API unavailable</div><p class="muted">'+e.message+'</p></article>'}}
async function analyze(signals){const box=document.getElementById('findings');box.innerHTML='Analyzing…';try{const r=await fetch('/api/v1/ai/analyze',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({scope:'network',signals})});const d=await r.json();box.innerHTML='<p><strong>Risk:</strong> <span class="'+d.risk+'">'+d.risk+'</span></p>'+(d.findings.length?d.findings.map(f=>'<div class="finding '+f.severity+'"><strong>'+f.code+'</strong> — '+f.summary+'<br><span class="muted">'+f.evidence+' • confidence '+f.confidence+'</span></div>').join(''):'<p class="ok">No rule-based anomaly detected.</p>')+'<pre>'+JSON.stringify(d.recommendations,null,2)+'</pre>'}catch(e){box.innerHTML='<div class="critical">AI analysis failed: '+e.message+'</div>'}}
loadAll();
</script></body></html>`
