package packaging

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// bindingManifest mirrors the tenant binding manifest (repo-bindings/v1) for
// the self-consistency proofs of the canonical adoption. The home-side proof
// against the canonical masters is owned by the verify-canonical tool; these
// tests bind the tenant files to the manifest.
type bindingManifest struct {
	Home struct {
		Repository string `json:"repository"`
		SHA        string `json:"sha"`
	} `json:"home"`
	Callers []struct {
		File   string `json:"file"`
		Master string `json:"master"`
		SHA256 string `json:"sha256"`
	} `json:"callers"`
	Files struct {
		Lefthook      fileBinding `json:"lefthook"`
		Gitattributes fileBinding `json:"gitattributes"`
		Gitignore     fileBinding `json:"gitignore"`
		Dependabot    fileBinding `json:"dependabot"`
	} `json:"files"`
	Codeowners struct {
		Path         string `json:"path"`
		DefaultOwner string `json:"defaultOwner"`
	} `json:"codeowners"`
}

type fileBinding struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

func readBindingManifest(t *testing.T) bindingManifest {
	t.Helper()
	var manifest bindingManifest
	if err := json.Unmarshal([]byte(readRepositoryFile(t, "repo-bindings.json")), &manifest); err != nil {
		t.Fatalf("repo-bindings.json is not valid JSON: %v", err)
	}
	if manifest.Home.Repository != "t33n-software/repository-governance" {
		t.Fatalf("the manifest binds home %q", manifest.Home.Repository)
	}
	return manifest
}

// hashRepositoryFile hashes the LF-normalized repository file; the canonical
// .gitattributes makes the checkout LF, and the normalization keeps the
// derivation tolerant as the second line of defense.
func hashRepositoryFile(t *testing.T, path string) string {
	t.Helper()
	normalized := strings.ReplaceAll(readRepositoryFile(t, path), "\r\n", "\n")
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}

func TestCanonicalCallersMatchTheBindingManifest(t *testing.T) {
	manifest := readBindingManifest(t)
	want := map[string]string{
		".github/workflows/ci.yml":                "hosting-platforms/github/workflows/callers/go/ci.yml",
		".github/workflows/codeql.yml":            "hosting-platforms/github/workflows/callers/go/codeql.yml",
		".github/workflows/dependency-review.yml": "hosting-platforms/github/workflows/callers/go/dependency-review.yml",
	}
	if len(manifest.Callers) != len(want) {
		t.Fatalf("the manifest carries %d callers, want %d", len(manifest.Callers), len(want))
	}
	for _, caller := range manifest.Callers {
		master, found := want[caller.File]
		if !found {
			t.Fatalf("the manifest carries an unexpected caller %q", caller.File)
		}
		if caller.Master != master {
			t.Fatalf("caller %q binds master %q, want %q", caller.File, caller.Master, master)
		}
		if hash := hashRepositoryFile(t, caller.File); hash != caller.SHA256 {
			t.Fatalf("the tenant caller %s hashes to %s, want the bound %s", caller.File, hash, caller.SHA256)
		}
		content := readRepositoryFile(t, caller.File)
		if !strings.Contains(content, "uses: "+manifest.Home.Repository+"/.github/workflows/reusable-") {
			t.Fatalf("the tenant caller %s does not reference a home payload", caller.File)
		}
		if !strings.Contains(content, "@"+manifest.Home.SHA) {
			t.Fatalf("the tenant caller %s does not pin the bound home SHA", caller.File)
		}
		if !strings.Contains(content, `branches: [main, develop, "release/**", "support/**"]`) {
			t.Fatalf("the tenant caller %s does not cover every shared line", caller.File)
		}
	}
}

func TestCanonicalFileFamilyMatchesTheBindingManifest(t *testing.T) {
	manifest := readBindingManifest(t)
	for _, topic := range []fileBinding{
		manifest.Files.Lefthook,
		manifest.Files.Gitattributes,
		manifest.Files.Dependabot,
	} {
		if hash := hashRepositoryFile(t, topic.Path); hash != topic.SHA256 {
			t.Fatalf("the canonical file %s hashes to %s, want the bound %s", topic.Path, hash, topic.SHA256)
		}
	}
	// The gitignore topic is prefix-mode in the home verifier: the canonical
	// core is a verbatim prefix and project additions live below the mark.
	gitignore := readRepositoryFile(t, manifest.Files.Gitignore.Path)
	if !strings.Contains(gitignore, "# -- project additions below this line --\n") {
		t.Fatal("the gitignore does not carry the canonical core with the project-block mark")
	}

	codeowners := readRepositoryFile(t, manifest.Codeowners.Path)
	if !strings.Contains(codeowners, "* "+manifest.Codeowners.DefaultOwner) {
		t.Fatalf("the ownership file does not bind the default owner %q", manifest.Codeowners.DefaultOwner)
	}
}

func TestConformanceWorkflowBindsTheVerifier(t *testing.T) {
	manifest := readBindingManifest(t)
	content := readRepositoryFile(t, ".github/workflows/canonical-conformance.yml")
	for _, required := range []string{
		"permissions: {}",
		"name: Canonical conformance",
		"uses: " + manifest.Home.Repository + "/.github/actions/verify-canonical-files@" + manifest.Home.SHA,
		`branches: [main, develop, "release/**", "support/**"]`,
	} {
		if !strings.Contains(content, required) {
			t.Fatalf("the canonical conformance workflow does not contain %q", required)
		}
	}
}

