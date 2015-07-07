package recordstore

import (
	"errors"

	cxt "golang.org/x/net/context"

	record "github.com/ipfs/go-iprs/record"
)

var (
	// ErrTimeout is returned when a Put or Get operation times out.
	ErrTimeout = errors.New("record store timeout")

	// ErrInvalid is returned when a Put or Get operation receives
	// an invalid record.
	ErrInvalid = errors.New("invalid record")
)

// Path is a record's location in the Store. It helps
// scope records, and to organize their storage. For example,
// Store implementations could use the path to derive a location:
// - fs-backed at that path
// - dns-backed under a domain derived from the path
// - dht-backed under hash(path)
type Path string

// Store is a record storage system, possibly networked.
// Its interface is very simple, it only allows Put and Get
// on a record. The key
type Store interface {
	// Put adds a record to the record.Store. Multiple records
	// may be put at once to the same path.
	// In networked Stores, this may be a blocking operation.
	// Some Stores may enforce strict consistency, others may be
	// eventually consistent, and some might simply be best effort.
	Put(cxt.Context, Path, record.Record) error

	// Get retrieves the "best" record for a given path from the
	// record.Store. Determining the "best" record is based on
	// the total ordering of records, given by record.Order().
	// In networked Stores, this may be a blocking operation.
	// Some Stores may enforce strict consistency, others may be
	// eventually consistent, and some might simply be best effort.
	Get(cxt.Context, Path) (record.Record, error)

	// GetChan return a channe of records for a given path from the
	// record.Store. The records are returned as they arrive, and
	// thus "better" records may follow. This gives the user control.
	// In networked Stores, this may be a blocking operation.
	// Some Stores may enforce strict consistency, others may be
	// eventually consistent, and some might simply be best effort.
	GetChan(cxt.Context, Path) (<-chan record.Record, error)
}
