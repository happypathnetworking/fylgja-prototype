"""The adapter's read of a real Infrahub, as `internal/intent/read_contract_test.go` held the
Go reader to it (T020; SC-002): the fixture branch reads to a snapshot whose id is the hash
of its own body (`0940db93…` on the development host; each install's envelope is its own,
so the id is never pinned here) and whose CTM is `testdata/ctm/three-node.json` but for the
envelope; the cables the read cannot place
are omitted; the platforms it refuses are named in one pass; a scope selecting nothing is
`scope.empty`; an unpinned read records `observed_at` and no `at`, a pinned one sends and
records its `at` verbatim; Infrahub's refusals of an `at` and of a branch are operational
failures naming what was asked; the schema endpoint names every kind and implementer the
mapping reads; and a wrong credential never reaches an error.

The fixture branch is only read. Each case that changes a branch seeds its own throwaway
`fylgja-test-py-*` through the fixture tool and deletes it when it ends.

    set -a; . local/.env; set +a; orchestrator/.venv/bin/pytest -q orchestrator -m contract -k test_read_fixture
"""

import json
import os
import time
from datetime import UTC, datetime, timedelta
from urllib.parse import parse_qs, urlsplit

import pytest

from fylgja_orchestrator.adapters.infrahub import __main__ as adapter
from fylgja_orchestrator.adapters.infrahub import mapping, packages, schema, snapshot
from fylgja_orchestrator.adapters.infrahub import read as reading
from fylgja_orchestrator.adapters.infrahub.client import Client, Config, ClientError, from_env
from fylgja_orchestrator.documents.findings import INFO, REJECTION, Finding

from .conftest import FIXTURE_BRANCH, REPO, SCOPE, artifacts_ready_at, read_under

pytestmark = pytest.mark.contract

LIBRARY_MAPPING = {"reader": mapping.READER, "model": mapping.MODEL, "version": mapping.VERSION}


def read(branch: str, at: str = "", scope: str = SCOPE) -> dict:
    """read_under, failing on anything the command would not write a CTM for."""
    ctm, found, err = read_under(branch, at, scope)
    assert err is None, f"reading {branch}: {err}"
    assert not found.rejected, f"reading {branch} was rejected: {list(found)}"
    return ctm


def test_the_fixture_reads_to_its_snapshot_and_the_three_node_ctm(tmp_path):
    out = tmp_path / "snap.json"
    answer = adapter.read(FIXTURE_BRANCH, scope=SCOPE, out=str(out), env=adapter.Env())
    assert answer.doc.status == "ok", list(answer.doc.findings)
    assert list(answer.doc.findings) == []
    taken = snapshot.read(out)
    assert taken["snapshot_id"] == snapshot.snapshot_id(taken)
    assert answer.doc.blocks["snapshot_id"] == taken["snapshot_id"]
    assert taken["probes"] == []

    got = taken["ctm"]
    want = snapshot.normalize(json.loads((REPO / "testdata" / "ctm" / "three-node.json").read_text()))
    e, w = got["envelope"], want["envelope"]
    assert e["scope"] == SCOPE == w["scope"]
    assert e["mapping"] == LIBRARY_MAPPING == w["mapping"]
    # The fixture's cables are all placed: the three fixed rows alone.
    assert len(e["fidelity"]) == 3 and e["fidelity"] == w["fidelity"]
    got["envelope"] = w
    assert got == want
    assert snapshot.marshal(snapshot.take(got)) == snapshot.marshal(snapshot.take(want))


def test_the_read_counts_the_fixture():
    # Three devices, each with its management port, two cabled ports and two probe-ready
    # ports, three cables and three artifacts.
    assert reading.summarize(read(FIXTURE_BRANCH)) == reading.Counts(devices=3, interfaces=15, links=3, artifacts=3)


