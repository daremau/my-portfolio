package content

import (
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

func validFrontMatter(title, date string) string {
	return "title: \"" + title + "\"\n" +
		"description: \"A project.\"\n" +
		"category: \"Web Dev\"\n" +
		"pubDate: " + date + "\n" +
		"tags: [\"Go\", \"HTML\"]\n" +
		"link: \"https://example.com\""
}

func TestLoadDirSortsNewestFirst(t *testing.T) {
	dir := t.TempDir()
	writeProject(t, dir, "old.md", validFrontMatter("Old", "2026-01-01"))
	writeProject(t, dir, "new.md", validFrontMatter("New", "2026-02-01"))
	writeProject(t, dir, "middle.md", validFrontMatter("Middle", "2026-01-15"))

	projects, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}

	want := []string{"New", "Middle", "Old"}
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
	writeProject(t, dir, "neural-network-visualizer.md", validFrontMatter("Visualizer", "2026-01-26"))

	projects, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}

	got := projects[0]
	if got.Slug != "neural-network-visualizer" {
		t.Errorf("Slug = %q", got.Slug)
	}
	if got.Category != "Web Dev" {
		t.Errorf("Category = %q", got.Category)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "Go" {
		t.Errorf("Tags = %#v", got.Tags)
	}
	if got.PubDate.IsZero() {
		t.Error("PubDate was not parsed")
	}
}

func TestLoadDirRejectsUnknownCategory(t *testing.T) {
	dir := t.TempDir()
	writeProject(t, dir, "bad.md", strings.Replace(validFrontMatter("Bad", "2026-01-01"), "Web Dev", "Design", 1))

	_, err := LoadDir(dir)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "unknown category") {
		t.Fatalf("error = %v, want unknown category", err)
	}
	if !strings.Contains(err.Error(), "bad.md") {
		t.Fatalf("error = %v, want the file name", err)
	}
}

func TestLoadDirRejectsInvalidLink(t *testing.T) {
	dir := t.TempDir()
	writeProject(t, dir, "bad.md", strings.Replace(validFrontMatter("Bad", "2026-01-01"), "https://example.com", "not-a-url", 1))

	_, err := LoadDir(dir)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "invalid link") {
		t.Fatalf("error = %v, want invalid link", err)
	}
}

func TestLoadDirRequiresDate(t *testing.T) {
	dir := t.TempDir()
	frontMatter := strings.Replace(validFrontMatter("No date", "2026-01-01"), "pubDate: 2026-01-01\n", "", 1)
	writeProject(t, dir, "bad.md", frontMatter)

	_, err := LoadDir(dir)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "pubDate") {
		t.Fatalf("error = %v, want pubDate", err)
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
