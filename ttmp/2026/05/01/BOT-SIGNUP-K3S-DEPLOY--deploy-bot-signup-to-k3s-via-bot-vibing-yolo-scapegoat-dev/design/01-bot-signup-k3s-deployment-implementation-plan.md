---
Title: Bot Signup k3s Deployment Implementation Plan
Ticket: BOT-SIGNUP-K3S-DEPLOY
Status: blocked
Topics:
    - deployment
    - kubernetes
    - argocd
    - oauth
    - discord
    - bot-signup
DocType: design
Intent: ""
Owners: []
RelatedFiles:
    - Path: ../../../../../../../2026-03-27--hetzner-k3s/docs/source-app-deployment-infrastructure-playbook.md
      Note: Canonical platform playbook for GHCR to GitOps to Argo CD deployments
    - Path: ../../../../../../../2026-03-27--hetzner-k3s/gitops/applications/bot-signup.yaml
      Note: Argo CD Application for first-time bootstrap
    - Path: ../../../../../../../2026-03-27--hetzner-k3s/gitops/kustomize/bot-signup/deployment.yaml
      Note: Bot-signup Kubernetes Deployment with OAuth env and health probes
    - Path: ../../../../../../../2026-03-27--hetzner-k3s/gitops/kustomize/pyxis/deployment.yaml
      Note: Reference app Deployment with Vault-synced OAuth env vars
    - Path: ../../../../../../../2026-04-23--pyxis/ttmp/2026/04/29/PYXIS-PRODUCTION-ARGOCD-GLAZED--turn-pyxis-into-a-deployed-glazed-application-on-argocd-k3s/reference/01-production-deployment-diary.md
    - Path: .github/workflows/publish-image.yaml
      Note: GHCR publish and GitOps PR automation workflow
    - Path: Dockerfile
      Note: Production multi-stage image build for embedded UI and Go server
    - Path: Makefile
      Note: Build and Docker smoke targets used during validation
    - Path: cmd/bot-signup/main.go
      Note: |-
        Defines the bot-signup serve command
        Contains the current ServeMux route conflict blocking live rollout
    - Path: deploy/gitops-targets.json
      Note: Source-repo target metadata for bot-signup GitOps image bumps
    - Path: internal/auth/discord_oauth.go
      Note: Defines Discord OAuth scopes and callback exchange behavior
    - Path: internal/server/server.go
      Note: Defines production HTTP routes including /api/health and OAuth endpoints
    - Path: internal/web/static.go
      Note: Defines embedded SPA fallback behavior relevant to the production image
ExternalSources: []
Summary: Plan for deploying bot-signup to the Hetzner k3s Argo CD setup at bot-vibing.yolo.scapegoat.dev, intentionally paused before implementation because the source repository already had uncommitted code changes.
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---



# Bot Signup k3s Deployment Implementation Plan

## Goal

Deploy `/home/manuel/code/wesen/2026-05-01--bot-signup` to the Hetzner k3s GitOps repository at `/home/manuel/code/wesen/2026-03-27--hetzner-k3s`, reachable as:

```text
https://bot-vibing.yolo.scapegoat.dev
```

The requested model is the same operator path used by Pyxis: app repository builds and publishes an immutable GHCR image, the k3s GitOps repository pins that image in `gitops/kustomize/<app>`, Argo CD reconciles it, and OAuth/runtime secrets are sourced from Vault through Vault Secrets Operator.

## Current stop condition

Implementation is paused before live rollout. The first stop condition, an already-dirty application tree, was cleared by the owner. I then added local deployment infrastructure and GitOps desired-state manifests, but stopped before applying Argo CD because the production container exits immediately during HTTP route registration.

The user explicitly requested: “We're still building out the web stuff a little bit, so if you encounter compiling issues or changed code, stop.” This is not a compile failure, but it is a web/server runtime failure in the same deployment-critical area, so the safe behavior is to stop rather than patching application code during deployment work.

