package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/builtbyrobben/trello-cli/internal/outfmt"
	"github.com/builtbyrobben/trello-cli/internal/secrets"
)

type AuthCmd struct {
	SetKey   AuthSetKeyCmd   `cmd:"" name:"set-key" help:"Set API key (uses --stdin by default)"`
	SetToken AuthSetTokenCmd `cmd:"" name:"set-token" help:"Set API token (uses --stdin by default)"`
	Status   AuthStatusCmd   `cmd:"" help:"Show authentication status"`
	Remove   AuthRemoveCmd   `cmd:"" help:"Remove stored credentials"`
}

type AuthSetKeyCmd struct {
	Stdin bool   `help:"Read API key from stdin (default: true)" default:"true"`
	Key   string `arg:"" optional:"" help:"API key (discouraged; exposes in shell history)"`
}

func (cmd *AuthSetKeyCmd) Run(ctx context.Context) error {
	apiKey, err := readCredential(cmd.Key, "API key")
	if err != nil {
		return err
	}

	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	if err := store.SetAPIKey(apiKey); err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "API key stored in keyring",
		})
	}

	fmt.Fprintln(os.Stderr, "API key stored in keyring")

	return nil
}

type AuthSetTokenCmd struct {
	Stdin bool   `help:"Read token from stdin (default: true)" default:"true"`
	Token string `arg:"" optional:"" help:"API token (discouraged; exposes in shell history)"`
}

func (cmd *AuthSetTokenCmd) Run(ctx context.Context) error {
	token, err := readCredential(cmd.Token, "API token")
	if err != nil {
		return err
	}

	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	if err := store.SetToken(token); err != nil {
		return fmt.Errorf("store token: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "API token stored in keyring",
		})
	}

	fmt.Fprintln(os.Stderr, "API token stored in keyring")

	return nil
}

type AuthStatusCmd struct{}

func (cmd *AuthStatusCmd) Run(ctx context.Context) error {
	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	hasKey, err := store.HasKey()
	if err != nil {
		return fmt.Errorf("check API key: %w", err)
	}

	hasToken, err := store.HasToken()
	if err != nil {
		return fmt.Errorf("check token: %w", err)
	}

	envKeyOverride := os.Getenv("TRELLO_API_KEY") != ""
	envTokenOverride := os.Getenv("TRELLO_TOKEN") != ""

	status := map[string]any{
		"has_key":            hasKey,
		"has_token":          hasToken,
		"env_key_override":   envKeyOverride,
		"env_token_override": envTokenOverride,
		"storage_backend":    "keyring",
	}

	if hasKey && !envKeyOverride {
		key, keyErr := store.GetAPIKey()
		if keyErr == nil && len(key) > 8 {
			status["key_redacted"] = key[:4] + "..." + key[len(key)-4:]
		}
	}

	if hasToken && !envTokenOverride {
		token, tokenErr := store.GetToken()
		if tokenErr == nil && len(token) > 8 {
			status["token_redacted"] = token[:4] + "..." + token[len(token)-4:]
		}
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, status)
	}

	fmt.Fprintf(os.Stderr, "Storage: %s\n", status["storage_backend"])

	// API Key status
	switch {
	case envKeyOverride:
		fmt.Fprintln(os.Stderr, "API Key: Using TRELLO_API_KEY environment variable")
	case hasKey:
		fmt.Fprintln(os.Stderr, "API Key: Configured")

		if redacted, ok := status["key_redacted"].(string); ok {
			fmt.Fprintf(os.Stderr, "  Key: %s\n", redacted)
		}
	default:
		fmt.Fprintln(os.Stderr, "API Key: Not configured")
	}

	// Token status
	switch {
	case envTokenOverride:
		fmt.Fprintln(os.Stderr, "Token: Using TRELLO_TOKEN environment variable")
	case hasToken:
		fmt.Fprintln(os.Stderr, "Token: Configured")

		if redacted, ok := status["token_redacted"].(string); ok {
			fmt.Fprintf(os.Stderr, "  Token: %s\n", redacted)
		}
	default:
		fmt.Fprintln(os.Stderr, "Token: Not configured")
	}

	if (!hasKey && !envKeyOverride) || (!hasToken && !envTokenOverride) {
		fmt.Fprintln(os.Stderr, "\nTo authenticate:")

		if !hasKey && !envKeyOverride {
			fmt.Fprintln(os.Stderr, "  trello-cli auth set-key --stdin")
		}

		if !hasToken && !envTokenOverride {
			fmt.Fprintln(os.Stderr, "  trello-cli auth set-token --stdin")
		}
	}

	return nil
}

type AuthRemoveCmd struct{}

func (cmd *AuthRemoveCmd) Run(ctx context.Context) error {
	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	if err := store.DeleteAPIKey(); err != nil {
		return fmt.Errorf("remove API key: %w", err)
	}

	if err := store.DeleteToken(); err != nil {
		return fmt.Errorf("remove token: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "All credentials removed",
		})
	}

	fmt.Fprintln(os.Stderr, "All credentials removed")

	return nil
}

func readCredential(argValue, label string) (string, error) {
	if argValue != "" {
		fmt.Fprintf(os.Stderr, "Warning: passing credentials as arguments exposes them in shell history. Use --stdin instead.\n")
		return strings.TrimSpace(argValue), nil
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintf(os.Stderr, "Enter %s: ", label)

		byteVal, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)

		if err != nil {
			return "", fmt.Errorf("read %s: %w", label, err)
		}

		return strings.TrimSpace(string(byteVal)), nil
	}

	byteVal, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("read %s from stdin: %w", label, err)
	}

	val := strings.TrimSpace(string(byteVal))
	if val == "" {
		return "", fmt.Errorf("%s cannot be empty", label)
	}

	return val, nil
}
