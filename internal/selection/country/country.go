package country

import "charm.land/bubbles/v2/list"

type Model struct {
	list list.Model
}

func InitialModel() Model {
	return Model{
		list: list.Model{},
	}
}
