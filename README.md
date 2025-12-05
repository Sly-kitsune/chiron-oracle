# Chiron wound inversion oracle

A minimal, global-ready web app that calculates natal Chiron placements (sign + house) and reframes traditional “wounds” into Left-Hand Path strengths. Runs as a single Go 1.22 binary with zero external dependencies, plus an optional Cloudflare Worker for always-on edge deployment.

---

## Features

- **Frontend:** Modern responsive dark theme with gradients, animations, and validation.
- **Global-ready:** Works with any latitude/longitude and supports time zones for local birth times.
- **Interpretations:** Complete mapping (12 signs × 12 houses = 144).
- **API:** REST JSON endpoints with CORS.
- **Operational:** Single binary (~8MB), low memory (<50MB), cross-platform.

---

## Quick start (Codespaces)

```bash
# 1) Enter the project
cd /workspaces/chiron-oracle

# 2) Run the server
go run main.go

# 3) Open in browser
# - Open the "Ports" tab
# - Find port 8080 (Forwarded)
# - Click the globe icon 🌐

