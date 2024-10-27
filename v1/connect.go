package gopherdb

import (
	. "github.com/wirvii/gopherdb/internal/consts"
	"github.com/wirvii/gopherdb/v1/options"
)

// Connect is a function that connects the client to the database.
func Connect(opts ...options.Option) (*Client, error) {
	clientOptions := &options.ClientOptions{
		CacheSize: OneHundredsMB,
	}

	for _, opt := range opts {
		opt.Apply(clientOptions)
	}

	c := newClient(clientOptions)
	if err := c.connect(); err != nil {
		return nil, err
	}

	return c, nil
}
