package record

import (
	"errors"
	"fmt"

	key "github.com/ipfs/go-ipfs/blocks/key"
	dag "github.com/ipfs/go-ipfs/merkledag"
)

var (
	linkType  = "@type"
	linkValue = "value"
)

// mustMarshal returns a byte representation of a record.
func mustMarshal(r Record) []byte {
	buf, err := Marshal(r)
	if err != nil {
		panic("marshal failed. record should never be invalid")
	}
	return buf
}

// Marshal returns a byte representation of a record.
func Marshal(r Record) ([]byte, error) {
	return r.Node().Marshal()
}

// requiredLinks is a set of links required from records
var requiredLinks = []string{
	linkType,
	linkValue,
}

// NodeHasRequiredLinks ensures a node has required links
func NodeHasRequiredLinks(nd *dag.Node, links []string) error {
	for _, n := range links {
		found := false
		for _, link := range nd.Links {
			if link.Name == n {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid record. missing link: %s", n)
		}
	}

	return nil
}

// unmarshalNode returns a Record instance from its byte representation.
func unmarshalNode(buf []byte) (*dag.Node, error) {
	nd, err := dag.Decoded(buf)
	if err != nil {
		return nil, err
	}

	if err := NodeHasRequiredLinks(nd, requiredLinks); err != nil {
		return nil, err
	}

	// ok looks good.
	return nd, nil
}

func unmarshalTypeNode(t Type, buf []byte) (*dag.Node, error) {
	nd, err := unmarshalNode(buf)
	if err != nil {
		return nil, err
	}

	if ok, err := NodeHasType(nd, t); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("record type mismatch")
	}

	return nd, nil
}

// UnmarshalType returns a Record instance from its byte representation,
// acord to given type.
func UnmarshalType(t Type, buf []byte) (Record, error) {
	nd, err := unmarshalTypeNode(t, buf)
	if err != nil {
		return nil, err
	}

	return t.New(nd)
}

// UnmarshalTypeSet returns a Record instance from its byte representation,
// acord to type chosen from given typeset.
func UnmarshalFromSet(ts *TypeSet, buf []byte) (Record, error) {
	nd, err := unmarshalNode(buf)
	if err != nil {
		return nil, err
	}

	nt, err := nodeType(nd)
	if err != nil {
		return nil, err
	}

	t := ts.Type(string(nt))
	if t == nil {
		return nil, fmt.Errorf("unsupported record type %s", key.Key(nt).Pretty())
	}

	return t.New(nd)
}

// nodeType returns the record node's type link value.
func nodeType(nd *dag.Node) ([]byte, error) {
	tl, err := nd.GetNodeLink(linkType)
	if err != nil {
		return nil, fmt.Errorf("invalid record. failed to get link %s (%s)", linkType, err)
	}

	return tl.Hash, nil
}