def test_cables_the_read_cannot_place(fixture_tool):
    """US2 scenarios 2 and 3 of 017: a console server cs1 carrying no tag, with a port, and a
    cable from a new n1:ethernet-1/3 to it; a cable joining three new ports. Each is omitted,
    an info finding and a fidelity row, and refuses nothing; cs1 is never read, and the end
    outside the scope is named by its device's id."""
    with fixture_tool.throwaway("omit") as branch:
        added = fixture_tool.change(branch, "-add-device", "cs1:port1")
        cs1 = added.split(": id ", 1)[1].split(",", 1)[0]
        for port in ("n1:ethernet-1/3", "n1:ethernet-1/4", "n2:ethernet-1/3", "n3:ethernet-1/3"):
            fixture_tool.change(branch, "-add-interface", port)
        fixture_tool.change(branch, "-add-cable", "c-out=n1:ethernet-1/3,cs1:port1")
        fixture_tool.change(branch, "-add-cable", "c-three=n1:ethernet-1/4,n2:ethernet-1/3,n3:ethernet-1/3")

        ctm, found, err = read_under(branch)
    assert err is None
    # The artifacts were rendered before the ports were added and are read as they stand: a
    # port the artifact does not name is no finding. So the omissions are the only findings.
    outside = cs1 + ":port1"
    assert list(found) == [
        Finding(INFO, reading.RULE_OMIT_CABLE_OUTSIDE_SCOPE, "c-out",
                f"cable c-out is omitted: an end outside the scope ({outside}); "
                f"its ends are {outside}, n1:ethernet-1/3"),
        Finding(INFO, reading.RULE_OMIT_CABLE_ENDS, "c-three",
                "cable c-three is omitted: other than two ends (3); its ends are n1:ethernet-1/4, n2:ethernet-1/3, "
                "n3:ethernet-1/3"),
    ]
    assert ctm is not None
    fidelity = ctm["envelope"]["fidelity"]
    for row in (
        {"disposition": "omitted", "object": "c-out", "source": "DcimCable",
         "reason": f"an end outside the scope ({outside})"},
        {"disposition": "omitted", "object": "c-three", "source": "DcimCable", "reason": "other than two ends (3)"},
    ):
        assert row in fidelity
    assert len(fidelity) == 5
    assert reading.summarize(ctm) == reading.Counts(devices=3, interfaces=19, links=3, artifacts=3, cables_omitted=2)
    devices = {d["name"]: d for d in ctm["devices"]}
    assert "cs1" not in devices
    for p in ("n1:ethernet-1/3", "n1:ethernet-1/4", "n2:ethernet-1/3", "n3:ethernet-1/3"):
        device, port = p.split(":")
        interfaces = {i["name"]: i for i in devices[device]["interfaces"]}
        assert port in interfaces, f"port {p} was not read"
        assert not interfaces[port].get("link"), f"port {p} read as {interfaces[port]}, cabled to something"


def test_platforms_the_read_refuses(fixture_tool):
    """n1 moved onto nokia / SR-Linux, which differs from Nokia / SR Linux in letter case and
    spelling, and n2's platform taken away: both refused by name in one pass, with no CTM."""
    with fixture_tool.throwaway("platform") as branch:
        fixture_tool.change(branch, "-set-platform", "n1=nokia/SR-Linux")
        fixture_tool.change(branch, "-set-platform", "n2=")
        ctm, found, err = read_under(branch)
    assert err is None
    assert ctm is None, "a refused read returned a CTM"
    want = [
        (REJECTION, reading.RULE_PLATFORM_UNMATCHED, "n1",
         'device n1: manufacturer "nokia", platform "SR-Linux"; the packages match: '
         "arista_eos [Arista / EOS], nokia_srlinux [Nokia / SR Linux]"),
        (REJECTION, reading.RULE_PLATFORM_MISSING, "n2", ""),
    ]
    assert len(found) == len(want), list(found)
    for f, (severity, rule, obj, prefix) in zip(found, want, strict=True):
        assert (f.severity, f.rule, f.object) == (severity, rule, obj)
        assert f.message.startswith(prefix)
    assert "device n2 carries no platform" in found[1].message


def test_a_scope_selecting_nothing_is_refused():
    ctm, found, err = read_under(FIXTURE_BRANCH, scope="tag:nothing")
    assert err is None
    assert ctm is None
    assert list(found) == [Finding(REJECTION, reading.RULE_SCOPE_EMPTY, "tag:nothing",
                                   f"scope tag:nothing selects no device on branch {FIXTURE_BRANCH}")]


