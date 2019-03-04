package radix

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/ppai-plivo/mnlookup/api"

	radix "github.com/armon/go-radix"
)

type RadixTree struct {
	tree           *radix.Tree
	countByNumType map[string]int
}

func (t *RadixTree) CountByNumType() map[string]int {
	return t.countByNumType
}

func (t *RadixTree) Lookup(number string) (*api.Record, error) {
	_, value, ok := t.tree.LongestPrefix(number)
	if !ok {
		return nil, fmt.Errorf("Not Found")
	}

	r, ok := value.(*api.Record)
	if !ok {
		return nil, fmt.Errorf("Entry corrupt")
	}

	return r, nil
}

func New(reader io.Reader) (*RadixTree, error) {

	r := csv.NewReader(reader)
	r.ReuseRecord = true

	// read column titles
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	tree := radix.New()

	countByNumType := make(map[string]int)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		countByNumType[record[7]]++

		prefix := record[0]
		mcc := record[8]
		mnc := record[9]

		if mcc == "" || mnc == "" {
			continue
		}

		_, _ = tree.Insert(prefix,
			&api.Record{
				Prefix: prefix,
				MCC:    mcc,
				MNC:    mnc,
			})
	}

	return &RadixTree{
		tree:           tree,
		countByNumType: countByNumType,
	}, nil
}
