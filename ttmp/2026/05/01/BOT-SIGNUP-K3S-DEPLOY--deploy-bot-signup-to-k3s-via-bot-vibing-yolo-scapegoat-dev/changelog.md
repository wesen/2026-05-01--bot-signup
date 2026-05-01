# Changelog

## 2026-05-01

- Initial workspace created


## 2026-05-01

Created deployment plan and diary, then stopped before implementation because the source repo already had uncommitted code/web changes per user guardrail.

### Related Files

- /home/manuel/code/wesen/2026-05-01--bot-signup/ttmp/2026/05/01/BOT-SIGNUP-K3S-DEPLOY--deploy-bot-signup-to-k3s-via-bot-vibing-yolo-scapegoat-dev/design/01-bot-signup-k3s-deployment-implementation-plan.md — Deployment implementation plan and stop-condition record
- /home/manuel/code/wesen/2026-05-01--bot-signup/ttmp/2026/05/01/BOT-SIGNUP-K3S-DEPLOY--deploy-bot-signup-to-k3s-via-bot-vibing-yolo-scapegoat-dev/reference/01-deployment-diary.md — Chronological diary of references found and stop decision


## 2026-05-01

Resumed after direnv fix: added Docker/GHCR/GitOps scaffolding, seeded Vault runtime and image-pull secrets, validated builds and Kustomize render, then stopped before rollout because the container panics on conflicting ServeMux SPA routes.

### Related Files

- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/bot-signup/deployment.yaml — Desired-state Deployment added but not applied
- /home/manuel/code/wesen/2026-05-01--bot-signup/Dockerfile — Production image build added and validated
- /home/manuel/code/wesen/2026-05-01--bot-signup/cmd/bot-signup/main.go — Runtime route conflict blocking rollout


## 2026-05-01

Promoted the first production Discord user to admin via a guarded SQLite maintenance Job and added a production SQLite access playbook.

### Related Files

- /home/manuel/code/wesen/2026-05-01--bot-signup/ttmp/2026/05/01/BOT-SIGNUP-K3S-DEPLOY--deploy-bot-signup-to-k3s-via-bot-vibing-yolo-scapegoat-dev/playbook/01-production-sqlite-access-guide.md — Reusable production DB access and admin promotion guide

