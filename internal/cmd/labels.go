package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/trello-cli/internal/outfmt"
)

type LabelsCmd struct {
	List LabelsListCmd `cmd:"" help:"List labels on a board"`
}

type LabelsListCmd struct {
	Board string `required:"" help:"Board ID"`
}

func (cmd *LabelsListCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	labels, err := client.Labels().ListByBoard(ctx, cmd.Board)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, labels)
	}

	if len(labels) == 0 {
		fmt.Fprintln(os.Stderr, "No labels found")
		return nil
	}

	for _, l := range labels {
		color := l.Color
		if color == "" {
			color = "(none)"
		}

		name := l.Name
		if name == "" {
			name = "(unnamed)"
		}

		fmt.Printf("%-24s  %-10s  %s\n", l.ID, color, name)
	}

	return nil
}
