---
Title: Deployment Diary
Ticket: BOT-SIGNUP-K3S-DEPLOY
Status: active
Topics:
    - deployment
    - kubernetes
    - argocd
    - oauth
    - discord
    - bot-signup
DocType: reference
Intent: ""
Owners: []
RelatedFiles:
    - Path: ../../../../../../../2026-03-27--hetzner-k3s/gitops/kustomize/bot-signup/kustomization.yaml
      Note: Local Kustomize render passed for the new GitOps package
    - Path: ../../../../../../../2026-03-27--hetzner-k3s/gitops/kustomize/bot-signup/runtime-secret.yaml
      Note: VaultStaticSecret for seeded bot-signup OAuth/session runtime secrets
    - Path: ../../../../../../../2026-04-23--pyxis/ttmp/2026/04/29/PYXIS-PRODUCTION-ARGOCD-GLAZED--turn-pyxis-into-a-deployed-glazed-application-on-argocd-k3s/playbooks/02-pyxis-argocd-gitops-playbook.md
      Note: Primary OAuth/Vault/GitOps playbook reference
    - Path: ../../../../../../../2026-04-23--pyxis/ttmp/2026/04/29/PYXIS-PRODUCTION-ARGOCD-GLAZED--turn-pyxis-into-a-deployed-glazed-application-on-argocd-k3s/reference/01-production-deployment-diary.md
      Note: Primary reference diary for the matching Pyxis deployment
    - Path: .github/workflows/publish-image.yaml
      Note: GHCR and GitOps PR workflow added but not yet pushed/run
    - Path: Dockerfile
      Note: Container build passed before server runtime smoke exposed route conflict
    - Path: Makefile
      Note: Read-only inspection of current build targets; file was already modified before deployment work
    - Path: cmd/bot-signup/main.go
      Note: |-
        Read-only inspection of runtime flags before stopping on dirty source tree
        Route registration panic encountered during container HTTP smoke
ExternalSources: []
Summary: Chronological diary for planning the bot-signup k3s deployment and recording the early stop on existing uncommitted source changes.
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---



# Diary

## Goal

Record the deployment planning work for putting `bot-signup` on the Hetzner k3s Argo CD setup as `bot-vibing.yolo.scapegoat.dev`, including references found, decisions made, and the exact reason implementation stopped.

## Step 1: Created the deployment ticket and found the relevant references

I created a dedicated docmgr ticket so the deployment work has a separate planning and diary trail from the product/build ticket. I then located the platform playbook in the k3s repository and the Pyxis production deployment ticket that contains the closest OAuth/Vault/GitOps precedent.

This step established the intended deployment model before touching code: source repo builds an image, GitOps repo owns the Kubernetes desired state, Vault provides runtime OAuth/session secrets, and Argo CD reconciles the app under its own namespace.

### Prompt Context

**User prompt (verbatim):** "Deploy this application to k3s using the bot-vibing.yolo.scapegoat.dev setup, using similar credential setup for oauth than ../2026-04-23--pyxis, in /home/manuel/code/wesen/2026-03-27--hetzner-k3s

Find playbook and pyxis diary related to the deploy, then set up the same setup here. We're still building out the web stuff a little bit, so if you encounter compiling issues or changed code, stop.

Create a new docmgr ticket to make an implementation plan for the deployment, and then write a detailed diary as you work"

**Assistant interpretation:** Create a ticket, find the platform/Pyxis deployment references, plan and set up a similar k3s/Argo/Vault deployment for bot-signup, but stop if the work runs into uncommitted code or compile/build instability.

**Inferred user intent:** Get bot-signup production deployment ready without trampling ongoing web work or attempting risky fixes while the app is still under active construction.

### What I did