Container smoke failure:

```text
panic: pattern "GET /{filepath...}" (registered at github.com/go-go-golems/bot-signup/cmd/bot-signup/main.go:79) conflicts with pattern "GET /" (registered at github.com/go-go-golems/bot-signup/cmd/bot-signup/main.go:78):
	GET /{filepath...} matches the same requests as GET /
```

Earlier observed source status before the owner fixed/committed the web work:

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

## Reference material found

### Platform playbook

Primary platform playbook:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/source-app-deployment-infrastructure-playbook.md
```

Relevant rules from that playbook:

- The app repository owns source, tests, Dockerfile, image publishing workflow, and deployment target metadata.
- The GitOps repository owns Kubernetes desired state.
- First deployment of a new app needs a one-time `kubectl apply -f gitops/applications/<app>.yaml`; merely committing the file is not enough unless an app-of-apps layer is added later.
- If GHCR publishing fails in CI, stop and ask for guidance instead of doing an unplanned local workaround.
- If runtime secrets or identity are needed, use the app runtime secrets / Vault path.

### Pyxis deployment diary and playbook

Primary Pyxis references:

```text
/home/manuel/code/wesen/2026-04-23--pyxis/ttmp/2026/04/29/PYXIS-PRODUCTION-ARGOCD-GLAZED--turn-pyxis-into-a-deployed-glazed-application-on-argocd-k3s/reference/01-production-deployment-diary.md
/home/manuel/code/wesen/2026-04-23--pyxis/ttmp/2026/04/29/PYXIS-PRODUCTION-ARGOCD-GLAZED--turn-pyxis-into-a-deployed-glazed-application-on-argocd-k3s/playbooks/02-pyxis-argocd-gitops-playbook.md
```

The matching deployed GitOps reference is:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/pyxis
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/applications/pyxis.yaml
```

Pyxis uses these important patterns:

- Namespace-scoped app package under `gitops/kustomize/pyxis`.
- `VaultConnection`, `VaultAuth`, and `VaultStaticSecret` resources for runtime and image pull secrets.
- Runtime secret path `kv/apps/pyxis/prod/runtime`.
- Ingress host `pyxis.yolo.scapegoat.dev` with Traefik and `letsencrypt-prod`.
- Sync waves: namespace/service accounts/Vault auth first, `VaultStaticSecret` next, workload and service next, ingress last.
- Private GHCR image pull secret via Vault when needed.

## Application runtime contract observed

The current `bot-signup` server is a Go/Cobra binary with a `serve` subcommand.

Important runtime flags and environment variables from `cmd/bot-signup/main.go`:

| Flag | Env var | Current default | Production value |
| --- | --- | --- | --- |
| `--addr` | `ADDR` | `:8080` | `:8080` or `0.0.0.0:8080` |
| `--db` | `DB_PATH` | `data/bot-signup.db` | `/data/bot-signup.db` |
| `--session-secret` | `SESSION_SECRET` | `dev-insecure-change-me` | Vault secret |
| `--discord-client-id` | `DISCORD_CLIENT_ID` | empty | Vault secret |
| `--discord-client-secret` | `DISCORD_CLIENT_SECRET` | empty | Vault secret |
| `--discord-redirect-url` | `DISCORD_REDIRECT_URL` | `http://localhost:8080/auth/discord/callback` | `https://bot-vibing.yolo.scapegoat.dev/auth/discord/callback` |
| `--secure-cookies` | `SECURE_COOKIES` | false unless exactly `true` | `true` |

HTTP surface from `internal/server/server.go`:

- Health: `GET /api/health`
- OAuth login/callback: `GET /auth/discord/login`, `GET /auth/discord/callback`
- Auth/session APIs: `/api/auth/*`
- User profile APIs: `/api/profile`
- Admin APIs: `/api/admin/*`
- SPA fallback is currently added by `cmd/bot-signup/main.go` when embedded web assets are available.

