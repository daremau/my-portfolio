// Package content loads and validates the portfolio's Markdown project files.
package content

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/adrg/frontmatter"
)

// Project is a single portfolio entry parsed from a Markdown file.
type Project struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Stack       []string `yaml:"stack"`
	Cover       string   `yaml:"cover"`
	Repo        string   `yaml:"repo"`
	Demo        string   `yaml:"demo"`
	Year        int      `yaml:"year"`
	Order       int      `yaml:"order"`

	// Slug is derived from the file name.
	Slug string `yaml:"-"`
}

// StackLine joins the project's stack for display.
func (p Project) StackLine() string {
	return strings.Join(p.Stack, " · ")
}

// LoadDir reads every .md file in dir, validates it, and returns the projects
// sorted by their display order.
func LoadDir(dir string) ([]Project, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read content directory: %w", err)
	}

	projects := make([]Project, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			continue
		}

		project, err := loadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	sort.SliceStable(projects, func(i, j int) bool {
		if projects[i].Order == projects[j].Order {
			return projects[i].Title < projects[j].Title
		}
		return projects[i].Order < projects[j].Order
	})

	return projects, nil
}

func loadFile(path string) (Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Project{}, fmt.Errorf("read %s: %w", path, err)
	}

	var project Project
	if _, err := frontmatter.Parse(bytes.NewReader(data), &project); err != nil {
		return Project{}, fmt.Errorf("parse %s: %w", filepath.Base(path), err)
	}

	project.Slug = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	if err := project.Validate(); err != nil {
		return Project{}, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}

	return project, nil
}

// Validate checks the fields required to render a project.
func (p Project) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(p.Description) == "" {
		return errors.New("description is required")
	}
	if len(p.Stack) == 0 {
		return errors.New("stack is required")
	}
	for i, item := range p.Stack {
		if strings.TrimSpace(item) == "" {
			return fmt.Errorf("stack[%d] is empty", i)
		}
	}
	if strings.TrimSpace(p.Cover) == "" {
		return errors.New("cover is required")
	}
	if !filepath.IsLocal(filepath.FromSlash(p.Cover)) {
		return fmt.Errorf("cover %q must be a local path inside public/", p.Cover)
	}
	if p.Order <= 0 {
		return errors.New("order must be greater than 0")
	}
	if p.Year != 0 && (p.Year < 1990 || p.Year > 2100) {
		return fmt.Errorf("year %d is out of range", p.Year)
	}
	if p.Repo != "" {
		if err := validateURL("repo", p.Repo); err != nil {
			return err
		}
	}
	if p.Demo != "" {
		if err := validateURL("demo", p.Demo); err != nil {
			return err
		}
	}
	return nil
}

func validateURL(name, raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid %s URL %q: %w", name, raw, err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("invalid %s URL %q: scheme must be http or https", name, raw)
	}
	if parsed.Host == "" {
		return fmt.Errorf("invalid %s URL %q: missing host", name, raw)
	}
	return nil
}
