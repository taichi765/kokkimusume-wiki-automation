package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/taichi765/kokkimusume-wiki-automation/common"
)

var sampleCharaData = []common.CharacterData{
	{
		Name:                "日本",
		Area:                "東アジア",
		FirstAppearenceDate: "2026/06/07",
	},
	{
		Name:                "韓国",
		Area:                "東アジア",
		FirstAppearenceDate: "2026/06/07",
	},
}

var multipleAreasCharaData = []common.CharacterData{
	{
		Name:                "日本",
		Area:                "東アジア",
		FirstAppearenceDate: "2026/06/07",
	},
	{
		Name:                "韓国",
		Area:                "東アジア",
		FirstAppearenceDate: "2026/06/07",
	},
	{
		Name:                "イタリア",
		Area:                "ヨーロッパ",
		FirstAppearenceDate: "2026/06/29",
	},
	{
		Name:                "カナダ",
		Area:                "北米",
		FirstAppearenceDate: "2026/06/13",
	},
}

func TestSplitLinesToEditNormal(t *testing.T) {
	testCases := []struct {
		desc        string
		src         string
		beforeStart string
		afterEnd    string
	}{
		{
			desc: "without whitespace",
			src: `TITLE:テストページ

//@generated_start
[[なんか]]
[[なんか2]]
//@generated_end
`,
			beforeStart: `TITLE:テストページ

//@generated_start
`,
			afterEnd: `//@generated_end
`,
		},
		{
			desc: "with whitespace",
			src: `TITLE:テストページ

// @generated_start
[[なんか]]
[[なんか2]]
//@generated_end

`,
			beforeStart: `TITLE:テストページ

// @generated_start
`,
			afterEnd: `//@generated_end

`,
		},
		{
			desc: "with comment after keyword",
			src: `TITLE:テストページ
//@generated_start はーじまーるよー
[[なんか]]
[[なんか2]]
//@generated_end おーわりーだよー
`,
			beforeStart: `TITLE:テストページ
//@generated_start はーじまーるよー
`,
			afterEnd: `//@generated_end おーわりーだよー
`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			beforeStart, afterEnd, err := splitLinesToEdit(tC.src)
			require.Nil(t, err, "should succeed")
			assert.Equal(t, tC.beforeStart, beforeStart, "beforeStart should be equal")
			assert.Equal(t, tC.afterEnd, afterEnd, "afterEnd should be equal")
		})
	}
}

func TestSplitLinesToEditError(t *testing.T) {
	testCases := []struct {
		desc string
		src  string
	}{
		{
			desc: "without @generated_end",
			src: `TITLE:テストページ

//@generated_start
[[なんか]]
[[なんか2]]
//@gonorotod_end
`,
		},
		{
			desc: "without @generated_start",
			src: `TITLE:テストページ
//@gonototod_start
[[なんか]]
[[なんか2]]
//@generated_end
`,
		},
		{
			desc: "multiple start",
			src: `TITLE:テストページ
//@generated_start
[[なんか]]
//@generated_start
[[なんか2]]
//@generated_end
`,
		},
		{
			desc: "multiple end",
			src: `TITLE:テストページ
//@generated_start
[[なんか]]
//@generated_end
[[なんか2]]
//@generated_end
`,
		},
		{
			desc: "end before start",
			src: `TITLE:テストページ
[[なんか]]
//@generated_end
[[なんか2]]
//@generated_start`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, _, err := splitLinesToEdit(tC.src)
			assert.NotNil(t, err)
		})
	}
}

func TestGenerateCharaListPage(t *testing.T) {
	testCases := []struct {
		desc   string
		old    string
		charas []common.CharacterData
		expect string
	}{
		{
			desc: "normal",
			old: `
TITLE:テストページ
//@generated_start
hogehoge
//@generated_end
fugafuga
`,
			charas: sampleCharaData,
			expect: `
TITLE:テストページ
//@generated_start
#tablesort{{
|名前|地域|初出|h
|[[日本]]|東アジア|2026/06/07|
|[[韓国]]|東アジア|2026/06/07|
}}
//@generated_end
fugafuga
`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := generateCharaListPage(tC.old, tC.charas)
			require.Nil(t, err)
			assert.Equal(t, tC.expect, got)
		})
	}
}

func TestGenerateMenuBar(t *testing.T) {
	testCases := []struct {
		desc   string
		old    string
		charas []common.CharacterData
		expect string
	}{
		{
			desc: "normal",
			old: `TITLE:テストページ
//@generated_start
hogehoge
//@generated_end
fugafuga`,
			charas: multipleAreasCharaData,
			expect: `TITLE:テストページ
//@generated_start
** 東アジア
- [[日本]]
- [[韓国]]
** 東南アジア・南アジア
** 中東
** ヨーロッパ
- [[イタリア]]
** オセアニア
** 北米
- [[カナダ]]
** 中南米
** アフリカ
//@generated_end
fugafuga
`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := generateMenuBar(tC.old, tC.charas)
			require.Nil(t, err, "should succeed")
			assert.Equal(t, tC.expect, got)
		})
	}
}
