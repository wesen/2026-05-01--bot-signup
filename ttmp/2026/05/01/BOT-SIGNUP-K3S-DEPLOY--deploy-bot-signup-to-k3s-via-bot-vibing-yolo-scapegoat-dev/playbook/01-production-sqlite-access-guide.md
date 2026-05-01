---
Title: Production SQLite Access Guide
Ticket: BOT-SIGNUP-K3S-DEPLOY
Status: active
Topics:
    - deployment
    - kubernetes
    - sqlite
    - bot-signup
DocType: playbook
Intent: ""
Owners: []
RelatedFiles:
    - Path: ../../../../../../../2026-03-27--hetzner-k3s/gitops/kustomize/bot-signup/deployment.yaml
      Note: Mounts the PVC at /data and sets DB_PATH for production
    - Path: ../../../../../../../2026-03-27--hetzner-k3s/gitops/kustomize/bot-signup/persistentvolumeclaim.yaml
      Note: Defines the bot-signup-data PVC that stores the production SQLite database
    - Path: internal/database/users.go
      Note: User role/status fields updated through SQLite maintenance
    - Path: internal/server/middleware.go
      Note: Admin access behavior to review after role promotion
ExternalSources: []
Summary: Operator guide for safely inspecting and updating the bot-signup production SQLite database stored on the bot-signup-data PVC.
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# Production SQLite Access Guide

## Purpose

This playbook shows how to inspect and carefully update the production `bot-signup` SQLite database on k3s.

The production database lives on the `bot-signup-data` PVC and is mounted by the app at:

```text
/data/bot-signup.db
```

Use this guide for one-off operator tasks such as promoting the first Discord user to `admin` after they have signed in once.

## Environment assumptions

Run these commands from an operator machine with cluster access:

```bash
export KUBECONFIG=/home/manuel/code/wesen/2026-03-27--hetzner-k3s/.cache/kubeconfig-tailnet.yaml
kubectl -n bot-signup get deploy,pods,pvc,secrets
```

Expected high-level state:

```text
deployment.apps/bot-signup available
persistentvolumeclaim/bot-signup-data Bound
secret/bot-signup-runtime present
```

The application image does not include `sqlite3`, so the recommended access pattern is a short-lived Kubernetes Job using `alpine` plus `apk add sqlite`, mounting the same PVC.

## Read users

Create a one-off read Job:

```bash
cat > /tmp/bot-signup-sqlite-read.yaml <<'EOF'
apiVersion: batch/v1
kind: Job
metadata:
  name: bot-signup-sqlite-read
  namespace: bot-signup
spec:
  ttlSecondsAfterFinished: 300
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: sqlite
          image: alpine:3.20
          command: ["sh", "-lc"]
          args:
            - |
              set -eux
              apk add --no-cache sqlite
              sqlite3 /data/bot-signup.db "select id, discord_id, email, display_name, role, status from users;"
          volumeMounts:
            - name: data
              mountPath: /data
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: bot-signup-data
EOF
kubectl apply -f /tmp/bot-signup-sqlite-read.yaml
kubectl -n bot-signup wait --for=condition=complete job/bot-signup-sqlite-read --timeout=120s
kubectl -n bot-signup logs job/bot-signup-sqlite-read
```

Cleanup if desired:

```bash
kubectl -n bot-signup delete job bot-signup-sqlite-read
```

## Promote the only local user to admin

Only use this broad update immediately after first login when there is exactly one user.

