package handlers

import (
	"strings"
)

// formatLocation processes a location string, applying specific formatting rules.
func formatLocation(location string) string {
	// Split the location string by hyphens
	parts := strings.Split(location, "-")

	// Process each part of the location
	for i, part := range parts {
		// Replace underscores with spaces
		part = strings.ReplaceAll(part, "_", " ")

		// Special handling for "USA" and "UK" to keep them uppercase
		switch strings.ToLower(part) {
		case "usa", "uk":
			parts[i] = strings.ToUpper(part)
		default:
			// Capitalize the first letter of each word for other parts
			parts[i] = capitalizeWords(part)
		}
	}

	// Rejoin the parts with hyphens
	return strings.Join(parts, "-")
}

