package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sureshkrishnan-v/lazynet/internal/actions"
	"github.com/sureshkrishnan-v/lazynet/internal/models"
)

// showHelpModal displays a colorful, detailed list of keybindings.
func (v *View) showHelpModal() {
	var builder strings.Builder
	builder.WriteString("[::b][yellow]LazyNet Keybindings[-:-:-]\n\n")
	builder.WriteString("[::u]Navigation[-:-]\n")
	builder.WriteString("[green]↑/↓      [white]Move selection up/down\n")
	builder.WriteString("[green]←/→      [white]Scroll table left/right\n")
	builder.WriteString("[green]Enter    [white]Show detailed info for selection\n\n")
	builder.WriteString("[::u]Actions[-:-]\n")
	builder.WriteString("[green]k        [white]Kill selected process (Graceful)\n")
	builder.WriteString("[green]K        [white]Force Kill selected process (SIGKILL)\n")
	builder.WriteString("[green]p        [white]Ping remote address of selection\n")
	builder.WriteString("[green]/        [white]Filter connections\n")
	builder.WriteString("[green]e        [white]Export visible connections to CSV\n\n")
	builder.WriteString("[::u]Sorting[-:-]\n")
	builder.WriteString("[green]s        [white]Cycle through sortable columns\n")
	builder.WriteString("[green]S        [white]Toggle sort order (ASC/DESC)\n\n")
	builder.WriteString("[::u]Application[-:-]\n")
	builder.WriteString("[green]h        [white]Show/Hide this help panel\n")
	builder.WriteString("[green]r        [white]Refresh connections manually\n")
	builder.WriteString("[green]q        [white]Quit LazyNet\n")

	textView := tview.NewTextView().SetDynamicColors(true).SetText(builder.String())
	textView.SetBorder(true).SetBorderPadding(1, 1, 1, 1)

	frame := tview.NewFrame(textView).
		AddText("Help", true, tview.AlignCenter, tview.Styles.TitleColor).
		AddText("Press 'h' or Esc to close", false, tview.AlignCenter, tview.Styles.SecondaryTextColor)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'h' {
			v.pages.RemovePage("help_modal")
			v.app.tviewApp.SetFocus(v.table)
			return nil
		}
		return event
	})

	v.pages.AddPage("help_modal", frame, true, true)
}

// showPingModal displays a modal and runs the ping command for the selected connection.
func (v *View) showPingModal() {
	c := v.GetSelectedConnection()
	if c == nil || c.Raddr == "" || strings.HasPrefix(c.Raddr, ":") {
		return // No remote address to ping
	}
	ip := strings.Split(c.Raddr, ":")[0]

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() { v.app.tviewApp.Draw() }) // Redraw on new content

	frame := tview.NewFrame(textView).
		AddText(fmt.Sprintf("Pinging %s...", ip), true, tview.AlignCenter, tview.Styles.TitleColor).
		AddText("Press Esc to close", false, tview.AlignCenter, tview.Styles.SecondaryTextColor)

	ctx, cancel := context.WithCancel(context.Background())

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			cancel() // Cancel the context to stop the ping command
			v.pages.RemovePage("ping_modal")
			v.app.tviewApp.SetFocus(v.table)
			return nil
		}
		return event
	})

	outputChan := make(chan string)
	go actions.Ping(ctx, ip, outputChan)

	go func() {
		for line := range outputChan {
			v.app.tviewApp.QueueUpdateDraw(func() {
				currentText := textView.GetText(false)
				textView.SetText(currentText + line + "\n")
				textView.ScrollToEnd()
			})
		}
	}()

	v.pages.AddPage("ping_modal", frame, true, true)
}

// showKillConfirmationModal displays a modal to confirm killing a process.
func (v *View) showKillConfirmationModal(force bool) {
	c := v.GetSelectedConnection()
	if c == nil || c.Pid == 0 {
		return
	}

	actionText := "kill"
	actionFunc := func() error { return actions.KillProcess(c.Pid) }
	if force {
		actionText = "forcefully kill"
		actionFunc = func() error { return actions.ForceKillProcess(c.Pid) }
	}

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to %s process '%s' (PID: %d)?", actionText, c.ProcessName, c.Pid)).
		AddButtons([]string{"Confirm", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Confirm" {
				if err := actionFunc(); err == nil {
					go v.app.loadData()
				}
			}
			v.pages.RemovePage("kill_confirm").ShowPage("main")
			v.app.tviewApp.SetFocus(v.table)
		})
	v.pages.AddPage("kill_confirm", modal, true, true)
}

// showDetailsModal creates and displays a modal with in-depth connection details.
func (v *View) showDetailsModal(c models.Connection) {
	details := models.GetDetailedInfo(c.Pid)
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[yellow]Process:[white]    %s\n", c.ProcessName))
	builder.WriteString(fmt.Sprintf("[yellow]PID:[white]        %d\n", c.Pid))
	builder.WriteString(fmt.Sprintf("[yellow]User:[white]       %s\n\n", details.Username))
	builder.WriteString(fmt.Sprintf("[yellow]Status:[white]     %s\n", c.Status))
	builder.WriteString(fmt.Sprintf("[yellow]Local Addr:[white] %s\n", c.Laddr))
	builder.WriteString(fmt.Sprintf("[yellow]Remote Addr:[white] %s\n\n", c.Raddr))
	builder.WriteString(fmt.Sprintf("[yellow]Command:[white]\n%s\n", details.Cmdline))

	textView := tview.NewTextView().SetDynamicColors(true).SetText(builder.String())
	textView.SetBorder(true).SetBorderPadding(1, 1, 1, 1)

	frame := tview.NewFrame(textView).
		AddText("Connection Details", true, tview.AlignCenter, tview.Styles.TitleColor).
		AddText("Press Enter or Esc to close", false, tview.AlignCenter, tview.Styles.SecondaryTextColor)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyEscape {
			v.pages.RemovePage("details_modal")
			v.app.tviewApp.SetFocus(v.table)
			return nil
		}
		return event
	})

	v.pages.AddPage("details_modal", frame, true, true)
}