- Loaded the `docmgr` and `diary` skills.
- Ran `docmgr status --summary-only` in `/home/manuel/code/wesen/2026-05-01--bot-signup`.
- Listed existing tickets with `docmgr ticket list`.
- Found the k3s repo and the Pyxis repo:
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s`
  - `/home/manuel/code/wesen/2026-04-23--pyxis`
- Found the canonical deployment playbook:
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/source-app-deployment-infrastructure-playbook.md`
- Found the Pyxis deployment references:
  - `/home/manuel/code/wesen/2026-04-23--pyxis/ttmp/2026/04/29/PYXIS-PRODUCTION-ARGOCD-GLAZED--turn-pyxis-into-a-deployed-glazed-application-on-argocd-k3s/design/01-pyxis-production-glazed-argocd-implementation-guide.md`
  - `/home/manuel/code/wesen/2026-04-23--pyxis/ttmp/2026/04/29/PYXIS-PRODUCTION-ARGOCD-GLAZED--turn-pyxis-into-a-deployed-glazed-application-on-argocd-k3s/reference/01-production-deployment-diary.md`
  - `/home/manuel/code/wesen/2026-04-23--pyxis/ttmp/2026/04/29/PYXIS-PRODUCTION-ARGOCD-GLAZED--turn-pyxis-into-a-deployed-glazed-application-on-argocd-k3s/playbooks/02-pyxis-argocd-gitops-playbook.md`
- Created the ticket:
  - `BOT-SIGNUP-K3S-DEPLOY`
- Added docs:
  - `design/01-bot-signup-k3s-deployment-implementation-plan.md`
  - `reference/01-deployment-diary.md`

### Why

- The deployment needs both implementation and operator documentation.
- The Pyxis rollout already solved the platform-specific secrets/GitOps/OAuth shape, so bot-signup should reuse that pattern instead of inventing a new one.
- A separate ticket makes the deployment plan reviewable and resumable after the web changes settle.

### What worked

- The docmgr ticket was created successfully.
- The Pyxis diary/playbook and the k3s source-app deployment playbook were found quickly.
- The current bot-signup runtime contract was identifiable from source: Cobra `serve`, SQLite path, session secret, Discord OAuth env vars, SPA embed handler, and `/api/health`.

### What didn't work

- `docmgr doc search --query deploy` returned:

```text
Error: fts5 not available (docs_fts missing)
```

This did not block the work because direct `find`, `grep`, and known repo paths located the deployment references.

### What I learned

- The platform playbook says a first-time app requires applying `gitops/applications/<app>.yaml` to the cluster once; committing it alone is not sufficient unless an app-of-apps layer exists.
- Pyxis stores runtime OAuth credentials in Vault and syncs them through VSO into a Kubernetes Secret.
- Pyxis uses `pyxis.yolo.scapegoat.dev`; bot-signup should mirror that ingress pattern as `bot-vibing.yolo.scapegoat.dev`.

### What was tricky to build

- The key trick was determining which Pyxis materials were relevant. Pyxis has a long production diary with packaging, PostgreSQL, staff app, Discord OAuth, and smoke-test steps. For bot-signup, the useful subset is the GitOps/Vault/OAuth/Ingress pattern, not the Pyxis-specific Postgres bootstrap and flyer storage machinery.

### What warrants a second pair of eyes

- Confirm whether bot-signup's image should be public or private in GHCR; this changes whether a Vault-backed `image-pull-secret.yaml` is needed.
- Confirm whether SQLite-on-PVC is the desired v1 state backend.

### What should be done in the future

- Resume implementation only after the source tree is clean or the owner confirms the current uncommitted changes are safe to build on.

### Code review instructions

- Start with the plan document:
  - `ttmp/2026/05/01/BOT-SIGNUP-K3S-DEPLOY--deploy-bot-signup-to-k3s-via-bot-vibing-yolo-scapegoat-dev/design/01-bot-signup-k3s-deployment-implementation-plan.md`
- Compare it against:
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/pyxis`
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/applications/pyxis.yaml`

### Technical details

Reference paths:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/source-app-deployment-infrastructure-playbook.md
/home/manuel/code/wesen/2026-04-23--pyxis/ttmp/2026/04/29/PYXIS-PRODUCTION-ARGOCD-GLAZED--turn-pyxis-into-a-deployed-glazed-application-on-argocd-k3s/reference/01-production-deployment-diary.md
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/pyxis
```

## Step 2: Stopped before deployment changes because the source tree is already dirty

Before writing app packaging or GitOps manifests, I checked the current source status. The repository already contains modified and untracked application/web files. Because the user explicitly asked to stop on changed code, I stopped deployment implementation at the planning stage and did not modify application source or k3s GitOps manifests.

This is the safest stopping point: the plan now captures exactly how to resume, but no deployment-specific changes have been layered on top of the in-progress web work.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Do not proceed into packaging/build/deployment if the app is in an unstable or actively changing state.

**Inferred user intent:** Avoid conflicting with ongoing web work and avoid debugging transient build failures caused by partially completed changes.

### What I did

- Ran source status in `/home/manuel/code/wesen/2026-05-01--bot-signup`.
- Observed this exact status:

```text
 M Makefile
 M cmd/bot-signup/main.go
 M go.mod
 M go.sum
 M ui/package.json
