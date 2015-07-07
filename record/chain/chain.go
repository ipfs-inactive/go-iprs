package record

import (
	proto "code.google.com/p/goprotobuf/proto"
	dag "github.com/ipfs/go-ipfs/merkledag"

	pb "github.com/ipfs/go-iprs/record/pb"
)

// links
var (
	parentLink = "parent"
	recordLink = "record"
)

var Type = &chainType{}

type chainType struct {
	nd *dag.Node
	v  *Validator
}

// Node returns a DAG.Node representing the record.Type.
func (t *chainType) Node() *dag.Node {
	return t.nd
}

// Validator returns an object that can determine the validity
// of a Record, and that can order records deterministically.
// The Validator is specific to this Type.
func (t *chainType) Validator() *Validator {
	return t.Validator
}

// New constructs a new record from given node.
func (t *chainType) New(nd *dag.Node) (Record, error) {
	if err := NodeHasRequiredLinks(nd, []string{parentLink, recordLink}); err != nil {
		return nil, err
	}

	// should get subrecord

	var rd *pb.Record
	err := proto.Unmarshal(nd.Data, rd)
	if err != nil {
		return nil, err
	}

	return &Record{rd, nd}, nil
}

type Record struct {
	rd *pb.Record
	nd *dag.Node
}

// Node returns a DAG.Node representing the Record.
func (r *Record) Parents(ctx cxt.Context, ds dag.DAGService) ([]*Record, error) {
	var parents []*Record
	for _, l := range r.nd.Links {
		if l.Name == parentLink {
			pnd, err := l.GetNode(ctx, ds)
			if err != nil {
				return nil, err
			}
			parents = append(parents, p)
		}
	}

}

// Node returns a DAG.Node representing the Record.
func (r *Record) Node() *dag.Node {
	return r.nd
}

// Type returns the type of the record.
func (r *Record) Type() record.Type {
	return Type
}

// Version returns the version number of a record.
func (r *Record) Version() int {
	return r.rd.Version
}

// Data returns data carried (or linked) by the record used in determining validity
func (r *Record) Validity() []byte {
	return r.rd.Validity
}

// Value returns data carried (or linked) by the record to be used
func (r *Record) Value() []byte {
	return r.rd.Value
}
