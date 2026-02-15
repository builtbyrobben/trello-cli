package trello

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/builtbyrobben/trello-cli/internal/api"
)

func newTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	return &Client{
		Client: api.NewClient("",
			api.WithBaseURL(srv.URL),
			api.WithQueryAuth(map[string]string{
				"key":   "testkey",
				"token": "testtoken",
			}),
		),
	}
}

func TestBoardsList(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/members/me/boards" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		boards := []Board{
			{ID: "board1", Name: "My Board"},
			{ID: "board2", Name: "Other Board"},
		}

		json.NewEncoder(w).Encode(boards)
	}))

	boards, err := client.Boards().List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(boards) != 2 {
		t.Fatalf("expected 2 boards, got %d", len(boards))
	}

	if boards[0].Name != "My Board" {
		t.Errorf("unexpected board name: %s", boards[0].Name)
	}
}

func TestBoardsGet(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards/board1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(Board{ID: "board1", Name: "My Board"})
	}))

	board, err := client.Boards().Get(context.Background(), "board1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if board.ID != "board1" {
		t.Errorf("unexpected board ID: %s", board.ID)
	}
}

func TestBoardsGet_EmptyID(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	_, err := client.Boards().Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestListsList(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards/board1/lists" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		lists := []List{
			{ID: "list1", Name: "To Do", IDBoard: "board1"},
			{ID: "list2", Name: "Done", IDBoard: "board1"},
		}

		json.NewEncoder(w).Encode(lists)
	}))

	lists, err := client.Lists().List(context.Background(), "board1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(lists) != 2 {
		t.Fatalf("expected 2 lists, got %d", len(lists))
	}
}

func TestListsCreate(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body CreateListRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.Name != "New List" {
			t.Errorf("unexpected name: %s", body.Name)
		}

		if body.IDBoard != "board1" {
			t.Errorf("unexpected board: %s", body.IDBoard)
		}

		json.NewEncoder(w).Encode(List{ID: "list3", Name: "New List", IDBoard: "board1"})
	}))

	list, err := client.Lists().Create(context.Background(), "board1", "New List")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if list.Name != "New List" {
		t.Errorf("unexpected list name: %s", list.Name)
	}
}

func TestCardsCreate(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body CreateCardRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.Name != "My Card" {
			t.Errorf("unexpected name: %s", body.Name)
		}

		json.NewEncoder(w).Encode(Card{ID: "card1", Name: "My Card", IDList: "list1"})
	}))

	card, err := client.Cards().Create(context.Background(), "list1", "My Card", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if card.Name != "My Card" {
		t.Errorf("unexpected card name: %s", card.Name)
	}
}

func TestCardsCreate_EmptyName(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	_, err := client.Cards().Create(context.Background(), "list1", "", "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCardsGet(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cards/card1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(Card{ID: "card1", Name: "Test Card"})
	}))

	card, err := client.Cards().Get(context.Background(), "card1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if card.ID != "card1" {
		t.Errorf("unexpected card ID: %s", card.ID)
	}
}

func TestCardsDelete(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))

	err := client.Cards().Delete(context.Background(), "card1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMembersMe(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/members/me" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(Member{ID: "member1", Username: "testuser", FullName: "Test User"})
	}))

	member, err := client.Members().Me(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if member.Username != "testuser" {
		t.Errorf("unexpected username: %s", member.Username)
	}
}

func TestLabelsList(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards/board1/labels" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		labels := []Label{
			{ID: "label1", Name: "Bug", Color: "red"},
			{ID: "label2", Name: "Feature", Color: "green"},
		}

		json.NewEncoder(w).Encode(labels)
	}))

	labels, err := client.Labels().ListByBoard(context.Background(), "board1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}

	if labels[0].Color != "red" {
		t.Errorf("unexpected color: %s", labels[0].Color)
	}
}
