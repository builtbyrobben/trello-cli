package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/trello-cli/internal/outfmt"
)

type MembersCmd struct {
	List MembersListCmd `cmd:"" help:"List members of a board"`
	Me   MembersMeCmd   `cmd:"" help:"Show current authenticated member"`
}

type MembersListCmd struct {
	Board string `required:"" help:"Board ID"`
}

func (cmd *MembersListCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	members, err := client.Members().ListByBoard(ctx, cmd.Board)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, members)
	}

	if len(members) == 0 {
		fmt.Fprintln(os.Stderr, "No members found")
		return nil
	}

	for _, m := range members {
		fmt.Printf("%-24s  @%-20s  %s\n", m.ID, m.Username, m.FullName)
	}

	return nil
}

type MembersMeCmd struct{}

func (cmd *MembersMeCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	member, err := client.Members().Me(ctx)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, member)
	}

	fmt.Printf("ID:       %s\n", member.ID)
	fmt.Printf("Username: %s\n", member.Username)

	if member.FullName != "" {
		fmt.Printf("Name:     %s\n", member.FullName)
	}

	if member.URL != "" {
		fmt.Printf("URL:      %s\n", member.URL)
	}

	return nil
}
