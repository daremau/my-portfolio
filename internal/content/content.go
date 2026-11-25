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
	"time"

	"github.com/adrg/frontmatter"
)

// Categories are the project categories supported by the site.
var Categories = []string{"Web Dev", "Data Science", "AI Engineering"}

// Project is a single portfolio entry parsed from a Markdown file.
type Project struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Category    string   `yaml:"category"`
	PubDate     Date     `yaml:"pubDate"`
	Image       string   `yaml:"image"`
	Tags        []string `yaml:"tags"`
	Link        string   `yaml:"link"`

	// Slug and Body are parsed from the file itself, not the front matter.
	Slug string `yaml:"-"`
	Body string `yaml:"-"`
}

// Date accepts both YAML timestamps and "YYYY-MM-DD" strings.
type Date struct {
	time.Time
}

// UnmarshalYAML implements the yaml.Unmarshaler interface used by
// gopkg.in/yaml.v2, which the frontmatter package relies on.
func (d *Date) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var t time.Time
	if err := unmarshal(&t); err == nil {
		d.Time = t
		return nil
	}

	var raw string
	if err := unmarshal(&raw); err != nil {
		return fmt.Errorf("expected a date: %w", err)
	}

	for _, layout := range []string{"2006-01-02", time.RFC3339} {
		if t, err := time.Parse(layout, raw); err == nil {
			d.Time = t
			return nil
		}
	}

	return fmt.Errorf("invalid date %q (want YYYY-MM-DD)", raw)
}

// LoadDir reads every .md file in dir, validates it, and returns the projects
// sorted newest first.
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
		if projects[i].PubDate.Equal(projects[j].PubDate.Time) {
			return projects[i].Title < projects[j].Title
		}
		return projects[i].PubDate.After(projects[j].PubDate.Time)
	})

	return projects, nil
}

func loadFile(path string) (Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Project{}, fmt.Errorf("read %s: %w", path, err)
	}

	var project Project
	body, err := frontmatter.Parse(bytes.NewReader(data), &project)
	if err != nil {
		return Project{}, fmt.Errorf("parse %s: %w", filepath.Base(path), err)
	}

	project.Slug = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	project.Body = string(body)

	if err := project.Validate(); err != nil {
		return Project{}, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}

	return project, nil
}

// Validate checks the required fields and their formats.
func (p Project) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(p.Description) == "" {
		return errors.New("description is required")
	}
	if !isCategory(p.Category) {
		return fmt.Errorf("unknown category %q (want one of %s)", p.Category, strings.Join(Categories, ", "))
	}
	if p.PubDate.IsZero() {
		return errors.New("pubDate is required")
	}
	if p.Link == "" {
		return errors.New("link is required")
	}
	if err := validateLink(p.Link); err != nil {
		return err
	}
	return nil
}

func isCategory(category string) bool {
	for _, candidate := range Categories {
		if category == candidate {
			return true
		}
	}
	return false
}

func validateLink(link string) error {
	parsed, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("invalid link %q: %w", link, err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("invalid link %q: scheme must be http or https", link)
	}
	if parsed.Host == "" {
		return fmt.Errorf("invalid link %q: missing host", link)
	}
	return nil
}
