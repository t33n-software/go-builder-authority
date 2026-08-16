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
		"00-push-protections.json",
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

func TestPushProtectionsRulesetBlocksCredentialShapedArtifacts(t *testing.T) {
	push := readRepositoryFile(t, filepath.Join("docs", "hosting-platforms", "github", "rulesets", "00-push-protections.json"))
	for _, required := range []string{
		"\"name\": \"push-protections: block secret and key shaped artifacts\"",
		"\"target\": \"push\"",
		"\"source\": \"t33n-software/go-builder-authority\"",
		"\"enforcement\": \"active\"",
		"\"conditions\": null",
		"\"bypass_actors\": []",
		"file_extension_restriction",
		"restricted_file_extensions",
		"file_path_restriction",
		"restricted_file_paths",
	} {
		if !strings.Contains(push, required) {
			t.Fatalf("00-push-protections.json does not contain %q", required)
		}
	}
	for _, extension := range []string{"pem", "key", "p12", "pfx", "jks", "keystore", "kdbx", "ppk", "gpg"} {
		if !strings.Contains(push, "\"*."+extension+"\"") {
			t.Fatalf("00-push-protections.json does not restrict the %q extension in glob form", extension)
		}
	}
	for _, path := range []string{"**/.env", "**/.env.*", "**/credentials", "**/credentials.*", "**/*.tfstate", "**/*.tfstate.*"} {
		if !strings.Contains(push, "\""+path+"\"") {
			t.Fatalf("00-push-protections.json does not restrict the %q path", path)
		}
	}
	for _, forbidden := range []string{"ref_name", "required_status_checks", "code_scanning", "code_quality", "code_coverage"} {
		if strings.Contains(push, forbidden) {
			t.Fatalf("00-push-protections.json unexpectedly contains %q; a push ruleset has no branch targets or check bindings", forbidden)
		}
	}

	readme := normalizeWhitespace(readRepositoryFile(t, filepath.Join("docs", "hosting-platforms", "github", "rulesets", "README.md")))
	for _, required := range []string{"00-push-protections.json", "fork network", "Team plan", "public"} {
		if !strings.Contains(readme, required) {
			t.Fatalf("Ruleset README does not document the push protections token %q", required)
		}
	}
}

func TestModuleIdentityMatchesOrganizationNamespace(t *testing.T) {
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
}

func TestGoToolchainAndBuildToolingContract(t *testing.T) {
	toolsMod := readRepositoryFile(t, filepath.Join("tools", "go.mod"))
	for _, required := range []string{
		"module github.com/t33n-software/go-builder-authority/tools",
		"toolchain go1.26.6",
		"github.com/evilmartians/lefthook/v2",
		"golang.org/x/vuln/cmd/govulncheck",
		"honnef.co/go/tools/cmd/staticcheck",
	} {
		if !strings.Contains(toolsMod, required) {
			t.Fatalf("tools/go.mod does not contain %q", required)
		}
	}
	if _, err := os.Stat(repositoryPath("tools", "go.sum")); err != nil {
		t.Fatalf("tools/go.sum is missing: %v", err)
	}

	ci := readRepositoryFile(t, ".github/workflows/ci.yml")
	for _, required := range []string{
		`go-version: "1.26.6"`,
		`test "$(go env GOVERSION)" = "go1.26.6"`,
		"schedule:",
		"cron:",
	} {
		if !strings.Contains(ci, required) {
			t.Fatalf("CI workflow does not contain %q", required)
		}
	}

	codeql := readRepositoryFile(t, ".github/workflows/codeql.yml")
	for _, required := range []string{
		`go-version: "1.26.6"`,
		`test "$(go env GOVERSION)" = "go1.26.6"`,
	} {
		if !strings.Contains(codeql, required) {
			t.Fatalf("CodeQL workflow does not contain %q", required)
		}
	}

	lefthook := readRepositoryFile(t, "lefthook.yml")
	for _, required := range []string{
		"commit-msg:",
		`git-governance --interactive never commit validate --message-file "{1}"`,
		"pre-push:",
		"go run -mod=readonly ./cmd/build",
	} {
		if !strings.Contains(lefthook, required) {
			t.Fatalf("lefthook.yml does not contain %q", required)
		}
	}

	traceability := readRepositoryFile(t, filepath.Join("docs", "TRACEABILITY.md"))
	if !strings.Contains(traceability, "GBA-4") {
		t.Fatal("TRACEABILITY.md does not contain GBA-4")
	}
}

func normalizeWhitespace(content string) string {
	return strings.Join(strings.Fields(content), " ")
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
