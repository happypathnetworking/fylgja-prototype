# Decision log

One file, one entry per decision that shapes Fylgja. Each entry records its context, the
options considered, the decision and its consequences. **The options considered are the
point**: [architecture.md](architecture.md) describes what was chosen, and this log keeps
the paths not taken and why, which is otherwise lost once the discussion that weighed them
ends. Before proposing an alternative, look for it here; it may already have been weighed.

**Statuses.** *Accepted*: in force, dated the day it was decided. *Superseded by D-NNN*:
replaced by a later entry, which says what changed and why. *Parked*: a direction neither
adopted nor rejected until the scope that needs it is planned (the [Parked](#parked)
section). The log is append-only: a change is a new entry that supersedes the old one, the
old one gains its *Superseded* status and nothing else, and a change to a principle is then
a constitution amendment that cites the new entry. This repository's log begins at its
first commit. Every entry below was decided before it and is restated as in force, with
the date it was decided and what later work established about it folded into its body,
without the record that produced it, which is private ([D-043](#d-043)). So no entry here
is superseded.

**Citing an entry.** Cite an entry by its number, `D-0NN`, in a code comment, a document
or a commit message, and link to it by its anchor, `decisions.md#d-0nn` from `docs/` or
`docs/decisions.md#d-0nn` from the repository root. The anchor is the number in lower case
and does not change when a title would; a number is never reused, and a new entry takes the
next one.

## Index

| ID | Decision | Area |
|---|---|---|
| [D-001](#d-001) | Infrahub as the network source of truth | NSoT |
| [D-002](#d-002) | Schema contract as Infrahub generics; a reference schema is shipped | Contract |
| [D-003](#d-003) | Interface contract: `iftype` is kind, `mgmt_only` is the one use-flag | Contract |
| [D-004](#d-004) | Reserve only confident structure: keep `parent`, defer the VLAN generic | Contract |
| [D-005](#d-005) | Go end to end; no second runtime | Language |
| [D-006](#d-006) | Temporal, self-hosted, for durable provisioning | Orchestration |
| [D-007](#d-007) | Short workflows with fixed ids; no lifecycle workflow | Orchestration |
| [D-008](#d-008) | Drive the `clab` CLI; do not link containerlab | Orchestration |
| [D-009](#d-009) | Fylgja never renders configuration; bootstrap from the PSP | Config |
| [D-010](#d-010) | Keep the CTM, thin | Data model |
| [D-011](#d-011) | The compiler is pure; the bundle is deterministic, content-addressed and deployable as emitted | Compiler |
| [D-012](#d-012) | The intent reference is `(branch, at?)`; an unpinned reference is recorded, not pinned, and its twin follows the branch | Intent |
| [D-013](#d-013) | One twin at a time | Twin |
| [D-014](#d-014) | State lives in containerlab, files and Temporal; no database | State |
| [D-015](#d-015) | One task queue until the host move | Orchestration |
| [D-016](#d-016) | Multi-vendor by requirement; platforms are data (PSPs) | Platforms |
| [D-017](#d-017) | Test architecture: pure core, Temporal test suite, real Infrahub, end-to-end on demand | Testing |
| [D-018](#d-018) | Single module, single binary in three roles, PSPs and schema embedded | Repository |
| [D-019](#d-019) | Noun-verb CLI; `twin create` is the front door; stages exposed | CLI |
| [D-020](#d-020) | One lab host, the development host, until the host move | Infrastructure |
| [D-021](#d-021) | SR Linux first, EOS second; later order parked | Platforms |
| [D-022](#d-022) | Every bundle carries a fidelity manifest | Fidelity |
| [D-023](#d-023) | `observed_at` is recorded outside the bundle | Intent |
| [D-024](#d-024) | `bundle_id` is the canonical-bytes hash; the bundle carries no build stamp | Compiler |
| [D-025](#d-025) | `fidelity.production_forwarding` is asserted per platform, one entry per package | Compiler |
| [D-026](#d-026) | A following twin detects change by a scheduled read and compare, and nothing else | Intent |
| [D-027](#d-027) | `--no-follow`, not `--pin`, is the opt-out from following | CLI |
| [D-028](#d-028) | Configuration is rendered in Infrahub and enters the bundle at read time | Config |
| [D-029](#d-029) | A support package declares whether readiness waits for the push transport | Provisioning |
| [D-030](#d-030) | The heartbeat cadence is 10s, and cancellation latency is bounded by it | Provisioning |
| [D-031](#d-031) | A mechanism is code keyed by a format value; a platform is data | Platform |
| [D-032](#d-032) | A waypoint is an operator-written name for a pinned reference | Intent |
| [D-033](#d-033) | A replace resets the candidate to the node's own baseline: S1 on SR Linux, E1′ on EOS | Config |
| [D-034](#d-034) | A refused SR Linux replace asks the node again for its reason, in a request that lands nothing | Config |
| [D-035](#d-035) | A commit-time refusal on SR Linux gives no reason under the replace; the request after it stays | Config |
| [D-036](#d-036) | A failed or cancelled step leaves the twin up and diverged; a diverged twin accepts only `twin destroy`; a restart needs `--allow-restart` | Provisioning |
| [D-037](#d-037) | Observed facts are a Fylgja-owned vocabulary, mapped to each platform's paths by its package | Verification |
| [D-038](#d-038) | Assertions are data derived from the bundle; no operator-written assertion is shipped | Verification |
| [D-039](#d-039) | Verify is advisory: no Fylgja operation reads the report, and the report is the findings document with a `verify` block | Verification |
| [D-040](#d-040) | The API is the boundary of the core: a user-facing client reaches Fylgja through it alone, and test code is excepted | API |
| [D-041](#d-041) | `fylgja serve` runs on the lab host, as a role of its own, and reads the host in its own process until the host move | API |
| [D-042](#d-042) | The API's token is one shared secret from the process environment, sent as a bearer credential and compared before anything else | API |
| [D-043](#d-043) | The development record stays private: this repository began as one commit of cleaned code and rewritten documents, and development continues here | Repository |
| [D-044](#d-044) | The artifacts template lives in this repository, and Infrahub registers it read-only on a ref with no credential | Config |
| [D-045](#d-045) | A stage is a moment of a proposed change: the branch point read on `main`, the operator's marks, and the head | Intent |
| [D-046](#d-046) | A stage's generation stamp is the later of its artifact's checksum stamp and the artifact validation; staleness per device | Intent |
| [D-047](#d-047) | A stage is orphaned when its proposed change no longer names it | Intent |
| [D-048](#d-048) | The head and the staleness rule also read the branch's intent events, so a removal is seen | Intent |
| [D-049](#d-049) | A check posts the rehearsal's report as a comment; Fylgja writes nothing | Verification |
| [D-050](#d-050) | The API's server listens on the host's Docker bridge beside loopback, in plaintext, under the one token | API |
| [D-051](#d-051) | `twin rehearse` is a fifth short workflow whose children are provision, step and destroy, under a derived deadline | Orchestration |
| [D-052](#d-052) | The check concludes from the report's status alone; Fylgja's verdicts stay advisory | Verification |
| [D-053](#d-053) | The check's comment is a thread on the proposed change, as its view shows one | Verification |
| [D-054](#d-054) | For a demo, `main` carries the fixture's three-node network for a window the fixture tool opens and closes | Testing |
| [D-055](#d-055) | The contracts' `$id` namespace is this repository's GitHub URL | Contract |
| [D-056](#d-056) | Spec Kit's files are installed on each host, not carried in the tree | Repository |
| [D-057](#d-057) | The CTM is the contract; a reader maps one source's model onto it; the Infrahub reader reads the schema library's DCIM at a pinned commit | Contract |
| [D-058](#d-058) | A stamp Infrahub's own account wrote is not intent's | Intent |
| [D-059](#d-059) | Infrahub's own account is the one that writes the template query's groups, read from Infrahub | Intent |
| [D-060](#d-060) | Infrahub's own account is the one that created the template query's groups | Intent |
| [D-061](#d-061) | The snapshot is the core's one input: a content-addressed wrap of the CTM, stored beside the bundles; the bundle's CTM sidecar retired | Contract |
| [D-062](#d-062) | Following leaves the core: a twin is frozen at its snapshot and changes only when the operator steps it | Orchestration |

---

<a id="d-001"></a>
## D-001 — Infrahub as the network source of truth

**Status:** Accepted 2026-09-14; superseded by [D-066](#d-066) on 2026-10-29 in the part that says NetBox has no native branching; Infrahub stays the first source of truth.

### Context

Fylgja builds a twin from network intent, so the network source of truth (NSoT) shapes the
whole design. Two requirements: **branching** (the source of truth itself is versioned) and
a **flexible schema**.

### Options considered

- **NetBox / Nautobot.** Mature, large ecosystem, fixed DCIM-shaped model. No native
  branching: a change is a configuration diff somewhere else, so the twin has no
  addressable object to be built *of*. Fails both requirements.
- **Infrahub.** Native branching with diff and merge, user-defined schema with generics,
  GraphQL, branch-scoped artifacts, time-travel queries. Younger ecosystem; operationally
  heavy (graph database, message bus, task workers, cache).
- **In-house.** Full control, no ecosystem; a source of truth is not the product.

### Decision

Infrahub.

### Consequences

- Intent has an address: `(branch, at?)`, the intent reference (D-012).
- "What did my change do?" is a twin of `main`, then a twin of the change's branch
  (D-013); their bundles can be compared at compile time regardless.
- Time travel is free: a twin of last Tuesday is an incident-reconstruction tool.
- Infrahub artifacts give "Fylgja never renders configuration" a home (D-009, D-028).
- A flexible schema means Fylgja can assume no kinds exist, hence the generics contract
  (D-002).
- A twin's lifetime is **not** coupled to its branch's. A twin of a merged branch may be
  kept for inspection; coupling is an [open question](#open-questions), opt-in if ever.

---

<a id="d-002"></a>
## D-002 — Schema contract as Infrahub generics; a reference schema is shipped

**Status:** Accepted 2026-09-14; superseded by [D-057](#d-057) on 2026-10-25.

### Context

Infrahub's schema is user-defined. Fylgja needs a stable shape to query and to generate
typed Go against. The install it was designed against was new, with no schema of its own.

### Options considered

- **Fylgja owns the schema.** Fastest to a first twin; presumes Fylgja dictates the NSoT's
  shape, which an NSoT that also serves provisioning, IPAM and monitoring should not
  accept.
- **Map onto the organisation's schema.** Good citizen; field renaming cannot handle
  structural differences without becoming a mini-ETL language.
- **A published profile with two implementations.** Honest; most work; the contract stays
  under-specified and mapping becomes the de facto path.
- **Contract as Infrahub generics.** Fylgja defines generics; the organisation's concrete
  kinds implement them; Fylgja queries the generic. The previous option done in
  Infrahub's own type system, with a stable GraphQL shape for `genqlient`.

### Decision

Generics. Because the install had no schema, Fylgja also ships a **reference concrete
schema**. Three rules: demand the least possible; graduate, don't gate; never reach past
the generic.

### Consequences

- `genqlient` generates against the generics at build time; a breaking schema change fails
  the build, not the twin.
- Read validates schema conformance *and* data completeness; `fylgja schema check` runs the
  first alone.
- The contract is versioned, and the version is recorded in every provenance block.
- Generics only help where the organisation's model can express the concept; a missing
  link object needs modelling work in Infrahub, not a mapping layer.
- What Infrahub does here is in [verified-facts.md](verified-facts.md): a generic query
  resolves across implementing kinds, `GET /api/schema`'s `used_by` is the conformance
  check, and an address's peer is the core `BuiltinIPAddress`.

---

<a id="d-003"></a>
## D-003 — Interface contract: `iftype` is kind, `mgmt_only` is the one use-flag

**Status:** Accepted 2026-09-14; superseded by [D-057](#d-057) on 2026-10-25.

### Context

The compiler must know which interfaces become veth links, which one maps onto the node's
containerlab management interface, and which are configured but never cabled. Names matter
because organisations map existing fields onto the generics, and a familiar name with
different semantics produces a wrong mapping that validates cleanly.

### Options considered

For the kind attribute: **`role`** (already means uplink/downlink/peering; left unclaimed
for that meaning); **`type`** (OpenConfig-aligned, but NetBox's media type, a Go keyword,
and refused by Infrahub as reserved); **`kind`** (Infrahub's own term for concrete kinds);
**`class`** (chosen first; refused by Infrahub as a Python keyword); **`wiring` /
`attachment`** (name the compiler's behaviour, pushing implementation into someone else's
schema); **`interface_type`** (redundant on `FylgjaInterface`); **`iftype`** (IANA's term
for exactly this; OpenConfig's interface `type` values are IANA ifType identities;
collides with nothing).

For management: **address on the device** (answers "how is it reached", not "which port is
OOB"; breaks in-band management; right later as a graduated `primary_address`);
**management-interface pointer on the device** (structure you are guessing at);
**`management` as a kind value** (a use masquerading as a kind, the conflation `role` was
rejected for); **a boolean on the interface** (NetBox's `Interface.mgmt_only` has the
identical meaning, so the naming rule allows it).

### Decision

`iftype` with values `physical | loopback | svi | subinterface` (later `lag`), kind only.
`mgmt_only`, boolean, default false, marks the OOB port. Management addressing is interface
addressing. In-band management is deferred; `primary_address` is the expected later
addition.

### Consequences

- **Naming rule:** never claim a name that carries established operational meaning unless
  the same meaning is intended. **Corollary:** check candidates against Infrahub's reserved
  names (`POST /api/schema/check`) before deciding.
- Wiring is derived by the compiler, never declared: `physical` cables unless `mgmt_only`;
  `loopback`, `svi` and `subinterface` are configured, never cabled.
- An unimplemented `iftype` value is an explicit compile-time rejection.
- A link with a `mgmt_only` endpoint is omitted and recorded in the fidelity manifest; a
  second `mgmt_only` interface on a device is omitted and recorded.
- Twin reachability is containerlab's assigned address, never intent's.

---

<a id="d-004"></a>
## D-004 — Reserve only confident structure: keep `parent`, defer the VLAN generic

**Status:** Accepted 2026-09-14; superseded by [D-057](#d-057) on 2026-10-25.

### Context

SVIs need a VLAN referent. The first contract sketch added a `FylgjaVLAN {id, name}` stub.
Both `parent` and a VLAN generic serve features that are not implemented.

### Options considered

- **Stub VLAN generic now.** Relationships are expensive to retrofit, but proper L2
  modelling (membership, trunks, native and allowed lists) is exactly the deferred part,
  and a placeholder entrenches a design that must be redone.
- **`vlan_id` integer on Interface.** Migrating an attribute into a relationship later is
  the expensive direction.
- **Defer; reserve the `svi` value.** Adding an optional relationship later is additive,
  not breaking.

### Decision

Defer the VLAN generic; `svi` rejects at compile time. `parent` stays: self-referential, no
new kind, and its shape will not change when subinterfaces, LAG members or breakouts are
modelled.

### Consequences

- **Rule:** reserve structure you are confident about; never structure you are guessing
  at.
- The contract is four generics, `FylgjaPlatform`, `FylgjaDevice`, `FylgjaInterface` and
  `FylgjaLink`, beside the `FylgjaContract` node that carries its version.

---

<a id="d-005"></a>
## D-005 — Go end to end; no second runtime

**Status:** Accepted 2026-09-14; superseded by [D-065](#d-065) on 2026-10-29 in its scope: it binds the core, not the orchestrator.

### Context

The network-automation ecosystem is Python. Temporal has first-class SDKs for both. The
question was Go plus a Python activity tier, or Go only.

### Options considered

- **Go + Python tier.** Clean seam via Temporal; costs two toolchains, two release
  pipelines with SDK version skew, a runtime and venv on every privileged lab host, and
  shared types in two places.
- **All-Go.** One static binary per host; `genqlient` gives compile-time safety against
  schema drift, the design's largest correctness risk; single-SDK releases; Go's gNMI
  tooling is the reference implementation. Weaker for CLI scraping.
- **All-Go with a scoped Python queue for production collection.** Justified only by
  brownfield collection, which is out of scope.

Decisive: everything that boots in containerlab is modern and speaks gNMI or NETCONF, and
configuration push in a twin is largely file-writing.

### Decision

Go only. Where a dependency choice affects the static binary, the pure-Go option is chosen.
CLI parsing, when a platform requires it, is a Go TextFSM implementation over vendored
templates.

### Consequences

- One toolchain, one test framework, one artifact; the compiler checks a refactor end to
  end.
- Honest counter-case, recorded: if most contributors write Python and few write Go,
  adoption beats elegance. Any operator-extensible surface (assertions, later) is
  declarative so the hedge exists.

---

<a id="d-006"></a>
## D-006 — Temporal, self-hosted, for durable provisioning

**Status:** Accepted 2026-09-14

### Context

Provisioning is long, multi-step, failure-prone, and holds real resources. A run that dies
partway must resume; what it created must be reclaimed; operators must be able to see where
it is.

### Options considered

- **Ad-hoc job state**: a state machine in a database plus workers. Retries, heartbeats,
  cancellation and history all become bespoke.
- **Prefect**: Infrahub ships it, so one fewer service. Weaker durability and cancellation
  semantics.
- **Argo / Kubernetes-native**: presumes Kubernetes; containerlab's privileged Docker sits
  awkwardly in it.
- **Temporal**: durable execution, heartbeats, disconnected-context cleanup, schedules, a
  first-class Go SDK and test suite. Another stateful service to run.
- **Temporal Cloud**: removes the service; puts an external dependency in the path.

### Decision

Temporal, self-hosted. `temporal server start-dev`, file-backed, until the host move.

### Consequences

- Provisioning is a sequence of idempotent activities in a short workflow (D-007).
- Cleanup runs on a disconnected context.
- Budgets and timeouts are per platform from the PSP, never global.
- Tier-1 tests cover workflow logic with the SDK test suite and no server (D-017).

---

<a id="d-007"></a>
## D-007 — Short workflows with fixed ids; no lifecycle workflow

**Status:** Accepted 2026-09-14; superseded by [D-062](#d-062) on 2026-10-28 in its `Reconcile`,
its Schedule `fylgja-follow` and its count of workflows and activities.

### Context

With one twin at a time (D-013) and its state in containerlab and files (D-014), what does
Temporal actually run? An earlier design had one long-running lifecycle workflow per twin
that provisioned, then parked on a TTL timer and signals (`destroy`, `extend_ttl`,
`intent_changed`), answered `status` queries, and used continue-as-new to bound history.

### Options considered

- **Long-running lifecycle workflow.** The workflow *is* the state machine; TTL as a
  durable timer; destroy as a signal; Temporal enforces the singleton through the workflow
  ID. Costs: signal and query handlers, continue-as-new, `GetVersion` discipline for code
  that runs for weeks, and a second source of truth beside containerlab. Every one of those
  serves a parked phase that, for a twin that does not follow its branch, does nothing but
  wait for `destroy`.
- **Provisioning-only workflow plus an external scheduler** for TTL and reconcile.
  Reinvents timers outside Temporal.
- **Short workflows plus a Temporal Schedule for reconcile.** `provision` runs to ready and
  cleans up on failure; `destroy` removes the lab. Following (D-026) is a Schedule that
  starts a `reconcile` workflow on an interval: Temporal's own cron, needing no parked
  workflow. No TTL: a forgotten twin costs the lab host's memory, not pool capacity, and
  the operator destroys it.

### Decision

Short workflows, each started by one command and ending when its work is done:
`fylgja-provision` and `fylgja-destroy` with fixed ids; `fylgja-step` with its fixed id
(D-036); and `Reconcile`, started by the Schedule `fylgja-follow` as
`fylgja-reconcile-<time>`. No TTL, no signals, no queries, no continue-as-new, and one task
queue (D-015).

### Consequences

- The fixed ids reject a concurrent second `provision`, `destroy` or step; the host check
  refuses when the lab already exists (D-013). A create, a provision or a step beside a
  run in flight is refused as `run.in_flight`, and `twin destroy` cancels a provisioning
  run or a step and waits for its cleanup or its record.
- `fylgja twin show` reads containerlab and the twin directory; while a run is in flight
  it describes the run through Temporal's API.
- Cleanup on a create's failure or cancellation runs `DestroyLab` and `UnstageTwin` on a
  disconnected context. There is no sweeper: an orphan lab is detected at `create` and at
  worker start and cleared by `destroy`.
- `fylgja-step` takes a running waypoint twin to another waypoint of its series in eight
  phases (`inspect`, `host check`, `plan reconcile`, `stage`, `reconcile`, `readiness`,
  `push`, `record`) and then waits for the twin to settle: `VerifyTwin` reads the twin
  every second until every assertion read from a node holds or the budget is spent (120s
  by default, or `twin step --wait`, with a 30s margin for its timeouts). The wait runs on
  the run's own context with `WaitForCancellation`, so a `twin destroy` during it ends it
  at the next heartbeat or read and has it recorded `cancelled`; a run cancelled before its
  wait begins skips it. The wait never changes the step's outcome, exit, status, phase or
  timings (D-039).
- The worker registers the four workflows and twenty activities.
- If a long-lived phase is ever needed (opt-in branch coupling, incremental reconcile), the
  lifecycle design described above is the starting point.

---

<a id="d-008"></a>
## D-008 — Drive the `clab` CLI; do not link containerlab

**Status:** Accepted 2026-09-14

### Context

Fylgja deploys, inspects and destroys containerlab labs, and must do so in a way that
survives containerlab's releases and that an operator can repeat by hand.

### Options considered

- **Link containerlab as a Go library.** Tighter control; its internal packages are not a
  stability-guaranteed API and churn between releases.
- **Shell out to `clab`** with structured output. Stable; identical to what an operator
  runs by hand.

### Decision

Shell out.

### Consequences

- A broken lab is reproducible by hand from what Fylgja logged.
- containerlab upgrades do not break Fylgja's build.
- containerlab is also Fylgja's source of truth for whether a twin exists (D-014).

---

<a id="d-009"></a>
## D-009 — Fylgja never renders configuration; bootstrap from the PSP

**Status:** Accepted 2026-09-14

### Context

A twin fed by a Fylgja-specific renderer tests a different artifact than the one
production receives. But a node cannot boot into a topology with no configuration.

### Options considered

- **Fylgja renders** with its own templates. A second renderer; twin and production diverge
  by construction.
- **Consume an external pipeline's output** (Ansible/Nornir in git). Same artifact as
  production; needs a second ingest path and a branch-to-run mapping.
- **Consume Infrahub artifacts.** Transformations render per-branch artifacts inside the
  NSoT: same templates, same data, same output, already branch-scoped.

### Decision

Infrahub artifacts for everything that is configuration (D-028). The one thing Fylgja emits
per node is **bootstrap** (hostname, cabled ports enabled, discovery on), templated from
the PSP, never from intent, carrying no addresses. Bootstrap is platform plumbing so that a
node can be reached and wired; it is not rendering.

### Consequences

- Twin and production converge on identical text once artifacts are pushed.
- Rendering lives in Infrahub, so the external-pipeline path is not needed (D-028).
- **The artifact is pushed unchanged to a node whose names may differ.** "Fylgja never
  renders" also means Fylgja never *rewrites*: production's artifact names production's
  interfaces, and the twin may have mapped one of them to a port the node calls something
  else, or omitted it for a lossy mapping. Translating the text would make the twin test a
  document production never receives, which is what this entry rejects. So the line is
  pushed as production wrote it, and the honest surface is a **warning** at compile time,
  `artifact.interface.unrepresented`, naming each such line and what the node calls it. It
  never refuses; the operator decides what an unrepresented interface means for their
  test. It can only speak about interfaces intent also carries.
- **Bootstrap is the one thing Fylgja emits, and on one platform it arrives through the
  push.** The `ceos` kind takes a startup-config as a node's whole configuration, so a
  partial bootstrap sent that way would be a replace; instead its lines go ahead of the
  artifact's in the same push request (`config.bootstrap_via: push`). What bootstrap
  contains is the same either way, and the artifact is still the node's configuration.

---

<a id="d-010"></a>
## D-010 — Keep the CTM, thin

**Status:** Accepted 2026-09-14

### Context

With generics as the contract and a reference schema shipped, the Canonical Topology Model
is nearly a copy of the `genqlient`-generated types. Spec Kit's simplicity gate would
reject it.

### Options considered

- **Keep it, thin, blessed in the constitution.** Justified by synthesized nodes: boundary
  stubs and traffic generators are not in Infrahub and cannot live in types generated from
  Infrahub.
- **Drop it; compile from generated types.** Generated types are shaped by query structure
  (edges, nullable pointers), and synthesized nodes have nowhere to go at the moment they
  are needed.
- **Put synthesized nodes in Infrahub** on the twin's branch, flagged. Flips Fylgja from
  read-only to writing into the NSoT and pollutes the branch diff.

### Decision

Keep it, thin.

### Consequences

- A near-copy of the generated types until a synthesized node needs more.
- Provenance (`intent` / `observed` / `synthesized`) is carried from the start.
- The constitution blesses the layer so it is not re-litigated per plan.

---

<a id="d-011"></a>
## D-011 — The compiler is pure; the bundle is deterministic, content-addressed and deployable as emitted

**Status:** Accepted 2026-09-14

### Context

The compiler holds the highest-risk logic in the system, interface mapping. It must be
testable without a lab and inspectable by a human. An earlier design added a separate
**bind** stage that joined a twin-agnostic bundle with per-twin values (lab name,
management network, subnet, host, image overrides) so that one bundle could back many
twins.

### Options considered

- **Compiler does I/O.** Every test needs infrastructure; a failure could be in any of
  three concerns.
- **Pure compiler + bind stage.** The bundle stays hashable; a binding file carries the
  deployment values. With one twin (D-013) the binding is a constant lab name, a default
  network and image overrides: nothing that varies per deployment.
- **Pure compiler, deployable output.** The compiler writes the fixed lab name and relies
  on containerlab's default management network; a locally imported image is supplied
  through the PSP override directory that already exists. The bundle is host-agnostic
  because nothing host-specific is in it.

### Decision

Pure compiler, deployable output. Purity is enforced by a test. `bundle_id` is the hash of
the bundle's canonical bytes (D-024).

### Consequences

- Golden-file testable with zero infrastructure; `fylgja twin compile` runs the same
  function as the tests.
- The same intent compiles to the same `bundle_id` anywhere; a changed `bundle_id` is the
  drift signal for `reconcile` (D-026).
- No bind stage, no binding file. If deployment values ever need to vary (M9), the bind
  design described above is the starting point.
- Synthesized nodes are added to the CTM before compilation, never looked up by the
  compiler.

---

<a id="d-012"></a>
## D-012 — The intent reference is `(branch, at?)`; an unpinned reference is recorded, not pinned, and its twin follows the branch

**Status:** Accepted 2026-09-14; superseded by [D-062](#d-062) on 2026-10-28 in the part that says
an unpinned twin follows its branch.

### Context

The operator names a branch and, optionally, a point in time. Infrahub's `at` filters
`created_at <= at`, and pinning a read to a second-granularity local "now" silently hides
objects written in that second ([verified-facts.md](verified-facts.md)). Fylgja must decide
what an unpinned reference means over the life of a twin.

### Options considered

For an absent `at`:

- **Pin at creation to a local "now".** Reproducible on paper; the behaviour above makes it
  lossy, and clock skew between Fylgja and Infrahub makes it worse.
- **Send no `at`, record the read time, freeze.** Honest and simple; nothing is ever live
  with respect to intent.
- **Send no `at`, record the read time, follow.** The twin tracks its branch: `reconcile`
  re-reads, re-compiles, and rebuilds if the `bundle_id` changed.

For a present `at`: pass it verbatim on every query. Nothing else keeps the twin
reproducible.

### Decision

Pinned-verbatim for a present `at`. For an absent `at`: no `at` is sent; `observed_at` is
recorded as informational provenance, outside the bundle (D-023); and the twin **follows**
its branch through a scheduled `reconcile` (D-007, D-026), unless it is created with
`--no-follow`, which freezes it (D-027). The provenance block records `(branch, at?,
observed_at, schema_hash, contract_version)`; the bundle manifest carries all of it except
`observed_at` (D-023).

A **waypoint** is the third way to name a reference (D-032). Beside `--branch` and
`--branch --at`, `--waypoint <series>/<sequence>` names a waypoint the operator wrote in
Infrahub. The CLI's request resolves it before any run starts, to `(branch, at)`: the
waypoint's `as_of` as written, or, when none is written, the `branch.updated_at` Infrahub
stamped on the waypoint. From there the run is a pinned run: `at` verbatim on every query,
and a twin that never follows. The name is not part of the reference and enters no bundle,
so the same intent hashes the same however it was named.

### Consequences

- A twin is pinned (`--at` or a waypoint: reproducible, never rebuilt), following (an
  unpinned reference: rebuilt on change) or frozen (an unpinned reference with
  `--no-follow`), or was provisioned from a bundle. `fylgja twin show` says which.
- Reproducing an unpinned reference with `--at <observed_at>` is best-effort; the manifest
  says so.
- A following twin's reconcile is a rebuild; a waypoint twin is moved in place by a step
  (D-036). Incremental apply for a following twin is a limitation, stated as one.
- A pinned `at` carries at most six fractional digits, and more is refused
  (`intent.at.precision`). The `at` edge cases that remain are an
  [open question](#open-questions).
- Change detection for a following twin is a scheduled read and compare (D-026).

---

<a id="d-013"></a>
## D-013 — One twin at a time

**Status:** Accepted 2026-09-14

### Context

An earlier design required concurrent walking twins on one host, from the same reference
or different ones, and derived a per-twin isolation model: unique lab names, a Docker
network per twin with a subnet allocated from a per-host pool and released on destroy,
capacity accounting across twins, human names, and end-to-end cases for concurrency.

### Options considered

- **Concurrent twins with per-twin isolation.** Parallel experiments, live A/B, per-operator
  sandboxes. Costs: allocation with release and reconciliation, disambiguation in every
  command, a sweeper reasoning about many labs, and an identity model to keep twins apart
  from their references. On the one lab host (D-020), which also runs Infrahub and
  Temporal and is bound by its memory, a second twin of any realistic size is rarely
  affordable, so most of the machinery would be exercised only by tests.
- **One twin per lab host.** Only meaningful after the host move; adoptable then without
  touching the single-host design.
- **One twin at a time.** Fixed lab name, default network, no allocation, no names, no
  identity beyond "the twin". A/B is sequential; compile-time comparison between
  references is unaffected because bundles are content-addressed (D-011).

### Decision

One twin at a time.

### Consequences

- `fylgja twin create` is refused while lab `fylgja` exists, with instructions; the
  operator destroys first. Fixed workflow ids (D-007) reject concurrent runs.
- The lab name `fylgja` is written into every bundle by the compiler (D-011).
- An orphan lab, present on the host without a completed provision, is detected at
  `create` and at worker start, named with what is there, and cleared by `destroy`.
- The capacity check sums PSP memory budgets against the host budget and refuses with a
  reason.
- Runtime comparison between two live twins is out of scope and listed as a limitation.
  Lifting the rule means re-adopting the isolation model described above.

---

<a id="d-014"></a>
## D-014 — State lives in containerlab, files and Temporal; no database

**Status:** Accepted 2026-09-14; superseded by [D-062](#d-062) on 2026-10-28 in "Temporal holds …
the Schedule `fylgja-follow`".

### Context

Fylgja must answer: does a twin exist, what is it built from, what bundles have been
compiled, and what happened. An earlier plan used SQLite (pure-Go driver, WAL, `STRICT`
tables, dual-dialect migrations) for run history, migrating to Postgres at the host move;
an earlier design kept it as a twin registry beside Temporal.

### Options considered

- **SQLite registry now, Postgres later.** A relational record of twins, references,
  bundles and history. With one twin the registry holds one row that must be kept
  consistent with containerlab and Temporal, plus a schema, a driver and migrations.
- **Temporal as the only record.** Run history and search attributes answer "what
  happened". Retention is bounded; the current twin's provenance would live only in a
  closed workflow's result.
- **containerlab plus files.** Whether the twin exists and what it contains is what `clab
  inspect` reports. Its provenance is the manifest staged beside it in a twin directory.
  Every compiled bundle sits in a bundle store by hash. Temporal keeps run history. Nothing
  to keep in sync, nothing to migrate.

### Decision

containerlab plus files, with Temporal keeping what is Temporal's. No database.

### Consequences

- **State lives in three places.** containerlab says whether the twin exists. The files
  under the state root (`FYLGJA_STATE_ROOT`) hold the twin directory, with its record
  `twin.json` and its staged bundle, and the bundle store. Temporal holds run history and
  the Schedule `fylgja-follow`, the only record of the branch a following twin follows and
  its interval (D-026): a following create makes it, and `twin destroy` or a failed
  rebuild deletes it. Constitution IX says "Temporal holds run history", which understates
  the Schedule but breaks none of its rules.
- `fylgja twin show` is `clab inspect`, `twin.json` and Temporal's runs and Schedule.
- **Whether the twin exists is containerlab's to say.** A record without its lab is no
  twin: `twin show` reports kind `none` beside the record, `twin step` refuses it as
  `step.twin.unsteppable`, and `twin verify` as `verify.twin.absent`.
- **The record, `twin.json` version 5,** carries the provenance; the waypoint the twin was
  built from (`{series, sequence, description, at_source}`), or an explicit `null` for a
  twin named by branch; `state` (`ready` or `diverged`); the last step (its sides, run,
  outcome, phase, containerlab's plan with each package's declaration beside it, each
  node's push outcome, the timings, and `wait`, how the step's wait ended, `null` until
  `VerifyTwin` writes it), or `null` for a twin that has not stepped; and on each node
  `holds`, the `bundle_id` last known to run on it, or `null` for a node containerlab
  restarted that the step never pushed. A `failed` push may have committed, and only
  `refused` is the node's own atomic answer, so a reader of `holds` reads the step's push
  outcome beside it. An earlier version of the record still reads, as a twin without the
  later fields.
- The waypoint stays Infrahub's data, read each time it is used; nothing but the twin's own
  record holds it (D-032).
- A step stages its target into the twin directory before it touches the lab. After a
  step that landed, the manifest beside the lab is the provenance of what runs; after a
  diverged one, the top level of the record stays where the twin came from, the step
  names where it was going, the staged manifest is the target's, and each node's `holds`
  says which bundle it runs (D-036).
- `twin verify` writes nothing. It reads the record's `nodes[].holds` as the record's claim
  and the staged manifest as what to assert (D-038).
- **The API's server keeps nothing between two requests** (D-040, D-041). What it answers
  comes from containerlab, the files under the state root, the nodes, Infrahub and the
  workflow service, read for each request; it holds no registry, session or cache. One
  thing lives as long as a request and no longer: while a run's request is open, the
  server holds the channel that request's interrupts arrive on, keyed by an identity the
  client chose, so that the operator's interrupt, a second request, reaches the handler
  that holds the run. No answer is made from it, and it is not a lock (D-007, D-013). A
  restarted server has lost nothing a run needs.
- The bundle store (`<state root>/bundles/<bundle_id>/`) sits behind a small interface so
  that object storage can replace the directory at M8; the twin directory stays local to
  the host.
- History older than Temporal's retention exists as bundles on disk. A question only a
  database can answer is the trigger to add one.
- The constitution's storage principle is "files only; no database".

---

<a id="d-015"></a>
## D-015 — One task queue until the host move

**Status:** Accepted 2026-09-14

### Context

Read and compile can run anywhere; deploy, readiness and destroy must run where the
containers are. Temporal routes by task queue. An earlier design used a control queue plus
one queue per lab host from the start, so that the single-host assumption would never
become load-bearing.

### Options considered

- **Control queue plus per-host queues now.** Correct for many hosts; on one host it is two
  names served by one worker, a scheduler stub, and a concept in every test.
- **Temporal Sessions.** Idiomatic host affinity; harder to inspect.
- **One queue, with the host-bound boundary kept in code.** Everything that must run on the
  lab host lives in `internal/lab`; dispatching that package to a per-host queue at M9 is
  a change to worker registration and one option on each activity call, not a refactor.

### Decision

One queue, `fylgja`, with `internal/lab` as the host-bound boundary.

### Consequences

- One worker process, one registration site.
- Per-host queues and a real scheduler arrive with M9, where they are needed.

---

<a id="d-016"></a>
## D-016 — Multi-vendor by requirement; platforms are data

**Status:** Accepted 2026-09-14

### Context

A slice of a real network spans vendors. An early design treated the second platform as the
design centre.

### Options considered

- **Single-vendor first, generalize later.** Abstractions validated at N=1 are fiction;
  generalization happens under pressure.
- **"Supports several vendors", one per twin.** A real slice is mixed.
- **Multi-vendor by requirement, heterogeneous twins.** Nothing may assume homogeneity; the
  vendor abstraction is a named, declarative object.

### Decision

The third. A **Platform Support Package** (PSP) is the unit of platform support; a platform
is supported when its PSP passes the conformance suite. If a new platform needs a compiler
change, the abstraction is wrong and the change is rejected.

### Consequences

- The PSP format is designed against the full expected platform range (public container,
  account-gated container, restricted CLI-first binary, licensed VM).
- Budgets, timeouts, readiness, image acquisition and fidelity are PSP facets.
- Two packages ship, `nokia_srlinux` and `arista_eos`, and a mixed twin of both compiles,
  provisions and steps with no platform branch in the code (D-031).
- The constitution blesses the abstraction so it is not re-litigated per plan.

---

<a id="d-017"></a>
## D-017 — Test architecture

**Status:** Accepted 2026-09-14

### Context

Fylgja integrates with Infrahub, Temporal, containerlab and network operating systems, each
slow or heavy to run. Its tests must be cheap enough to run on every change, and real where
the integrations are most likely to surprise.

### Options considered

- **Pure core, tiered by cost.** Compiler golden-tested; workflow logic with the Temporal
  test suite; contract tests against a real Infrahub; end-to-end with a booted NOS on
  demand.
- **Everything real, always.** A loop of ten minutes and more on every change, which a
  single lab host cannot run often; people stop running it.
- **Mock-heavy.** Fakes encode assumptions about the two integrations most likely to
  surprise.
- **Recorded interactions.** Cassettes go stale as the schema evolves.

### Decision

Pure core, tiered by cost. Temporal's `testsuite` is tier 1: it runs workflow code without
a server and is where retry, cancellation and cleanup paths are tested.

### Consequences

- Compiler purity is enforced by a test.
- Contract tests use a real Infrahub, never a fake.
- End-to-end never gates a change. Its cases include a refused second `create`, an orphan
  lab cleared, destroy-then-create across references, a following twin rebuilt, a mixed
  twin and a step there and back.
- Recorded workflow histories, replayed in tier 1, hold what the test suite cannot show:
  how a run behaves when the service cancels it at a given point.
- The populated fixture branch the contract tier needs is also the demonstration topology.

---

<a id="d-018"></a>
## D-018 — Single module, single binary in three roles, PSPs and schema embedded

**Status:** Accepted 2026-09-14; superseded by [D-065](#d-065) on 2026-10-29 in its scope: it binds the core, not the orchestrator.

### Context

Fylgja has a command-line client, a worker that runs on the lab host, and platform
packages and a schema that both read. How they are versioned and shipped decides what can
skew.

### Options considered

- **Single module, two binaries** (`fylgja`, `fylgja-worker`). Two artifacts to version and
  ship.
- **Single module, one binary** with `fylgja worker run`. One artifact per host; the CLI
  and the worker cannot skew.
- **Core module plus a PSP repository.** The destination *if* outside contributors ship
  PSPs; premature before one thing works.

### Decision

One module, one binary. PSPs and schema embedded, with a runtime override directory.

### Consequences

- The conformance suite is `go test`.
- **One binary, three roles**: the client's commands, `fylgja worker run` and `fylgja
  serve` (D-040, D-041). A lab host needs the binary and nothing else, and runs two of its
  roles as two processes beside the workflow service.
- **The worker and the server cannot skew**: they start from one binary on one host. A
  client on another machine can differ from them, so the API carries its version in every
  path and its build in every answer; a client refuses a version it does not speak, and
  says so in one line when the builds differ.
- **The client's commands live in a package of their own**, which a test over imports
  holds to the API's client, the document types and a file helper. The binary links the
  core for its two roles, and no command of the client's can reach it.
- The override directory (`--psp-dir`) is the server's and the worker's; a client has
  none. It is also how a locally imported image reaches the compiler (D-011).

---

<a id="d-019"></a>
## D-019 — Noun-verb CLI; `twin create` is the front door; stages exposed

**Status:** Accepted 2026-09-14

### Context

Fylgja is driven from the command line. The CLI serves operators ("give me a twin") and
developers ("show me one stage's output"). There is one twin.

### Options considered

- **`twin build`** as the front door. Reads as producing an artifact; sits oddly beside
  `destroy`.
- **`twin create` / `show` / `destroy`** with stage commands `intent read`, `twin compile`,
  `twin provision`, and `schema check` for conformance alone. No twin arguments: there is
  one. No `resolve` command: conformance is `schema check`, completeness is part of
  `intent read`.
- **Verb-first** (`fylgja build`). Messy once nouns multiply.
- **Stages only.** Wrong for operators, who want a twin, not four steps.

### Decision

The second.

### Consequences

- Stage commands share code with their tests; divergence is a defect. Each stage runs in
  the API's server (D-040): the stage a command asks the API for is the code its tests
  call, and the command's own tests reach it through a server in the test's process.
- `--json` on every read command; `--dry-run` on `create`, `provision` and `step`.
- Later commands keep the shape: `twin step` and `twin verify` (D-036, D-039), `waypoint
  list` and `waypoint plan` (D-032), `psp validate`, and the two roles, `worker run` and
  `serve` (D-018).
- Flags arrive with the feature that needs them. Following's opt-out is `--no-follow`
  (D-027); `--psp-dir` is a flag of `serve` and `worker run` alone.
- The stage commands need a server, and neither Infrahub nor the workflow service when
  they read none: `psp validate` and `twin compile` run against the server's packages.
- A client reaches a lab host by the API's address, `FYLGJA_API_ADDRESS`. Whether M9 needs
  a `--host` flag beside it is M9's to say.

---

<a id="d-020"></a>
## D-020 — One lab host, the development host, until the host move

**Status:** Accepted 2026-09-14

### Context

A twin runs on a lab host: Docker, containerlab and the node images, beside the worker.
Where that host is decides what a twin can cost and which platforms can boot.

### Options considered

The development host (no infrastructure beyond it; bound by its memory); cloud VMs
(elastic, a cost per twin-hour, needs a scheduler); on-prem hosts (predictable, must be
provisioned first); Kubernetes via clabernetes (presumes a cluster).

### Decision

One lab host, the development host, until M9.

### Consequences

- Scope is a hard constraint: the host's memory caps twin size.
- Native-container platforms run on it. A vrnetlab platform needs KVM on the lab host and
  arrives with the host move.
- One task queue is sufficient (D-015).
- The worker and the API's server run on it, beside the workflow service (D-041).

---

<a id="d-021"></a>
## D-021 — SR Linux first, EOS second; later order parked

**Status:** Accepted 2026-09-14

### Context

Fylgja needs a first platform to prove the mechanics on, and a second to prove that the
platform abstraction (D-016) is real.

### Options considered

- **Production platform first (EOS).** Valuable from day one; cEOS is account-gated and
  needs an import path, and the harder interface mapping arrives before the loop is
  proven.
- **SR Linux first.** Public image, containerlab's home platform, YANG-modelled throughout;
  zero vendor friction. M1 and M2 then test nothing about the actual network.

### Decision

SR Linux for M1 and M2; EOS at M7, because an abstraction validated at N=1 is fiction. The
order after that is [parked](#parked).

### Consequences

- M1 and M2 are a mechanics proof, not a delivery.
- SR Linux's one-to-one mapping must not shape the mapping format. **The format was
  designed against the lossy case first**: PSP `0.4`'s mapping profile was shaped on a
  test-only package with EOS-style naming, `testdata/psp/lossy/chassisos.yaml`, which
  expresses many-to-one, collapsing breakout and a node name unlike the production name,
  and compiles beside SR Linux in one heterogeneous golden, `testdata/golden/lossy/`. Its
  names come from public documentation, verified against no image; it is never shipped or
  deployed. SR Linux's mapping is expressed in the same format.
- **EOS is not the lossy case.** cEOS carries production's names one to one: the chassis
  reports `Ethernet1` … `Ethernet511`, and modular and breakout names such as
  `Ethernet2/1` and `Ethernet3/1/1` are the node's own. So `psp/arista_eos.yaml` declares
  no lossy rule, `chassisos` alone proves the lossy paths, and the first test of the
  format against a lossy platform that boots is still ahead. EOS was worth taking second
  anyway: a second vendor, CLI syntax, push mechanism and probe found what a one-to-one
  twin of one vendor could not (D-029, D-031).
- **Bootstrap is line-at-a-time.** The template repeats a line holding `{interface}` or
  `{port}` once per cabled port and emits every other line once. SR Linux's `set /
  interface {interface} admin-state enable` fits one line; a modal CLI's block does not
  (`chassisos`'s `interface {interface}` then `no shutdown` enables only the last cabled
  port, and the lossy golden keeps those lines, since a one-line form invented for a
  platform that never boots would prove nothing). cEOS needs no per-port enable: the
  `ceos` kind's own default configuration sets the host name from the node name and brings
  every port up. The limit stands for the next platform to meet, as a change to the
  bootstrap format, never a compiler branch for the platform (D-016).
- **Two things only a booted EOS showed.** cEOS reads `!` as a comment and refuses `#` at
  token 0, so the generated header's marker is package data (`config.comment_prefix`);
  and the `ceos` kind takes a startup-config as the node's **whole** configuration, so the
  bootstrap reaches the node through the push (`config.bootstrap_via: push`, D-009). Both
  are format values keying code that already existed (D-031).

---

<a id="d-022"></a>
## D-022 — Every bundle carries a fidelity manifest

**Status:** Accepted 2026-09-14

### Context

Nothing verifies that the twin resembles production; fidelity is asserted, never measured.
A twin that quietly dropped three devices is worse than no twin.

### Options considered

- **Disclose in documentation.** Nobody reads it at the moment it matters.
- **Disclose per build, in the artifact.** The manifest travels with the bundle and is
  printed by `twin show`; it is the only trust signal Fylgja has about production.

### Decision

Per build, in the bundle manifest: modelled exactly, approximated, stubbed, omitted;
software- versus hardware-forwarding per platform (D-025); every omission with its reason;
every synthesized node.

### Consequences

- Skips and omissions are recorded, never silent.
- Each mapping row carries intent's `enabled`, so the manifest says what intent says of
  each port. A port intent disables is still cabled and its bootstrap still enables it;
  intent's disabling reaches the node through the artifact alone.
- **The compiler's omissions are verify's context, never its findings** (D-038). `twin
  verify` asserts what the manifest says is cabled: each node's host name, each cabled port
  intent enables, and each link from both ends. It asserts nothing of an uncabled port, an
  omitted interface or a lossy mapping's production name, and a port intent disables is
  skipped with its link, each skip recorded with its reason. It measures the twin against
  intent; fidelity to production is still asserted by this manifest and never measured
  (Constitution X).

---

<a id="d-023"></a>
## D-023 — `observed_at` is recorded outside the bundle

**Status:** Accepted 2026-09-14

### Context

D-012 records `observed_at`, the time an unpinned read was made. D-011 makes `bundle_id`
the hash of the bundle's canonical bytes and uses a changed `bundle_id` as the drift signal
for `reconcile`. If `observed_at` were in the bundle manifest, both could not hold: two
reads of unchanged intent a minute apart would produce different manifests, different
hashes, and a spurious rebuild on every check.

### Options considered

- **Keep `observed_at` in the manifest; exclude it from the hash.** Define the canonical
  form as the bundle with `provenance.observed_at` stripped. "Canonical bytes" then carries
  a hidden exception that every hasher and every reader of the manifest has to know about.
- **Record `observed_at` outside the bundle.** It is a fact about the read, not about the
  intent. `read` writes it into the CTM file's envelope; `provision` writes it into
  `twin.json`. The manifest's provenance block keeps branch, `at`, schema hash and contract
  version, all properties of the intent. Compile depends only on intent in the strong
  sense.

### Decision

Outside the bundle.

### Consequences

- The manifest's provenance block is `(branch, at?, schema_hash, contract_version)`. The
  CTM envelope and `twin.json` carry `observed_at` as well.
- The same intent compiles to the same `bundle_id` regardless of when it was read;
  `reconcile` compares hashes without special cases.
- The best-effort reproducibility statement for an unpinned reference is made by
  `twin.json` and `fylgja twin show`, not by the manifest.
- Constitution principle VI says so.

---

<a id="d-024"></a>
## D-024 — `bundle_id` is the canonical-bytes hash; the bundle carries no build stamp

**Status:** Accepted 2026-09-14

### Context

D-011 says the bundle is deterministic and content-addressed, but not how the address is
computed. `bundle_id` is load-bearing: `provision` verifies a staged bundle before
deploying it, and `reconcile` compares hashes to decide whether to rebuild. Both need a
rule they can cite, and one an operator can check without running Fylgja.

The same requirement settles a second question. An earlier manifest carried a top-level
`fylgja_version` stamp and a `provenance.source` field. A build stamp in the bundle means
every Fylgja release changes `bundle_id` for unchanged intent: the defect D-023 found in
`observed_at`, arriving from the build rather than from the clock.

### Options considered

- **Hash a canonical serialization of the manifest.** Small and fast, but it makes
  `bundle_id` an identity for the manifest rather than for the bundle: a corrupted or
  edited configuration file would keep its id. What `provision` has to verify is the bytes
  it is about to deploy.
- **Hash an archive of the directory (tar, zip).** Ties the id to an archive format's own
  choices (ordering, padding, recorded mtimes and permissions), so the rule becomes
  "whatever this library does this year".
- **Hash a file listing, in the `sha256sum` format.** The id is the SHA-256 of the lines
  `<sha256 of the file>  <relative path>` over every regular file, in byte-wise order of
  path. Fully specified in prose, reproducible from the directory with standard tools, and
  deliberately blind to everything that is not a regular file's path and bytes.
- **Keep the build stamp and exclude it from the hash.** The hidden-exception shape D-023
  already rejected.
- **Keep the build stamp and accept a new `bundle_id` per release.** Honest about
  provenance and fatal to `reconcile`, which would rebuild every twin on upgrade.

### Decision

`bundle_id` is the lowercase hex SHA-256 of the file listing in `sha256sum`'s output
format: one `<hex>  <path>\n` line per regular file, in byte-wise path order, paths
relative to the bundle root. It is recomputable with no Fylgja binary:

```sh
( cd <bundle> && find . -type f | sed 's|^\./||' | LC_ALL=C sort | xargs sha256sum | sha256sum | cut -d' ' -f1 )
```

The bundle carries no build stamp. The manifest's provenance block is `(branch, at?,
schema_hash, contract_version)`: every field a property of the intent, none of the build or
the clock.

### Consequences

- Only a regular file's path and bytes affect the id. Directories, permissions, ownership,
  mtimes and symlinks do not, so a bundle copied between machines keeps its id.
- `bundle_id` is never written into the bundle; it would change itself. `twin compile`
  prints it, and `--json` carries it as the document's top-level `bundle_id`.
- `internal/bundle` implements the rule twice on purpose, over the compiler's in-memory
  output and over a directory on disk, with a tier-1 test asserting the two agree.
  `twin provision` and a step verify a staged directory with the second.
- The binary stamps its own `--version`; that identifies the binary and reaches no bundle.
- Upgrading Fylgja does not change `bundle_id` for unchanged intent. A compiler change that
  *should* invalidate every bundle says so as a format or contract version change, never
  silently with a release.
- Principle V requires the bundle to be "content-addressed by the hash of its canonical
  bytes"; this entry says what canonical means, so the principle can be checked rather than
  merely agreed with.

---

<a id="d-025"></a>
## D-025 — `fidelity.production_forwarding` is asserted per platform, one entry per package

**Status:** Accepted 2026-09-14

### Context

D-022 says the fidelity manifest records software- versus hardware-forwarding **per
platform**. The first manifest had a single `fidelity.production_forwarding` value, which
the compiler assigned inside its per-device loop, so in a heterogeneous twin the last
device's platform would silently have won. With one platform supported that could not
happen, but the code claimed the manifest carried every platform's answer, and an assertion
about fidelity that is not true is what Constitution X exists to prevent.

### Options considered

- **Keep one value and let the last platform win.** Cheap and dishonest: a mixed-vendor
  twin would assert hardware forwarding for a software-forwarding node, the one kind of
  claim a fidelity manifest must never make.
- **Keep one value while one platform is supported, and refuse a second answer.** Honest,
  and what the single value did until a second platform arrived: a CTM whose devices
  disagreed was refused as `fidelity.forwarding.mixed`. It defers the format change to the
  moment it can be tested, and is wrong as soon as a mixed twin is meant to work.
- **Per platform, one entry per package, with a bundle-format version change.** Chosen once
  a second platform shipped.

### Decision

`fidelity.production_forwarding` is an object keyed by platform id, one entry per package
in the bundle, so a mixed twin asserts each vendor's own answer. The same format change
made the rest of the block unconditional: `mapping[].node_name` is on every row (`null`
where the row has no port), `fidelity.lossy` is always present (`[]` when empty), and every
node carries a `bootstrap` entry (`file`, `via`).

### Consequences

- `fidelity.forwarding.mixed` is retired with its tests: the refusal existed only because
  one field had to hold two answers ([D-031](#d-031)).
- `fidelity.approximations` is per platform too, and is deduplicated on the claim as
  emitted, platform included, so two platforms asserting the same approximation produce
  two entries.
- **Every rejection the compiler can raise is mirrored in `internal/validate`**, worded
  identically, so that a branch that reads clean will compile: a rule only the compiler
  knew would let `intent read` write a CTM that `twin compile` then refuses. Validating
  over the whole CTM also names every offending device, where the compiler's loop returns
  at the first.
- A bundle-format version change is made once and rebuilds a following twin once; that is
  why it is made rarely. `twin provision` refuses a bundle of any version but the one this
  build deploys (`4`), naming both.
- Principle X requires the manifest to say that fidelity is asserted, which
  `fidelity.basis` does; this entry makes the rest of the block as true as that field
  claims.

---

<a id="d-026"></a>
## D-026 — A following twin detects change by a scheduled read and compare, and nothing else

**Status:** Accepted 2026-09-16; superseded by [D-062](#d-062) on 2026-10-28.

### Context

D-007 and D-012 decide that an unpinned twin follows its branch, through a Temporal
Schedule that starts a short `reconcile` workflow. What decides that the branch changed?
Infrahub's event stream and polling its branch diff were the candidates; both were checked
against Infrahub 1.11.2 on 2026-09-16.

The compare needs nothing new. A check reads the branch without `at`, compiles it on the
code path `twin create` uses (D-011), and compares the new `bundle_id` with the one
`twin.json` records. `bundle_id` is the hash of the bundle's bytes (D-024), so it moves for
every change that reaches the bundle and for no other (D-023).

### Options considered

- **Infrahub's event log, as a trigger or as a filter before the read.** `InfrahubEvent` is
  a GraphQL query with `since`, `branches` and `level` filters. It is not a stream: the one
  `Subscription` field is a polled query, and push exists only as webhooks, which need an
  HTTP listener. Used as a filter ("any root event on branch B since the last check?"), it
  would save one read and compile per quiet interval: about a second with the three-node
  branch (0.77–1.4s measured). It cannot say whether a change reaches the bundle, since an
  attribute outside the generics fires `infrahub.node.updated` too, so the compare still
  has to run. It never sees a change that moves `bundle_id` without an event on the branch:
  a compiler upgrade or a support package's image. And it needs a `since` cursor kept
  between checks, a second piece of state beside the Schedule, which holds only the branch
  and the interval, and D-014 keeps Fylgja free of state like that.
- **Polling the branch diff.** `DiffTree`/`DiffTreeSummary` return `null` until a
  `DiffUpdate` mutation computes the diff (3.5s on a throwaway branch), and Fylgja does not
  write to Infrahub (Constitution III). The diff also answers a different question: it
  compares the branch with `main` from the branch's creation, so a branch whose data
  changed by three objects reported 1219 additions, because its schema load was in the
  diff. Asking what changed since the twin was built needs `from_time` set to the last
  check's time: the same cursor, and still no answer about the bundle.
- **The scheduled read and compare, alone.** One read and one compile per interval,
  detecting by construction everything the bundle can differ by: data, the schema hash
  (which moves for a schema load on the branch, never for data), the compiler, the support
  packages. It keeps no state beyond the Schedule.

### Decision

The third. Only the Schedule starts a check. A check reads, compiles and compares every
time. Nothing from the event log or the branch diff decides whether a check reads, or
starts one sooner than the interval.

### Consequences

- A change is seen at the next check, never sooner. The interval is the operator's
  (`--interval`, default 5 minutes, floor 10 seconds), held by the Schedule.
- A quiet interval costs one read and one compile. The compiled bundle is filed in the
  store as every compile is. An unchanged check touches nothing on the host.
- A change is seen as it stands when the check reads. A change made by several writes, or
  a write whose artifacts are not yet generated ([D-028](#d-028)), can be read half made:
  the check rebuilds to that intermediate bundle, and the next check rebuilds again.
- A compiler or support package upgrade on the worker rebuilds a following twin at its next
  check, when the upgrade moves `bundle_id`. That is intended: the twin follows what the
  product would build now.
- No cursor, memo or search attribute is kept anywhere.
- The API's server ([D-040](#d-040)) is the long-lived listener a webhook needs. A webhook
  that starts a check early is on the [roadmap](roadmap.md#backlog); it
  would add to the interval, never replace the compare.

---

<a id="d-027"></a>
## D-027 — `--no-follow`, not `--pin`, is the opt-out from following

**Status:** Accepted 2026-09-16; superseded by [D-062](#d-062) on 2026-10-28.

### Context

A twin built without `at` follows its branch (D-012). It needs an opt-out, a flag that
stops it from following. The glossary and D-012 use **pinned** for an intent reference that
carries an `at`, which is reproducible bit for bit (Constitution VI). A twin read at the
branch head and never checked again carries no `at` and is not reproducible: its
`observed_at` is a local clock reading, and reproducing it with `--at <observed_at>` is
best-effort.

### Options considered

- **`--pin`.** Short, and it reads as the opposite of following. It would give "pinned" two
  meanings: a reference with `at`, and a twin that was only frozen. `twin show` would have
  to tell a pinned-with-`at` twin from a pinned-without-`at` one, and an operator reading
  "pinned" could believe a frozen twin reproducible.
- **`--freeze`.** Names the result, but introduces a verb the rest of the surface does not
  use, and hides that the thing turned off is following.
- **`--no-follow`.** Named for what it does. "Pinned" keeps its one meaning, and the frozen
  twin gets its own term.

### Decision

`--no-follow`. With it, `twin create` builds from the branch head and never checks the twin
again. With `--at` a twin never follows, whether or not `--no-follow` is given, and
`--interval` with either is refused (`follow.flags.conflict`). A twin built from a
waypoint never follows either, and `--waypoint` beside `--interval` is refused as
`waypoint.flags.conflict` ([D-032](#d-032)).

### Consequences

- `twin show` reports the kind of twin: **pinned** (an `at` given, reproducible, never
  follows), **following** (no `at`, checked every interval), **frozen** (no `at`, not
  following: `--no-follow`, a Schedule lost or stopped, or a following of another branch),
  **from a bundle** (provisioned from a bundle, never follows), **diverged** (a step failed
  or was cancelled, [D-036](#d-036)), or **none**.
- The glossary defines both **pinned** and **frozen twin**.
- Principle VI reserves reproducibility for a present `at`, and this entry keeps the flag
  from blurring that.

---

<a id="d-028"></a>
## D-028 — Configuration is rendered in Infrahub and enters the bundle at read time

**Status:** Accepted 2026-09-18

### Context

D-009 decides that Fylgja never renders configuration and that everything a node runs
beyond bootstrap is an artifact Infrahub rendered for the branch under test. It left open
whether rendering happens in Infrahub at all; if not, an external-pipeline ingest path
would have had to come first. On 2026-09-18 the Infrahub 1.11.2 Fylgja was built against
rendered nothing, and nothing rendered anywhere else, so the question was not where
rendering lives but whether to build it in Infrahub. The operator did: one Jinja2
transformation over the Fylgja generics, one query and one artifact definition
(`device-config`, `text/plain`, one per device in a Standard Group).

What Infrahub 1.11.2 does with it ([verified-facts.md](verified-facts.md)):

- The target's kind must inherit `CoreArtifactTarget`. Content is served at `GET
  /api/storage/object/{storage_id}`, and `CoreArtifact` carries `checksum` (the MD5 of the
  bytes served) and `storage_id`. A listing pinned to `at` gives the storage id as of `at`,
  and the older object is still served.
- **It regenerates nothing on its own.** Not on a data change, a group join or a commit
  import. Generation is an explicit `POST /api/artifact/generate/<definition id>?branch=`,
  or a proposed change's checks. An intent change nobody generated for leaves the artifact
  `Ready` with its old bytes, and nothing marks it stale. A template commit reaches a
  branch without git sync only by a rebase, and then a generate.
- **It reads anonymously**: a request with no token is served, and a wrong token is 401
  everywhere.

A second question sits beside the first. D-026 detects change by comparing `bundle_id`, so
how a following twin sees a changed artifact depends on where the artifact enters the
pipeline.

### Options considered

- **Consume an external pipeline's output.** D-009's fallback. A second ingest path and a
  branch-to-run mapping, for a renderer this deployment does not have.
- **Render in Infrahub, fetch at push time.** A post-boot activity asks Infrahub for each
  node's artifact during provisioning. The bundle is unchanged, so an artifact-only change
  never moves `bundle_id`: a following twin runs stale configuration until something else
  changes, and a pinned twin's configuration is whatever Infrahub holds when it is
  provisioned, not at `at`. Constitution V's "deployable as emitted" would no longer
  describe the bundle.
- **Render in Infrahub, fetch at read time.** The read fetches each device's artifact with
  the intent, and the compiler places it in the bundle beside bootstrap. The bundle is the
  whole of what the twin will run; `bundle_id` moves for a configuration change as for a
  topology change; D-026 detects it with no new mechanism; a pinned read fetches the
  artifact as it stood at `at`, so reproducibility (Constitution VI) extends to
  configuration.

### Decision

The third. Rendering is built in Infrahub, and the external-pipeline path is rejected, not
deferred. Artifacts enter the pipeline at read time and travel in the bundle; delivery to
the node is post-boot, by the mechanism the platform's support package declares
([architecture §4.7](architecture.md#47-platform-support-packages)). Fylgja pushes the
artifact's bytes unchanged (Constitution IV).

### Consequences

- The bundle carries each node's artifact (`configs/<node>.<artifact_name>`) and the
  manifest its `artifact` entry, with the checksum. Configuration is among the things
  `bundle_id` covers, and among the things a following twin is rebuilt for.
- **Whoever changes a branch generates its artifacts**: the operator, or the fixture tool
  and the end-to-end script for their own branches. Fylgja never does, since generating is
  a write to Infrahub (Constitution III). A twin built from a branch nobody generated for
  runs the branch's last rendering, which is what Infrahub holds as current.
- The reference schema's device kind inherits `CoreArtifactTarget`, and the artifact
  definition needs a Standard Group of devices (`fylgja-devices`), which the fixture seeds.
  An organisation using its own kinds adds the inheritance to its device kind, as it adds
  the Fylgja generics (D-002); conformance names it when it is missing
  (`schema.artifact_target.missing`).
- The support package names the artifact that is a device's configuration (the Config
  facet: artifact name and accepted content type, beside the delivery mechanism), so a
  second platform brings its own name and no flag is added. One template renders both
  shipped platforms, branching on the device's platform, so every artifact is
  `device-config`.
- A read refuses, with a finding naming the device and the artifact, when a device has no
  artifact of that name or more than one, when Infrahub does not hold it `Ready`, when its
  content type is not one the package accepts, when its content is not served, or when its
  bytes do not hash to Infrahub's checksum. A refusal of a stale artifact reads Infrahub's
  own status, never a staleness Fylgja infers.
- The token Fylgja reads intent with must also read artifacts and storage objects. An
  Infrahub that reads anonymously cannot show a token that reads intent and not content;
  that refusal is proven in tier 1 from a faked 401 and 403.
- Where the template lives, and how Infrahub registers its repository, is [D-044](#d-044):
  in this repository, registered read-only with no credential. A private repository needs
  a `CorePasswordCredential` over HTTPS, since Infrahub 1.11.2's task workers carry no SSH
  client; that is Infrahub's concern, not Fylgja's.
- Revisit if an organisation renders outside Infrahub. That would reopen the first option
  as a feature of its own, not a change to this one.

---

<a id="d-029"></a>
## D-029 — A support package declares whether readiness waits for the push transport

**Status:** Accepted 2026-09-21

### Context

The push is how an EOS node gets its bootstrap at all: containerlab's `ceos` kind takes a
startup-config as the whole configuration, so a partial one would cost the node its login,
its eAPI and its gNMI, and the bootstrap therefore travels with the artifact in one eAPI
request after the node is ready ([D-009](#d-009)).

Readiness gates on the probe the package declares: for `arista_eos`, a gNMI Get on 6030.
Three mixed creates in a row failed at `push`, both EOS nodes answering `dial tcp
<ip>:443: connect: connection refused`, while SR Linux's JSON-RPC push to the same port
number succeeded every time. Measured by deploying the mixed topology by hand and polling
every 0.5s: gNMI opened at deploy+11.1s on both EOS nodes and eAPI at +12.1s and +13.6s.
So a create's push began about a second after the probe succeeded, and its three attempts
spanned 8.9s without reaching a listening eAPI, never spending the 30s `push_timeout_s` the
package budgets.

**A node was reported `ready` while the transport its configuration must arrive on was
shut.** That is a defect in what `ready` means, not in the push.

### Options considered

1. **Raise the push's retry budget.** One constant, no format change. Rejected: it
   lengthens every push failure's path, including a dead node's, and it leaves `ready`
   meaning something untrue. It also depends on the gap being smaller than whatever budget
   is picked.
2. **Retry inside the push activity until `push_timeout_s` is spent.** Uses a budget the
   package already declares. Rejected for the same two reasons.
3. **Gate readiness on the push transport for every package, in code.** Fixes the meaning
   of `ready` whatever the gap's size. Rejected because every platform pays for one
   platform's behaviour: SR Linux's push transport is up when its probe answers, and a
   mechanism that assumes otherwise is a branch on a platform wearing a mechanism's clothes
   (Constitution II).
4. **The package declares the wait.** Chosen.
5. **Do nothing and record it.** Rejected: a mixed twin could not be created.

### Decision

`readiness.await_push_transport`, a boolean, absent meaning false. When a package sets it,
the node is ready only once its probe has answered **and** the endpoint its `config.push`
names accepts a TCP connection, both inside the one `readiness.timeout_s` the package
budgets. `arista_eos` sets it; `nokia_srlinux` does not. The wait needs a scheme and a port
and nothing else: **no login crosses the queue for it**, and the connection is opened and
closed without a request.

### Consequences

- `ready` means, for a package that asks, that the node can be configured.
  `ready_after_s` in `twin.json` covers both waits for such a node.
- A transport that never opens fails under `readiness.timeout`, its message naming which
  of the two waits did not end.
- A package that asks for the wait and declares no `config.push` is refused at load under
  `psp.readiness.await_push_transport`: the request would be unsatisfiable.
- The push's retry policy is untouched. A refusal is not retried, and a push that could not
  be made is retried three times, against a node whose transport has already accepted
  once.
- Nothing in the mechanism names a platform, and the platform grep ([D-031](#d-031))
  passes.

---

<a id="d-030"></a>
## D-030 — The heartbeat cadence is 10s, and cancellation latency is bounded by it

**Status:** Accepted 2026-09-21

### Context

A host-bound activity records a heartbeat every `lab.HeartbeatInterval`, and the worker caps
how long the SDK may hold one back with `provision.HeartbeatThrottle`. The Temporal Go SDK
abandons an activity attempt at the **first heartbeat call that goes unanswered for about
half the throttle**, whether the workflow service is unreachable or merely slow
([verified-facts.md](verified-facts.md)).

The pair was 2s, then 5s. At 5s, a lab host that also carries Infrahub (whose repository
sync takes a whole CPU every minute and whose message queue was measured at 150–203% of a
CPU) cut working deploys short: four of six end-to-end attempts lost the three-node
fixture's `clab deploy` at its first failed heartbeat, once 17.8s into the deploy, while
containerlab was still working. That twin boots three SR Linux nodes at 2048 MiB each,
about 5.5 GiB, the heaviest thing the end-to-end tier builds.

### Options considered

1. **Raise the throttle alone.** Rejected: a test pins `HeartbeatThrottle <=
   lab.HeartbeatInterval`, for a reason: a throttle above the interval delays the
   heartbeats a cancellation travels on, so cancellation latency would grow anyway and the
   invariant would be broken rather than honoured.
2. **Shrink the three-node fixture to two nodes.** It addresses the load directly. Rejected:
   the three-node fixture is the goldens', the live fixture branch's and the recorded
   histories', and shrinking it weakens every three-node claim the tiers make: a test made
   easier rather than a mechanism made right.
3. **Raise the heartbeat timeout.** Rejected: it is the backstop for a worker that has
   died, not the tolerance for a slow answer, and it is not what abandons the attempt.
4. **Double both constants.** Chosen.

### Decision

`lab.HeartbeatInterval` and `provision.HeartbeatThrottle` are both **10s**.
`HeartbeatTimeout` is 30s, so the invariants hold: the throttle does not exceed the
interval, and it stays well inside the timeout, with three heartbeats to a timeout window.

### Consequences

- **A heartbeat call survives a stall of about 5s.** A deploy that is working is not
  abandoned because the host was busy for three seconds.
- **A cancellation reaches a running `clab` up to about 10s after it is asked.** That is the
  price, paid by `twin destroy` cancelling a create or a step, and by a check's rebuild.
  Nothing asserts a bound it breaks: the end-to-end tier allows a destroy 60s, and observed
  destroys take 2–9s.
- A host so loaded that it cannot answer a heartbeat for five seconds still cuts an
  activity short, under its own rule, retryable, and named as a lost heartbeat.
- No workflow, activity or wire type depends on either constant, and no recorded history
  encodes one.

---

<a id="d-031"></a>
## D-031 — A mechanism is code keyed by a format value; a platform is data

**Status:** Accepted 2026-09-21

### Context

[D-016](#d-016) says adding a platform is a data change, never a compiler change. With one
platform shipped and a second existing only as a test-only package that never booted, that
was easy to hold. A real second vendor forced the question. cEOS does four things SR Linux
does not:

- its gNMI server speaks **plaintext**, not TLS;
- its configuration arrives over **eAPI** in a named session, not JSON-RPC over a
  candidate;
- its containerlab kind takes a startup-config as the node's **whole** configuration, so a
  partial bootstrap cannot be delivered that way;
- its image is **account-gated**, so nothing may pull it.

Each is a genuine difference in how a node is reached, and none is expressible by a package
that only names strings for the existing code to interpolate. Four times, the cheapest
change was a branch on the platform (`if kind == "ceos"`), and four times that branch
would have made D-016 a slogan.

### Options considered

1. **Branch on the platform where a difference appears.** One `if` per difference, in the
   driver, the probe, the compiler. Honest about the cost of the first one and invisible
   about the cost of the tenth: every later platform pays for it, the compiler stops being
   the pure function [D-011](#d-011) claims, and there is no line at which someone is told
   to stop. Rejected. It is what D-016 exists to forbid, and the first branch is the one
   that is always defensible.
2. **A plugin interface per platform**, a Go interface each vendor implements. It keeps the
   branch out of shared code but puts the platform back into code: a platform is then a
   compiled artifact, the override directory ([D-018](#d-018)) can no longer add one, and
   two vendors that differ in one field carry two implementations. Rejected.
3. **Make the package a program**: an embedded expression language, or the push's request
   body as a template. Maximum expressiveness, and no new arm ever needed. Rejected: a
   package would become untestable by inspection, the conformance suite could no longer say
   what a package *means*, and the failure mode moves from "refused at load" to "wrong at
   the node".
4. **A fixed, named set of mechanisms, each selected by a format value.** Chosen.

### Decision

**A mechanism is code; a platform is data; the format is the seam between them.**

- The driver implements a fixed, named set of mechanisms: two push mechanisms (`json_rpc`,
  `eapi`), each with its replace ([D-033](#d-033)); a probe that dials with TLS or without
  (`readiness.tls`); a bootstrap applied at deploy or sent in the push
  (`config.bootstrap_via`); a presence check for an image that may not be pulled
  (`image.acquisition`); a wait for the push transport
  (`readiness.await_push_transport`, [D-029](#d-029)); and a reader of a booted node keyed
  by the package's `conformance` facet ([D-037](#d-037)). Each is reached only by a package
  declaring the value that selects it.
- **No code may branch on what a platform *is***: its id, vendor, NOS, containerlab kind or
  image reference. A tier-1 test (`TestCompilerNamesNoPlatform`) greps every non-test
  source under `internal/compiler`, `psp`, `lab`, `provision`, `stage`, `bundle`,
  `validate`, `step`, `waypoint`, `verify`, `api`, `server`, `cli` and `tree` for every
  platform name the tree knows (`srlinux`, `nokia`, `eos`, `arista`, `ceos`, `chassisos`),
  **comments included**: a comment that explains code by naming a platform is the same
  defect one step from being code.
- **A value this build has no arm for is refused at load, naming the arms there are**:
  never skipped, never defaulted at the moment of use. `psp.config.delivery.unimplemented`
  is the shape; `psp.config.bootstrap_via` and `psp.readiness.await_push_transport` refuse
  combinations no arm can honour, and `verify.package.unreadable` refuses, before any node
  is read, a package with no `conformance` facet or whose probe has no reader arm.
- **A platform that needs something new is a format change and a new arm**, weighed on its
  merits and available to every package, not a special case for the vendor that asked.

### Consequences

- **The grep is the enforcement.** Every mechanism above was implemented under it, and it
  is empty. A defect it catches is caught at `make test`, not at review.
- **The format grows.** That is the cost this decision accepts: the format is where
  platform differences are allowed to accumulate, and [psp/README.md](../psp/README.md)
  records each field under the version it arrived in, so the growth is legible.
- **A package can be wrong in a new way**, by asking for a combination no arm honours, such
  as `await_push_transport` with no push to wait on. Each is a consistency rule that
  refuses at load with both fields named, which is cheaper than a node refusing a request
  at 3 a.m.
- **The mechanisms are the evidence for the rule, not designed for one vendor.** `tls` is a
  boolean any package may set, `eapi` is one dispatch arm beside `json_rpc`, and a third
  vendor that speaks either needs no code. The test of this entry is the *third* platform,
  and what it should cost is a YAML file.
- **`fidelity.forwarding.mixed` retiring is the same principle read backwards**
  ([D-025](#d-025)): that refusal existed because the bundle had one field for a value that
  is per platform. Where the data model admits the difference, the refusal is not needed.
- Constitution II (Platform Abstraction Is Absolute) is tested by this rule; the entry adds
  no exception to it.

---

<a id="d-032"></a>
## D-032 — A waypoint is an operator-written name for a pinned reference

**Status:** Accepted 2026-09-28

### Context

[D-012](#d-012) gives an operator two ways to name intent: a branch, or a branch at a time.
Walking a change as a story needs a third. Think of a fabric before a migration, after its
first link moved, after its second. The operator wants to mark each chapter once, in order,
and build a twin from any of them. Without a name, the story lives in shell history, as
`--at` values composed by hand from a clock that is not Infrahub's, which is how a local
"now" of one-second granularity comes to hide objects written in that second.

On 2026-09-28 the operator decided that the operator writes waypoints and Fylgja reads
them. This entry records that, and what was built on it. What Infrahub 1.11.2 does with
such a kind ([verified-facts.md](verified-facts.md)):

- A branch-agnostic kind loads with a uniqueness constraint over two attributes.
- An object written from one branch is read from every branch whose schema has the kind.
  It outlives the branch it names. A branch whose schema lacks the kind cannot query it at
  all.
- At creation every attribute carries one `updated_at`. Only a real change of `branch`
  moves `branch.updated_at`, and a read at that stamp sees every write made before the
  waypoint.
- An attribute name is 3–64 characters.

### Options considered

Where the waypoint lives:

- **A branch-aware kind on the branch it marks.** Every waypoint would land in that
  branch's diff and its merge, and a series that marks several branches in turn would be
  scattered across them. Rejected for [D-010](#d-010)'s reason: an object Fylgja needs,
  flagged on a branch under test, pollutes that branch's diff.
- **A branch-agnostic kind.** Chosen. A waypoint marks a branch without landing in its
  diff, and one series may mark several branches.

Who writes it:

- **Fylgja writes waypoints**, through a command that marks the branch's head. Rejected:
  every operation the product makes on Infrahub is a query (Constitution III), and a writer
  would be its first mutation.
- **The operator writes them, by Infrahub's own means**: an object file, the UI, or
  `infrahubctl`. Chosen. Test support writes them for tests, behind a build tag.

What it is:

- **A provenance field in the bundle**, carrying the waypoint's name. Rejected: the same
  intent would hash differently according to what it was called. [D-023](#d-023) keeps
  `observed_at` out of the bundle for the same reason.
- **Inside the generics contract, as a new contract version.** Rejected: the goldens, the
  fixture ids and the fixture branch would move, all for a kind nobody implements. The
  waypoint lives beside the contract instead, as `FylgjaContract` does, and the contract
  stays `0.2`.
- **A required `at`.** Rejected: the operator would have to compose a timestamp, which is
  the very act waypoints remove.
- **An unpinned waypoint**, where an absent `at` means the branch's head. Rejected: the
  waypoint would name a moving thing, and two creates from it could differ.
- **An absent `at` means the waypoint's own write time.** Chosen.

### Decision

- **A waypoint is a name for a pinned reference.** The kind is `FylgjaWaypoint`
  ([schema/fylgja-waypoint.yaml](https://github.com/happypathnetworking/fylgja-prototype/blob/1982b7263e0e79c80988ea2055266f3fcc697c72/schema/fylgja-waypoint.yaml), as the tree held it before 014 removed it). It has a series, a
  sequence, a branch, an optional `as_of` and a description. It is unique on series and
  sequence, and branch-agnostic. `--waypoint <series>/<sequence>` on `twin create`, `intent
  read` and `twin step` is resolved to `(branch, at)` before any run starts, at the
  findings step `resolve`. The run is then a pinned run.
- **The operator writes it, and Fylgja only reads it.** Fylgja reads it from the default
  branch through the unnamed endpoints (`/api/schema`, `/graphql`), so it names no branch
  to do so.
- **An absent `at` is `branch.updated_at`, as Infrahub returns it.** That is the stamp
  Infrahub wrote when the waypoint last said which branch it seals. A description edit
  leaves it where it is, and re-pointing the waypoint moves it. A written `at` is used as
  written. Either one is passed verbatim.
- **The attribute is `as_of`, not `at`**: Infrahub refuses an attribute name shorter than
  three characters. Everything Fylgja prints, records or accepts says `at`.
- **The kind lives beside the contract.** It is loaded on the default branch, where
  Fylgja's schema is, and nobody is asked to implement it. Without it the waypoint commands
  refuse under `waypoint.kind.absent`, naming the file. Nothing else reads the kind.
- **The name is not provenance.** It enters neither the CTM nor the bundle. The twin's
  record names it ([D-014](#d-014)).
- **A step's `unchanged` compares content, never ids.** `at` is in the hashed provenance
  ([D-024](#d-024)), so two waypoints never share a `bundle_id`, even when the branch did
  not move between them. The step's comparison leaves `provenance` out, and `waypoint plan`
  prints both ids beside every step.

### Consequences

- A twin built from a waypoint is pinned. It never follows ([D-026](#d-026)), and no check
  resolves a waypoint. It walks its series by `twin step` ([D-036](#d-036)), without a
  rebuild.
- `waypoint list` and `waypoint plan` need no lab, worker or workflow service. The plan
  reads and compiles each waypoint through the stage pipelines, files each bundle in the
  store ([D-011](#d-011)), and prints the step between consecutive pairs.
- **The operator has one rule: write the waypoint last**, after the branch's writes and
  after generating ([D-028](#d-028)). A waypoint written before a generate seals the old
  artifact; the plan then shows a topology step with no artifact change. Fylgja reports
  this, and never infers staleness.
- Editing a waypoint changes no twin built from it. The record keeps the reference the
  twin was built from, and `waypoint list` warns `waypoint.twin.moved` beside a row that
  has moved.
- The resolution refusals (`waypoint.ref.invalid`, `waypoint.flags.conflict`,
  `waypoint.kind.absent`, `waypoint.unknown`, `waypoint.duplicate`,
  `waypoint.at.unresolved`) are all made before any run. A given `at` later than now is
  refused until it has passed: Fylgja reads its clock to compare with a written `at` and
  never sends that reading, and it never compares Infrahub's own `updated_at` stamp with
  it.
- Every operation Fylgja makes on Infrahub is still a query, held by
  `internal/intent/readonly_test.go`. An unwritten `at` is Infrahub's stamp of the
  operator's write, read rather than composed, so Constitution VI's "Fylgja MUST NOT send
  an `at` of its own" holds.
- Revisit if Fylgja must write a waypoint itself, for instance a step that marks where it
  stopped. That reverses this entry's first decision, and needs a new entry.

---

<a id="d-033"></a>
## D-033 — A replace resets the candidate to the node's own baseline: S1 on SR Linux, E1′ on EOS

**Status:** Accepted 2026-10-01

### Context

A step takes a running twin from one waypoint's bundle to another ([D-036](#d-036)). A
`merge` push from artifact A to artifact B would leave behind whatever A added that B
lacks, so a step needs a **replace**, and [D-009](#d-009) bounds it: Fylgja never renders or
rewrites configuration, so the device must compute the difference, not Fylgja. A replace by
`delete /` before the artifact commits an empty configuration and takes the management
plane dark. The artifact carries no management: on both platforms, containerlab's
configuration does.

containerlab leaves that management on every node as a **baseline**:

- SR Linux's startup configuration and its checkpoint `clab-initial` hold management plus
  the deploy's bootstrap.
- cEOS's `flash:startup-config` holds containerlab's management configuration.

Nothing Fylgja does writes either one, because no push saves.

Five SR Linux methods, five EOS methods and the Python tools were weighed, and nine of the
methods tried against a lab of SR Linux 24.7.1 and cEOS 4.32.0.2F under containerlab 0.79.0
on 2026-10-01.

### Options considered

SR Linux:

- **S1: load the baseline into the private candidate, then send the bootstrap and the
  artifact, then commit.** Chosen. `load` replaces the candidate, uncommitted edits
  included. It works by checkpoint name or as `load startup`, through the JSON-RPC `cli`
  method, in 2.2–4.5s. `diff flat` gives the device's diff.
- **S2: a gNMI Set, with a CLI-origin replace carrying the baseline text and an update
  carrying the artifact.** It works, but "blank" is empty, so every replace must carry
  containerlab's management text. It took 19.3s for a 573-line baseline, and a client
  deadline does not stop the server's commit.
- **S3: scoped deletes kept as package data.** It leaves behind anything A added outside the
  scopes.
- **S4: a JSON-RPC `set` replace at `/`.** It needs JSON, and the whole tree, management
  included.
- **S5: NETCONF, `copy-config` startup→candidate, then the artifact as XML.** It works, in
  about 0.7s. Against it: it commits through the shared candidate, which a confirmed-commit
  revert leaves dirty until `discard-changes`; a confirmed commit needs `persist` to outlive
  its request; NETCONF gives no device diff; and SR Linux's artifact would become XML while
  EOS's stays CLI, where one template renders both.
- **napalm-srlinux.** Its replace is `delete /` plus the lines, which darkens the node, and
  it is Python ([D-005](#d-005)).

EOS:

- **E1′: in a named configuration session, `rollback clean-config` then `copy
  startup-config session-config`, then the bootstrap and the artifact, then commit.**
  Chosen. The pair makes the session exactly the baseline; the copy alone would merge. It
  took 1.6s. `show session-config diffs` gives the device's diff.
- **E1: the same session, with the baseline's lines read from `show startup-config` and
  replayed.** It works, in 1.4s. But Fylgja would read, strip and replay the node's own
  configuration, hashed secret included.
- **E2: `configure replace flash:<file>`**, with the file written by the worker into the
  lab directory. It ties the push to the lab host's paths, which is the seam M8 and M9
  change.
- **E3: `configure replace startup-config`, then a merge.** It is not atomic, and sent
  inside a session it applies at once and ends the session.
- **E4: a gNMI Set, with a CLI-origin replace at the root.** It works on 4.32, but the blob
  must carry the baseline, and cEOS silently ignores gNMI's commit-confirmed extension.
- **E5: NETCONF.** It has no reset into the candidate; its `arista-cli` channel refuses any
  repeated line, and every artifact repeats lines; an OpenConfig root replace leaves behind
  any configuration OpenConfig does not model; a confirmed commit is silently ignored; and
  a refused edit drops the session.
- **NAPALM, Ansible `arista.eos`, AVD and scrapligocfg.** They use E1's pattern, but they
  are Python or SSH screen-scraping, and each also saves the running configuration, which
  would overwrite the baseline.

Both platforms:

- **Infrahub renders the whole configuration, management included**, so each vendor's
  native full replace works as written. Rejected: Infrahub would hold twin facts (the
  management address, the TLS profiles and the server settings), a login would enter the
  artifact and the bundle, and it would reverse D-009, Constitution IV and
  [D-028](#d-028)'s artifact.

### Decision

**A replace resets the candidate to the baseline the node already holds, sends the bundle's
bootstrap and then the artifact, and commits.** It is one request with one atomic commit,
and the device computes the difference.

- **SR Linux (`delivery: json_rpc`):** `enter candidate private`, `load startup`, the
  bootstrap's lines, the artifact's lines, `diff flat`, `commit now`. `load startup` is
  chosen over `load checkpoint name clab-initial`, which gave the same candidate: both
  packages then reset to the same thing, the startup configuration, under the one rule that
  nothing saves it; the checkpoint is pruned once ten newer ones exist and its id shifts as
  older ones go; and `load startup` measured 2.2s against the checkpoint's 4.5s.
- **EOS (`delivery: eapi`):** `enable`, `configure session fylgja-<attempt>-<ns>`,
  `rollback clean-config`, `copy startup-config session-config`, the bootstrap's lines, the
  artifact's lines, `show session-config diffs`, `commit`, sent with `format: "text"`,
  since the diff has no JSON model and `format: "json"` fails the whole request. A refused
  request's session is aborted before the next attempt.
- **The reset is code keyed by `delivery`** ([D-031](#d-031)). A package asks for a replace
  with `config.mode: replace`, and both shipped packages do, so a create and a step share
  one push: at create the node holds only its baseline, so a replace lands what a merge
  would. A package left on `merge` can create but cannot step, and is refused naming it
  (`step.package.merge`).
- **Under `mode: replace` the bootstrap is always sent after the reset**, whatever
  `config.bootstrap_via` says, since the baseline holds the bootstrap as it was deployed and
  a topology step changes it, and containerlab's reconcile does not apply a changed startup
  snippet to a running SR Linux node. `bootstrap_via` says only how a deploy delivers the
  bootstrap: containerlab's startup configuration, or the push.
- **Fylgja never saves the running configuration, on any platform.** The startup
  configuration is the baseline every replace resets to.
- **What a step records of a node's change is the device's diff** (`diff flat`, `show
  session-config diffs`). Fylgja computes none.
- **A refused `json_rpc` replace may send a second request**, which asks the node for the
  reason its first answer did not carry and lands nothing ([D-034](#d-034)). A replace that
  commits is always one request.

### Consequences

- **No baseline text is read, carried or replayed.** The one request joins only the
  bootstrap and the artifact. Joining them is not rendering, so D-009 holds. No
  configuration is read back from the node at push time, so nothing new can leak.
- **The replace commits at once.** A safety net, if one is ever wanted, is the device's own
  timer (SR Linux's `commit confirmed`, confirmed with `/tools system configuration
  confirmed-accept`; EOS's `commit timer`, confirmed by committing the session again), never
  gNMI's or NETCONF's commit-confirmed, which cEOS 4.32 accepts and ignores.
- **A topology step has an order**: the reconcile, by `clab deploy`; readiness of every
  node containerlab restarted (a restarted cEOS node answered gNMI after about 49s and eAPI
  after about 55s); then the replace, on every node whose content changed **and on every
  node that was restarted**, since a restarted node's push is gone even when its artifact
  is unchanged. Which kinds restart on a link change is fidelity a package declares
  (`fidelity.link_change`); `clab deploy --dry-run` names those nodes before anything is
  applied.
- **The push meets two node behaviours.** A refused EOS request leaves its session pending
  until it is aborted, and EOS allows five pending sessions. A refused SR Linux line stays
  in the login's private candidate across requests, and the next replace's `load startup`
  clears it, so no `discard` is sent.
- **The baseline is containerlab's**, so it is re-verified on a containerlab or image
  upgrade: what each node's startup configuration holds, and the E1′ pair.
- Constitution I (the existing transports and Go clients, no second runtime), II (keyed by
  `delivery`; what a link change does to a kind is package data), IV (the push carries
  exactly the artifact and the PSP's bootstrap; the reset is a device command, not
  configuration Fylgja wrote) and X (the device's diff is the record of what a step changed,
  and a restarted node's lost push is stated) hold unamended.
- **Revisit S5** if a structured artifact is ever adopted for SR Linux for its own sake,
  **and E4** if an EOS release honours gNMI's commit-confirmed extension.

---

<a id="d-034"></a>
## D-034 — A refused SR Linux replace asks the node again for its reason, in a request that lands nothing

**Status:** Accepted 2026-10-02

### Context

A refused SR Linux push keeps only the reason, because the JSON-RPC error echoes every
command sent, the artifact included. Fylgja finds the reason after the echo, at `' failed
with error '`. SR Linux cuts the error message at about 1 KiB.

Under the replace ([D-033](#d-033)), `enter candidate private`, `load startup` and the
bootstrap's lines come before the artifact. A real replace of 33 commands made an echo of
1604 bytes alone; the node sent back a 1032-byte message that ended inside the echo, so the
reason (`invalid port 99, max is 58`) never arrived, and the finding said the message was
withheld. So every real SR Linux artifact refused under a replace would give `push.refused`
with no reason. Nothing lands, and nothing leaks.

### Options considered

- **Accept it.** Reword the withheld sentence to say the node cut its own message. It keeps
  one request, but `push.refused` would say nothing useful on the platform where refusals
  are most common (the package's port ranges).
- **Commit in a request of its own.** Send load, bootstrap, artifact and `diff flat`, then
  `commit now` alone; the private candidate persists across requests, and the commit's echo
  is short, so a reason given at commit would fit. Rejected: every replace becomes two
  requests, the successful ones too, and a line refused at parse time in the first request
  still echoes all of them.
- **Re-commit the refused candidate.** Rejected: if the first request was refused at parse
  time, the candidate holds the reset and only part of the content, and committing it would
  land that partial configuration.
- **Ask again with `commit validate`.** Chosen.

### Decision

**A replace that commits stays one request with one atomic commit. A `json_rpc` replace
whose refusal gives no reason gets one more request:** `enter candidate private`, `commit
validate`.

- **It lands nothing.** `commit validate` checks the private candidate the refused request
  left without committing it. Its echo is two commands, so a reason given at validation
  fits under the node's cap.
- **The reason is kept as any refusal's is**, through the same reading of the answer. If
  that answer gives no reason either, the finding says so in [D-035](#d-035)'s words.
- **It runs inside the push's own budget** (`push_timeout_s`). It is never retried. It is
  not a new step, activity or identifier. The refusal stays a refusal, and the second
  request cannot change that.
- **The candidate is left as a refused replace leaves it.** The next replace's `load
  startup` clears it.
- **It is code keyed by `delivery`** ([D-031](#d-031)). `eapi` names the refused command in
  its own answer, so it gains nothing.

### Consequences

- A refusal costs one more round trip, about a second, inside a budget of 30s.
- `commit validate` recovers the reason of a schema constraint, and not of a check a
  commit makes as it applies, which [D-035](#d-035) words.
- The refusal is read before the message's first empty line, since no command sent is
  empty, so a configuration value that happens to carry the marker in the diff that
  follows the echo never names the wrong command.
- Constitution I (the existing transport, no new client), II (keyed by `delivery`), IV (the
  second request carries no configuration, and only a reason is kept) and X (a refusal says
  why, where it had said nothing) hold unamended.

---

<a id="d-035"></a>
## D-035 — A commit-time refusal on SR Linux gives no reason under the replace; the request after it stays

**Status:** Accepted 2026-10-02

### Context

[D-034](#d-034) sends one more request, `enter candidate private`, `commit validate`, after
a `json_rpc` replace whose refusal gives no reason, on the assumption that `commit validate`
reports the same reason as the refused commit. Live, on SR Linux 24.7.1, that holds for one
class of refusal and not for another:

- **Schema constraints.** `commit validate` refuses a VLAN encapsulation on a port without
  `vlan-tagging`, and its 287-byte answer carries the reason. At a real replace's length
  (26 commands, a 1238-byte echo), the first message was cut and the second request
  recovered the reason.
- **Checks made as a commit applies.** For the port-range check (`invalid port 99, max is
  58`), SR Linux answers `All changes are valid.` Only a commit carries that reason; a short
  `commit now` on the refused candidate gave it in 196 bytes.

Re-committing is ruled out ([D-034](#d-034)): a cut message cannot be told apart from a
parse-time refusal, and committing a partial candidate would land it. A commit in a request
of its own remains the only way to get commit-time reasons back.

### Options considered

- **Say the reason is missing, and keep the request.** Chosen.
- **Say the reason is missing, and remove the request.** Rejected: the request recovers
  schema-constraint reasons for the cost of one request after a refusal. Taking it out
  would refuse those reasons for no gain.
- **The commit in a request of its own on every replace.** Rejected, as in D-034: every
  successful replace becomes two requests, and a refusal at parse time in the first
  request still echoes every command.

### Decision

- **A refused `json_rpc` replace whose message gives no reason is asked again by D-034's
  `commit validate`, unchanged.**
- **When that answer gives the reason, the reason is the finding's**, worded as given at
  `commit validate`.
- **When it gives none, the finding says so in one of two sentences, and quotes nothing:**
  - the candidate validates: `the node's message ends before any reason Fylgja can read,
    and the candidate it refused passes `commit validate`, which does not run every check a
    commit does (D-035)`;
  - the answer is cut, refused in a shape Fylgja does not read, fails or never comes back:
    `the node's message ends before any reason Fylgja can read, and asking again with
    `commit validate` gave none (D-035)`.
- **A merge, and a refusal whose reason fits, are as they were**: the merge keeps its
  withheld sentence.

### Consequences

- **A commit-time refusal on SR Linux names the node, the artifact and its checksum, but
  not the reason.** The operator finds the reason by pushing the artifact by hand, or from
  the package's port ranges, which `twin compile` already enforces for any interface the
  profile maps. A port out of range reaches the push only from an artifact naming an
  interface the CTM does not carry.
- Both sentences are the CLI's contract's, and tier 1 holds them word for word.
- Constitution IV (neither sentence quotes the node) and X (a refusal states that its
  reason is missing, and why) hold unamended.
- **Revisit** if an SR Linux release makes `commit validate` run the commit-time checks, or
  lifts the 1 KiB cap on its error message. Either would bring the reason back with no
  change to the code: the first through the existing request, the second through the
  first answer.

---

<a id="d-036"></a>
## D-036 — A failed or cancelled step leaves the twin up and diverged; a diverged twin accepts only `twin destroy`; a restart needs `--allow-restart`

**Status:** Accepted 2026-10-02

### Context

`twin step` takes a running waypoint twin to another waypoint of its series without a
rebuild. Every other run that touches the host ends in one of two states: the twin is
ready, or nothing is left. That rule is that **cleanup survives everything**: a provisioning
run that fails or is cancelled after its host check tears the lab down and removes the twin
directory, on a disconnected context ([D-007](#d-007)).

A step cannot follow that rule as written. Before the step there was a ready twin the
operator built and may have been working on for an hour. A step reconciles the running lab
(`clab deploy` without `--reconfigure`), waits for the nodes containerlab restarted, and
replaces each changed node's configuration ([D-033](#d-033)). Any of these can fail part
way. Tearing the twin down on failure destroys the operator's work for a fault in one push.

Two of the step's facts were verified live before the design depended on them
([verified-facts.md](verified-facts.md)):

- containerlab's dry run (`clab deploy --dry-run --format json`, with `CLAB_LABDIR_BASE` at
  the twin directory) names every node it will restart, recreate or create, and writes
  nothing.
- A cEOS node restarts in place when a link of its is added or removed, and returns about a
  minute later on its startup configuration, with its push lost. SR Linux re-cables live. A
  node whose image or kind changes is recreated.

On 2026-09-28 the operator decided that a failed step leaves the twin up and diverged, so
that a later rollback can start from it, and that a cancelled step is a failed step. On
2026-10-01 the operator answered what a diverged twin may do next, what a step that would
restart a node does, and where the device's diff lives.

### Options considered

What a failed or cancelled step leaves:

- **Tear the twin down, as a provisioning run does.** Rejected: it destroys a working twin
  for a fault in one node's push, and the operator loses the state the step was meant to
  keep.
- **Roll back automatically** to the waypoint the twin came from. Rejected: a rollback is a
  second step, which can fail in its turn, and a step that restarted a node would restart
  it again. It is planned for later, from the record this entry keeps.
- **Cancel to the previous waypoint**, so that `twin destroy` during a step first undoes it.
  Rejected: the operator asked to destroy, and the destroy follows.
- **Leave the twin up and diverged, and record where it stopped.** Chosen.

What a diverged twin may do next:

- **Step again, to any waypoint of its series.** Rejected: the step would start from a twin
  no bundle describes, and its diff and push plan would be computed against the wrong side.
- **Repair: step again to the same target.** Rejected for the same reason, and because a
  repair that fails again leaves two diverged records to explain.
- **Only `twin destroy`.** Chosen. The record keeps what a later rollback or repair needs.

A step whose plan restarts or recreates a node:

- **Proceed and record the restart.** Rejected: a reboot the operator did not expect loses
  the node's running state, and a minute of a cEOS node's time.
- **Refuse on any kind that restarts.** Rejected: EOS could then never take a topology step.
- **Require `--allow-restart`.** Chosen. Without it the step is refused before the host is
  touched, naming each node with containerlab's reason and its package's declaration, so
  the dry run is the normal first step of a topology change.

### Decision

- **A step that fails or is cancelled after it has touched the host leaves the twin up and
  diverged.** The host is touched from the **stage**, the swap of the twin directory's
  bundle to the target's. A refusal at `inspect`, `host check` or `plan reconcile` touches
  nothing and is a rejection, exit 1. From the stage on, every failure and every
  cancellation ends `diverged`, exit 4, with `step.diverged` naming the phase
  (`reconcile`, `readiness`, `push` or `record`) and which nodes the push reached. Nothing
  is torn down and nothing is unstaged.
- **What cleanup means for a step is the record.** `RecordStep` runs on a disconnected
  context after the stage, whatever happened, and writes `twin.json` whole. The diverged
  record keeps the previous waypoint and bundle at its top level, names the target, the run
  and the phase in its `step` block, and gives each node the bundle it holds, or `null` for
  a node containerlab restarted that the step never pushed ([D-014](#d-014)). `RecordStep`
  is idempotent: a record whose step block names this run is returned as written.
- **A diverged twin accepts only `twin destroy`.** `twin step` and its dry run refuse it
  under `step.twin.diverged`, naming the run, the target, the phase and the remedy. `twin
  show` names it as the kind `diverged`.
- **A step whose plan restarts or recreates a node needs `--allow-restart`**
  (`step.restart.required` without it). A node containerlab creates for the step needs no
  flag, because it loses nothing. Every restarted, recreated or created node is awaited
  ready and then pushed, whatever its content.
- **The device's diff lives in the run's history alone**, as the push activity's result,
  bounded by `lab.DiffLimit` and stripped before it enters any other input or result. No
  record, document, finding or log line carries it, and `twin.json` holds no count of it.
- **A run in flight is `run.in_flight`**, naming the run, for a step beside a provision, a
  destroy or another step, and for a create or provision beside a step.
- **An image or kind change is a recreate**, decided by containerlab from the lab
  directory's state. A change the plan would not apply (a package change under the same
  image and kind) is refused as `step.node.unapplied` before anything is touched.

### Consequences

- **Exit 4 has two status words.** It says something remains on the host and the report
  says how to clear it: `unclean` for a create or destroy, `diverged` for a step.
- **`twin destroy` during a step cancels it, waits for its record, then destroys.** The
  record is written diverged on the disconnected context, and then the destroy removes the
  lab and the directory, the diverged record with them. The recorded history
  `internal/provision/testdata/step-cancelled.history.json` is such a run.
- **A diverged twin is described in one line**, in `twin show`, in `waypoint list`'s mark
  and in the refusal a second step gets.
- **Rollback and repair start from the diverged record.** It holds the waypoint and bundle
  the twin came from, the target, the phase and each node's state. A later feature adds
  the step it needs, as a new entry.
- **A topology step costs a restarted node's minute.** The steps that restarted a cEOS node
  took 62–81s end to end, and an artifact-only step 5.5s.
- `fylgja-step` has a fixed id, is refused beside any other run and is cancelled by `twin
  destroy` (Constitution VII). Its activities are idempotent, pass paths and ids, and
  heartbeat under the packages' budgets, and its cleanup is the record (VIII). Whether the
  twin exists is containerlab's to say, and the lab directory's own state is read by
  containerlab's dry run, never by Fylgja (IX).
- Revisit if a step ever needs to leave nothing behind, for instance when twins share a
  host after M9. That reverses this entry's first decision.

---

<a id="d-037"></a>
## D-037 — Observed facts are a Fylgja-owned vocabulary, mapped to each platform's paths by its package

**Status:** Accepted 2026-10-03

### Context

An earlier design had verification read nodes through a Fylgja-owned normalized state
model shaped by OpenConfig, and parked that until verification was planned. The
conformance suite's boot half reads five facts of a booted node (host name, port enabled,
port discovering, neighbour, version) over each package's readiness transport, by paths the
package's `conformance` facet declares, with `{node_name}` rendered, a declared meaning for
an absent value, and neighbour leaves relative to one entry. `twin verify` reads three of
the five against intent. The question is what model those reads are of.

### Options considered

- **Adopt OpenConfig as the model**, translating each platform's answers into it.
  Rejected: SR Linux answers in its native model, cEOS in OpenConfig, and a translation
  layer is a second model to maintain for five facts.
- **Assertion code per platform.** Rejected by [D-031](#d-031) and Constitution II.
- **The facet as built.** Chosen: the facts are Fylgja's names, each package maps them to
  its own paths and values, and the reader is one mechanism keyed by the facet. "Shaped by
  OpenConfig" means the names alone.

### Decision

Fylgja owns a vocabulary of observed facts: `host_name`, `port.enabled`,
`port.discovering`, `port.neighbor` (with `system_name` and `port_id`) and `version`. A
package maps each to a path, a value and what absence means, through its `conformance`
facet; the paths are data ([D-031](#d-031)), and no model is adopted. `twin verify`
asserts three (host name, port enabled, neighbour) against intent; the suite asserts all
five against the package. The readers are product code, `internal/verify`, which `twin
verify`, the step's `VerifyTwin` and the suite share; `internal/conformance` stays a
`go test` package that no product code imports.

### Consequences

- A new fact is a new name in the facet and one reader, never a model import.
- A platform whose transport is not gNMI needs a new reader arm under [D-031](#d-031), and
  is refused by verify until it exists (`verify.package.unreadable`).
- The operational state diff on the [roadmap](roadmap.md#backlog)
  extends the vocabulary (adjacencies, routes) rather than adopting a model.
- Every assertion carries a stable identifier and names what was read; fidelity is still
  asserted, never measured (Constitution X).
- Revisit if a third platform's answers cannot be named by a path and a value.

---

<a id="d-038"></a>
## D-038 — Assertions are data derived from the bundle; no operator-written assertion is shipped

**Status:** Accepted 2026-10-03

### Context

An earlier design decided that assertions would be declarative, never code, and parked it.
`twin verify` needs three kinds of assertion about a running twin: each node's host name is
its node name, every cabled port intent enables is enabled, and every link is an LLDP
adjacency seen from both ends. Each is derivable from what the twin was built from.

### Options considered

- **Operator-written assertion files.** Rejected: nothing asks for them, and
  [D-005](#d-005)'s hedge already requires that any later operator surface be declarative.
- **Assertions as Go code per check.** Rejected: it would branch on a platform and bypass
  the manifest.
- **Derived from the CTM.** Rejected: a twin provisioned from a bundle has none, and the
  CTM sidecar is written by a read, not by `twin provision`.
- **Derived from the staged bundle's manifest.** Chosen: it is beside every twin, carries
  the names the node uses (`node_name`) and the links, and carries each interface's intent
  `enabled` state, so a port intent disables is skipped and recorded, never asserted.
- **Read the artifact back from the node.** Rejected by the operator: a per-package
  running-configuration read is a new mechanism under [D-031](#d-031) and a research
  question of its own, so configuration conformance is reported from the record, labelled
  as its claim.

### Decision

`twin verify`'s assertions are `verify.Derive(manifest)`: a pure function of the staged
bundle's manifest, three kinds and no more, with a port intent disables and its link
skipped at both ends and each skip recorded. A node's configuration is reported from
`twin.json`'s `nodes[].holds` against the staged bundle, as the record's claim, and a claim
not held is a finding (`verify.record.holds`), as is a node the staged bundle names that
the record does not. No operator writes an assertion; a later surface, if one comes, is
declarative ([D-005](#d-005)).

### Consequences

- Each mapping row carries intent's `enabled` (bundle `4`, [D-022](#d-022)).
- A diverged twin is verified against the staged bundle, the target it was going to, with
  its divergence named.
- Probes and the operational state diff on the
  [roadmap](roadmap.md#backlog) add assertion kinds as data under this
  entry.
- Nothing is rendered and the artifact is never read back (Constitution IV); the compiler
  stays pure (V); the manifest beside the lab is the provenance of what runs, and the
  assertions come from it (IX); skips are recorded with their reason and the record's claim
  is labelled as such (X).
- Revisit when an operator needs to assert something intent does not carry.

---

<a id="d-039"></a>
## D-039 — Verify is advisory: no Fylgja operation reads the report, and the report is the findings document with a `verify` block

**Status:** Accepted 2026-10-03

### Context

An earlier design decided that verification reports and never blocks, and parked it until
proposed-change integration. `twin verify` builds the report; posting it on a proposed
change comes later. The report needs a shape, a status, a place in the step workflow, and a
budget for its waiting form. The operator answered these on 2026-10-03.

### Options considered

- **A gate** (a nonconforming twin refuses a step, stops following, or fails a proposed
  change). Rejected: advisory means nothing in Fylgja acts on the report, and a consumer
  reports it and never blocks on it.
- **A document of its own with its own version.** Rejected: the findings document `1` with
  a `verify` block, as the `plan` and `step` blocks are.
- **Exit 0 whenever the read succeeded**, the findings alone saying so. Rejected: a script
  would parse JSON for the one thing it asks.
- **Exit 1 as a rejection.** Rejected: it conflates intent not met on the node with a
  command that refused to run.
- **A status of its own, `nonconforming`, exit 5.** Chosen; a node that could not be read
  is `operation.failed`, exit 2.
- **A per-package wait budget** (a PSP field per NOS's LLDP timers). Rejected: one default
  of 120 seconds, shared by `twin verify --wait` and the step, overridden by `--wait` on
  `twin step`.
- **The record's claim as information only.** Rejected: it counts toward the status, since
  a diverged twin could otherwise end `ok` while running the wrong configuration.
- **Skipping the wait after a divergence.** Rejected: it runs after every record.

### Decision

`twin verify` reads and reports, and changes nothing. Its report is the findings document,
operation `twin.verify`, with a `verify` block under `verify.schema.json`: one finding per
failed assertion and per record claim not held, status `nonconforming` (exit 5) when every
node was read and any failed, and `operation.failed` (exit 2) when a node could not be
read. `--wait[=<duration>]` reads every second until the twin settles or the budget (2m0s
bare) expires. No Fylgja operation reads the report: `twin create`, `twin step`, `twin
destroy`, following and the step's outcome are unchanged by it. The step workflow runs
verify's waiting form after its record as one activity, `VerifyTwin`, under a version
marker, bounded by the default budget or `twin step --wait`, and never fails the step for
it ([D-007](#d-007)); `twin.json` records how the wait ended and when the twin settled
([D-014](#d-014)).

### Consequences

- Posting the report on a proposed change, which the
  [roadmap](roadmap.md#backlog) plans, reports the block and never
  blocks on it.
- The budget is the operator's, by flag: it bounds how long the operator will wait, not how
  a platform behaves, which is why it is not a package field. Constitution II's "budgets
  from the PSP" governs platform behaviour; this one governs patience.
- A twin that never conforms costs a step its budget, once, and the record says so; a wait
  that did not settle is the warning `verify.wait.unsettled`.
- A read stops dialling a node that does not answer (a gRPC `Unavailable` or
  `DeadlineExceeded`), so a silent node costs one probe deadline per read and the wait
  stays inside its bound. A node a read could not dial is looked up again from containerlab
  before the next read, and an address that is not an IP address is none.
- `VerifyTwin` is one bounded, heartbeating, idempotent activity on the run's own context;
  the record is unchanged by it but for the wait block (Constitution VIII). Every finding is
  identified, the record's claim is labelled, skips are recorded, and no report carries
  configuration content or a credential (X).
- Revisit if a consumer needs a gate: that is a new entry, not a flag.

---

<a id="d-040"></a>
## D-040 — The API is the boundary of the core: a user-facing client reaches Fylgja through it alone, and test code is excepted

**Status:** Accepted 2026-10-04

### Context

Before the API, the CLI was the core's only client, and it linked the core. Each command
ran part of the system in the operator's shell: the read and the compile, the waypoint
resolution, the filing of a bundle in the store, the host check, containerlab's plan, the
record, the reads of the nodes, and the workflow service's client. So the operator's shell
held every credential the system has: Infrahub's token, both node logins and access to the
workflow service.

An API was first planned beside the CLI, with the stage commands left in the CLI's process.
Six items on the [roadmap](roadmap.md#backlog) stand on an API: two
proposed-change checks, a webhook, an agent surface, observation and a web UI. Each is
another client, and a rule that lives in the CLI is one each of them must repeat or go
without.

On 2026-10-04 the operator set the rule this entry records: the API is the primary boundary
of the core, a user-facing client (the CLI or a web UI) interacts with Fylgja only through
it, and testing is an exception, because not every test case can, or should, run through
the API.

### Options considered

- **An API beside the CLI.** `fylgja serve` exposes the lifecycle to later clients, and the
  CLI goes on linking the core. Rejected: every operation then has two paths, a guard
  exists twice or differs, the tiers prove the path the CLI takes while every other client
  takes the other, and the operator's shell keeps every credential.
- **The CLI as the API's client for the lifecycle alone**, with `intent read`, `schema
  check`, `twin compile` and `psp validate` in the CLI's own process as developer tools.
  Rejected: two of the four read Infrahub with Infrahub's token, so the shipped CLI would
  keep a path to the core that a web UI could not use, and `schema check` would be built
  again for the first client that wanted it.
- **A split by what a command touches**: the two that read Infrahub behind the API, and
  `twin compile` and `psp validate` in the CLI as pure functions over local files, so that
  a package author needs no server. Rejected by the operator: the rule would name its
  exceptions, the import test would let the compiler into the client's package, and a
  client's packages could differ from the server's, so a bundle compiled in a shell need
  not be the one the server compiles.
- **The API as the only path for a user-facing client, for every command, with test code
  excepted.** Chosen.
- **No exception: every test through the API.** Rejected. Tier 1 calls the pure core, the
  workflows and the activities by design ([D-017](#d-017)). The contract tier reads a real
  Infrahub through the packages it tests. The end-to-end tier is worth its minutes because
  its oracles are independent: it reads the nodes with `gnmic`, the containers with Docker
  and the record from disk, and an oracle that goes through the system under test proves
  less. The hand checks stop a worker or the workflow service, which no API call should be
  able to do.
- **A direct mode in the CLI for tests**: a flag or a variable under which a command links
  the core. Rejected: it is a second path in a shipped client, which is what the rule
  forbids, and a test that takes it proves a path no operator takes.
- **The later clients' endpoints built with the API.** Rejected: each is a contract with
  one consumer that does not exist yet.

### Decision

The API is the boundary of the core.

- **A user-facing client reaches Fylgja through the API alone.** The CLI is the first
  client. A web UI, an agent surface and a check running inside Infrahub are later ones. A
  client sends requests to `fylgja serve`. It does not import, link or run the read, the
  compile, the waypoint resolution, the bundle store, the workflow service's client or
  `internal/lab` by any other path.
- **Every refusal is the server's.** A guard that concerns intent, a twin, a bundle or the
  host runs in the server, so every client is refused alike and under the same identifier.
  A client parses its own arguments, sends the request and renders the answer. "Before any
  connection" in the CLI's contract means the server's connections to Infrahub and to the
  workflow service.
- **The API's payloads are the documents Fylgja already versions**: the findings document a
  command prints under `--json`, with its blocks, and the CTM, the bundle and the support
  package where a command takes or returns one. The API adds no second representation of a
  twin, a plan, a step or a report.
- **The API exposes what the CLI needs and nothing else.** Each later client adds the
  endpoint it needs, in its own feature.
- **Test code and test tools are excepted, and nothing else is.** They may call the core
  directly, and may read the nodes, the containers and the record on their own: tier 1's
  and tier 2's Go tests, which call the packages they test; `fylgja-fixture`, which writes
  to Infrahub and is no client of Fylgja; the conformance suite's boot half, which reads
  the record and the nodes; and the end-to-end tier's and the hand checks' own reads and
  interventions (`gnmic` against a node, Docker and containerlab against the lab,
  `twin.json` and the store on disk, a stopped worker or workflow service). The commands
  they run are the CLI's, so they go through the API. A shipped client carries no second
  path for a test: no flag and no variable makes a command link the core.
- **A tier-1 test over imports holds the boundary**
  (`TestClientReachesTheCoreThroughTheAPIAlone`), as `internal/compiler/purity_test.go`
  holds the compiler's purity. Go imports are per package, and the one binary also carries
  the roles that are the core, `serve` and `worker run` ([D-018](#d-018)). So the client's
  commands live in a package of their own, `internal/cli`, and the test holds it to
  `internal/api`, the findings document types and `internal/tree`, and to nothing of the
  core.
- **No command of the CLI runs the core in its own process.** The stage commands are
  requests, as `twin create` is. A stage command still reads and writes the operator's
  files: the client sends the CTM, the bundle or the support package a command takes, and
  writes the CTM or the bundle it returns.

### Consequences

- **A client holds one credential: the API's own token, beside the API's address**
  ([D-042](#d-042)). Infrahub's token and the node logins stay with the processes that are
  the core.
- **`fylgja serve` listens on loopback unless told otherwise, and terminates no TLS.** A
  key file would be the first credential Fylgja reads from a file. A client on another
  machine reaches the server through an SSH tunnel or a proxy the operator runs.
- **The CLI and the server can differ.** A client on another machine can run another
  build. The API is versioned in its path, the server reports its build, and the client
  refuses an API version it does not know and warns when the builds differ
  ([D-018](#d-018)).
- **The stage commands need a server**, and the packages they are checked against are the
  server's ([D-019](#d-019)).
- **A run is two requests.** The command's own request starts the run and stays open as its
  progress: it answers with the run's identity as soon as there is one, then the step
  events, then the document. The start and the progress share a request because the
  closing report is built from what the server decided before the start (the run's input,
  the command's own findings, the subject), which only that request's handler holds while
  the server keeps nothing, and because a progress stream apart from its start would be a
  re-attach to a run, which the API does not expose. The operator's interrupt is a second
  request, delivered to the open one: the first cancels the run, or abandons a start still
  waiting for a worker, and the second ends the answer. A client that goes away, or a
  server that stops, leaves the run running ([D-014](#d-014) says what the server holds
  for it).
- **The server renders on request.** By default an answer is the findings document, and a
  run's progress is its events; a client may ask for the text as well, and the CLI does,
  printing the lines as received.
- **A client words what happened to itself**: an unreachable server, a refused token, an
  unknown version and a transfer over the bound (`api.unreachable`, `api.token.refused`,
  `api.version.unknown`, `api.transfer.too_large`), and a file on its own disk that it
  cannot read or write. None concerns intent, a twin, a bundle or the host, which stay the
  server's to refuse. A command that ran is answered as a success of the transport,
  carrying its document; the transport's error statuses are for its own faults.
- **No new lock.** The fixed workflow ids stay the only guard against two operations at
  once ([D-007](#d-007), [D-013](#d-013)), now that two clients can ask at once.
- **The CLI's tests keep what they assert.** They run each command against a server in the
  test's process over faked services, so a wording that moves is a defect, not a
  re-baseline. The end-to-end tier runs its commands with the API's address and token in
  their environment and nothing else of Fylgja's.
- **Every later client is bound by this entry**: the web UI, the agent surface's MCP
  server and the checks. [D-026](#d-026)'s webhook has its listener.
- Constitution XI states this rule, and principle I names the binary's three roles; the
  server is Go's standard library, so there is still no second runtime.
- Revisit if a client needs an operation the server cannot offer without holding state:
  that is [D-014](#d-014)'s question before it is this entry's.

---

<a id="d-041"></a>
## D-041 — `fylgja serve` runs on the lab host, as a role of its own, and reads the host in its own process until the host move

**Status:** Accepted 2026-10-04

### Context

[D-040](#d-040) puts every command behind the API. Most commands read the lab host before
they start a run, or instead of starting one:

- `twin show` and `twin verify` read containerlab, `twin.json` and the staged manifest, and
  verify reads each node over its package's transport with that package's login.
- Every dry run and `twin step` run the host check, and a step reads containerlab's plan
  with `CLAB_LABDIR_BASE` at the twin directory.
- `twin step`, `twin provision`, `waypoint plan` and every dry run file a bundle in the
  store. A step and a provision then hand the run that bundle's local path, which the
  worker opens. That is the seam the roadmap's M8 names.

Two promises depend on where the server runs. `twin show` reports the host and the record
with the workflow service unreachable, within five seconds. `twin verify` needs no run and
no worker ([D-039](#d-039)).

### Options considered

- **Anywhere, with every host read behind the task queue.** The server would need no
  containerlab and no node login. Rejected: a client of the workflow service cannot ask a
  worker for an activity without a workflow, so each read would become a short workflow
  (show, verify, the host check, containerlab's plan and the record); `twin show` could no
  longer answer with the workflow service down, nor `twin verify` with no worker; the
  filing for a step, a provision and a plan would have to move into activities, or M8's
  object store would have to come first; and it designs remote reads before M9's runtime is
  chosen (an [open question](#open-questions)).
- **On the lab host, reading the host in its own process.** Chosen.
- **One process for the server and the worker.** One thing to start, one environment, and
  no skew between the two. Rejected: a restart of the API, to change its token or where it
  listens, would cut every activity in flight, and a fault in a handler would take a
  deploy's heartbeats with it.
- **Two roles of the one binary, two processes.** Chosen.

### Decision

`fylgja serve` runs on the lab host, beside `fylgja worker run`, as a process of its own.
It shares the worker's state root. It reads containerlab, the twin directory, the bundle
store and the nodes in its own process, through `internal/lab`, and starts runs through the
workflow service. This holds until the host move (M9), which decides how a server reaches a
host it does not run on.

### Consequences

- `twin show` answers for the host and the record with the workflow service unreachable,
  and `twin verify` needs no run and no worker. No workflow, activity or wire type exists
  for a read.
- **The server's environment is the worker's and one more.** It holds Infrahub's address
  and token, both node logins, the workflow service's address, the state root and the
  override directory, and the API's own token. `docker` and `containerlab` are on its
  `PATH`.
- **The server and the worker share one filesystem.** A bundle's path crosses the task
  queue as a local path until M8.
- **A client on another machine drives the lab host through the API**, which is the
  operator's access into a remote twin. Whether `twin show` and the dry run move behind the
  queue is M9's to decide, with its runtime.
- **The server's host reads sit behind one seam**, so that M9 changes what is behind it and
  not the handlers.
- Three processes run on a lab host: the workflow service, the worker and the server.
- One twin; the server names nothing per twin (Constitution VII). No new workflow, signal
  or queue, and host-bound code stays in `internal/lab` (VIII). The server keeps no state
  (IX). It takes its credentials from its environment alone and puts none in a document or
  a log line (X).
- Revisit at M9, or sooner if the server must run where containerlab does not.

---

<a id="d-042"></a>
## D-042 — The API's token is one shared secret from the process environment, sent as a bearer credential and compared before anything else

**Status:** Accepted 2026-10-04

### Context

[D-040](#d-040) makes the API the boundary of the core, and [D-041](#d-041) puts its server
on the lab host with every credential the core has: Infrahub's token, both node logins and
access to the workflow service. Whoever can send the server a request can read intent,
build a twin and destroy it. So the API needs a credential of its own. It is the first
credential Fylgja checks itself; the others are Infrahub's and the nodes' to check.

Constitution X binds it: a credential comes from the process environment only, is never
persisted, and appears in no bundle, finding or log line.

### Options considered

- **No credential on loopback.** Rejected: every user and process on the lab host could
  drive the twin, and an SSH tunnel or a proxy carries loopback to wherever the operator
  points it.
- **One shared token from the process environment.** Chosen.
- **Roles, or several tokens**: a read-only token for `twin show` and `twin verify`, or one
  per client. Rejected by the operator: there is one operator and one client, and a role
  model drawn before a second client exists is a guess.
- **A token file**, named by a flag. Rejected: it would be the first credential Fylgja reads
  from a file.
- **A token the server makes at start and prints.** Rejected: it prints a credential, and a
  restart would change it under every client.
- **The server terminating TLS, or mutual TLS.** Rejected by the operator: the key is a
  credential in a file.
- **The token in the address or in the request's body.** Rejected: a proxy logs addresses,
  and the body must not be read before the token is checked.
- **HTTP Basic.** Rejected: there is no user to name. A bearer credential is the form every
  HTTP client already has.
- **A Unix socket, guarded by file permissions, in place of a token.** Rejected: a client's
  SSH tunnel and a later client in a container both need an address that TCP reaches.

### Decision

- **One token.** `FYLGJA_API_TOKEN`, read from the process environment by the server at
  start and by a client at each command. It has no default, and any non-empty value a
  request header can carry is accepted: its strength is the operator's.
- **The server refuses to start without it**, naming the variable, exit 2, before it
  listens.
- **It is sent as `Authorization: Bearer <token>` and compared on every request**, in
  constant time, before the route is looked up, the body read, a guard run or anything
  dialled. A request without it, or with another, is `401`.
- **It is never persisted and never printed**: in no document, finding, frame, problem
  body, log line, error or start-up report. A client's finding names the variable and
  never its value.
- **The server listens on loopback unless `--listen` says otherwise**, at
  `127.0.0.1:7650`, and terminates no TLS. Told otherwise, its start-up report says the
  token is sent unencrypted on that address. No variable sets the server's address.
- **A client whose variable is unset, or holds a value no request header can carry (a
  control character, or a trailing space or tab), refuses without connecting**, under the
  same identifier as a refused token (`api.token.refused`).

### Consequences

- The token is kept with the other credentials the operator loads into the environment
  (`local/.env` on the development host), and [.env.example](../.env.example) names it.
- The token crosses in the clear. Keeping it confidential is the transport's work, and the
  transport is the operator's: loopback, an SSH tunnel, or a proxy that terminates TLS,
  which a client addresses by an `https://` URL. A client follows no redirect, so the token
  is never sent to an address the operator did not give.
- Rotating it is a restart of the server with another value. No run is touched, since the
  server and the worker are two processes ([D-041](#d-041)).
- A request that carries a file sends `Expect: 100-continue`, so a refused token costs no
  upload.
- The end-to-end tier makes a token for each run and greps every output, record, history
  and log for it, the server's log included.
- The token grants everything the API does. There is no read-only client until a new entry
  says what one is. It is what closes the boundary of Constitution XI to whoever is not a
  client.
- Revisit when a second client needs less than everything (the web UI, the check inside
  Infrahub), or when a server has to listen off loopback as a matter of course: that needs
  TLS, and a decision about where its key lives under Constitution X.

---

<a id="d-043"></a>
## D-043 — The development record stays private: this repository began as one commit of cleaned code and rewritten documents, and development continues here

**Status:** Accepted 2026-10-05

### Context

Fylgja was built under the Spec Kit workflow in a private repository. Besides the code, that
repository holds the record of how it was built: its history, each feature's
specification, plan, research log and tasks, the command log of every session, the working
notes of each pass, and earlier forms of this log and the roadmap, with every correction
and estimate. Those are records of how the project was built, not documents a reader needs
to understand or stand up what was built, and they cite commit hashes, a development
host's process ids and local paths. Publishing them would also mean auditing every commit
for what it carried.

### Options considered

- **Make the private repository public.** Rejected: the first public view would be the
  build records, not the system, and the whole history would need an audit.
- **Filter its history** (`git filter-repo`) and publish the result. Rejected: it keeps the
  provenance, but still needs every commit audited.
- **Keep developing in private and copy cleaned snapshots out.** Rejected: two trees to
  keep in step.
- **Copy cleaned code and rewritten documents into a new repository in one commit, and move
  development there.** Chosen.

### Decision

- **This repository began with one commit**: the code, tests and test data, and the
  documents rewritten for its reader. No commit of the private repository crossed.
- **Development continues here.** Spec Kit, the converge loop and the command log run here.
  `specs/` holds the features specified here; their numbers continue the private record's,
  so the first is `013`.
- **The development record stays private**: the history, the features specified before this
  repository, their research and command logs, and the earlier forms of this log and the
  roadmap. A write-up of how the project was built summarises them and does not cite them.
- **The bring-up script and CI are built here**, after the first commit.

### Consequences

- **The contracts live outside `specs/`**, in [contracts/](../contracts/), which the tests
  read. A test that validates an older format against an older schema reads a frozen copy
  under its package's `testdata/contracts/`.
- **No code comment or document cites the private record**: a citation of a requirement, a
  research run or a task became the fact it cited, or went. A bare milestone marker (`M7`)
  stays, and `CLAUDE.md`'s milestone table and the roadmap's Built table resolve it.
- **The decision log is append-only in each repository.** This log is a new document: it
  states each decision in force, its rejected alternatives and its reason, without the
  private log's supersession chains and notes, and begins at this repository's first commit
  (the head says so). Constitution 1.1.1's Documents gate says the same.
- **The artifacts template lives here** ([D-044](#d-044)), so a reader can have Infrahub
  render the configuration Fylgja pushes with no credential.
- **The tooling crossed with the code**: `.specify/`, `.claude/` (its hooks and skills) and
  `scripts/`. The command log is git-ignored here.
- M14 is two features: the cut, which made this repository's first commit, specified and
  built in the private record; and the launch, feature `013-launch`, the first feature here,
  which builds the bring-up script and CI and makes the repository public
  ([roadmap](roadmap.md)).

---

<a id="d-044"></a>
## D-044 — The artifacts template lives in this repository, and Infrahub registers it read-only on a ref with no credential

**Status:** Accepted 2026-10-06

### Context

[D-028](#d-028) puts the rendering of a device's configuration in Infrahub: one Jinja2
transformation over the Fylgja generics, one query, one artifact definition
(`device-config`, `text/plain`, one per device in the group `fylgja-devices`). They lived in
a private repository, registered as a read-write `CoreRepository` with a
`CorePasswordCredential` holding a GitHub personal access token. A reader of this
repository cannot stand the system up without the template: without a rendered
`device-config`, every read refuses (`artifact.missing`), the fixture tool's seed refuses,
and tier 2 and the end-to-end tier cannot run.

Verified on 2026-10-06 against Infrahub 1.11.2 ([verified-facts.md](verified-facts.md)):
both repository kinds are branch-agnostic, inherit `CoreGenericRepository`, and declare
`credential` as an **optional** relationship to `CoreCredential`. A
`CoreReadOnlyRepository` created on `main` for a public HTTPS repository with no credential
reached `internal_status active` with the remote's head commit in 14.5s and
`operational_status online` in 34.8s; its import failed only because that repository had
no `.infrahub.yml`. So the credential-less clone and connectivity check pass. A
`CoreRepository` (read-write) mirrors the remote's branches as Infrahub branches and pushes
worktrees back for Infrahub branches created with `sync_with_git` on; a
`CoreReadOnlyRepository` tracks one `ref`, never pushes, imports at creation, and takes a
new commit by a `ref` update or by re-creation. The fixture tool looks the definition and
the group up by name and never names the repository object, so a registration under another
name serves it.

### Options considered

- **A second public repository for the template**, cut the same way. Rejected: two things to
  clone, read and register, for a template that renders from the schema in the repository
  beside it, and a second place for the reader to be told about.
- **Keep the template private and hand the reader the credential path.** Rejected: it keeps
  every reader out of the one thing they need to render a configuration.
- **The template here, registered as a `CoreRepository` with no credential.** Rejected: a
  read-write registration mirrors every branch of a public repository into Infrahub, and
  pushes for any Infrahub branch created with `sync_with_git` on, which no credential
  allows; the kind's point is writing back, and nothing is written back.
- **The template here, registered as a `CoreReadOnlyRepository` on a ref with no
  credential.** Chosen.

### Decision

- **The template, its query and `.infrahub.yml` are files of this repository**:
  `.infrahub.yml` at the root, where Infrahub requires it, and `infrahub/queries/device_config.gql`
  and `infrahub/templates/device_config.j2`. The transformation is named
  `srlinux_device_config` (the one definition renders both platforms, [D-028](#d-028)) and
  the artifact is `device-config`.
- **Infrahub registers this repository as a `CoreReadOnlyRepository` on `main` with no
  credential.** With the schema on `main` and the group `fylgja-devices`, that is what a
  reader does to render the configuration Fylgja pushes, and
  [development.md](development.md) describes it.
- **A private copy registers as a `CoreRepository` with a `CorePasswordCredential`**, which
  development.md names as the private alternative.

### Consequences

- A reader needs no GitHub credential to render: the repository, the schema and the group
  are the whole of it.
- The template's bytes are those the private repository rendered with, so registering this
  repository in its place moves no artifact and no fixture id (`b9d53ebc…`).
- The fixture tool's error "the fylgja-artifacts repository must be connected and in-sync on
  main" still names the private repository; it is reworded when this repository is
  registered, since nothing but the message changes.
- A read-only registration imports at creation. Whether Infrahub 1.11.2 pulls a new commit
  on a `ref` update, or only on re-creation, is verified when this repository is
  registered; until then a template change is a documented re-registration, not a push
  Infrahub notices.
- The template reads `FylgjaDevice` and its generics alone, so where it lives changes
  nothing it queries (Constitution III); Infrahub renders and Fylgja pushes the bytes it
  fetched (IV).

---

<a id="d-045"></a>
## D-045 — A stage is a moment of a proposed change: the branch point read on `main`, the operator's marks, and the head pinned one microsecond after the branch's latest stamp

**Status:** Accepted 2026-10-16. Supersedes [D-032](#d-032). Built by
`014-proposed-change-stages`; until it is, [architecture.md](architecture.md) §3's
Waypoints paragraph and the glossary describe what runs.

### Context

[D-032](#d-032) gave the operator a name for a pinned reference: a `FylgjaWaypoint`, in a
series, naming a branch by free text. The operator found it clumsy on two counts. It is
not Infrahub-native: nothing in Infrahub knows a waypoint, and its own view shows none. And
it is separate from the change: a branch-agnostic side object that names its branch by
text, outlives it, and is invisible where the schema lacks the kind.

Infrahub 1.11.2 has no named point in time — no branch tag, checkpoint or saved
time-travel view; the time selector writes a URL parameter and saves nothing — and
`Branch` is not a kind anything can relate to. Its one unit of change that a user kind
may relate to is the **proposed change**, which can be opened as a draft while the change
is built. So a stage is a moment of a proposed change.

Verified on 2026-10-16 against 1.11.2 ([verified-facts.md](verified-facts.md);
`specs/014-proposed-change-stages/research.md` §2), on throwaway branches and then on
`main`:

- A branch-agnostic `FylgjaStage` with a required relationship to `CoreProposedChange`
  and an extension giving the proposed change its `fylgja_stages` loads on a branch and on
  `main`; the uniqueness constraint `[proposed_change, label__value]` refuses a duplicate
  (`Violates uniqueness constraint 'proposed_change-label'`) and accepts the same label on
  another proposed change. A human-friendly id through the proposed change is refused,
  since two proposed changes may share a name.
- One open proposed change per source branch; `main` as source refused. Closing one
  frees the branch for another; deleting the branch cancels it about 16 s later, and it
  still names the branch. Deleting a proposed change succeeds under its stages, which
  are left pointing at nothing.
- A read pinned at exactly `branched_from` answers with `main` as it then stood, back to
  the default branch's creation, and is refused before it.
- **A read pinned at exactly an attribute's `updated_at` sees the node but not the
  value**; one microsecond later sees it. The boundary is a strict comparison, not a
  truncated stamp.
- **A rebase re-stamps the branch's whole content at the new `branched_from`** and
  discards its earlier history from pinned reads; the branch at `branched_from` then
  reads the branch's own changes, not the base.
- The proposed change's pipeline **regenerates the branch's artifacts at its creation,
  draft or not, and on a run-check**, and not on a branch update or the end of a draft.
  A regeneration moves `checksum.updated_at` and `storage_id.updated_at` together and
  leaves `status.updated_at` where it was.
- The diff tree is a cached calculation, refreshed by the pipeline or the `DiffUpdate`
  mutation, so neither the head nor `main`'s movement is read from it without a write.

### Options considered

Where the mark lives:

- **A branch tag or checkpoint.** None exists in 1.11.2, and no issue asks for one.
- **A tag on intent**: each object marked with the stage it goes live in, stage N
  compiled from the head by filtering. Rejected: it expresses additions only, never a port
  disabled or a description changed in stage 2, and the artifact for a partial topology
  does not exist — Infrahub renders the head, and Fylgja never renders ([D-009](#d-009),
  [D-028](#d-028)).
- **A stack of proposed changes**, one per step. Rejected: Infrahub branches from `main`
  alone, so a second step cannot be expressed until the first merges; the two derived
  ends already rehearse "`main` plus this change".
- **A side object naming a branch by text**, as D-032. Rejected: the two counts above.
- **A kind related to the proposed change.** Chosen.

The ends:

- **Marked, by the operator.** Rejected: every proposed change has a branch point and a
  head whether or not anyone marks them.
- **Derived.** Chosen: the branch point and the head are stages of every proposed
  change, and only the moments between are written.

The head's pin:

- **The branch's latest stamp, verbatim**, as the design first said. Rejected: a read
  pinned at exactly that stamp reads a half-visible write — the node without the value.
- **Unpinned**: the head read live and recorded by `observed_at`. Rejected: every other
  stage is pinned, and a twin from the head must be reproducible until the branch moves.
- **The latest stamp plus one microsecond**, the smallest increment Fylgja can send.
  Chosen.

The branch point's read:

- **On the branch at `branched_from`.** Rejected: correct for a fresh branch and wrong
  after a rebase, where it reads the branch's re-stamped content, half-visible.
- **On `main` at `branched_from`.** Chosen: the branch point by definition, before and
  after a rebase, and `main` was not written at that instant.

A stale stage, one marked after an intent change and before a regenerate:

- **Left to the operator**, as D-032's "write the waypoint last". Rejected: 015's check
  would fail a merge on a stage whose configuration is merely stale.
- **Inferred around**, by regenerating or by reading the artifact at a later time.
  Rejected: Fylgja writes nothing (Constitution III), and a later artifact is another
  stage's.
- **Refused**, naming the device and the remedy. Chosen.

Waypoints kept beside stages: rejected; Fylgja is 0.x, and two ways to name a pinned
reference would split every command.

### Decision

- **A stage is a moment of a proposed change.** The kind is `FylgjaStage`, loaded on
  `main` beside the contract ([D-028](#d-028)), branch-agnostic, with a required
  relationship to its `CoreProposedChange`, a `label` unique within it, an optional `as_of`
  and an optional description, and no branch of its own: its branch is the proposed
  change's source. An extension gives `CoreProposedChange` its `fylgja_stages`.
- **The operator writes stages; Fylgja only reads them** (Constitution III), by
  Infrahub's means: the UI, an object file, `infrahubctl`. A proposed change is found by
  its source branch, since Infrahub allows one open proposed change per branch.
- **A proposed change's stages, in order**: the branch point, read on **`main`** pinned
  at the branch's `branched_from`; each `FylgjaStage` in order of its `at` — `as_of` as
  written, or the stamp Infrahub returned for its writing; and the head, read on the
  branch pinned at **the latest `updated_at` across its intent (attributes and
  relationship edges) and its artifacts' checksums, plus one microsecond**. Every stage is
  pinned, and Fylgja composes no `at` from a clock of its own: each is a stamp Infrahub
  returned, or the first instant after one. Constitution VI says so from 1.2.0.
- **A stage is refused** when its `at` is before the branch point or later than now; when
  the branch's `branched_from` has moved past the stage's stamp while its `created_at` has
  not, the refusal says a rebase did it, since a rebase refuses every mark at once.
- **A stale stage is refused.** At each stage Fylgja reads, pinned at the stage's `at`, the
  latest intent stamp on the branch and each artifact's `checksum.updated_at`; intent
  newer than its artifact is `stage.artifact.stale`, naming the device and saying to
  regenerate, or run the proposed change's checks, and mark again.
- **It warns when `main` has moved since the branch point**, read as the head is: `main`'s
  latest stamp against `branched_from`.
- **`--stage` replaces `--waypoint`** on `twin create`, `twin step` and `intent read`;
  `stage list` and `stage plan` replace the waypoint commands, naming the proposed change
  by its source branch; `twin step` goes to a stage of the twin's own proposed change
  ([D-036](#d-036)'s target). The findings' `waypoint.*` rules become `stage.*`. The name
  is still not provenance: it enters neither the CTM nor the bundle, and the twin's
  record names it ([D-014](#d-014)).
- **A story along `main`'s own history keeps no named marks**: it is `--branch main --at`,
  as [D-012](#d-012) has it, since a proposed change cannot have `main` as its source.
- **`FylgjaWaypoint` is removed**, with its commands, rules, contract and test support.

### Consequences

- Constitution VI is amended (1.2.0): "MUST NOT send an `at` it composed from a clock of
  its own; one read from Infrahub, or the first instant after one, is passed as the
  operator's would be". D-032's reading, that a stamp read is not composed, is now the
  text, and the microsecond is the one addition.
- The twin record, the API and the CLI move for stages while the contract stays `0.2`;
  the goldens, the fixture and its id do not move. `main`'s schema hash moved with the
  kind (`0aba31a5…` → `e7386419…` on this install); a branch's hash is its own, so the
  fixture branch keeps `0aba31a5…` and compiles to `47b2c449…` as before.
- A test that makes a proposed change deletes its stages first, then the proposed
  change, then the branch; a branch deleted under an open proposed change leaves it
  canceled, naming the branch. A stage whose proposed change is gone is a finding, not a
  crash.
- A stage marked before a regenerate is the case the draft workflow invites, since the
  pipeline runs at creation and on a run-check and not as the branch is edited; the
  finding's remedy names both.
- A rebase ends every twin's story on that proposed change until the operator marks
  again; the record keeps what the twin was built from, as [D-036](#d-036) has it.
- A stage, like a waypoint, needs no series: one proposed change, one story. Two
  proposed changes cannot be stepped between, since a twin steps within its own.
- Infrahub's proposed-change view does not show stages. 015's comment lists them on the
  proposed change, and the upstream request for named points in time, drafted once the
  pass works, asks Infrahub for a native mark.
- Revisit if Infrahub gains a named point in time or a rebase that keeps history (both
  would retire the microsecond and the branch point's read on `main`), if the strict
  boundary of §2.6 moves on an upgrade (re-verify it on every upgrade), or if a second
  proposed change on one branch becomes expressible.

<a id="d-046"></a>
## D-046 — A stage's generation stamp is the later of its artifact's checksum stamp and the proposed change's artifact validation; staleness is judged per device

**Status:** Accepted 2026-10-16. Refines [D-045](#d-045). Built by
`014-proposed-change-stages`; superseded by [D-058](#d-058) on 2026-10-26 in what a device's
intent stamp is.

### Context

[D-045](#d-045) refuses a stale stage: at each stage, pinned at its `at`, intent newer
than its artifact's `checksum.updated_at`. Verified at the plan on 1.11.2 (2026-10-16,
`specs/014-proposed-change-stages/research.md` §6.5): a regeneration to identical bytes
writes no checksum and no storage id, on a hand generate and on the proposed change's
run-check alike, while the run-check re-runs the proposed change's artifact validator and
moves its `completed_at`, identical bytes or not. So under D-045's rule as written an
intent edit the template does not render (a device's `site`) would leave a stage stale with
no remedy that clears it, and 015's check would fail a merge on a harmless edit.

### Options considered

- **D-045's rule as written.** The hole above.
- **Branch-wide intent against the artifacts.** A hand generate after a change leaves the
  unchanged devices' bytes identical and their checksum stamps old, so one added link
  would make every untouched device stale for ever after.
- **A warning instead of a refusal.** Contradicts D-045 and its reason: a rehearsal of a
  stage whose configuration is not the intent's proves nothing.
- **The artifact's node metadata as the generation stamp.** Whether an identical
  regeneration touches it is unverified; the validator's completion is.

### Decision

A device is stale at a stage when the latest stamp of its own intent, its attributes and
relationship edges, read pinned at the stage's `at`, is later than its generation stamp:
the later of its artifact's `checksum.updated_at` (the artifact of its package's
`config.artifact_name`, read on the branch pinned at `at`) and the proposed change's
artifact validator's `completed_at` value (read on the default branch pinned at `at`;
absent when none has completed). The comparison is per device: exactly the devices whose
intent moved are named. The one case a per-device reading misses, a neighbour's change
that alters this device's rendered lines while none of its own stamps moves, is covered by
the run-check, which re-stamps the validator for every device. `stage.artifact.stale`
names the device and both remedies: a generate, or the proposed change's checks. Nothing
else in D-045 moves.

### Consequences

- The stage reader reads `CoreArtifactValidator` pinned, beside `CoreProposedChange` and
  `Branch` (Constitution III as amended, 1.2.1).
- Tier 2 proves the hole closed: a non-rendered intent change, a stage refused stale, a
  generate that moves nothing, a run-check, a fresh stage planned.
- Revisit if an upgrade makes an identical regeneration move the checksum's stamp, or
  gives an artifact a generation stamp of its own.

<a id="d-047"></a>
## D-047 — A stage is orphaned when its proposed change no longer names it; Infrahub refuses the delete that D-045 gave as the cause

**Status:** Accepted 2026-10-17. Refines [D-045](#d-045). Built by
`014-proposed-change-stages`.

### Context

[D-045](#d-045)'s context says a proposed change's delete succeeds under its stages and
leaves them pointing at nothing, and the feature's spec, `stage.orphaned`'s message and the
data model give that delete as the orphan's cause. That was verified with the stage kind
loaded on a throwaway branch alone. With the kind on `main`, as it ships, Infrahub 1.11.2
refuses the delete while a stage names it (2026-10-17,
`specs/014-proposed-change-stages/research.md` §8.8). An orphan arises only when the stage
leaves the proposed change's `fylgja_stages` (`RelationshipRemove`), after which the stage
reads `proposed_change { node: null }` despite `optional: false`, and the proposed change
can then be deleted.

### Options considered

- **Keep the deletion wording.** It names a cause Infrahub refuses, so an operator reading
  the finding looks for a delete that never happened.
- **Name both causes** ("no proposed change names it", covering a delete a later Infrahub
  might allow). Rejected on the operator's word: it describes a path not reachable on the
  version Fylgja verifies against.
- **Name the removal.** Chosen.

### Decision

A stage is orphaned when its proposed change no longer names it: it was removed from the
proposed change's `fylgja_stages`. `stage.orphaned`'s message, the spec's edge case and US3
AC3, the data model and the code's comments say so, and none gives a deleted proposed
change as the cause. What an orphan does is unchanged: a warning in `stage list`, planned
and created from by nothing, and a reference naming it refused as unknown. Nothing else in
D-045 moves.

### Consequences

- The message and its pins change together (`internal/stages/list.go`, its tests,
  `internal/cli/stage_list_test.go`, `internal/findings/stages_test.go`), with
  `contracts/stage-reference.md`'s quoted wording.
- An orphan names no branch, so it outlives its branch's cleanup until it is deleted by
  id (research §8.8); a test that makes one deletes it itself.
- Revisit on an Infrahub upgrade that lets a proposed change be deleted under its stages.

<a id="d-048"></a>
## D-048 — The head and the staleness rule also read the branch's intent events, so a removal that no stamp records is seen

**Status:** Accepted 2026-10-17. Refines [D-045](#d-045) and [D-046](#d-046); [D-026](#d-026)
stands. Built by `014-proposed-change-stages` (T078); until it is,
[architecture.md](architecture.md) §3's Stages paragraph describes what runs.

### Context

[D-045](#d-045) pins the head one microsecond after the branch's latest stamp, and
[D-046](#d-046) judges a device stale when the latest stamp of its own intent is later than
its generation stamp. Both read the stamps of what is present. Verified on Infrahub 1.11.2
(2026-10-17, `specs/014-proposed-change-stages/research.md` §8.17): deleting a link as the
branch's latest write moves no stamp the stamps query reads. So the head is pinned before
the delete and still reads the link, a stage marked after it is not refused though its
artifacts still configure the link, and `stage plan`'s last step puts the link back. Only a
regenerate, which moves the endpoint devices' checksums, hides it. Three readers were tried
read-only (§8.18):

### Options considered

- **The diff tree** (`DiffTreeSummary`, `DiffTree`). It is a cached calculation. After the
  delete only `num_untracked_diff_changes` moved, with no stamp. Refreshing it needs
  `DiffUpdate`, a write to Infrahub. Rejected, as D-026 rejected it.
- **An unpinned read compared with the head's.** It needs no new root, and it differs by
  exactly the deleted link. But it carries no stamp, so it can warn that a head misses a
  change and cannot place a head after it. It also costs a second read per head.
- **Recording the limitation**, with a regenerate before marking after a removal. Rejected
  on the operator's word: the head is the stage a rehearsal steps to last, and 015's check
  would pass a change whose removal its twin never saw.
- **Infrahub's event log**, `InfrahubEvent`. Chosen. It is a read with no write.
  `branches: [<branch>]` with `since: <branched_from>` bounds it to the branch's own
  history, so no cursor is kept between reads, which was D-026's objection to it for
  following. The delete appears as `infrahub.node.deleted`, its `relationships` naming the
  removed peers, with child `node.updated` events on each endpoint interface. A read pinned
  one microsecond after the event's `occurred_at` already reads the removal.

### Decision

The stage reader also reads the source branch's events since `branched_from`, through
`InfrahubEvent`, pinned to nothing, paged at its limit of 50, newest first, and stopping
once the events are older than what the read needs. An event is **intent's** when its
primary node's kind is `used_by` a Fylgja generic in the branch's schema, mapped as
conformance maps a kind, at any level. Infrahub's own churn (query groups, artifacts,
memberships) is never intent. A concrete kind is read only to be mapped this way and is
never queried.

- **The head** is pinned one microsecond after the later of the branch's latest stamp
  (D-045) and the `occurred_at` of its latest intent event.
- **A device is stale** at a stage (D-046) also when an intent event later than its
  generation stamp, and no later than the stage's `at`, names it. An event names a device
  when the device's id, or the id of one of its interfaces as read at the stage, is the
  event's primary node or a peer in its `relationships`. Its message names the event's
  action and kind, beside the stamp's.
- `stage.main.moved`, the branch point and every mark's `at` are as they were. Following
  is as D-026 decides: no event starts a check, and none decides whether one reads.

### Consequences

- Constitution III is amended (1.3.0): `InfrahubEvent` becomes the fourth Infrahub-owned
  root the reader may read, named in `internal/intent/readonly_test.go` with its reason.
  An event is Infrahub's record of a write, not the operator's model.
- Each head and each staleness judgement costs a page or more of events (0.11–0.29 s a
  page on the development host). A branch with a long history pages further. The stop
  condition keeps a head's read to the pages newer than its latest stamp.
- Tier 2 proves it: a link deleted as the branch's latest write moves the head past the
  delete, a stage marked after it is refused stale for both endpoint devices and no
  other, and a regenerate clears it.
- Revisit if an upgrade stamps a removal in the graph (then the event read can go), if
  `InfrahubEvent` changes its shape or its retention, or if a removal on `main` must move
  `stage.main.moved`.

---

<a id="d-049"></a>
## D-049 — A check posts the rehearsal's report as a comment; Fylgja writes nothing

**Status:** Accepted 2026-10-19; superseded by [D-053](#d-053) on 2026-10-20 in the kind it
names. Beside [D-026](#d-026), which stands. To be built by `015-proposed-change-rehearsal`. Superseded by [D-065](#d-065) on 2026-10-29 in "the check posts" and "Fylgja writes nothing"; the thread stands.

### Context

F-007 and F-008 bring the rehearsal inside Infrahub's proposed-change workflow, whose unit of
writing is a comment rendered as markdown, while Constitution III keeps every operation
Fylgja makes on Infrahub a query. Verified on Infrahub 1.11.2 (2026-10-17,
`specs/015-proposed-change-rehearsal/research.md` §1.3): a check runs in Infrahub's task
worker with the source branch and the worker's own client, may create a `CoreChangeComment`
on the proposed change, and its own message is plain text.

### Options considered

- **Fylgja posting the comment through the API's server.** Rejected: a write to Infrahub,
  against Constitution III.
- **The report in the check's message alone.** Rejected: a check's message renders as code,
  not markdown.
- **One comment edited on each run.** Rejected: Infrahub's comment carries no stamp, and a
  history of rehearsals should stay readable.

### Decision

The two checks this repository ships, `fylgja-step` and `fylgja-twin`, are clients of the
API and glue. `fylgja-twin` posts one new `CoreChangeComment` per run, carrying the report's
markdown headed by the run and when the stages were resolved; it edits and deletes none, and
writes nothing else. `fylgja-step` writes nothing. A comment that cannot be posted is logged
and fails nothing.

### Consequences

- A proposed change's comments are its rehearsals' history.
- A test's proposed change goes with its comments when deleted (research §1.6).
- D-026 stands: Fylgja itself still writes nothing to Infrahub.

---

<a id="d-050"></a>
## D-050 — The API's server listens on the host's Docker bridge beside loopback, in plaintext, under the one token

**Status:** Accepted 2026-10-19. Supersedes the loopback default of [D-042](#d-042) for this
deployment; D-042's token and its comparison stand. To be built by
`015-proposed-change-rehearsal`. Superseded by [D-065](#d-065) on 2026-10-29 in the listener it puts on the bridge, which is the orchestrator's platform API's; the core's server listens on loopback alone.

### Context

Infrahub's task workers run on a Compose bridge where neither loopback nor
`host.docker.internal` reaches the lab host (`specs/015-proposed-change-rehearsal/research.md`
§1.1). Verified (2026-10-17): with `extra_hosts: host.docker.internal:host-gateway`, a
container on that network resolves the name to the host's docker0 address (`172.17.0.1` by
Docker's default), reaches a server bound there with the token, and is refused without it.

### Options considered

- **A relay container.** Rejected by the feature's brief: a second process to run and keep.
- **Server TLS.** Rejected by the brief and by D-042: a key file would be the first
  credential read from a file.
- **The task workers on the host's network.** Rejected: Infrahub's Compose file stays
  Infrahub's.

### Decision

`fylgja serve --listen` is repeatable. On the lab host the server listens on
`127.0.0.1:7650` and on the bridge's address, `<docker0>:7650`. The Compose override gives
the task workers the host-gateway entry, `FYLGJA_API_ADDRESS` and `FYLGJA_API_TOKEN`, and the
bring-up script starts the server so.

### Consequences

- Any container on the host can reach the API, the twin's nodes included, and the bearer
  check is what guards it. `SECURITY.md` says so.
- The token is in the task workers' environment as well as in `local/.env`.
- D-042's revisit, for a client that needs less than everything, stays deferred.

---

<a id="d-051"></a>
## D-051 — `twin rehearse` is a fifth short workflow, `fylgja-rehearse`, that runs provision, step and destroy as children under their fixed ids, with a deadline derived from the stages' budgets

**Status:** Accepted 2026-10-19. Beside [D-007](#d-007) and [D-036](#d-036). Constitution
VIII is satisfied by this entry and its text is not amended: it admits a workflow through a
decision-log entry, as the step and following workflows were admitted. To be built by
`015-proposed-change-rehearsal`. Superseded by [D-065](#d-065) on 2026-10-29 whole.

### Context

A rehearsal is a create, a verify, a step and a verify per later stage, and a destroy, run
unattended for minutes, and it must destroy its twin whatever happens
(`specs/015-proposed-change-rehearsal/spec.md`, FR-005).

### Options considered

- **The server driving the runs one by one.** Rejected: the server keeps nothing between
  requests ([D-014](#d-014)), and [D-040](#d-040) leaves a run running when its client goes.
- **One workflow scheduling the activities itself.** Rejected: it repeats the three
  workflows' logic and their cleanup.
- **A long-running workflow.** Rejected by [D-007](#d-007).

### Decision

`fylgja-rehearse` has a fixed id. Its children are `fylgja-provision` (already compiled),
`fylgja-step` (allowing restarts, with the default wait) and `fylgja-destroy` on a
disconnected context, each under its own fixed id; its baseline wait is `VerifyTwin`. It is
refused beside any run in flight, and `twin destroy` cancels it and waits for its destroy.
Its deadline is the sum of its stages' provision, step and wait budgets, from the packages
and the mechanism constants, under a ceiling of Fylgja's own (3,600 s), enforced by a timer
that cancels the children and then destroys. Its twin's record names it, and a rehearsal
whose run is no longer running is destroyed by the next (the reclaim).

### Consequences

- The worker serves five workflows.
- The twin record moves to `7`; its rehearsal mark is the only new state.
- A step cut by the deadline records diverged ([D-036](#d-036)) before the destroy.
- `fylgja-twin`'s check timeout (4,200 s) exceeds the ceiling by the destroy's budget, the
  work before the start and a margin.

---

<a id="d-052"></a>
## D-052 — The check concludes from the report's status alone; Fylgja's verdicts stay advisory

**Status:** Accepted 2026-10-19. Beside [D-039](#d-039), which stands. To be built by
`015-proposed-change-rehearsal`.

### Context

A stage that fails must block the merge (F-008), while D-039 keeps every Fylgja report
advisory and no Fylgja operation reads one.

### Options considered

- **Fylgja refusing something on a failed stage.** Rejected: D-039.
- **The check reading the stages' verdicts.** Rejected: two places would decide what fails.

### Decision

`twin.rehearse`'s document carries the verdicts and the status, and the check maps the status
alone: `ok` is success and every other status is failure, with the report's text as the
message either way. A draft proposed change is skipped with success. `fylgja-step` fails only
on a refusal as intent, the stages unlistable, or Fylgja unreachable, a branch point refused
for want of a network on `main` excepted: the one rule the Python holds, read from the
finding's rule and `at_source`. A later `stage.plan` that carries the distinction itself
retires it.

### Consequences

- The merge is gated by Infrahub's own rule over the check's conclusion; nothing in Fylgja
  acts on the rehearsal.
- The roadmap's F-007 row is reworded: `fylgja-step` posts no comment.

---

<a id="d-053"></a>
## D-053 — The check's comment is a thread on the proposed change, as its view shows one

**Status:** Accepted 2026-10-20. Supersedes [D-049](#d-049) in the kind it names; the rest
of D-049 stands: the checks as clients of the API and glue, one new comment per run
carrying the report's markdown headed by the run, none edited, `fylgja-step` writing
nothing, a comment that cannot be posted failing nothing. [D-026](#d-026) stands. Built by
`015-proposed-change-rehearsal`. Superseded by [D-065](#d-065) on 2026-10-29 in "the check posts" and "Fylgja writes nothing"; the thread stands.

### Context

D-049 named `CoreChangeComment`, which a check may create, and the reviewer reads the
report in the proposed change's view (US2). Verified on Infrahub 1.11.2 (2026-10-20,
`specs/015-proposed-change-rehearsal/research.md` §5.9): the overview lists the proposed
change's events, and a comment appears there only in a thread, through the event a
thread's creation emits; a `CoreChangeComment` is counted in `total_comments` and shown
nowhere in the view; the operator read no comment in the overview of a proposed change
that held one.

### Options considered

- **Keep `CoreChangeComment`.** Rejected: no reviewer sees it.
- **An object or artifact thread.** Rejected: each is anchored to an object or an
  artifact's line, and the report is the whole proposed change's.
- **A `CoreChangeThread` holding one `CoreThreadComment`**, as the view's own *Add
  comment* writes one. Chosen.

### Decision

`fylgja-twin` posts each run's report as one new `CoreChangeThread` on the proposed change,
labelled `Conversation` and unresolved as the view labels its own, holding one
`CoreThreadComment` whose text is the report's markdown, both on `main`. It deletes the
thread it made when the comment cannot be saved, as the view does, and otherwise edits and
deletes nothing and writes nothing else.

### Consequences

- Each rehearsal is a thread in the proposed change's overview, in the order run.
- A reviewer may reply to it or resolve it, which nothing of Fylgja's reads.
- A test's proposed change goes with its threads and their comments when deleted
  (research §5.9).

---

<a id="d-054"></a>
## D-054 — For a demo, `main` carries the fixture's three-node network for a window the fixture tool opens and closes

**Status:** Accepted 2026-10-21. Beside [D-028](#d-028) and [D-044](#d-044), whose writes
to `main` stand, and [D-045](#d-045), whose branch point is read on `main` as before. To
be built for `v0.2.0`'s demo, in `cmd/fylgja-fixture`.

### Context

A rehearsal boots the first stage it can build ([D-051](#d-051)), and the first is the
branch point, read on `main` at the branch's `branched_from` ([D-045](#d-045)). `main`
holds no network on the development host or on CI's runner, so the branch point is
refused `contract.node.missing` and passed over, and the rehearsal starts at the
operator's first mark. Release 0.2's recording is meant to show the whole story the
[roadmap](roadmap.md#release-02) tells: a twin booted at the branch point where `main`
carries a network, then stepped through the marked stages to the head.

The tiers assume a `main` with no network. A `fylgja-test-*` branch made from `main` sees
`main` as it was at its creation, and test support's three-node seed does nothing when
the branch already holds a device. So on a seeded `main` every throwaway branch would hold
`main`'s devices instead of its own, and the assertions that read the branch point refused
(`internal/stages`, `internal/checks`, tier 3's case 9 and the README's walk-through)
would fail. CI's contract job brings up its own Infrahub, whose `main` stays empty
whatever this host holds.

### Options considered

- **No baseline: show the branch point refused.** Rejected: the recording would show the
  rehearsal starting at a mark, which the README's walk-through already shows, and not the
  branch point booting, which is what 0.2 is meant to show.
- **A permanent baseline on `main`.** Rejected for 0.2: the tiers would need to read their
  branches' devices beside `main`'s, and CI's contract job would need to seed `main` the
  same way so that both test one setup. That is a feature, with its own spec, and not
  release preparation.
- **A second Infrahub instance for the demo.** Rejected: the host runs one Compose project
  ([D-020](#d-020)), and a second one would need its own registration, override and
  tokens, plus about 5 GiB beside a twin.
- **An object file loaded by `infrahubctl`.** Rejected: the fixture tool is the one binary
  that writes intent, and its three-node seed is already the topology the goldens and
  `fylgja-fixture` describe. A second description of it would drift from the first.
- **Fixture-tool flags that write the three-node seed to `main` and remove exactly what
  they wrote.** Chosen.

### Decision

For a demo or a recording, and only then, `main` carries a **baseline**: test support's
three-node seed (`SeedThreeNode`), three SR Linux nodes n1–n3 in a triangle of three
links under contract `0.2`, with the devices joined to `fylgja-devices` and their
`device-config` artifacts generated on `main`.

- **`-seed-main` opens the window.** Like `-prepare-main`, it takes no `-branch`, and it
  takes no other flag. It refuses when `main` already holds anything of the baseline's
  kinds (the contract, a platform, a device, an interface or a link).
- **`-clear-main` closes the window.** It deletes what `-seed-main` wrote, links first, then
  any of the devices' artifacts Infrahub left, and leaves `fylgja-devices` empty. It
  refuses while any `fylgja-test-*` branch exists, and when `main` holds anything of those
  kinds beyond the baseline, such as a merged change, naming what it found. Either refusal
  deletes nothing. A baseline left half-written by a failed seed is still removed.
- **The order of the window.** The gates run first. Then the baseline is seeded, and only
  after that is any demo branch created, so that its stages carry no `stage.main.moved`.
  Before the window closes, every demo branch is deleted by `-delete` (stages, proposed
  changes, branches).
- **While the window is open**, no tier 2 or tier 3 runs, nobody follows the README's
  walk-through, and no proposed change is merged into `main`.

Fylgja itself does not change: nothing under `internal/` knows about the window.

### Consequences

- The demo's branch point compiles and boots, and its rehearsal reports a baseline built
  at the branch point.
- `stage.main.moved` can be tried with a real write on `main`, which 014 left open: a
  throwaway branch with a proposed change and a stage, made before `-seed-main`, should be
  warned about after it.
- `main`'s history keeps the window. A read pinned inside it sees the baseline after
  `-clear-main`, and Infrahub offers no purge, so this Infrahub stays not throwaway, as
  `CLAUDE.md` says.
- `CLAUDE.md`'s list of writes made to `main` on purpose gains the baseline, and
  `docs/development.md`'s fixture-tool paragraph gains both flags.
- **Owed to the running 1.11.2, in the window, then to
  [verified-facts.md](verified-facts.md):**
  - `fylgja-fixture`, made before the seed, sees none of it and still compiles to
    `47b2c449…`;
  - a branch made after the seed reads the baseline at its branch point;
  - after `-clear-main`, no artifact of a deleted device is left on `main`, and a new
    branch's branch point is refused `contract.node.missing` again.
- Revisit if a release wants a baseline on `main` at all times. That would be a feature,
  as above.

<a id="d-055"></a>
## D-055 — The contracts' `$id` namespace is this repository's GitHub URL

**Status:** Accepted 2026-10-22. Found by the release's hygiene check before `v0.2.0`; no
earlier entry chose the namespace.

### Context

Every schema under `contracts/` (and `psp/psp.schema.json`, its copy of the package
schema) carries a JSON Schema `$id`, and since the cut each was a URL under
`https://fylgja.dev/`: `…/contracts/twin/v7`, `…/schemas/show.schema.json`, `…/psp/v0.6`.
The schemas refer to one another by those URLs (`findings.schema.json` to the blocks'
schemas, `api.schema.json` to findings), the tests register each document under its `$id`
before compiling, and `internal/psp` compares the package schema's `$id` with a constant.
The frozen copies under each package's `testdata/contracts/` carry the same ids.
`fylgja.dev` was never this project's: the name belongs to someone else and answers on the
web, so a validator that resolved `$id`s by fetching would fetch a stranger's document
under this project's name. JSON Schema promises no retrieval by `$id`, and nothing here
fetches; the id is the name a document is registered and referred to by. No bundle, record,
finding or golden carries it, so the namespace is metadata of the contracts, not part of
any format.

### Options considered

- **Keep `fylgja.dev`.** Rejected: the project names a domain it does not control, in
  sixty-six files, and a reader cannot tell.
- **The company's domain.** Rejected: nothing is served there either, and the repository
  already has a home whose name it uses as its module path.
- **A URN (`urn:fylgja:contracts:twin:v7`).** Rejected: valid and never resolvable, but
  unfamiliar, and the ids' path shape (`contracts/<name>/v<N>`, `schemas/<file>`), which
  the documents and tests read, would be rewritten rather than re-rooted.
- **The repository's URL, `https://github.com/happypathnetworking/fylgja/`, as the
  prefix.** Chosen.

### Decision

Every `$id` and cross-document `$ref` that began `https://fylgja.dev/` begins
`https://github.com/happypathnetworking/fylgja/`, and nothing after the prefix changes:
`https://github.com/happypathnetworking/fylgja/contracts/twin/v7`,
`…/schemas/show.schema.json`, `…/psp/v0.6`. The tests' registrations and `internal/psp`'s
constant follow. The ids remain identifiers: no document is served at those paths, and
nothing resolves them. The frozen copies are renamed with the rest, since the namespace is
not part of the formats they freeze.

### Consequences

- No format version moves: not `twin.json` `7`, the bundle manifest `4`, the CTM `1`, PSP
  `0.6`, findings `1` or the API `1`. The goldens and every bundle id are unchanged.
- A client that validates a document against these schemas by `$id` updates one prefix.
- A new schema takes an id under the same prefix, and `psp/psp.schema.json` stays
  identical to `contracts/psp.schema.json`, as `internal/psp`'s contract test requires.

<a id="d-056"></a>
## D-056 — Spec Kit's files are installed on each host, not carried in the tree

**Status:** Accepted 2026-10-22. Found by the release's hygiene check before `v0.2.0`.
Beside [D-043](#d-043), which keeps the development record private and development here.

### Context

The cut carried what `specify init` had written on 2026-09-12 (Spec Kit 0.14.4.dev0): ten
`speckit-*` skills under `.claude/skills/`, five scripts under `.specify/scripts/bash/`,
five templates, and the install's metadata. Nineteen of those files were byte-identical to
what the tool wrote, as its manifests record; one, `.specify/templates/plan-template.md`,
was amended here for the Constitution Check. They are GitHub's work under the MIT licence,
whose notice did not travel with them, so the public tree redistributed them without it.
The workflow needs them on disk: the command-log hook and `scripts/converge-loop.sh` run
`check-prerequisites.sh`, and every Spec Kit command is one of the skills.
`scripts/bring-up.sh` and CI never touch them.

### Options considered

- **Keep them tracked and add the notice.** Rejected: compliant and exactly reproducible,
  but the tree would go on shipping another project's files, and each Spec Kit upgrade
  would be a commit of them.
- **Vendor a pinned copy under a `third_party/` path.** Rejected: the same, moved.
- **A bring-up part that installs Spec Kit.** Rejected for now: the script sets up the lab
  host, and a network install it cannot pin to the version that built this project widens
  it for one developer workflow.
- **Untrack the unmodified files and the metadata, ignore them, document the install.**
  Chosen.

### Decision

The nineteen unmodified files and Spec Kit's install metadata are untracked and
git-ignored; they stay on the development host's disk. A fresh clone installs Spec Kit
(`uv tool install specify-cli`, then `specify init --here --integration claude --script
sh`) before its first Spec Kit command, and gets the current release's files. The
repository keeps what is its own: the constitution, the amended plan template, with Spec
Kit's MIT notice beside it in `.specify/LICENSE.spec-kit`, and `.specify/extensions.yml`.

### Consequences

- The hook, the loop and the skills are unchanged on this host. Elsewhere, until `specify
  init` runs, the hook records no feature directory and the loop stops at its first line.
- A newer Spec Kit writes different prompt text; the files this project was built with
  remain readable in the history before this entry.
- `specs/013-launch/brief.md`'s link to `.specify/templates/spec-template.md` points at
  the file as the launch's commit holds it, in this repository's history.
- Revisit if the bring-up script grows a developer-workflow part.

---

<a id="d-057"></a>
## D-057 — The CTM is the contract; a reader maps one source's model onto it; the Infrahub reader reads the schema library's DCIM at a pinned commit

**Status:** Accepted 2026-10-25, by the operator, at `017-dcim-reader`'s converge loop
(drafted by its plan, research §3.1). Supersedes [D-002](#d-002), [D-003](#d-003) and
[D-004](#d-004) whole. Beside [D-010](#d-010) (the CTM, now the contract),
[D-016](#d-016) and [D-031](#d-031) (the package, which now also declares its platforms),
[D-028](#d-028) and [D-044](#d-044) (unchanged), [D-045](#d-045) (the stage kind,
untouched). To be built by `017-dcim-reader`.

### Context

D-002 made Fylgja's four generics the contract so that any Infrahub schema could implement
them and Infrahub would resolve one query across the implementers. They met a model
somebody else wrote on 2026-10-16, OpsMill's schema library, and could inherit into its
DCIM only after adaptation on one side or the other ([verified
facts](verified-facts.md), "Inheriting from two generics"). Release 0.3 adds NetBox as a
second source of truth, read directly; a reader of NetBox maps a fixed, DCIM-shaped model
onto the CTM whatever the contract says, and once one reader works that way the generics
make the Infrahub reader the odd one out. Nothing in the compiler, the bootstrap or verify
ever consumed an address, a loopback, a role or a site: the configuration Infrahub renders
carries them.

### Options considered

- **A layer-1 generics contract carried into the library by an overlay schema.** Designed
  first; rejected by the operator on 2026-10-25: an overlay is still something every user
  loads, and NetBox cannot load it.
- **A mapping declared as data per install.** Rejected: D-002's "mini-ETL language", now
  with two sources.
- **The CTM as the contract and one documented mapping per source.** Chosen.

### Decision

- The CTM is the contract, versioned by `ctm_version`. A reader maps one source's model
  onto it, by one documented mapping, names the model version it reads in every CTM's
  envelope, records every place the source's meaning is near the CTM's and not identical
  as an approximation and every object it leaves out as an omission, in the envelope's
  fidelity list, which the compiler copies into the manifest unchanged, and reads nothing
  the mapping does not name. The read-only test holds the list of kinds a reader queries.
- The Infrahub reader reads OpsMill's schema library's DCIM at commit `26336e9`, vendored
  under `schema/dcim/` with its licence, by the mapping in
  [`specs/017-dcim-reader/contracts/mapping.md`](../specs/017-dcim-reader/contracts/mapping.md):
  a device is a `DcimDevice` carrying the scope's tag, a port an `InterfacePhysical`,
  `mgmt_only` its `role: management` and `enabled` its `status != disabled` (approximated),
  a link a `DcimCable` with exactly two physical ends in the scope, any other cable omitted
  and recorded; the platform's version is absent. Conformance is per kind and per field
  against what the pinned commit loads, reported in one pass; completeness is strict inside
  the scope.
- A twin is of a scope, a tag on each device the request names (`--scope tag:<name>`) or
  the server's default (`FYLGJA_DEFAULT_SCOPE`), recorded in the CTM, the manifest and the
  record. A package declares the manufacturer and platform pairs it is (`platform.matches`),
  matched exactly; nothing is added to anyone's model to key a package or a scope.
- What survives of D-003 and D-004: wiring is derived by the compiler from what the reader
  mapped (`mgmt_only`, the cables), never declared in a schema; a reader never claims a
  meaning the source does not carry, and says so where it approximates.

### Consequences

- Constitution III is rewritten and VI reworded (2.0.0); the generics contract leaves the
  blessed layers and the mapping takes its place.
- `schema/`'s generics and reference schema go; `FylgjaStage` stays. The CTM moves to `2`,
  the manifest to `5`, `twin.json` to `8`, the package to `0.7`; the goldens move.
- The development host's `main` is migrated by two removal documents and `-prepare-main`;
  a DCIM user loads nothing of Fylgja's but the stage kind, and `fylgja schema check` says
  whether their DCIM is the one the reader reads; a diverged DCIM needs a reader of its
  own, which the [production shape](roadmap.md#the-production-shape) makes a seam.
- The library calls itself experimental: the commit is pinned, and moving it is a change
  to the mapping's table, verified against the running Infrahub.
- Revisit when the library publishes a release, when NetBox's reader (F-025) shows the CTM
  wants a field this one cannot fill, or when the library's `DcimEndpoint` gains a second
  implementer.

---

<a id="d-058"></a>
## D-058 — A stamp Infrahub's own account wrote is not intent's

**Status:** Accepted 2026-10-26, by the operator, at `017-dcim-reader`'s converge loop
(Run 19's question, T062). Supersedes [D-046](#d-046) in what a device's intent stamp is;
its generation stamp, its per-device comparison and its remedies stand. Beside
[D-048](#d-048), which stands. To be built by `017-dcim-reader`.

### Context

[D-046](#d-046) reads a device's intent stamp as the latest stamp of its own intent, and
017 reads that stamp as each node's `node_metadata.updated_at` (R-10). Verified on 1.11.2
on 2026-10-26 (`specs/017-dcim-reader/research.md` §5.11): about 3–7 s after a generation
first reads a port or a cable, Infrahub adds it to the template query's groups, and that
membership moves the node's `updated_at`, past the generation. `updated_by` reads `agent`,
Infrahub's own account, on those writes, and `admin` on the operator's. So a new device,
interface or cable reads stale after its hand generate, and no further generate clears it:
`TestChecksInsideInfrahub` fails, tier 3's cases 7 and 9 would, and `TestStagesListAndPlan`
passes only by reading before the join. 0.2 left the head fresh on the same route, and
FR-024 keeps 0.2's behaviour. [D-048](#d-048) already holds that Infrahub's own churn
(query groups, artifacts, memberships) is never intent; its event rule leaves those writes
out, and the stamps did not.

### Options considered

- **Read stamps from the intent fields alone**, leaving out `member_of_groups` as the
  event rule does. Rejected, as R-10 rejected it: the query names every field of the
  library's kinds and no longer covers a user's extra ones.
- **Accept the over-report**, the remedy naming the proposed change's artifact check.
  Rejected: it amends FR-024 and US5 for Infrahub's bookkeeping.
- **Leave out a stamp Infrahub's own account wrote.** Chosen.

### Decision

- A node's `node_metadata.updated_at` counts toward a device's intent stamp only when its
  `updated_by` is not Infrahub's own account. When it is, the node contributes its
  `created_at` instead: a membership is written only after a generation has read the node,
  so its creation is no later than that generation.
- Infrahub's own account is identified as Infrahub itself identifies it, verified against
  the running 1.11.2 before this is built. If it can be identified only by a name an
  install may change, the build stops and asks.
- An intent write by any account, Infrahub's own included (a generator's), is still seen
  by [D-048](#d-048)'s event read, which counts creations and updates by any account and
  leaves out only Infrahub's own relationships. The head and the staleness rule read both,
  as before.

### Consequences

- The stamps' queries read `updated_by` and `created_at` beside `updated_at` (R-10's one
  wrapper per node stays); the rule is pinned by a tier-1 test.
- `TestChecksInsideInfrahub` and `TestStagesListAndPlan` read the head fresh after a hand
  generate, without a race; tier 3's cases 7 and 9 are unchanged.
- Revisit if Infrahub records a membership on the group alone, writes it under the
  operator's account, or stamps nodes on a later generation that reads them.

---

<a id="d-059"></a>
## D-059 — Infrahub's own account is the one that writes the template query's groups, read from Infrahub

**Status:** Accepted 2026-10-26, by the operator, at `017-dcim-reader`'s converge loop
(Run 19's second question, T062). Refines [D-058](#d-058), which left how to identify the
account to be verified, and stopped the build when only a name could. To be built by
`017-dcim-reader`. Superseded by [D-060](#d-060) on 2026-10-26 in the field it reads
(`created_by` in place of `updated_by`), before it was built.

### Context

Verified on 1.11.2 on 2026-10-26 (`specs/017-dcim-reader/research.md` §5.12): Infrahub's
own identity, `__system__`, writes none of the group joins. They are written by the flow
`graphql-query-group-update` on a task worker, through the SDK, as whichever account owns
the token the install gives its task workers; on this install that is `agent`. Infrahub
creates `agent` under a fixed but editable name, with `admin`'s `account_type`, and an
install may give its workers another account's token. So nothing Infrahub records about
the account marks it as Infrahub's own: only its name does, and an install may change it.

### Options considered

- **The name `agent`, fixed.** Rejected: an install that renames the account or gives its
  workers another one reads every newly rendered node stale again.
- **A server setting naming the account**, defaulting to `agent`. Rejected: one more value
  for an operator to keep in step with Infrahub's own configuration.
- **Read it from Infrahub: the account that writes the template query's internal
  groups.** Chosen. Only Infrahub's flow writes them, on any install, with no setting.

### Decision

- Infrahub's own account, in [D-058](#d-058)'s sense, is the `updated_by` of the internal
  query groups (`CoreGraphQLQueryGroup`, `group_type: internal`) of the query the device's
  artifact definition renders through, read on the branch pinned at the stage's `at`, as
  the stamps are.
- With no such group, nothing has been generated through that query on the branch, so no
  join has been written and no stamp is left out.
- `CoreGraphQLQueryGroup` joins the Infrahub-owned kinds the stage reader may read,
  named in the read-only test with its reason (Constitution III, 2.1.0).

### Consequences

- The stamps' read gains one query of the groups' `updated_by`; the rule is pinned by a
  tier-1 test over both an `agent`-named and a renamed account.
- Revisit if Infrahub writes the groups as `__system__`, or marks its workers' account in
  the account itself.

---

<a id="d-060"></a>
## D-060 — Infrahub's own account is the one that created the template query's groups

**Status:** Accepted 2026-10-26, by the operator, at `017-dcim-reader`'s converge loop
(Run 19's third question, T062). Supersedes [D-059](#d-059) in the field it reads: its
groups, its pinned read, its no-group case and Constitution III's read (2.1.0) stand.
Beside [D-058](#d-058), which stands. To be built by `017-dcim-reader`.

### Context

[D-059](#d-059) read Infrahub's own account as the `updated_by` of the template query's
internal groups, on the premise that only Infrahub's flow writes them. Verified on 1.11.2
on 2026-10-26 (`specs/017-dcim-reader/research.md` §5.13), before anything was built, that
premise does not hold: deleting a cable as `Admin` moved its two groups' `updated_by` to
`Admin`, and `main`'s three groups have read `Admin` since a platform's delete on
2026-10-25. Under D-059 the operator's own account would then count as Infrahub's, and
D-058's rule would hide the operator's writes: `TestStagesSeeARemoval` would fail, and at
the branch point, which reads no events, a device edited after its render would read
fresh. Each group's `created_by`, and each remaining member's `created_by`, stayed `Agent`
through the delete.

### Options considered

- **Each member link's `created_by`.** Only the flow adds a member, so it stays exact
  across a change of the workers' token. Rejected: it reads every member of every group,
  one group per rendered device.
- **The groups' `created_by`.** Chosen: the flow creates each group and a delete does not
  move it, in the same query with one field changed.

### Decision

- Infrahub's own account, in [D-058](#d-058)'s sense, is the `created_by` of the internal
  query groups (`CoreGraphQLQueryGroup`, `group_type: internal`) of the query the device's
  artifact definition renders through, read on the branch pinned at the stage's `at`.
  Everything else D-059 decides stands.
- If an install moves its task workers to another account after the groups were created,
  the new account's joins count as intent: a device reads stale when it is not, never
  fresh when it is stale.

### Consequences

- The build reads `created_by` in place of `updated_by`; its tier-1 test also covers a
  group whose `updated_by` is the operator's account.
- Revisit if Infrahub re-creates the groups on a change of the workers' token, or if a
  stale reading from such a change is met in practice.

---

<a id="d-061"></a>
## D-061 — The snapshot is the core's one input: a content-addressed wrap of the CTM, stored beside the bundles, and the bundle's CTM sidecar is retired

**Status:** Accepted 2026-10-27, by the operator, at `018-snapshots`'s `analyze` (T002).
Beside [D-023](#d-023), whose rule the id reuses, and [D-057](#d-057), whose contract it wraps
unchanged. To be built by `018-snapshots`.

### Context

Two of the three live paths, `twin step` and `twin rehearse`, read, compile and file on
the server and hand the worker a bundle id; `twin create` still read the branch inside
its provision workflow (`ReadIntent`), through a CTM under `bundles/.reads/`, and the
bundle store kept the CTM a create read as a sidecar beside its bundle. Until every
workflow the core runs takes one input that is not a branch, the production shape's
core (`docs/roadmap.md`, "The production shape"), which knows no branch, proposed
change, stage or source of truth, is a drawing; and F-025's reader of NetBox, which
cannot replay time, has nothing to address a stage by.

### Options considered

- **A CTM `3` that gains the probe slot and loses `observed_at`.** One noun fewer, but
  the CTM is the reader's contract (D-057) and would carry the core's concerns.
- **A snapshot that holds the configurations beside the CTM**, as the roadmap first
  wrote. The CTM has carried each device's artifact since M5; a second copy would be a
  second place for one fact.
- **A thin document wrapping an unchanged CTM `2`**, with an empty slot for probe
  declarations and an id of its own: the hash of the body minus the read's
  `observed_at`, by the rule the manifest already uses for `bundle_id` (D-023).

### Decision

The snapshot, version `1`, wraps CTM `2` unchanged: `snapshot_version`, `snapshot_id`,
the take's `observed_at`, the `ctm` and `probes`, required and empty until 016. The id
is the SHA-256 of the body's canonical bytes with the inner `observed_at` left out and
nothing else left out: the branch, the `at`, the schema hash, the scope, the mapping,
the fidelity list, every artifact and the probe list are in it. It is filed by id under
the state root's `snapshots/`, beside `bundles/`, behind a store of the bundle store's
shape (D-014) that never overwrites, so two takes of one intent share one document, the
first take's. The server takes the snapshot for every create, as it reads for a step
and a rehearsal, and every workflow the core runs but the reconcile is started with a
snapshot id and a bundle id and never a branch. `twin create --snapshot <id|file>` and
`twin compile --snapshot <file>` take one directly; `intent read --out` writes one;
`twin compile --ctm` is retired. A twin built from a snapshot is frozen by rule: its
reference is a document, which does not move. Wherever a record names a bundle id it
names the snapshot id beside it, except in the bundle, so bundle `5` and the goldens do
not move. The bundle's CTM sidecar is retired: no build writes one, and those on disk
are never read. The reconcile workflow is untouched until F-026.

### Consequences

- The same intent read twice has one snapshot id, and one stage resolved twice has one,
  while the branch's history stands: a rebase re-stamps the branch and moves every
  stage's snapshot, and an artifact re-rendered to different bytes for unchanged data,
  never seen, would move it. Both are limits, documented and not tested.
- Where the source of truth cannot replay time, the snapshot's id addresses intent, and
  a stage there is the snapshot Fylgja took when the operator marked it: Constitution VI
  says so from 2.2.0, for F-025 to inherit.
- `twin.json` `9` names the snapshot, or none for a twin `twin provision` built or a
  following rebuild made; the five documents that name a bundle id name it too.
- The compile that `twin compile` asks for runs on the snapshot's CTM, the goldens'
  path (D-011); the CTM stays the reader's contract (D-057).
- A snapshot holds configuration, as the CTM does: it lives in the store and crosses as
  a file, never in a finding, a log line or a document.

---

<a id="d-062"></a>
## D-062 — Following leaves the core: a twin is frozen at its snapshot and changes only when the operator steps it

**Status:** Accepted 2026-10-28, by the operator, at `019-following-removed`'s plan.
Supersedes [D-012](#d-012) in the part that says an unpinned twin follows its branch,
[D-026](#d-026) and [D-027](#d-027) whole, [D-007](#d-007) in its `Reconcile`, its Schedule
`fylgja-follow` and its count of workflows and activities, and [D-014](#d-014) in "Temporal
holds … the Schedule `fylgja-follow`". Beside [D-061](#d-061), whose sentences on the
reconcile and on a snapshot twin frozen by rule it makes universal. To be built by
`019-following-removed`.

### Context

Since D-061 every workflow the core runs but one starts with a snapshot id and a bundle id
and never reads a branch. The one left is following: a twin built from a branch with no
`at` and no `--no-follow` is checked every interval by the `Reconcile` workflow, which a
Temporal Schedule starts, reads and compiles the branch in the worker, and rebuilds the
twin when its bundle id has moved (D-012, D-026, D-027). It is the core's last reader of a
branch, and the production shape's core knows none. A twin that may tear itself down
between two measurements is also the wrong ground for probes (F-014), which measure loss
across a step. Entries that cite following stand as history: D-011, D-017, D-023, D-024,
D-025, D-028, D-032, D-039, D-048, D-049, D-053.

### Options considered

- **Keep following, fed by snapshots**: the Schedule's check takes a snapshot instead of a
  CTM. The core would still hold a timer and a branch, and probes would still measure a twin
  that rebuilds itself.
- **Keep following as an opt-in flag**: the same, behind a choice nobody defaults to.
- **Remove following from the core**, and keep the step: a twin is frozen at its snapshot
  and changes only by `twin step` or `twin destroy`; any twin whose record names a snapshot
  steps to a stage of any proposed change or to a snapshot of any branch, but never to
  another scope or mapping; the orchestrator of the production shape may bring following
  back outside the core as take, compare, step.

### Decision

Following leaves the core, and not Fylgja. The `Reconcile` workflow, the Schedule
`fylgja-follow`, the five control activities (`ReadIntent`, `Compile`, `RunsInFlight`,
`StartFollowing`, `StopFollowing`), `bundles/.reads/`, `--interval` and `--no-follow` go;
the worker serves four workflows (provision, destroy, step, rehearse) and the lab's fifteen
activities, and reads no source of truth. A twin is frozen at its snapshot, and nothing
rebuilds it but the operator. **Walking means stepping**: a twin walks when the operator
steps it from one snapshot to another, through a proposed change's stages or to any
snapshot the operator takes; following was the way of walking the operator did not drive.
`twin step` steps any twin whose record names a snapshot, to `--stage <branch>/<label>` of
any proposed change or to `--snapshot <id|file>` of any branch; a twin's scope and mapping
are fixed for its life, and its branch is not. `twin show` names the snapshot a twin was
built from in place of pinned, frozen or following; "pinned" remains a word for a
reference with an `at`. Short workflows with fixed ids stand (D-007), and Temporal holds
run history alone (D-014).

### Consequences

- The core's surface is create from a snapshot, step to one, destroy, show and verify, with
  no timer of its own; an orchestrator that follows a branch calls the step.
- A twin of a 0.2 build left following across an upgrade keeps a Schedule that starts a
  workflow no worker serves: the release notes say to destroy it first, and the build
  carries no code to clear it.
- The rehearsal's baseline stops nothing, so every rehearsal's deadline is 30 s shorter.
- A twin whose record names no snapshot does not step: one built from a bundle directory,
  one whose record is of `twin.json` 8 or earlier (018's allowance for an 8 that names a
  stage is retired), and one an earlier build's rebuild made. `twin destroy` clears each.
- The constitution is not amended: no principle names following, a reconcile, a Schedule
  or an interval. VI, VIII and IX cite D-012, D-007 and D-014, each superseded here in part;
  as when D-060 superseded D-059 (cited by III) in part, Governance's sentence on a changed
  cited decision is read as a change to what a principle relies on, and none does.
- D-061's sentences that waited for this feature ("every workflow the core runs but the
  reconcile", "the reconcile workflow is untouched until F-026") now hold without the
  exception; D-061 is not superseded.

---

<a id="parked"></a>
## Parked

Directions decided before this log began, neither adopted nor rejected until the scope that
needs them is planned. Each becomes an entry when that scope is planned.

- **Firewall staging.** A firewall platform is staged: its network configuration first,
  its policy later, traffic validation last. Network configuration is what the twin already
  pushes for any platform; policy and traffic validation need more than a pushed artifact
  and a read of the node. Parked because no firewall platform is planned; needed by the
  platform track.
- **Production state collection.** Fylgja proceeds without collecting state from production
  devices: a twin is built from intent alone, and nothing compares it with the network it
  models. That is a stated limitation
  ([architecture §6](architecture.md#6-known-limitations)); read-only substitutes for it
  are a later decision.
- **The platform order's tail.** After SR Linux and EOS ([D-021](#d-021)), Cisco IOL is
  next, then PAN-OS. IOL needs no KVM and its image is account-gated, as cEOS's is, so it
  is the third platform the [roadmap](roadmap.md#backlog) plans; PAN-OS
  is a licensed VM and waits for a lab host with KVM ([D-020](#d-020)). Needed by the
  platform track.

---

<a id="open-questions"></a>
## Open questions

- ~~**Infrahub `at` edge cases**~~ (raised 2026-09-14; [D-012](#d-012)). Settled by
  [D-045](#d-045)'s verification on 2026-10-16: a read pinned before the branch's
  creation answers with `main` as it then stood, back to the default branch's creation,
  and is refused before that; a read pinned at exactly an attribute's stamp sees the
  node and not the value, so an `at` that must see a write is strictly after its stamp;
  Infrahub offers no server-side marker for the read time, so `observed_at` stays
  best-effort. Sub-second precision was settled earlier: a pinned `at` carries at most six
  fractional digits.
- **Opt-in branch coupling** (raised 2026-09-14; [D-001](#d-001)). Whether a twin's lifetime
  may be tied to its branch's, so that merging or deleting the branch destroys the twin.
  Deferred; opt-in if ever.
- **Which lab runtime does M9 target?** (raised 2026-09-22). Decide before M9 is specified.
  M9 as planned runs containerlab on remote hosts, with per-host task queues, a scheduler,
  and binaries and images shipped to each host. A Kubernetes runtime would make most of
  that redundant, because the cluster schedules pods and the kubelet pulls images. So the
  runtime has to be chosen before M9 is designed, even if the runtime itself arrives after
  it. Read from source and documentation on 2026-09-22, not verified live:
  - **clabernetes.** It runs containerlab's own kinds from a pinned containerlab module
    (0.78.0 on its `main`), so its vendor coverage is containerlab's, and Fylgja's bundle
    is already a containerlab topology. It re-cables a running lab per node
    (`link-apply-mode`: `live`, `restart`, `recreate`), taking each kind's default from that
    module. It needs a cluster ([D-020](#d-020)). The runtime described is the unreleased
    0.9, whose upgrade from 0.8 is destructive.
  - **KNE.** Its meshnet has hot-added and hot-deleted links since meshnet v0.4.0 (July
    2026), but the `kne` CLI does not expose it. It has its own topology format, so it
    would need a second compiler output. Each vendor comes through a separately installed
    operator, so every support package would need a KNE section. Its strength is
    OpenConfig's test ecosystem (featureprofiles, Ondatra), not vendor breadth. It stays an
    option for a consumer that needs that ecosystem.
  - **containerlab on remote hosts**, as M9 stands. It has re-cabled a running lab since
    0.78 (`deploy` converges a running lab: SR Linux links live, cEOS restarts the node,
    vrnetlab VMs are recreated), so incremental updates need no runtime change.

  **Recommended: M9 targets clabernetes**, as a second lab driver behind `internal/lab`,
  with the cluster standing in for "the host". M8 is unaffected: object storage and a
  Temporal cluster are needed either way. The M9 spec waits for clabernetes 0.9's release
  and a cluster to run it on. Choosing it supersedes [D-008](#d-008) (Kubernetes resources
  instead of the `clab` CLI) and [D-014](#d-014) (clabernetes's resources say whether a
  twin exists), and may supersede [D-013](#d-013) and [D-020](#d-020) (one twin per
  namespace would let two waypoint series run in parallel). It amends Constitution VII (the
  fixed lab name, and no resources named per twin). Running both drivers, containerlab on
  the development host and clabernetes remote, is the fallback, at the cost of a second
  end-to-end tier.

---

<a id="d-063"></a>
## D-063 — The feature chain: the script runs every Spec Kit stage, the operator's work precedes `specify`, and the route's lanes run at once in worktrees merged back mechanically

**Status:** Accepted 2026-10-28, by the operator, in the review of the workflow after
019's close. Supersedes development.md's rejection of worktree isolation (2026-10-17, "a
merge is a judgement"). Beside [D-056](#d-056), which stands: Spec Kit's own files stay
installed and untracked, and this repository's templates sit beside them. Built the same
day, outside any feature's spec.

### Context

Three features (017, 018, 019; 98 logged sessions, about 32 active hours) ran through the
workflow as it stood: `specify` to `analyze` by hand with `/clear` and `/model`, then
`scripts/converge-loop.sh` for the implement runs and the convergence passes. The records
showed four things. The operator's decisions leaked into the run: 017 carried 12
`[operator]` tasks and raised 11 mid-loop questions, and the review sat before converge in
every route and was deferred every time. Only half the chain was driven, and the models
drifted from development.md's table twice. The house style of `tasks.md` and of a dozen
other documents lived nowhere but the previous feature, since the tasks prompt was empty
each time and Spec Kit's template is git-ignored. Nothing guided the split into runs: they
were contiguous id ranges along user story and phase, which gave one run of five packages
at 524k context and one of nine record tasks that exhausted its budget, and the route was
strictly serial. The attended session also woke only when its Monitor expired, since its
filter took a question, a pause or a stop and nothing else, and the command log kept the
operator's typed answers verbatim.

### Options considered

- **Keep the chain by hand and codify the formats alone.** Rejected: the model drift and
  the operator questions mid-run are the chain's shape, not a document's.
- **An agent that directs the chain by judgement.** Rejected, as ideas.md's assessment of
  2026-10-04 rejected it: a model directing a long run fills its context and drifts; the
  script owns the order, the gates and the commits, and the session converses and relays.
- **Lanes with a merge the attended session judges.** Rejected: the loop's value is that
  nothing in it judges.
- **Lanes with a mechanical merge.** Chosen: git merges the branches, `tasks.md` by a
  driver that resolves ticks and notes and nothing else, and a merge that would need a
  judgement stops the loop for the operator.

### Decision

- The operator's work precedes the chain. A feature's `brief.md`, written to its template,
  settles every decision in `/feature-open`, so `clarify` has nothing to ask; an
  `[operator]` task sits in the route only where it could not precede the chain (the
  review, a push, an import that waits on the pushed tree), kept to a minimum, with its
  reason on its line.
- `scripts/converge-loop.sh` runs every stage from the one `--from` names, each as a
  headless session on the model and effort of a table the script alone holds, each gated
  on the tree before the next and its artifacts committed by the script; a stage that
  would ask the operator ends its report with a wait marker, which attended is a text
  question answered through `reply-<n>` and unattended the `operator` stop. Every event
  is a line in the run's `events.jsonl`, which the attended session watches whole and
  relays, every session's report verbatim.
- The tasks template is this repository's, tracked under `.specify/templates/overrides/`;
  a run is the tasks that share one context, with the files to read first on its line, and
  a lane mark names the run it may run beside. Runs that mark each other run at once, up
  to `--lanes`, each in its own worktree seeded with Spec Kit's ignored files, and merge
  back in route order; a lane is never live and never names tier 2.
- The command log records a command's own prompt and every segment's metrics, never a
  follow-up's text.
- The other documents Spec Kit has no template for have section lists under
  `docs/templates/`, and the forms inferred from neighbours until now (commits, task notes,
  schema histories, the code forms) are development.md's "Conventions".
- The model table: specify, clarify, plan, tasks, analyze and converge on Fable 5.1 (plan at
  `--plan-effort`, xhigh for new design), planned implement runs on Opus 5.5 at
  `--implement-effort`, high by default, live and converge-appended runs at xhigh. 016
  measures high against 017–019's per-task figures, and the flag flips the default if xhigh
  wins.

### Consequences

- development.md's "Spec Kit workflow" is rewritten around the chain; its model table is
  the script's; `CLAUDE.md`'s Workflow says the same in three bullets; the three ideas.md
  entries (a directing agent, a feature open and close skill, templates) are closed.
- A merge that stops the loop is the operator's to make by hand, with the lane's worktree
  and branch kept under `local/lanes/`.
- Converge's stopping rule and its tags are unchanged, by the operator's choice, though
  the records show convergence took about a third of each feature's time and never ended
  by the rule; a later entry may take that up.
- Revisit if a lane's merge stops the loop often, if the attended session's relay of every
  report proves too much for the operator's chat, or when 016's measurement decides the
  implement effort.

---

<a id="d-064"></a>
## D-064 — A probe is a declaration outside intent: a file first, an Infrahub kind second, carried by the snapshot and compiled into synthesized hosts by a probe package keyed by its format

**Status:** Accepted 2026-10-28, by the operator, at `016-probes`'s `analyze`. Beside
[D-010](#d-010) (the CTM, whose `synthesized` provenance it fills), [D-031](#d-031) (a
mechanism is code keyed by a format value, which the probe package is a second instance
of), [D-038](#d-038) (an assertion kind added as data), [D-045](#d-045) (the kind's pattern),
[D-057](#d-057), [D-061](#d-061) and [D-062](#d-062), which stand. Supersedes nothing.
D-010's "everything is `intent` today" is history from this entry. To be built by
`016-probes`.

### Context

Since 018 the snapshot carries `probes`, a list required and empty; since 019 a twin is
frozen at its snapshot, so loss across a step is the step's. The roadmap's F-014, re-shaped
on 2026-10-25, names the headline demo, loss during a change applied to a running twin,
and the production shape names a declaration outside intent as the third seam 0.3 proves.
The CTM carries no address (D-057), so a host's address must come from somewhere intent
does not own; the reader must not read addresses (D-057); and Infrahub renders every
device's configuration (D-028, Constitution IV).

### Options considered

- **The kind primary, the file `twin compile`'s input** (the 2026-10-16 brief). Rejected
  by the operator on 2026-10-25: F-025's NetBox reader takes a file alone, and the snapshot
  is the core's one input.
- **Probes derived from the topology.** Rejected: Fylgja guessing at intent.
- **The declaration a file first, the kind one source of it, carried by the snapshot,
  with the host's address and gateway in the declaration.** Chosen.
- **A generator type per probe, in code.** Rejected: D-031's branch on what a thing is.
- **The probe package a facet of the PSP.** Rejected: a host is of no platform.

### Decision

- A probe declaration is a document of its own (`contracts/probes.schema.json`, version
  `1`): each probe a name, a format and two endpoints, each a device, an interface, the
  address its host takes with its prefix length, and its gateway, the device interface's
  address. It carries nothing else; the rate and the burst count are the package's.
- The declaration lives in a file the operator writes, and in Infrahub as `FylgjaProbe`,
  one object per probe on the default branch, branch-agnostic, on D-045's pattern, written
  by the operator and only read by Fylgja. The kind is one source of the file, not the
  declaration: the NetBox reader takes the file alone.
- A request reads its probes once, from `--probes <file>`, else the kind, else none, with
  no warning, and every snapshot it takes carries that list. The snapshot, version `2`,
  carries the declarations beside the CTM and hashes them; a `1` reads as a `2` with none.
- The compile adds one synthesized `linux` host per endpoint to the CTM it compiles, on a
  synthesized link to the device interface, configured from the declaration alone, and
  refuses a probe it cannot place; the manifest marks every synthesized node, link and
  port. The device side is intent's: Infrahub's template renders the interface's address.
- A probe package is a mechanism keyed by a format value beside the PSP (`probes/`):
  image, acquisition, how a host is configured, how a flow runs and is read, how a burst
  runs, rate, count, resources, timeouts; `ping` is the first format.
- Transient loss across a step is measured and reported, labelled the container's, never
  asserted; unrecovered loss is a verify assertion kind, `burst`, under D-038, which fails
  a stage through D-052's gate unchanged.

### Consequences

- Bundle `6`, snapshot `2`, `twin.json` `11`; the step, verify and rehearsal documents
  carry loss and the burst; the goldens move twice. The worker serves four workflows and
  seventeen activities, which amends D-062's count.
- Constitution II is expanded and the Check's gate blesses the probe package (2.3.0).
- The reference transform renders addresses and no routing: a probe across two devices
  needs routing intent renders, which the user's own transform does; the fixture's probe
  lies across one device.
- A probe kept on `main` blocks a change that decommissions its device until it is
  deleted; the remedy is the declaration's.
- Revisit at F-020 (a second format and element shape), when the reference transform
  gains routing, or if a probe needs a rate of its own.

---

<a id="d-065"></a>
## D-065 — Fylgja is two deployables: a core that runs a twin from a snapshot, and an orchestrator on Prefect that owns the sources, the stages and the rehearsal and drives the core through the API alone

**Status:** Accepted 2026-10-29, by the operator, at `020-core-orchestrator-split`'s
feature-open. Supersedes [D-051](#d-051) whole; [D-049](#d-049) and [D-053](#d-053) in "the
check posts" and "Fylgja writes nothing"; [D-050](#d-050) in the listener it puts on the
bridge, which becomes the orchestrator's platform API; and [D-005](#d-005) and
[D-018](#d-018) in their scope, which is the core's. Beside [D-040](#d-040),
[D-041](#d-041), [D-042](#d-042) and [D-052](#d-052), which stand, and [D-061](#d-061),
[D-062](#d-062) and [D-064](#d-064), whose seam it completes. To be built by
`020-core-orchestrator-split`.

### Context

Since D-061 every workflow starts with a snapshot id and a bundle id, since D-062 the
worker reads no source of truth, and since D-064 the one input outside intent reaches the
core inside the snapshot. But the server still reads Infrahub for `twin create --branch`
and `--stage`, `twin step --stage`, `intent read`, `stage list`, `stage plan`, `twin
rehearse` and `schema check`, resolves stages, renders the rehearsal's report and holds
Infrahub's token, and `fylgja-rehearse` is a workflow of the worker (D-051). The roadmap's
production shape, discussed and tabled on 2026-10-25, named two components as a direction.
F-025's NetBox reader, settled on 2026-10-29 as a Python adapter outside the core, would
invent the orchestrator's toolchain if built first; its spike showed the core takes an
adapter's snapshot unchanged and the snapshot id reproduces outside Go
([verified facts](verified-facts.md#netbox-474-with-netbox-branching-122), 10-29). And a
platform whose resume is a rerun must start a run and poll it, which the API's open stream
does not offer ([api.md](../contracts/api.md)).

### Options considered

- **One deployable, the readers in the server.** *Rejected*: the core keeps knowing a
  branch, and every new source is Go inside it.
- **The orchestrator on Temporal, in Go.** The cheapest split and the most durable.
  *Rejected*: it defers the platform question, leaves Constitution I's freedom unused and
  makes F-025's adapter the one Python in the system.
- **The orchestrator on Kestra.** *Rejected*: declarative YAML fits an orchestrator that is
  a contract others' flows fill, which is not decided.
- **The orchestrator on Prefect, in Python.** *Chosen*: the ecosystem's language and
  F-025's, flows as code, Apache 2, no Kubernetes, and the platform Infrahub's own checks
  run on.
- **The Infrahub reader extracted unchanged as a Go tool.** *Rejected*: the orchestrator
  would hold two toolchains. *Chosen*: the reader ported to Python on `infrahub-sdk`.
- **The check posts the report, as D-053 has it.** *Rejected*: F-025's report on a NetBox
  branch has no check to post it, and a second posting path would follow.

### Decision

- **Two deployables.** The **core** is the `fylgja` binary in its three roles: it builds,
  steps, destroys, shows and verifies a twin from a snapshot, compiles one, validates a
  package and answers the step between two snapshots, and knows no branch, proposed change,
  stage or source of truth. Its server holds no Infrahub credential and listens on loopback
  alone; its worker serves provision, step and destroy. The **orchestrator** is a Python
  project on Prefect in this repository: its adapters read one source each to a snapshot
  file, the Infrahub adapter first, a port of the Go reader on `infrahub-sdk`; its flows
  read, create, step, plan, rehearse and check a schema; it posts the rehearsal's report
  back to the source; it keeps the snapshots it takes and one record of the twin it drove,
  as files; and its CLI, `fylgja-cli`, runs a flow and prints the findings document.
- **The orchestrator drives the core through the API alone**, a client under D-040's rule,
  and never starts a core workflow. Temporal stays in the core, where mid-flight durability
  matters, and is not the orchestrator's integration surface.
- **The API's long operations** may end an answer at the run's identity, and a run is read
  and interrupted by its identity, both built from the run's history so that the server
  keeps nothing between two requests (D-014). The stream stays for the core's own CLI.
- **The rehearsal is a flow of the orchestrator**, not a workflow. It posts its report as a
  thread on the proposed change, as the view shows one; the checks trigger a flow, wait for
  it and conclude from the report's status alone (D-052). The checks reach the
  orchestrator's platform API on the host's Docker bridge, which takes D-050's place, and
  the core's token never crosses the bridge.
- **Constitution I binds the core and frees the orchestrator**, which calls tools written
  in anything; D-005 and D-018 are the core's.

### Consequences

- The core's API and CLI lose `twin.rehearse`, `stage.list`, `stage.plan`, `intent.read`,
  `schema.check` and the arguments `branch`, `at`, `stage`, `scope` and `probes`;
  `twin.json` moves to `12`; the worker serves three workflows and seventeen activities,
  which amends D-062's and D-064's count.
- The `rehearsal`, `stages` and `plan` documents and the markdown report move to the
  orchestrator's contract unchanged but for `fylgja-cli` in their reproduce lines.
- The development host runs five processes; CI runs the Python tests and linter beside
  Go's.
- The constitution is amended MAJOR, citing this entry: I redefined in its scope, VIII and
  IX the core's, XI gaining the orchestrator as a client.
- **Revisit** if Prefect's open-source server lacks what the flows need (verified at plan),
  or once the orchestrator is decided a contract rather than a product.

---

<a id="d-066"></a>
## D-066 — A stage resolves to a snapshot: by a pinned read where the source replays time, else the snapshot the adapter took when the operator marked it

**Status:** Accepted 2026-10-29, by the operator, at `020-core-orchestrator-split`'s
feature-open. Supersedes [D-001](#d-001) in the part that says NetBox has no native
branching; Infrahub stays the first source of truth and the one the development host runs.
Beside [D-045](#d-045) to [D-048](#d-048), which stand for a source that replays time, and
[D-065](#d-065), which makes stages the orchestrator's. To be built by
`020-core-orchestrator-split` for Infrahub and by F-025 for NetBox.

### Context

D-045 makes a stage a moment of a proposed change, resolved to `(branch, at)` and read
pinned, because Infrahub reads a branch as it stood at any past time. NetBox 4.7.4 with
netbox-branching 1.2.2, probed live on 2026-10-29, has branches with a diff, merge, revert
and a change log of their own, so D-001's "no native branching" no longer holds; but
nothing reads a branch as it was at an earlier change
([verified facts](verified-facts.md#netbox-474-with-netbox-branching-122), 10-29). A stage
there has no `at` to resolve to. Since D-061 the core takes a snapshot alone, and since
D-065 stages are the orchestrator's.

### Options considered

- **Stages for Infrahub alone.** *Rejected*: a twin of a NetBox branch could be built but
  neither staged nor rehearsed.
- **A NetBox branch's past rebuilt from its change log.** *Rejected*: it re-makes the
  source's history outside the source.
- **A stage is a snapshot, however it is got.** *Chosen*.

### Decision

- A stage resolves to a snapshot id. Where the source replays time (Infrahub), the adapter
  resolves the mark by D-045 to D-048 and takes the snapshot by a read pinned there, at
  each request. Where it cannot (NetBox), the stage is the snapshot the adapter took when
  the operator marked it, which the orchestrator holds under its snapshot directory, by id,
  never overwritten.
- The flows that consume stages (the plan, a step to a stage, the rehearsal) take snapshot
  ids and do not know how each was got.
- NetBox with its branching plugin is a source Fylgja may read (F-025); Infrahub stays the
  first.

### Consequences

- The orchestrator's snapshot directory is the stage store's seed; F-025 adds the marking
  and the store for NetBox.
- A NetBox stage is fixed when it is marked, and a later edit to its branch does not move
  it; an Infrahub stage is read again at each request and refused when stale (D-046,
  D-048).
- **Revisit** if NetBox gains a read of a branch at a past time, or a source offers named
  points in time (014's upstream request).

---
