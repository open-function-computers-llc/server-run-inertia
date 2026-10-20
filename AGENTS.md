# AGENTS.md — server-run-inertia

## What this is
Single-binary Go web app managing website-hosting **accounts** on a server.
`urfave/cli` entrypoint; the production subcommand is `serve` (alias `s`), which
boots an HTTP server that (1) renders a Vue 3 **Inertia** app (SSR-less: Go
injects a `data-page` JSON blob into `dist/index.html`), (2) exposes JSON APIs,
and (3) runs whitelisted shell scripts behind WebSocket streaming.

> **Do not confuse with `../server-run`** — that is the OLD, abandoned Angular
> version. This repo is the active rewrite. (The MySQL admin feature was built
> there first, then ported here; the old repo is reference only.)

## Build / run / verify
- `./build.sh` — canonical loop: `killall server-run` → `npm exec mix --production`
  (SKIPPED when `APP_ENV=development`) → `CGO_ENABLED=0 go build
  -ldflags "-X main.Version=$(git rev-parse HEAD)" -o cmd/server-run` → runs it.
- **Frontend is embedded in the binary**: `//go:embed dist/*` (`main.go:18`).
  Build the frontend (`npx mix --production`) BEFORE `go build` or Vue changes
  are invisible. Mix emits `dist/app.js|app.css|index.html|mix-manifest.json`.
- `./dev.sh` — live-reload dev mode (port 8983). `APP_ENV=development`
  **disables session validation** (see Auth) — never in prod.
- Binary is `cmd/server-run`; `.gitignore` excludes `cmd/*` and `dist/*` —
  unlike the old repo the binary is NOT committed.
- No CI, no linter config. "Tests" = `go build` + `go vet` + `mix --production`
  compile clean. One Go test exists: `account/account_test.go`.
- Toolchain trap: bare `go get`/`go build` may try to download toolchain
  `go1.25` and die with "toolchain not available". Use `GOTOOLCHAIN=go1.26.9`.

## Env vars
`.env` (godotenv/autoload) + `/etc/server-run.env` if present. `main.go:verifyValidENV()`
hard-fails at boot if ANY of these are empty: `APP_PORT, APP_URL, AUTH_USER,
AUTH_PASSWORD, UPTIME_USER, UPTIME_PASS, UPTIME_DOMAIN,
THIRD_PARTY_ACCESS_TOKEN, WEBSITES_ROOT, ACCOUNTS_ROOT, APP_SOCKET_PATH`.
Also read: `APP_ENV` (`development` = debug logs + auth bypass + build.sh skips
mix), `APP_LIVERELOAD_PORT`, `SCRIPTS_ROOT` (not in the required list, but
`scriptrunner.VerifyAllScriptsExist()` panics at boot if any registered script
file is missing — see `script-runner/{argument,plain,env}Scripts.go`; 16 files
today, including the extensionless `f2bstatus`), `VCS_PROVIDER`
(+`GITEA_TOKEN`/`GITEA_URL`).
**Trailing slashes on `WEBSITES_ROOT`/`ACCOUNTS_ROOT`/`SCRIPTS_ROOT` are
load-bearing** — paths are built with `+` (e.g. `ACCOUNTS_ROOT + name`).

## Top-level layout
```
main.go                 — cli app; embeds dist/* + assets/logo; env validation
app/                    — cli commands (serve, listSetting, update, envExport, repo cmds)
server/                 — HTTP layer: server struct, routes, ALL handlers, middleware
  routes.go             — the ONLY place routes are bound (bindRoutes)
  server.go             — server struct + New(): session bag, inertia manager, VCS
  handlePage.go         — generic Inertia page renderer
  findAccount.go        — lookup against the boot snapshot of accounts
account/                — domain model: Account, Database, settings.json IO
session/                — in-memory session store (cookie key -> expiry)
script-runner/          — whitelisted shell-script registries + streaming runner
errx/                   — error wrapping helper (errx.Wrapf)
vcs/                    — VCS providers (dummy default, gitea, github)
frontend/               — Vue 3 Inertia app (Laravel Mix)
  app.js, Pages/, Components/, Components/Account/, Layouts/, Icons/, scss/, composables/
server/index.html       — Inertia root template (mix copies to dist/index.html)
cmd/server-run          — built binary (gitignored)
dist/                   — mix build output (gitignored, embedded)
```