def test_an_unpinned_read_records_observed_at_and_no_at():
    before = datetime.now(UTC) - timedelta(seconds=1)
    e = read(FIXTURE_BRANCH)["envelope"]
    assert e["branch"] == FIXTURE_BRANCH
    assert "at" not in e, "no at was supplied, so none may be recorded"
    # Six fractional digits, always, so an observed_at can be handed back as --at.
    observed = datetime.strptime(e["observed_at"], reading.OBSERVED_AT_FORMAT).replace(tzinfo=UTC)
    assert snapshot.OBSERVED_AT_FORM.match(e["observed_at"])
    assert observed >= before
    assert reading.check_at_precision(e["observed_at"]) == ""
    conformance, _ = schema.check_conformance(Client(from_env(FIXTURE_BRANCH)))
    assert e["schema_hash"] and e["schema_hash"] == conformance.schema_hash
    assert e["mapping"] == LIBRARY_MAPPING


def test_a_pinned_read_sends_and_records_its_at_verbatim(proxy):
    # The whole second after the fixture's last artifact became Ready, so the read is
    # complete, with six fractional digits that nothing may normalise; taken from the running
    # Infrahub, since a constant names one install's fixture.
    pinned = artifacts_ready_at(FIXTURE_BRANCH).astimezone(UTC).replace(microsecond=0) + timedelta(seconds=1)
    at = pinned.strftime("%Y-%m-%dT%H:%M:%S") + ".123456Z"
    wait = (pinned + timedelta(microseconds=123456) - datetime.now(UTC)).total_seconds()
    if wait > 0:
        time.sleep(wait)

    token = os.environ["INFRAHUB_API_TOKEN"]
    client = Client(Config(address=proxy.address, token=token, branch=FIXTURE_BRANCH, at=at))
    ctm, found = reading.read(client, SCOPE, packages.load(), reading.observed_now())
    assert ctm is not None and not found.rejected, list(found)
    graphql = [path for method, path, _ in proxy.recorded if method == "POST" and path.startswith("/graphql/")]
    assert graphql, "the read asked GraphQL nothing"
    for path in graphql:
        assert parse_qs(urlsplit(path).query) == {"at": [at]}, f"the request {path} does not carry the at verbatim"
    assert ctm["envelope"]["at"] == at
    assert ctm["envelope"]["observed_at"]
    assert len(ctm["devices"]) == 3


def test_an_at_before_the_default_branch_is_an_operational_failure():
    at = "2000-01-01T00:00:00Z"
    ctm, _, err = read_under(FIXTURE_BRANCH, at)
    assert err is not None, "an instant before the default branch was not refused"
    assert ctm is None
    for want in (FIXTURE_BRANCH, at, "before branch"):
        assert want in str(err)


def test_an_unparseable_at_is_an_operational_failure():
    ctm, _, err = read_under(FIXTURE_BRANCH, "yesterday")
    assert err is not None and ctm is None
    assert "Invalid time format" in str(err)


def test_an_unknown_branch_is_reported():
    branch = "fylgja-no-such-branch"
    ctm, _, err = read_under(branch)
    assert err is not None and ctm is None
    # Both halves: the branch named, and that it does not exist, in the schema endpoint's
    # words, the first thing a read touches.
    assert branch in str(err)
    assert "not found" in str(err)


def test_the_schema_endpoint_names_the_kinds_and_implementers():
    info = schema.schema_info(Client(from_env(FIXTURE_BRANCH)))
    for table in mapping.TABLES:
        kinds = info.generics if table.generic else info.nodes
        assert table.kind in kinds, f"kind {table.kind} is absent"
        for r in table.relationships:
            assert r.name in kinds[table.kind].relationships, f"{table.kind} has no relationship {r.name}"
    for u in mapping.USED_BY:
        assert u.implementer in info.implementers(u.generic), f"{u.generic} is not implemented by {u.implementer}"
    assert info.hash


def test_a_wrong_credential_does_not_leak():
    cfg = from_env(FIXTURE_BRANCH)
    bogus = "deliberately-wrong-token-9f3a"
    with pytest.raises(ClientError) as raised:
        schema.schema_info(Client(Config(address=cfg.address, token=bogus, branch=FIXTURE_BRANCH)))
    assert bogus not in str(raised.value)
    assert "INFRAHUB_API_TOKEN" in str(raised.value)
