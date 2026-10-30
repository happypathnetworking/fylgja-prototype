# Fylgja — Verified facts

A **verified fact** is something about a dependency that was checked against the running
system, not taken from its documentation or from memory: how Infrahub, containerlab, a
network operating system, Temporal, Prefect, GitHub's runner or the host actually behaves. Each fact below names the
version it was verified against and the date. Re-verify a fact on any upgrade of the
system it names, and before a change depends on one that is not here, verify it and add
it, as the [development guide](development.md#verification-is-just-in-time) says. A fact a
session needs in every task is also in `CLAUDE.md`; this document holds the rest, read on
demand.

Terms are the [glossary](glossary.md)'s. Dates are 2026. A date alone up to 10-08 is a
verification on the development host's first release, Ubuntu 24.04.1, and a later one on
Ubuntu 26.04, where the host has run since; *(26.04, 10-08)* beside an earlier date is a
re-verification on Ubuntu 26.04 on that date, made when the bring-up script set the
development host up again there, or on a fresh 26.04 VM before it.

---

## Infrahub 1.11.2

**Reading intent**

- Infrahub filters on `created_at <= at`, so a read pinned to a second-granularity local
  "now" hides objects written in that second: send no `at` unless the operator supplied
  one. *(09-14)*
- **A read pinned at exactly an attribute's `updated_at` sees the node and not the
  value**: `count 1`, every attribute `null`; one microsecond later, the values; one
  microsecond earlier, no node; one nanosecond later (nine digits, which Infrahub takes
  and Fylgja refuses), the values. A strict comparison on the attribute's stamp, not a
  truncated stamp. The same for an artifact pinned at exactly its `checksum.updated_at`,
  which reads the previous version. So an `at` that must see a write is strictly after
  the write's stamp, and the head of a branch is pinned one microsecond after its latest
  stamp ([D-045](decisions.md#d-045)), or after its latest intent event when that is later
  ([D-048](decisions.md#d-048)). *(10-16)*
- **A read pinned before the branch's creation** reaches back through the branch point
  into `main`'s history: on a branch at exactly `branched_from`, at its `created_at`
  (140 µs earlier), a second and an hour before, the answer is `main` as it then stood,
  down to the default branch's creation; before that it is refused, `Requested time '…'
  is before branch 'main' was created at '…'`, HTTP 200 with `http_status: 422`, on any
  branch. For a fresh branch, the branch at `branched_from` and `main` at `branched_from`
  read the same. *(10-16)*
- **A relationship edge carries its own `updated_at`**: `interfaces { edges { properties
  { updated_at } } }`, `device { properties { updated_at } }`, `endpoints { edges {
  properties { updated_at } } }`. `NetworkLink` has no `name`: `id`, `hfid`,
  `display_label`, `endpoints` and the groups. Infrahub returns no instant of its own:
  `InfrahubInfo` is `version` and `deployment_id`, and HTTP's `Date` is whole seconds.
  *(10-16)* `NetworkLink` was the generics' link kind, removed at 017
  ([D-057](decisions.md#d-057)). The library's `DcimCable` has no `name` either: its
  `label` (optional, unique) is its `human_friendly_id`, and a tag added to a device
  carries the write's stamp on its edge as the edges above do. *(10-25;
  `specs/017-dcim-reader/research.md` §1.3, §1.7)*
- **A branch read pinned at its own `branched_from` is `main` as it then stood**, so on an
  install whose `main` holds the schema and no intent it finds no contract object, no
  device and no link: `fylgja-fixture` pinned at its `branched_from` reads `FylgjaDevice 0,
  FylgjaLink 0, FylgjaContract 0`, and Fylgja's read refuses it `contract.node.missing`, as
  it refuses `main` itself. *(10-16)* That was on the generics' kinds and their contract
  node, retired at 017 ([D-057](decisions.md#d-057)). On the library's kinds the branch
  point of a branch made from a `main` holding no network reads no device carrying the
  scope's tag, and the read refuses it `scope.empty` (`scope tag:fylgja selects no device
  on branch main at <branched_from>`). *(10-26; `specs/017-dcim-reader/research.md` §5.5,
  §5.15)*
- `Branch(name:)` answers `[Branch!]!`, `[]` for a name that does not exist. Its
  `branched_from` and `created_at` are strings in the `Z` layout
  (`2026-10-08T05:47:52.538186Z`), where an attribute's stamp reads `+00:00`; a fresh
  branch's `branched_from` is 98–140 µs after its `created_at`. *(10-16)*
- **An unset filter is no filter**: an attribute filter (`name__value`) sent as a `null`
  variable, omitted, or written as a literal `null` answers every object; `""` answers
  none. So a query whose filter is optional omits the variable rather than sending `""`.
  *(10-16)*
- GraphQL takes the branch in the path, `POST /graphql/<branch>`; a `?branch=` parameter
  there is ignored and answers for `main` with no error. *(09-30)*
  Reported upstream: **dropped**, on the operator's word. It is
  [opsmill/infrahub#10686](https://github.com/opsmill/infrahub/issues/10686), closed as
  completed on 2026-09-28 with no linked change, no mention in the 1.11.3 or 1.11.4
  release notes, and the GraphQL app on `stable` and `develop` still reading the branch
  from the path alone. *(10-08)*
- `POST /graphql` and `GET /api/schema` with no branch answer for the default branch.
  *(09-28)*
- A pinned `at` carries at most six fractional digits: Infrahub takes nine, and Fylgja
  refuses more than six (`intent.at.precision`). *(09-16)*
- History is append-only and stamped by the server's clock: no mutation takes an `at`,
  and nothing writes at a past instant. `?at=` on a mutation's request is parsed and
  validated as a query time, then ignored: a create and an update sent with `at` an hour
  earlier both answered `ok: true`, each attribute's `updated_at` was the moment of the
  write, and a read pinned between the create and the update saw the value before it. The
  SDK's `save`, `create` and `update` take no `at`; its `execute_graphql` does, and builds
  the URL the same way for a query or a mutation. In the source, mutations look the node
  up and save it with no `at`, then overwrite the context's `at` with now for the
  resolvers that follow. *(26.04, 10-13)*
  **Bug, to follow up on**: the silence. The parameter is accepted and validated on a
  request it has no effect on, and the client is not told. A malformed `at` is refused on
  a mutation while a well-formed one is discarded. Not yet reported upstream.
- On a host where Infrahub runs in Docker beside Fylgja, a local clock read after a seed
  returns sees the whole seed, since both share the kernel's clock (verified on an earlier
  host where Docker ran in a VM).
  *(09-16)*
- This Infrahub reads anonymously: no token is 200 on GraphQL, the schema and storage,
  and a wrong token is 401, so no token exists that reads intent and cannot read content.
  *(09-18)*
- **The library's roots filter one level deep.** On OpsMill's schema library at
  `26336e9`, `DcimDevice` takes `tags__name__value`, `tags__ids` and `tags__isnull`, beside
  `ids`, `name__value`, `platform__name__value`, `platform__ids`, `status__value` and
  `node_metadata__updated_at__after`. `InterfacePhysical` takes `device__ids`,
  `device__name__value`, `device__isnull`, `connector__ids`, `connector__isnull`,
  `role__value` and `status__value`, and no `device__tags__*`: `Unknown argument
  'device__tags__name__value' on field 'Query.InterfacePhysical'`. `DcimCable` takes
  `connected_endpoints__ids` and `connected_endpoints__isnull`. A nested `interfaces` takes
  attribute filters and `ids` alone, nothing that narrows it to a kind, and its edges carry
  both `InterfacePhysical` and `InterfaceVirtual`. *(10-25;
  `specs/017-dcim-reader/research.md` §1.4)*
- **Devices by tag, then their ports by device.** `DcimDevice(tags__name__value: "fylgja")`
  answers exactly the devices carrying the tag, and a tag no device carries answers 0.
  `InterfacePhysical(device__ids: […])` answers every physical port of those devices and no
  `InterfaceVirtual`. A port's `connector` reaches its cable through the generic
  `DcimConnector` (`... on DcimCable`), with every end in `connected_endpoints` and each
  end's device; an uncabled port's `connector.node` is `null`. `DcimCable(connected_endpoints__ids:
  […])` answers every cable with an end among the ports. A device read by name with its
  `interfaces` answers each physical port with its far end's device and port, and a
  loopback as `InterfaceVirtual` with no `connector`. *(10-25;
  `specs/017-dcim-reader/research.md` §1.6)*
- **Selecting an uncabled port's `connector { node_metadata }` refuses the whole page.**
  The SDL declares `NestedEdgedDcimConnector.node_metadata` non-null, and an uncabled
  port's is null, so the answer is `Cannot return null for non-nullable field
  NestedEdgedDcimConnector.node_metadata`, once per uncabled port.
  `InterfacePhysical(device__ids: …, connector__isnull: false)` reads the cabled ports
  alone and is accepted. *(10-26; `specs/017-dcim-reader/research.md` §5.8)*
- **A pinned read filters by the tag as it stood at `at`**: a device tagged after `at` is
  absent from a read pinned at `at` and present in the unpinned one. *(10-26;
  `specs/017-dcim-reader/research.md` §5.7, §5.11, §5.14)*
- **Two unpinned reads of an unchanged branch differ in `observed_at` alone.** `intent
  read --branch fylgja-fixture --out` twice, 14.4 s apart, wrote two CTM files whose `diff`
  was one line, the envelope's `observed_at`: Infrahub served the same artifact bytes to
  both, and the reader's output is canonical. Taken as snapshots by 018's code, two reads
  6 s apart had one snapshot id, `576933ec51da1866cdb7344a1f60ec2468c7d04aafbee42eacd171e352bd7482`,
  the fixture branch's, their files differing in the outer and the inner `observed_at`
  alone, and the store kept the first take byte for byte. *(10-27;
  `specs/018-snapshots/research.md` §1.2, `tasks.md` T051)*
- **A stage resolved twice has one snapshot.** On a throwaway branch with a draft proposed
  change and one mark, `stage plan` run twice on the server's path named one snapshot id
  for the mark and one for the head in both plans, left the stored documents' bytes and
  back-dated modification times as the first plan filed them, and wrote no
  `bundles/*.ctm.json`: `4345c370…` and `e8c31ea7…` alone, `ffa1dac8…` and `98c33061…`
  within the whole contract tier (`TestStageTakenTwiceHasOneSnapshot`). Three reads pinned
  at one T, before a change, after n1's regeneration and after the change, had one
  snapshot id with different `observed_at`s, and the unpinned read after the change
  another (`TestPinnedReadReproducesItsBundle`). *(10-27; `specs/018-snapshots/tasks.md`
  T045, T046, T053)*
- **A snapshot id is reproduced by `jq`.** jq 1.8.1 keeps a file's key order, indents two
  spaces as Go's `json.MarshalIndent` does, writes `[]` for an empty list and ends with a
  newline, so `jq --indent 2 '{ctm: del(.envelope.observed_at), probes: []}' <ctm> |
  sha256sum` over a bare CTM, and the schema's recipe over a stored snapshot, `jq --indent
  2 'del(.snapshot_version, .snapshot_id, .observed_at, .ctm.envelope.observed_at)' |
  sha256sum`, give the id Go's encoding gives: `576933ec…` for the fixture's two reads and
  its stored snapshot, and `6a8a77fc…`, `91f06cd6…` and `37d7378a…` for the three wrapped
  fixtures. Exact only while no string holds `<`, `>`, `&`, U+2028 or U+2029, which Go's
  encoder escapes and `jq` writes raw, or DEL (U+007F), which `jq` writes `\u007f` and Go
  raw; no fixture artifact holds any. *(10-27; the last three 10-28;
  `specs/018-snapshots/research.md` §1.3)*
- **A read pinned before a kind's load answers as if the kind had no node then.** A query
  pinned at an `at` is checked against the branch's current GraphQL schema.
  - On `main`, pinned at `2026-10-26T12:00:00Z`, before the schema library's load at
    12:29:43Z: `DcimDevice(tags__name__value: "fylgja") { count }` answers `count 0`,
    HTTP 200, no error.
  - That took about 10 s, against 0.1 s pinned after the load.
  - So Fylgja's read of `main` pinned then passes conformance, which reads the schema
    unpinned, and is refused `scope.empty`, naming the `at`.
  *(10-26; `specs/017-dcim-reader/research.md` §5.16)*
- `BuiltinTag` is branch-aware: tags made on a branch read there and not on `main`, before
  and after the branch's deletion. *(10-25; `specs/017-dcim-reader/research.md` §1.5)*

**The schema**

- `POST /api/schema/load` rebuilds the branch's GraphQL schema asynchronously: poll for
  the new mutations before seeding. *(09-14)*
- Attribute names `class` and `type` are refused, descriptions are capped at 128
  characters, and strict mode is on; check a candidate with `POST /api/schema/check`.
  *(09-14)*
- An attribute name is 3–64 characters of `[a-z0-9_]`, so `at` is refused by the JSON
  schema before it reaches the server; the stage's attribute is `as_of`. *(09-28)*
  Reported upstream: **dropped**, on the operator's word. Infrahub's schema reference
  documents the rule (`Length: min 3, max 64` in
  [attribute.mdx](https://github.com/opsmill/infrahub/blob/infrahub-v1.11.2/docs/docs/reference/schema/attribute.mdx)
  at 1.11.2). *(10-08)*
- On a branch without the Fylgja schema the generics are not GraphQL types at all: read
  `GET /api/schema` first, whose `used_by` is the conformance check. *(09-14)* The
  generics were retired at 017 ([D-057](decisions.md#d-057)); a kind a branch lacks is no
  GraphQL type there either, and the library's conformance reads the same endpoint
  (below).
- **OpsMill's schema library at `26336e9` loads beside the stage kind.** Its six files
  (`base/dcim.yml`, `base/ipam.yml`, `base/location.yml`, `base/organization.yml`,
  `extensions/cable/cable.yml`, `extensions/location_minimal/location_minimal.yml`;
  Apache 2.0, its README calling it examples "in an experimental phase") define 27 kinds.
  `schema check` passes them, every kind `added`, and `infrahubctl schema load` sends every
  file it is given in one request: 27.5 s and 24.1 s on a branch, `DcimDevice` on `GET
  /api/schema` a second later. They add 14 nodes and 21 generics, on a `main` that held
  the generics and on one that did not. *(10-25, 10-26;
  `specs/017-dcim-reader/research.md` §1.2, §5.1, §5.2, §5.10)*
- **`GET /api/schema?branch=` reports a node with its inherited attributes and
  relationships merged in**, each marked `inherited`, and each generic with its `used_by`.
  An attribute with a `default_value` is reported `optional: true`, whatever its file says:
  `DcimDevice.status`, written `optional: false` with a default, reads `optional: true`. So
  a conformance table is what the endpoint answers, not the YAML. What it answered for the
  five kinds the reader reads is the table in `specs/017-dcim-reader/research.md` §1.3:
  among them, `DcimDevice` inherits `CoreArtifactTarget` with `artifacts` merged in,
  `InterfacePhysical.device` is a required `Parent` on `DcimGenericDevice`, and
  `DcimCable.connected_endpoints` has `min_count 0, max_count 0`. `used_by` reads
  `DcimGenericDevice` [`DcimDevice`], `DcimInterface` [`InterfacePhysical`,
  `InterfaceVirtual`], `DcimEndpoint` [`InterfacePhysical`] and `DcimConnector`
  [`DcimCable`], the same on a whole branch's untrimmed document (76 nodes, 56 generics).
  *(10-25, 10-26; `specs/017-dcim-reader/research.md` §1.3, §5.2, §5.3)*
- **What the library demands of a write** (raw GraphQL on a branch). A `LocationSite`
  needs its whole parent chain (`parent is mandatory for LocationSite at parent`): a
  region, a country, a metro and the site, in four creates. A `DcimDevice` needs its
  `location` (`location is mandatory for DcimDevice at location`), and loads with no
  platform, or with a platform that has no manufacturer. An `IpamIPAddress` loads with
  `address` and `interface` alone. A `DcimCable` takes zero, one, two or three ends; an end
  on an `InterfaceVirtual` is refused (`… cannot be added to relationship, must be of type:
  ['InterfacePhysical'] at connected_endpoints`), and a port already cabled refuses a
  second cable (`Node <id> has 2 peers for connector__endpoint, maximum of 1 allowed`).
  `DcimCableUpdate` with one peer in `connected_endpoints` replaces the set, and the former
  end's `connector` reads `{"node": null}`. *(10-25, 10-26;
  `specs/017-dcim-reader/research.md` §1.5, §5.9)*
- `POST /api/schema/load` answers `diff`, `hash`, `previous_hash`, `schema_updated` and
  `warnings`, `hash` equal to `previous_hash` for a load that changed nothing. *(10-26;
  `specs/017-dcim-reader/research.md` §5.2)*
- The IPAM core generic is `BuiltinIPAddress`; `IpamIPAddress` is a demo-schema name.
  *(09-14)*
- **An `IpamIPAddress` written with no `ip_namespace` lands in the default namespace**:
  `IpamIPAddressCreate(data: {address: {value: "10.1.10.1/24"}, interface: {id: <port>}})`
  reads back the one `BuiltinIPNamespace`, an `IpamNamespace` named `default` with
  `default: true` (a field of `IpamNamespace` alone, so a query of the generic reads it in
  an inline fragment). Its `interface` reads back the port by id, name and device, the
  port's fields in a `... on InterfacePhysical` fragment, since `InterfaceLayer3` has no
  `name`. Deleting the baseline (`-clear-main`'s order) takes the addresses first, and
  `BuiltinIPAddress` then counts 0. *(10-29; `specs/016-probes/research.md` §5, T033)*
- The SDL is at `GET /schema.graphql?branch=`, and its field order is not stable between
  two fetches of one branch: compare two SDLs structurally, never with `diff`. *(09-28)*
  Reported upstream: **dropped**, on the operator's word, before any search. *(10-08)*
- `GET /api/schema/summary?branch=` returns the branch's schema hash as `main`; the hash
  covers the branch's whole schema, so two branches whose schemas differ, and two installs,
  have different hashes. *(09-14)* Three installs that loaded the same `schema/` gave
  three: `4d5b37aa…` on 24.04, `3b686463…` on a fresh 26.04 VM and `0aba31a5…` on the
  26.04 development host, so a bundle compiled from the same fixture differs by install.
  *(26.04, 10-07, 10-08)* The hash is the branch's, not the content's: the stage kind
  loaded on `main` moved it `0aba31a5…` → `e7386419…`, the same candidate on a throwaway
  branch gave `17934826…`, and the fixture branch kept `0aba31a5…`, reading
  `schema_differs_from_default_branch: true` and compiling to `47b2c449…` as before.
  *(10-16)* The final stage kind, a migration of the candidate, moved `main` to
  `a8045729…`, and the removal of M10's kind to `e43a1cde…`; the fixture branch kept
  `0aba31a5…` through both, and the removed kind with it, until 017. *(10-16, 10-17)* At
  017 the removal of the reference kinds moved `main` to `6a17b5b8…`, the generics' to
  `66edc133…`, and the library's load to **`7c413536…`** (71 nodes, 52 generics). The
  fixture branch, re-seeded from that `main`, carries its `7c413536…` and compiles to
  `f764e1fa…`. Each throwaway branch that loaded the library had a hash of its own
  (`908871c3…`, `4161a1ab…`, `18e5f023…`, `5dbe3902…`). *(10-25, 10-26;
  `specs/017-dcim-reader/research.md` §1.2, §5.1–§5.3, §5.5, §5.10)*
- A branch sees `main` as it was when the branch was created, so a branch made before the
  schema, the group or the template existed on `main` cannot render. *(09-18)* Nor can it
  query a kind loaded on `main` after it; a branch created after the load carries
  `main`'s hash and can. *(10-16)*
- **`BranchRebase` re-stamps the branch's whole content at the new `branched_from`** (its
  `created_at` stays) and discards the branch's earlier history from pinned reads: after
  it, every attribute at the head carries the rebase's stamp, a read pinned at any
  earlier instant on the branch finds none of the branch's nodes (`main`'s and
  branch-agnostic objects still read), a read at exactly the new `branched_from` finds
  the nodes with every value `null`, and one microsecond later the branch's own values,
  not the base's. `wait_until_completion: true` returned after 67 s for a three-node
  branch. *(10-16)* With nothing on `main` to take, a rebase still moves `branched_from`
  to its own moment, and returned after 17–22 s on a seeded branch with an open draft
  proposed change and five stages. *(10-16, 10-17)* Over the library's kinds a rebase
  moves `branched_from` and leaves `created_at`, as before. *(10-26;
  `specs/017-dcim-reader/research.md` §5.11, §5.14)*
- **A kind is removed by a document naming it with `state: absent`.** `schema check`
  prints `removed: <kind>`; the load takes about 40 s (38 s on a branch, 40.8 s on
  `main`); the kind then leaves `GET /api/schema` and GraphQL (`Cannot query field '<kind>'
  on type 'Query'`); and a second check prints `No diff` and a second load is already up to
  date in 1.8 s, so the document is idempotent. *(10-16, 10-17)*
- **Implementers and their generics go in two documents.** One naming the four `Network*`
  kinds, `FylgjaContract` and the four generics `state: absent` is refused by `schema
  check` and `schema load` alike (`NetworkDevice: Relationship 'interfaces' is referring
  an invalid peer 'FylgjaInterface'`): the implementers are validated against a schema the
  generics have already left. The implementers' document, then the generics', each loads:
  47.2–50.6 s and 42.2–48.0 s, on a branch and on `main`, and each checks `No diff` after.
  Then `CoreArtifactTarget`'s `used_by` is empty, the removed devices' `CoreArtifact`s
  stay, and a branch-agnostic kind beside them and its extension are untouched. A removal
  that leaves a generic with no implementer (`DcimCable`, `DcimConnector`'s one) loads too,
  in 1 m 48.8 s while a second branch loaded the library. *(10-25, 10-26;
  `specs/017-dcim-reader/research.md` §1.9, §5.7, §5.10)*
- A change of an attribute's `parameters.regex` shows in `schema check`'s diff as
  `parameters.regex, regex`: `regex` is Infrahub's mirror of it, read equal on `main`
  after the load, not a second change. *(10-16)*
- `infrahubctl schema format` (SDK 1.23.2) does not keep the order the stage kind's
  contract calls canonical: `--check` exits 1 on `schema/fylgja-stage.yaml` and would put
  `order_by` before `display_label`, `generate_profile` and `branch` after
  `uniqueness_constraints`, and an attribute's `optional` after its `parameters`. No gate
  runs it. *(10-16)*

**Branch-agnostic kinds**

Verified with M10's kind, since removed, and held by the stage kind (below).

- A branch-agnostic kind loads with a uniqueness constraint over two members, and a
  duplicate is refused from any branch as `Violates uniqueness constraint '<a>-<b>'`:
  HTTP 200 on the wire, with `http_status: 422` in the extensions. *(09-28)*
- A branch-agnostic object is read from every branch whose schema has its kind, and
  outlives the branch it names. *(09-28)*
- A branch whose schema lacks the kind cannot query it at all: `Cannot query field
  '<kind>' on type 'Query'`, HTTP 200 with a 500 in the extensions that names nothing, so
  read the kind's presence from `GET /api/schema` first. *(09-28)*
- At creation every attribute carries one `updated_at` (six fractional digits, `+00:00`),
  and only a real change of its value moves it: a description edit leaves another
  attribute's stamp where it was, and a change of that attribute's value moves it. No
  node-level `_updated_at` is queryable on a Fylgja kind. *(09-28)* The node type has no
  stamp of its own on the library's kinds either (`Cannot query field 'node_metadata' on
  type 'DcimDevice'`). The `Edged*` wrapper around it does: `node_metadata { created_at
  updated_at created_by updated_by }`, on a root's edges and on every nested relationship's
  (**Artifacts and repositories**, below). *(10-25; `specs/017-dcim-reader/research.md`
  §1.4)*
- A DateTime attribute keeps a written value character for character and refuses only
  nonsense; GraphQL types it as text, so an unwritten one reads `""`. A Number attribute's
  value is a `BigInt`, which genqlient binds to `int`. *(09-28)*
- An object file's top-level `kind` is `Object`, with the node's kind under `spec.kind`;
  `infrahubctl object validate` checks one without writing. *(09-28)*
- **The probe kind `FylgjaProbe` loads on `main`** (016): `schema check --branch main
  schema/fylgja-probe.yaml` answers `Valid!` with a diff that adds the kind alone, and
  `schema load` takes about 24 s, moving `main`'s hash from `7c413536…` to
  **`8a4f3bee8a72549e97450a414da18107`**. An object written by `infrahubctl object create
  FylgjaProbe --set …` reads back every `Text` attribute as written (addresses with their
  prefix length, gateways without), `format` defaulting to `ping`, and `object delete
  FylgjaProbe <name>` removes it by its human-friendly id. *(infrahub-sdk 1.23.2, 10-28,
  10-29; `specs/016-probes/research.md` §5, T009, T017)*

**Proposed changes and stages (014)**

- A branch-agnostic kind with a required `cardinality: one` relationship to
  `CoreProposedChange`, and an `extensions` entry giving the proposed change the reverse
  relationship under the same identifier, passes `schema/check` (202, no warning) and
  loads (200) on a branch and on `main`. A human-friendly id through the proposed change
  is refused: its `name` is not unique, and two proposed changes did carry one name.
  *(10-16)*
- The uniqueness constraint `[proposed_change, label__value]` refuses a duplicate from
  the branch as `Violates uniqueness constraint 'proposed_change-label'`, HTTP 200 with
  `http_status: 422`, and accepts the same label on another proposed change. A stage's
  write stamps every attribute and the relationship edge alike; an unwritten `as_of`
  reads `null` with its own `updated_at`. The reverse relationship lists the stages, and
  the kind filters by `proposed_change__ids` and by any attribute of the peer
  (`proposed_change__source_branch__value`, `proposed_change__state__value`). Where the
  schema lacks the kind, both the kind and the extension's field are unqueryable (`Cannot
  query field 'fylgja_stages' on type 'CoreProposedChange'`). Infrahub knows nothing of
  the labels Fylgja reserves: a stage labelled `head` is accepted. *(10-16)*
- **A stage is written three ways, each on the default branch** (`infrahubctl` 1.23.2,
  connecting as Admin, with no `--branch`): an object file (`kind: Object`, `spec.kind:
  FylgjaStage`) whose `proposed_change` is the proposed change's id as a bare string,
  since the peer has no human-friendly id, checked by `object validate` and written by
  `object load`; `object create FylgjaStage --set label=… --set proposed_change=<id>`; and
  the UI, whose menu lists `Stage` as an entry of its own in the object section
  (`/objects/FylgjaStage`), under no Fylgja group. `object delete FylgjaStage <id> --yes`
  deletes one. *(10-17)*
- **Editing a stage moves only the stamp of the attribute whose value changed**
  (`FylgjaStageUpdate` on the default branch): a `description` edit leaves
  `label.updated_at` where the create put it; a relabel moves it to the relabel's write, so
  a stage with no `as_of`, whose `at` is that stamp, then names the moment of the relabel;
  an `as_of` written after it reads back verbatim and leaves `label.updated_at` where the
  relabel put it. *(10-17)*
- **One open proposed change per source branch**: a second is refused, `An open proposed
  change already exists for branch '…'` (422); `main` as source, `The source and
  destination branch can't be the same`; an unknown source, `The specified source branch
  for this proposed change was not found`. Created `state: open`, draft or not;
  `CoreProposedChangeUpdate` with `state: closed` closes it and frees the branch for
  another, and the closed one's stages still read. *(10-16)*
- **Deleting the source branch cancels its open proposed change asynchronously**: the
  `BranchDelete` returned after 8 s, the proposed change read `open` 2 s later and
  `canceled` 16 s after the call, still naming the branch. *(10-16)*
- **A network on `main`** ([D-054](decisions.md#d-054)'s window, written by
  `-seed-main` and removed by `-clear-main`):
  - A branch made before a write on `main` does not see it. `fylgja-fixture`, and a
    throwaway branch made two minutes before the seed, kept their own three devices, and
    the fixture still compiled to `47b2c449…`.
  - A branch made after the seed reads `main`'s network as its own data and at its branch
    point (3 devices, 3 links).
  - `stage.main.moved` is warned against a real write: one warning on the earlier
    branch's `stage list` naming `main`'s latest write (the seed's end, `19:24:40.538763Z`)
    and the branch point (`19:22:44.514878Z`), exit 0, nothing refused. The branch point
    still reads `main` as it was, refused `contract.node.missing`.
  - A proposed change from a branch with no change from `main` gets no artifact
    validator, only Data Integrity and the user checks. The fixture tool's `-propose`,
    which waits for one whenever the group has members, times out after 6 minutes there.
  - Deleting a device on `main` deletes its artifacts with it (3 of 3) and leaves the
    group empty. A branch made after `-clear-main` has its branch point refused
    `contract.node.missing` again. *(10-21)* That was on the generics' kinds; on the
    library's, a branch point with no tagged device is refused `scope.empty` (**Reading
    intent**, above).
- **`CoreProposedChangeDelete` is refused while a `FylgjaStage` names the proposed change**,
  with the kind on `main`: `Cannot delete CoreProposedChange '<id>'. It is linked to
  mandatory relationship proposed_change on node FylgjaStage '<id>'`
  (`ATTRIBUTE_CONSTRAINT_VIOLATION`, 422), nothing changed. Delete stages first. The
  10-16 reading, that the delete succeeds and leaves the stage naming nothing, held only
  while the kind was loaded on the throwaway branch alone. `RelationshipRemove` of a stage
  from the proposed change's `fylgja_stages` is accepted, though, and leaves the stage
  with `proposed_change { node: null }` despite `optional: false`: that is how an orphan
  arises, and nothing then ties it to its branch. *(10-17)*
- **The proposed change's pipeline regenerates the source branch's artifacts** at its
  creation, draft (`is_draft: true`) or not — the checksums moved 59 s after the create,
  `Data Integrity`, `Schema Integrity` and `Artifact Validator: <transform>` each
  `completed`/`success` — and on `CoreProposedChangeRunCheck(data: {id, check_type:
  ARTIFACT})`, 22–25 s after the call. An intent change on the branch under an open
  proposed change, and `is_draft` set to `false`, start nothing (100 s each). `CheckType`
  is `ARTIFACT, DATA, GENERATOR, REPOSITORY, SCHEMA, TEST, USER, ALL`. *(10-16)*
- **A draft's creation pipeline** is one validator, `Data Integrity`, on a branch with
  the schema alone, concluded 7–17 s after the create; on a seeded branch three, the
  artifact validator last, concluded 36–63 s after it. Infrahub creates the validators one
  after another as the pipeline goes, so "every validator completed" holds in the gaps
  between them; nor do the pipeline's tasks (`InfrahubTask(related_node__ids:)`) end it,
  since a child task is listed only once it starts, about 2 s after its parent ends, and
  the first validator has no task. So a wait for the creation's rendering waits for a
  completed `CoreArtifactValidator` too, when the branch's devices are in the target
  group. Infrahub once dropped a poll's connection mid-pipeline (`EOF`), so a poll reads
  again until its deadline. *(10-16)*
- **A run-check re-runs the same artifact validator, identical bytes or not.** After an
  edit the template renders nowhere (`site`), a hand generate moves no checksum, and
  `CoreProposedChangeRunCheck(check_type: ARTIFACT)` leaves the checksums where they were
  while the proposed change's one `CoreArtifactValidator` node reads `in_progress`, then
  `completed`/`success` with a new `completed_at`, 17–24 s after the call. *(10-16)*
- **Pinned, the validator's `completed_at` value is the one in force at `at`, and its
  `updated_at` the current one**: compare the value, never the metadata.
  `CoreArtifactValidator(proposed_change__ids:)` answers pinned, with `definition { node {
  name } }` naming the definition. *(10-16)*
- **The diff tree is a cached calculation.** `DiffTreeSummary(branch)` and
  `DiffTree(branch)` carry `from_time` (= `branched_from`), `to_time`,
  `num_untracked_base_changes`, `num_untracked_diff_changes` and nodes with
  `last_changed_at`; after the run-check, `to_time` stayed at the run-check's moment
  through four intent writes and a regenerate (`num_untracked_diff_changes: 6`), and the
  `DiffUpdate` mutation (`branch, name, from_time, to_time`) is what refreshes it. A
  rebase resets `from_time`. So neither a branch's head nor `main`'s movement is read
  from it without a write. *(10-16)* A link deleted after it was computed leaves `to_time`,
  `num_removed` and every node's `last_changed_at` where they were, the link still
  `ADDED`; only `num_untracked_diff_changes` moves (0 → 46 for the one delete), a count
  with no stamp. *(10-17)*
- **A removal leaves no stamp the intent's read sees.** A link deleted
  (`NetworkLinkDelete`) as a branch's latest write moves no `updated_at` the stamps read
  selects: the link's endpoint edges go with it, and each endpoint interface's `link` edge
  reads `properties.updated_at` `null`. So a head computed from the stamps is pinned before
  the delete and reads the link, while an unpinned read does not; a mark after the delete
  reads it gone, with artifacts that still render it. A generate after it moves the
  checksums of the devices whose rendering named the link, and only then the head.
  *(10-17)* That was `NetworkLink`'s, since removed; a `DcimCable`'s delete is stamped
  (below).
- **A `DcimCable`'s delete stamps its two end ports.** On the schema library's kinds,
  `DcimCableDelete` moves each end port's `node_metadata.updated_at` to one instant
  between the mutation's send and its return. The stamp reads alike through the device's
  `interfaces` edge and the `InterfacePhysical` node. Nothing else moves: not either
  device's own stamp, not the platform's, not a third device's, not another cable or its
  ports. So a stamps read sees the delete on both end devices. A cable's re-ending to one
  end (`DcimCableUpdate`) is still stamped nowhere, and the event log names it
  ([D-048](decisions.md#d-048)). The delete's own `node.updated` events on the two ports
  are stamped 0.1–1.1 s later than the ports, so the head passes the ports' stamp by
  them. *(10-26; `specs/017-dcim-reader/research.md` §5.11)*
- **A generation's query-group joins are written by the task workers' account.** About
  2–7 s after a generation's checksums move, each node the template query read for the
  first time joins that query's `CoreGraphQLQueryGroup` (`device_config__<hash>`,
  `group_type` `internal`). The join moves the node's `node_metadata.updated_at` past the
  checksum. After a seed's own render, that is every device, platform, manufacturer,
  interface and cable. A node already a member does not move. Infrahub's flow
  `graphql-query-group-update` writes the join on a task worker through the SDK, as the
  account that owns the worker's `INFRAHUB_API_TOKEN`. The account then reads as the
  node's `updated_by`, and as the group's `created_by` and `updated_by`. Infrahub's
  published Compose file gives the workers `INFRAHUB_INITIAL_AGENT_TOKEN`, whose account
  Infrahub's initialization names `agent`, with `account_type` `User`, the same as
  `admin`'s. Infrahub's own identity is `__system__` (shown as `[Infrahub System]`), which
  created both accounts and wrote none of the joins. `CoreGenericAccount.name` is not
  read-only. So nothing but that account's name, which an install may change, tells it
  from an operator's. Read in the image's source and live. *(10-26;
  `specs/017-dcim-reader/research.md` §5.11, §5.12)*
- **A member's delete restamps the template query's groups as the deleting account.**
  When a node that a generation read is deleted (a cable, a port, a device, a platform),
  it leaves each `CoreGraphQLQueryGroup` it was a member of. Each of those groups then
  reads `node_metadata.updated_by` as the deleting account, at the delete's instant. The
  event log records it as `infrahub.node.updated` and `infrahub.group.member_removed` on
  the group, under that account's `account_id`. The group's `created_by` does not move,
  and neither does any remaining member edge's `relationship_metadata.created_by`: both
  stay the task workers' account. So the groups' `updated_by` is not that account alone.
  Seen on `main` after the 10-25 platform delete, and on a throwaway branch after a
  `DcimCableDelete`. *(10-26; `specs/017-dcim-reader/research.md` §5.13)*
- **The groups' creator is read beside the stamps, and a joined node keeps its
  creation.** In one stamps query, each artifact reaches its query by
  `definition.node.transformation.node.query.node.id`. The transformation's `__typename`
  is `CoreTransformJinja2` for the template. `CoreGraphQLQueryGroup(group_type__value:
  "internal", query__ids: [<id>])` then answers each group's `node_metadata.created_by {
  id }`, a `CoreAccount`. A node the join restamped keeps a `created_at` no later than the
  generation that first read it. So a stamp read as the node's `created_at`, whenever its
  `updated_by` is the groups' creator ([D-058](decisions.md#d-058),
  [D-060](decisions.md#d-060)), reads a branch fresh after a hand generate. The contract
  tier is green with it, including `TestChecksInsideInfrahub`, and `TestStagesListAndPlan`
  marked after the join. *(10-26; `specs/017-dcim-reader/research.md` §5.14)*
- **The event log records a removal, read-only.** `InfrahubEvent(branches: [<branch>],
  since: <branched_from>, order:, limit: ≤50)` (above 50 is refused) lists
  `infrahub.node.deleted` for a deleted link, primary node the concrete kind
  (`NetworkLink`) with `endpoints` `REMOVED`, at an `occurred_at` between the mutation's
  send and its return, then `node.updated` on each endpoint interface (`link` `REMOVED`)
  and Infrahub's own query-group events. A read pinned at that `occurred_at` −1 µs already
  reads the removal. It answers alike through `/graphql/main` and `/graphql/<branch>`, in
  0.11–0.29 s a page. *(10-17)*
- **The event log keeps a deleted branch's events, and not what is older than it keeps.**
  A deleted test branch's events still read by its name, so `since: <branched_from>` is
  what keeps an earlier branch of the same name out. On this install the oldest event is
  `2026-10-13T16:48:04`, so `fylgja-fixture`, seeded on 10-08, reads none (`count: 0`);
  `since` defaults to 180 days back. *(10-17)*
- **Infrahub writes its own churn on intent's kinds.** After each generation, about three
  to four seconds after the checksums, `node.updated` events on every `NetworkDevice`,
  `NetworkInterface`, `NetworkLink` and `NetworkPlatform` change `member_of_groups` alone
  (the query groups' membership); an artifact's creation adds one on its device changing
  `artifacts` alone. The schema endpoint lists `member_of_groups`, `subscriber_of_groups`
  (kind `Group`) and `profiles` (kind `Profile`) on each Fylgja generic. So an event's kind
  does not tell intent from churn; the names of the fields an update changed
  (`attributes { name }`, `relationships { name }` on `NodeMutatedEvent`) do. *(10-17)*
  Those were the generics' kinds; on the library's the group join is the same churn,
  `CoreStandardGroup updated, members → …` at level 0 and one `member_of_groups →
  <group>` update a device at level 1. *(10-25; `specs/017-dcim-reader/research.md` §1.8)*
- **The log lags its writes, by a varying few seconds.** It is ordered strictly by
  `occurred_at`. A deleted link's own `node.deleted` is stamped just before the mutation
  returns, and its child events (each endpoint interface's `link` REMOVED, the query
  groups') up to about 25 ms after; each lands in the log later. A read half a second after
  the return held one child event and not the other; under the contract tier's load reads
  1.5 s after held none; on an idle host the first read to hold one was 6.6 s after the
  return. So a head read in that window does not see a removal yet. *(10-17)*
- **What the log names on the library's kinds.** Since a branch point it lists
  `infrahub.node.created`, `infrahub.node.updated`, `infrahub.group.member_added`,
  `infrahub.group.member_removed` and `infrahub.schema.updated`, whose `primary_node` is
  `null`. `NodeMutatedEvent` carries `primary_node { id kind }`, `related_nodes` (`{id,
  kind}`), `attributes { name }`, `relationships { name peer { id kind } }` and `level`. A
  creation names every attribute, `human_friendly_id` and `display_label` among them, and
  every relationship set; an attribute edit names the attribute; a tag added or removed is
  the device's `relationships: [tags → <tag>]`. A cable re-ended names its former port
  twice, 1.4–3 ms apart: the cable's event (level 0, `connected_endpoints → <port>`) and the
  port's (level 1, `connector → <cable>`), read 2 s after the write. A device created on a
  platform is `DcimPlatform updated` naming the device. *(10-25, 10-26;
  `specs/017-dcim-reader/research.md` §1.4, §1.8, §5.8)*
  - **The order of a device's creation.** The device's own `created` comes first, naming
    its platform and site as peers. About 1.6 ms later comes `DcimPlatform updated`
    (`devices → <device>`), then `LocationSite updated` 2.4 ms after that.
  - So the platform's event is the latest one naming the platform, about 70 ms after the
    platform's stamp moved.
  *(10-26; `specs/017-dcim-reader/research.md` §5.17)*
  - **A device type's writes, and a device's delete.**
    - **A device type created.** `DcimDeviceType created` names its manufacturer, and
      `OrganizationManufacturer updated` (`device_types` ADDED) follows about 3 ms later.
    - **A device type moved.** The type's own `updated` names both manufacturers, and each
      manufacturer's `updated` follows, `device_types` ADDED on the new one and REMOVED on
      the old.
    - **A device deleted.** `DcimDevice deleted` names its ports, platform and site, each
      REMOVED, and each port's own `deleted` follows, since `interfaces` is a `Component`
      edge. Then `DcimPlatform updated` (`devices` REMOVED) comes about 13 ms after the
      device's, and `LocationSite updated` after that.
    - So the manufacturer's or the platform's update is the latest event naming it.
    - Each relationship and attribute carries an `action` beside its name (`ADDED`,
      `REMOVED`, `UNCHANGED`).
    - The events were stamped 33 ms to 0.97 s after the write's stamp.
  *(10-26; `specs/017-dcim-reader/research.md` §5.18)*

**Checks, and the proposed change's view (015)**

Read in the image's source (`/source/backend`, `/source/python_sdk`,
`/source/frontend/app/src`) and seen live on the 26.04 development host, through throwaway
registrations served from the host's Docker bridge, the tree's checks and the host's
registration `fylgja`.

- **A check is an `InfrahubCheck` the task worker builds** with the source branch, the
  worker's own client (`http://infrahub-server:8000`, default branch `main`, timeout 60)
  and no initializer, so its `proposed_change_id` is always `""`: a check finds its
  proposed change by its source branch (`filters(kind="CoreProposedChange", branch="main",
  source_branch__value=…, state__value="open")`, with `is_draft` readable). `query` is
  required, a cheap `{ Branch { name } }` will do, and its answer reaches `validate`
  unwrapped. It has `log_error` and `log_info` and no `log_warning`; the check's message is
  its entries as `Message: … Level: …` blocks, a long multi-line entry stored whole, with
  `Check succesfully completed` (sic) appended on a pass, and an exception out of
  `validate` reads `FATAL Error/n:` (sic), conclusion `failure`. *(10-17)*
- **The module is `commits.<sha>.infrahub.checks.<stem>`** for a check at
  `infrahub/checks/<stem>.py`, with the worktree's root on `sys.path`, so the repository's
  `infrahub/` directory does not shadow Infrahub's own package, and a relative import
  between two files of one directory resolves with no `__init__.py`. *(10-17, 10-19)*
- **A check's bound is its own `timeout`**, read from the class at import into
  `CoreCheckDefinition.timeout` (a change needs a re-import) and applied as a Prefect task
  timeout; a check longer than the SDK's default of 60 s runs to its own and the pipeline
  waits for it. `INFRAHUB_MISC_MAXIMUM_VALIDATOR_EXECUTION_TIME`, which the published
  Compose file sets, is read by nothing. *(10-17)*
- **A timed-out check leaves its validator `in_progress` for ever**, conclusion
  `unknown`, with no `CoreStandardCheck`, across re-runs, and nothing in the task workers'
  logs names it. The merge rule wants every validator but the data validator concluded
  `success` (`Unable to merge proposed change containing failing checks`), so it blocks the
  merge. Past the cut an `await` runs nothing; a blocking call made in a `BaseException`
  handler still reaches a server. *(10-17)*
- **Validators and checks are reused across runs**: one `CoreUserValidator` per
  repository and definition, labelled `Check: <name>`, restarted at each run, holding one
  `CoreStandardCheck` rewritten in place (`message`, `conclusion` by `passed`, `severity`
  `info` or `critical`). *(10-17)*
- **When the checks run.** At a proposed change's creation, draft or not, every registered
  repository's checks run (the user validators 43–54 s after the create, a data-only change
  included); `CoreProposedChangeRunCheck(check_type: USER)` runs them again alone, in about
  20 s, with no artifact validator and no checksum moved; `is_draft` set to `false` starts
  nothing. A run runs, for each repository in the branch's diff, that repository's checks
  as the source branch sees them: a definition imported after the branch's creation never
  runs on it, and a branch made while a since-deleted registration lived still lists its
  definitions, `repository` null, which never run. *(10-17, 10-19)*
- **An import reaches the task workers after it.** `InfrahubReadOnlyRepositoryImportLastCommit`
  on `fylgja` made `fylgja_branch` and both check definitions 15–21 s after the call, and
  `in-sync` at 25 s (10-19), or `in-sync` at 7.3 s (10-20); a second `Import last commit`
  and `Pull latest commit` follow each import, unasked and unattributed. Each task worker
  then fetches the commit in a flow of its own (`Fetch git repository <name> on <worker>`),
  both done 13–17 s after the import's own flow; a run-check 6 s after an import failed on
  the worker that had not fetched (`SHA … could not be resolved`) and started no
  validator. *(10-17, 10-19, 10-20)* On 10-26 the import's task ran 28.4 s, `commit`
  moving at its end, and the second task worker held the commit 4.6 s after the flow
  ended. *(10-26; `specs/017-dcim-reader/research.md` §5.4)* **A worker's pull can fail
  and is never retried.** On 10-27 the import of `cb0a8d7` moved `commit` 25 s after the
  call and task-worker-2 held `commits/cb0a8d7…`; task-worker-1's `Pull latest commit`,
  21 s after the call, failed `CommitNotFoundError` (`git worktree add … fatal: invalid
  reference`): its clone's `origin/main` was still the previous import's commit, and no
  `Fetch git repository` flow ran on either worker. Nothing repaired it: a second import
  (`ok: true`, `commit` already the head) and a `docker restart` of the worker each left
  its clone unfetched. A `git fetch origin` in that worker's `/opt/infrahub/git/<repository
  id>/main`, then the `git worktree add` its flow had failed, gave it the commit's checks
  byte for byte. After an import, read both workers' `commits/` before relying on a
  check. *(10-27)*
- **A query's `models` are derived at its import**, from every kind it reaches, generics
  through fragments included: `device_config` over the library reads `DcimCable`,
  `DcimConnector`, `DcimDevice`, `DcimEndpoint`, `DcimGenericDevice`, `DcimInterface`,
  `DcimPlatform`, `InterfacePhysical`, `InterfaceVirtual` and `OrganizationManufacturer`,
  and `fylgja_branch`, which reads no kind, has none. It imported with no error on a
  `main` holding the library. A kind's removal leaves a stored query's `models` naming it
  until the next import. *(10-26; `specs/017-dcim-reader/research.md` §5.4, §5.10)*
- **`CoreGraphQLQuery` and `CoreCheckDefinition` names are unique on `main`**: a second
  registration declaring a query name the first holds goes `error-import` (`Violates
  uniqueness constraint 'name'`, 422) and imports nothing at all, its checks included.
  Deleting a `CoreReadOnlyRepository` deletes its definitions and its query. *(10-17,
  10-19)*
- **The task worker's client writes comments and threads** on `main`. A
  `CoreChangeComment` (`text`, `change`; no `created_at` or `created_by` field) is counted
  in `total_comments` and **shown nowhere in the proposed change's overview**, which lists
  the proposed change's events: a comment appears there only in a thread, through
  `infrahub.proposed_change_thread.created`, which the creation of a `CoreChangeThread` (or
  an object, artifact or file thread) on the proposed change emits, read within 16 s. The
  view's *Add comment* writes a `CoreChangeThread` (`label` `Conversation`, `resolved`
  false), then a `CoreThreadComment` in it, and deletes the thread when the comment fails.
  Written so by a check, the thread is authored by the task worker's account, `Agent`, and
  its markdown renders with GFM tables, single line breaks kept (`remark-breaks`) and a
  Mermaid block with no error banner; a check's message shows as plain text under
  `Message`. *(10-20)*
- **How the view renders a comment's markdown**, tried in the comment box's *Preview*,
  which renders without posting:
  - A Mermaid node's `style` attribute is dropped, so a `style sN …` line does nothing. The
    drawing's own `<style>` block and its classes are kept, so a `classDef` with `class sN …`
    restyles a node.
  - A drawing is scaled to the thread's width. Its zoom buttons scale it inside a frame
    that does not grow, so a zoomed drawing is cropped.
  - A code block's `code` is `display: inline` with `padding-left: 5.6px`, so its first
    line alone is indented, whatever the info string (none, `text`, `sh`, `console`) and
    for an indented block too. *(10-21)*
- **A comment's text holds at most 4,092 bytes**: a text attribute is refused, `Text
  attribute length should be less than 4096 characters.`, when `3 + len(utf8) >= 4096`
  (`core/attribute.py`), a thread's comment as a plain one. *(10-20)*
- **`CoreProposedChangeDelete` takes its comments, threads with their comments,
  validators and checks with it**, a validator stuck `in_progress` included, since each
  hangs on it by a cascading relationship. *(10-17, 10-20)*

**The step (019)**

- **A step to a snapshot sends Infrahub nothing.** With a counting reverse proxy in front
  of the Infrahub `INFRAHUB_ADDRESS` names, through which every reader and stamper the
  server builds sends its requests, `stage plan` on a throwaway branch sent 33 requests and
  the `intent read --out --stage` that took the mark's snapshot sent 24. The step's
  snapshot target, built by the server's own code from that file and then from its id
  (resolved, its scope and mapping held to the record's, compiled and filed), sent 0 each,
  named the plan's snapshot and compiled to its bundle. *(10-28;
  `TestStepTargetFromASnapshot`, `specs/019-following-removed/tasks.md` T036)*

**Artifacts and repositories**

- Infrahub regenerates no artifact on its own: not on a data change, a group join or a
  commit import. Generation is `POST /api/artifact/generate/<definition id>?branch=`,
  which returns `null` in about 0.1s; three artifacts were `Pending` at 6.4s and `Ready`
  at 7.6s. *(09-18)*
- Nothing marks an artifact stale. A regeneration to identical bytes writes nothing, not
  even a storage id; one that changes bytes returns before the rewrite, and the artifact
  stays `Ready` throughout, so wait for its checksum to move. *(09-18)* The proposed
  change's pipeline is the one thing that regenerates unasked (above). *(10-16)*
- **An artifact's stamp is `checksum.updated_at`.** A first generation stamps
  `checksum`, `storage_id` and `status` alike; a regeneration moves the checksum's and
  the storage id's together and leaves `status.updated_at` where it was (`Ready` →
  `Ready` is no change: unmoved through four regenerations). A pinned read returns the
  stamps of the versions in force at `at`, not the current ones, so intent newer than its
  artifact at a pinned instant is read from Infrahub's own stamps: the staleness rule of
  [D-045](decisions.md#d-045). *(10-16)*
  Reported upstream: **dropped**, on the operator's word. Infrahub documents both
  behaviours as the design ([opsmill/infrahub#10437](https://github.com/opsmill/infrahub/pull/10437):
  no data or membership change starts a run, and there is no staleness indicator).
  *(10-08)*
- **The stamps read through the generics**, unpinned and pinned: `FylgjaDevice`'s
  attributes' `updated_at`, its `platform` edge's `properties.updated_at`, each interface's
  edge and attributes, its `link` edge (the cabling's stamp) and its `parent` edge, whose
  `properties.updated_at` is `null` when unset; `FylgjaLink.endpoints.edges.properties
  .updated_at`; `FylgjaContract.version.updated_at`. `artifacts` must be spread `... on
  CoreArtifactTarget`: selected on `FylgjaDevice` it is refused `Cannot query field
  'artifacts' on type 'FylgjaDevice'`. *(10-16)* The platform node's own attributes, through
  `FylgjaPlatform`, carry `updated_at` too: the fixture's `vendor`, `nos`, `version` and
  unset `model` all read `05:47:56.812998`, the platform's creation, unpinned and pinned at
  the head, before every device's stamp, and pinned at exactly a device's stamp the edge
  and the node are `null`. *(10-17)* Those were the generics, retired at 017
  ([D-057](decisions.md#d-057)); on the library's kinds a node's stamp is its wrapper's
  (next).
- **A node's stamp is its `Edged*` wrapper's `node_metadata.updated_at`**, the latest
  stamp among the node's current attributes and relationship edges, on the library's
  kinds. An attribute's `updated_at` moves on a change of its own value alone, and a second
  write of the same value moves nothing. The wrapper follows: a description edit moved it,
  a tag added moved it again, and the tag's removal moved it **back** to the edit's stamp,
  since a removal leaves no stamp; a no-op write moved nothing. An interface's edit moved
  the interface's and not its device's, a cable's the cable's, a manufacturer's the
  manufacturer's and not its platform's. A device created on a platform moves the
  platform's, by its reverse `devices` edge, and an edit of that device then moves nothing
  of the platform's. A node's first join of a group moves its stamp; an upsert of a member
  already there moves nothing (`RelationshipManager.update` keeps the peer). A port's
  `connector { node_metadata { updated_at } }` reads its cable's own stamp, which the
  cable's edit moves while the edge's `properties.updated_at` stays at the cable's
  creation. A cable re-ended to one end moves nothing stamped: the former port reads
  `connector { node: null, properties { updated_at: null } }`. The root filter
  `node_metadata__updated_at__after` answers by the same stamp. Pinned at exactly a
  write's stamp the wrapper reads the previous latest while the attribute reads its new
  value with its old stamp; one microsecond later both read the write. Pinned at exactly a
  node's creation its edge is absent, and one microsecond later present with `updated_at`
  equal to `created_at`. *(10-25, 10-26; `specs/017-dcim-reader/research.md` §1.4, §1.7,
  §5.8)*
  - **A device type and a device's delete.**
    - A `DcimDeviceType` created under a manufacturer, or moved onto it, moves the
      manufacturer's stamp to the write, through the reverse `device_types` edge.
    - Moved off it, the type moves the former manufacturer's stamp **back**, as a removal
      does.
    - A device deleted (`DcimDeviceDelete`) moves its platform's stamp to the delete, as a
      cable's delete stamps its ends. Its manufacturer's stamp does not move: the device
      reaches it only through the platform.
  *(10-26; `specs/017-dcim-reader/research.md` §5.18)*
- **A generation stamps each device's artifact apart.** The fixture's three checksums read
  `…38.543597`, `…38.544382` and `…38.544482`, each after its own device's intent, so the
  latest stamp of a branch is one device's checksum. A regeneration that changes bytes
  lands device by device, seconds apart: a generate returned in 0.44 s, and the three
  checksums moved 4.6–7.5 s after a mark made as it returned; after `-add-link`, n1's and
  n3's moved 7.6–10 s after it returned. *(10-16, 10-17)*
- `CoreArtifact.checksum` is the MD5 hex of the bytes `GET /api/storage/object/<storage_id>`
  serves. A listing pinned to `at` gives the storage id as of `at`, and the older object
  is still served. *(09-18)*
- **The pipeline renders as `infrahubctl render` does.** `.venv/bin/infrahubctl render
  <transform> device=<d> --branch <b>` (SDK 1.23.2, read-only) gave bytes whose MD5 is each
  device's live `CoreArtifact.checksum`; and the renders of 017's template predicted from
  it held byte for byte once Infrahub's task workers rendered the imported commit. The
  SDK's file-based environment has `trim_blocks`, `lstrip_blocks` and the default
  `Undefined`, so an attribute read on a null node fails the render. Its filter allowlist
  (`infrahub_sdk/template/filters.py`) admits `sort` in the worker's and the local context,
  not the core's, and artifacts render in the worker. *(10-25, 10-26;
  `specs/017-dcim-reader/research.md` §5.1)*
- A render for a device the query does not find stops in the template (`list object has no
  element 0` at `data.DcimDevice.edges[0]`), the query answering `edges: []` with no
  error; the pipeline renders the target group's members alone. A member that the imported
  template's query cannot read (a template reading `FylgjaDevice` over a `DcimDevice`)
  has its artifact end `Error`. *(10-26; `specs/017-dcim-reader/research.md` §5.2, §5.7)*
- **`infrahubctl render <transform> device=<name> --branch <b> --out <file>` renders the
  tree's template against a live branch** before any import: it reads `.infrahub.yml`'s
  `jinja2_transforms` entry by the transform's name (`srlinux_device_config`, not the
  query's), runs the tree's query on the branch and renders the tree's `.j2` locally. A
  port's `ip_addresses { edges { node { address { value } } } }` inside `... on
  InterfacePhysical` answers with no further fragment. *(infrahub-sdk 1.23.2, 10-29;
  `specs/016-probes/research.md` §5, R-18)*
- **The template's address lines** (016), rendered and then pushed: an addressed SR Linux
  port gains `subinterface 0 ipv4 admin-state enable` and `subinterface 0 ipv4 address
  <a/l>` after its `subinterface 0 admin-state enable`; an EOS device gains `ip routing`
  after `hostname`, and each addressed port `no switchport` and `ip address <a/l>` after
  `no shutdown` (a second address renders `secondary` on EOS, where a plain second `ip
  address` replaces the first, and plainly on SR Linux, in value order). A device with no
  addressed port renders byte for byte what the template rendered before: the fixture's
  three before its re-seed, their stored checksums. SR Linux 24.7.1 accepts both lines in
  the replace (pushes 2.8–2.9 s; `info from running` reads `ipv4 { admin-state enable
  address 10.1.10.1/24 { } }`), and cEOS 4.32.0.2F all three in the eAPI replace (pushes
  2.3–2.5 s; `show running-config` reads them). *(10-29; `specs/016-probes/research.md`
  §5, T034, T043)*
- **An artifact moves for every non-management port, cabled or not**, since the template
  renders each, so a port added to a seeded device changes its artifact before any address
  does, and its address changes it again. The fixture re-seeded under the imported
  template carries each device's two addressed ports, 1,327 bytes an artifact (`n1`
  `28881d3d…`, `n2` `092ca9f6…`, `n3` `95a3a917…`, the MD5 of the stored bytes); the
  reader reads no address, so the CTM moves for the ports alone. *(10-29;
  `specs/016-probes/research.md` §5, R-19, T042, T044)*
- A branch created with `sync_with_git: false`, as every test branch is, keeps the
  repository commit it forked with: a template commit reaches it only by `BranchRebase`,
  then a generate. *(09-18)*
- A `BranchCreate` straight after a `BranchDelete` of the same name can fail with
  `graphql: None`; run it again. *(09-18)* The branch is created all the same: the
  delete's task purge, which picks runs by the branch's name, deletes the new create's
  flow run while the API polls it (2 of 12 re-creates, 1.11.2). Reported upstream:
  **dropped**, on the operator's word; the same cause is
  [opsmill/infrahub#10694](https://github.com/opsmill/infrahub/pull/10694)'s, from the
  delete's side. *(10-08)*
- The repository sync (Prefect's `git_repositories_sync`) runs every minute with one
  concurrency slot, so one run stuck `PENDING` blocks every import until it is cancelled.
  *(09-18)*
- `CoreRepository` and `CoreReadOnlyRepository` are both branch-agnostic and take an
  optional credential. A public HTTPS repository registers as a `CoreReadOnlyRepository`
  on `main` with no credential: the clone and the connectivity check pass, and the import
  then reads the repository's `.infrahub.yml`. *(10-06)* This repository's registration
  reached `active`, `online` and its query, transform and definition on `main` 22–48 s
  after its creation, with no sync stuck: 24 s on the 24.04 development host, 22 s on
  GitHub's runner, 48 s and 43 s on fresh 26.04 VMs and 47 s on the 26.04 development
  host. *(10-07; 26.04, 10-07, 10-08)*
- A `CoreReadOnlyRepository` imports at its creation, its `commit` the head of its `ref`,
  and afterwards only when asked: a commit pushed to `ref` is not picked up by the minute
  sync (three minutes on, `commit` unmoved and `sync_status in-sync`), nor by a
  `CoreReadOnlyRepositoryUpdate` writing `ref`'s own value (`ok: true`, nothing written).
  `InfrahubReadOnlyRepositoryImportLastCommit(data: {id})` imports the head of `ref`:
  `commit` moved 8s after the call. *(10-07)*
- Deleting a `CoreRepository` deletes the query, the transform and the artifact
  definition its import made; a later registration's import makes them anew, with new
  ids. *(10-07)*

**Its published Compose**

- `https://infrahub.opsmill.io/1.11.2` serves the Compose file whose three Infrahub
  services read `${VERSION:-1.11.2}`; the root URL serves the latest release's, and
  `/1.11` is 404. *(10-07)*
- The published file carries defaults, readable by anyone, for
  `INFRAHUB_INITIAL_ADMIN_TOKEN`, `INFRAHUB_INITIAL_AGENT_TOKEN` and
  `INFRAHUB_SECURITY_SECRET_KEY`, and `INFRAHUB_INITIAL_ADMIN_PASSWORD` defaults to
  `infrahub`. It publishes `8000` (the server), `2004` and `6362` (Neo4j) and `15692`
  (RabbitMQ's metrics) on every interface. Its seven services' images total about
  2.1 GB. *(10-07)*
- Compose 5 runs it with an override beside it: 5.5.1 on 24.04, and 5.6.0 on 26.04, which
  warns once per start that the task worker's `deploy.mode is only honored in Swarm mode`.
  *(10-07; 26.04, 10-07)*
- **Inheriting from two generics.** A node may inherit two generics that declare the same
  attribute or relationship only when the two declarations match in every normalized
  property. `POST /api/schema/check` refuses the first mismatch it meets and names it: a
  `role` that is `Text` in one and `Dropdown` in the other (`must be the same kind`), then a
  relationship's `peer`, `identifier`, `min_count` and `on_delete` in turn, each as
  `inherits from 'A' & 'B' with different '<property>' on the '<name>' relationship`. The
  default identifier of a relationship with none written is the two kinds lowercased and
  joined (`fylgjadevice__fylgjaplatform`), so two generics peering on different kinds never
  agree on it unless one writes the other's. Verified by checking OpsMill's schema library
  (`base/dcim.yml` with its ipam, location and organization bases, read on 10-16) against
  `main`'s Fylgja generics: the four files load beside them unchanged, and the DCIM kinds
  inherit the generics only once `DcimGenericDevice`'s `platform` and `interfaces` are
  exact copies of `FylgjaDevice`'s and `DcimDevice` drops its `role`. *(10-16)* 017
  removed the generics instead ([D-057](decisions.md#d-057)).
- **Two Parent relationships are refused, even when identical.** A node inheriting
  `device` as `kind: Parent` from two generics fails `schema/check` with `Only one
  relationship of type parent is allowed, but all the following are of type parent:
  ['device', 'device']`, before the declarations are compared, so matching them does not
  help: one generic must not declare it. Dropping it from `DcimInterface` then fails its
  `uniqueness_constraints`, `human_friendly_id` and `order_by`, which reference `device`,
  until those are rewritten too. *(10-16)*
- **A schema check applies a submitted generic as a patch.** Submitting `FylgjaDevice`
  without its `role`, or `FylgjaInterface` without its `device`, removes nothing: the
  check reports the live attribute still present (and still conflicting). A removal needs
  `state: absent` on the member, as a load does. A check with `?branch=main` validates
  against `main`'s schema and writes nothing; `main`'s hash read `0aba31a5…` before and
  after seven rounds of checks. *(10-16)*

**The embedded Prefect**

- **Infrahub's Prefect is a separate server, on the Compose network alone.** The service
  `task-manager` runs `uvicorn --host 0.0.0.0 --port 4200 --factory
  infrahub.prefect_server.app:create_infrahub_prefect` against `task-manager-db` (Postgres
  18) and publishes no port on the host, whose `4200` stays free; from inside a task
  worker its API answers `GET /api/version` `0.8.4` and `POST /api/flows/filter` (Infrahub's
  own flows) with no authentication. So an orchestrator's Prefect server on the host
  collides with nothing. *(10-29)*
- **The task workers are Prefect workers pointed at it**: `prefect worker start --type
  infrahubasync --pool infrahub-worker` with `PREFECT_API_URL=http://task-manager:4200/api`.
  Inside `fylgja-infrahub-task-worker-1`: Prefect 3.7.5 on Python 3.14.7, `httpx` 0.28.1,
  `requests` 2.34.2, `infrahub_sdk` 1.11.2 (the image's own), `pydantic` 2.12.5, `yaml`
  6.0.3, no `aiohttp`. A check that drives another Prefect server must not use the worker's
  Prefect client, whose environment names Infrahub's; plain HTTP with `httpx` reaches it.
  *(10-29)*
- **A recreated task worker has no checkout, and a user check cannot make one for a
  read-only registration.** The task workers mount no volume, so `up -d` empties
  `/opt/infrahub/git`. The minute's `git_repositories_sync` reads `CoreRepository` alone and
  never clones a `CoreReadOnlyRepository`; `run_user_check` (`infrahub/git/tasks.py:1193`)
  passes `InfrahubKind.REPOSITORY` hard-coded, so on a worker with no clone it asks for a
  `CoreRepository` named `fylgja`, gets `NodeNotFoundError`, and the validator stays
  `in_progress`: the first proposed change opened after a recreate had both user validators
  `in_progress` for ever, and the fixture tool's `-propose` lapsed its 6 minutes. The
  artifact check, which passes the registration's own kind, cloned it on both workers about
  90 s on (`Initialized the local directory for fylgja because it was missing`), and a
  re-run of the user checks then concluded both; every later check ran. **This is a
  defect of Infrahub 1.11.2**, not a setting: a user check of a read-only registration on a
  worker with no clone can never succeed, since the kind it asks for is fixed. Fylgja does
  not work around it; a recreate is followed by an artifact check or a re-run of the user
  checks, and the report upstream is the operator's (`CLAUDE.md`, State). *(10-30;
  `specs/020-core-orchestrator-split/research.md` §5.28)*
- Infrahub removes no clone of a deleted repository: each task worker keeps the checkout
  directory of a registration a test made and deleted. *(10-30)*

**The adapter (020)**

The orchestrator's Infrahub adapter, Python 3.14.4 in `orchestrator/.venv`.

- **`infrahub-sdk` 1.23.2's `execute_graphql` rewrites `at`**: its `_graphql_url` parses
  the value with the SDK's `Timestamp` and sends `to_string()`, UTC with six fractional
  digits (`2026-10-29T12:00:00.123456+02:00` goes as `…T10:00:00.123456Z`, `12:00:00Z` as
  `12:00:00.000000Z`, a zone-less `2026-10-29 12:00:00` as UTC), so the operator's `at`
  is not what reaches Infrahub. The branch goes in the path, `quote(branch, safe='')`.
  Its exceptions drop the HTTP status: a 401 or 403 is `AuthenticationError` with
  Infrahub's messages, a 404 `URLNotFoundError` naming the URL alone (not GraphQL's
  `Branch: <b> not found.`), and any other status whose body holds `errors`, 400 and 422
  among them, a `GraphQLError` as if under 200. The client's `custom_recorder` sees every
  response, and its `sync_requester` replaces the transport. Its settings read
  `INFRAHUB_*` from the environment: a `username` there turns on password login unless
  `username` is passed explicitly. The adapter's `client.py` overrides `_graphql_url` to
  send `at` verbatim, form-encoded, and reads the status from a recorder; against this
  host's 1.11.2 a pinned `at` with a `+00:00` offset and six digits read `main`, and a
  missing branch answered 400 from the schema endpoint and 404 from GraphQL, both
  `Branch: <b> not found.`. *(10-29; `specs/020-core-orchestrator-split/research.md`
  §5.10)*
- **The snapshot id in Python**: SHA-256 of the hashed body as
  `json.MarshalIndent(body, "", "  ")` plus a trailing newline, keys in the contract's
  order and the lists sorted as `ctm.Normalize` sorts them; Python's `json.dumps(body,
  indent=2) + "\n"` gave `1a29d9dc…` and `6a4a91ff…`, the ids the core's server filed for
  the NetBox spike's two snapshots; the compact and sorted-key forms do not. *(10-29)*
- **Python writes Go's bytes once five characters are escaped**: Python 3.14.4's
  `json.dumps(…, indent=2, ensure_ascii=False)`, with `<`, `>`, `&`, U+2028 and U+2029
  translated to Go's `\u003c`, `\u003e`, `\u0026`, `\u2028`, `\u2029` and a lone surrogate
  to U+FFFD, plus a newline, is go1.26.8's `json.MarshalIndent` byte for byte: control
  characters both write as `\b`, `\f`, `\n`, `\r`, `\t` or a lowercase `\u00XX`, DEL and
  every character beyond ASCII both raw. The default `ensure_ascii=True` escapes beyond
  ASCII, which Go never does, so the form above is exact only for ASCII text. The adapter's
  `snapshot.py` gave the id of all 203 snapshots under this host's `local/snapshots/` (133
  of version `1`), wrote the 70 of version `2` back byte for byte, and did the same for the
  four wrapped fixtures and a Go-written file holding every such character
  (`orchestrator/tests/data/snapshot/escapes.json`). *(10-29;
  `specs/020-core-orchestrator-split/research.md` §5.9)*
- **The adapter's read of the fixture** (`python -m fylgja_orchestrator.adapters.infrahub
  read --branch fylgja-fixture`) took 4.0–4.6 s wall, the interpreter's start included,
  three times, each to snapshot `0940db93…` (3 devices, 15 interfaces, 3 links, 3
  artifacts, schema `8a4f3bee…`). Beside the Go reader before it went, Go's through the
  server and the adapter's in the test's process, the fixture read 3.3 s and 2.8 s, a
  pinned read 2.1 s and 2.3 s, and seven reads 39.8 s and 34.7 s in all, every file equal
  byte for byte once its two `observed_at` lines were equalised. *(10-30;
  `specs/020-core-orchestrator-split/research.md` §5.15, §5.28)*

---

## NetBox 4.7.4 with netbox-branching 1.2.2

Verified on the development host on 10-29, before F-025 was briefed, by a throwaway spike:
netbox-docker 5.1.1 (`docker.io/netboxcommunity/netbox:v4.7.4-5.1.1`, `postgres:18-alpine`,
`valkey/valkey:9.1-alpine` twice) as the Compose project `fylgja-netbox` under
`local/netbox/`, on loopback port 8800 (8000 is Infrahub's), with a one-layer image adding
`netboxlabs-netbox-branching==1.2.2` by `uv pip install` (`uv` is at `/usr/local/bin` in
the image, the venv at `/opt/netbox/venv`), and `configuration/plugins.py` carrying
`PLUGINS`, `DATABASES` in the plugin's `DynamicSchemaDict`, `DATABASE_ROUTERS` and
`EVENTS_PIPELINE` as NetBox Labs' netbox-docker page gives them; and a Python adapter of
about 170 lines that read a branch, mapped it onto CTM 2 and wrote snapshot 2 files, which
the core took unchanged. The adapter and the install files are not in this repository.

**Running it**

- **The first start takes about 11 minutes**, nearly all of it `manage.py migrate`;
  Compose marks the `netbox` container unhealthy after its 90 s start period and never
  starts `netbox-worker`, whose dependency failed, so a second `up -d` once `netbox` is
  healthy starts it. The worker is needed: a branch is provisioned by an RQ job. The
  plugin image built in 4 s. *(10-29)*
- **The stack idles at about 1.3 GiB**: `netbox` 1.03 GiB, `netbox-worker` 173 MiB,
  `postgres` 89 MiB, the two valkeys 4 MiB each; the host's `free -m` lost 1.2 GiB of
  available memory with it up. *(10-29)*
- **API tokens are v2**: a token is presented as `Authorization: Bearer nbt_<key>.<token>`,
  the 12-character `key` and 40-character `token` that `POST /api/users/tokens/provision/`
  (username and password) returns once, `version 2`; `Authorization: Token <token>` is
  parsed as a v1 token and refused `Invalid v1 token`. v1 tokens are deprecated since 4.6
  and go in 5.0. netbox-docker 5.1.1's entrypoint ignores `SUPERUSER_API_TOKEN` ("No API
  token was created…", with the variable in the container's environment). *(10-29)*
- **`/api/status/`** names `netbox-version 4.7.4`, `django-version 6.1.2`,
  `python-version 3.14.4` and `plugins {netbox_branching: 1.2.2}`. *(10-29)*

**Branches**

- **A branch is selected by `X-NetBox-Branch: <schema_id>` on REST and GraphQL requests
  alike**, and on GraphQL by `?_branch=<schema_id>` too; the `schema_id` is eight
  lower-case alphanumerics (`p9d46y7z`), read from the branch's own representation. A
  header naming no branch is `400 Invalid branch identifier`; one naming a branch that is
  `new`, `merged` or `archived` is `400 Selected branch is not ready.`, for `render-config`
  as for any read, so a merged branch cannot be read again. *(10-29)*
- **`POST /api/plugins/branching/branches/ {name}`** answers `status {value: new}`, then
  `ready` once the worker has provisioned it, 5.7 s on this host (a `provisioned` row in
  `branch-events`); `merge` with `commit: true` is a 2.2 s job that leaves the branch
  `merged`, `main` carrying its changes and a `merged` event; `archive` refuses an unmerged
  branch ("Only merged branches can be archived") and `revert` an archived one ("Only
  merged branches can be reverted"); `sync` with `commit: false` enqueues a `Sync branch`
  job. The plugin's API root lists `branch-events`, `branchable-models`, `branches` and
  `changes`; there is no `changediffs`. *(10-29)*
- **`render-config` honours the header**: with `n1 ethernet-1/2` set `enabled: false` on
  the branch, `POST /api/dcim/devices/1/render-config/` with `Accept: text/plain` renders
  `admin-state disable` under the header and `enable` without it, in about 0.5 s;
  `Accept: application/json` answers `{configtemplate, content}`; a template edited under
  the header is what the branch renders, `main` rendering its own. The template is
  resolved device, then role, then platform. *(10-29)*
- **GraphQL answers the mapping's fields under the header**:
  `device_list(filters: {tags: {name: {exact: "fylgja"}}})` with `name`, `platform {name
  manufacturer {name}}` and `interfaces {name type mgmt_only enabled cable {id}
  ip_addresses {address} connected_endpoints {... on InterfaceType {name device {name}}}}`,
  and `cable_list {id status terminations {cable_end termination {__typename ... on
  InterfaceType {id name device {name}}}}}` (`CableTerminationTerminationType` is a
  union); three devices and every cable in about 0.45 s. *(10-29)*
- **Branch-aware models** (`/api/plugins/branching/branchable-models/`, a bare list of
  120) include `extras.configtemplate`, `extras.configcontext`, `extras.journalentry`,
  `extras.tag` and `extras.taggeditem`; `extras.eventrule` and `extras.webhook` are
  exempt, so a webhook posted under the header lands on `main`, while a journal entry
  posted under it lives in the branch alone (200 there, 404 on `main`). A journal entry
  attaches to the Branch object itself (`assigned_object_type: netbox_branching.branch`,
  201). `OPTIONS /api/extras/event-rules/` lists no `branch_*` event types among its
  choices, though the plugin documents five. *(10-29)*
- **A branch's state has a marker but no earlier reading**:
  `/api/plugins/branching/changes/?branch_id=` rows carry `action`, `object_type`,
  `object_id`, `original_data`, `modified_data`, `current_data` (`main`'s value now),
  `diff` and `conflicts`, and `/api/core/object-changes/` under the header is the branch's
  own changelog, with monotonic ids; nothing reads the branch as it was at an earlier
  change. *(10-29)*
- **NetBox Changes (`netbox_changes`), the change-request layer, is a Pro and Premium
  feature** (NetBox Labs' v1.0 post, May 2026), not installable on the community edition;
  there the merge gate is the plugin's `protect_main` and permissions. Read, not run.
  *(10-29)*

**Through the core**

- **An adapter's snapshot builds and steps a twin unchanged**: the file, with no
  `snapshot_id`, passed `twin compile --snapshot` and `twin create --dry-run --snapshot`
  (verdict clear), `twin create --snapshot` built the three-node SR Linux twin in 43–44 s
  (deploy 39–40 s, pushes 3.0–3.1 s), `twin show` reads `read by netbox-dcim from
  netbox@4.7.4`, kind `snapshot`, and `twin verify` held 15 assertions over 3 nodes. The
  adapter wrote `mapping {reader: netbox-dcim, model: netbox, version: 4.7.4}`, no `at`,
  `schema_hash` the SHA-256 of the GraphQL introspection result, `mgmt_only` and `enabled`
  exact (no approximation rows) and one fidelity row, `omitted platform.version
  dcim.Platform`; a read with three `render-config` calls took 2.2 s. *(10-29)*
- **`twin step --snapshot <branch file> --wait 2m`** from that twin found nodes and links
  unchanged and one artifact moved, pushed `n1` in 3.3–3.4 s and settled in 2.9–3.4 s,
  8–9 s end to end; `verify` then skipped the disabled port and its link's two ends and
  failed none, and gNMI read `n1`'s `admin-state` `enable` before and `disable` after.
  A twin's node addresses are reassigned at each create: read them from `twin show`.
  *(10-29)*
- **The snapshot id reproduces outside Go**, as the spike's Python gave it, which 020's
  adapter then took: the rule is under Infrahub 1.11.2, **The adapter (020)**. *(10-29)*

---

## containerlab 0.79.0

- `clab deploy` creates `clab-<lab>/` beside the topology file, or under
  `CLAB_LABDIR_BASE` when it is set, mostly root-owned; `clab destroy --cleanup` removes
  it without sudo in about 1s. *(09-15; 26.04, 10-08)*
- containerlab applies an SR Linux `.cli` startup snippet over the default configuration
  and commits it itself, so the snippet carries no `commit` line; LLDP and gNMI are on by
  default. *(0.77, 09-15; 26.04, 10-08)*
- `inspect --all --format json` gives each container's `absLabPath`, always absolute;
  `labPath` is relative to the caller's directory when the topology file is under it, so
  decode `absLabPath`. *(09-16; 26.04, 10-08)*
- `inspect` reports an exited container's address as `N/A`. *(10-03; 26.04, 10-08)*
- `/usr/bin/containerlab` is setuid root, so the kernel clears `Pdeathsig` when it starts,
  and a `clab deploy` outlives a worker killed with `kill -9`: stop a stray `clab` before
  acting. *(09-15; 26.04, 10-08)*
- A plain re-deploy of a partial lab skips containerlab's post-deploy and leaves the
  bootstrap unapplied; deploy over a partial lab with `--reconfigure`. *(09-15; 26.04,
  10-08)*
- `clab deploy --dry-run --format json` prints a plan on stdout (`added-nodes`,
  `deleted-nodes`, `recreated-nodes`, `restarted-nodes`, `added-links`,
  `deleted-endpoints`, `node-change-reasons`) and writes nothing. *(10-01; 26.04, 10-08)*
  On 26.04 every plan also carried `started-nodes`, and no `deployed-lab: false` plan
  carried `lab-name`, which every such plan recorded on 24.04 at the same release
  (`internal/lab/testdata/`) does. Fylgja reads `lab-name` only in its refusal of a
  `deployed-lab: true` answer, which still carries it. *(26.04, 10-08)*
- That dry run must run with `CLAB_LABDIR_BASE` at the running lab's directory:
  containerlab decides a recreate from `.state.clab.yaml` there, so from any other base an
  image or kind change is invisible. *(10-01; 26.04, 10-08)* Refined on 26.04: from
  another base while the lab runs, the answer is still a plan, `deployed-lab: false`, that
  misses the change (a kind change showed only as a restart, for `deleted endpoint`);
  `deployed-lab: true` with `lab-name` and every list `null` comes only when no container
  of the lab exists. *(26.04, 10-08)*
- `clab deploy` without `--reconfigure` reconciles a running lab: a link added or removed
  re-cables SR Linux live and restarts a cEOS node in place (3.5–14.6s for the apply); a
  node added is created with containerlab's post-deploy (31.3s); an image or kind change
  is a recreate (20.3s, 58.7s with a far cEOS end's restart); a node removed is deleted
  (0.85s); nothing to do takes 0.39s. *(10-01)* On 26.04: an SR Linux link added 0.66s
  and removed 0.59s, live, with neither node's `StartedAt` moved; a link between cEOS and
  SR Linux 3.2s; a node added 37.0s; a kind change 53.6s; a node removed 1.04s; nothing to
  do 0.53s. *(26.04, 10-08)*
- A changed SR Linux `.cli` is never applied to a running node by a reconcile. *(10-01;
  26.04, 10-08)*
- Disabling one end of a containerlab link takes the link down at both ends: both are
  oper `down`, neither has a neighbour, and the far end's admin-state stays `enable`.
  *(09-19; 26.04, 10-08)*
- Its `.deb` installs `/usr/bin/containerlab` setuid root, creates the group
  `clab_admins` and puts the installing user in it. *(26.04, 10-07)*
- It does not need `clab_admins` among the calling process's groups: a worker running
  under `sg docker`, with `docker` as its group and no `clab_admins` among its groups
  (the user listed in it by `getent`), deployed and destroyed every twin of tier 3.
  *(26.04, 10-08)*
- Its setup script (`https://containerlab.dev/setup`, `install-docker`) pins Docker per
  distribution release in its own text, `29.8.1` for Ubuntu 26.04, and fails, removing
  nothing, when that release is not in Docker's repository. `CLAB_VERSION` names
  containerlab's own release, and the copy at its `v0.79.0` tag predates 26.04. On 26.04 it
  installed `docker-ce` 29.8.1, `containerd.io` 2.3.6 and the Compose plugin 5.6.0.
  *(10-07; 26.04, 10-07)*

**The `ping` host (016)**

- **The image** is `alpine:3.24.2`, pulled by hand: its id and its repository digest are
  both `sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6` on Docker
  29.8.1, its busybox `v1.37.0`. *(10-29)*
- **busybox 1.37.0's `ping` and `ip` are the mechanism's tools**: `ping` takes `-c`, `-i`,
  `-W` (seconds to wait for the first response after the `-c` packets are sent), `-w`,
  `-q`, `-I` and `-t`, and `-i 0.2` is accepted as root with no minimum-interval refusal;
  `ip` has `address`, `route` and `link`; `pgrep`, `pkill` and a fractional `sleep` are
  applets of the same image. A reply line reads `64 bytes from <a>: seq=0 ttl=64 time=…`
  (`seq=`, where iputils writes `icmp_seq=`), and the summary `N packets transmitted, M
  packets received, P% packet loss`. A background `ping` sent `SIGINT` prints the same
  summary and exits: an interrupt ends a flow and reads its counters. *(Docker 29.8.1,
  10-28)*
- **A `linux` node with `exec` lines plans and deploys beside SR Linux.** `clab deploy
  --dry-run --format json` checks a topology of a `linux` node and a `nokia_srlinux` node
  without the `linux` node's image and pulls nothing. In lab `fylgja`, both hosts of a
  probe deploy beside three privileged SR Linux nodes (46.7 s for 5 nodes), a `linux` node
  privileged too, and each runs its two `exec` lines after its links are made: inside
  `clab-fylgja-n1-loop-a`, `eth1` holds `10.1.10.2/24` and the route `10.1.11.2 via
  10.1.10.1 dev eth1` sits beside the connected one, 0.1–0.2 s after the deploy returns.
  A host's container holds about 700 KiB, and its pings to the far host go through the
  device (TTL 63; `traceroute` reads the gateway, then the far host). *(SR Linux 24.7.1,
  alpine 3.24.2, 10-28, 10-29)*
- **A failed `exec` line is logged, not failed**: a route through an unreachable gateway
  is logged `ERRO Failed to execute command … rc=2` with its stderr (`ip: RTNETLINK
  answers: Network unreachable`), and the deploy exits 0 with the node `running`, so the
  host's readiness must read the lines' effect itself. *(10-29)*
- **`docker exec -d` and `pkill -INT` in a containerlab-run container**: `docker exec -d
  <host> sh -c 'pkill -x ping; exec ping -i 0.2 <far> > /run/flow.log 2>&1'` exits 0 and
  leaves `ping` running; `pkill -INT -x ping` makes it print its counters at once (the
  read 0.15 s), as outside containerlab. *(Docker 29.8.1, 10-29)*

**The step's flows and a burst (016)**

- A burst, `ping -c 5 -i 0.2 -W 2`, takes 0.93–0.96 s when every echo is answered; `twin
  verify` on the probed three-node twin took 5.3 s (22 assertions over 5 nodes), 6.0 s
  with the burst failing. *(10-29)*
- Starting and reading the flows adds no measurable time to a step: three steps of a
  probed twin, each pushing `n1` alone, took 12.0, 40.8 and 18.3 s at the client, the
  second waiting out a 30 s wait on a failing burst. The loss read across them: `0/45`
  with `n1:ethernet-1/2` disabled, `176/190` (longest gap 35.2 s) with the probe's own
  port disabled, and `41/77` (8.2 s) back. Under tier 3, a step across a cEOS node's
  restart in place lost `369/403` (gap 73.8 s) and `321/356` (64.2 s). *(SR Linux 24.7.1,
  cEOS 4.32.0.2F, 10-29; `specs/016-probes/tasks.md` T043, T045)*

---

## SR Linux 24.7.1

**Booting** (one node, Infrahub idle)

- From deploy start the container runs at 0.6s and `sr_cli` answers at 2s, which is
  **not** readiness; `clab deploy` returns at 12.4s, and gNMI `:57400` opens and completes
  a TLS handshake at 15.2–15.4s. *(09-15)*
- One booted node holds about 1.57 GiB (cgroup `memory.current`, steady over a minute).
  *(09-15)*
- The chassis is the `7220 IXR-D2L`, containerlab's default `type` for `nokia_srlinux`.
  Its front ports are `ethernet-1/1` to `ethernet-1/58`, beside `lo0` and `mgmt0`: a push
  naming `ethernet-1/99` is refused `invalid port 99, max is 58`, and one naming
  `ethernet-2/1` `invalid slot 2`. *(09-19)* Seen again as a rehearsal's step pushing an
  artifact that names `ethernet-1/99`: `push.refused` on that node alone. *(26.04, 10-19,
  10-20)*
- The node's name for a port is `ethernet-1/1`, not containerlab's endpoint name `e1-1`:
  `set / interface e1-1 admin-state enable` is refused. *(09-15)*
- Every LLDP adjacency forms before `clab deploy` returns, so a read after the twin is
  `ready` never waits for a neighbour. *(09-19)*
- SR Linux writes to containerlab's node directory on its own
  (`clab-fylgja/<node>/config/gnoi/healthz/events/events.data` moves on every node), so a
  digest of the whole twin directory moves with nothing read. *(10-03)*
- It boots on Ubuntu only once AppArmor allows its `rsyslogd` ([the host](#the-host)).
  *(09-30; 26.04, 10-07, 10-08)*

**gNMI**

- gNMI listens in the management namespace: reach it from the host at the node's
  management address, or by name through containerlab's `/etc/hosts` entry, not from
  inside the container. *(09-15)*
- The encoding must be set: gnmic's default JSON gets `Unimplemented` ("Only ASCII,
  PROTO/TYPED_VAL, JSON_IETF, … supported"), so a probe on the default never succeeds;
  `json_ietf` works. *(09-15)*
- The paths: `/system/name/host-name`; `/system/information/version`
  (`v24.7.1-330-g38f237abfe`); `/interface[name=X]/admin-state`;
  `/system/lldp/interface[name=X]/admin-state`; and `…/neighbor`, whose `system-name` is
  the far node and whose `port-id` is the far port under the far node's own name
  (`port-id-type INTERFACE_NAME`). *(09-19)*
- A leaf comes back as one update at the leaf; a subtree comes back as one update at the
  list entry, with the subtree inside its value, the update path's elements carrying
  module prefixes and the keys inside the value carrying none. *(09-19)*
- An absent value is no update and no error: an uncabled port's neighbour, or a port the
  node lacks, gives `updates: []` with exit 0. *(09-19)*

**The push (JSON-RPC)**

- The push is JSON-RPC `cli` over HTTPS 443 with basic auth: `enter candidate private`,
  the lines, `commit now`. Without the candidate, `set` is refused. It takes 0.62–0.96s
  for an 839-byte artifact, and is idempotent. *(09-18)*
- A refused line comes back as HTTP 200 with a JSON-RPC error that echoes every command
  sent, the artifact included, so keep only the reason. *(09-18)*
- A refused line stays in the login's private candidate across requests, and every later
  `commit now` fails on it, until `discard now` or `load startup` clears it. *(09-19,
  10-01)*
- Under a replace, a refusal's message is `At line N: Error: <reason>`, a blank line, then
  the request's output so far, `diff flat`'s `insert /` lines included; the reason ends at
  the first blank line. *(10-01)*
- SR Linux cuts that message at about 1 KiB, inside its echo of the commands, so a real
  replace's reason never arrives. `commit validate` after it recovers a schema
  constraint's reason, but answers `All changes are valid.` for the port-range check,
  which only a commit makes. *(10-02)*
- `delete /` as a replace darkens the management plane under the shipped template;
  `load startup` resets to the baseline containerlab left and does not. *(09-18, 10-01)*

---

## cEOS 4.32.0.2F

- The image is account-gated: downloaded from Arista with an account and imported with
  `docker import`, which leaves one layer and no `Cmd`, and gives each import a new image
  id. containerlab's `ceos` kind boots it with nothing else stated and supplies the
  command. *(09-20; 26.04, 10-07, 10-08)*
- The file Arista's portal gives, `cEOS-lab-4.32.0.2F.tar`, is xz-compressed despite its
  name, and `docker import` reads the xz stream as it is. The imported image's one layer
  has the decompressed tar's sha256 as its digest, so a present image can be told from the tar it
  came from without importing again: the tar whose sha256 is `89a567d5…` gives the layer
  `sha256:09ab9635…`. *(10-07; 26.04, 10-07, 10-08)*
- It carries production's names one to one: `Ethernet1` … `Ethernet511`
  (`/interfaces/interface/name`, beside `Loopback0` and `Management0`), and modular and
  breakout forms such as `Ethernet2/1` and `Ethernet3/1/1` are the node's own names. A
  push naming a port beyond 511 is refused at parse; one inside the range that the node
  lacks is accepted with a warning. *(09-20)*
- gNMI is plaintext gRPC on 6030, not TLS on 57400 (containerlab's default is `transport
  grpc default` with no SSL profile), and it answers about a second before eAPI does.
  *(09-20, 09-21)*
- The startup-config is the node's whole configuration, not a part merged over a default.
  The kind's default already sets the host name from the node name and brings every port
  up. *(09-20)*
- The push is eAPI `runCmds` over HTTPS to `/command-api` with basic auth, in a named
  configuration session, since a committed session name cannot be reused; a refused
  session is aborted before the next attempt. *(09-20)*
- `#` is refused at token 0, and `!` is read as a comment. *(09-20)*
- `show session-config diffs` has no JSON model: `format: "json"` fails the whole request
  with code 1003 and leaves the session pending, and `format: "text"` answers. The diff
  carries its own bare `!` separators, never the artifact's `!` comment. *(10-01, 10-02)*
- The paths: `/system/state/hostname`; `/system/state/software-version`
  (`4.32.0.2F-41889544.43202F (engineering build)`);
  `/interfaces/interface[name=X]/state/enabled`; `/lldp/interfaces/interface[name=X]/…`.
  *(09-20)*
- An absent value is not a default: OpenConfig reports nothing where SR Linux reports a
  value. *(09-20)*
- A node restarted in place by a reconcile answers gNMI about 49s and eAPI about 55s
  later, on its startup configuration: its push is lost. *(10-01)*
- A gNMI Set is refused (`Unavailable desc = system not yet initialized`) more than a
  minute after the node is ready, while its Gets answer; shut a port through the node's
  own CLI instead: `docker exec clab-fylgja-<node> Cli -p 15 -c 'configure / interface
  Ethernet1 / shutdown'`. *(10-03)*
- A shut port lists a far end as its LLDP neighbour after that end restarted in place,
  while the far end sees none. *(10-03)*

---

## Temporal Server 1.32 / Go SDK 1.49

- The SDK abandons an activity attempt at the first heartbeat call that fails, not after
  `HeartbeatTimeout`, cancelling its context with that call's error as the cause.
  *(09-15)*
- A heartbeat call fails once it has gone unanswered for about half of
  `MaxHeartbeatThrottleInterval`, unreachable or slow alike: 1.0s at 2s, 2.5s at 5s.
  *(09-15)*
- `context.Cause(ctx)` is a `*temporal.CanceledError` for a cancellation the workflow
  asked for, and the RPC's service error for a lost heartbeat. *(09-15)*
- Under `WaitForCancellation`, an activity that never sees the cancellation (one that
  does not heartbeat, or finishes before its next heartbeat) returns its result with no
  error, so the workflow must check `ctx.Err()` after it. *(09-15)*
- The SDK's `testsuite` ignores `WaitForCancellation` and settles a cancelled activity at
  once; that ordering shows only against a live service, or by replaying a history
  recorded from one with `worker.WorkflowReplayer`. *(09-15)*
- Without `WorkflowExecutionErrorWhenAlreadyStarted` the SDK swallows the server's
  already-started error, so a second start waits on the first run and reports its result
  as its own; the `USE_EXISTING` conflict policy makes a second start attach to the run in
  flight. *(09-15)*
- A worker that stopped stays listed as polling its queue for about five minutes.
  *(09-15)*
- **A run whose worker is gone does not end by itself**: with no execution timeout it
  stays running, and only a terminate ends it, which runs nothing of the workflow's own
  cleanup. A child started under `ParentClosePolicy: ABANDON` lives on after its parent's
  terminate, and runs to its end once a worker returns. *(Server 1.32, 10-19)*
- The CLI's installer (`https://temporal.download/cli.sh`, `--version 1.9.1`) puts it in
  `~/.temporalio/bin` and adds that to no shell file's `PATH`, so a login shell has no
  `temporal` until something else adds it. *(10-07; 26.04, 10-07)*
- **A dev server restarted after days down purges its expired runs first**, those past
  the namespace's 24 h retention, at about nine a second. Until it ends, a visibility list
  over many runs times out (`context deadline exceeded`, `RST_STREAM … CANCEL`) while
  `temporal workflow count` answers, so `twin show` and `twin verify` warn
  `show.service.unreachable`. Before tier 3 on a restarted dev server, wait for `workflow
  count` to stop falling. *(CLI 1.9.1, 10-16)*
- **A host of this build makes no Schedule.** After 019's tier 3 whole (nine cases, nine
  twins, `E2E-OK`), `temporal schedule list --address localhost:7233` printed its header
  alone (`ScheduleId  Action  Paused  NextRunTime  LastRunTime`), exit 0, as it did after
  the hand step and again before this record was written. An earlier build's create of a
  following twin made the Schedule `fylgja-follow`. Over the hand step, `temporal workflow
  list` listed exactly the operator's four runs: the provision, two steps and the destroy.
  *(Server 1.32.0, CLI 1.9.1, 10-28; `specs/019-following-removed/tasks.md` T049, T051)*

---

## Prefect 3.8.8

The orchestrator's platform, in `orchestrator/.venv` on Python 3.14.4. §1 of
`specs/020-core-orchestrator-split/research.md` verified 3.9.0 on 10-29; the pins below
made it 3.8.8, and each fact says which it holds on.

**The install**

- **Prefect 3.9.0 installs and runs on Python 3.14.4** alone (`pydantic` 2.14.2, a SQLite
  server), but **does not co-install with `infrahub-sdk` 1.23.2**: `ResolutionImpossible`,
  the SDK requiring `whenever<0.10.0,>=0.9.3` and 3.9.0 `whenever<0.11.0,>=0.10.5`. No SDK
  release, pre-releases included, admits `whenever` 0.10; every Prefect from 3.6.26 to 3.8.8
  says `whenever<0.11.0,>=0.7.3`. *(10-29)*
- **Prefect 3.8.8 with the SDK 1.23.2** resolves (`whenever` 0.9.5), but its server does
  not start against the FastAPI pip picks for it: `fastapi<1.0.0,>=0.139.0` resolved to
  0.162.0, and `prefect server start` fails at import (`concurrency_limits_v2.py`:
  `TypeError: Body() got an unexpected keyword argument 'min_items'`); 3.9.0 narrowed it to
  `fastapi<0.151`. With `fastapi==0.150.0` pinned beside 3.8.8 the server starts and `pip
  check` is clean; `make orchestrator-venv` makes that venv in about 60 s. *(10-29)*

**The server**

- `prefect server start --host <address> --port <port> --no-ui` under a `PREFECT_HOME`
  with `PREFECT_SERVER_API_AUTH_STRING=<user>:<password>` has **basic authentication**:
  without the credential `GET /api/version`, `/api/admin/version` and `POST
  /api/flows/filter` are `401` and `/api/health` stays `200`; with `Authorization: Basic`
  (`httpx` `auth=(user, password)`, or the client's `PREFECT_API_AUTH_STRING`) the API
  answers. Its API version is `0.8.4`, as Infrahub's 3.7.5 reports. `/api/health` answered
  `200` 23.6–32 s after the start, on 3.8.8 as on 3.9.0. `prefect config view` reads
  `PREFECT_RUNNER_POLL_FREQUENCY='10' (from defaults)`, and no profile persists an auth
  string (the home holds no `profiles.toml`). *(3.9.0, 10-29; 3.8.8, 10-29, 10-30)*
- **A `503` at start**: the first `create_flow_run` straight after the deployments
  appeared answered `503 Service Unavailable`, the same request seconds later `201`. *(3.9.0,
  10-29)*

**Deployments, runs and their documents over HTTP**

- `serve(flow.to_deployment(name="main"), …)` from one process (the runner) serves each
  flow's deployment: `GET /api/deployments/name/<flow>/main` answers its `id` within 4 s of
  the runner's start, `404` for one not served, `401` without the credential. `POST
  /api/deployments/<id>/create_flow_run {"parameters": {…}}` answers `201` with the run
  (`SCHEDULED`) in 0.05–0.08 s; `GET /api/flow_runs/<id>` gives `state.type`, `state.name`
  and `end_time`, and `404` for an unknown id. Nothing locks two runs of one deployment.
  *(3.9.0, 10-29; 3.8.8, 10-30)*
- **The runner's poll is 10 s**, Prefect's default: `created` to `start_time` 4.2–9.5 s
  over six runs (12.4 s once, a poll missed), the 2.6–3.1 s between `PENDING` and `RUNNING`
  the runner's subprocess starting. *(3.8.8, 10-30)*
- **`POST /api/flow_runs/filter`'s default order is not newest first**: it is
  `FlowRunSort.ID_DESC` (`prefect/server/api/flow_runs.py:587`) over ids that are `uuid4`
  (`orm_models.py:119`), so with two runs of one flow either can come first;
  `sort: START_TIME_DESC` orders them by time. *(3.8.8, 10-30)*
- **A cancellation**: `POST /api/flow_runs/<id>/set_state {"state": {"type":
  "CANCELLING"}}` without `force` answers `ACCEPT` for a running run, `CANCELLED` about
  4 s later (3.7 s; 4.1 s with `force`); `REJECT` for a scheduled one, which the server ends `CANCELLED` at once; and
  `ABORT` (`Run is already in terminal state COMPLETED.`) for an ended one, which stays as
  it was. **With `force`, an ended `COMPLETED` run is moved to `CANCELLING` and left
  there** (`201 ACCEPT`). *(3.9.0, 10-29; 3.8.8, 10-30)*
- **A flow's document crosses as an artifact**: `Artifact(type="fylgja-findings", key=…,
  data={…}).create()` inside the flow, read by `POST /api/artifacts/filter {"artifacts":
  {"flow_run_id": {"any_": [<id>]}, "type": {"any_": ["fylgja-findings"]}}, "sort":
  "CREATED_DESC", "limit": 1}`, its `data` the dict **as a JSON string**; visible within
  0.02 s of the run's final state. A flow's return value is never read, since its
  persistence needs a result store. *(3.9.0, 10-29; 3.8.8, 10-30)*
- **A run's log**: `POST /api/logs/filter {"logs": {"flow_run_id": {"any_": […]}}, "sort":
  "TIMESTAMP_ASC", "limit", "offset"}` answers its records (`name`, `level`, `message`),
  landing about 1 s after the run reads `COMPLETED`; 250 lines came back whole across two
  pages of 200. Prefect's own lines and `get_run_logger()`'s share the name
  `prefect.flow_runs`; `get_run_logger().getChild("fylgja")` reaches the API under
  `prefect.flow_runs.fylgja`, which tells a flow's lines apart. `prefect.flow_runs`
  propagates to the root logger's console handler (3.8.8's `logging.yml`), so a process
  that prints its own lines sets `PREFECT_LOGGING_HANDLERS_CONSOLE_LEVEL=CRITICAL`.
  *(3.8.8, 10-30)*
- A flow called in a process whose `PREFECT_API_URL` names a server is a run recorded
  there **with no deployment**, beside the deployments' runs. *(3.8.8, 10-30)*

**A timeout and a cancellation**

- **A timeout**: `@flow(timeout_seconds=3)` sleeping in 1 s steps got
  `prefect.flow_engine.FlowRunTimeoutError`, a `TimeoutError`, inside the body 3 s after it
  began, between two sleeps; its `finally` ran, the `on_failure` hook ran, and the run ended
  `FAILED` named `TimedOut` (`Flow run exceeded timeout of 3.0 second(s)`). The runner's
  log also shows `Engine execution exited with unexpected exception … CancelledError`,
  noise. *(3.9.0, 10-29; 3.8.8, 10-30)*
- **A cancellation runs the `finally` and the `on_cancellation` hook, inside a fixed 30 s
  grace**: the runner sends SIGTERM, which the flow sees as `TerminationSignal` (a
  `BaseException`), then SIGKILL 30 s on. A `finally` sleeping 60 s began 0.09 s after
  `set_state CANCELLING` and never ended; the run read `CANCELLED` within 31 s. The 30 s is
  `grace_seconds = 30.0` in the source (`runner/_cancellation_manager.py:102`,
  `runner/runner.py:1275`), read from no setting. A rehearsal cancelled 15 s into its
  baseline's provision asked the core's destroy 0.5 s after, the lab was gone 6.3 s after,
  and the run read `CANCELLED` 15.7 s after, its report posted, inside the grace.
  *(3.8.8, 10-30; `specs/020-core-orchestrator-split/research.md` §5.18, §5.28)*

**The ephemeral server and the test harness**

- **With no `PREFECT_API_URL`, a flow runs on an ephemeral server** under `PREFECT_HOME`,
  recorded in that home's database with its artifacts. Its first start in a fresh home
  lapsed the default `PREFECT_SERVER_EPHEMERAL_STARTUP_TIMEOUT_SECONDS` of 20; at 180 it
  started in about 10–31 s. A `fylgja-cli` command took 39.8–40.4 s on a fresh home and
  16.8–17.1 s warm, every command starting the server in its own process. *(3.8.8,
  10-30)*
- **`prefect.testing.utilities.prefect_test_harness`** exists (3.9.0 and 3.8.8); a flow
  called under it with no `PREFECT_API_URL` is a run it records, `COMPLETED`, its artifact
  read back through Prefect's client. Its first start in a fresh home took 26.9–28.3 s in
  five starts, under its 30 s default but close to it. *(3.9.0, 10-29; 3.8.8, 10-30)*

**The flow's overhead, on the development host**

- `create_flow_run` to a flow's first request to the core: 9.3 s warm, 29.9 s when its
  adapter read was Infrahub's first since a restart; for the checks' rehearsals 13.6–15.4 s
  (10-30, by hand) and 13.9–23.0 s (10-30, tier 3's case 9, a poll missed in the second),
  beside `fylgja-twin`'s 4,260 s. *(3.8.8, 10-30; `specs/020-core-orchestrator-split/research.md`
  §5.28, §5.29)*

---

## GitHub Actions, `ubuntu-26.04`

- The hosted runner for a public repository has 4 vCPUs, 15,983 MiB and 3 GiB of swap,
  Ubuntu 26.04.1 and kernel `7.0.0-1012-azure` (image `20260927.149.1`), and holds Infrahub
  1.11.2 brought up by the bring-up script's Infrahub part: Compose's pulls and start in
  56s, `infrahub-server` healthy 86s later, the registration's import 22s, the seed 24s.
  Tier 2 then took 2 min 46s, and the contract job 6 min 55s in all. Its log shows no peak
  memory. *(10-07)*
- On the Free plan, rulesets and branch protection need the repository public: on a
  private one the API answers `Upgrade to GitHub Pro or make this repository public`. No
  setting disables pull requests; interaction limits last at most six months and cover
  issues too. *(10-07)*
- A `pull_request_target` workflow with no checkout commented on and closed a pull request
  10–15s after it opened, and a job under `if: github.event_name == 'push'` was skipped on
  it. *(10-07, 10-08)*

---

## The host

Ubuntu 26.04.1 LTS, kernel `7.0.0-38-generic`, Docker 29.8.1 (Compose 5.6.0, cgroup v2),
as a QEMU/KVM guest with 10 vCPUs, 32 GiB (31,065 MiB) and 8 GiB of swap, set up by
`scripts/bring-up.sh` *(10-08)*. The facts were first verified on the same guest's earlier
install, Ubuntu 24.04.1, kernel 6.8, Docker 27.5.1, and each says whether it was
re-verified here.

- Ubuntu's `/etc/apparmor.d/usr.sbin.rsyslogd` attaches to the `rsyslogd` inside every
  privileged SR Linux container; unwidened, it fails `log_mgr`, and every `nokia_srlinux`
  deploy fails at containerlab's post-deploy commit (`Applications have failed:
  log_mgr`). Add `/opt/srlinux/** mr,` and `/run/srlinux/** rw,` to
  `/etc/apparmor.d/local/usr.sbin.rsyslogd` and reload it with `apparmor_parser -r`.
  *(09-30)* On 26.04 the profile ships and still attaches: inside a running twin SR
  Linux's `rsyslogd` runs `(enforce)`, as the host's own does, and with the three lines the
  bring-up script writes, those two and `/run/syslogd.pid* rw,`, every SR Linux deploy
  passed. Whether the two alone suffice there was not tried. *(26.04, 10-08)*
- With those lines the profile still denies SR Linux's `rsyslogd` the `net_admin`
  capability (`capable`: 10 denials over a tier-3 run), denies no file, and every deploy
  passes. Nobody looked on 24.04. *(26.04, 10-08)*
- A guest that does not nest has no `/dev/kvm`, which only a `vrnetlab_vm` package would
  need. *(09-30; 26.04, 10-08)*
- **The Docker bridge is how Infrahub's task workers reach the host.** `docker0` is
  `172.17.0.1/16`; Infrahub's Compose network `fylgja-infrahub_default` is
  `172.18.0.0/16`, gateway `172.18.0.1`. With no `/etc/docker/daemon.json`,
  `host-gateway` resolves to `172.17.0.1` on Docker 29.8.1. In a task worker as published,
  `host.docker.internal` does not resolve; with the override's `extra_hosts:
  ["host.docker.internal:host-gateway"]` it resolves to `172.17.0.1`, and a server listening
  there answers it, `200` with the API's token and `401` without, Docker's inter-network
  isolation not applying since the address is the host's own. The published task worker's
  environment holds `INFRAHUB_ADDRESS` and `INFRAHUB_API_TOKEN` and nothing of Fylgja's;
  it has `git` 2.47.3, Python 3.14.7 and `httpx` 0.28.1. *(10-17, 10-19)* Since 020 the bridge carries
  the orchestrator's Prefect server, and the core's server listens on loopback alone
  (below). *(10-30)*
- **A changed Compose override reaches a service only when it is recreated**: `docker
  compose … up -d` recreated the two task workers, whose entry and environment had moved,
  and left the other six services and the volumes as they were. *(10-19)*
- **The host's five processes (020)**: Temporal's dev server on `:7233`; the core's worker;
  the core's server, `fylgja serve --listen 127.0.0.1:7650` alone, started with no
  `INFRAHUB_*` name and no `FYLGJA_DEFAULT_SCOPE` in its environment (`ss` shows `7650` on
  loopback alone); the orchestrator's Prefect server on the bridge, `prefect server start
  --host 172.17.0.1 --port 4200 --no-ui` under `PREFECT_HOME=local/orchestrator/prefect`;
  and its runner, `fylgja-orchestrator serve`, serving six deployments 4.1 s after its
  start. The task workers, recreated with the override's `FYLGJA_ORCHESTRATOR_ADDRESS`
  (`host.docker.internal:4200`) and `FYLGJA_ORCHESTRATOR_AUTH` (16 s, the two alone), read
  `GET …/api/deployments/name/twin-rehearse/main` `200` with the auth and `401` without,
  and `POST http://host.docker.internal:7650/v1/twin/show` fails `ConnectError [Errno 111]
  Connection refused`. A core started from a shell that loaded `local/.env` with only
  `INFRAHUB_ADDRESS`, `INFRAHUB_API_TOKEN` and the scope unset still holds Compose's
  `INFRAHUB_INITIAL_*` and `INFRAHUB_SECURITY_SECRET_KEY`, the admin token equal to the API
  token among them, so every `INFRAHUB_*` name is removed. *(10-30;
  `specs/020-core-orchestrator-split/research.md` §5.28)*
- Infrahub idles at about 5 GiB, and its Neo4j, uncapped, once took the headroom a twin
  needs and killed an end-to-end run; capped at heap 1g/2g and page cache 1g, it does not.
  *(09-21)* On 26.04: 5.2 GiB idle across its eight containers, Neo4j 2.12 GiB under the
  cap and the server 1.56 GiB, leaving 24.5 GiB available with the three processes
  running. *(26.04, 10-08)* After 017's contract tier, before tier 3: Infrahub's
  containers held about 14 GiB, its server 7.5 GiB of that, with 12.1 GiB available, and
  tier 3 passed whole. *(10-26; `specs/017-dcim-reader/research.md` §5.15)* With the
  five processes: 15,770 MiB available before 020's live work, 21,196 MiB after a Compose
  `restart` of Infrahub's containers (12 s); 18,033 MiB after tier 2 whole with those
  containers at about 6.5 GiB, 19,750 MiB after a restart, and tier 3 passed whole.
  *(10-30; `specs/020-core-orchestrator-split/research.md` §5.28, §5.29)*
- **020's tiers on this host**: tier 2 whole (`make test-contract`) 2,692.6 s, its Go half
  about 370 s with no `internal/stages` (the first `go test`'s 21 packages, `internal/cli`
  21.2 s the longest, then `internal/checks` alone 338.3 s) and `pytest -m contract`
  2,322.3 s; tier 3's nine cases 2,798.6 s, 562 s more than 016's 2,236.6 s, since every
  `fylgja-cli` command starts Prefect's ephemeral server in its own process (17–40 s a
  start). The tier-3 figure is a copy of `scripts/e2e.sh` whose case-9 precheck compared
  the checks' syntax trees where the script compares bytes, the script itself refused at
  that precheck. *(10-30; `specs/020-core-orchestrator-split/research.md` §5.29)*
- On a host with little memory to spare, a twin booted beside the contract tier starves
  the worker: it misses heartbeats, its deploys are cut short, and Infrahub times out the
  tier. Read `free -m` before booting a twin. *(09-15)* Not re-verified on 26.04, on
  purpose: proving it starves the host again.
- Under memory pressure SR Linux's own `sr_cli` can segfault, or refuse containerlab's
  post-deploy checkpoint, failing the deploy. *(09-16, 09-17)* Seen on 26.04 too, unasked:
  in a tier-3 run at a load of 5–6 with about 10.2 GiB available, n3's `sr_cli` segfaulted
  in containerlab's post-deploy (`python3.11` in the kernel log), containerlab cancelled
  the other two nodes' post-deploy and the deploy failed; the run after a rest passed
  whole. *(26.04, 10-27; `specs/018-snapshots/tasks.md` T054)*
- A group added to the user reaches no process that began before it: a shell, an agent's
  session or the user's `systemd --user` started before the bring-up script's lab part
  lacks `docker` and `clab_admins` while `getent` lists the user in both, until the next
  login. Under `sg`, nesting `sg docker` and `sg clab_admins` keeps only the inner group.
  *(26.04, 10-08)*
- Ubuntu 26.04's `sudo` is sudo-rs 0.2.13. `sudo -n -v` asks for a password unless every
  rule naming the user is `NOPASSWD`, while `sudo -n true` passes on one (`verifypw=all`),
  and `-E` is ignored with a warning. *(26.04, 10-07)*
- Ubuntu 26.04's `golang-go` (`2:1.26~1`) is go1.26.0, the release `go.mod` names, so no
  toolchain is fetched; by 10-27 an update had made it go1.26.8 (`golang-1.26-go
  1.26.8-1~ubuntu0.26.04.2`), and `go.mod`'s `go 1.26.0` still fetches nothing; its `python3` is 3.14.4, inside `infrahub-sdk` 1.23.2's `>=3.10,
  <3.15`. *(26.04, 10-07, 10-08)*
