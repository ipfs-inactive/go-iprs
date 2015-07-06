package records

import (
	dag "github.com/ipfs/go-ipfs/merkledag"
)

type Type interface {
	// Node returns a DAG.Node representing the record.Type.
	Node() dag.Node

	// Validator returns an object that can determine the validity
	// of a Record, and that can order records deterministically.
	// The Validator is specific to this Type.
	Validator() Validator
}

// Record is an object that stores some data relevant to a distributed
// system. It is a necesary part of most distributed systems -- it is a sort of
// "glue" that improves how they operate.
type Record interface {
	// Node returns a DAG.Node representing the Record.
	Node() dag.Node

	// Type returns the type of the record.
	Type() Type

	// Version returns the version number of a record.
	Version() int

	// Data returns data carried by the record used in determining validity
	Validity() []byte

	// Value returns data carried by the record to be used
	Value() []byte
}

// Validator represents an algorithm for determining the validity of a
// Record, and to order records deterministically.
type Validator interface {

	// Type returns the type of the validator.
	Type() Type

	// Valid returns whether a Record is valid in present circumstances.
	// This could include checking correctness of the record (checking cryptographic
	// signatures, and the like), or some validity regarding external infrastructure
	// such as:
	// - PKIs (Public Key Infrastructures)  -- is a signature chain valid?
	// - TIs (Time Infrastructures) -- is this record valid _right now_?
	Valid(Record) (bool, error)

	// Order returns {-1, 0, 1} to order (and pick from) {a, b}.
	Order(a, b Record) int
}

// IsValid checks a record's validity with the appropriate validator,
// and returns whether a Record is valid. err may be non nil if there
// is an error checking the validity. An error counts as invalid, so
// IsValid will always return ok == false when err != nil.
func IsValid(r Record) (ok bool, err error) {
	return r.Type().Validator().Valid(r)
}

// Order returns {-1, 0, 1} to order (and pick from) {a, b}. Orders by:
//
//   ( cmp(a.Version(), b.Version()),     // 1) version numbers always take precedence
//     validator.Order(a, b),             // 2) user's validator.Order(.) function
//     cmp(Marshal(a), Marshal(b) )       // 3) worst case, order by raw bytes.
//
// 1) A higher version number _always_ takes precedence over a lower version number.
// Record systems could use version numbers primarily for delivering updates, but
// SHOULD still address ordering records with equal version numbers (multipler
// writer problem).
//
// 2) The user's validator's Order function is used next to determine the order
// of records. Thus the user may define ordering based on timestamps on the record,
// or on some (pure) computation based on the record.
//
// 3) In the worst case, records are orderered by cmp( Marshal(a), Marshal(b) )
// to ensure there is _always_ a deterministic way to order records. This also
// lets the user define Validator.Order(.) functions to always return 0 and
// the record system _will still be deterministic_ (a very important property).
func Order(validator Validator, a, b Record) int {

	// 1) version numbers always take precedence
	av := a.Version()
	bv := b.Version()
	switch {
	case av < bv:
		return -1
	case av > bv:
		return 1
	default:
	}

	// 2) user's validator.Order(.) function
	if vo := validator.Order(a, b); vo != 0 {
		return vo
	}

	// 3) worst case, order by raw bytes.
	am := mustMarshal(a)
	bm := mustMarshal(b)
	return bytes.Compare(am, bm)
}

// mustMarshal returns a byte representation of a record.
func mustMarshal(Record) []byte {
	buf, err := Marshal()
	if err != nil {
		panic("marshal failed. record should never be invalid")
	}
	return buf
}

// Marshal returns a byte representation of a record.
func Marshal(r Record) ([]byte, error) {
	return r.Node().Marshal()
}

// Unmarshal returns a Record instance from its byte representation.
func Unmarshal([]byte) (Record, error) {
	panic("not yet implemented")
}
