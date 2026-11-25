package content

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeProject(t *testing.T, dir, name, frontMatter string) {
	t.Helper()

	body := "---\n" + frontMatter + "\n---\n"
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
}

func validFrontMatter(title string, order int) string {
	return fmt.Sprintf(`title: %q
description: "A project."
stack: ["Go", "HTML"]
cover: "projects/example.png"
repo: "https://github.com/daremau/example"
year: 2026
order: %d`, title, order)
}

func TestLoadDirSortsByOrder(t *testing.T) {
	dir := t.TempDir()
	writeProject(t, dir, "third.md", validFrontMatter("Third", 3))
	writeProject(t, dir, "first.md", validFrontMatter("First", 1))
	writeProject(t, dir, "second.md", validFrontMatter("Second", 2))

	projects, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}

	want := []string{"First", "Second", "Third"}
	if len(projects) != len(want) {
		t.Fatalf("got %d projects, want %d", len(projects), len(want))
	}
	for i, title := range want {
		if projects[i].Title != title {
			t.Errorf("project %d = %q, want %q", i, projects[i].Title, title)
		}
	}
}

func TestLoadDirParsesFields(t *testing.T) {
	dir := t.TempDir()
	writeProject(t, dir, "atlas-crm.md", validFrontMatter("Atlas CRM", 1))

	projects, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}

	got := projects[0]
	if got.Slug != "atlas-crm" {
		t.Errorf("Slug = %q", got.Slug)
	}
	if got.Cover != "projects/example.png" {
		t.Errorf("Cover = %q", got.Cover)
	}
	if got.Order != 1 {
		t.Errorf("Order = %d", got.Order)
	}
	if got.Year != 2026 {
		t.Errorf("Year = %d", got.Year)
	}
	if got.StackLine() != "Go · HTML" {
		t.Errorf("StackLine = %q", got.StackLine())
	}
}

func TestLoadDirRejectsMissingStack(t *testing.T) {
	dir := t.TempDir()
	frontMatter := strings.Replace(validFrontMatter("Bad", 1), "stack: [\"Go\", \"HTML\"]\n", "", 1)
	writeProject(t, dir, "bad.md", frontMatter)

	_, err := LoadDir(dir)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "stack is required") {
		t.Fatalf("error = %v, want stack is required", err)
	}
	if !strings.Contains(err.Error(), "bad.md") {
		t.Fatalf("error = %v, want the file name", err)
	}
}

func TestLoadDirRejectsInvalidRepo(t *testing.T) {
	dir := t.TempDir()
	frontMatter := strings.Replace(validFrontMatter("Bad", 1), "https://github.com/daremau/example", "not-a-url", 1)
	writeProject(t, dir, "bad.md", frontMatter)

	_, err := LoadDir(dir)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "invalid repo URL") {
		t.Fatalf("error = %v, want invalid repo URL", err)
	}
}

func TestLoadDirRejectsCoverOutsidePublic(t *testing.T) {
	dir := t.TempDir()
	frontMatter := strings.Replace(validFrontMatter("Bad", 1), "projects/example.png", "../secret.png", 1)
	writeProject(t, dir, "bad.md", frontMatter)

	_, err := LoadDir(dir)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "must be a local path") {
		t.Fatalf("error = %v, want cover path error", err)
	}
}

func TestLoadDirRequiresOrder(t *testing.T) {
	dir := t.TempDir()
	frontMatter := strings.Replace(validFrontMatter("No order", 1), "order: 1", "", 1)
	writeProject(t, dir, "bad.md", frontMatter)

	_, err := LoadDir(dir)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "order") {
		t.Fatalf("error = %v, want order error", err)
	}
}

func TestLoadDirAllowsEmptyDirectory(t *testing.T) {
	projects, err := LoadDir(t.TempDir())
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	if len(projects) != 0 {
		t.Fatalf("got %d projects, want 0", len(projects))
	}
}
