package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/trello-cli/internal/outfmt"
	"github.com/builtbyrobben/trello-cli/internal/trello"
)

type CardsCmd struct {
	List   CardsListCmd   `cmd:"" help:"List cards on a board or in a list"`
	Get    CardsGetCmd    `cmd:"" help:"Get a card by ID"`
	Create CardsCreateCmd `cmd:"" help:"Create a new card"`
	Update CardsUpdateCmd `cmd:"" help:"Update a card"`
	Move   CardsMoveCmd   `cmd:"" help:"Move a card to another list"`
	Delete CardsDeleteCmd `cmd:"" help:"Delete a card"`
}

type CardsListCmd struct {
	Board string `required:"" help:"Board ID"`
	List  string `optional:"" help:"List ID (filter cards to a specific list)"`
}

func (cmd *CardsListCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	var cards []trello.Card

	if cmd.List != "" {
		cards, err = client.Cards().ListByList(ctx, cmd.List)
	} else {
		cards, err = client.Cards().ListByBoard(ctx, cmd.Board)
	}

	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, cards)
	}

	if len(cards) == 0 {
		fmt.Fprintln(os.Stderr, "No cards found")
		return nil
	}

	for _, c := range cards {
		fmt.Printf("%-24s  %s\n", c.ID, c.Name)
	}

	return nil
}

type CardsGetCmd struct {
	ID string `arg:"" required:"" help:"Card ID"`
}

func (cmd *CardsGetCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	card, err := client.Cards().Get(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, card)
	}

	fmt.Printf("ID:     %s\n", card.ID)
	fmt.Printf("Name:   %s\n", card.Name)

	if card.Desc != "" {
		fmt.Printf("Desc:   %s\n", card.Desc)
	}

	fmt.Printf("List:   %s\n", card.IDList)

	if card.URL != "" {
		fmt.Printf("URL:    %s\n", card.URL)
	}

	fmt.Printf("Closed: %v\n", card.Closed)

	return nil
}

type CardsCreateCmd struct {
	List string `required:"" help:"List ID"`
	Name string `required:"" help:"Card name"`
	Desc string `optional:"" help:"Card description"`
}

func (cmd *CardsCreateCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	card, err := client.Cards().Create(ctx, cmd.List, cmd.Name, cmd.Desc)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, card)
	}

	fmt.Fprintf(os.Stderr, "Created card\n\n")
	fmt.Printf("ID:   %s\n", card.ID)
	fmt.Printf("Name: %s\n", card.Name)

	if card.URL != "" {
		fmt.Printf("URL:  %s\n", card.URL)
	}

	return nil
}

type CardsUpdateCmd struct {
	ID   string  `arg:"" required:"" help:"Card ID"`
	Name *string `optional:"" help:"New card name"`
	Desc *string `optional:"" help:"New card description"`
}

func (cmd *CardsUpdateCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	card, err := client.Cards().Update(ctx, cmd.ID, cmd.Name, cmd.Desc)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, card)
	}

	fmt.Fprintf(os.Stderr, "Updated card\n\n")
	fmt.Printf("ID:   %s\n", card.ID)
	fmt.Printf("Name: %s\n", card.Name)

	return nil
}

type CardsMoveCmd struct {
	ID   string `arg:"" required:"" help:"Card ID"`
	List string `required:"" help:"Target list ID"`
}

func (cmd *CardsMoveCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	card, err := client.Cards().Move(ctx, cmd.ID, cmd.List)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, card)
	}

	fmt.Fprintf(os.Stderr, "Moved card to list %s\n\n", cmd.List)
	fmt.Printf("ID:   %s\n", card.ID)
	fmt.Printf("Name: %s\n", card.Name)

	return nil
}

type CardsDeleteCmd struct {
	ID string `arg:"" required:"" help:"Card ID"`
}

func (cmd *CardsDeleteCmd) Run(ctx context.Context) error {
	client, err := getTrelloClient()
	if err != nil {
		return err
	}

	if err := client.Cards().Delete(ctx, cmd.ID); err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Card deleted",
		})
	}

	fmt.Fprintln(os.Stderr, "Card deleted")

	return nil
}
