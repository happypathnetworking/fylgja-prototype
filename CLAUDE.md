# Fylgja — working notes for Claude

A walking twin of network intent: Infrahub branch (+ optional point in time) → snapshot →
pure compiler → deployable bundle → one containerlab twin, provisioned through Temporal.
One twin at a time. No database. **Two deployables** ([D-065](docs/decisions.md#d-065)):
the **core**, `fylgja`, which builds, steps, destroys, shows and verifies a twin from a
snapshot and knows no branch, stage or source of truth, every command of it a request to
the API's server, `fylgja serve`, on the lab host; and the **orchestrator**, a Python
project on Prefect under `orchestrator/`, whose CLI `fylgja-cli` runs every command that
names a branch or a stage (`intent read`, `schema check`, `twin create`, `twin step`,
`twin rehearse`, `stage plan`, each a flow, and `stage list`): its Infrahub adapter reads
the branch to a snapshot, and the flow drives the core through the API alone.

Read in this order, then the current feature under `specs/`:
`docs/brief.md` (what we are building) · `docs/architecture.md` · `docs/decisions.md`
(rejected alternatives live here — check before proposing one) · `docs/roadmap.md`
(the releases and the backlog of features, `F-NNN`; what v0.1 built is in
`docs/archive/`) · `docs/development.md` (conventions).
Terms are in `docs/glossary.md`. What has been verified about a dependency is in
`docs/verified-facts.md`, read when a task touches that dependency. The constitution at
`.specify/memory/constitution.md` (3.0.0) binds every plan's Constitution Check. The schema
Fylgja reads is the OpsMill schema library's DCIM at a pinned commit, by one documented
mapping ([D-057](docs/decisions.md#d-057)); `schema/README.md` says which kinds and fields.

This repository began with one commit, M14's cut ([D-043](docs/decisions.md#d-043)). The
record of how M1–M13 were built (their specs, research logs and command logs) is private
and is not cited here; a bare milestone marker (`M7`) is resolved by the table below.

## Workflow

- **Spec Kit, as one chain**: `/speckit-specify` → `clarify` → `plan` → `tasks` →
  `analyze` → `implement` → `converge` (repeat until nothing is added), run by
  `scripts/converge-loop.sh` from the stage asked for (`docs/development.md`, "The feature
  chain"). **Before the chain, the operator's work is done**: the feature's `brief.md`,
  written to `docs/templates/brief-template.md`, settles every decision in `/feature-open`,
  so `clarify` has nothing to ask and an `[operator]` task sits in the route only where it
  could not precede the chain (the review, a push, an import), with its reason. Then
  `/converge-loop --from specify --feature specs/NNN-slug` runs every command as a headless
  session on its row's model and effort, gated on the tree before the next, and the
  attended session shows each session's report and carries the questions. **Each command
  runs in a fresh context, except `clarify`, which runs in `specify`'s** (by `--resume`).
  `specify` creates the feature's branch, by the git extension's hook, named by
  `GIT_BRANCH_NAME=<NNN-slug>` in the brief's prompt, which also names the directory with
  `SPECIFY_FEATURE_DIRECTORY=`. Every command reads its inputs from disk; `tasks` reads the
  tracked `.specify/templates/overrides/tasks-template.md`, whose route groups runs by
  shared context and marks the lanes the loop runs at once. `.specify/feature.json`
  (git-ignored) names the current feature's directory.
- **Commands are logged automatically** by the Claude Code hooks in `.claude/settings.json`
  (`.claude/hooks/command-log.py`): the command's prompt, report, model, effort, time,
  context and tokens, kept in `local/command-log/pending.jsonl` so no turn dirties the
  tree; **never a follow-up's text**, which is the operator's answer. `make command-log`
  publishes them into `commands/` or `specs/NNN-slug/commands/`, which are **git-ignored
  here**: it commits nothing. Run it after a hand session, once the session is closed or
  cleared. Do not write the log files by hand. The loop's stopping rule reads `tasks.md`'s
  tags, never the log.
- **Commit manually** with a conventional-commit message whose body says *why*
  (`docs/development.md`, "Conventions"). Never commit `local/`. The one exception is
  `scripts/converge-loop.sh`, which commits as it goes (each stage's artifacts, each run's
  work, each lane's merge, each convergence phase) and never pushes
  (`docs/development.md`, "The converge loop"). Its sessions get Infrahub's address and
  token from `local/.env` and run tier 2, never the node logins; tier 3, twins, the worker
  and the hand scenarios stay the operator's. Boot no twin while it runs. Two runs the route
  marks `[∥]` run at once in worktrees under `local/lanes/` and merge back mechanically;
  one that will not stops the loop at `merge` with the worktree kept. **Attended**
  (`/converge-loop`; "The attended loop"), its live sessions get the logins and run tier 3,
  the worker and the hand scenarios, every session's report reaches the operator's chat as
  it lands, and what would stop it is a question put to the operator by the session that
  started it; that session writes the answer into the tree only while the loop says it is
  paused, or once it has stopped, and a chain stage's answer through `reply-<n>`.
- **A session never pushes.** `origin` is pushed by the operator; a moved remote-tracking
  ref stops the converge loop at `contract`.
- **A changed decision** is a new entry in `docs/decisions.md` that supersedes the old one
  (append-only; the log begins at this repository's first commit), then grep
  `docs/architecture.md` for the old position. The constitution is amended only by citing
  an entry, with a version bump.
- **Verify against the running system** before asserting how Infrahub, containerlab, a
  network operating system, Temporal or Prefect behaves; record the result in the feature's
  `research.md` and, once it holds, in `docs/verified-facts.md`.
- **A changed package, import or component is redrawn** in the same change: edit
  `docs/c4/workspace.dsl` beside the prose it illustrates, then `make diagrams`.
- **The OpsMill `infrahub-*` skills** (the user's plugin `infrahub@opsmill`, not in this
  repository; its skills are named `infrahub:infrahub-managing-schemas` and so on) are
  loaded, not left to chance: when a feature designs or changes something in Infrahub — a
  schema kind, object files, a transform, a check — `plan` loads the matching skill before
  writing `data-model.md` or `research.md`, and `specify` does when it names a new kind.
  Most often `infrahub-managing-schemas`, then `-objects`, `-transforms`, `-checks`; and
  before a change to the adapter's GraphQL (`orchestrator/src/fylgja_orchestrator/adapters/
  infrahub/queries.py`, the queries the Go reader sent, ported) or its SDK client, read
  `infrahub-common`'s GraphQL references and `-transforms`' union-fragment rule. Its
  start-up hook's path triggers hold here (`infrahub/templates/`, `infrahub/queries/`,
  a check's `.py`): read the skill it names before editing. Apply their rules, and do not
  obey them where this repository already decides otherwise. The hook and its rules say
  otherwise in three places, and this file wins in each:
  - **`infrahubctl` is `.venv/bin/infrahubctl`.** Skip the hook's `uv run` / `poetry run`
    detection; neither tool is used here, and nothing is `uv add`ed or `poetry add`ed.
    The venv's `infrahub-sdk` is pinned (Environment, Python).
  - **The schema Fylgja reads loads on `main`** ([D-028](docs/decisions.md#d-028)),
    against `workflow-branch-first`: the schema library's DCIM, vendored under
    `schema/dcim/` ([D-057](docs/decisions.md#d-057)), and beside it the branch-agnostic
    kinds, `FylgjaStage` and any added after it.
  - **Some writes go to `main` on purpose**, against `workflow-branch-for-crud`: the
    library and the stage kind, the group `fylgja-devices` and the repository registration
    (`fylgja-fixture -prepare-main`, which loads the stage and probe kinds wherever absent,
    then the library only on a `main` holding none of its kinds);
    stages and other branch-agnostic objects, which land in no branch's diff; and what
    a check posts on a proposed change; and a demo's three-node baseline, for a window
    `-seed-main` opens and `-clear-main` closes ([D-054](docs/decisions.md#d-054)), with no
    tier 2 or tier 3 while it is open. Every other write goes to a branch, as the rule
    says: the fixture's own branch, a test's `fylgja-test-*` one. This Infrahub is not
    throwaway, since its volumes hold the admin token.

  They do not replace verifying against the running 1.11.2.

## Environment

**The development host** is one machine: Ubuntu 26.04.1 LTS, kernel 7.0, under QEMU/KVM,
10 vCPUs, 32 GiB, 8 GiB swap, set up by `scripts/bring-up.sh` from this clone. It is the
only lab host until M9. Ubuntu 26.04 is the one release the script supports; the host ran
24.04.1 until the launch, and `docs/verified-facts.md` marks each fact re-verified on 26.04.

- **Credentials**: `local/.env` (git-ignored, mode `0600`, scaffolded by the bring-up
  script) holds `INFRAHUB_ADDRESS` and `INFRAHUB_API_TOKEN`, which the orchestrator reads
  (`fylgja-cli`, the runner's flows, the adapter), with the fixture tool and the contract
  tier, and never the core; both node logins — the SR Linux
  probe and push login `FYLGJA_SRLINUX_USERNAME`/`_PASSWORD` and the EOS one
  `FYLGJA_EOS_USERNAME`/`_PASSWORD`, each the image's default, as `scripts/e2e.sh` sets
  them — and `FYLGJA_API_TOKEN`, the API's token, which the server and every client's
  command share, the orchestrator's flows among the clients, and the worker never reads
  ([D-042](docs/decisions.md#d-042)); `FYLGJA_DEFAULT_SCOPE`, `tag:fylgja`, the scope a read
  that names none is read under, which the orchestrator alone reads; and the orchestrator's
  four, appended on 2026-10-30 (020's T040): `FYLGJA_ORCHESTRATOR_ROOT`, the absolute
  `local/orchestrator` (its record `twin.json` and its `snapshots/`), read by `fylgja-cli`
  and the runner; `PREFECT_API_URL`, `http://172.17.0.1:4200/api`, the platform, read by
  the same two; `PREFECT_API_AUTH_STRING`, the platform's basic auth `fylgja:<password>`,
  which their Prefect client sends and the Prefect server is given as
  `PREFECT_SERVER_API_AUTH_STRING`; and `FYLGJA_ORCHESTRATOR_AUTH`, the same value, which
  Compose hands the task workers for the checks and nothing on the host reads. It also
  holds Infrahub's Compose values, which Compose reads and Fylgja never does:
  `COMPOSE_PROJECT_NAME`, `INFRAHUB_INITIAL_ADMIN_TOKEN` (equal to `INFRAHUB_API_TOKEN`),
  `INFRAHUB_INITIAL_ADMIN_PASSWORD`, `INFRAHUB_INITIAL_AGENT_TOKEN` and
  `INFRAHUB_SECURITY_SECRET_KEY`. Never re-make it: Infrahub's volumes hold the admin
  token it carries. Both login names are their package's data, not code. Load with
  `set -a; . local/.env; set +a`. Fylgja reads the process environment only, never a
  file; `local/.env` is the only file that holds a value, and `.env.example` documents
  the names. **Never echo, log, or persist a token, a password or the secret key.**
- **Infrahub 1.11.2** runs as the Docker Compose project `fylgja-infrahub`, the host's one,
  from `local/infrahub/` (git-ignored): Infrahub's published Compose file, fetched by
  version from `https://infrahub.opsmill.io/1.11.2`, and this repository's
  `scripts/infrahub/docker-compose.override.yml` copied beside it; at
  http://localhost:8000. The override pins the three Infrahub services to `1.11.2`, since a
  stray `VERSION` would move them, caps Neo4j at heap 1g/2g and page cache 1g, and gives
  the two task workers `extra_hosts: host.docker.internal:host-gateway`,
  `FYLGJA_ORCHESTRATOR_ADDRESS: host.docker.internal:4200` and `FYLGJA_ORCHESTRATOR_AUTH:
  ${FYLGJA_ORCHESTRATOR_AUTH}`, which Compose reads from `local/.env`, so the file holds no
  value ([D-065](docs/decisions.md#d-065), in D-050's place): the checks reach the
  orchestrator's platform on the bridge, and the core not at all (the core's port refuses
  the connection from inside each). A changed override reaches them only when `up -d`
  recreates them (last on 2026-10-30, 020's T040), which the bring-up script's `infrahub`
  part runs whenever it copies the override. **A recreated task worker has no checkout**
  (they mount no volume), and Infrahub 1.11.2's user check cannot clone a read-only
  registration, so the first proposed change opened after a recreate can leave its user
  validators `in_progress` until its artifact check has cloned `fylgja` on both workers;
  then re-run its user checks (`docs/verified-facts.md`, 10-30). Run
  Compose from the repository root as the script does: `docker compose -p fylgja-infrahub
  --env-file local/.env -f local/infrahub/docker-compose.yml -f
  local/infrahub/docker-compose.override.yml <command>`. Volumes persist across `down`.
- **Infrahub's `main`** carries the schema Fylgja reads: the schema library's DCIM, the six
  files of `opsmill/schema-library@26336e9` vendored under `schema/dcim/`, loaded by
  `-prepare-main` on 2026-10-26 once `specs/017-dcim-reader/migration/`'s two `state:
  absent` documents had removed the reference kinds, then the four generics and their
  contract node; and the stage kind `FylgjaStage`, branch-agnostic (`label`, `as_of`,
  `description`, a required `proposed_change`), with its `fylgja_stages` extension of
  `CoreProposedChange`, loaded from `schema/fylgja-stage.yaml` on 2026-10-16 (M10's kind
  was removed on 2026-10-17 by `specs/014-proposed-change-stages/migration/`); and the
  probe kind `FylgjaProbe`, branch-agnostic (`name`, `format` `ping`, and each end's device,
  interface, address and gateway), loaded from `schema/fylgja-probe.yaml` on 2026-10-29,
  which holds no probe once a tier or a hand scenario has run. It holds no device. Hash
  **`8a4f3bee…`** on this install (`7c413536…` before the probe kind, `e43a1cde…` before
  the library,
  `a8045729…` before M10's kind's removal, `e7386419…` with
  [D-045](docs/decisions.md#d-045)'s candidate kind, `0aba31a5…` before any stage kind; a
  branch's hash is its own, so each install and each branch has one; the 24.04 install's
  was `4d5b37aa…`); the empty group
  `fylgja-devices`; and the registration `fylgja`, a `CoreReadOnlyRepository` at
  `https://github.com/happypathnetworking/fylgja-prototype.git` (this repository, renamed
  from `fylgja` on 2026-11-03, the repository cut's step 6; the registration's `location`
  updated in place on that day, and both task workers' clones' `origin` set to it by hand,
  since the update leaves them), `ref main`, with no credential
  ([D-044](docs/decisions.md#d-044)), imported at `c006bf9` on 2026-10-26, at
  `cb0a8d7` on 2026-10-27, at `842c893` on 2026-10-29 (the template's address lines), at
  `504ec98` on 2026-10-30 (020's T039: the checks as the platform's clients) and at
  `4dd6405` on 2026-11-03 from the renamed URL, its checks this tree's byte for byte,
  whose import
  made the query `device_config` (over the library's kinds), the transform
  `srlinux_device_config` and the definition `srlinux_device_config` → artifact
  `device-config`, `text/plain`, and the query `fylgja_branch` and the two checks beside
  the template in `.infrahub.yml`:
  `fylgja-step` (`infrahub/checks/fylgja_step.py`, timeout 300) and `fylgja-twin`
  (`fylgja_twin.py`, timeout 4260), global, which Infrahub's task workers run on every
  proposed change, draft or not, at its creation and on a re-run of its user checks, as
  clients of the orchestrator's platform on the bridge (`infrahub/checks/fylgja_api.py`,
  `contracts/checks.md`): each triggers a flow's deployment (`stage-plan/main`,
  `twin-rehearse/main`), waits for its run and concludes from its document, the run's
  artifact. `fylgja-twin` skips a draft and has a ready one rehearsed on a twin by the
  `twin-rehearse` flow, which posts the report as a thread on it; the check writes nothing
  ([D-052](docs/decisions.md#d-052), [D-065](docs/decisions.md#d-065)).
  `fylgja-fixture -prepare-main` loaded the library and the stage kind, made the group and
  registered the repository, each only when absent. It loads the stage and probe kinds
  first, wherever either is absent, then on a `main` holding any of the library's kinds
  refuses, exit 1, loading nothing of the library, which the bring-up script reads as the
  library present; this host's probe kind was loaded by hand
  (`.venv/bin/infrahubctl schema load --branch main schema/fylgja-probe.yaml`) before the
  tool took that order (016's T053). No credential is registered. **One template renders both platforms**: it branches on the device's
  platform, and renders each physical port's addresses (SR Linux's `subinterface 0 ipv4`
  lines; EOS's `ip routing`, `no switchport` and `ip address`). **A template or check change reaches Infrahub
  only when asked**: once it is on `origin`'s `main`, which the operator pushes,
  `InfrahubReadOnlyRepositoryImportLastCommit(data: {id})` on `main` imports it (neither
  the minute sync nor rewriting `ref` does), and both task workers fetch it about 15 s
  after the import's own flow, before which a run-check fails; a worker's fetch can fail
  and is never retried, so read both workers' `commits/` after an import
  (`docs/verified-facts.md`, 10-27); a branch then sees it only
  after `BranchRebase` or re-creation, so a branch made before a check's import never runs
  it, and each artifact changes only once generated again. Prefect's `git_repositories_sync` runs every minute with
  one concurrency slot, so a run stuck `PENDING` blocks every import until it is
  cancelled.
- **The fixture branch `fylgja-fixture`**: three devices in the library's kinds, re-seeded
  on 2026-10-29: the tag `fylgja` on each, the manufacturer `Nokia`, the platform `SR
  Linux`, four locations, 15 physical ports (`ethernet-1/1`, `ethernet-1/2`, `mgmt0`, role
  `management`, and the two probe-ready ports `ethernet-1/10` and `ethernet-1/11`,
  uncabled, addressed `10.<d>.10.1/24` and `10.<d>.11.1/24` in the default namespace, on
  each), 6 addresses and 3 cables; the devices are in group `fylgja-devices`, each with
  one `device-config` artifact, `Ready`, 1,327 bytes carrying its two ports' address lines
  (`n1` `28881d3d…`, `n2` `092ca9f6…`, `n3` `95a3a917…`), and each `description` is the
  marker. Its schema hash is `main`'s, `8a4f3bee…`, since a branch sees `main`'s schema as it was
  at the branch's creation; it is never rebased, only deleted and re-seeded. It has no
  proposed change. It compiles to **`c3872f5c…`** on this install (`f764e1fa…` before
  016's ports, `47b2c449…` in the generics' kinds before 017), and its snapshot id is
  **`0940db93…`** (`0940db93e2ad1244a097824ececdaa6979500b1ddeeef15ca3cd612da68656ef`, the
  same for every read while the branch stands and no probe is declared; `576933ec…` before
  016's ports); its CTM equals `testdata/ctm/three-node.json` in
  everything but the envelope. It is written by its seed (`make infrahub-seed`, which the bring-up
  script runs when the branch is absent) and by no test. Contract tests make their own
  throwaway `fylgja-test-*` branches, with their draft proposed changes and stages, and
  the test that wrote them deletes the stages, then the proposed changes, then the
  branch. The fixture tool refuses to open a proposed change on `fylgja-fixture` or mark a
  stage of it.
- **Stages** are the operator's to mark on a proposed change, on the default branch: by
  object file (`kind: Object`, `spec.kind: FylgjaStage`, `proposed_change` the proposed
  change's id), the UI (`Stage`, an object entry of its own) or `.venv/bin/infrahubctl
  object create FylgjaStage --set label=… --set proposed_change=<id>`; the guide is
  `schema/README.md`. After the contract tier, `CoreProposedChange`, `FylgjaStage`, `FylgjaProbe`,
  `CoreChangeThread` and `CoreUserValidator` count 0, `CoreReadOnlyRepository` holds
  `fylgja` alone with its two check definitions, and `Branch` lists `main` and
  `fylgja-fixture` alone. **Every proposed change opened here runs both checks** at its
  creation, and once it is ready for review a re-run of its user checks rehearses it on a
  twin through the orchestrator's platform on the bridge, whose flow drives the core:
  never mark one ready and re-run its checks while tier 2 runs.
- **containerlab 0.79.0** (its `.deb`; `/usr/bin/containerlab` setuid root) and **Docker
  29.8.1** with Compose 5.6.0, from Docker's repository. The user is in `clab_admins` and
  `docker`: `clab` runs without sudo. A shell begun before the bring-up script added the
  groups lacks them until the next login (read `id -nG`): run Docker, containerlab, the
  worker, the server and tier 3 from one under `sg docker`, which is all containerlab
  needs. **No `/dev/kvm`**: the host does not nest, which only a `vrnetlab_vm` package
  would need. `gnmic` 0.49.0 is installed (`/usr/local/bin`) as a hand tool.
- **The images.** SR Linux `ghcr.io/nokia/srlinux:24.7.1`, pulled. cEOS `4.32.0.2F`,
  imported by the bring-up script from the account-gated tar `cEOS-lab-4.32.0.2F.tar`
  (xz-compressed despite its name; sha256 `89a567d5…`, the script's recorded value), which
  sits in the clone's parent directory: the script looks at `--ceos-tar`, `local/` and
  there, never in the clone's root, from which the loop commits. `ceos:4.32.0.2F` has one
  layer, `sha256:09ab9635…`, the decompressed tar's sha256, and no `Cmd`, as `docker
  import` leaves it; each import gets a new image id. **cEOS is never pulled**: its
  package says `acquisition: account_gated`, so the host check verifies presence under
  exactly that reference and refuses the run when it is absent — the same image under
  another tag is absent. `alpine:3.24.2`, the `ping` probe package's image (id and
  repository digest `sha256:294b683c…`, busybox 1.37.0), pulled by hand: a run never pulls
  a probe image, whatever its acquisition, and the host check refuses a probed run without
  it.
- **The host's AppArmor must allow SR Linux's `rsyslogd`.** Ubuntu's
  `/etc/apparmor.d/usr.sbin.rsyslogd` attaches to the `rsyslogd` inside every privileged
  SR Linux container; unwidened, every `nokia_srlinux` deploy fails at containerlab's
  post-deploy commit (`Applications have failed: log_mgr`).
  `/etc/apparmor.d/local/usr.sbin.rsyslogd` carries `/opt/srlinux/** mr,`,
  `/run/srlinux/** rw,` and `/run/syslogd.pid* rw,`, which the bring-up script's lab part
  writes and reloads with `apparmor_parser -r`. The kernel log's `net_admin` denials of
  SR Linux's `rsyslogd` are harmless: every deploy passes.
- **Go**: Ubuntu 26.04's `golang-go`, go1.26.8 since an update (`golang-1.26-go
  1.26.8-1~ubuntu0.26.04.2`; go1.26.0 at the launch); `go.mod` names `go 1.26.0`, so
  nothing is fetched. The module is `github.com/happypathnetworking/fylgja`.
- **golangci-lint 2.14.0** (`/usr/local/bin`; the v2 line — `.golangci.yml` declares
  `version: "2"`, which v1 rejects). **ShellCheck 0.11.0**, 26.04's package: CI's lint job
  runs it over `scripts/bring-up.sh` and `scripts/e2e.sh`.
- **Temporal CLI 1.9.1** (Server 1.32.0) at `~/.temporalio/bin/temporal`, where its
  installer puts it and which no login shell's `PATH` names: `make temporal-dev` and
  `scripts/e2e.sh` look there when the `PATH` has none. The dev server is `temporal server
  start-dev --db-filename local/temporal.db` on `:7233`, UI `:8233` (`make temporal-dev`).
  Start it detached, as the worker is, logging to `local/temporal.log`, and restart it if
  it is down.
- **Python**: `.venv` at the repository root (Python 3.14.4, `python3.14-venv`) holds
  `infrahub-sdk[ctl]` 1.23.2, which connects as Admin: `infrahubctl` and nothing else.
  **`orchestrator/.venv`** (`make orchestrator-venv`, the package installed editable from
  `orchestrator/pyproject.toml`'s pins) holds the orchestrator: Prefect 3.8.8 with `fastapi`
  0.150.0 beside it (3.9.0 does not co-install with the SDK, and 3.8.8's server does not
  start on the FastAPI pip would pick), `infrahub-sdk` 1.23.2, `httpx` 0.28.1, `pytest`
  9.2.0 and `ruff` 0.17.3, and the console scripts `fylgja-cli` and `fylgja-orchestrator`.
  The repository does not ignore either venv; a `.gitignore` of `*` inside each keeps the
  tree clean, so a new venv needs one too.
- **sudo** is sudo-rs 0.2.13, and a session cannot answer its prompt (`sudo -n` asks for
  a password): a step that needs root is the operator's. Every part of the bring-up script
  runs `sudo -v` first, so a session takes `--part orchestrator`'s steps by hand.
- **Git**: this clone's `user.name`/`user.email` are set in its own config, as the
  history's author.
- **Memory**: 32 GiB leaves tier 3 room. Still read `free -m` before tier 3. Infrahub
  idles at about 5.2 GiB, leaving about 24 GiB available with the core's three processes
  running, and its Neo4j is capped (above) because uncapped it once took the headroom a
  twin needs and killed a tier-3 run. Infrahub grows with uptime: after 14 days its
  containers held about 15.4 GiB (its server 7.97 GiB), leaving about 10.5 GiB, and a
  tier-3 attempt lost a deploy to an SR Linux `sr_cli` segfault (2026-10-27). A Compose
  `restart` of its containers brought it back to about 21.5 GiB available, keeping its
  volumes, both task workers' repository checkouts and their reach to the platform on the
  bridge; with the five processes, 15.8 GiB became 21.2 GiB, and 18.0 GiB 19.8 GiB
  (2026-10-30). When `free -m` reads well below about 20 GiB, restart it before a live run.
  Do not boot a twin while the contract tier runs, and rest the host between tier-3 runs.

**Five processes on the host**: Temporal's dev server, the core's worker and server, and
the orchestrator's Prefect server and runner ([D-065](docs/decisions.md#d-065)). Each is
yours to run.

**The worker is yours to run.** Start, stop and restart `fylgja worker run` yourself
whenever you need to — after every rebuild, before any live run, and whenever its binary
reads `(deleted)`. Do not leave it for the operator and do not ask. Start it detached so
it outlives the session, from a shell that has loaded `local/.env`, at the repository
root:

```
set -a; . local/.env; set +a; setsid nohup env FYLGJA_STATE_ROOT="$PWD/local" bin/fylgja worker run >> local/worker.log 2>&1 < /dev/null &
```

Then verify what you started, never assume it: `/proc/<pid>/exe` is not `(deleted)`, the
environment names carry both logins, and the start-up report names both platform packages and the probe package. The
worker reads no source of truth since 019, so it needs no `INFRAHUB_*`, which it still
inherits from `local/.env` and never reads. A rebuild does not always replace the binary,
so read `/proc/<pid>/exe` rather than assuming either way. The worker serves three
workflows (provision, destroy and step; the rehearsal is the orchestrator's flow since
020), each started with a snapshot id and a bundle id, and seventeen activities; a worker
on an older binary than the tree may lack one and fail a run at it. The start-up report
names the packages, the probe package among them (`ping (embedded): image alpine:3.24.2
(public_registry, present on this host), flow ping every 0.2s, burst 5`), never workflows
or activities. A stopped worker lingers as a poller for about
five minutes — harmless, and not a second worker. The dev server is the same: restart it
if it is down.

**The API's server is yours to run, as the worker is.** Every command of the core but
`worker run` and `serve` itself is a request to `fylgja serve`
([D-040](docs/decisions.md#d-040), [D-041](docs/decisions.md#d-041)); `--version` and
`--help` alone need no server. Start, stop and restart it yourself on the same occasions,
detached, from a shell that has loaded `local/.env` (which carries `FYLGJA_API_TOKEN`;
without it the server refuses to start), **on loopback alone and with every `INFRAHUB_*`
name and `FYLGJA_DEFAULT_SCOPE` removed from its environment**: it reads none of them, and
`local/.env`'s Compose values include `INFRAHUB_INITIAL_ADMIN_TOKEN`, the API token's value,
so the bring-up script removes the whole pattern (`CORE_UNSET`) and tier 3's case 9 refuses
a core holding any:

```
set -a; . local/.env; set +a; setsid nohup env $(compgen -e | grep -E '^(INFRAHUB_[A-Za-z0-9_]*|FYLGJA_DEFAULT_SCOPE)$' | sed 's/^/-u /') FYLGJA_STATE_ROOT="$PWD/local" bin/fylgja serve --listen 127.0.0.1:7650 >> local/server.log 2>&1 < /dev/null &
```

Then verify it: `/proc/<pid>/exe` is not `(deleted)`; its environment carries both logins,
`FYLGJA_STATE_ROOT` and `FYLGJA_API_TOKEN`, and no `INFRAHUB_*` and no
`FYLGJA_DEFAULT_SCOPE` (read the names alone from `/proc/<pid>/environ`); `ss -ltn` shows
`7650` on `127.0.0.1` alone; and its report (`local/server.log`) says `listening on
127.0.0.1:7650 (API version 1, build …)`, names `bundles:`, `snapshots:` and `twin:` in
its state root line (the worker's names the same), says `the workflow service
(localhost:7233, namespace default), Docker, containerlab and the nodes are dialled when a
request needs them`, with no default scope line and nothing of Infrahub, and names both
platform packages with their logins set and the `ping` probe package with its image
present. **Find both with `pgrep -x fylgja`** (`pgrep -f` also matches the shell that runs
it) and tell them apart by `/proc/<pid>/cmdline`. The server dials nothing at start and
keeps nothing between two requests, so a restart loses nothing but the answers it was
sending, whose runs go on. It logs one line a request to `local/server.log`, never a
token. `docker` and `containerlab` must be on the `PATH` of both the worker and the server,
which makes every dry run; a client needs neither.

**The orchestrator's Prefect server and its runner are yours to run, as the worker is.**
The Prefect server is the platform the runner and Infrahub's task workers meet on, on the
Docker bridge alone (`PREFECT_API_URL`, `http://172.17.0.1:4200/api`; uvicorn takes one
host, and the host's processes reach the bridge as the task workers do), its home and
SQLite database under `local/orchestrator/prefect`; the runner, `fylgja-orchestrator
serve`, serves the six flows' deployments (`intent-read`, `twin-create`, `twin-step`,
`stage-plan`, `twin-rehearse` and `schema-check`, each `main`) that the checks trigger.
Start each detached, the server first, from a shell that has loaded `local/.env`:

```
set -a; . local/.env; set +a; mkdir -p local/orchestrator/prefect; PREFECT_HOME="$PWD/local/orchestrator/prefect" PREFECT_SERVER_API_AUTH_STRING="$PREFECT_API_AUTH_STRING" setsid nohup orchestrator/.venv/bin/prefect server start --host "$(docker network inspect bridge -f '{{(index .IPAM.Config 0).Gateway}}')" --port 4200 --no-ui >> local/prefect.log 2>&1 < /dev/null &
set -a; . local/.env; set +a; PREFECT_HOME="$PWD/local/orchestrator/prefect" setsid nohup orchestrator/.venv/bin/fylgja-orchestrator serve >> local/orchestrator.log 2>&1 < /dev/null &
```

Then verify each: the server's `/api/health` answers `200` (23–32 s after its start),
`/api/admin/version` `"3.8.8"` with the auth string and `401` without; the runner's line
in `local/orchestrator.log`, `fylgja-orchestrator: serving intent-read/main,
twin-create/main, twin-step/main, stage-plan/main, twin-rehearse/main, schema-check/main
on the platform at http://172.17.0.1:4200/api` (about 4 s after its start), and each
deployment read by an authenticated `GET /api/deployments/name/<flow>/main`; neither log
holds the auth string, its password or its base64. The runner holds the code it started
with (the package is installed editable): restart it after any change under
`orchestrator/src/`. From inside each task worker, the deployment
`…/api/deployments/name/twin-rehearse/main` with `FYLGJA_ORCHESTRATOR_AUTH` answers `200`,
and `POST http://host.docker.internal:7650/v1/twin/show` refuses to connect
(`specs/020-core-orchestrator-split/quickstart.md` §6). Find them with `pgrep -f
'orchestrator/.venv/bin/prefect server start'` and `pgrep -f
'orchestrator/.venv/bin/fylgja-orchestrator serve'`. `make prefect-server` and `make
orchestrator-serve` run each in the foreground. A cancelled flow's `finally` has Prefect's
fixed 30 s grace before SIGKILL.

**A client of the core** (`fylgja`) needs `FYLGJA_API_TOKEN`, and `FYLGJA_API_ADDRESS`
when the server is off `127.0.0.1:7650`, and nothing else: run a client's command from a
shell that holds the token alone, or with `local/.env` loaded, which it ignores but for the
token. A client on another machine comes through `ssh -N -L 7650:127.0.0.1:7650 <lab
host>`, with a binary built from this tree. **`fylgja-cli`**
(`orchestrator/.venv/bin/fylgja-cli`) runs its flow in its own process and needs the API's
two, `INFRAHUB_ADDRESS` and `INFRAHUB_API_TOKEN`, `FYLGJA_DEFAULT_SCOPE` (or `--scope`)
and `FYLGJA_ORCHESTRATOR_ROOT`; with `PREFECT_API_URL` and `PREFECT_API_AUTH_STRING` its
runs are recorded on the platform, and without them on Prefect's ephemeral server under
`PREFECT_HOME`, started by each command (17–40 s a start), as tier 3 runs it. Loaded
`local/.env` gives it all of them.

**Formats**: bundle `6` (a probe's two hosts, their links and the ports they take under
the manifest's `fidelity.synthesized`; a bundle of another format is refused), CTM `2` (its version is
the contract's: contract `0.2` is retired, [D-057](docs/decisions.md#d-057)), PSP `0.7`
(each package's `platform.matches`, matched exactly; both shipped packages on `mode:
replace`), snapshot `2` (CTM `2` wrapped with its probe list, its id the hash of the body without
the read's `observed_at`, filed under `snapshots/`, [D-061](docs/decisions.md#d-061); a `1`
read as a `2` with no probes, under the same id; the core's store and the orchestrator's
each file one), probes `1` (a declaration's file,
`contracts/probes.schema.json`, [D-064](docs/decisions.md#d-064)), probe package `0.1`
(`probes/ping.yaml`, keyed by its `format`), `twin.json` `12` (the core's record: its
snapshot and bundle, no stage and no rehearsal block; a step compares the twin's snapshot's
scope and mapping with the target's; an `11` or earlier an earlier build left is shown and
destroyed, never stepped from), the orchestrator's record `1`
(`$FYLGJA_ORCHESTRATOR_ROOT/twin.json`: the stage and the proposed change the twin was
driven to, its snapshot, the flow and its run, trusted only while `twin show`'s snapshot
is its, else `record.twin.stale`), findings `1` (the core's and the orchestrator's, the
operation enum keeping `twin.rehearse`, `stage.list`, `stage.plan`, `intent.read` and
`schema.check` for the orchestrator's documents), the API `1` (`contracts/api.md`, the
version in every path; `snapshot` on `twin.create`, `twin.compile` and `twin.step`;
`detach` on `twin.create`, `twin.provision`, `twin.step` and `twin.destroy`, a detached run
read and interrupted by its identity through `run.read` and `run.interrupt`; `twin.diff`,
the step between two snapshots; the five operations `twin.rehearse`, `stage.list`,
`stage.plan`, `intent.read` and `schema.check` and the arguments `branch`, `at`, `stage`,
`scope` and `probes` gone, pre-1.0, without a version move). **The orchestrator owns**
(`contracts/orchestrator.md`) `fylgja-cli`'s commands and documents: the `stages` and
`plan` blocks (`contracts/stages.schema.json`, whose `$defs/step` is also the core's
`twin.diff` block), the `rehearsal` block (`contracts/rehearsal.schema.json`), the
rehearsal's text and markdown report (at most 4,092 bytes, its reproduce lines naming
`fylgja-cli` for the create and the steps and `fylgja` for the verify and the destroy) and
the deadline, held byte for byte to goldens captured from the Go build under
`orchestrator/tests/golden/`; and the checks' contract, `contracts/checks.md`.
`twin compile --ctm`, the bundle's CTM sidecar, `--interval` and `--no-follow` are retired,
and `fylgja`'s `--branch`, `--at`, `--stage`, `--scope` and `--probes` with the five
commands. The four goldens are three-node
**`c70589c3…`**, lossy **`10099944…`**, mixed **`44ac283f…`** and three-node-probed
**`5222063b…`**, their wrapped fixtures under `testdata/snapshot/` compiling to each, and
the live fixture compiles to `c3872f5c…` on this install, from snapshot `0940db93…`. The current contracts are in `contracts/`; a test of an
older format reads a frozen copy under its package's `testdata/contracts/`.

## Verified facts a session acts on (trust these over recollection)

Each is in `docs/verified-facts.md` with the version and date it was verified; the rest
of what has been verified — the SR Linux and cEOS nodes' paths and answers, the push and
the replace, containerlab's reconcile, the branch-agnostic kinds, proposed changes and
stages, verify's reads — is there too.
Re-verify on any upgrade.

- **Infrahub filters `created_at <= at`**: send no `at` unless the operator supplied one.
  A pinned `at` carries at most six fractional digits.
- **GraphQL takes the branch in the path**, `POST /graphql/<branch>`; a `?branch=`
  parameter there is ignored and answers for `main` with no error.
- **The schema.** `POST /api/schema/load` rebuilds the branch schema asynchronously: poll
  before seeding. Check candidates with `POST /api/schema/check`. Attribute names `class`
  and `type` are rejected, a name is 3–64 characters of `[a-z0-9_]`, a description ≤128
  characters; strict mode is on. A kind a branch's schema lacks is no GraphQL type there:
  read `GET /api/schema` first, as conformance does, kind by kind and field by field.
  The SDL is at `GET /schema.graphql?branch=`, and its field order is not stable between
  two fetches: compare SDLs structurally.
- **A node's stamp is its `Edged*` wrapper's `node_metadata.updated_at`**; the node type
  has none. A generation's query-group joins move it: about 2–7 s after a generation's
  checksums, each node the template query read for the first time joins that query's
  internal `CoreGraphQLQueryGroup`, written by the task workers' account, which then reads
  as the node's `updated_by`. Such a stamp is read as the node's `created_at`
  ([D-058](docs/decisions.md#d-058), [D-060](docs/decisions.md#d-060)). A `DcimCable`'s
  delete stamps its two end ports.
- **Infrahub regenerates no artifact on its own**: not on a data change, a group join or
  a commit import. Generate with `POST /api/artifact/generate/<definition id>?branch=`,
  then wait for the checksum to move; a regeneration to identical bytes writes nothing.
  `CoreArtifact.checksum` is the MD5 of the stored bytes.
- **A proposed change's pipeline** regenerates its source branch's artifacts at its
  creation, draft or not, and on a run-check (`CoreProposedChangeRunCheck`, `check_type:
  ARTIFACT`), never as the branch is edited; a run-check completes the artifact validator
  anew even when the bytes are identical. Infrahub refuses a proposed change's delete
  while a stage names it, and cancels an open one about 16 s after its branch is deleted:
  delete the stages, then the proposed change, then the branch (`fylgja-fixture -delete`).
- **A check runs as the registration imported it**: every registered repository's checks
  on a proposed change at its creation, draft or not, and again on
  `CoreProposedChangeRunCheck(check_type: USER)`, as the source branch sees them, so a
  definition imported after the branch's creation never runs on it. Its bound is its own
  `timeout`; one Infrahub times out leaves its validator `in_progress` for ever, blocking
  the merge. A proposed change's overview shows a comment only in a thread
  (`CoreChangeThread` holding a `CoreThreadComment`), whose text holds at most 4,092
  bytes, and deleting the proposed change takes its threads, comments and validators.
- **A read pinned at exactly an attribute's `updated_at` sees the node and not the
  value**, and one microsecond later sees it. A rebase re-stamps a branch's whole content
  at the new `branched_from` and discards its earlier history from pinned reads.
- **A `BranchCreate` straight after a `BranchDelete`** of the same name can fail with
  `graphql: None`; run it again.
- **containerlab**: decode `absLabPath` from `inspect --all --format json`, never
  `labPath`. `clab deploy --dry-run --format json` is a plan, and must run with
  `CLAB_LABDIR_BASE` at the twin directory. `/usr/bin/containerlab` is setuid root, so a
  `clab` outlives a worker killed with `kill -9`.
- **gNMI**: set the encoding (`json_ietf`), or SR Linux answers `Unimplemented`. SR
  Linux's gNMI is TLS on 57400 in the management namespace, probed at the node's
  management address; cEOS's is plaintext on 6030. cEOS refuses a gNMI Set: shut a port
  through its CLI (`docker exec clab-fylgja-<node> Cli -p 15 -c 'configure / interface
  Ethernet1 / shutdown'`).
- **The SR Linux push**: a refused line stays in the login's private candidate across
  requests, and the replace's `load startup` clears it. A refusal's message is cut at
  about 1 KiB.
- **Temporal Go SDK v1.49 abandons an activity attempt at the first heartbeat call that
  fails**, not after `HeartbeatTimeout`; a call fails once unanswered for about half of
  `MaxHeartbeatThrottleInterval`. `context.Cause(ctx)` is a `*temporal.CanceledError` for
  a cancellation the workflow asked for, and the RPC's error for a lost heartbeat.
- **No twin while tier 2 runs**: the host starves, deploys lose their heartbeats, and
  Infrahub times out the tier.
- **Prefect 3.8.8**: basic auth guards every path but `/api/health`; a flow's document
  crosses as its run's artifact, `data` a JSON string, never its return value;
  `flow_runs/filter`'s default order is by `uuid4` id, not time (sort `START_TIME_DESC`);
  `set_state CANCELLING` with `force` moves an ended run to `CANCELLING` and leaves it
  there, so send none; a cancellation's grace before SIGKILL is a fixed 30 s; the
  ephemeral server's default 20 s start-up bound lapses on a fresh home (180 is set).

## Milestones

M1 — Read and compile: the intent reader, the CTM and the pure compiler to a content-addressed bundle.
M2 — Provision and destroy: the provision and destroy workflows, the worker, the containerlab driver, readiness; one SR Linux twin.
M3 — Operate: refusal while a twin exists, orphan detection, a pinned `at` that reproduces its bundle.
M4 — Walking: following a branch by a scheduled read and compare, `--no-follow`, `twin show`; following removed by 019 ([D-062](docs/decisions.md#d-062)), so a twin walks by stepping.
M5 — Configuration: each device's artifact as Infrahub renders it, in the bundle and pushed after readiness.
M6 — PSP format and lossy mapping: the mapping profile and the conformance suite.
M7 — EOS: the second vendor as data, the account-gated image path, a mixed twin.
M8 — Object storage and Temporal cluster: not yet built.
M9 — Remote hosts: not yet built.
M10 — Waypoints: `FylgjaWaypoint`, `twin create --waypoint`, `waypoint list` and `waypoint plan`.
M11 — Stepping: `twin step` from one waypoint to another without a rebuild, the replace push, a diverged twin.
M12 — Verify: `twin verify` and the step's wait.
M13 — API: `fylgja serve` on the lab host, every command its client.
M14 — Reproducible and visible: the cut (this repository's first commit) and the launch.

## Build and test

`scripts/bring-up.sh` (sets a fresh Ubuntu 26.04 host up: the toolchain, the lab host,
Infrahub with `main` prepared and the fixture seeded, the orchestrator's venv, its
auth string and its two processes (`--part orchestrator`), then the build and the core's
processes, started before tiers 1 and 2 so that tier 2 reaches this build's server, and
left running: five processes; `--part NAME` runs one part, and `--part infrahub` is what
CI's contract job runs; a second run installs nothing; it asks for sudo, so it is the
operator's) · `make build` → `bin/fylgja` (CGO disabled; must stay static; `make build
VERSION=…` stamps `--version`) · `make orchestrator-venv` (`orchestrator/.venv` at
`orchestrator/pyproject.toml`'s pins, the package editable; CI's lint, unit and contract
jobs and `release.yml` run it after `actions/setup-python` at `3.14`) · `make test` (tier
1, no infrastructure: `go test`, then `orchestrator/.venv/bin/pytest -q orchestrator -m
'not contract'`; under 90s wall uncached, raised from 65s on 2026-10-29 by the operator's
answer when 020's Run 17 measured 68.0s and 67.0s at a load of 2.1–4.7 with the
orchestrator's 710 Python tests beside Go; at 020's close 55.9s at a load of 2.1, Go's
`internal/cli` 21.9s of it and 1,120 Python tests in 18.3s, since the reader's packages
left; CI's unit job on every push and pull request, beside its lint job, which runs
golangci-lint, `make lint-python` and `shellcheck scripts/bring-up.sh scripts/e2e.sh`) ·
`make test-contract` (tier 2, real Infrahub, never cached, in three parts that each run
whatever the one before answers: the Go packages under the tag, `internal/checks` alone,
then `pytest -m contract`; needs `INFRAHUB_ADDRESS` and `INFRAHUB_API_TOKEN`; CI's
contract job on every push to `main`, against an Infrahub the script brings up on the
runner, with no secret; **2,692.6 s whole on 2026-10-30**, the Go half about 370 s (the
first `go test`'s 21 packages, `internal/cli` 21.2 s the longest, then `internal/checks`
338.3 s, whose `TestChecksInsideInfrahub` runs the tree's checks inside Infrahub's pipeline
under a registration of its own through the orchestrator's platform on the bridge, and
skips, saying why, where the task workers reach none, as in CI), and `pytest -m contract`
2,322.3 s, the Go reader's contract cases carried over with their assertions and the flows
on `prefect_test_harness`, each against a core server the test starts with no
`INFRAHUB_*`; every artifact wait 180 s, the schema waits 90 s) · `make lint`
(golangci-lint with build tags `contract,fixture,e2e`, then `make lint-python`, ruff over
`orchestrator` and `infrahub/checks`) · `make linkcheck` (`scripts/linkcheck.py`: every
relative link and anchor in the tracked Markdown resolves in a public clone; a hand tool,
not a CI step) · `go test ./internal/compiler -update` regenerates the four goldens —
review the diff, and commit a re-baseline on its own; the orchestrator's goldens under
`orchestrator/tests/golden/` are re-baselined the same way, in a commit of their own · the
conformance suite's pure half is part of `make test`; its boot half is
`FYLGJA_STATE_ROOT=$PWD/local go test -count=1 -tags e2e ./internal/conformance -run
'^TestBootHalf$' -v` against a ready twin · `make temporal-dev` (the file-backed dev
server) · `make worker` and `make serve` (each with the state root made absolute; load
`local/.env` first) · `make prefect-server` and `make orchestrator-serve` (the platform and
the runner in the foreground; load `local/.env` first) · `make infrahub-seed` and `make
infrahub-clean` (the fixture branch) · `make diagrams` (needs Docker; no tier needs it) ·
**every `fylgja` command but `worker run` and `serve` needs a running server** (`psp
validate` and `twin compile` need neither Infrahub nor the workflow service, but they need
the server), and `fylgja-cli`'s commands need Infrahub, and the core's server wherever
they ask the core ·
`make test-e2e` (tier 3, `scripts/e2e.sh`: nine cases, nine twins, the boot half on three
of them (cases 1, 5 and 6, case 5's twin probed), `twin verify` on three of them and on
the steps' waits, case 5 a twin left 330 s with nothing rebuilding it and then stepped to
a snapshot by file and by id, probes in cases 5, 8 and 9 (case 5's `n1-loop` by file, case
8's `e1-loop` from the kind across `e1`'s restart, case 9's stage `probe-port-down` failing
on `verify.burst`); in cases 1 to 8 every command naming a branch or a stage through
`fylgja-cli` on Prefect's ephemeral server under the run's own home, the rest through
`fylgja`, against a core the script starts itself with a token it makes and no
`INFRAHUB_*`; case 9, a proposed change rehearsed twice by `fylgja-twin` inside Infrahub,
through the orchestrator's platform on the bridge and the operator's core on loopback; the
core's server logs asserted to name no Infrahub operation; needs the dev server and a
worker running from the repository root, Docker, containerlab, **both** images,
`alpine:3.24.2`, `orchestrator/.venv` and Infrahub with `FylgjaProbe` on `main` holding no
probe, and for case 9 the platform and the runner from this tree, the operator's core on
loopback from this tree with no `INFRAHUB_*`, and both checks imported on `main` byte for
byte as this tree has them, each checked before any case; about 47 minutes (2,798.6 s on
2026-10-30, by a copy whose case-9 precheck compared the checks' syntax trees, the
operator's choice, since `make test-e2e` refused this tree at that precheck; 562 s more
than 016's 2,236.6 s, every `fylgja-cli` command starting its ephemeral server, 17–40 s a
start); `WORKER_LOG=<file>` adds the worker's log to its credential greps;
`PLATFORMS=nokia_srlinux` narrows it to the seven cases SR Linux alone runs, about 19
minutes before 020 (1,156.1 s on 2026-10-27), ending `E2E-PARTIAL: platforms …; ran …;
skipped …` with exit 0, never `E2E-OK`; never in CI). A command that runs longer than one
tool call (tier 3) starts detached with its exit status in a file, and is waited on in
foreground calls of under ten minutes.

## State

Rewritten whole at each release's close, never appended to; rewritten whole at 020's.

- **Release 0.1** (`v0.1.0`, 2026-10-08): M1–M7 and M10–M14 built, M8 and M9 not; the cut,
  this repository's one first commit ([D-043](docs/decisions.md#d-043)), and the launch
  (`013-launch`), as [the archive](docs/archive/v0.1-built.md) records.
- **Release 0.2** (`v0.2.0`, closed 2026-10-21; tagged at `450c380` and published as a
  pre-release on 2026-10-22): a proposed change rehearsed on a twin before it merges
  ([the archive](docs/archive/v0.2-built.md)): **014** (F-023), a proposed change's stages
  in place of M10's waypoints ([D-045](docs/decisions.md#d-045)–[D-048](docs/decisions.md#d-048));
  **015** (F-007, F-008), `twin rehearse`, the checks `fylgja-step` and `fylgja-twin` and
  the report as a thread ([D-049](docs/decisions.md#d-049)–[D-053](docs/decisions.md#d-053));
  and the release's own work: the report made legible in Infrahub's view,
  [D-054](docs/decisions.md#d-054)'s demo window and the fixture tool's `-seed-main`,
  `-clear-main` and `-add-fourth-node`, the recording and four screenshots, the schemas'
  `$id` namespace moved to this repository's URL ([D-055](docs/decisions.md#d-055)), Spec
  Kit's own files installed per host ([D-056](docs/decisions.md#d-056)) and `make
  linkcheck`. Verified at the close: CI green on `dc7dcea`, tier 2 whole, tier 3 `E2E-OK`
  in 1,650.6 s.
- **Release 0.3**, set on 2026-10-25 as exploratory coding toward the production refactor
  ([roadmap](docs/roadmap.md#release-03)), six specs in order: F-011, F-024, F-026, F-014,
  F-027 (added on 2026-10-29 and placed before F-025 by the operator) and F-025; F-012 and
  F-022 left it. The [production shape](docs/roadmap.md#the-production-shape) it explored
  toward is decided since 020. Five specs are built:
  - **017** (`017-dcim-reader`, F-011, closed 2026-10-27): the reader read the schema
    library's DCIM at `26336e9` by one documented mapping, in place of the generics
    ([D-057](docs/decisions.md#d-057)–[D-060](docs/decisions.md#d-060); Constitution 2.0.0,
    then 2.1.0); a request's scope `tag:<name>`; PSP `0.7`'s `platform.matches`; CTM `2`,
    the contract version retired; `main` migrated and both fixtures re-seeded.
    `infrahubctl schema format` is never run on the vendored files (017 research R-22).
  - **018** (`018-snapshots`, F-024, closed 2026-10-28): the snapshot, the core's one input
    ([D-061](docs/decisions.md#d-061); Constitution 2.2.0), content-addressed and filed
    under the state root's `snapshots/` by a store that never overwrites; `twin create
    --snapshot <id|file>` and `twin compile --snapshot`; the bundle's CTM sidecar retired.
  - **019** (`019-following-removed`, F-026, closed 2026-10-28): following removed from the
    core ([D-062](docs/decisions.md#d-062)): a twin is frozen at its snapshot, the reconcile
    workflow and the Schedule `fylgja-follow` gone, `twin step --snapshot <id|file>`.
  - **016** (`016-probes`, F-014, closed 2026-10-29): probes, `ping` only
    ([D-064](docs/decisions.md#d-064); Constitution 2.3.0): a declaration outside intent,
    by file or from `FylgjaProbe` on `main`, carried in snapshot `2`; two synthesized hosts
    a probe under bundle `6`'s `fidelity.synthesized`; `probes/ping.yaml` (`0.1`,
    `alpine:3.24.2`, never pulled); a step's `step.loss` and verify's burst; the fixture
    re-seeded (`0940db93…`, `c3872f5c…`) and the goldens moved twice beside a fourth.
    Merged into `main` at `195af3a` and pushed; CI green on it (run `40894683064`).
  - **020** (`020-core-orchestrator-split`, F-027, closed 2026-10-30 but for the operator's
    review, merge and push, T048): **two deployables** ([D-065](docs/decisions.md#d-065),
    [D-066](docs/decisions.md#d-066); Constitution 3.0.0, Principle I binding the core and
    freeing the orchestrator). **The core narrowed to the snapshot**: `twin.rehearse`,
    `stage.list`, `stage.plan`, `intent.read` and `schema.check` and the arguments
    `branch`, `at`, `stage`, `scope` and `probes` gone from the API and the CLI, with the
    Go reader and genqlient, `make sdl` and `make generate`, `internal/rehearsal` and the
    rehearsal's workflow (so 015's defect, a client following a rehearsal terminated
    outside Fylgja, went with it); `detach` on the long operations, `run.read` and
    `run.interrupt` by a run's identity, and `twin.diff`, the step between two snapshots;
    `twin.json` `12`; the worker three workflows and seventeen activities; the server on
    loopback alone with no Infrahub credential. **The orchestrator** under `orchestrator/`,
    Python 3.14.4 on Prefect 3.8.8 (3.9.0 does not co-install with `infrahub-sdk` 1.23.2):
    the Infrahub adapter ported on the SDK and held to the Go reader byte for byte before
    the reader went (research §5.15); six flows, `fylgja-cli` and the runner; its record
    `1` and its snapshot store; the plan, stages and rehearsal documents and the markdown
    report held to goldens captured from the Go build (T003), their reproduce lines
    re-baselined alone (`927561b`); the checks the platform's clients, imported at
    `504ec98` (T039); the override's two variables and the task workers recreated with
    them; five processes on the host; CI running `make orchestrator-venv`, ruff and the
    Python tests. **Unmoved**: CTM `2`, snapshot `2`, bundle `6`, probes `1`, PSP `0.7`,
    the probe package `0.1`, the API `1`, findings `1`, the four goldens, the fixture's
    `0940db93…` and `c3872f5c…`, the seven replay histories, `main` and both fixture
    branches. Built by the converge loop, T001–T047, the operator's T002 (the amendment)
    and T039 (the checks' push and import) among them; 46 sessions to Run 38, about 16.1
    hours, against R-22's 24–36. **Found live**: a recreated task worker has no checkout,
    and the first proposed change after it can leave its user validators `in_progress`
    (Environment, Infrahub). Tier 1 55.9 s uncached at the close. Tier 2 whole 2,692.6 s on
    2026-10-30, green once a Python helper's run order was fixed (`a6f3fe3`) and its four
    tests re-run together. Tier 3: `make test-e2e` refused this tree at case 9's precheck,
    since the imported `infrahub/checks/fylgja_api.py` differs by T035's one comment; a
    copy of the script comparing the checks' syntax trees ran all nine cases `E2E-OK` in
    2,798.6 s, which the operator took as T043's evidence (option c), the script's byte
    rule unchanged.
- **Next**: T048, the operator's review of 020, its merge into `main`, the push and `make
  command-log`, then `converge`; then F-025, the NetBox reader, Release 0.3's last spec,
  the orchestrator's second adapter on 020's adapter contract (a tool that reads one
  source under a scope and writes a snapshot `2` file, its probes from a file), a stage as
  the snapshot the adapter took when the operator marked it
  ([D-066](docs/decisions.md#d-066)). It has no brief yet. NetBox was probed live on
  2026-10-29 (`docs/verified-facts.md`, "NetBox 4.7.4 with netbox-branching 1.2.2"): NetBox
  4.7.4 with the plugin 1.2.2 is installed under `local/netbox/` as the Compose project
  `fylgja-netbox` on loopback 8800, stopped with its volumes kept (its `main` carries the
  spike's merged change), and its token and admin password are in
  `local/netbox-spike/.env`, beside the throwaway adapter, the first draft of F-025's.
- **Open, the operator's**: whether the bring-up
  script's `infrahub` part, which recreates the task workers when it copies a changed
  override, should first have the registration cloned, and whether Infrahub's
  `run_user_check` hard-coding `CoreRepository` (`infrahub/git/tasks.py:1193`) is reported
  upstream; file 014's upstream request for named points in time
  (`specs/014-proposed-change-stages/upstream-request.md`, a public file) and a report of
  Infrahub's view: a code block's first line indented (`pre > code` inline with a left
  padding) and a zoomed Mermaid drawing cropped (`docs/verified-facts.md`, 10-21); read
  Arista's own checksum for the cEOS tar at the next download.
- **Open, for a session**: the fixture tool's `-propose` waits its 6 minutes for an
  artifact validator a branch with no change from `main` never gets (D-054's window);
  bound with `timeout` there, as the release's demo did.
