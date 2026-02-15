package trello

import (
	"context"
	"errors"
	"fmt"

	"github.com/builtbyrobben/trello-cli/internal/api"
)

var (
	errBoardIDRequired = errors.New("board ID is required")
	errListIDRequired  = errors.New("list ID is required")
	errCardIDRequired  = errors.New("card ID is required")
	errNameRequired    = errors.New("name is required")
)

const defaultBaseURL = "https://api.trello.com/1"

// Client wraps the API client with Trello-specific methods.
type Client struct {
	*api.Client
}

// NewClient creates a new Trello API client with query param auth.
func NewClient(apiKey, token string) *Client {
	return &Client{
		Client: api.NewClient("",
			api.WithBaseURL(defaultBaseURL),
			api.WithUserAgent("trello-cli/1.0"),
			api.WithQueryAuth(map[string]string{
				"key":   apiKey,
				"token": token,
			}),
		),
	}
}

// Boards provides methods for the Boards API.
func (c *Client) Boards() *BoardsService {
	return &BoardsService{client: c}
}

// Lists provides methods for the Lists API.
func (c *Client) Lists() *ListsService {
	return &ListsService{client: c}
}

// Cards provides methods for the Cards API.
func (c *Client) Cards() *CardsService {
	return &CardsService{client: c}
}

// Members provides methods for the Members API.
func (c *Client) Members() *MembersService {
	return &MembersService{client: c}
}

// Labels provides methods for the Labels API.
func (c *Client) Labels() *LabelsService {
	return &LabelsService{client: c}
}

// BoardsService handles board operations.
type BoardsService struct {
	client *Client
}

// List returns all boards for the authenticated member.
func (s *BoardsService) List(ctx context.Context) ([]Board, error) {
	var result []Board
	if err := s.client.Get(ctx, "/members/me/boards", &result); err != nil {
		return nil, fmt.Errorf("list boards: %w", err)
	}

	return result, nil
}

// Get returns a board by ID.
func (s *BoardsService) Get(ctx context.Context, id string) (*Board, error) {
	if id == "" {
		return nil, errBoardIDRequired
	}

	var result Board

	path := fmt.Sprintf("/boards/%s", id)
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get board: %w", err)
	}

	return &result, nil
}

// ListsService handles list operations.
type ListsService struct {
	client *Client
}

// List returns all lists on a board.
func (s *ListsService) List(ctx context.Context, boardID string) ([]List, error) {
	if boardID == "" {
		return nil, errBoardIDRequired
	}

	var result []List

	path := fmt.Sprintf("/boards/%s/lists", boardID)
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list lists: %w", err)
	}

	return result, nil
}

// Create creates a new list on a board.
func (s *ListsService) Create(ctx context.Context, boardID, name string) (*List, error) {
	if boardID == "" {
		return nil, errBoardIDRequired
	}

	if name == "" {
		return nil, errNameRequired
	}

	req := CreateListRequest{
		Name:    name,
		IDBoard: boardID,
	}

	var result List
	if err := s.client.Post(ctx, "/lists", req, &result); err != nil {
		return nil, fmt.Errorf("create list: %w", err)
	}

	return &result, nil
}

// CardsService handles card operations.
type CardsService struct {
	client *Client
}

// ListByBoard returns all cards on a board.
func (s *CardsService) ListByBoard(ctx context.Context, boardID string) ([]Card, error) {
	if boardID == "" {
		return nil, errBoardIDRequired
	}

	var result []Card

	path := fmt.Sprintf("/boards/%s/cards", boardID)
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list cards by board: %w", err)
	}

	return result, nil
}

// ListByList returns all cards in a list.
func (s *CardsService) ListByList(ctx context.Context, listID string) ([]Card, error) {
	if listID == "" {
		return nil, errListIDRequired
	}

	var result []Card

	path := fmt.Sprintf("/lists/%s/cards", listID)
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list cards by list: %w", err)
	}

	return result, nil
}

// Get returns a card by ID.
func (s *CardsService) Get(ctx context.Context, id string) (*Card, error) {
	if id == "" {
		return nil, errCardIDRequired
	}

	var result Card

	path := fmt.Sprintf("/cards/%s", id)
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}

	return &result, nil
}

// Create creates a new card.
func (s *CardsService) Create(ctx context.Context, listID, name, desc string) (*Card, error) {
	if listID == "" {
		return nil, errListIDRequired
	}

	if name == "" {
		return nil, errNameRequired
	}

	req := CreateCardRequest{
		Name:   name,
		IDList: listID,
		Desc:   desc,
	}

	var result Card
	if err := s.client.Post(ctx, "/cards", req, &result); err != nil {
		return nil, fmt.Errorf("create card: %w", err)
	}

	return &result, nil
}

// Update updates a card's name and/or description.
func (s *CardsService) Update(ctx context.Context, id string, name, desc *string) (*Card, error) {
	if id == "" {
		return nil, errCardIDRequired
	}

	req := make(map[string]string)

	if name != nil {
		req["name"] = *name
	}

	if desc != nil {
		req["desc"] = *desc
	}

	var result Card

	path := fmt.Sprintf("/cards/%s", id)
	if err := s.client.Put(ctx, path, req, &result); err != nil {
		return nil, fmt.Errorf("update card: %w", err)
	}

	return &result, nil
}

// Move moves a card to a different list.
func (s *CardsService) Move(ctx context.Context, cardID, listID string) (*Card, error) {
	if cardID == "" {
		return nil, errCardIDRequired
	}

	if listID == "" {
		return nil, errListIDRequired
	}

	var result Card

	path := fmt.Sprintf("/cards/%s", cardID)
	if err := s.client.Put(ctx, path, map[string]string{"idList": listID}, &result); err != nil {
		return nil, fmt.Errorf("move card: %w", err)
	}

	return &result, nil
}

// Delete deletes a card.
func (s *CardsService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errCardIDRequired
	}

	path := fmt.Sprintf("/cards/%s", id)
	if err := s.client.Delete(ctx, path); err != nil {
		return fmt.Errorf("delete card: %w", err)
	}

	return nil
}

// MembersService handles member operations.
type MembersService struct {
	client *Client
}

// Me returns the authenticated member.
func (s *MembersService) Me(ctx context.Context) (*Member, error) {
	var result Member
	if err := s.client.Get(ctx, "/members/me", &result); err != nil {
		return nil, fmt.Errorf("get current member: %w", err)
	}

	return &result, nil
}

// ListByBoard returns all members of a board.
func (s *MembersService) ListByBoard(ctx context.Context, boardID string) ([]Member, error) {
	if boardID == "" {
		return nil, errBoardIDRequired
	}

	var result []Member

	path := fmt.Sprintf("/boards/%s/members", boardID)
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list board members: %w", err)
	}

	return result, nil
}

// LabelsService handles label operations.
type LabelsService struct {
	client *Client
}

// ListByBoard returns all labels on a board.
func (s *LabelsService) ListByBoard(ctx context.Context, boardID string) ([]Label, error) {
	if boardID == "" {
		return nil, errBoardIDRequired
	}

	var result []Label

	path := fmt.Sprintf("/boards/%s/labels", boardID)
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list board labels: %w", err)
	}

	return result, nil
}
