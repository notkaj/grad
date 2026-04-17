package selector

type Item struct {
	title       string
	description string
	ID          string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.description }
func (i Item) FilterValue() string { return i.title }

var AllItem Item = Item{title: "ALL", ID: "ALL"}

func NewItem(title string, desc string, id string) *Item {
	return &Item{
		title:       title,
		description: desc,
		ID:          id,
	}
}
