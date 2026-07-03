package ui

import (
	"fmt"
	"net/http"
)

var DashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>git-share</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #fff; color: #333; display: flex; justify-content: center; align-items: center; min-height: 100vh; padding: 1rem; }
.container { width: 100%; max-width: 600px; margin: 0 auto; }
.logo { text-align: center; margin-bottom: 2.5rem; }
.logo h1 { font-size: 1.75rem; font-weight: 700; color: #f05133; letter-spacing: -0.02em; }
.logo p { color: #999; font-size: 0.875rem; margin-top: 0.25rem; }
.card { background: #fff; padding: 1.25rem 0; }
.info-row { display: flex; justify-content: space-between; align-items: center; padding: 0.875rem 0; }
.label { color: #999; font-size: 0.875rem; }
.value { color: #333; font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', monospace; font-size: 0.875rem; }
.clone-box { margin-top: 1.5rem; background: #fafafa; border-radius: 8px; padding: 1rem 1.25rem; font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', monospace; font-size: 0.8125rem; word-break: break-all; color: #333; line-height: 1.5; }
.clone-box strong { color: #f05133; font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.04em; display: block; margin-bottom: 0.25rem; }
.badge { display: inline-flex; align-items: center; gap: 0.375rem; background: #fef3f0; color: #f05133; font-size: 0.75rem; font-weight: 600; padding: 0.25rem 0.75rem; border-radius: 999px; }
.badge .dot { width: 6px; height: 6px; border-radius: 50%; background: #f05133; }
.clone-box code { display: block; margin-top: 0.125rem; }
</style>
</head>
<body>
<div class="container" id="app">
<div class="logo">
<h1>git share</h1>
<p>Repository sharing is active</p>
</div>
<div class="card">
<div class="info-row"><span class="label">Repository</span><span class="value" id="repo">-</span></div>
<div class="info-row"><span class="label">Branch</span><span class="value" id="branch">-</span></div>
<div class="info-row"><span class="label">Mode</span><span class="value"><span class="badge"><span class="dot"></span><span id="mode">-</span></span></span></div>
<div class="info-row"><span class="label">Clients</span><span class="value" id="clients">0</span></div>
<div class="info-row"><span class="label">Uptime</span><span class="value" id="uptime">-</span></div>
<div class="clone-box"><strong>Clone</strong><code id="cloneUrl">git clone -</code></div>
</div>
</div>
<script>
async function loadInfo() {
try {
const resp = await fetch('/info');
const data = await resp.json();
document.getElementById('repo').textContent = data.repository;
document.getElementById('branch').textContent = data.branch;
document.getElementById('mode').textContent = data.readonly ? 'Read Only' : 'Read / Write';
document.getElementById('cloneUrl').textContent = 'git clone ' + (data.lan_url || data.clone_url);
} catch(e) { document.getElementById('repo').textContent = 'Error loading'; }
}
loadInfo();
setInterval(loadInfo, 5000);
</script>
</body>
</html>`

func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, DashboardHTML)
}
