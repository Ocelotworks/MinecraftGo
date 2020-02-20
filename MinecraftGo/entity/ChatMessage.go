package entity

type ChatMessage struct {
	Translate string                  `json:"translate"`
	With      *[]ChatMessageComponent `json:"with"`
}

type ChatMessageComponent struct {
	Text       string                  `json:"text"`
	Bold       *bool                   `json:"bold,omitempty"`
	Italic     *bool                   `json:"italic,omitempty"`
	Underlined *bool                   `json:"underlined,omitempty"`
	Obfuscated *bool                   `json:"obfuscated,omitempty"`
	Colour     *ChatMessageColour      `json:"color,omitempty"`
	Insertion  *string                 `json:"insertion,omitempty"`
	ClickEvent *ChatMessageClickEvent  `json:"clickEvent,omitempty"`
	HoverEvent *ChatMessageHoverEvent  `json:"hoverEvent,omitempty"`
	Children   *[]ChatMessageComponent `json:"extra,omitempty"`
}

type ChatMessageColour string

const (
	Black         ChatMessageColour = "black"
	DarkBlue      ChatMessageColour = "dark_blue"
	DarkGreen     ChatMessageColour = "dark_green"
	DarkAqua      ChatMessageColour = "dark_aqua"
	DarkRed       ChatMessageColour = "dark_red"
	Purple        ChatMessageColour = "dark_purple"
	Gold          ChatMessageColour = "gold"
	Grey          ChatMessageColour = "gray"
	Gray          ChatMessageColour = "gray"
	DarkGrey      ChatMessageColour = "dark_gray"
	DarkGray      ChatMessageColour = "dark_gray"
	Blue          ChatMessageColour = "blue"
	BrightGreen   ChatMessageColour = "green"
	Cyan          ChatMessageColour = "aqua"
	Red           ChatMessageColour = "red"
	Pink          ChatMessageColour = "light_purple"
	Yellow        ChatMessageColour = "yellow"
	White         ChatMessageColour = "white"
	Random        ChatMessageColour = "obfuscated"
	Bold          ChatMessageColour = "bold"
	StrikeThrough ChatMessageColour = "strikethrough"
	Underline     ChatMessageColour = "underline"
	Italic        ChatMessageColour = "italic"
)

type ChatMessageClickEvent struct {
	OpenUrl        *string `json:"open_url,omitempty"`
	OpenFile       *string `json:"open_file,omitempty"`
	RunCommand     *string `json:"run_command,omitempty"`
	SuggestCommand *string `json:"suggest_command,omitempty"`
	ChangePage     *int    `json:"change_page"`
}

type ChatMessageHoverEvent struct {
	ShowChatMessage *ChatMessageComponent `json:"show_text"`
	ShowText        *string               `json:"show_text"`
	//ShowItem TODO: NBT STRUCT JSON
	//ShowEntity same as above
}
