package options

// DeleteOptions es un struct que contiene las opciones para una eliminación.
type DeleteOptions struct {
	Limit *int
}

// Delete crea una nueva instancia de deleteOptions.
func Delete() *DeleteOptions {
	return &DeleteOptions{}
}

// Merge combina las opciones de varias consultas.
func (o *DeleteOptions) Merge(opts ...*DeleteOptions) *DeleteOptions {
	for _, opt := range opts {
		if opt.Limit != nil {
			o.Limit = opt.Limit
		}
	}

	return o
}
