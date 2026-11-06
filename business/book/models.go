package book

import "time"

// Book represents an individual book.
type Book struct {
	ID          string    `json:"id"`           // Unique identifier.
	Title       string    `json:"title"`        // Display title of the book.
	Author      string    `json:"author"`       // Display the author of the book.
	ISBN        int       `json:"isbn"`         // International Standard Book Number.
	Cost        int       `json:"cost"`         // Price for one item in cents.
	Quantity    int       `json:"quantity"`     // Original number of items available.
	Sold        int       `json:"sold"`         // Aggregate field showing number of items sold.
	Revenue     int       `json:"revenue"`      // Aggregate field showing total cost of sold items.
	UserID      string    `json:"user_id"`      // ID of the user who created the product.
	DateCreated time.Time `json:"date_created"` // When the product was added.
	DateUpdated time.Time `json:"date_updated"` // When the product record was last modified.
}
