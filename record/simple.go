package record

import (
// dag "github.com/ipfs/go-ipfs/merkledag"

// pb "github.com/ipfs/go-iprs/record/pb"
)

// type untypedRecord struct {
// 	rd *pb.Record
// 	nd *dag.Node
// }

// func newUntypedRecord(nd *dag.Node) (*untypedRecord, error) {
// 	rd, err := pb.Un(nd.Data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &untypedRecord{rd, nd}
// }

// // Node returns a DAG.Node representing the Record.
// func (ur *untypedRecord) Node() *dag.Node {
// 	return ur.nd
// }

// // Type returns the type of the record.
// func (ur *untypedRecord) Type() Type {
// 	return nil
// }

// // Version returns the version number of a record.
// func (ur *untypedRecord) Version() int {
// 	return ur.rd.Version
// }

// // Data returns data carried (or linked) by the record used in determining validity
// func (ur *untypedRecord) Validity() []byte {
// 	return ur.rd.Validity
// }

// // Value returns data carried (or linked) by the record to be used
// func (ur *untypedRecord) Value() []byte {
// 	return ur.rd.Value
// }
