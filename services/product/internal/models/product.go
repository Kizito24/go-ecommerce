package models

// Product represents the document in MongoDB
type Product struct {
	ID          int64   `bson:"_id"` // We are forcing int64 to match Proto
	Name        string  `bson:"name"`
	Description string  `bson:"description"`
	Price       float32 `bson:"price"`
	Stock       int64   `bson:"stock"`
}
