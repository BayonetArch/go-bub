package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	name   []string
	text   textinput.Model
	output string
	cmd    string
	err    error
}

func InitialModel() *model {
	ti := textinput.New()
	ti.Placeholder = "enter something"
	ti.Focus()
	ti.Prompt = "\033[1;32m$\033[0m "
	ti.CharLimit = 200
	ti.Width = 100

	return &model{
		name:   make([]string, 10),
		output: "",
		text:   ti,
		cmd:    "",
		err:    nil,
	}
}

func Time() tea.Cmd {

	return func() tea.Msg {

		t := time.Now()
		loc, _ := time.LoadLocation("Asia/Kathmandu")
		tf := t.In(loc).Format("2006 01 02 03:04:05PM")

		time.Sleep(time.Second * 3)
		return OUTPUT(fmt.Sprintf("\r%s", tf))
	}
}

func Cmd(command *string) tea.Cmd {

	return func() tea.Msg {
		if command == nil {
			return nil
		}

		out, err := exec.Command(
			"sh",
			"-c",
			fmt.Sprintf("%s", *command),
		).
			CombinedOutput()

		/*



		   1./system/bin/pm list packages  --user 0   | rg "$userinput"  | sed 's/package://'

		   2. /system/bin/pm resolve-activity --brief  --user 0 $fstep | sed ' 1d' | sed 's/ //'
		   3./system/bin/am  start --user 0  -n  $sstep
		*/

		if len(out) > 400 {
			return CMD("\noutput is too long")
		}

		if *command == "" {
			return ERROR(nil)
		}

		if err != nil {
			return ERROR(fmt.Errorf("%s", out))
		}
		return CMD(fmt.Sprintf("%s\n", out))
	}
}

type OUTPUT string
type CMD string
type ERROR error

func (m model) Init() tea.Cmd {

	return tea.Batch(
		Time(),
		Cmd(nil),
	)

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case OUTPUT:
		m.output = string(msg)
		return m, Time()

	case CMD:
		m.cmd = string(msg)
		m.err = nil
		return m, nil

	case ERROR:
		m.err = msg
		m.cmd = ""
		return m, nil
	case tea.KeyMsg:

		switch msg.String() {

		case "alt+e":
			return m, tea.Quit

		case "enter", "alt+m":
			var val string
			val = m.text.Value()
			m.text.Reset()

			return m, Cmd(&val)
		}
	}
	var cmd tea.Cmd
	m.text, cmd = m.text.Update(msg)

	return m, cmd
}

func (m model) View() string {

	err := ""

	if m.err != nil {
		styles := lipgloss.NewStyle()
		styles = styles.Underline(true).UnderlineSpaces(false).Foreground(lipgloss.Color("9")).Width(40)

		err += styles.Render(fmt.Sprintf("\n%s", m.err))
	} else {
		err = ""
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		m.output,
		"\n"+m.text.View(),
		m.cmd,
		err,
	)
}

func main() {

	f, e := tea.LogToFile("debug.log", "[debug]:")

	if e != nil {
		log.Fatal(e)
	}

	defer f.Close()
	_, err := tea.NewProgram(InitialModel(), tea.WithAltScreen()).Run()
	if err != nil {
		log.Fatal(err)
	}
}
