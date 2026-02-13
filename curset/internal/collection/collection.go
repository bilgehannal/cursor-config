package collection

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// CollectionFile represents the top-level collection.json structure.
type CollectionFile struct {
	Collections map[string]Collection `json:"collections"`
}

// Collection represents a single named collection with dynamic object types.
// Keys are object types like "rules", "commands", etc.
// Values are lists of entry names.
type Collection map[string][]string

// Parse parses the collection.json bytes into a CollectionFile.
func Parse(data []byte) (*CollectionFile, error) {
	var cf CollectionFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("failed to parse collection.json: %w", err)
	}
	return &cf, nil
}

// SortedNames returns the collection names sorted alphabetically.
func (cf *CollectionFile) SortedNames() []string {
	names := make([]string, 0, len(cf.Collections))
	for name := range cf.Collections {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// PrintTables renders each collection as its own rounded Unicode table.
func (cf *CollectionFile) PrintTables() {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	names := cf.SortedNames()

	for i, name := range names {
		col := cf.Collections[name]

		// Get sorted object type keys
		keys := make([]string, 0, len(col))
		for k := range col {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Build rows
		rows := make([][]string, 0, len(keys))
		for _, k := range keys {
			rows = append(rows, []string{k, strings.Join(col[k], ", ")})
		}

		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
			Headers(name, "").
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == table.HeaderRow {
					return headerStyle
				}
				return cellStyle
			}).
			Rows(rows...)

		fmt.Println(t)

		if i < len(names)-1 {
			fmt.Println()
		}
	}
}

// ToJSON returns the CollectionFile as pretty-printed JSON.
func (cf *CollectionFile) ToJSON() (string, error) {
	data, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal collection: %w", err)
	}
	return string(data), nil
}
