# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Photography/filmmaking portfolio website for Bhavpreet Singh, built with Go's standard library (no external dependencies). Deployed to Heroku.

## Build & Run Commands

```bash
make all              # Clean and build (outputs bin/web)
make clean            # Remove bin/ directory
./bin/web -root .     # Run locally on :8080
```

**Flags:** `-host` (default ""), `-port` (default 8080), `-root` (path to project root, default "."), `-rtimeout` / `-wtimeout` (seconds, default 15)

On Heroku, the `PORT` env var overrides the `-port` flag.

**Tests:** The project uses GOPATH-style builds (`go install web`), not Go modules. Tests are in `src/web/web_test.go`.

## Architecture

Single-file Go server (`src/web/web.go`) using `net/http` with no framework or external dependencies.

**Startup flow (`init()`):**
1. Parse CLI flags for host/port/root/timeouts
2. Compile Go HTML templates from `templates/`
3. Load portfolio data from `data.json` into the global `Info` (MainData struct)
4. Register routes dynamically based on `data.json` categories/projects

**Data model:** `MainData` -> `[]Category` -> `[]Project`. Projects are either `"internal"` (rendered with image carousel via `proj.html`) or `"external"` (link to outside URL).

**Routing:** Routes are registered dynamically from data.json. Static routes: `/`, `/about`, `/contact`. Project routes: `/{category}/{project_link}` (e.g., `/work/wires_and_pigeons`). URL matching uses regex `^/(.*)/(.*)$`.

**Templates:** Go `html/template` with a custom `url_for` function. Templates: `layout.html` (shared head/tail), `index.html`, `proj.html`, `about.html`, `contact.html`. Named template blocks: `layout_head`, `layout_tail`, `page_*`.

**Static assets** served from `/static/` with path stripping. Frontend uses jQuery + Slick carousel for image galleries.
