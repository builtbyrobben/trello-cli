package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/trello-cli/internal/outfmt"
)

type ListsCmd struct {
	List   ListsListCmd   `cmd:"" help:"List all lists on a board"`
	Create ListsCreateCmd `cmd:"" help:"Create a new list on a board"`
}

type ListsListCmd struct {
	Board string `required:"" help:"Board ID"`
}

func (cmd *ListsListCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	lists, err := client.Lists().List(ctx, cmd.Board)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, lists)
	}
	if outfmt.IsPlain(ctx) {
		headers := []string{"ID", "NAME", "CLOSED"}
		var rows [][]string
		for _, l := range lists {
			rows = append(rows, []string{l.ID, l.Name, fmt.Sprintf("%v", l.Closed)})
		}
		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(lists) == 0 {
		fmt.Fprintln(os.Stderr, "No lists found")
		return nil
	}

	for _, l := range lists {
		fmt.Printf("%-24s  %s\n", l.ID, l.Name)
	}

	return nil
}

type ListsCreateCmd struct {
	Board string `required:"" help:"Board ID"`
	Name  string `required:"" help:"List name"`
}

func (cmd *ListsCreateCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	list, err := client.Lists().Create(ctx, cmd.Board, cmd.Name)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, list)
	}
	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout,
			[]string{"ID", "NAME"},
			[][]string{{list.ID, list.Name}},
		)
	}

	fmt.Fprintf(os.Stderr, "Created list\n\n")
	fmt.Printf("ID:   %s\n", list.ID)
	fmt.Printf("Name: %s\n", list.Name)

	return nil
}
