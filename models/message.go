package models

// Messages structure
// Note that you need to change the DB configuration to change this structure

type Message struct {
	ID      int
	Content string
	Link    string
}
