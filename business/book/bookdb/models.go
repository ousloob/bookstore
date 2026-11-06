package bookdb

import "time"

// Book represents an individual book.
type Book struct {
	Title       string    `db:"title"`        // Display title of the book.
	Author      string    `db:"author"`       // Display the author of the book.
	ISBN        int       `db:"isbn"`         // International Standard Book Number.
	Cost        int       `db:"cost"`         // Price for one item in cents.
	Quantity    int       `db:"quantity"`     // Original number of items available.
	Sold        int       `db:"sold"`         // Aggregate field showing number of items sold.
	Revenue     int       `db:"revenue"`      // Aggregate field showing total cost of sold items.
	UserID      string    `db:"user_id"`      // ID of the user who created the product.
	DateCreated time.Time `db:"date_created"` // When the product was added.
	DateUpdated time.Time `db:"date_updated"` // When the product record was last modified.
	ID          string    `db:"id"`           // Unique identifier.
}
