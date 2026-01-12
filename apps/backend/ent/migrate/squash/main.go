//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"
)

const migrationsDir = "ent/migrate/migrations"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))

	checkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	uncheckStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true)
)

type model struct {
	files    []string
	selected map[int]bool
	cursor   int
	done     bool
	aborted  bool
}

func initialModel(files []string) model {
	selected := make(map[int]bool)
	for i := range files {
		selected[i] = true
	}
	return model{
		files:    files,
		selected: selected,
		cursor:   0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.aborted = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		case " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "a":
			for i := range m.files {
				m.selected[i] = true
			}
		case "n":
			for i := range m.files {
				m.selected[i] = false
			}
		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.done || m.aborted {
		return ""
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("🗜️  Select migrations to squash"))
	b.WriteString("\n\n")

	for i, file := range m.files {
		cursor := "  "
		if m.cursor == i {
			cursor = "▸ "
		}

		checked := uncheckStyle.Render("○")
		if m.selected[i] {
			checked = checkStyle.Render("●")
		}

		style := itemStyle
		if m.cursor == i {
			style = selectedItemStyle
		}

		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, checked, style.Render(file)))
	}

	// Count selected
	count := 0
	for _, s := range m.selected {
		if s {
			count++
		}
	}

	b.WriteString(helpStyle.Render(fmt.Sprintf("\n%d of %d selected", count, len(m.files))))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/k up • ↓/j down • space toggle • a all • n none • enter confirm • q quit"))

	return b.String()
}

func main() {
	squashName := "squashed"
	if len(os.Args) >= 2 {
		squashName = os.Args[1]
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		logrus.Fatalf("failed to read migrations directory: %v", err)
	}

	var sqlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}

	if len(sqlFiles) <= 1 {
		fmt.Println("Only 1 or no migrations found, nothing to squash.")
		return
	}

	sort.Strings(sqlFiles)

	// Run bubbletea program
	p := tea.NewProgram(initialModel(sqlFiles))
	finalModel, err := p.Run()
	if err != nil {
		logrus.Fatalf("error running program: %v", err)
	}

	m := finalModel.(model)

	if m.aborted {
		fmt.Println("Aborted.")
		return
	}

	// Collect selected files
	var selectedFiles []string
	for i, f := range m.files {
		if m.selected[i] {
			selectedFiles = append(selectedFiles, f)
		}
	}

	if len(selectedFiles) < 2 {
		fmt.Println("Need at least 2 migrations to squash.")
		return
	}

	fmt.Printf("\n%s\n", warningStyle.Render("⚠️  WARNING: Only do this if no other environments have applied these migrations!"))
	fmt.Printf("\nSquashing %d migrations:\n", len(selectedFiles))
	for _, f := range selectedFiles {
		fmt.Printf("  • %s\n", f)
	}
	fmt.Print("\nContinue? [y/N]: ")

	var response string
	fmt.Scanln(&response)
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		fmt.Println("Aborted.")
		return
	}

	// Combine selected SQL content
	var combinedSQL strings.Builder
	combinedSQL.WriteString("-- Squashed migration combining the following migrations:\n")
	for _, f := range selectedFiles {
		combinedSQL.WriteString(fmt.Sprintf("-- %s\n", f))
	}
	combinedSQL.WriteString("\n")

	for _, f := range selectedFiles {
		content, err := os.ReadFile(filepath.Join(migrationsDir, f))
		if err != nil {
			logrus.Fatalf("failed to read migration file %s: %v", f, err)
		}
		combinedSQL.WriteString(fmt.Sprintf("-- === From: %s ===\n", f))
		combinedSQL.WriteString(string(content))
		combinedSQL.WriteString("\n")
	}

	timestamp := time.Now().Format("20060102150405")
	newFilename := fmt.Sprintf("%s_%s.sql", timestamp, squashName)
	newPath := filepath.Join(migrationsDir, newFilename)

	backupDir := filepath.Join(migrationsDir, fmt.Sprintf(".backup_%s", timestamp))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		logrus.Fatalf("failed to create backup directory: %v", err)
	}

	for _, f := range selectedFiles {
		oldPath := filepath.Join(migrationsDir, f)
		backupPath := filepath.Join(backupDir, f)
		if err := os.Rename(oldPath, backupPath); err != nil {
			logrus.Fatalf("failed to backup migration %s: %v", f, err)
		}
	}

	if err := os.WriteFile(newPath, []byte(combinedSQL.String()), 0644); err != nil {
		logrus.Fatalf("failed to write squashed migration: %v", err)
	}

	fmt.Printf("\n✅ Created squashed migration: %s\n", newFilename)
	fmt.Printf("📁 Backed up old migrations to: %s\n", backupDir)

	fmt.Println("\n🔄 Regenerating atlas.sum...")
	cmd := exec.Command("atlas", "migrate", "hash", "--dir", fmt.Sprintf("file://%s", migrationsDir))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Fatalf("failed to regenerate atlas.sum: %v", err)
	}

	fmt.Println("\n🎉 Squash complete!")
	fmt.Println("\n📋 Next steps:")
	fmt.Println("   1. Reset your local atlas_schema_revisions table")
	fmt.Println("   2. Run: pnpm run db:migrate:apply")
	fmt.Println("   3. If issues occur, restore migrations from the backup directory")
}