## Auth (prod depends on this — do not "fix" it)
- `GET /login` renders Inertia page `Login` (`frontend/Pages/Login.vue`).
- `POST /handle-auth` (`server/handleFormProcessAuth.go`) expects a **JSON**
  body `{"userName":..., "password":...}` — NOT form-encoded — compares against
  `AUTH_USER`/`AUTH_PASSWORD`, then `session.Create()` (uniuri 64) sets cookie
  **`server-run-auth`** (30 min) and 303-redirects to `/dashboard`.
- `ProtectRequest` (`middleware_protected.go`): missing/invalid cookie → 303 →
  `/login`. Browser routes never see a 401.
- **`session.Validate` returns true for ANY cookie value when
  `APP_ENV=development`** (`session/session.go:42`). For scripted API tests in
  dev mode any `-b "server-run-auth=x"` works.
- **Known bug (verified):** login failures report via
  `s.inertiaManager.Share("error", ...)`, which is **process-global and never
  cleared** — after one bad login, EVERY render (including a fresh `GET /login`)
  carries the stale `error` prop until the process restarts. Never mistake a
  stale "Invalid Credentials" for a live auth failure; restart to clear, and
  check the server log for whether a `POST /handle-auth` actually arrived.
- `session.Initialize()`'s GC goroutine sweeps once and exits (no outer loop),
  and the map has no mutex. Don't add goroutines touching `sBag.sessions`.
- Credentials for MySQL live in the same plaintext `settings.json` as everything
  else (pre-existing pattern) — keep `$ACCOUNTS_ROOT` permissions tight.

## Routes (`server/routes.go` — single source of truth)
Go 1.22+ patterned `http.ServeMux` (`"GET /account/{account}"`, params via
`r.PathValue("account")`). NOT gorilla/mux. Four groups:
- **open**: `/login`, `POST /handle-auth`, `/error`, `/logo`, `/`.
- **protected** (pages + `/api/...`): `inertiaManager.Middleware(LogRequest(ProtectRequest(h)))`.
- **external** (`/external-access/*`): `thirdPartyProtection` checks static
  header `X-Server-Run-Access: $THIRD_PARTY_ACCESS_TOKEN`.
- **websocket** (`/stream/system-load`, `/stream/script-runner`): ProtectRequest;
  browsers can't set headers on a WS handshake, so the token rides in the query string.
Also: unix-socket server on `APP_SOCKET_PATH` (`server/socket.go`) and a
`GET /dist/*` file server for mix assets.

**Adding a JSON API endpoint**: write `handleFoo()` in `server/`, add one line
to `protectedRoutes`. Use `sendJSON(w, data)` / `sendJSONError(w, status,
{"error": "..."})` (`server/sendJSON.go`). `sendJSON` takes no status arg
(always 200, even on marshal failure).

## The account model (`account/`)
- One directory = one account: `$ACCOUNTS_ROOT/<name>/settings.json` is the only
  persistent state; the `json` tags ARE the file format.
- `account.Find(name)` hydrates from settings.json (creating dir/file with
  defaults if absent — that IS how accounts get created). `bootstrapAccounts()`
  shells out to `ls $ACCOUNTS_ROOT` at boot into `s.accounts`, a **snapshot**;
  `handleAccountDetails` re-runs it per request. `findAccount()` reads the
  snapshot only — re-`Find` when freshness matters.
- `Account`: `Name`, `PrimaryDomain` (json `domain`), `AlternateDomains`,
  `Databases []Database{Username,Password,Host,Name}`, `IsLocked`,
  `IsLockedDownDisabled`, `UptimeURL`, `AlwaysUnlockedDirectories`, `PubKey`
  (json `sshPubKey`), `CreatedAt`, `UpdatedAt`.
- Every mutator = set field + `writeSettingsFile()` (Lock, Unlock,
  SetPrimaryDomain, AddDatabase, ...). Copy one; receivers are pointers, `Find`
  returns `*Account` but `s.accounts` stores values.

## Frontend (Vue 3 + Inertia, `frontend/`)
- **Laravel Mix** (`webpack.mix.js`): entry `frontend/app.js`, `@` → `frontend/`,
  sass entry `frontend/scss/app.scss`, copies `server/index.html` →
  `dist/index.html`, `.version()` cache-busting via `mix-manifest.json`.
- Pages in `frontend/Pages/`; shared bits in `frontend/Components/`;
  account-scoped tabs in `frontend/Components/Account/`.
