// internal/shell/history/manager.go
package history

import (
	"bufio"
	"os"
	"path/filepath"
)

type Manager struct {
	entries    []string
	maxEntries int
	filePath   string
}

func NewManager(historyFile string) (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(homeDir, historyFile)
	manager := &Manager{
		entries:    make([]string, 0),
		maxEntries: 1000,
		filePath:   filePath,
	}

	// Create history file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	// Load existing history
	if err := manager.load(); err != nil {
		return nil, err
	}

	return manager, nil
}

func (m *Manager) Add(command string) error {
	m.entries = append(m.entries, command)
	if len(m.entries) > m.maxEntries {
		m.entries = m.entries[1:]
	}
	return m.save()
}

func (m *Manager) Save() error {
	return m.save()
}

func (m *Manager) Close() error {
	return m.save()
}

func (m *Manager) load() error {
	file, err := os.Open(m.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m.entries = append(m.entries, scanner.Text())
	}
	return scanner.Err()
}

func (m *Manager) save() error {
	file, err := os.Create(m.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, entry := range m.entries {
		if _, err := writer.WriteString(entry + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}
