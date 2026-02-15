package cmd

import (
	"fmt"
	"os"

	"github.com/builtbyrobben/trello-cli/internal/secrets"
	"github.com/builtbyrobben/trello-cli/internal/trello"
)

func getTrelloClient() (*trello.Client, error) {
	// Check for environment variable overrides first
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")

	if apiKey == "" || token == "" {
		store, err := secrets.OpenDefault()
		if err != nil {
			return nil, fmt.Errorf("open credential store: %w", err)
		}

		if apiKey == "" {
			apiKey, err = store.GetAPIKey()
			if err != nil {
				return nil, fmt.Errorf("get API key: %w (set TRELLO_API_KEY or run 'trello-cli auth set-key --stdin')", err)
			}
		}

		if token == "" {
			token, err = store.GetToken()
			if err != nil {
				return nil, fmt.Errorf("get token: %w (set TRELLO_TOKEN or run 'trello-cli auth set-token --stdin')", err)
			}
		}
	}

	return trello.NewClient(apiKey, token), nil
}
