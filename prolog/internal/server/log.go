package server

import (
	"fmt"
	"sync"
)

type Log struct {
	mu      sync.Mutex
	records []Record
}

type Record struct {
	Value  []byte `json:"value"`
	OffSet uint64 `json:"offSet"`
}

var ErrOfsetNotFound = fmt.Errorf("offset not found")

func NewLog() *Log {
	return &Log{}
}

func (c *Log) Append(record Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record.OffSet = uint64(len(c.records))
	c.records = append(c.records, record)
	return record.OffSet, nil
}

func (c *Log) Read(offset uint64) (Record, error) {

	c.mu.Lock()
	defer c.mu.Unlock()
	if offset > uint64(len(c.records)) {
		return Record{}, ErrOfsetNotFound
	}
	record := c.records[offset]
	return record, nil
}
