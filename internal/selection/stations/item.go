package stations

type item struct {
	title       string
	description string
	id          string
	url         string
	homepage    string
	codec       string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func NewItem(title string, desc string, id string, url string, homepage string, codec string) item {
	return item{
		title:       title,
		description: desc,
		id:          id,
		url:         url,
		homepage:    homepage,
		codec:       codec,
	}
}
