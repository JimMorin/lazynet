package tui

import (
	"time"

	"github.com/rivo/tview"
	"github.com/sureshkrishnan-v/lazynet/internal/conn"
)

// App encapsulates the entire TUI application, including its state and components.
// It holds references to all UI views and manages the application's data state.
type App struct {
	tviewApp *tview.Application
	view     *View
	state    *AppState
}

// NewApp creates and initializes a new Lazynet application.
func NewApp() *App {
	a := &App{
		tviewApp: tview.NewApplication(),
		state:    NewAppState(),
	}
	a.view = NewView(a) // Pass the app instance to the view
	return a
}

// Run starts the TUI application and the main data refresh loop.
func (a *App) Run() error {
	a.view.Init()
	a.setKeybindings()

	// Start the background data refresh loop
	go a.refreshDataLoop()

	return a.tviewApp.Run()
}

// Stop gracefully shuts down the application.
func (a *App) Stop() {
	a.tviewApp.Stop()
}

// refreshDataLoop periodically fetches new connection data and redraws the UI.
// This runs in a separate goroutine to avoid blocking the UI.
func (a *App) refreshDataLoop() {
	// Initial data load
	a.loadData()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		a.loadData()
	}
}

// loadData fetches, sorts, and filters connections, then queues a UI update.
func (a *App) loadData() {
	// Preserve the currently selected row index to maintain selection after refresh
	selectedRow, _ := a.view.table.GetSelection()

	// Fetch fresh connection data
	conns, err := conn.FetchConnections()
	if err != nil {
		// In a real app, you might want to display this error in the UI
		return
	}

	// Update the application state with the new data
	a.state.SetConnections(conns)

	// Queue an update on the main UI thread
	a.tviewApp.QueueUpdateDraw(func() {
		a.view.Refresh()
		// Restore selection if it's within the new bounds of the table
		if selectedRow < a.view.table.GetRowCount() {
			a.view.table.Select(selectedRow, 0)
		} else if a.view.table.GetRowCount() > 1 {
			a.view.table.Select(1, 0) // Select the first data row if previous selection is out of bounds
		}
	})
}
