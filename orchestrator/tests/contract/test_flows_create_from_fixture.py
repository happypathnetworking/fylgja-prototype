"""The create flow from the fixture, as the core's dry run (T028; R-15, FR-028, SC-002).

`twin-create --branch fylgja-fixture --dry-run` runs on Prefect's test harness against the real
Infrahub and a core server of the test's own, started with no INFRAHUB_* and no node login:
the flow reads the fixture and files its snapshot under the orchestrator's root, and the core
compiles it to a bundle and checks the host (snapshot `0940db93…` and bundle `c3872f5c…` on
the development host; each install's ids are its own, so neither is pinned here), which
refuses only the logins the test's core lacks. No twin boots: a dry run starts no run, and the host check could not
pass. The document is the run's artifact, as a check reads it.

    set -a; . local/.env; set +a; orchestrator/.venv/bin/pytest -q orchestrator -m contract \
        -k test_flows_create_from_fixture -s
"""

import json

import pytest

from fylgja_orchestrator.adapters.infrahub import snapshot
from fylgja_orchestrator.documents.findings import REJECTION, STATUS_REJECTED

from .conftest import FIXTURE_BRANCH, LOGIN_NAMES, SCOPE, artifact_of, flow_runs_of, use_core

pytestmark = pytest.mark.contract

LOGIN_UNSET = "host.probe.login_unset"


def test_flows_create_from_fixture(platform, own_core, tmp_path, monkeypatch):
    from fylgja_orchestrator.flows.twin_create import twin_create

    flows = use_core(monkeypatch, own_core, tmp_path / "orchestrator")
    doc, _ = flows.run(twin_create, branch=FIXTURE_BRANCH, scope=SCOPE, dry_run=True)
    text = flows.text()

    # The snapshot the flow filed, under its root, and the one the core was sent.
    filed = sorted(p.name for p in (flows.root / "snapshots").iterdir())
    assert len(filed) == 1, f"the orchestrator's store holds {filed}"
    snapshot_id = filed[0].removesuffix(".json")
    assert snapshot.snapshot_id(snapshot.read(flows.root / "snapshots" / filed[0])) == snapshot_id
    body = json.loads((flows.root / "snapshots" / filed[0]).read_text())
    assert body["snapshot_version"] == "2" and body["probes"] == [], body.keys()
    assert not (flows.root / "twin.json").exists(), "a dry run wrote the orchestrator's record"

    # The core's dry run: the bundle compiled, the host refused for the logins alone.
    assert doc.operation == "twin.create", doc.operation
    bundle_id = doc.blocks.get("bundle_id", "")
    assert bundle_id and doc.blocks.get("snapshot_id") == snapshot_id, doc.blocks
    assert snapshot_id[:8] in text and bundle_id[:8] in text, text
    rejections = [f for f in doc.findings if f.severity == REJECTION]
    assert rejections and {f.rule for f in rejections} == {LOGIN_UNSET}, [f.to_dict() for f in rejections]
    assert sorted({f.object for f in rejections}) == sorted(n for n in LOGIN_NAMES if "SRLINUX" in n), rejections
    assert doc.status == STATUS_REJECTED and doc.exit_code == 1, doc.status

    # Nothing ran: the core was asked one dry run, which starts no run.
    assert own_core.operations() == [("twin.create", "rejected")], own_core.operations()
    assert not (own_core.state / "twin").exists(), "the dry run made a twin directory"

    # The run on the harness, its document filed as the run's artifact.
    [run] = flow_runs_of("twin-create")[:1]
    assert run.state.is_completed(), run.state
    assert artifact_of(str(run.id)) == doc.to_dict()
