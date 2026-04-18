package country

import (
	"fmt"
)

type item struct {
	title       string
	description string
	ID          string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func AllItem(acc int) item {
	return item{
		title:       "All",
		description: fmt.Sprintf("%d Stations", acc),
		ID:          "ALL",
	}
}

func NewItem(title string, desc string, id string) item {
	return item{
		title:       title,
		description: desc,
		ID:          id,
	}
}
