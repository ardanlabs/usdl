// Package ui provides the TUI for the application.
package ui

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"regexp"
	"strings"
	"time"

	"github.com/ardanlabs/usdl/api/clients/tui/ui/client"
	"github.com/ardanlabs/usdl/foundation/agents/ollamallm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// =============================================================================

type TUI struct {
	tviewApp *tview.Application
	flex     *tview.Flex
	list     *tview.List
	textView *tview.TextView
	textArea *tview.TextArea
	aiToggle *tview.Button
	button   *tview.Button
	app      *client.App
	agent    *ollamallm.Agent
	history  *history
	aiMode   bool
}

func New(myAccountID common.Address, agent *ollamallm.Agent) *TUI {
	ui := TUI{
		agent:   agent,
		history: NewHistory(5),
	}

	tApp := tview.NewApplication()

	// -------------------------------------------------------------------------

	textView := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetWordWrap(true).
		SetChangedFunc(func() {
			tApp.Draw()
		})

	textView.SetBorder(true)
	textView.SetTitle(fmt.Sprintf("*** %s ***", myAccountID))

	// -------------------------------------------------------------------------

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Users")

	list.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftDoubleClick {
			ui.establishUserConnection()
		}

		return action, event
	})

	list.SetChangedFunc(func(idx int, name string, id string, shortcut rune) {
		textView.Clear()

		if ui.app == nil {
			return
		}

		actID := id
		if i := strings.Index(id, "]"); i > 0 {
			actID = id[i+1:]
		}

		addrID := common.HexToAddress(actID)

		user, err := ui.app.QueryContactByID(addrID)
		if err != nil {
			textView.ScrollToEnd()
			fmt.Fprintln(textView, "-----")
			fmt.Fprintln(textView, err.Error()+":"+addrID.Hex())
			return
		}

		for i, msg := range user.Messages {
			fmt.Fprintf(textView, "%s: %s\n", msg.Name, client.StitchMessages(msg.Content))
			if i < len(user.Messages)-1 {
				fmt.Fprintln(textView, "-----")
			}
		}

		textView.ScrollToEnd()

		name = strings.ReplaceAll(name, "* ", "")
		ui.list.SetItemText(idx, name, id)
	})

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

	aiToggle := tview.NewButton("Agent Off")
	aiToggle.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed).Bold(true))
	aiToggle.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed).Bold(true))
	aiToggle.SetBorder(true)
	aiToggle.SetBorderColor(tcell.ColorRed)

	if agent == nil {
		aiToggle.SetLabel("Agent Disabled")
	}

	// -------------------------------------------------------------------------

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(list, 0, 85, false).
			AddItem(aiToggle, 0, 15, false),
			20, 1, false).
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
			tApp.Stop()
			return nil
		}

		return event
	})

	ui.tviewApp = tApp
	ui.flex = flex
	ui.list = list
	ui.textView = textView
	ui.textArea = textArea
	ui.aiToggle = aiToggle
	ui.button = button

	// -------------------------------------------------------------------------

	aiToggleHandler := func() {
		ui.aiToggleHandler(agent != nil)
	}

	aiToggle.SetSelectedFunc(aiToggleHandler)

	buttonHandler := func() {
		ui.buttonHandler(common.Address{})
	}

	button.SetSelectedFunc(buttonHandler)

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			buttonHandler()
			return nil
		}
		return event
	})

	return &ui
}

func (ui *TUI) SetApp(app *client.App) {
	ui.app = app

	for i, user := range app.Contacts() {
		shortcut := rune(i + 49)
		switch user.Key {
		case "":
			ui.list.AddItem(user.Name, "[red]"+user.ID.Hex(), shortcut, nil)
		default:
			ui.list.AddItem(user.Name, "[green]"+user.ID.Hex(), shortcut, nil)
		}
	}
}

func (ui *TUI) Run() error {
	ui.updateState()

	return ui.tviewApp.SetRoot(ui.flex, true).EnableMouse(true).Run()
}

