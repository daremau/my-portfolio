package components

import "strconv"

// Skills shown in the About section.
func Skills() []string {
	return []string{
		"Python", "JavaScript", "TypeScript", "React", "Node.js", "PostgreSQL",
		"TensorFlow", "PyTorch", "Pandas", "Docker", "AWS", "Git",
	}
}

// FilterClass returns the classes for a project filter tab.
func FilterClass(index int) string {
	if index == 0 {
		return "filter-btn filter-btn-active"
	}
	return "filter-btn"
}

// AriaSelected returns the aria-selected value for a filter tab.
func AriaSelected(index int) string {
	if index == 0 {
		return "true"
	}
	return "false"
}

// ProjectClass returns the wrapper classes for a project card.
func ProjectClass(index int) string {
	return "project-item fade-in stagger-" + strconv.Itoa(index%6+1)
}

// LimitTags returns at most n tags.
func LimitTags(tags []string, n int) []string {
	if len(tags) <= n {
		return tags
	}
	return tags[:n]
}

// CategoryClass returns the badge classes for a project category.
func CategoryClass(category string) string {
	switch category {
	case "Web Dev":
		return "badge badge-web"
	case "Data Science":
		return "badge badge-data"
	case "AI Engineering":
		return "badge badge-ai"
	default:
		return "badge badge-default"
	}
}
