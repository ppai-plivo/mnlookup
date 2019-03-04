package store

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type Store struct {
	tree *Tree
}

func (t *Store) Len() int {
	return t.tree.Len()
}

func (t *Store) Lookup(number string) (Value, error) {
	_, value, ok := t.tree.LongestPrefix(number)
	if !ok {
		return value, fmt.Errorf("Not Found")
	}

	return value, nil
}

func New(reader io.Reader) (*Store, error) {

	r := csv.NewReader(reader)
	r.ReuseRecord = true

	// read column titles
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	tree := NewTree()

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if record[8] == "" || record[9] == "" {
			// MCC and MNC not empty
			continue
		}

		prefix := string([]byte(record[0])) // force GC
		mcc, _ := strconv.Atoi(record[8])
		mnc, _ := strconv.Atoi(record[9])

		_, _ = tree.Insert(prefix,
			Value{
				MCC: uint16(mcc),
				MNC: uint16(mnc),
			})
	}

	return &Store{
		tree: tree,
	}, nil
}