func (ui *TUI) WriteText(msg client.Message) {
	ui.textView.ScrollToEnd()

	switch msg.From {
	case common.Address{}:
		fmt.Fprintln(ui.textView, "-----")
		fmt.Fprintf(ui.textView, "%s: %s\n", msg.Name, client.StitchMessages(msg.Content))

	case ui.app.ID():
		idx := ui.list.GetCurrentItem()
		_, currentID := ui.GetItemText(idx)

		if msg.To.Hex() == currentID {
			fmt.Fprintln(ui.textView, "-----")
			fmt.Fprintf(ui.textView, "%s: %s\n", msg.Name, client.StitchMessages(msg.Content))
		}

	default:
		idx := ui.list.GetCurrentItem()

		_, currentID := ui.GetItemText(idx)
		if currentID == "" {
			fmt.Fprintln(ui.textView, "-----")
			fmt.Fprintln(ui.textView, "id not found: "+msg.From.Hex())
			return
		}

		msgContent := fmt.Sprintf("%s: %s", msg.Name, client.StitchMessages(msg.Content))

		ui.history.add(msg.From, msgContent)

		if msg.From.Hex() == currentID {
			fmt.Fprintln(ui.textView, "-----")
			fmt.Fprintf(ui.textView, "%s\n", msgContent)

			if ui.aiMode {
				ui.agentResponse(msg.From)
			}

			return
		}

		for i := range ui.list.GetItemCount() {
			name, idStr := ui.GetItemText(i)
			if msg.From.Hex() == idStr {
				if !strings.Contains(name, "*") {
					ui.list.SetItemText(i, "* "+name, idStr)
					ui.tviewApp.Draw()
				}

				if ui.aiMode {
					ui.agentResponse(msg.From)
				}

				return
			}
		}
	}
}

func (ui *TUI) AddContact(id common.Address, name string) {
	shortcut := rune(ui.list.GetItemCount() + 49)
	ui.list.AddItem(name, id.Hex(), shortcut, nil)
}

var re = regexp.MustCompile(`\s{2,}`)

func (ui *TUI) ApplyContactPrefix(id common.Address, option string, add bool) {
	for i := range ui.list.GetItemCount() {
		name, idStr := ui.GetItemText(i)

		if id.Hex() == idStr {
			hasStar := strings.Contains(name, "*")
			hasLeftArrow := strings.Contains(name, "<-")
			hasRightArrow := strings.Contains(name, "->")

			name = strings.ReplaceAll(name, "*", "")
			name = strings.ReplaceAll(name, ">", "")
			name = strings.ReplaceAll(name, "<", "")
			name = strings.ReplaceAll(name, "-", "")

			switch add {
			case true:
				switch option {
				case "->":
					finalName := fmt.Sprintf("-> %s", name)
					if hasLeftArrow {
						finalName = fmt.Sprintf("<- %s", finalName)
					}
					if hasStar {
						finalName = fmt.Sprintf("* %s", finalName)
					}
					finalName = strings.ReplaceAll(finalName, "<- ->", "<->")
					finalName = strings.TrimSpace(re.ReplaceAllString(finalName, " "))
					ui.list.SetItemText(i, finalName, idStr)

				case "<-":
					finalName := name
					if hasRightArrow {
						finalName = fmt.Sprintf("-> %s", finalName)
					}
					finalName = fmt.Sprintf("<- %s", finalName)
					if hasStar {
						finalName = fmt.Sprintf("* %s", finalName)
					}
					finalName = strings.ReplaceAll(finalName, "<- ->", "<->")
					finalName = strings.TrimSpace(re.ReplaceAllString(finalName, " "))
					ui.list.SetItemText(i, finalName, idStr)
				}

			case false:
				switch option {
				case "->":
					finalName := name
					if hasLeftArrow {
						finalName = fmt.Sprintf("<- %s", finalName)
					}
					if hasStar {
						finalName = fmt.Sprintf("* %s", finalName)
					}
					finalName = strings.TrimSpace(re.ReplaceAllString(finalName, " "))
					ui.list.SetItemText(i, finalName, idStr)

				case "<-":
					finalName := name
					if hasRightArrow {
						finalName = fmt.Sprintf("-> %s", finalName)
					}
					if hasStar {
						finalName = fmt.Sprintf("* %s", finalName)
					}
					finalName = strings.TrimSpace(re.ReplaceAllString(finalName, " "))
					ui.list.SetItemText(i, finalName, idStr)
				}
			}

			return
		}
	}
}

// =============================================================================

