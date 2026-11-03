//go:build contract || fixture

// Package testsupport is the contract-test harness.
//
// It is the one place Fylgja code writes to Infrahub — creating branches, loading
// schema, seeding fixtures. The product must never do that (Constitution III),
// so this package is build-tagged and is never imported by cmd/fylgja; a test in
// internal/compiler asserts as much.
package testsupport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/happypathnetworking/fylgja/internal/probe"
)

// Client talks to Infrahub with write access.
type Client struct {
	Address string
	Token   string
	HTTP    *http.Client
}

// NewClient builds a client from INFRAHUB_ADDRESS and INFRAHUB_API_TOKEN.
func NewClient() (*Client, error) {
	addr := os.Getenv("INFRAHUB_ADDRESS")
	token := os.Getenv("INFRAHUB_API_TOKEN")
	if addr == "" || token == "" {
		return nil, fmt.Errorf("INFRAHUB_ADDRESS and INFRAHUB_API_TOKEN must be set")
	}
	return &Client{
		Address: strings.TrimRight(addr, "/"),
		Token:   token,
		HTTP:    &http.Client{Timeout: 120 * time.Second},
	}, nil
}

func (c *Client) do(method, path string, body any) (int, []byte, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.Address+path, rdr)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-INFRAHUB-KEY", c.Token)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, nil, c.redact(err)
	}
	defer func() { _ = resp.Body.Close() }()
	out, err := io.ReadAll(resp.Body)
	return resp.StatusCode, out, err
}

// redact strips the token from an error, so a harness failure cannot leak it either.
func (c *Client) redact(err error) error {
	if err == nil || c.Token == "" {
		return err
	}
	return fmt.Errorf("%s", strings.ReplaceAll(err.Error(), c.Token, "<redacted>"))
}

type gqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// GraphQL runs a query or mutation against a branch.
func (c *Client) GraphQL(branch, query string) (json.RawMessage, error) {
	status, body, err := c.do(http.MethodPost, "/graphql/"+branch, map[string]any{"query": query})
	if err != nil {
		return nil, err
	}
	var r gqlResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("graphql %d: %s", status, string(body))
	}
	if len(r.Errors) > 0 {
		return r.Data, fmt.Errorf("graphql: %s", r.Errors[0].Message)
	}
	return r.Data, nil
}

// CreateBranch creates a branch, treating "already exists" as success so seeding is
// repeatable.
func (c *Client) CreateBranch(name string) error {
	q := fmt.Sprintf(`mutation { BranchCreate(data: {name: %q, sync_with_git: false}) { ok } }`, name)
	if _, err := c.GraphQL("main", q); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	return nil
}

// DeleteBranch removes a branch.
func (c *Client) DeleteBranch(name string) error {
	q := fmt.Sprintf(`mutation { BranchDelete(data: {name: %q}) { ok } }`, name)
	_, err := c.GraphQL("main", q)
	return err
}

// BranchExists reports whether Infrahub lists a branch by name. It writes nothing.
func (c *Client) BranchExists(name string) (bool, error) {
	raw, err := c.GraphQL("main", fmt.Sprintf(`{ Branch(name: %q) { name } }`, name))
	if err != nil {
		return false, fmt.Errorf("reading branch %s: %w", name, err)
	}
	var r struct {
		Branch []struct {
			Name string `json:"name"`
		} `json:"Branch"`
	}
	if err := json.Unmarshal(raw, &r); err != nil {
		return false, err
	}
	for _, b := range r.Branch {
		if b.Name == name {
			return true, nil
		}
	}
	return false, nil
}

// LibraryDir is where schemaDir keeps the schema library's six files, vendored byte for
// byte at the commit the reader reads (schema/dcim/README.md).
const LibraryDir = "dcim"

