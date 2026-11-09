package book

import (
	"time"
	"uuid"
)

// Book is a book as the business layer sees one.
type Book struct {
	ID     uuid.UUID
	Title  string
	Author string
	// ISBN must stay a string: ISBN-10 can start with 0 and end with the
	// letter X, both of which a numeric type destroys.
	ISBN        string
	Cost        int // In cents.
	Quantity    int
	Sold        int
	Revenue     int // In cents.
	UserID      uuid.UUID
	DateCreated time.Time
	DateUpdated time.Time
}
