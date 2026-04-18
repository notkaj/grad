package world

import (
	"fmt"
)

type item struct {
	title string
	ID    string
	Count int
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return fmt.Sprintf("%d Stations", i.Count) }
func (i item) FilterValue() string { return i.title }

func AllItem(acc int) item {
	return item{
		title: "All",
		ID:    "ALL",
		Count: acc,
	}
}

func NewItem(title string, desc string, id string, size int) item {
	return item{
		title: title,
		ID:    id,
		Count: size,
	}
}
