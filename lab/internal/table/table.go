package table

import (
	"encoding/csv"
	"os"
)

type Table struct {
	w     *csv.Writer
	close func() error
}

func New(filename string) (*Table, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	w := csv.NewWriter(f)
	return &Table{
		w: w,
		close: func() error {
			return f.Close()
		}}, nil
}

func (t *Table) Write(record []string) error {
	return t.w.Write(record)
}

func (t *Table) Flush() error {
	t.w.Flush()

	return t.w.Error()
}

func (t *Table) Close() error {
	return t.close()
}
