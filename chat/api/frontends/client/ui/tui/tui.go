package tui

import (
	"fmt"

	"github.com/ardanlabs/usdl/chat/api/frontends/client/app"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App interface {
	SendMessageHandler(to common.Address, msg []byte) error
}

type Storage interface {
	Contacts() []app.User
	QueryContactByID(id common.Address) (app.User, error)
}

// =============================================================================

type TUI struct {
	tviewApp *tview.Application
	flex     *tview.Flex
	list     *tview.List
	textView *tview.TextView
	textArea *tview.TextArea
	button   *tview.Button
	app      App
	db       Storage
}

func New(myAccountID common.Address, db Storage) *TUI {
	var ui TUI

	app := tview.NewApplication()

	// -------------------------------------------------------------------------

	textView := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetBorder(true)
	textView.SetTitle(fmt.Sprintf("*** %s ***", myAccountID))

	// -------------------------------------------------------------------------

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Users")
	list.SetChangedFunc(func(idx int, name string, id string, shortcut rune) {
		textView.Clear()

		if ui.app == nil {
			return
		}

		addrID := common.HexToAddress(id)

		user, err := ui.db.QueryContactByID(addrID)
		if err != nil {
			textView.ScrollToEnd()
			fmt.Fprintln(textView, "-----")
			fmt.Fprintln(textView, err.Error())
			return
		}

		for i, msg := range user.Messages {
			fmt.Fprintln(textView, string(msg))
			if i < len(user.Messages)-1 {
				fmt.Fprintln(textView, "-----")
			}
		}

		list.SetItemText(idx, user.Name, user.ID.Hex())
	})

	for i, user := range db.Contacts() {
		shortcut := rune(i + 49)
		list.AddItem(user.Name, user.ID.Hex(), shortcut, nil)
	}

	// -------------------------------------------------------------------------

	button := tview.NewButton("SUBMIT")
	button.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen).Bold(true))
	button.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen).Bold(true))
	button.SetBorder(true)
	button.SetBorderColor(tcell.ColorGreen)

	// -------------------------------------------------------------------------

	textArea := tview.NewTextArea()
	textArea.SetWrap(false)
	textArea.SetPlaceholder("Enter message here...")
	textArea.SetBorder(true)
	textArea.SetBorderPadding(0, 0, 1, 0)

	// -------------------------------------------------------------------------

	flex := tview.NewFlex().
		AddItem(list, 20, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(textView, 0, 5, false).
			AddItem(tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(textArea, 0, 90, false).
				AddItem(button, 0, 10, false),
				0, 1, false),
			0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlQ:
			app.Stop()
			return nil
		}

		return event
	})

	ui.tviewApp = app
	ui.flex = flex
	ui.list = list
	ui.textView = textView
	ui.textArea = textArea
	ui.button = button
	ui.db = db

	button.SetSelectedFunc(ui.buttonHandler)

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			ui.buttonHandler()
			return nil
		}
		return event
	})

	return &ui
}

func (ui *TUI) SetApp(app App) {
	ui.app = app

	// We need to get the change function to execute, so remove the last
	// items and add it back.

	lastItem := ui.list.GetItemCount() - 1
	main, secondary := ui.list.GetItemText(lastItem)
	ui.list.RemoveItem(lastItem)
	shortcut := rune(ui.list.GetItemCount() + 49)
	ui.list.AddItem(main, secondary, shortcut, nil)
	ui.list.SetCurrentItem(lastItem)
	ui.list.SetCurrentItem(0)
}

func (ui *TUI) Run() error {
	return ui.tviewApp.SetRoot(ui.flex, true).EnableMouse(true).Run()
}

func (ui *TUI) WriteText(id string, msg []byte) {
	ui.textView.ScrollToEnd()

	switch id {
	case "system":
		fmt.Fprintln(ui.textView, "-----")
		fmt.Fprintln(ui.textView, string(msg))

	default:
		idx := ui.list.GetCurrentItem()

		_, currentID := ui.list.GetItemText(idx)
		if currentID == "" {
			fmt.Fprintln(ui.textView, "-----")
			fmt.Fprintln(ui.textView, "id not found: "+id)
			return
		}

		if id == currentID {
			fmt.Fprintln(ui.textView, "-----")
			fmt.Fprintln(ui.textView, string(msg))
			return
		}

		for i := range ui.list.GetItemCount() {
			name, idStr := ui.list.GetItemText(i)
			if id == idStr {
				ui.list.SetItemText(i, "* "+name, idStr)
				ui.tviewApp.Draw()
				return
			}
		}
	}
}

func (ui *TUI) UpdateContact(id string, name string) {
	shortcut := rune(ui.list.GetItemCount() + 49)
	ui.list.AddItem(name, id, shortcut, nil)
}

// =============================================================================

func (ui *TUI) buttonHandler() {
	_, to := ui.list.GetItemText(ui.list.GetCurrentItem())

	msg := ui.textArea.GetText()
	if msg == "" {
		return
	}

	id := common.HexToAddress(to)

	if err := ui.app.SendMessageHandler(id, []byte(msg)); err != nil {
		ui.WriteText("system", fmt.Appendf(nil, "Error sending message: %s", err))
		return
	}

	ui.textArea.SetText("", false)
}
