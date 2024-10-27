package gopherdb

// Database is a struct that holds the database information.
type Database struct {
	name   string
	client *Client
}

// Name is a function that returns the name of the database.
func (d *Database) Name() string {
	return d.name
}

// Client is a function that returns the client.
func (d *Database) Client() *Client {
	return d.client
}

// Collection is a function that returns a collection.
func (d *Database) Collection(name string) *Collection {
	return &Collection{
		name:     name,
		database: d,
	}
}