// LoadSchema loads schemaDir onto a branch, then waits for the branch's GraphQL schema to
// be rebuilt. The rebuild is asynchronous: seeding immediately after a load races it.
//
// schemaDir is the repository's schema/: the schema library's files under dcim/ go in one
// request, since they are one model whose kinds name each other (R-12), then the files
// beside them, the stage kind, in a second, since the stage kind needs none of the library
// and its extension of CoreProposedChange is already on main. A directory with no dcim/
// loads its own files alone, as a test's additive load does. The wait is for the Create
// mutation of every node the loaded files define and do not remove (state: absent), so
// for DcimDeviceCreate and FylgjaStageCreate among them.
func (c *Client) LoadSchema(branch, schemaDir string) error {
	library, err := libraryFiles(schemaDir)
	if err != nil {
		return err
	}
	own, err := ownFiles(schemaDir)
	if err != nil {
		return err
	}
	if len(library)+len(own) == 0 {
		return fmt.Errorf("%s holds no schema file: neither %s/*.yml nor *.yaml", schemaDir, LibraryDir)
	}
	var mutations []string
	for _, paths := range [][]string{library, own} {
		if len(paths) == 0 {
			continue
		}
		created, _, err := c.loadSchemaFiles(branch, paths)
		if err != nil {
			return err
		}
		mutations = append(mutations, created...)
	}
	if len(mutations) == 0 {
		return nil
	}
	return c.waitForKinds(branch, mutations...)
}

func libraryFiles(schemaDir string) ([]string, error) {
	return filepath.Glob(filepath.Join(schemaDir, LibraryDir, "*.yml"))
}

func ownFiles(schemaDir string) ([]string, error) {
	return filepath.Glob(filepath.Join(schemaDir, "*.yaml"))
}

// definedKinds is what one schema file defines: its nodes and its generics, each with the
// state that removes it.
type definedKinds struct {
	Nodes    []definedKind `yaml:"nodes"`
	Generics []definedKind `yaml:"generics"`
}

type definedKind struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	State     string `yaml:"state"`
}

func (k definedKind) kind() string { return k.Namespace + k.Name }

// LibraryKinds names every node and generic the schema library's files under
// schemaDir/dcim define, sorted: the kinds -prepare-main's gate looks for (R-13). They are
// read from the files the gate guards, so the list cannot drift from what is loaded; with
// the repository's schema/ those are the embedded files, byte for byte
// (TestVendoredLibraryIsThePinnedCommit).
func LibraryKinds(schemaDir string) ([]string, error) {
	paths, err := libraryFiles(schemaDir)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("%s holds no file of the schema library", filepath.Join(schemaDir, LibraryDir))
	}
	seen := map[string]bool{}
	for _, p := range paths {
		raw, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var defined definedKinds
		if err := yaml.Unmarshal(raw, &defined); err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		for _, k := range append(defined.Nodes, defined.Generics...) {
			if k.State != "absent" {
				seen[k.kind()] = true
			}
		}
	}
	kinds := make([]string, 0, len(seen))
	for k := range seen {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	return kinds, nil
}

