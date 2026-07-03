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
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #0d1117; color: #c9d1d9; display: flex; justify-content: center; align-items: center; min-height: 100vh; }
.container { max-width: 600px; width: 100%; padding: 2rem; }
.card { background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 2rem; margin-bottom: 1rem; }
h1 { font-size: 1.5rem; margin-bottom: 0.5rem; color: #58a6ff; }
.subtitle { color: #8b949e; margin-bottom: 1.5rem; }
.info-row { display: flex; justify-content: space-between; padding: 0.75rem 0; border-bottom: 1px solid #21262d; }
.info-row:last-child { border-bottom: none; }
.label { color: #8b949e; }
.value { color: #c9d1d9; font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', monospace; }
.clone-url { background: #0d1117; border: 1px solid #30363d; border-radius: 6px; padding: 0.75rem 1rem; font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', monospace; margin: 1rem 0; word-break: break-all; color: #58a6ff; }
.status { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-right: 0.5rem; }
.status.online { background: #3fb950; }
.status.offline { background: #f85149; }
</style>
</head>
<body>
<div class="container" id="app">
<div class="card">
<h1>git-share</h1>
<p class="subtitle">Repository sharing is active</p>
<div class="info-row"><span class="label">Repository</span><span class="value" id="repo">-</span></div>
<div class="info-row"><span class="label">Branch</span><span class="value" id="branch">-</span></div>
<div class="info-row"><span class="label">Status</span><span class="value"><span class="status online"></span>Running</span></div>
<div class="info-row"><span class="label">Mode</span><span class="value" id="mode">-</span></div>
<div class="info-row"><span class="label">Clients</span><span class="value" id="clients">0</span></div>
<div class="info-row"><span class="label">Uptime</span><span class="value" id="uptime">-</span></div>
<div class="clone-url" id="cloneUrl">-</div>
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
document.getElementById('cloneUrl').textContent = data.clone_url || data.lan_url;
document.getElementById('cloneUrl').innerHTML = '<strong>Clone:</strong> git clone ' + (data.lan_url || data.clone_url);
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