func TestOrganizationRulesetAdoptionHasNoLocalLegacyDefinitions(t *testing.T) {
	if _, err := os.Stat(repositoryPath("docs", "hosting-platforms")); !os.IsNotExist(err) {
		t.Fatalf("legacy ruleset location must not exist")
	}

	conventions := readRepositoryFile(t, filepath.Join("docs", "conventions", "hosting-plattform", "github", "rule-sets", "README.md"))
	for _, required := range []string{
		"git-governance",
		"quality-gates=linux-only",
		"~ALL",
	} {
		if !strings.Contains(conventions, required) {
			t.Fatalf("rule-set conventions README does not contain %q", required)
		}
	}
}

func TestModuleIdentityAndQualityContract(t *testing.T) {
	goMod := readRepositoryFile(t, "go.mod")
	for _, required := range []string{
		"module github.com/t33n-software/go-builder-authority",
		"go 1.26",
		"toolchain go1.26.6",
	} {
		if !strings.Contains(goMod, required) {
			t.Fatalf("go.mod does not contain %q", required)
		}
	}

	quality := readRepositoryFile(t, "git-governance.quality.json")
	var qualityConfig struct {
		SchemaVersion int `json:"schemaVersion"`
		Toolchain     struct {
			Language string `json:"language"`
			Version  string `json:"version"`
		} `json:"toolchain"`
		Extends  []string        `json:"extends"`
		Defaults json.RawMessage `json:"defaults"`
		Gates    []struct {
			Name    string   `json:"name"`
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"gates"`
		Project json.RawMessage `json:"project"`
	}
	if err := json.Unmarshal([]byte(quality), &qualityConfig); err != nil {
		t.Fatalf("git-governance.quality.json is not valid JSON: %v", err)
	}
	if qualityConfig.SchemaVersion != 4 {
		t.Fatalf("git-governance.quality.json carries schemaVersion %d, want 4", qualityConfig.SchemaVersion)
	}
	if qualityConfig.Toolchain.Language != "go" || qualityConfig.Toolchain.Version != "1.26.6" {
		t.Fatalf("git-governance.quality.json carries toolchain %#v, want go 1.26.6", qualityConfig.Toolchain)
	}
	if qualityConfig.Extends == nil || len(qualityConfig.Extends) != 0 {
		t.Fatalf("git-governance.quality.json carries extends %#v, want an explicit empty list", qualityConfig.Extends)
	}
	if len(qualityConfig.Defaults) != 0 {
		t.Fatalf("git-governance.quality.json must rely on the schema-owned defaults, got %s", qualityConfig.Defaults)
	}
	if len(qualityConfig.Project) != 0 {
		t.Fatalf("git-governance.quality.json must not declare repository-local gate-chain binaries, got %s", qualityConfig.Project)
	}
	if len(qualityConfig.Gates) != 1 {
		t.Fatalf("git-governance.quality.json carries %d gates, want the canonical gate chain", len(qualityConfig.Gates))
	}
	gate := qualityConfig.Gates[0]
	if gate.Name != "go-builder-authority-source-quality" ||
		gate.Command != "go" ||
		!slices.Equal(gate.Args, []string{"tool", "-modfile", "tools/go.mod", "quality-gate"}) {
		t.Fatalf("the quality gate does not invoke the canonical gate chain through the tooling module pin: %#v", gate)
	}
	if _, err := os.Stat(repositoryPath("cmd")); !os.IsNotExist(err) {
		t.Fatalf("the repository must not retain a repo-local gate-chain directory: %v", err)
	}

	lefthook := readRepositoryFile(t, "lefthook.yml")
	if !strings.Contains(lefthook, "git-governance --interactive never validate pre-push --remote") {
		t.Fatal("lefthook.yml does not bind the canonical pre-push validation")
	}
}

func TestGoToolchainAndBuildToolingContract(t *testing.T) {
	toolsMod := readRepositoryFile(t, filepath.Join("tools", "go.mod"))
	for _, required := range []string{
		"module github.com/t33n-software/go-builder-authority/tools",
		"toolchain go1.26.6",
		"github.com/evilmartians/lefthook/v2",
		"golang.org/x/vuln/cmd/govulncheck",
		"honnef.co/go/tools/cmd/staticcheck",
		"github.com/t33n-software/go-quality-authority/cmd/quality-gate",
		"github.com/t33n-software/go-quality-authority/cmd/check-coverage",
		"github.com/t33n-software/repository-governance/cmd/verify-canonical",
	} {
		if !strings.Contains(toolsMod, required) {
			t.Fatalf("tools/go.mod does not contain %q", required)
		}
	}
	if _, err := os.Stat(repositoryPath("tools", "go.sum")); err != nil {
		t.Fatalf("tools/go.sum is missing: %v", err)
	}

	manifest := readBindingManifest(t)
	for _, caller := range []string{"ci.yml", "codeql.yml"} {
		content := readRepositoryFile(t, ".github/workflows/"+caller)
		if !strings.Contains(content, "uses: "+manifest.Home.Repository+"/.github/workflows/reusable-") {
			t.Fatalf("the caller %s does not reference a home payload", caller)
		}
	}

	lefthook := readRepositoryFile(t, "lefthook.yml")
	for _, required := range []string{
		"commit-msg:",
		`git-governance --interactive never commit validate --message-file "{1}"`,
		"pre-push:",
		`git-governance --interactive never validate pre-push --remote "{1}"`,
	} {
		if !strings.Contains(lefthook, required) {
			t.Fatalf("lefthook.yml does not contain %q", required)
		}
	}

	traceability := readRepositoryFile(t, filepath.Join("docs", "TRACEABILITY.md"))
	for _, ticket := range []string{"GBA-4", "GBA-8"} {
		if !strings.Contains(traceability, ticket) {
			t.Fatalf("TRACEABILITY.md does not contain %s", ticket)
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