The current health endpoint is `/api/health`, not `/health`. Kubernetes probes and smoke tests should use `/api/health` unless the app adds a top-level `/health` route.

## Desired deployment architecture

```text
bot-signup app repo
  -> tests/builds embedded Go + Vite artifact
  -> publishes ghcr.io/wesen/bot-signup:sha-<commit>
  -> optional GitOps target metadata opens PR
      -> ../2026-03-27--hetzner-k3s/gitops/kustomize/bot-signup
      -> Argo CD Application bot-signup
      -> namespace bot-signup
      -> https://bot-vibing.yolo.scapegoat.dev
```

Runtime state:

- One Deployment, one replica initially.
- One Service on port 80 to container port 8080.
- One Ingress with TLS for `bot-vibing.yolo.scapegoat.dev`.
- One PVC mounted at `/data` for SQLite persistence.
- One Vault-synced runtime Kubernetes Secret named `bot-signup-runtime`.
- Optional Vault-synced GHCR pull secret if the image is private.

## Proposed Vault credential setup

Use the same shape as Pyxis, but smaller because this app currently uses SQLite and Discord OAuth only.

Vault path:

```text
kv/apps/bot-signup/prod/runtime
```

Keys:

```text
session_secret
sqlite_db_path
website_url
discord_client_id
discord_client_secret
discord_redirect_url
secure_cookies
```

Recommended values:

```text
sqlite_db_path=/data/bot-signup.db
website_url=https://bot-vibing.yolo.scapegoat.dev
discord_redirect_url=https://bot-vibing.yolo.scapegoat.dev/auth/discord/callback
secure_cookies=true
```

If the GHCR package is private, add the same private image pull pattern used by Pyxis:

```text
kv/apps/bot-signup/prod/ghcr-pull
```

with a Kubernetes dockerconfigjson payload shape matching the existing Pyxis `image-pull-secret.yaml` convention.

## GitOps files created but not applied

In `/home/manuel/code/wesen/2026-03-27--hetzner-k3s`:

```text
gitops/applications/bot-signup.yaml
gitops/kustomize/bot-signup/namespace.yaml
gitops/kustomize/bot-signup/serviceaccount.yaml
gitops/kustomize/bot-signup/vault-connection.yaml
gitops/kustomize/bot-signup/vault-auth.yaml
gitops/kustomize/bot-signup/runtime-secret.yaml
gitops/kustomize/bot-signup/image-pull-secret.yaml     # only if private image
gitops/kustomize/bot-signup/persistentvolumeclaim.yaml
gitops/kustomize/bot-signup/deployment.yaml
gitops/kustomize/bot-signup/service.yaml
gitops/kustomize/bot-signup/ingress.yaml
gitops/kustomize/bot-signup/kustomization.yaml
```

Deployment environment should map the Vault secret into the app's existing env contract:

```yaml
env:
  - name: ADDR
    value: :8080
  - name: DB_PATH
    valueFrom:
      secretKeyRef:
        name: bot-signup-runtime
        key: sqlite_db_path
  - name: SESSION_SECRET
    valueFrom:
      secretKeyRef:
        name: bot-signup-runtime
        key: session_secret
  - name: DISCORD_CLIENT_ID
    valueFrom:
      secretKeyRef:
        name: bot-signup-runtime
        key: discord_client_id
  - name: DISCORD_CLIENT_SECRET
    valueFrom:
      secretKeyRef:
        name: bot-signup-runtime
        key: discord_client_secret
  - name: DISCORD_REDIRECT_URL
    valueFrom:
      secretKeyRef:
        name: bot-signup-runtime
        key: discord_redirect_url
  - name: SECURE_COOKIES
    value: "true"
```

Readiness/liveness probes should initially use:

```yaml
readinessProbe:
  httpGet:
    path: /api/health
    port: http
livenessProbe:
  httpGet:
    path: /api/health
    port: http
```