// KindsHeld lists which of kinds the branch's schema holds, as a node or a generic, sorted.
// It reads GET /api/schema?branch= and writes nothing.
func (c *Client) KindsHeld(branch string, kinds []string) ([]string, error) {
	status, body, err := c.do(http.MethodGet, "/api/schema?branch="+url.QueryEscape(branch), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("reading the schema of %s: HTTP %d: %s", branch, status, c.redact(fmt.Errorf("%s", body)))
	}
	var doc struct {
		Nodes []struct {
			Kind string `json:"kind"`
		} `json:"nodes"`
		Generics []struct {
			Kind string `json:"kind"`
		} `json:"generics"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, fmt.Errorf("reading the schema of %s: %w", branch, err)
	}
	has := map[string]bool{}
	for _, n := range doc.Nodes {
		has[n.Kind] = true
	}
	for _, g := range doc.Generics {
		has[g.Kind] = true
	}
	var held []string
	for _, k := range kinds {
		if has[k] {
			held = append(held, k)
		}
	}
	sort.Strings(held)
	return held, nil
}

// LibraryHeldError is -prepare-main's refusal (R-13): the branch already holds some of the
// schema library's kinds, so nothing of the library was loaded. Its sentence's opening, up
// to "an empty main alone", is the one the bring-up script reads as present; the rest says
// the tool is not an install path (US6 scenario 2).
type LibraryHeldError struct {
	Branch string
	Held   []string
}

func (e *LibraryHeldError) Error() string {
	return fmt.Sprintf("%s holds %d of the schema library's kinds (%s): -prepare-main is the development host's and CI's "+
		"tool, which loads the library on an empty main alone; it is not an install path: a main that already holds "+
		"the library is yours, and fylgja schema check says whether it is the one the reader reads",
		e.Branch, len(e.Held), strings.Join(e.Held, ", "))
}

// LoadLibrary is -prepare-main's gate and load (R-13). Before any write it reads which of
// the library's kinds the branch holds: any, and it loads nothing and returns a
// *LibraryHeldError naming them; none, and it loads schemaDir/dcim's files in one
// request, waits for their Create mutations and returns how many files it loaded.
func (c *Client) LoadLibrary(branch, schemaDir string) (files int, err error) {
	kinds, err := LibraryKinds(schemaDir)
	if err != nil {
		return 0, err
	}
	held, err := c.KindsHeld(branch, kinds)
	if err != nil {
		return 0, err
	}
	if len(held) > 0 {
		return 0, &LibraryHeldError{Branch: branch, Held: held}
	}
	paths, err := libraryFiles(schemaDir)
	if err != nil {
		return 0, err
	}
	mutations, _, err := c.loadSchemaFiles(branch, paths)
	if err != nil {
		return 0, err
	}
	return len(paths), c.waitForKinds(branch, mutations...)
}

// PrepareMainSchema is -prepare-main's schema: Fylgja's own kinds first, then the library,
// on the default branch (D-028, D-064). The own kinds load in one request, a no-op where
// both are present, so a main holding the library and lacking one gains it; then
// LoadLibrary, which refuses a main holding any of the library's kinds with a
// *LibraryHeldError and loads nothing of it. own is whether the own kinds' load moved the
// hash, set whenever that load succeeded, the library's refusal included; files is the
// library's files loaded, 0 on its refusal.
func (c *Client) PrepareMainSchema(schemaDir string) (own bool, files int, err error) {
	own, err = c.LoadOwnSchema(DefaultBranch, schemaDir)
	if err != nil {
		return false, 0, err
	}
	files, err = c.LoadLibrary(DefaultBranch, schemaDir)
	return own, files, err
}

// LoadOwnSchema loads schemaDir's own files, the stage kind, in one request and waits for
// their Create mutations. loaded is false when the branch's schema hash did not move: the
// kind was there already, as the files define it.
func (c *Client) LoadOwnSchema(branch, schemaDir string) (loaded bool, err error) {
	paths, err := ownFiles(schemaDir)
	if err != nil {
		return false, err
	}
	if len(paths) == 0 {
		return false, fmt.Errorf("%s holds no *.yaml schema file", schemaDir)
	}
	mutations, loaded, err := c.loadSchemaFiles(branch, paths)
	if err != nil {
		return false, err
	}
	return loaded, c.waitForKinds(branch, mutations...)
}

// LoadSchemaFile loads one schema file onto a branch, as a test changes a branch's schema,
// then waits until the branch's GraphQL schema exposes the Create mutation of every node the
// file defines and no longer exposes that of every node it removes (state: absent): a
// removal is rebuilt asynchronously too, and a read made before the rebuild still sees the
// kind. It never runs on main.
func (c *Client) LoadSchemaFile(branch, path string) error {
	if branch == DefaultBranch {
		return fmt.Errorf("LoadSchemaFile refuses %s: a test changes a throwaway branch's schema, never the default branch's", branch)
	}
	created, _, err := c.loadSchemaFiles(branch, []string{path})
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var defined definedKinds
	if err := yaml.Unmarshal(raw, &defined); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	var removed []string
	for _, n := range defined.Nodes {
		if n.State == "absent" {
			removed = append(removed, n.kind()+"Create")
		}
	}
	if err := c.waitForKinds(branch, created...); err != nil {
		return err
	}
	return c.waitForKindsGone(branch, removed...)
}

// loadSchemaFiles loads files onto a branch in one request and returns the Create
// mutation of each node they define and do not remove, and whether the branch's schema
// hash moved, as infrahubctl tells an updated schema from one already in place.
func (c *Client) loadSchemaFiles(branch string, paths []string) (mutations []string, updated bool, err error) {
	sort.Strings(paths)
	schemas := make([]any, 0, len(paths))
	for _, p := range paths {
		raw, err := os.ReadFile(p)
		if err != nil {
			return nil, false, err
		}
		var doc any
		if err := yaml.Unmarshal(raw, &doc); err != nil {
			return nil, false, fmt.Errorf("%s: %w", p, err)
		}
		var defined definedKinds
		if err := yaml.Unmarshal(raw, &defined); err != nil {
			return nil, false, fmt.Errorf("%s: %w", p, err)
		}
		for _, n := range defined.Nodes {
			if n.State != "absent" {
				mutations = append(mutations, n.kind()+"Create")
			}
		}
		schemas = append(schemas, doc)
	}
	status, body, err := c.do(http.MethodPost, "/api/schema/load?branch="+branch, map[string]any{"schemas": schemas})
	if err != nil {
		return nil, false, err
	}
	if status != http.StatusOK && status != http.StatusAccepted {
		return nil, false, fmt.Errorf("schema load of %d file(s) on %s, %d: %s", len(paths), branch, status, string(body))
	}
	var hashes struct {
		Hash         string `json:"hash"`
		PreviousHash string `json:"previous_hash"`
	}
	if err := json.Unmarshal(body, &hashes); err != nil {
		return nil, false, fmt.Errorf("schema load of %d file(s) on %s: %w", len(paths), branch, err)
	}
	updated = hashes.Hash != "" && hashes.PreviousHash != "" && hashes.Hash != hashes.PreviousHash
	return mutations, updated, nil
}

// waitForKinds polls until the named mutations appear on the branch's GraphQL schema.
func (c *Client) waitForKinds(branch string, mutations ...string) error {
	const q = `{ __type(name: "Mutation") { fields { name } } }`
	deadline := time.Now().Add(90 * time.Second)
	missing := mutations
	for time.Now().Before(deadline) {
		data, err := c.GraphQL(branch, q)
		if err == nil {
			var r struct {
				Type struct {
					Fields []struct {
						Name string `json:"name"`
					} `json:"fields"`
				} `json:"__type"`
			}
			if json.Unmarshal(data, &r) == nil {
				have := map[string]bool{}
				for _, f := range r.Type.Fields {
					have[f.Name] = true
				}
				missing = nil
				for _, m := range mutations {
					if !have[m] {
						missing = append(missing, m)
					}
				}
				if len(missing) == 0 {
					return nil
				}
			}
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("branch %s: GraphQL schema did not expose %v within 90s", branch, missing)
}

// waitForKindsGone polls until none of the named mutations is on the branch's GraphQL
// schema.
func (c *Client) waitForKindsGone(branch string, mutations ...string) error {
	if len(mutations) == 0 {
		return nil
	}
	const q = `{ __type(name: "Mutation") { fields { name } } }`
	deadline := time.Now().Add(90 * time.Second)
	still := mutations
	for time.Now().Before(deadline) {
		data, err := c.GraphQL(branch, q)
		if err == nil {
			var r struct {
				Type struct {
					Fields []struct {
						Name string `json:"name"`
					} `json:"fields"`
				} `json:"__type"`
			}
			if json.Unmarshal(data, &r) == nil {
				have := map[string]bool{}
				for _, f := range r.Type.Fields {
					have[f.Name] = true
				}
				still = nil
				for _, m := range mutations {
					if have[m] {
						still = append(still, m)
					}
				}
				if len(still) == 0 {
					return nil
				}
			}
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("branch %s: GraphQL schema still exposed %v after 90s", branch, still)
}

// create runs a Create mutation and returns the new object's id.
func (c *Client) create(branch, kind, data string) (string, error) {
	q := fmt.Sprintf(`mutation { %sCreate(data: %s) { ok object { id } } }`, kind, data)
	raw, err := c.GraphQL(branch, q)
	if err != nil {
		return "", fmt.Errorf("creating %s: %w", kind, err)
	}
	var r map[string]struct {
		Object struct {
			ID string `json:"id"`
		} `json:"object"`
	}
	if err := json.Unmarshal(raw, &r); err != nil {
		return "", err
	}
	return r[kind+"Create"].Object.ID, nil
}

// ProbeKind is Fylgja's probe kind (016, D-064): branch-agnostic, on the default branch,
// where the server reads it when a request sends no probe file.
const ProbeKind = "FylgjaProbe"

// WriteProbes writes each probe of a declaration as one object of the probe kind on the
// default branch, its endpoints the a_ and b_ attributes, every value verbatim (R-26). The
// kind is branch-agnostic, so a probe lands in no branch's diff and every branch's request
// reads it: whoever writes one deletes it (DeleteProbes) before another test reads the kind.
func (c *Client) WriteProbes(decl *probe.Declaration) error {
	if decl == nil {
		return fmt.Errorf("writing probes: no declaration")
	}
	for _, p := range decl.Probes {
		a, b := p.Endpoints[0], p.Endpoints[1]
		data := fmt.Sprintf(`{name: {value: %q}, format: {value: %q}, `+
			`a_device: {value: %q}, a_interface: {value: %q}, a_address: {value: %q}, a_gateway: {value: %q}, `+
			`b_device: {value: %q}, b_interface: {value: %q}, b_address: {value: %q}, b_gateway: {value: %q}}`,
			p.Name, p.Format, a.Device, a.Interface, a.Address, a.Gateway, b.Device, b.Interface, b.Address, b.Gateway)
		id, err := c.create(DefaultBranch, ProbeKind, data)
		if err != nil {
			return fmt.Errorf("writing probe %s: %w", p.Name, err)
		}
		if id == "" {
			return fmt.Errorf("writing probe %s: %sCreate returned no object", p.Name, ProbeKind)
		}
	}
	return nil
}

// DeleteProbes deletes every object of the probe kind on the default branch and returns
// how many, so a test or the fixture tool leaves the kind as it found it: empty.
func (c *Client) DeleteProbes() (int, error) {
	ids, err := c.ids(DefaultBranch, ProbeKind, "")
	if err != nil {
		return 0, err
	}
	n := 0
	for _, id := range ids {
		if err := c.deleteNode(DefaultBranch, ProbeKind, id); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// json_Unmarshal is a tiny indirection so seed.go can decode without importing
// encoding/json again under a different name.
func json_Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

// ArtifactDefinitionName is the definition this repository's .infrahub.yml declares, which
// renders a device's configuration once the repository is registered on main (D-044). It
// is the name, not an id: the id differs per install, and a branch inherits main's
// definition (D-028).
const ArtifactDefinitionName = "srlinux_device_config"

// ArtifactGroupName is the standard group whose members the definition renders for. A
// device outside it has no artifact, which is a refusal the tiers prove.
const ArtifactGroupName = "fylgja-devices"

// ArtifactWait bounds every wait for generation, WaitForArtifacts' and each contract
// test's wait for a checksum to move: 180s, about twenty-four times the 7.6s measured for
// three devices alone. 90s, twelve times, lapsed three times under the whole tier on
// 2026-10-29 (016's T044: one artifact still Pending, each test green alone) and once at
// 019's close, as Infrahub's task workers queue the tier's generations behind one another;
// the schema waits keep their 90s, which never lapsed.
const ArtifactWait = 180 * time.Second

// unregistered is what a lookup by name says when main lacks what the registration brings,
// or a branch predates it: the definition comes with the repository's import, and the
// group with fylgja-fixture -prepare-main, which makes the registration too.
const unregistered = "register this repository (github.com/happypathnetworking/fylgja-prototype) on main as a " +
	"read-only repository with no credential and wait for its import " +
	`(docs/development.md, "What Infrahub needs"); a branch created before the import cannot see it`

// lookupByName returns the one object of a kind with the given name, or an error naming
// how many were found. Used for the definition and the group, which are looked up by
// name because their ids differ per install.
func (c *Client) lookupByName(branch, kind, name string) (string, error) {
	found, err := c.ids(branch, kind, fmt.Sprintf(`(name__value: %q)`, name))
	if err != nil {
		return "", err
	}
	if len(found) != 1 {
		return "", fmt.Errorf("branch %s holds %d %s named %q, want one; "+unregistered, branch, len(found), kind, name)
	}
	return found[0], nil
}

// ArtifactGroup returns the id of the group the artifact definition renders for, on a
// branch. A branch created before the group existed on main cannot see it, which is a
// real failure rather than an empty result: the seed would silently render nothing.
func (c *Client) ArtifactGroup(branch string) (string, error) {
	return c.lookupByName(branch, "CoreStandardGroup", ArtifactGroupName)
}

// JoinArtifactGroup adds devices to the group the definition renders for. Membership is
// what makes a device a target; it does not by itself render anything.
func (c *Client) JoinArtifactGroup(branch string, deviceIDs []string) error {
	if len(deviceIDs) == 0 {
		return nil
	}
	group, err := c.ArtifactGroup(branch)
	if err != nil {
		return err
	}
	nodes := make([]string, 0, len(deviceIDs))
	for _, id := range deviceIDs {
		nodes = append(nodes, fmt.Sprintf(`{id: %q}`, id))
	}
	raw, err := c.GraphQL(branch, fmt.Sprintf(
		`mutation { RelationshipAdd(data: {id: %q, name: "members", nodes: [%s]}) { ok } }`,
		group, strings.Join(nodes, ", ")))
	if err != nil {
		return fmt.Errorf("adding %d device(s) to group %s on %s: %w", len(deviceIDs), ArtifactGroupName, branch, err)
	}
	var r struct {
		Add struct {
			OK bool `json:"ok"`
		} `json:"RelationshipAdd"`
	}
	if err := json.Unmarshal(raw, &r); err != nil {
		return err
	}
	if !r.Add.OK {
		return fmt.Errorf("adding devices to group %s on %s: mutation returned ok: false", ArtifactGroupName, branch)
	}
	return nil
}

// GenerateArtifacts renders every group member's artifact on a branch and waits until
// each is Ready.
//
// Infrahub 1.11.2 regenerates nothing on its own: not on a data change, not on a group
// join, not on a commit. Whoever changes a branch regenerates it, and
// for a test branch that is this package. Fylgja itself never generates (Constitution
// III), which is why this lives behind a build tag.
//
// The POST returns in about 0.1s, before anything exists; the artifacts appear Pending
// and reach Ready about a second later, 7.6s from the call for three devices. So
// the wait is for the count and the status, not for the call. The count is the group's
// members, since the definition renders for the group and a device that left it loses
// its artifact on this generation: counting the branch's devices instead let the
// old artifacts satisfy the wait before a removal took effect, and failed it once the
// removal had.
//
// What the wait cannot see is a rewrite. A generation after an intent change rewrites a
// member's artifact in one step, Ready before and Ready after, so it returns as
// soon as the count holds, which it already did. A caller that changed intent waits for
// the checksums to move, as scripts/e2e.sh and the tier-2 tests do.
func (c *Client) GenerateArtifacts(branch string) error {
	definition, err := c.lookupByName(branch, "CoreArtifactDefinition", ArtifactDefinitionName)
	if err != nil {
		return err
	}
	want, err := c.groupMemberCount(branch)
	if err != nil {
		return err
	}
	if want == 0 {
		return fmt.Errorf("group %s on branch %s has no members to generate artifacts for", ArtifactGroupName, branch)
	}

	code, body, err := c.do(http.MethodPost,
		"/api/artifact/generate/"+definition+"?branch="+branch, map[string]any{"nodes": []any{}})
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return fmt.Errorf("generating artifacts on %s: HTTP %d: %s", branch, code, string(body))
	}
	return c.WaitForArtifacts(branch, want)
}

// WaitForArtifacts polls until the branch holds want artifacts and every one is Ready.
// An artifact that stays Pending, or an Error, ends the wait at the deadline naming what
// it saw, so a test fails on the generation rather than later on a missing artifact.
func (c *Client) WaitForArtifacts(branch string, want int) error {
	const q = `{ CoreArtifact { count edges { node { status { value } } } } }`
	deadline := time.Now().Add(ArtifactWait)
	last := "nothing yet"
	for time.Now().Before(deadline) {
		raw, err := c.GraphQL(branch, q)
		if err == nil {
			var r struct {
				A struct {
					Count int `json:"count"`
					Edges []struct {
						Node struct {
							Status struct {
								Value string `json:"value"`
							} `json:"status"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"CoreArtifact"`
			}
			if json.Unmarshal(raw, &r) == nil {
				states := map[string]int{}
				ready := 0
				for _, e := range r.A.Edges {
					states[e.Node.Status.Value]++
					if e.Node.Status.Value == "Ready" {
						ready++
					}
				}
				if r.A.Count == want && ready == want {
					return nil
				}
				last = fmt.Sprintf("%d artifact(s), states %v", r.A.Count, states)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("branch %s: %d artifact(s) did not all reach Ready within %s (last saw %s)",
		branch, want, ArtifactWait, last)
}

// groupMemberCount is how many artifacts a fully generated branch should hold: one per
// member of the group the definition renders for, and none for any other device.
func (c *Client) groupMemberCount(branch string) (int, error) {
	raw, err := c.GraphQL(branch, fmt.Sprintf(
		`{ CoreStandardGroup(name__value: %q) { edges { node { members { count } } } } }`, ArtifactGroupName))
	if err != nil {
		return 0, err
	}
	var r struct {
		G struct {
			Edges []struct {
				Node struct {
					Members struct {
						Count int `json:"count"`
					} `json:"members"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"CoreStandardGroup"`
	}
	if err := json.Unmarshal(raw, &r); err != nil {
		return 0, err
	}
	if len(r.G.Edges) != 1 {
		return 0, fmt.Errorf("branch %s holds %d CoreStandardGroup named %q, want one; "+unregistered,
			branch, len(r.G.Edges), ArtifactGroupName)
	}
	return r.G.Edges[0].Node.Members.Count, nil
}
