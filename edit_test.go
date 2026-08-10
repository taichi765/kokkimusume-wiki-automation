package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindLinesToEditNormal(t *testing.T) {
	testCases := []struct {
		desc  string
		src   string
		start int
		end   int
	}{
		{
			desc: "without whitespace",
			src: `
			TITLE:テストページ

			//@generated_start
			[[なんか]]
			[[なんか2]]
			//@generated_end
			`,
			start: 3,
			end:   6,
		},
		{
			desc: "with whitespace",
			src: `
			TITLE:テストページ

			// @generated_start
			[[なんか]]
			[[なんか2]]
			//@generated_end

			`,
			start: 3,
			end:   6,
		},
		{
			desc: "with comment after keyword",
			src: `
			TITLE:テストページ
			//@generated_start はーじまーるよー
			[[なんか]]
			[[なんか2]]
			//@generated_end おーわりーだよー
			`,
			start: 2,
			end:   5,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			start, end, err := findLinesToEdit(tC.src)
			require.Nil(t, err, "should succeed")
			assert.Equal(t, tC.start, start)
			assert.Equal(t, tC.end, end)
		})
	}
}

func TestFindLinesToEditError(t *testing.T) {
	testCases := []struct {
		desc string
		src  string
	}{
		{
			desc: "without @generated_end",
			src: `
			TITLE:テストページ

			//@generated_start
			[[なんか]]
			[[なんか2]]
			//@gonorotod_end
			`,
		},
		{
			desc: "without @generated_start",
			src: `
			TITLE:テストページ
			//@gonototod_start
			[[なんか]]
			[[なんか2]]
			//@generated_end
			`,
		},
		{
			desc: "multiple start",
			src: `
			TITLE:テストページ
			//@generated_start
			[[なんか]]
			//@generated_start
			[[なんか2]]
			//@generated_end
			`,
		},
		{
			desc: "multiple end",
			src: `
			TITLE:テストページ
			//@generated_start
			[[なんか]]
			//@generated_end
			[[なんか2]]
			//@generated_end
			`,
		},
		{
			desc: "end before start",
			src: `
			TITLE:テストページ
			[[なんか]]
			//@generated_end
			[[なんか2]]
			//@generated_start`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, _, err := findLinesToEdit(tC.src)
			assert.NotNil(t, err)
		})
	}
}
