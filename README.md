# gatus-admin

[![Release](https://github.com/onatakduman/gatus-admin/actions/workflows/release.yml/badge.svg)](https://github.com/onatakduman/gatus-admin/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![GHCR](https://img.shields.io/badge/ghcr.io-onatakduman%2Fgatus--admin-blue?logo=github)](https://github.com/onatakduman/gatus-admin/pkgs/container/gatus-admin)
[![Docker Hub](https://img.shields.io/badge/docker-iamonat%2Fgatus--admin-blue?logo=docker)](https://hub.docker.com/r/iamonat/gatus-admin)

A lightweight web admin panel for [Gatus](https://gatus.io). Manage endpoints, alerting providers, branding, announcements, and maintenance windows from the browser — no more SSH + YAML edits.

- **~20 MB RAM, ~60 MB image** — Go binary + Alpine + Docker CLI.
- **Every change writes YAML** to the same directory Gatus reads; one click restarts the Gatus container.
- **Brandable admin panel** — upload your own logo (square + full), pick colors, fonts. Inherits site name / description from Gatus by default.
- **Tooltips everywhere** — every field shows a short explanation pulled from the Gatus docs.
- **Multi-file config** — endpoints, alerting, UI, announcements each live in their own YAML file for clean diffs.

## Quick start

```bash
# 1. Clone and enter the project
git clone https://github.com/onatakduman/gatus-admin && cd gatus-admin

# 2. Copy the example stack files
cp docker-compose.example.yml docker-compose.yml
cp Caddyfile.example Caddyfile
cp .env.example .env

# 3. Generate a password hash and paste it into .env
docker run --rm iamonat/gatus-admin:latest /hashpw mypassword
# or: go run ./cmd/hashpw mypassword

# 4. Edit .env: set ADMIN_PASSWORD_HASH, SECRET_KEY, STATUS_URL.
#    Edit Caddyfile: replace status.example.com with your real domain.

# 5. Start
docker compose up -d

# 6. Open https://status.example.com/admin and sign in.
```

## What's inside

The admin panel exposes these sections, each mapped 1:1 to a Gatus config file:

| Tab | Writes | What it manages |
|---|---|---|
| **Endpoints** | `endpoints.yaml` | HTTP / TCP / ICMP / DNS / SSH / starttls monitors, per-endpoint alerts, conditions, headers, timeouts. Bulk-add several URLs with shared settings in one form. |
| **Alerting** | `alerting.yaml` | ntfy, Slack, Discord, Telegram, Email (SMTP), and custom webhook providers. Per-provider thresholds, send-on-resolved. |
| **UI / Public Page** | `ui.yaml` | Gatus public page branding: title, header, logo URL, dashboard copy, custom CSS, favicons, sort/filter defaults, nav buttons. |
| **Admin Brand** | `brand/admin.yaml` + `brand/assets/` | This admin panel's own branding — separate from the public page. Upload a square logo (sidebar), a full logo (login), favicon, and tweak colors / fonts. |
| **Announcements** | `announcements.yaml` | Banner messages (outage / warning / info / operational). |
| **Maintenance** | `maintenance.yaml` | Scheduled downtime windows. Alerts are suppressed inside them. |
| **Core** | `core.yaml` | Web server address/port, storage backend (memory/sqlite/postgres), basic auth, concurrency, metrics. |
| **Advanced YAML** | every file | Raw editor for each config file, with YAML validation on save. |

After saving any change, click **Apply & Restart Gatus** to make Gatus reload.

## Architecture

```
Browser
   ↓
Caddy  (handles TLS, routes /admin/* to gatus-admin, everything else to Gatus)
   ├─ /admin/* → gatus-admin :8000    — form UI, writes to ./config/*.yaml
   └─ /*       → gatus      :8080    — public status page
                   ↑
                   └── reads ./config/*.yaml
```

Both services mount the same `./config/` volume. gatus-admin restarts the Gatus container via the Docker socket after each save.

## Security

- Password is a bcrypt hash (`cost=12`) supplied via `.env`.
- Session cookie is `HttpOnly + Secure + SameSite=Strict` in production.
- CSRF tokens on every POST.
- Uploaded logos are sanity-checked by content sniffing + extension and capped at 2 MB.
- **The admin container mounts `/var/run/docker.sock`** so it can restart Gatus. This means compromise of the admin → host root. Use a strong password and keep the panel behind a reverse-proxy IP allowlist if it runs on a public server.
- Brand color inputs are validated; `custom-css` fields are explicitly trusted (document, use with care).

## Directory layout

```
./
├── docker-compose.yml        ← you create from example
├── Caddyfile                 ← you create from example
├── .env                      ← you create from example (gitignored)
└── config/                   ← shared between gatus + gatus-admin
    ├── core.yaml
    ├── ui.yaml
    ├── alerting.yaml
    ├── endpoints.yaml
    ├── announcements.yaml
    ├── maintenance.yaml
    └── brand/
        ├── admin.yaml        ← admin panel branding (NEW)
        └── assets/
            ├── logo_full.png
            ├── logo_square.png
            └── favicon.ico
```

## Local development

```bash
# Build and run the dev stack (air hot-reload, port 8000 for admin, 8080 for gatus)
docker compose -f docker-compose.example.yml -f docker-compose.dev.yml up --build

# Admin   → http://localhost:8000/admin
# Gatus   → http://localhost:8080
```

Go files and templates are mounted into the container; `air` rebuilds on save.

## License

MIT — see [LICENSE](LICENSE).

## Acknowledgements

- Built on top of the excellent [Gatus](https://github.com/TwiN/gatus) by TwiN.
- Default dark theme inspired by minimalist dashboards; brand your own colors/fonts via the Admin Brand tab.
