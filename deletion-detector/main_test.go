package main

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/taichi765/kokkimusume-wiki-automation/wikiwiki"
)

func TestCountPageChanges(t *testing.T) {
	testCases := []struct {
		desc    string
		prev    []string
		curr    []string
		created int
		deleted int
	}{
		{
			desc: "deleted only",
			prev: []string{
				"Page 1",
				"Page 2",
				"Page 3",
			},
			curr: []string{
				"Page 1",
			},
			created: 0,
			deleted: 2,
		},
		{
			desc: "created only",
			prev: []string{
				"Page 1",
			},
			curr: []string{
				"Page 1",
				"Page 2",
				"Page 3",
			},
			created: 2,
			deleted: 0,
		},
		{
			desc: "both deleted and created",
			prev: []string{
				"Page 1",
				"Page 2",
			},
			curr: []string{
				"Page 2",
				"Page 3",
			},
			created: 1,
			deleted: 1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			prev := make([]wikiwiki.GeneralPageInfo, len(tC.prev))
			for _, v := range tC.prev {
				prev = append(prev, wikiwiki.GeneralPageInfo{
					Name: v,
				})
			}

			curr := make([]wikiwiki.GeneralPageInfo, len(tC.curr))
			for _, v := range tC.curr {
				curr = append(curr, wikiwiki.GeneralPageInfo{
					Name: v,
				})
			}

			created, deleted := countPageChanges(prev, curr)
			assert.Equal(t, tC.created, created)
			assert.Equal(t, tC.deleted, deleted)
		})
	}
}
