package recordstore_datastore

import (
	"errors"

	dag "github.com/ipfs/go-ipfs/merkledag"
	ds "github.com/jbenet/go-datastore"
	cxt "golang.org/x/net/context"

	record "github.com/ipfs/go-iprs/record"
	store "github.com/ipfs/go-iprs/store"
)

type Store struct {
	ts *record.TypeSet
	ds ds.Datastore
}

// New constructs a record.Store from a given TypeSet and
// Datastore.
func New(ts *record.TypeSet, ds ds.Datastore) *Store {
	return &Store{ts, ds}
}

// Put adds a record to the record.Store. Multiple records
// may be put at once to the same path.
// In networked Stores, this may be a blocking operation.
// Some Stores may enforce strict consistency, others may be
// eventually consistent, and some might simply be best effort.
func (s *Store) Put(ctx cxt.Context, p Path, r record.Record) error {
	k := ds.NewKey(string(p))
	m, err := record.Marshal(r)
	if err != nil {
		return err
	}

	return s.ds.Put(k, m)
}

// Get retrieves the "best" record for a given path from the
// record.Store. Determining the "best" record is based on
// the total ordering of records, given by record.Order().
// In networked Stores, this may be a blocking operation.
// Some Stores may enforce strict consistency, others may be
// eventually consistent, and some might simply be best effort.
func (s *Store) Get(ctx cxt.Context, p Path) (record.Record, error) {
	k := ds.NewKey(string(p))
	m, err := s.ds.Get(k)
	if err != nil {
		return nil, err
	}

	return s.ts.Unmarshal(m)
}

// GetChan return a channe of records for a given path from the
// record.Store. The records are returned as they arrive, and
// thus "better" records may follow. This gives the user control.
// In networked Stores, this may be a blocking operation.
// Some Stores may enforce strict consistency, others may be
// eventually consistent, and some might simply be best effort.
func (s *Store) GetChan(ctx cxt.Context, p Path) (<-chan record.Record, error) {
	r, err := s.Get(ctx, p)
	if err != nil {
		return nil, err
	}

	ch := make(chan record.Record, 1)
	ch <- r
	close(ch)
	return ch, nil
}
