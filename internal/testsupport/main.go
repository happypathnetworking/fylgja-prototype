//go:build contract || fixture

package testsupport

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// The registration Infrahub renders from (D-044): this repository, read-only, on main,
// with no credential, which a public repository needs none of.
const (
	RepositoryName     = "fylgja"
	RepositoryLocation = "https://github.com/happypathnetworking/fylgja-prototype.git"
	RepositoryRef      = "main"
)

// The names .infrahub.yml declares. The import is complete when main holds each of them
// once: until then a branch created from main cannot render an artifact.
const (
	QueryName     = "device_config"
	TransformName = "srlinux_device_config"
	ArtifactName  = "device-config"
)

// repositoryDescription is the registration's description: what it is, and why.
const repositoryDescription = "This repository, read-only on main with no credential (D-044)"

// EnsureGroup creates the standard group name on main when it is absent. More than one
// is an error, since every lookup by name wants exactly one.
func (c *Client) EnsureGroup(name string) (created bool, err error) {
	found, err := c.ids("main", "CoreStandardGroup", fmt.Sprintf(`(name__value: %q)`, name))
	if err != nil {
		return false, err
	}
	switch len(found) {
	case 0:
	case 1:
		return false, nil
	default:
		return false, fmt.Errorf("main holds %d CoreStandardGroup named %q, want one", len(found), name)
	}
	if _, err := c.create("main", "CoreStandardGroup",
		fmt.Sprintf(`{name: {value: %q}, label: {value: %q}}`, name, name)); err != nil {
		return false, err
	}
	return true, nil
}

// registeredRepository is a CoreGenericRepository as main holds it.
type registeredRepository struct {
	Kind              string
	Location          string
	Ref               string
	InternalStatus    string
	OperationalStatus string
}

// repository looks a repository up by name on main. It returns nil when there is none.
func (c *Client) repository(name string) (*registeredRepository, error) {
	raw, err := c.GraphQL("main", fmt.Sprintf(`{ CoreGenericRepository(name__value: %q) { edges { node {
		__typename
		location { value }
		internal_status { value }
		operational_status { value }
		... on CoreReadOnlyRepository { ref { value } }
	} } } }`, name))
	if err != nil {
		return nil, fmt.Errorf("looking up repository %s: %w", name, err)
	}
	type value struct {
		Value string `json:"value"`
	}
	var r struct {
		R struct {
			Edges []struct {
				Node struct {
					Typename          string `json:"__typename"`
					Location          value  `json:"location"`
					Ref               value  `json:"ref"`
					InternalStatus    value  `json:"internal_status"`
					OperationalStatus value  `json:"operational_status"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"CoreGenericRepository"`
	}
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, err
	}
	switch len(r.R.Edges) {
	case 0:
		return nil, nil
	case 1:
	default:
		return nil, fmt.Errorf("main holds %d repositories named %q, want one", len(r.R.Edges), name)
	}
	n := r.R.Edges[0].Node
	return &registeredRepository{
		Kind:              n.Typename,
		Location:          n.Location.Value,
		Ref:               n.Ref.Value,
		InternalStatus:    n.InternalStatus.Value,
		OperationalStatus: n.OperationalStatus.Value,
	}, nil
}

// EnsureReadOnlyRepository registers a read-only repository on main, with no credential,
// when none of that name exists. One that exists at another location or ref, or that is
// not read-only, is refused and left as it is: this tool updates and deletes nothing.
func (c *Client) EnsureReadOnlyRepository(name, location, ref string) (created bool, err error) {
	have, err := c.repository(name)
	if err != nil {
		return false, err
	}
	if have != nil {
		if have.Kind != "CoreReadOnlyRepository" {
			return false, fmt.Errorf("repository %s is registered at %s as a %s, not read-only; "+
				"remove it or name another with -repository-name", name, have.Location, have.Kind)
		}
		if have.Location != location || have.Ref != ref {
			return false, fmt.Errorf("repository %s is registered at %s on %s, not %s on %s; "+
				"remove it or name another with -repository-name", name, have.Location, have.Ref, location, ref)
		}
		return false, nil
	}
	if _, err := c.create("main", "CoreReadOnlyRepository", fmt.Sprintf(
		`{name: {value: %q}, location: {value: %q}, ref: {value: %q}, description: {value: %q}}`,
		name, location, ref, repositoryDescription)); err != nil {
		return false, err
	}
	return true, nil
}

// AwaitImport polls once a second until the repository is active and online and main
// holds the query, the transform and the definition .infrahub.yml declares, each once.
// Infrahub imports a registration asynchronously, through its minute-long
// git_repositories_sync, so the wait is bounded by the caller and names what is missing.
func (c *Client) AwaitImport(name string, wait time.Duration) (took time.Duration, err error) {
	start := time.Now()
	deadline := start.Add(wait)
	for {
		missing, status, lookupErr := c.importMissing(name)
		if lookupErr != nil {
			return time.Since(start), lookupErr
		}
		if len(missing) == 0 {
			return time.Since(start), nil
		}
		if !time.Now().Add(time.Second).Before(deadline) {
			return time.Since(start), fmt.Errorf("main did not hold %s within %s; the repository is %s; "+
				"a git_repositories_sync run stuck PENDING in Infrahub's task manager blocks every import until it is cancelled",
				strings.Join(missing, ", "), wait, status)
		}
		time.Sleep(time.Second)
	}
}

// importMissing names what the import has not yet put on main, and the repository's
// internal_status/operational_status.
func (c *Client) importMissing(name string) (missing []string, status string, err error) {
	repo, err := c.repository(name)
	if err != nil {
		return nil, "", err
	}
	if repo == nil {
		return []string{"repository " + name}, "absent", nil
	}
	status = repo.InternalStatus + "/" + repo.OperationalStatus
	if repo.InternalStatus != "active" || repo.OperationalStatus != "online" {
		missing = append(missing, "repository "+name+" active and online")
	}
	for _, o := range []struct{ kind, name, what string }{
		{"CoreGraphQLQuery", QueryName, "query"},
		{"CoreTransformJinja2", TransformName, "transform"},
	} {
		found, err := c.ids("main", o.kind, fmt.Sprintf(`(name__value: %q)`, o.name))
		if err != nil {
			return nil, "", err
		}
		if len(found) != 1 {
			missing = append(missing, fmt.Sprintf("the %s %s", o.what, o.name))
		}
	}
	found, err := c.ids("main", "CoreArtifactDefinition",
		fmt.Sprintf(`(name__value: %q, artifact_name__value: %q)`, ArtifactDefinitionName, ArtifactName))
	if err != nil {
		return nil, "", err
	}
	if len(found) != 1 {
		missing = append(missing, fmt.Sprintf("the definition %s (artifact %s)", ArtifactDefinitionName, ArtifactName))
	}
	return missing, status, nil
}
