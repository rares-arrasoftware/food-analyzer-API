package database

// Database defines a generic interface for CRUD operations on any data model T.
// Intended to be implemented by concrete backends like SQLite, Postgres, etc.
type Database[T any] interface {
	// Create inserts a new item into the database.
	Create(item T) error

	// Update modifies an existing item in the database.
	Update(item T) error

	// GetByID retrieves an item by its primary key (ID).
	GetByID(id uint) (*T, error)

	// GetByField retrieves a single item by a specific field and value.
	// Example: GetByField("email", "test@example.com")
	GetByField(field string, value any) (*T, error)

	// Delete removes an item by its ID.
	Delete(id uint) error
}

func NewDatabase[T any](dsn string, model T) (Database[T], error) {
	// here we can decide if we want to support other database(postgresql, etc.)
	return newSqlite(dsn, model)
}