?? .github/
?? README.md
?? cmd/build-web/
?? internal/web/
```

- Checked the k3s repo status and observed unrelated untracked local directories:

```text
?? .codex
?? .playwright-mcp/
```

- Read bot-signup runtime files to build the implementation plan without changing them:
  - `cmd/bot-signup/main.go`
  - `internal/server/server.go`
  - `internal/auth/discord_oauth.go`
  - `internal/web/static.go`
  - `Makefile`
- Read Pyxis GitOps files as deployment templates:
  - `gitops/kustomize/pyxis/deployment.yaml`
  - `gitops/kustomize/pyxis/runtime-secret.yaml`
  - `gitops/kustomize/pyxis/ingress.yaml`
  - `gitops/kustomize/pyxis/kustomization.yaml`
  - `gitops/applications/pyxis.yaml`
- Wrote the implementation plan with status `blocked` instead of creating manifests.

### Why

- The user gave a clear stop condition for changed code.
- App packaging and GitOps manifests depend on the final runtime contract. The current app already has ongoing changes to the Makefile, command, web embed path, and UI package metadata.
- Building, fixing, or deploying on top of that state could interfere with the web work.

### What worked

- The current deployment shape could still be planned from read-only inspection.
- The app already exposes enough production configuration to map into Kubernetes env vars:
  - `ADDR`
  - `DB_PATH`
  - `SESSION_SECRET`
  - `DISCORD_CLIENT_ID`
  - `DISCORD_CLIENT_SECRET`
  - `DISCORD_REDIRECT_URL`
  - `SECURE_COOKIES`
- The app has a production-relevant health endpoint at `/api/health`.

### What didn't work

- I did not continue to implement or validate a Dockerfile, GHCR workflow, or GitOps manifests because the dirty tree triggered the requested stop condition.
- I did not run `go test ./...`, `pnpm --dir ui build`, or `make build`; doing so after detecting changed code would have gone beyond the requested guardrail.

### What I learned

- Bot-signup is simpler than Pyxis for v1 deployment because it currently uses SQLite on a local file path instead of a shared PostgreSQL database/bootstrap job.
- The health probe must use `/api/health`, unlike many existing app examples that use `/health`.
- The OAuth credential setup should mirror Pyxis's Vault/VSO pattern but with bot-signup's env var names.

### What was tricky to build

- The main sharp edge is that bot-signup's ongoing web embedding work is exactly the area deployment packaging would touch. The modified files include `Makefile`, `cmd/bot-signup/main.go`, `cmd/build-web/`, and `internal/web/`, which are the same files a production Dockerfile/build flow would depend on. That makes proceeding risky until the current owner finishes or commits the web setup.

### What warrants a second pair of eyes

- Review the uncommitted changes before deciding whether to resume deployment work.
- Review whether SQLite on a PVC is acceptable for the first production deployment and document the one-replica constraint.
- Confirm the Discord app redirect URL is registered exactly as `https://bot-vibing.yolo.scapegoat.dev/auth/discord/callback`.

### What should be done in the future

- Once the source tree is clean, implement the GitOps files listed in the plan.
- Add or verify production image publishing in the app repository.
- Validate with `kubectl kustomize` and server-side dry-run before applying the Argo CD Application.

### Code review instructions

- First review the source tree status and identify ownership of the existing uncommitted files.
- Then review `cmd/bot-signup/main.go` for final env/flag names before writing Kubernetes manifests.
- Then compare any future `gitops/kustomize/bot-signup` package against the existing Pyxis package.

Validation commands to run after the stop condition is cleared:

```bash
go test ./...
pnpm --dir ui lint
pnpm --dir ui build
make build
```