Ingress should use:

```yaml
spec:
  ingressClassName: traefik
  tls:
    - hosts:
        - bot-vibing.yolo.scapegoat.dev
      secretName: bot-signup-tls
  rules:
    - host: bot-vibing.yolo.scapegoat.dev
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: bot-signup
                port:
                  number: 80
```

## App repository work performed after direnv/source cleanup

Added deployment infrastructure in the app repository:

```text
Dockerfile
.dockerignore
.github/workflows/publish-image.yaml
deploy/gitops-targets.json
scripts/open_gitops_pr.py
Makefile docker-build/docker-smoke targets
```

Validation that passed before the runtime stop:

```text
go test ./...
pnpm --dir ui lint
pnpm --dir ui build
BUILD_WEB_LOCAL=1 make build
make docker-build IMAGE_REPOSITORY=bot-signup IMAGE_TAG=k3s-local
make docker-smoke IMAGE_REPOSITORY=bot-signup IMAGE_TAG=k3s-local  # image build and --help checks passed
```

The deeper container server smoke failed when running `bot-signup serve`, so live rollout is blocked until the ServeMux conflict is fixed in application code.

Remaining app repository work:

1. Fix the conflicting SPA route registration in `cmd/bot-signup/main.go`.
2. Re-run:
   - `go test ./...`
   - `pnpm --dir ui lint`
   - `pnpm --dir ui build`
   - `BUILD_WEB_LOCAL=1 make build`
   - `make docker-smoke IMAGE_REPOSITORY=bot-signup IMAGE_TAG=k3s-local`
   - a container HTTP smoke against `/api/health`.
3. Create/push the GitHub repository or add the correct remote. The current local repository has no `origin` remote.
4. Configure `GITOPS_PR_TOKEN` for the source repository if GitOps PR automation should run.
5. Push to `main` so GHCR publishes the initial image.

## First rollout checklist

Do not run this until the container server smoke is fixed and a real GHCR image exists. After the source app is stable, image publishing works, and GitOps manifests are merged:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
kubectl kustomize gitops/kustomize/bot-signup
kubectl apply --dry-run=server -k gitops/kustomize/bot-signup
export KUBECONFIG=$PWD/.cache/kubeconfig-tailnet.yaml
kubectl apply -f gitops/applications/bot-signup.yaml
kubectl -n argocd annotate application bot-signup argocd.argoproj.io/refresh=hard --overwrite
kubectl -n argocd get application bot-signup
kubectl -n bot-signup get pods,svc,ingress,pvc,secrets
kubectl -n bot-signup logs deploy/bot-signup
```

HTTP smoke:

```bash
curl -i https://bot-vibing.yolo.scapegoat.dev/api/health
curl -i https://bot-vibing.yolo.scapegoat.dev/
```

OAuth smoke:

```text
https://bot-vibing.yolo.scapegoat.dev/auth/discord/login
```

Expected OAuth behavior:

- Discord OAuth app has exact redirect URL `https://bot-vibing.yolo.scapegoat.dev/auth/discord/callback`.
- Browser returns to the application callback.
- Secure session cookie is set.
- `/api/auth/me` returns the logged-in user.

## Blockers / questions for the next operator

- Live rollout is blocked by the `bot-signup serve` runtime panic caused by conflicting `GET /` and `GET /{filepath...}` ServeMux patterns.
- The local source repository has no Git remote, so GHCR publishing cannot run until a GitHub repository/remote exists.
- The GitOps manifests currently point at `ghcr.io/go-go-golems/bot-signup:main`; once CI publishes an immutable SHA tag, the GitOps image should be bumped to `sha-<short-sha>`.
- The app currently uses SQLite on a PVC. This is simple, but only one writer replica should run. Do not scale above one replica without moving state to a server database or verifying SQLite locking semantics over the storage class.
- The health endpoint is `/api/health`, while several existing GitOps examples use `/health`.
