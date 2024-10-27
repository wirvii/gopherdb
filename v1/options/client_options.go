package options

import "strings"

// ClientOptions is a struct that holds the options for the client.
type ClientOptions struct {
	FilePath      string
	CacheSize     int64
	Username      string
	Password      string
	EncryptionKey string
}

// Option is an interface that is used to apply options to the client.
type Option interface {
	Apply(o *ClientOptions)
}

type optionFunc func(*ClientOptions)

func (f optionFunc) Apply(o *ClientOptions) {
	f(o)
}

// WithFilePath is an option that sets the file path for the client.
func WithFilePath(fp string) Option {
	return optionFunc(func(o *ClientOptions) {
		o.FilePath = strings.TrimSpace(fp)
	})
}

// WithCacheSize is an option that sets the cache size for the client.
func WithCacheSize(cs int64) Option {
	return optionFunc(func(o *ClientOptions) {
		o.CacheSize = cs
	})
}

// WithUsername is an option that sets the username for the client.
func WithUsername(u string) Option {
	return optionFunc(func(o *ClientOptions) {
		o.Username = strings.TrimSpace(u)
	})
}

// WithPassword is an option that sets the password for the client.
func WithPassword(p string) Option {
	return optionFunc(func(o *ClientOptions) {
		o.Password = strings.TrimSpace(p)
	})
}

// WithEncryptionKey is an option that sets the encryption key for the client.
func WithEncryptionKey(ek string) Option {
	return optionFunc(func(o *ClientOptions) {
		o.EncryptionKey = strings.TrimSpace(ek)
	})
}