- **Account details page**: `Pages/Account/Show.vue` holds a `tabComponents` map
  (tab name → component); `Components/Account/SubNav.vue` renders the bar and
  emits `update-tab`. A new tab = component + SubNav entry + Show.vue entry.
- SCSS: components `@import "../../scss/variables.scss"` — path is
  depth-sensitive (Components/Account is two levels below `frontend/`).
- API errors: check `err.response.data.error` (server shape is `{"error":...}`);
  `Databases.vue`'s `apiError()` helper is the pattern.

## MySQL database admin (branch `adminer-in-go`)
Backend `server/handleDatabases.go`; routes under `protectedRoutes` in
`routes.go`; UI `frontend/Components/Account/Databases.vue`.
- **Security boundary** `openDBForAccount(account, db)`: connects ONLY if `db`
  is listed in that account's OWN `settings.json` (`Account.Databases`), and the
  account must be registered with this server. It re-reads settings.json via
  `account.Find` per request so newly added databases work without a restart.
  **Never widen this.** The picker lists only the viewed account's own databases.
- Endpoints: `GET .../databases/{db}/tables`, `GET .../tables/{table}?limit=`
  (limit clamped 1–5000; server builds the SELECT itself — no arbitrary SQL),
  `POST .../tables/{table}/rows` (insert), `PUT` (update by full PK), `DELETE`
  (delete by full PK; tables without a PK are refused).
- **UPDATE rewrites ONLY the columns you changed.** The server builds its SET
  clause purely from the `col_*` fields that arrive, so a column the client
  omits is never written; `Databases.vue` therefore sends only fields whose
  value differs from the `original` baseline captured in `startEdit` (INSERT
  still sends everything, and an edit that changes nothing sends nothing and is
  refused with "Nothing to save — no fields were changed."). This is deliberate:
  a blind full-row rewrite silently reverts concurrent edits by others
  (lost-update), re-applies lossy round-trips of binary/long-text values, and
  needlessly churns `ON UPDATE CURRENT_TIMESTAMP` and generated columns.
  `validateFields()` matches this and skips untouched columns, so a stored value
  the UI cannot represent (legacy `0000-00-00`, a blob) cannot block an unrelated
  edit. Verified by mutating a column directly in MySQL *while the edit modal was
  open* and confirming the save left it intact.
- **Safety by construction**: identifiers validated against
  `information_schema.COLUMNS` then backtick-quoted; user input arrives ONLY as
  `?`-bound parameters. Form convention: `col_<name>=value` for SET/VALUES,
  `pk_<name>=value` for WHERE; unknown columns silently ignored; empty value on
  a NULLABLE column stores NULL.
- Gotchas baked in — do not regress:
  - Driver returns `[]byte` for BLOB/some TEXT → `normalizeRowValues` converts
    to string, else the UI shows base64 garbage.
  - **`parseTime` is deliberately OFF in the DSN.** With it on, the driver turns
    DATE/DATETIME/TIME into `time.Time`, which JSON-marshals as RFC3339 with a
    bogus `Z` (`"2024-03-15T00:00:00Z"`); the client then does `new Date()`
    timezone math on a value that has no timezone and shifts it on save, and
    TIME values outside 00:00:00–23:59:59 can't round-trip at all. Left off, the
    driver hands back MySQL's own canonical text (`"08:30:00"`,
    `"2024-03-15 00:00:00"`) which `normalizeRowValues` makes plain strings, so
    the UI writes back exactly what it read. `toInputValue` treats these as
    strings — never reintroduce `new Date()` there.
  - `fetchColumnInfo` also selects **`COLUMN_TYPE`** — the client needs it for
    varchar length, decimal scale and enum members (unknowable from DATA_TYPE).
  - **EVERY form field renders as `type="text"`** (plus `<textarea>` for
    text/blob/json), validated by `validateValue()` against `COLUMN_TYPE`.
    Native typed inputs — `type="number"` AND `type="date"`/`"time"`/
    `"datetime-local"` — are banned here: on an unparseable entry the browser
    reports `""`, the server maps `""` to NULL for nullable columns, and the
    typo silently wipes the stored value with no error (verified live: `"abc"`
    into an INT and `"25:99"` into a TIME both saved NULL). Fixing only the
    numeric case leaves date/time wide open — it is the same flaw. Expected
    format lives in the `placeholder` instead (enum columns list their allowed
    values there). Empty means NULL only where the column allows it; NOT NULL
    rejects empty rather than letting MySQL's non-strict mode insert `0`/`''`.
    The server remains the authority; client validation only stops obvious
    mistakes.
  - Validation constants that are easy to get wrong: unsigned max is
    `2*signedMax+1` (tinyint unsigned = 255, **not** the signed 127 — reusing
    the signed bound blocks legitimate writes); MySQL TIME spans
    −838:59:59..838:59:59 so hours may exceed 24 but minutes/seconds must still
    be 0–59 (a digit-count-only regex let `"25:99"` through); use `BigInt` for
    range checks so bigint bounds don't lose precision.
  - MariaDB backs JSON with `longtext` + a CHECK constraint, so
    `information_schema` reports `longtext` and the client cannot know a column
    is JSON: bad JSON is rejected **server-side** (`Error 4025 ... CONSTRAINT
    ... failed`) with the typed value intact. That is the correct outcome —
    fail loudly, never blank a value.
  - Add-mode: only **auto-increment** columns are readonly/skipped. A plain
    (non-auto-inc) PK IS user-entered — skipping all PKs breaks tables like
    `nullable_test(id NOT NULL, no AI)` with MySQL error 1364. Edit-mode locks
    the PK (row identity); tables with no PK are view-only by design.
  - Rows expand on **click**, not hover (hover fired on incidental mouse
    movement and fought the horizontal scrollbar).
  - The db picker must load tables on **mount** too: a pre-filled
    `<select v-model>` never fires `@change`, so a `watch(x, fn, {immediate:true})`
    is required or the tab shows "No tables" forever while the API works fine.
