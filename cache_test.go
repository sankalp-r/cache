package cache

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheCreation(t *testing.T) {
	tests := []struct {
		desc     string
		capacity int
		err      error
	}{
		{
			desc:     "valid capacity",
			capacity: 3,
			err:      nil,
		},
		{
			desc:     "invalid capacity",
			capacity: 0,
			err:      errors.New("invalid capacity: capacity should be greater than zero"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			_, err := NewCache(test.capacity)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestLRUCacheOperations(t *testing.T) {

	tests := []struct {
		desc      string
		capacity  int
		operation []testOperation
		// expected result for Put Operation is empty string, for Get it can be empty/non-empty value
		expected []string
	}{
		{
			"cache operation with key-eviction",
			5,
			[]testOperation{
				{opr: "Put", args: []string{"1", "1"}},
				{opr: "Put", args: []string{"2", "2"}},
				{opr: "Get", args: []string{"1"}},
				{opr: "Put", args: []string{"3", "3"}},
				{opr: "Put", args: []string{"4", "4"}},
				{opr: "Put", args: []string{"5", "5"}},
				{opr: "Put", args: []string{"6", "6"}},
				{opr: "Get", args: []string{"2"}}, // should be evicted
			},
			[]string{"", "", "1", "", "", "", "", ""},
		},

		{
			"cache operation with key-overwrite",
			3,
			[]testOperation{
				{opr: "Put", args: []string{"1", "1"}},
				{opr: "Put", args: []string{"2", "2"}},
				{opr: "Get", args: []string{"1"}},
				{opr: "Put", args: []string{"1", "1.5"}}, //overwrite value
				{opr: "Get", args: []string{"1"}},
				{opr: "Get", args: []string{"1.5"}}, // non-existent key
			},
			[]string{"", "", "1", "", "1.5", ""},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			testCache, err := NewCache(test.capacity)
			assert.Equal(t, err, nil)
			var actual []string
			for _, opr := range test.operation {
				if opr.opr == "Put" {
					testCache.Put(opr.args[0], opr.args[1])
					actual = append(actual, "")
				} else if opr.opr == "Get" {
					val := testCache.Get(opr.args[0])
					actual = append(actual, val)
				}
			}
			assert.Equal(t, test.expected, actual)
		})
	}
}

type testOperation struct {
	opr  string
	args []string
}
