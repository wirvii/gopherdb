package gopherdb

import (
	. "github.com/wirvii/gopherdb/internal/consts"
	"github.com/wirvii/gopherdb/internal/hash"
	"github.com/wirvii/gopherdb/internal/log"
	"github.com/wirvii/gopherdb/internal/storage"
	"github.com/wirvii/gopherdb/v1/options"
)

// Client is a struct that holds the client information.
type Client struct {
	id            string
	logger        log.Logger
	clientOptions *options.ClientOptions
	storage       storage.Storage
}

// newClient is a function that creates a new client.
func newClient(opts *options.ClientOptions) *Client {
	return &Client{
		id:            hash.UUIDv4(),
		logger:        log.NewLogger(),
		clientOptions: opts,
	}
}

// connect is a function that connects the client to the database.
func (c *Client) connect() error {
	s, err := storage.New(
		storage.WithFilePath(c.clientOptions.FilePath),
		storage.WithCacheSize(c.clientOptions.CacheSize),
		storage.WithEncryptionKey(c.clientOptions.EncryptionKey),
	)
	if err != nil {
		c.logger.Errorw(
			"failed to connect to storage",
			"error", err,
		)

		return err
	}

	c.storage = s

	return nil
}

// Database is a function that returns a database.
func (c *Client) Database(name string) *Database {
	panic("unimplemented")
}

// initializeAdminDatabase is a function that initializes the admin database.
func (c *Client) initializeAdminDatabase() error {
	err := c.storage.CreateDatabase(AdminDBName)
	if err != nil {
		c.logger.Errorw(
			"failed to create admin database",
			"error", err,
		)

		return err
	}

	return nil
}
