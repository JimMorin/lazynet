package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sureshkrishnan-v/lazynet/internal/models"
)

// View represents the main UI of the application. It holds all the tview components.
type View struct {
	app         *App
	table       *tview.Table
	detailsView *tview.TextView
	filterInput *tview.InputField
	hintView    *tview.TextView
	pages       *tview.Pages
}

func NewView(app *App) *View {
	v := &View{app: app}

	v.table = tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false).
		SetFixed(1, 0)

	details := tview.NewTextView()
	details.SetDynamicColors(true)
	details.SetWordWrap(true)
	details.SetBorder(true)
	details.SetTitle(" Details ")
	v.detailsView = details

	v.filterInput = tview.NewInputField().
		SetLabel("Filter: ").
		SetLabelColor(tcell.ColorYellow).
		SetFieldBackgroundColor(tcell.ColorDarkSlateGray)

	hint := tview.NewTextView()
	hint.SetDynamicColors(true)
	hint.SetText("[::b]Keys:[-:-] [yellow]/[white]Filter [yellow]s/S[white]Sort [yellow]k/K[white]Kill [yellow]p[white]Ping [yellow]e[white]Export [yellow]h[white]Help [yellow]q[white]Quit")
	v.hintView = hint

	v.pages = tview.NewPages()

	return v
}

// Init sets up the layout and makes the view visible.
func (v *View) Init() {
	mainFlex := tview.NewFlex().
		AddItem(v.table, 0, 1, true).
		AddItem(v.detailsView, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, true).
		AddItem(v.filterInput, 1, 0, false).
		AddItem(v.hintView, 1, 0, false)

	v.pages.AddPage("main", layout, true, true)
	v.app.tviewApp.SetRoot(v.pages, true).EnableMouse(true)
}

// Refresh re-draws the main table and details view with the latest data from the app state.
func (v *View) Refresh() {
	v.populateTable()
	v.updateHeaderIndicator()
	selectedRow, _ := v.table.GetSelection()
	v.updateDetailsView(selectedRow)
}

// GetSelectedConnection returns the connection object for the currently selected row in the table.
func (v *View) GetSelectedConnection() *models.Connection {
	row, _ := v.table.GetSelection()
	if row < 1 {
		return nil
	}
	conns := v.app.state.GetFilteredConnections()
	if row > len(conns) {
		return nil
	}
	return &conns[row-1]
}

// ShowInfoModal displays a short-lived informational message to the user.
func (v *View) ShowInfoModal(message string, duration int) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"})

	// Automatically close the modal after a duration
	go func() {
		<-time.After(time.Duration(duration) * time.Second) // corrected
		v.app.tviewApp.QueueUpdateDraw(func() {
			v.pages.RemovePage("info_modal")
		})
	}()

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		v.pages.RemovePage("info_modal")
	})

	v.pages.AddPage("info_modal", modal, true, true)
}

// SetStatusMessage updates the hint view with a temporary message.
func (v *View) SetStatusMessage(message string) {
	originalText := v.hintView.GetText(false)
	v.hintView.SetText(fmt.Sprintf("[yellow]Status: [white]%s", message))

	go func() {
		<-time.After(3 * time.Second) // corrected
		v.app.tviewApp.QueueUpdateDraw(func() {
			v.hintView.SetText(originalText)
		})
	}()
}