- Confirmations are a self-dismissing toast naming table + row identity. Because
  it self-dismisses, assert it in the SAME browser call as the triggering action.
- v1 scope: single-row insert/update/delete, no DDL. Extend = new handler +
  route + form block in Databases.vue; keep everything parameterized.

## Verifying UI changes
curl proving an API returns data does NOT prove the component calls it. Drive a
real browser: fill `#username`/`#password`, click Sign in, goto
`/account/<name>`, click the tab by text, then read `.table-item .tname`,
`.table-scroll thead th`, `tbody tr` counts. When writing validation tests, pick
fixture columns that actually match the case (checking a NOT NULL rejection on a
nullable column passes for the wrong reason), and clean up junk rows through the
feature under test.
- `validateValue`/`parseColumnType` are pure — test them as a truth table in
  node (extract the block from the SFC, `[label, column, value, 'ok'|'err']`)
  rather than clicking the modal dozens of times. Cover BOTH directions: false
  positives that block legitimate data (255 into tinyint unsigned) only surface
  in the `'ok'` rows. Then spot-check the shipped bundle in a browser.
- "No error shown for field X" is NOT proof X passed — a different invalid field
  can block the save while you read only X's error. Prove a verdict with
  modal-still-open + ALL `.field-error`s + the server log (a real write logs
  `db-admin: UPDATE ...`) + a direct `mysql` SELECT. Saving with a valid value
  performs a REAL write, so restore any fixture row you clobber.
- Save blocked with zero visible field errors usually means a server-side error
  landed in the form-level message — dump `.row-form` innerText before assuming
  the client logic is wrong.

## Conventions / gotchas (general)
- Handlers are `func (s *server) handleX() http.HandlerFunc` closures (`server`
  is lowercase/unexported, package `server`).
- Older endpoints use **form-encoded** bodies (`r.ParseForm`, `PostFormValue`);
  login sends JSON and db-admin writes use multipart. Match the client.
- Error wrapping: `errx.Wrapf("context", err, sentinelError)` — see
  `account/errors.go`, `errx/wrap.go`.
- `fmt.Println` debugging is common in committed code (e.g. `account.Find`) —
  noise, not yours to fix piecemeal.
- Inertia `Share` props persist process-wide (see the login bug); prefer
  `WithProp(r.Context(), ...)` for per-request props.
- `Components/Account/previousAccountPageFromAngular.component.html` is leftover
  reference material from the port — unused.
- Dev loop / test rig details live in the `server-run-inertia-dev` skill, not
  here (scratch dirs get pruned after 24h idle). One trap belongs here: never
  recreate a test env file by transcribing secrets from chat history —
  transcript redaction turns a password into literal `***`, the server boots
  with garbage, and every real login is correctly rejected while your own curl
  (sending the same `***`) appears to pass.