```bash
cat > /tmp/bot-signup-sqlite-admin-oneshot.yaml <<'EOF'
apiVersion: batch/v1
kind: Job
metadata:
  name: bot-signup-sqlite-admin-oneshot
  namespace: bot-signup
spec:
  ttlSecondsAfterFinished: 300
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: sqlite
          image: alpine:3.20
          command: ["sh", "-lc"]
          args:
            - |
              set -eux
              apk add --no-cache sqlite
              echo 'Users before:'
              sqlite3 /data/bot-signup.db "select id, discord_id, email, display_name, role, status from users;"
              count=$(sqlite3 /data/bot-signup.db "select count(*) from users;")
              if [ "$count" != "1" ]; then
                echo "Expected exactly one user, found $count; refusing broad promote" >&2
                exit 2
              fi
              sqlite3 /data/bot-signup.db "update users set role = 'admin', updated_at = datetime('now');"
              echo 'Users after:'
              sqlite3 /data/bot-signup.db "select id, discord_id, email, display_name, role, status from users;"
          volumeMounts:
            - name: data
              mountPath: /data
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: bot-signup-data
EOF
kubectl apply -f /tmp/bot-signup-sqlite-admin-oneshot.yaml
kubectl -n bot-signup wait --for=condition=complete job/bot-signup-sqlite-admin-oneshot --timeout=120s
kubectl -n bot-signup logs job/bot-signup-sqlite-admin-oneshot
```

Cleanup if desired:

```bash
kubectl -n bot-signup delete job bot-signup-sqlite-admin-oneshot
```

## Promote one Discord account by ID

This is safer once more than one user exists. Replace `YOUR_DISCORD_ID` first.

```bash
DISCORD_ID='YOUR_DISCORD_ID'
cat > /tmp/bot-signup-sqlite-promote-one.yaml <<EOF
apiVersion: batch/v1
kind: Job
metadata:
  name: bot-signup-sqlite-promote-one
  namespace: bot-signup
spec:
  ttlSecondsAfterFinished: 300
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: sqlite
          image: alpine:3.20
          command: ["sh", "-lc"]
          args:
            - |
              set -eux
              apk add --no-cache sqlite
              sqlite3 /data/bot-signup.db "select id, discord_id, email, display_name, role, status from users;"
              sqlite3 /data/bot-signup.db "update users set role = 'admin', updated_at = datetime('now') where discord_id = '${DISCORD_ID}';"
              sqlite3 /data/bot-signup.db "select id, discord_id, email, display_name, role, status from users where discord_id = '${DISCORD_ID}';"
          volumeMounts:
            - name: data
              mountPath: /data
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: bot-signup-data
EOF
kubectl apply -f /tmp/bot-signup-sqlite-promote-one.yaml
kubectl -n bot-signup wait --for=condition=complete job/bot-signup-sqlite-promote-one --timeout=120s
kubectl -n bot-signup logs job/bot-signup-sqlite-promote-one
```

Cleanup:

```bash
kubectl -n bot-signup delete job bot-signup-sqlite-promote-one
```

## Current first-admin promotion result

On 2026-05-01, the first production user was promoted with the guarded one-user Job.

Before:

```text
1|363877777977376768|wesen@ruinwesen.com|slono|user|waiting
```

After:

```text
1|363877777977376768|wesen@ruinwesen.com|slono|admin|waiting
```

## Access the admin UI

After promotion, refresh:

```text
https://bot-vibing.yolo.scapegoat.dev/admin
```

If the old session still appears non-admin, log out and sign in with Discord again so the frontend refetches session/user state.

## Failure modes

### Job cannot mount the PVC

Check that the PVC exists and is bound:

```bash
kubectl -n bot-signup get pvc bot-signup-data
```

The app uses a `ReadWriteOnce` local-path PVC. On a single-node cluster this one-off Job can mount it, but for future multi-node layouts, run maintenance on the same node as the app or scale the app down first.

### Broad promote refuses to run

If the Job exits with:

```text
Expected exactly one user, found N; refusing broad promote
```

use the Discord-ID-specific promotion flow instead.

### App appears stale after DB update

Try:

```bash
kubectl -n bot-signup rollout restart deploy/bot-signup
kubectl -n bot-signup rollout status deploy/bot-signup
```

Then refresh the browser or log out and back in.

## Exit criteria

- `kubectl -n bot-signup logs job/<job-name>` shows the target user with `role=admin`.
- `https://bot-vibing.yolo.scapegoat.dev/admin` loads for the promoted Discord account.
