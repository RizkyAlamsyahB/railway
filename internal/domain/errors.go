package domain

import "errors"

var (
	// ErrStockUnavailable signals that checkout stock changed before the final
	// stock decrement transaction could reserve it.
	ErrStockUnavailable = errors.New("stock unavailable")
)
