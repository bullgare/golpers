package testhelpers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type args[T any] struct {
	in T
}
type testCase[T any] struct {
	name string
	args args[T]
	want []T
}

func runTest[T any](t *testing.T, tc testCase[T], logSuffix string) {
	t.Run(tc.name, func(t *testing.T) {
		// regular
		got := doGenerate(tc.args.in)
		assert.Equal(t, tc.want, got, logSuffix+": regular")

		// pointer
		got2 := doGenerate(&tc.args.in)
		require.Equal(t, len(tc.want), len(got2), logSuffix+": pointer: lengths")
		for i, v := range got2 {
			require.Equal(t, reflect.ValueOf(v).Kind(), reflect.Ptr, logSuffix+": pointer %d is not a pointer", i)
			assert.Equal(t, tc.want[i], *v, logSuffix+": pointer, %d: %v != %v", i, v, tc.want[i])
		}
	})
}

func TestGenerateTestPermutations(t *testing.T) {
	type level1 struct {
		Field1 string
		Field2 []string
		Field3 []int
	}

	type top1 struct{ Field1 string }
	type top2 struct{ Field2 []string }
	type top3 struct{ Field3 []int }

	type level2 struct {
		Top1 top1
		Top2 top2
		Top3 top3
	}

	tests1 := []testCase[level1]{
		{
			name: "level1 field 2",
			args: args[level1]{
				in: level1{
					Field1: "field 1",
					Field2: []string{"one", "two", "three"},
					Field3: []int{1},
				},
			},
			want: []level1{
				{
					Field1: "field 1",
					Field2: []string{"one", "two", "three"},
					Field3: []int{1},
				},
				{
					Field1: "field 1",
					Field2: []string{"one", "three", "two"},
					Field3: []int{1},
				},
				{
					Field1: "field 1",
					Field2: []string{"two", "one", "three"},
					Field3: []int{1},
				},
				{
					Field1: "field 1",
					Field2: []string{"two", "three", "one"},
					Field3: []int{1},
				},
				{
					Field1: "field 1",
					Field2: []string{"three", "one", "two"},
					Field3: []int{1},
				},
				{
					Field1: "field 1",
					Field2: []string{"three", "two", "one"},
					Field3: []int{1},
				},
			},
		},
		{
			name: "level1 field 3",
			args: args[level1]{
				in: level1{
					Field1: "field 1",
					Field2: []string{"one"},
					Field3: []int{1, 2, 3},
				},
			},
			want: []level1{
				{
					Field1: "field 1",
					Field2: []string{"one"},
					Field3: []int{1, 2, 3},
				},
				{
					Field1: "field 1",
					Field2: []string{"one"},
					Field3: []int{1, 3, 2},
				},
				{
					Field1: "field 1",
					Field2: []string{"one"},
					Field3: []int{2, 1, 3},
				},
				{
					Field1: "field 1",
					Field2: []string{"one"},
					Field3: []int{2, 3, 1},
				},
				{
					Field1: "field 1",
					Field2: []string{"one"},
					Field3: []int{3, 1, 2},
				},
				{
					Field1: "field 1",
					Field2: []string{"one"},
					Field3: []int{3, 2, 1},
				},
			},
		},
		{
			name: "level1 complicated one - field 2 and field 3",
			args: args[level1]{
				in: level1{
					Field1: "field 1",
					Field2: []string{"one", "two"},
					Field3: []int{1, 2},
				},
			},
			want: []level1{
				{
					Field1: "field 1",
					Field2: []string{"one", "two"},
					Field3: []int{1, 2},
				},
				{
					Field1: "field 1",
					Field2: []string{"two", "one"},
					Field3: []int{1, 2},
				},
				{
					Field1: "field 1",
					Field2: []string{"one", "two"},
					Field3: []int{2, 1},
				},
				{
					Field1: "field 1",
					Field2: []string{"two", "one"},
					Field3: []int{2, 1},
				},
			},
		},
		{
			name: "level1 empty slices are skipped",
			args: args[level1]{
				in: level1{
					Field1: "field 1",
					Field2: []string{},
					Field3: []int{},
				},
			},
			want: []level1{
				{
					Field1: "field 1",
					Field2: []string{},
					Field3: []int{},
				},
			},
		},
	}

	tests2 := []testCase[level2]{
		{
			name: "level2 field 2",
			args: args[level2]{
				in: level2{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one", "two", "three"}},
					Top3: top3{Field3: []int{1}},
				},
			},
			want: []level2{
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one", "two", "three"}},
					Top3: top3{Field3: []int{1}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one", "three", "two"}},
					Top3: top3{Field3: []int{1}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"two", "one", "three"}},
					Top3: top3{Field3: []int{1}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"two", "three", "one"}},
					Top3: top3{Field3: []int{1}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"three", "one", "two"}},
					Top3: top3{Field3: []int{1}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"three", "two", "one"}},
					Top3: top3{Field3: []int{1}},
				},
			},
		},
		{
			name: "level2 field 3",
			args: args[level2]{
				in: level2{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one"}},
					Top3: top3{Field3: []int{1, 2, 3}},
				},
			},
			want: []level2{
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one"}},
					Top3: top3{Field3: []int{1, 2, 3}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one"}},
					Top3: top3{Field3: []int{1, 3, 2}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one"}},
					Top3: top3{Field3: []int{2, 1, 3}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one"}},
					Top3: top3{Field3: []int{2, 3, 1}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one"}},
					Top3: top3{Field3: []int{3, 1, 2}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one"}},
					Top3: top3{Field3: []int{3, 2, 1}},
				},
			},
		},
		{
			name: "level2 complicated one - field 2 and field 3",
			args: args[level2]{
				in: level2{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one", "two"}},
					Top3: top3{Field3: []int{1, 2}},
				},
			},
			want: []level2{
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one", "two"}},
					Top3: top3{Field3: []int{1, 2}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"two", "one"}},
					Top3: top3{Field3: []int{1, 2}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"one", "two"}},
					Top3: top3{Field3: []int{2, 1}},
				},
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{"two", "one"}},
					Top3: top3{Field3: []int{2, 1}},
				},
			},
		},
		{
			name: "level2 empty slices are skipped",
			args: args[level2]{
				in: level2{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{}},
					Top3: top3{Field3: []int{}},
				},
			},
			want: []level2{
				{
					Top1: top1{Field1: "field 1"},
					Top2: top2{Field2: []string{}},
					Top3: top3{Field3: []int{}},
				},
			},
		},
	}

	tests3 := []testCase[string]{
		{
			name: "string",
			args: args[string]{
				in: "str",
			},
			want: []string{
				"str",
			},
		},
	}

	for _, tt := range tests1 {
		runTest(t, tt, "level 1")
	}
	for _, tt := range tests2 {
		runTest(t, tt, "level 2")
	}
	for _, tt := range tests3 {
		runTest(t, tt, "regular string")
	}
}