GitOps validation commands to run after manifests exist:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
kubectl kustomize gitops/kustomize/bot-signup
kubectl apply --dry-run=server -k gitops/kustomize/bot-signup
```

### Technical details

Planned Vault runtime path:

```text
kv/apps/bot-signup/prod/runtime
```

Planned runtime keys:

```text
session_secret
sqlite_db_path
website_url
discord_client_id
discord_client_secret
discord_redirect_url
secure_cookies
```

Planned smoke checks:

```bash
curl -i https://bot-vibing.yolo.scapegoat.dev/api/health
curl -i https://bot-vibing.yolo.scapegoat.dev/
```

## Step 3: Added deployment scaffolding, seeded Vault, then stopped on container runtime panic

After the owner fixed direnv and cleared the earlier source-tree stop condition, I resumed deployment setup. I added the local packaging and GitHub Actions pieces in the app repository, created the matching Argo CD/Kustomize desired state in the k3s repo, and seeded Vault runtime/image-pull secrets using the Discord OAuth values now available through direnv.

Build and package validation passed, but a real container server smoke exposed a production-blocking runtime panic in the web route registration. I stopped before applying the Argo CD Application or attempting a live rollout.

### Prompt Context

**User prompt (verbatim):** "fixed direnv"

**Assistant interpretation:** Direnv is now approved, so resume deployment work past the previous environment blocker, while still honoring the instruction to stop on unstable app/web code.

**Inferred user intent:** Continue the k3s deployment setup now that OAuth env vars are available, but do not paper over application build/runtime problems.

### What I did

- Verified direnv now loads Discord env vars.
- Confirmed the application source tree no longer had the earlier app-code modifications; only my docmgr files were untracked/modified at that point.
- Ran validation:
  - `go test ./...`
  - `pnpm --dir ui lint`
  - `pnpm --dir ui build`
- Added app deployment files:
  - `Dockerfile`
  - `.dockerignore`
  - `.github/workflows/publish-image.yaml`
  - `deploy/gitops-targets.json`
  - `scripts/open_gitops_pr.py`
  - Makefile `docker-build` and `docker-smoke` targets
- Ran:
  - `BUILD_WEB_LOCAL=1 make build`
  - `make docker-build IMAGE_REPOSITORY=bot-signup IMAGE_TAG=k3s-local`
  - `make docker-smoke IMAGE_REPOSITORY=bot-signup IMAGE_TAG=k3s-local`
- Created k3s GitOps desired-state files:
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/applications/bot-signup.yaml`
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/bot-signup/*`
- Validated local Kustomize render:
  - `kubectl kustomize gitops/kustomize/bot-signup`
- Seeded Vault:
  - `kv/apps/bot-signup/prod/runtime`
  - `kv/apps/bot-signup/prod/image-pull`
  - policy `bot-signup-prod`
  - Kubernetes auth role `bot-signup-prod`
- Ran a container HTTP server smoke and stopped on a runtime panic.

### Why

- The app needs a production image path before Argo CD can deploy it.
- The k3s repo needs a namespace-scoped Kustomize package and Argo CD Application matching the Pyxis deployment model.
- Vault needs to contain OAuth/session secrets before VSO can materialize the Kubernetes runtime Secret.
- A container-level server smoke is necessary because `--help` and image build success do not prove the embedded web server can start.

### What worked

- Direnv now loads successfully.
- Go tests passed.
- Frontend lint and build passed.
- Embedded binary build passed.
- Docker image build passed.
- Docker `--help` and `serve --help` smoke passed.
- Kustomize rendered the bot-signup package successfully.
- Vault runtime and GHCR pull secrets were seeded without printing secret values.

### What didn't work

Server-side dry-run against the new namespace failed because the namespace does not exist yet and dry-run does not make later namespaced resources see the dry-run namespace:

```text
namespace/bot-signup created (server dry run)
Error from server (NotFound): error when creating "gitops/kustomize/bot-signup": namespaces "bot-signup" not found
```

The blocking failure was the container server smoke. Running `bot-signup serve` in the container exits immediately:

```text
panic: pattern "GET /{filepath...}" (registered at github.com/go-go-golems/bot-signup/cmd/bot-signup/main.go:79) conflicts with pattern "GET /" (registered at github.com/go-go-golems/bot-signup/cmd/bot-signup/main.go:78):
	GET /{filepath...} matches the same requests as GET /

goroutine 1 [running]:
net/http.(*ServeMux).register(...)
	net/http/server.go:2911
net/http.(*ServeMux).Handle(0xc000091b28?, {0x9f776d?, 0xc0000e2270?}, {0xb53560?, 0xc000096500?})
	net/http/server.go:2874 +0x4d
main.runServe({0xc000012015, 0x5})
	github.com/go-go-golems/bot-signup/cmd/bot-signup/main.go:79 +0x34f
```

### What I learned

- The app can compile and package, but the embedded SPA route registration is not currently safe in Go's `http.ServeMux` because `GET /` and `GET /{filepath...}` overlap.
- This is exactly the deployment-critical web edge the user warned about, so deployment should stop here.
- The local repository still has no GitHub remote, so even after the runtime panic is fixed, source publishing needs either a new GitHub repository or an added remote before GHCR CI can publish.

### What was tricky to build

- The Docker `--help` smoke was not enough. The binary starts Cobra successfully, but the HTTP server path only initializes the conflicting ServeMux patterns when `serve` actually runs and embedded assets are present.
- Kustomize local render and server dry-run have different failure modes for first-time namespaces. Local render proves YAML structure; server dry-run of all resources fails until the namespace exists or resources are applied in waves by Argo.

### What warrants a second pair of eyes

- Review `cmd/bot-signup/main.go` route registration before changing it. The likely fix is to register only one SPA catch-all pattern, but the final shape should preserve `/api` and `/auth` exclusions from `internal/web/static.go`.
- Review whether the image repository should be `ghcr.io/go-go-golems/bot-signup` or another owner before pushing/rolling out.
- Review whether `image: ghcr.io/go-go-golems/bot-signup:main` is acceptable as a bootstrap value or should be replaced by an immutable `sha-<short-sha>` before merge.

### What should be done in the future

- Fix the ServeMux route conflict in application code.
- Re-run container HTTP smoke against `/api/health`.
- Add or create the GitHub remote, configure `GITOPS_PR_TOKEN`, and push the source repository.
- Let CI publish the image and update GitOps to the immutable short SHA tag.
- Only then apply `gitops/applications/bot-signup.yaml` and watch Argo CD.

### Code review instructions

Start review in these files:

- `cmd/bot-signup/main.go` for the runtime panic.
- `Dockerfile` for the production image build.
- `.github/workflows/publish-image.yaml` for GHCR and GitOps PR automation.
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/bot-signup/deployment.yaml` for runtime env and health probes.
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/bot-signup/runtime-secret.yaml` for the Vault path.

Validation commands after fixing the runtime panic:

```bash
go test ./...
pnpm --dir ui lint
pnpm --dir ui build
BUILD_WEB_LOCAL=1 make build
make docker-smoke IMAGE_REPOSITORY=bot-signup IMAGE_TAG=k3s-local
```

Then run an actual container server smoke:

```bash
docker run --rm -P \
  -e ADDR=:8080 \
  -e DB_PATH=/data/test.db \
  -e SESSION_SECRET=local-smoke-secret \
  -e DISCORD_CLIENT_ID=dummy \
  -e DISCORD_CLIENT_SECRET=dummy \
  -e DISCORD_REDIRECT_URL=http://localhost/auth/discord/callback \
  bot-signup:k3s-local serve
```

### Technical details

Vault runtime path:

```text
kv/apps/bot-signup/prod/runtime
```

Vault image-pull path:

```text
kv/apps/bot-signup/prod/image-pull
```

Kubernetes auth role:

```text
bot-signup-prod
```

Planned production host:

```text
https://bot-vibing.yolo.scapegoat.dev
```

## Step 4: Promoted first production Discord user and documented SQLite access

After the first Discord login created a production SQLite user row, I used a short-lived Kubernetes Job to mount the `bot-signup-data` PVC and run `sqlite3` against `/data/bot-signup.db`. The Job guarded against accidentally promoting multiple users by checking that exactly one user existed before running the broad `update users set role = 'admin'` statement.

I also wrote a repeatable production SQLite access playbook so future operator tasks do not require ad hoc shell reconstruction.

### Prompt Context

**User prompt (verbatim):** "I added a user, can you promote it to admin, and also make a user guide on how to access the prod db?

 After you sign in once with Discord, your user should exist in SQLite.                              
                                                                                                     
 If there is only one local user, run:                                                               
                                                                                                     
 ```bash                                                                                             
   sqlite3 data/bot-signup.db \                                                                      
     \"update users set role = 'admin', updated_at = datetime('now');\"                                
 ```                                                                                                 
                                                                                                     
 Then refresh:                                                                                       
                                                                                                     
 ```text                                                                                             
   http://localhost:5179/admin                                                                       
 ```                                                                                                 
                                                                                                     
 Safer: promote only your Discord account                                                            
                                                                                                     
 Find users:                                                                                         
                                                                                                     
 ```bash                                                                                             
   sqlite3 data/bot-signup.db \                                                                      
     \"select id, discord_id, email, display_name, role, status from users;\"                          
 ```                                                                                                 
                                                                                                     
 Then promote your user:                                                                             
                                                                                                     
 ```bash                                                                                             
   sqlite3 data/bot-signup.db \                                                                      
     \"update users set role = 'admin', updated_at = datetime('now') where discord_id =               
 'YOUR_DISCORD_ID';\"                                                                                 
 ```                                                                                                 
                                                                                                     
 Then go to:"

**Assistant interpretation:** Promote the newly-created production user to admin in the k3s SQLite database and write a reusable operator guide for production DB access.

**Inferred user intent:** Bootstrap the first admin account and leave behind safe instructions for future admin/user maintenance.

### What I did

- Verified the production pod and `/data/bot-signup.db` exist.
- Confirmed the app image does not include `sqlite3`.
- Created and ran a short-lived Job named `bot-signup-sqlite-admin-oneshot` in namespace `bot-signup` using `alpine:3.20` plus `apk add sqlite`.
- Mounted the `bot-signup-data` PVC into the Job at `/data`.
- Queried users before promotion.
- Verified exactly one user existed.
- Updated that user to `role='admin'`.
- Queried users after promotion.
- Added playbook `playbook/01-production-sqlite-access-guide.md`.

### Why

- The production app stores SQLite on a PVC, not inside the container filesystem.
- The runtime image is intentionally small and does not ship the `sqlite3` CLI.
- A Kubernetes Job gives a reproducible, auditable way to access the DB without modifying the running app image.

### What worked

The guarded promotion Job completed successfully.

Before:

```text
1|363877777977376768|wesen@ruinwesen.com|slono|user|waiting
```

After:

```text
1|363877777977376768|wesen@ruinwesen.com|slono|admin|waiting
```

### What didn't work

An initial thought was to run `sqlite3` directly in the app container, but the binary is not installed there. That is expected for the slim production image.

### What I learned

- The production DB file is present at `/data/bot-signup.db` and accompanied by SQLite WAL files.
- For one-off access, mounting the PVC into a temporary Job is cleaner than copying the SQLite database out and back in, especially with WAL files.

### What was tricky to build

- SQLite plus WAL files should not be casually copied while the app is running. The temporary Job approach avoids having to manually synchronize `bot-signup.db`, `bot-signup.db-wal`, and `bot-signup.db-shm`.
- Because the PVC is `ReadWriteOnce`, this works in the current single-node cluster. In a future multi-node cluster, maintenance Jobs may need node affinity or an app scale-down window.

### What warrants a second pair of eyes

- Confirm whether first-admin bootstrap should also set `status='approved'`. The requested operation only changed `role` to `admin`, and the resulting row still has `status='waiting'`.
- Consider adding a first-admin bootstrap command or admin CLI to avoid direct SQLite writes.

### What should be done in the future

- Delete the completed Job after TTL cleanup if it lingers.
- Add an app-level admin bootstrap path if this becomes a repeated operation.

### Code review instructions

- Review `playbook/01-production-sqlite-access-guide.md` for operational accuracy.
- Review admin middleware behavior to determine whether `role=admin` alone is enough or whether `status` should also change.

### Technical details

The successful command path was:

```bash
kubectl apply -f /tmp/bot-signup-sqlite-job.yaml
kubectl -n bot-signup wait --for=condition=complete job/bot-signup-sqlite-admin-oneshot --timeout=120s
kubectl -n bot-signup logs job/bot-signup-sqlite-admin-oneshot
```
