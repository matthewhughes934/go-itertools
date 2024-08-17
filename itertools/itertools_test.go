package itertools_test

import (
	"iter"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/matthewhughes934/go-itertools/itertools"
)

func TestMap(t *testing.T) {
	slice := []int{1, 2, 3}
	expected := []string{"1", "2", "3"}

	got := make([]string, 0, len(expected))
	iter := itertools.Map(strconv.Itoa, slices.Values(slice))

	for s := range iter {
		got = append(got, s)
	}

	require.Equal(t, expected, got)
}

func TestFilter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	filterFunc := func(i int) bool { return i%2 == 1 }
	expected := []int{1, 3, 5, 7, 9}

	got := make([]int, 0, len(expected))
	iter := itertools.Filter(filterFunc, slices.Values(slice))

	for n := range iter {
		got = append(got, n)
	}

	require.Equal(t, expected, got)
}

func TestEnumerate(t *testing.T) {
	slice := []string{"foo", "bar", "wat", "baz"}
	expectedIdx := []int{10, 11, 12, 13}
	expectedVals := []string{"foo", "bar", "wat", "baz"}

	gotIdx := make([]int, 0, len(expectedIdx))
	gotVals := make([]string, 0, len(expectedVals))
	iter := itertools.Enumerate(slices.Values(slice), 10)

	for i, v := range iter {
		gotIdx = append(gotIdx, i)
		gotVals = append(gotVals, v)
	}

	require.Equal(t, expectedIdx, gotIdx)
	require.Equal(t, expectedVals, gotVals)
}

func TestAnyFunc(t *testing.T) {
	data := []int{100, -1, 25, 13, 2, 4}

	for _, tc := range []struct {
		desc     string
		checker  func(int) bool
		expected bool
	}{
		{
			"always true",
			func(int) bool { return true },
			true,
		},
		{
			"never true",
			func(int) bool { return false },
			false,
		},
		{
			"too big",
			func(i int) bool { return i > 200 },
			false,
		},
		{
			"negative",
			func(i int) bool { return i < 0 },
			true,
		},
		{
			"positive",
			func(i int) bool { return i > 0 },
			true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			got := itertools.AnyFunc(tc.checker, slices.Values(data))

			require.Equal(t, tc.expected, got)
		})
	}
}

func TestAllFunc(t *testing.T) {
	data := []int{100, -1, 25, 13, 2, 4}

	for _, tc := range []struct {
		desc     string
		checker  func(int) bool
		expected bool
	}{
		{
			"always true",
			func(int) bool { return true },
			true,
		},
		{
			"never true",
			func(int) bool { return false },
			false,
		},
		{
			"fail at end",
			func(i int) bool { return i != 4 },
			false,
		},
		{
			"fail middle",
			func(i int) bool { return i != 25 },
			false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			got := itertools.AllFunc(tc.checker, slices.Values(data))

			require.Equal(t, tc.expected, got)
		})
	}
}

func TestZip(t *testing.T) {
	first := []int{1, 2, 3}
	second := []int{11, 12, 13}
	third := []int{21, 22, 23, 24}
	fillValue := 42
	expected := []int{1, 11, 21, 2, 12, 22, 3, 13, 23, 42, 42, 24}

	seqs := []iter.Seq[int]{
		slices.Values(first),
		slices.Values(second),
		slices.Values(third),
	}

	got := make([]int, 0, len(expected))
	iter := itertools.ZipLongest(fillValue, seqs...)

	for v := range iter {
		got = append(got, v)
	}

	require.Equal(t, expected, got)
}
