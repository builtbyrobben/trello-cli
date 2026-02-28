package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/trello-cli/internal/outfmt"
)

type BoardsCmd struct {
	List BoardsListCmd `cmd:"" help:"List all boards"`
	Get  BoardsGetCmd  `cmd:"" help:"Get a board by ID"`
}

type BoardsListCmd struct{}

func (cmd *BoardsListCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	boards, err := client.Boards().List(ctx)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, boards)
	}
	if outfmt.IsPlain(ctx) {
		headers := []string{"ID", "NAME", "CLOSED"}
		var rows [][]string
		for _, b := range boards {
			rows = append(rows, []string{b.ID, b.Name, fmt.Sprintf("%v", b.Closed)})
		}
		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(boards) == 0 {
		fmt.Fprintln(os.Stderr, "No boards found")
		return nil
	}

	for _, b := range boards {
		fmt.Printf("%-24s  %s\n", b.ID, b.Name)
	}

	return nil
}

type BoardsGetCmd struct {
	ID string `arg:"" required:"" help:"Board ID"`
}

func (cmd *BoardsGetCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	board, err := client.Boards().Get(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, board)
	}
	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout,
			[]string{"ID", "NAME", "URL", "CLOSED"},
			[][]string{{board.ID, board.Name, board.URL, fmt.Sprintf("%v", board.Closed)}},
		)
	}

	fmt.Printf("ID:   %s\n", board.ID)
	fmt.Printf("Name: %s\n", board.Name)

	if board.Desc != "" {
		fmt.Printf("Desc: %s\n", board.Desc)
	}

	if board.URL != "" {
		fmt.Printf("URL:  %s\n", board.URL)
	}

	fmt.Printf("Closed: %v\n", board.Closed)

	return nil
}
