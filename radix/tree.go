package radix

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/ppai-plivo/mnlookup/api"

	radix "github.com/hashicorp/go-immutable-radix"
)

type RadixTree struct {
	tree *radix.Tree
}

func (t *RadixTree) Dump() {
	walkfn := func(k []byte, v interface{}) bool {
		fmt.Println(string(k), v)
		return false
	}

	t.tree.Root().Walk(walkfn)
}

func (t *RadixTree) Lookup(number string) (*api.Record, error) {
	_, value, ok := t.tree.Root().LongestPrefix([]byte(number))
	if !ok {
		return nil, fmt.Errorf("Not Found")
	}

	r, ok := value.(api.Record)
	if !ok {
		return nil, fmt.Errorf("Entry corrupt")
	}

	return &r, nil
}

func New(reader io.Reader) (*RadixTree, error) {

	r := csv.NewReader(reader)

	// read column titles
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	tree := radix.New()

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		prefix := record[0]
		number_type := record[7]
		mcc := record[8]
		mnc := record[9]

		if (number_type != "MOB") || mcc == "" || mnc == "" {
			continue
		}

		tree, _, _ = tree.Insert([]byte(prefix),
			api.Record{
				Prefix: prefix,
				MCC:    mcc,
				MNC:    mnc,
			})
	}

	return &RadixTree{
		tree: tree,
	}, nil
}
