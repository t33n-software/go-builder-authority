package packaging

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestSourceWorkflowsEmitOnlyEstablishedSharedLineChecks(t *testing.T) {
	for _, workflow := range []string{
		".github/workflows/ci.yml",
		".github/workflows/codeql.yml",
		".github/workflows/dependency-review.yml",
	} {
		content := readRepositoryFile(t, workflow)
		for _, forbidden := range []string{"release/**", "support/**"} {
			if strings.Contains(content, forbidden) {
				t.Fatalf("%s unexpectedly targets %q before the governed builder release lifecycle exists", workflow, forbidden)
			}
		}
		for _, required := range []string{"main", "develop"} {
			if !strings.Contains(content, required) {
				t.Fatalf("%s does not target %q", workflow, required)
			}
		}
	}

	ci := readRepositoryFile(t, ".github/workflows/ci.yml")
	for _, required := range []string{
		"Quality gates (linux-amd64)",
		"go run -mod=readonly ./cmd/build",
		"actions/checkout@9f698171ed81b15d1823a05fc7211befd50c8ae0",
		"actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c",
	} {
		if !strings.Contains(ci, required) {
			t.Fatalf("CI workflow does not contain %q", required)
		}
	}

	dependencyReview := readRepositoryFile(t, ".github/workflows/dependency-review.yml")
	for _, required := range []string{
		"Dependency admission review",
		"fail-on-severity: low",
		"fail-on-scopes: runtime,development,unknown",
		"actions/dependency-review-action@2031cfc080254a8a887f58cffee85186f0e49e48",
	} {
		if !strings.Contains(dependencyReview, required) {
			t.Fatalf("dependency review workflow does not contain %q", required)
		}
	}
}

func TestRulesetsCoverOnlyEstablishedBranchFamilies(t *testing.T) {
	rulesetDirectory := repositoryPath("docs", "hosting-platforms", "github", "rulesets")
	entries, err := os.ReadDir(rulesetDirectory)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", rulesetDirectory, err)
	}

	jsonNames := make([]string, 0)
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			jsonNames = append(jsonNames, entry.Name())
		}
	}
	sort.Strings(jsonNames)
	wantNames := []string{
		"01-ticket-working-branches.json",
		"02-develop.json",
		"03-main.json",
	}
	if !reflect.DeepEqual(jsonNames, wantNames) {
		t.Fatalf("ruleset JSON files = %#v, want %#v", jsonNames, wantNames)
	}

	for _, name := range jsonNames {
		content := readRepositoryFile(t, filepath.Join("docs", "hosting-platforms", "github", "rulesets", name))
		for _, forbidden := range []string{"code_quality", "code_coverage"} {
			if strings.Contains(content, forbidden) {
				t.Fatalf("%s contains unsupported %q", name, forbidden)
			}
		}
		if !strings.Contains(content, "\"bypass_actors\": []") {
			t.Fatalf("%s does not prohibit Ruleset bypass actors", name)
		}
	}

	readme := readRepositoryFile(t, filepath.Join("docs", "hosting-platforms", "github", "rulesets", "README.md"))
	for _, required := range []string{
		"Do not import a `release/*` or `support/*` Ruleset yet.",
		"Quality gates (linux-amd64)",
		"Dependency admission review",
	} {
		if !strings.Contains(readme, required) {
			t.Fatalf("Ruleset README does not contain %q", required)
		}
	}
}

func TestAuthorityDocumentationPreservesTenantAndBuilderBoundaries(t *testing.T) {
	for _, path := range []string{
		"README.md",
		"docs/architecture/ADR-0001-GO-BUILDER-AUTHORITY.md",
		"docs/specification/BUILDER-AUTHORITY-CONTRACT.md",
		"docs/operations/BUILDER-EVIDENCE-OPERATIONS.md",
	} {
		content := readRepositoryFile(t, path)
		for _, required := range []string{"tenant", "builder"} {
			if !strings.Contains(strings.ToLower(content), required) {
				t.Fatalf("%s does not document %q boundary", path, required)
			}
		}
	}

	adr := readRepositoryFile(t, "docs/architecture/ADR-0001-GO-BUILDER-AUTHORITY.md")
	for _, forbidden := range []string{
		"tenant configuration, secrets, identities, runtimes, or deployments",
		"application or credential-broker platform release artifacts",
	} {
		if !strings.Contains(adr, forbidden) {
			t.Fatalf("ADR does not exclude %q", forbidden)
		}
	}
}

func readRepositoryFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(repositoryPath(filepath.FromSlash(path)))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	return string(content)
}

func repositoryPath(parts ...string) string {
	return filepath.Join(append([]string{"..", ".."}, parts...)...)
}
