package jira

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/cherryramatisdev/github"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/fs/dir"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
	"golang.design/x/clipboard"
)

const progressprefix string = "TEC"

//go:embed doc.md
var helpdoc string

var Cmd = &Z.Cmd{
	Name:        `jira`,
	Description: helpdoc,
	Commands:    []*Z.Cmd{help.Cmd, progressCmd, reviewCmd, viewCmd},
}

//go:embed progress.md
var progressdoc string

var progressCmd = &Z.Cmd{
	Name:        `progress`,
	Description: progressdoc,
	Commands:    []*Z.Cmd{help.Cmd},
	MinArgs:     1,
	MaxArgs:     2,
	Call: func(x *Z.Cmd, args ...string) error {
		ticketid := args[0]
		// TODO: Check this later
		// flag := args[1]
		prefix_and_ticket := fmt.Sprintf("%s-%s", progressprefix, ticketid)

		// if flag != "nojira" {
		if err := MoveTicketStatus(prefix_and_ticket, Transitions.Progress); err != nil {
			return err
		}

		if err := AutoAssignTicket(prefix_and_ticket); err != nil {
			return err
		}
		// }

		// if flag != "nogit" {
		if !dir.Exists("worktrees") {
			if err := dir.Create("worktrees"); err != nil {
				return err
			}
		}

		if dir.Exists("worktrees/" + ticketid) {
			fmt.Println("The branch exist locally at ./worktrees/" + ticketid)
			return nil
		}

		pr := github.Pr{
			Org:         github.GetCurrentOrg(),
			Repo:        github.GetCurrentRepo(),
			Prefixtoken: prefix_and_ticket,
		}

		if pr.Exist() {
			return Z.Exec("git", "worktree", "add", fmt.Sprintf("./worktrees/%s", ticketid), prefix_and_ticket)
		}

		originbranch := term.Prompt("From which branch you want to merge from? (main by default): ")
		if originbranch == "" {
			originbranch = "main"
		}

		return Z.Exec("git", "worktree", "add", fmt.Sprintf("./worktrees/%s", ticketid), "-b", fmt.Sprintf("feature/%s", prefix_and_ticket), originbranch)
		// }
		// return nil
	},
}

//go:embed review.md
var reviewdoc string

//go:embed prtemplate.md
var prtemplatedoc string

type templateData struct {
	PrUrl    string
	TicketId string
}

//go:embed view.md
var viewdoc string

var viewCmd = &Z.Cmd{
	Name:        `view`,
	Description: viewdoc,
	MinArgs:     1,
	Call: func(x *Z.Cmd, args ...string) error {
		ticketid := args[0]
		prefix_and_ticket := fmt.Sprintf("%s-%s", progressprefix, ticketid)

		ticket, err := GetTicket(prefix_and_ticket)

		if err != nil {
			return err
		}

		fmt.Print(ticket.Fields.Description)

		return nil
	},
}

var reviewCmd = &Z.Cmd{
	Name:        `review`,
	Usage:       `[help|<ticketid>]`,
	Description: reviewdoc,
	Commands:    []*Z.Cmd{help.Cmd},
	MinArgs:     1,
	Call: func(x *Z.Cmd, args ...string) error {
		ticketid := args[0]
		prefix_and_ticket := fmt.Sprintf("%s-%s", progressprefix, ticketid)

		if err := MoveTicketStatus(prefix_and_ticket, Transitions.Review); err != nil {
			return err
		}

		pr := github.Pr{
			Org:         github.GetCurrentOrg(),
			Repo:        github.GetCurrentRepo(),
			Prefixtoken: prefix_and_ticket,
		}

		url := pr.GetUrl()

		t, err := template.New("t").Parse(prtemplatedoc)
		if err != nil {
			return err
		}

		var data = templateData{
			TicketId: ticketid,
			PrUrl:    url,
		}

		var buf strings.Builder
		if err := t.Execute(&buf, data); err != nil {
			return err
		}

		if err := clipboard.Init(); err != nil {
			return err
		}

		clipboard.Write(clipboard.FmtText, []byte(buf.String()))

		return nil
	},
}
