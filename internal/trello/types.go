package trello

// Board represents a Trello board.
type Board struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Desc           string `json:"desc,omitempty"`
	Closed         bool   `json:"closed"`
	URL            string `json:"url,omitempty"`
	ShortURL       string `json:"short_url,omitempty"`
	IDOrganization string `json:"id_organization,omitempty"`
}

// List represents a Trello list on a board.
type List struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Closed  bool   `json:"closed"`
	IDBoard string `json:"id_board"`
	Pos     int    `json:"pos"`
}

// Card represents a Trello card.
type Card struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Desc      string   `json:"desc,omitempty"`
	Closed    bool     `json:"closed"`
	IDBoard   string   `json:"id_board"`
	IDList    string   `json:"id_list"`
	URL       string   `json:"url,omitempty"`
	ShortURL  string   `json:"short_url,omitempty"`
	Pos       float64  `json:"pos"`
	Due       string   `json:"due,omitempty"`
	IDMembers []string `json:"id_members,omitempty"`
	IDLabels  []string `json:"id_labels,omitempty"`
}

// Member represents a Trello member.
type Member struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name,omitempty"`
	URL      string `json:"url,omitempty"`
}

// Label represents a Trello label.
type Label struct {
	ID      string `json:"id"`
	IDBoard string `json:"id_board"`
	Name    string `json:"name"`
	Color   string `json:"color,omitempty"`
}

// CreateListRequest is the body for creating a new list.
type CreateListRequest struct {
	Name    string `json:"name"`
	IDBoard string `json:"idBoard"` //nolint:tagliatelle // Trello API uses camelCase
}

// CreateCardRequest is the body for creating a new card.
type CreateCardRequest struct {
	Name   string `json:"name"`
	IDList string `json:"idList"` //nolint:tagliatelle // Trello API uses camelCase
	Desc   string `json:"desc,omitempty"`
}

// UpdateCardRequest is the body for updating a card.
type UpdateCardRequest struct {
	Name string `json:"name,omitempty"`
	Desc string `json:"desc,omitempty"`
}
