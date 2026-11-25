package components

import "strings"

const (
	// GitHubURL is the public GitHub profile.
	GitHubURL = "https://github.com/daremau"
	// AllReposURL lists every public repository.
	AllReposURL = GitHubURL + "?tab=repositories"
	// Email is the public contact address.
	Email = "daremau@pm.me"
	// EmailURL is the mailto link used by the templates.
	EmailURL = "mailto:" + Email
)

// OtherProject is a compact entry in the "More on GitHub" list.
type OtherProject struct {
	Name        string
	Description string
	URL         string
	Language    string
}

// OtherProjects returns a few additional public repositories.
func OtherProjects() []OtherProject {
	return []OtherProject{
		{
			Name:        "Local News Summarizer",
			Description: "Scrapes and summarizes top news from Paraguay's two largest outlets.",
			URL:         "https://github.com/daremau/local-news-summarizer",
			Language:    "JavaScript",
		},
		{
			Name:        "Godot MCP Extended",
			Description: "MCP server for the Godot engine: inspect, edit, validate, run, and QA games.",
			URL:         "https://github.com/daremau/godot-mcp-extended",
			Language:    "TypeScript",
		},
		{
			Name:        "Agent Chat Share",
			Description: "Shares an in-progress chat between coding agents.",
			URL:         "https://github.com/daremau/agent-chat-share",
			Language:    "Rust",
		},
		{
			Name:        "SQL Data Warehouse",
			Description: "SQL Server data warehouse with ETL, data modeling, and analytics.",
			URL:         "https://github.com/daremau/sql-data-warehouse",
			Language:    "T-SQL",
		},
	}
}

// SkillGroup is one row in the "Tools I use" list.
type SkillGroup struct {
	Label string
	Items []string
}

// List joins a group's items for display.
func (g SkillGroup) List() string {
	return strings.Join(g.Items, " · ")
}

// SkillGroups returns the tools grouped by area.
func SkillGroups() []SkillGroup {
	return []SkillGroup{
		{Label: "Languages", Items: []string{"TypeScript", "JavaScript", "Python", "SQL"}},
		{Label: "Frontend", Items: []string{"React", "Next.js", "Tailwind CSS"}},
		{Label: "Backend", Items: []string{"Node.js", "NestJS", "Prisma"}},
		{Label: "Data", Items: []string{"PostgreSQL", "Supabase", "Redis", "pgvector"}},
		{Label: "AI / ML", Items: []string{"OpenAI", "RAG", "PyTorch", "TensorFlow", "Pandas"}},
		{Label: "Tooling", Items: []string{"Docker", "AWS", "Git", "Playwright", "Vitest"}},
	}
}

// ProjectCardClass returns the layout class for a project card.
func ProjectCardClass(featured bool) string {
	if featured {
		return "project-card project-card-featured"
	}
	return "project-card"
}
