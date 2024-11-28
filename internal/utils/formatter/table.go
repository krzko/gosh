// internal/utils/formatter/table.go
package formatter

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/krzko/gosh/internal/utils/color"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

type TableFormatter struct {
	theme *color.Theme
}

type FileInfo struct {
	Name        string
	Size        int64
	Mode        os.FileMode
	ModTime     time.Time
	IsDir       bool
	Owner       string
	Group       string
	Permissions string
}

func New() *TableFormatter {
	return &TableFormatter{
		theme: color.DefaultTheme(),
	}
}

func stripANSI(str string) string {
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansi.ReplaceAllString(str, "")
}

func (t *TableFormatter) FormatLongList(entries []FileInfo) error {
	table := tablewriter.NewWriter(os.Stdout)

	// Configure table
	table.SetHeader([]string{"Permissions", "Owner", "Group", "Size", "Modified", "Name"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("  ")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	// Add entries
	for _, entry := range entries {
		row := []string{
			t.theme.ColorizePermissions(entry.Permissions),
			entry.Owner,
			entry.Group,
			formatSize(entry.Size),
			entry.ModTime.Format("Jan _2 15:04"),
			t.theme.ColorizeName(entry.Name, entry.IsDir, entry.Mode),
		}
		table.Append(row)
	}

	table.Render()
	return nil
}

// FormatSimpleList provides a simpler listing format
func (t *TableFormatter) FormatSimpleList(entries []FileInfo) error {
	for _, entry := range entries {
		fmt.Println(t.theme.ColorizeName(entry.Name, entry.IsDir, entry.Mode))
	}
	return nil
}

func (t *TableFormatter) FormatCompact(entries []FileInfo, width int) error {
	if len(entries) == 0 {
		return nil
	}

	// Get terminal width if not provided
	if width <= 0 {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			width = w
		} else {
			width = 80 // fallback width
		}
	}

	// First pass: calculate maximum name length
	maxLen := 0
	for _, entry := range entries {
		if len(entry.Name) > maxLen {
			maxLen = len(entry.Name)
		}
	}

	// Add minimum spacing between columns
	colWidth := maxLen + 2
	cols := width / colWidth
	if cols < 1 {
		cols = 1
	}

	// Calculate number of rows needed
	numEntries := len(entries)
	rows := (numEntries + cols - 1) / cols

	// Create a matrix of entries
	matrix := make([][]string, rows)
	for i := range matrix {
		matrix[i] = make([]string, cols)
	}

	// Fill the matrix
	for i, entry := range entries {
		col := i / rows
		row := i % rows
		name := entry.Name
		if entry.IsDir {
			name = t.theme.ColorizeName(name, true, entry.Mode)
		} else {
			name = t.theme.ColorizeName(name, false, entry.Mode)
		}
		if col < cols {
			matrix[row][col] = name
		}
	}

	// Print the matrix
	for _, row := range matrix {
		for colIdx, name := range row {
			if name == "" {
				continue
			}
			if colIdx > 0 {
				fmt.Print(strings.Repeat(" ", colWidth-len(stripANSI(row[colIdx-1]))))
			}
			fmt.Print(name)
		}
		fmt.Println()
	}

	return nil
}

// FormatCompactList provides a compact multi-column listing
func (t *TableFormatter) FormatCompactList(entries []FileInfo) error {
	termWidth := 80 // You might want to get actual terminal width
	maxNameLength := 0

	// Find the longest name
	for _, entry := range entries {
		if len(entry.Name) > maxNameLength {
			maxNameLength = len(entry.Name)
		}
	}

	// Add padding
	columnWidth := maxNameLength + 2
	columns := termWidth / columnWidth
	if columns == 0 {
		columns = 1
	}

	// Print entries in columns
	for i, entry := range entries {
		if i > 0 && i%columns == 0 {
			fmt.Println()
		}
		fmt.Printf("%-*s", columnWidth, t.theme.ColorizeName(entry.Name, entry.IsDir, entry.Mode))
	}
	fmt.Println()

	return nil
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(size)/float64(div), "KMGTPE"[exp])
}
