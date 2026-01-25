package api

import "github.com/rogeriopvl/fizzy/internal/colors"

type Board struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AllAccess bool   `json:"all_access"`
	CreatedAt string `json:"created_at"`
	URL       string `json:"url"`
	Creator   User   `json:"creator"`
}

type CreateBoardPayload struct {
	Name               string `json:"name"`
	AllAccess          bool   `json:"all_access"`
	AutoPostponePeriod int    `json:"auto_postpone_period"`
	PublicDescription  string `json:"public_description"`
}

type Column struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Color     ColorObject `json:"color"`
	CreatedAt string      `json:"created_at"`
}

type ColorObject struct {
	Name  string `json:"name"`
	Value Color  `json:"value"`
}

type CreateColumnPayload struct {
	Name  string `json:"name"`
	Color *Color `json:"color,omitempty"`
}

type Card struct {
	ID              string   `json:"id"`
	Number          int      `json:"number"`
	Title           string   `json:"title"`
	Status          string   `json:"status"`
	Description     string   `json:"description"`
	DescriptionHTML string   `json:"description_html"`
	ImageURL        string   `json:"image_url"`
	Tags            []string `json:"tags"`
	Golden          bool     `json:"golden"`
	LastActiveAt    string   `json:"last_active_at"`
	CreatedAt       string   `json:"created_at"`
	URL             string   `json:"url"`
	Board           Board    `json:"board"`
	Creator         User     `json:"creator"`
	CommentsURL     string   `json:"comments_url"`
}

type CardFilters struct {
	BoardIDs         []string
	TagIDs           []string
	AssigneeIDs      []string
	CreatorIDs       []string
	CloserIDs        []string
	CardIDs          []string
	IndexedBy        string
	SortedBy         string
	AssignmentStatus string
	CreationStatus   string
	ClosureStatus    string
	Terms            []string
}

type CreateCardPayload struct {
	Title        string   `json:"title"`
	Description  string   `json:"description,omitempty"`
	Status       string   `json:"status,omitempty"`
	ImageURL     string   `json:"image_url,omitempty"`
	TagIDS       []string `json:"tag_ids,omitempty"`
	CreatedAt    string   `json:"created_at,omitempty"`
	LastActiveAt string   `json:"last_active_at,omitempty"`
}

// UpdateCardPayload image not included because we don't support files yet
type UpdateCardPayload struct {
	Title        string   `json:"title,omitempty"`
	Description  string   `json:"description,omitempty"`
	Status       string   `json:"status,omitempty"`
	TagIDS       []string `json:"tag_ids,omitempty"`
	LastActiveAt string   `json:"last_active_at,omitempty"`
}

type GetMyIdentityResponse struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	User      User   `json:"user"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email_address"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	URL       string `json:"url"`
}

type Notification struct {
	ID        string        `json:"id"`
	Read      bool          `json:"read"`
	ReadAt    string        `json:"read_at"`
	CreatedAt string        `json:"created_at"`
	Title     string        `json:"title"`
	Body      string        `json:"body"`
	Creator   User          `json:"creator"`
	Card      CardReference `json:"card"`
	URL       string        `json:"url"`
}

type CardReference struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

type Tag struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	URL       string `json:"url"`
}

type Comment struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      struct {
		PlainText string `json:"plain_text"`
		HTML      string `json:"html"`
	} `json:"body"`
	Creator      User          `json:"creator"`
	Card         CardReference `json:"card"`
	ReactionsURL string        `json:"reactions_url"`
	URL          string        `json:"url"`
}

type Color string

// Color constants using centralized definitions
var (
	Blue   Color = Color(colors.Blue.CSSValue)
	Gray   Color = Color(colors.Gray.CSSValue)
	Tan    Color = Color(colors.Tan.CSSValue)
	Yellow Color = Color(colors.Yellow.CSSValue)
	Lime   Color = Color(colors.Lime.CSSValue)
	Aqua   Color = Color(colors.Aqua.CSSValue)
	Violet Color = Color(colors.Violet.CSSValue)
	Purple Color = Color(colors.Purple.CSSValue)
	Pink   Color = Color(colors.Pink.CSSValue)
)
