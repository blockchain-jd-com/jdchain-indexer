package handler

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapFromtoToOffset(t *testing.T) {
	var max int64 = 30
	units := []struct {
		from, to      int64
		count, offset int64
	}{
		{
			from: -1, to: 10,
			count: 11, offset: 0,
		},
		{
			from: -1, to: 0,
			count: 1, offset: 0,
		},
		{
			from: -1, to: -1,
			count: 0, offset: 0,
		},
		{
			from: 0, to: -1,
			count: 0, offset: 0,
		},
		{
			from: 2, to: -1,
			count: 0, offset: 1,
		},
		{
			from: 0, to: 10,
			count: 11, offset: 0,
		},
		{
			from: 1, to: 10,
			count: 10, offset: 0,
		},
		{
			from: 2, to: 10,
			count: 9, offset: 1,
		},
		{
			from: 2, to: 100000,
			count: max, offset: 1,
		},
	}

	for _, unit := range units {
		count, offset := mapFromToToOffset(unit.from, unit.to, max)
		assert.Equal(t, unit.count, count, "count", spew.Sdump(unit))
		assert.Equal(t, unit.offset, offset, "offset", spew.Sdump(unit))
	}
}

func TestParseQuery(t *testing.T) {
	var max int64 = 30

	units := []struct {
		from, to                 int64
		expectedFrom, expectedTo int64
		result                   bool
	}{
		{
			from: -1, to: -1,
			expectedFrom: 0, expectedTo: 0,
			result: false,
		},
		{
			from: -1, to: 1000,
			expectedFrom: 0, expectedTo: 0,
			result: false,
		},
		{
			from: 0, to: 1000,
			expectedFrom: 0, expectedTo: max,
			result: true,
		},
		{
			from: 0, to: max - 1,
			expectedFrom: 0, expectedTo: max - 1,
			result: true,
		},
		{
			from: 10, to: 1000,
			expectedFrom: 10, expectedTo: max + 10,
			result: true,
		},
		{
			from: 2, to: 1,
			expectedFrom: 0, expectedTo: 0,
			result: false,
		},
	}

	for _, unit := range units {
		from, to, result := validateFromTo(unit.from, unit.to, max)
		assert.Equal(t, unit.result, result, spew.Sdump(unit))
		assert.Equal(t, unit.expectedFrom, from, spew.Sdump(unit))
		assert.Equal(t, unit.expectedTo, to, spew.Sdump(unit))
	}
}