func (ui *TUI) updateState() {
	state, err := ui.app.GetState(context.Background())
	if err != nil {
		fmt.Fprintln(ui.textView, "-----")
		fmt.Fprintln(ui.textView, "failed to get state: "+err.Error())
	}

	for _, conn := range state.TCPConnections {
		ui.ApplyContactPrefix(conn, "->", true)
	}
}

func (ui *TUI) agentResponse(from common.Address) {
	ctx := context.TODO()

	msgs := ui.history.retrieve(from)

	input := msgs[len(msgs)-1]
	history := msgs[:len(msgs)-1]

	start := time.Now()

	resp, err := ui.agent.Chat(ctx, input, history)
	if err != nil {
		fmt.Fprintln(ui.textView, "-----")
		fmt.Fprintln(ui.textView, "failed ollama response: "+err.Error())
	}

	// Create some artificial delay to simulate thinking
	// We don't care about negative calculations because
	// that means we've already waited long enough.

	dur := int(math.Ceil(time.Since(start).Seconds()))
	delayInSeconds := max(rand.IntN(11), 3)
	delayInSeconds -= dur
	time.Sleep(time.Duration(delayInSeconds) * time.Second)

	ui.textArea.SetText(resp, true)
	ui.buttonHandler(from)
}

func (ui *TUI) buttonHandler(to common.Address) {
	if to == (common.Address{}) {
		_, id := ui.GetItemText(ui.list.GetCurrentItem())
		to = common.HexToAddress(id)
	}

	msg := ui.textArea.GetText()
	if msg == "" {
		return
	}

	// Found a issue where the AI started a response with a /.
	if ui.aiMode && msg[0] == '/' {
		msg = msg[1:]
	}

	if err := ui.app.SendMessageHandler(to, []byte(msg)); err != nil {
		msg := client.Message{
			Name:    "system",
			Content: [][]byte{fmt.Appendf(nil, "Error sending message: %s", err)},
		}
		ui.WriteText(msg)
		return
	}

	ui.history.add(to, msg)

	ui.textArea.SetText("", false)
}

func (ui *TUI) aiToggleHandler(agent bool) {
	if !agent {
		return
	}

	switch ui.aiToggle.GetLabel() {
	case "Agent Off":
		ui.aiToggle.SetLabel("Agent On")
		ui.aiToggle.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen).Bold(true))
		ui.aiToggle.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen).Bold(true))
		ui.aiToggle.SetBorderColor(tcell.ColorGreen)
		ui.aiMode = true

	case "Agent On":
		ui.aiToggle.SetLabel("Agent Off")
		ui.aiToggle.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed).Bold(true))
		ui.aiToggle.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed).Bold(true))
		ui.aiToggle.SetBorderColor(tcell.ColorRed)
		ui.aiMode = false
	}
}

func (ui *TUI) establishUserConnection() {
	idx := ui.list.GetCurrentItem()
	name, currentID := ui.GetItemText(idx)

	fmt.Fprintln(ui.textView, "-----")
	fmt.Fprintf(ui.textView, "Establishing Peer Connection with %s\n", name)

	if err := ui.app.EstablishTCPConnection(context.Background(), ui.app.ID(), common.HexToAddress(currentID)); err != nil {
		if errors.Is(err, client.ErrConnectionDropped) {
			fmt.Fprintln(ui.textView, "-----")
			fmt.Fprintf(ui.textView, "TCP connection dropped for %s\n", name)
			ui.ApplyContactPrefix(common.HexToAddress(currentID), "<-", false)
			return
		}

		fmt.Fprintln(ui.textView, "-----")
		fmt.Fprintf(ui.textView, "Failed to establish TCP connection: %s\n", err)
		ui.ApplyContactPrefix(common.HexToAddress(currentID), "<-", false)
		return
	}

	fmt.Fprintln(ui.textView, "-----")
	fmt.Fprintf(ui.textView, "TCP connection established with %s\n", name)

	ui.ApplyContactPrefix(common.HexToAddress(currentID), "<-", true)
}

func (ui *TUI) GetItemText(idx int) (string, string) {
	name, id := ui.list.GetItemText(idx)

	i := strings.Index(id, "]")
	if i == -1 {
		return name, id
	}

	return name, id[i+1:]
}
