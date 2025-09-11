package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/sureshkrishnan-v/lazynet/internal/actions"
)

// setKeybindings configures all the key event handlers for the application.
func (a *App) setKeybindings() {
	// Global keybindings
	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If the filter input has focus, let it handle the key event
		if a.view.filterInput.HasFocus() {
			return event
		}

		switch event.Rune() {
		case 'q':
			a.Stop()
			return nil
		case 's':
			a.state.CycleSortColumn()
			a.view.Refresh()
			return nil
		case 'S':
			a.state.ToggleSortOrder()
			a.view.Refresh()
			return nil
		case 'r':
			go a.loadData() // Refresh data in the background
			return nil
		case 'k':
			a.view.showKillConfirmationModal(false) // false for regular kill
			return nil
		case 'K':
			a.view.showKillConfirmationModal(true) // true for force kill
			return nil
		case 'p':
			a.view.showPingModal()
			return nil
		case 'h':
			a.view.showHelpModal()
			return nil
		case 'e':
			a.handleExport()
			return nil
		case '/':
			a.tviewApp.SetFocus(a.view.filterInput)
			return nil
		}
		return event
	})

	// Handler for when the selection in the table changes
	a.view.table.SetSelectionChangedFunc(func(row, column int) {
		a.view.updateDetailsView(row)
	})

	// Handler for when the user presses Enter on a table row
	a.view.table.SetSelectedFunc(func(row, column int) {
		if conn := a.view.GetSelectedConnection(); conn != nil {
			a.view.showDetailsModal(*conn)
		}
	})

	// Handlers for the filter input field
	a.view.filterInput.SetChangedFunc(func(text string) {
		a.state.SetFilterText(text)
		a.view.Refresh()
	})

	a.view.filterInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter || key == tcell.KeyEscape {
			a.tviewApp.SetFocus(a.view.table)
		}
	})
}

// handleExport triggers the CSV export action and shows a status message.
func (a *App) handleExport() {
	conns := a.state.GetFilteredConnections()
	if len(conns) == 0 {
		a.view.SetStatusMessage("No connections to export.")
		return
	}

	filename, err := actions.ExportToCSV(conns)
	if err != nil {
		a.view.SetStatusMessage("Error exporting to CSV: " + err.Error())
	} else {
		// Use backticks for the filename to make it stand out
		a.view.SetStatusMessage("Exported " + strings.TrimSpace(filename) + "")
	}
}
