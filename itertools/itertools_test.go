package itertools_test

import (
	"context"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/matthewhughes934/go-itertools/itertools"
)

func collectPairs[K comparable](seq iter.Seq2[K, K]) [][]K {
	var res [][]K //nolint:prealloc
	for k1, k2 := range seq {
		res = append(res, []K{k1, k2})
	}
	return res
}

func TestSlice(t *testing.T) {
	data := slices.Collect(itertools.RangeUntil(10, 1))

	for _, tc := range []struct {
		start    int
		end      int
		step     int
		expected []int
	}{
		{
			0,
			5,
			1,
			[]int{0, 1, 2, 3, 4},
		},
		{
			0,
			10,
			2,
			[]int{0, 2, 4, 6, 8},
		},
		{
			10,
			5,
			1,
			nil,
		},
		{
			0,
			len(data),
			1,
			data,
		},
		{
			len(data) + 1,
			len(data) + 2,
			1,
			nil,
		},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			seq := itertools.Slice(slices.Values(data), tc.start, tc.end, tc.step)

			got := slices.Collect(seq)

			require.Equal(t, tc.expected, got)
		})
	}
}

func TestSlice_panicsOnBadStep(t *testing.T) {
	require.PanicsWithValue(
		t,
		"step for Slice must be a positive integer",
		func() { itertools.Slice(slices.Values([]int{}), 0, 0, -1) },
	)
}

func TestSlice_earlyStop(t *testing.T) {
	data := slices.Collect(itertools.RangeUntil(20, 1))
	takeUntil := 4
	expected := data[:takeUntil]

	seq := itertools.Slice(slices.Values(data), 0, len(data), 1)
	// an iterator that stops consuming early
	takeIter := func(yield func(int) bool) {
		next, stop := iter.Pull(seq)
		defer stop()

		for range takeUntil {
			v, _ := next()
			yield(v)
		}
	}

	got := slices.Collect(takeIter)

	require.Equal(t, expected, got)
}

func TestSlice2(t *testing.T) {
	dataLen := 10
	data := itertools.ZipPair(
		itertools.Range(0, dataLen, 1),
		itertools.Range(1, dataLen+1, 1),
	)

	for _, tc := range []struct {
		start    int
		end      int
		step     int
		expected [][]int
	}{
		{
			0,
			5,
			1,
			[][]int{
				{0, 1},
				{1, 2},
				{2, 3},
				{3, 4},
				{4, 5},
			},
		},
		{
			0,
			10,
			2,
			[][]int{
				{0, 1},
				{2, 3},
				{4, 5},
				{6, 7},
				{8, 9},
			},
		},
		{
			10,
			0,
			1,
			nil,
		},
		{
			0,
			dataLen,
			1,
			[][]int{
				{0, 1},
				{1, 2},
				{2, 3},
				{3, 4},
				{4, 5},
				{5, 6},
				{6, 7},
				{7, 8},
				{8, 9},
				{9, 10},
			},
		},
		{
			dataLen + 1,
			dataLen + 2,
			1,
			nil,
		},
		{
			dataLen - 1,
			2 * dataLen,
			1,
			[][]int{{9, 10}},
		},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			seq := itertools.Slice2(data, tc.start, tc.end, tc.step)

			got := collectPairs(seq)

			require.Equal(t, tc.expected, got)
		})
	}
}

func TestSlice2_panicsOnBadStep(t *testing.T) {
	require.PanicsWithValue(
		t,
		"step for Slice2 must be a positive integer",
		func() { itertools.Slice2(slices.All([]int{}), 0, 0, -1) },
	)
}

func TestSlice2_earlyStop(t *testing.T) {
	data := slices.Collect(itertools.RangeUntil(10, 1))
	inSeq := slices.All(data)
	takeUntil := 4
	expected := map[int]int{0: 0, 1: 1, 2: 2, 3: 3}

	seq := itertools.Slice2(inSeq, 0, len(data), 1)
	// an iterator that stops consuming early
	takeIter := func(yield func(int, int) bool) {
		next, stop := iter.Pull2(seq)
		defer stop()

		for range takeUntil {
			k, v, _ := next()
			yield(k, v)
		}
	}

	got := maps.Collect(takeIter)
	require.Equal(t, expected, got)
}

func TestRange(t *testing.T) {
	for _, tc := range []struct {
		start    int
		end      int
		step     int
		expected []int
	}{
		{
			0,
			0,
			1,
			nil,
		},
		{
			10,
			0,
			1,
			nil,
		},
		{
			0,
			5,
			1,
			[]int{0, 1, 2, 3, 4},
		},
		{
			4,
			-1,
			-1,
			[]int{4, 3, 2, 1, 0},
		},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			seq := itertools.Range(tc.start, tc.end, tc.step)

			require.Equal(t, tc.expected, slices.Collect(seq))
		})
	}
}

func TestRange_earlyStop(t *testing.T) {
	seq := itertools.Range(0, 5, 1)
	takeLen := 3
	expected := []int{0, 1, 2}

	got := slices.Collect(itertools.SliceUntil(seq, takeLen, 1))

	require.Equal(t, expected, got)
}

func TestRange_panicsOnZeroStep(t *testing.T) {
	require.PanicsWithValue(
		t,
		"step for Range must be non-zero",
		func() { itertools.Range(0, 10, 0) },
	)
}

func TestMap(t *testing.T) {
	slice := []int{1, 2, 3}
	expected := []string{"1", "2", "3"}

	iter := itertools.Map(strconv.Itoa, slices.Values(slice))

	require.Equal(t, expected, slices.Collect(iter))
}

func TestMap_earlyStop(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	takeLen := 3
	expected := []string{"1", "2", "3"}

	seq := itertools.Map(strconv.Itoa, slices.Values(slice))
	got := slices.Collect(itertools.SliceUntil(seq, takeLen, 1))

	require.Equal(t, expected, got)
}

func TestMap2_earlyStop(t *testing.T) {
	inSeq := itertools.Enumerate(itertools.RangeUntil(5, 1), 1)
	takeLen := 3
	mapFunc := func(x1 int, x2 int) (int, int) { return x1 * 2, x2 * 2 }
	expected := [][]int{{2, 0}, {4, 2}, {6, 4}}

	seq := itertools.Map2(mapFunc, inSeq)
	shortSeq := itertools.SliceUntil2(seq, takeLen, 1)
	got := collectPairs(shortSeq)

	require.Equal(t, expected, got)
}

func testFilter[V any](t *testing.T, data []V, filterFunc func(V) bool, expected []V) {
	t.Helper()
	seq := itertools.Filter(filterFunc, slices.Values(data))
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestFilter(t *testing.T) {
	t.Run("basic int filter", func(t *testing.T) {
		testFilter(
			t,
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			func(i int) bool { return i%2 == 1 },
			[]int{1, 3, 5, 7, 9},
		)
	})

	t.Run("basic string filter", func(t *testing.T) {
		testFilter(
			t,
			[]string{"foo", "skip", "bar", "wat", "skip"},
			func(s string) bool { return s != "skip" },
			[]string{"foo", "bar", "wat"},
		)
	})

	t.Run("trivial true", func(t *testing.T) {
		testFilter(
			t,
			[]int{1, 2, 3, 4, 5},
			func(int) bool { return true },
			[]int{1, 2, 3, 4, 5},
		)
	})
}

func TestFilter_earlyStop(t *testing.T) {
	data := itertools.RangeUntil(5, 1)
	filterFunc := func(int) bool { return true }
	takeLen := 3
	expected := []int{0, 1, 2}

	iter := itertools.Filter(filterFunc, data)
	got := slices.Collect(itertools.SliceUntil(iter, takeLen, 1))

	require.Equal(t, expected, got)
}

func TestFilter2_earlyStop(t *testing.T) {
	inSeq := itertools.Enumerate(itertools.RangeUntil(5, 1), 1)
	filterFunc := func(int, int) bool { return true }
	takeLen := 3
	expected := [][]int{{1, 0}, {2, 1}, {3, 2}}

	seq := itertools.Filter2(filterFunc, inSeq)
	shortSeq := itertools.SliceUntil2(seq, takeLen, 1)
	got := collectPairs(shortSeq)

	require.Equal(t, expected, got)
}

func TestKeys(t *testing.T) {
	inSeq := maps.All(map[string]int{"one": 1, "two": 2, "three": 3})
	expected := []string{"one", "two", "three"}

	seq := itertools.Keys(inSeq)
	got := slices.Collect(seq)

	// ordering of map iteration is not deterministic
	require.ElementsMatch(t, expected, got)
}

func TestKeys_earlyStop(t *testing.T) {
	inSeq := itertools.Enumerate(itertools.RangeUntil(5, 1), 1)
	takeLen := 3
	expected := []int{1, 2, 3}

	seq := itertools.Keys(inSeq)
	got := slices.Collect(itertools.SliceUntil(seq, takeLen, 1))

	require.Equal(t, expected, got)
}

func TestValues(t *testing.T) {
	inSeq := maps.All(map[string]int{"one": 1, "two": 2, "three": 3})
	expected := []int{1, 2, 3}

	seq := itertools.Values(inSeq)
	got := slices.Collect(seq)

	// ordering of map iteration is not deterministic
	require.ElementsMatch(t, expected, got)
}

func TestValues_earlyStop(t *testing.T) {
	inSeq := itertools.Enumerate(itertools.RangeUntil(5, 1), 1)
	takeLen := 3
	expected := []int{0, 1, 2}

	seq := itertools.Values(inSeq)
	got := slices.Collect(itertools.SliceUntil(seq, takeLen, 1))

	require.Equal(t, expected, got)
}

func TestEnumerate(t *testing.T) {
	slice := []string{"foo", "bar", "wat", "baz"}
	expectedIdx := []int{10, 11, 12, 13}
	expectedVals := []string{"foo", "bar", "wat", "baz"}

	seq := itertools.Enumerate(slices.Values(slice), 10)
	gotIdx := slices.Collect(itertools.Keys(seq))
	gotVals := slices.Collect(itertools.Values(seq))

	require.Equal(t, expectedIdx, gotIdx)
	require.Equal(t, expectedVals, gotVals)
}

func TestEnumerate_earlyStop(t *testing.T) {
	inSeq := itertools.RangeUntil(10, 1)
	start := 5
	takeLen := 3
	expectedIdx := []int{5, 6, 7}
	expectedVals := []int{0, 1, 2}

	seq := itertools.SliceUntil2(itertools.Enumerate(inSeq, start), takeLen, 1)
	gotIdx := slices.Collect(itertools.Keys(seq))
	gotVals := slices.Collect(itertools.Values(seq))

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

func TestChain_earlyStop(t *testing.T) {
	firstSeq := itertools.RangeUntil(5, 1)
	secondSeq := itertools.Range(5, 10, 1)
	takeLen := 6
	expected := []int{0, 1, 2, 3, 4, 5}

	seq := itertools.Chain(firstSeq, secondSeq)
	got := slices.Collect(itertools.SliceUntil(seq, takeLen, 1))

	require.Equal(t, expected, got)
}

func TestChain2_earlyStop(t *testing.T) {
	firstSeq := itertools.Enumerate(itertools.RangeUntil(5, 1), 1)
	secondSeq := itertools.Enumerate(itertools.Range(5, 10, 1), 6)
	takeLen := 6
	expected := [][]int{{1, 0}, {2, 1}, {3, 2}, {4, 3}, {5, 4}, {6, 5}}

	seq := itertools.Chain2(firstSeq, secondSeq)
	got := collectPairs(itertools.SliceUntil2(seq, takeLen, 1))

	require.Equal(t, expected, got)
}

func TestZipPair_earlyStop(t *testing.T) {
	firstSeq := itertools.RangeUntil(5, 1)
	secondSeq := itertools.Range(5, 10, 1)
	takeLen := 3
	expected := [][]int{{0, 5}, {1, 6}, {2, 7}}

	t.Run("first sequence stopped", func(t *testing.T) {
		firstSeqShort := itertools.SliceUntil(firstSeq, takeLen, 1)

		seq := itertools.ZipPair(firstSeqShort, secondSeq)
		got := collectPairs(seq)

		require.Equal(t, expected, got)
	})

	t.Run("second sequence stopped", func(t *testing.T) {
		secondSeqShort := itertools.SliceUntil(secondSeq, takeLen, 1)

		seq := itertools.ZipPair(firstSeq, secondSeqShort)
		got := collectPairs(seq)

		require.Equal(t, expected, got)
	})
}

func TestZipLongest(t *testing.T) {
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

	seq := itertools.ZipLongest(fillValue, seqs...)
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestZipLongest_earlyStop(t *testing.T) {
	first := []int{1, 2, 3}
	second := []int{11, 12, 13, 14}
	fillValue := -1
	takeLen := 4
	expected := []int{1, 11, 2, 12}

	seqs := []iter.Seq[int]{slices.Values(first), slices.Values(second)}

	seq := itertools.ZipLongest(fillValue, seqs...)
	shortSeq := itertools.SliceUntil(seq, takeLen, 1)
	got := slices.Collect(shortSeq)

	require.Equal(t, expected, got)
}

func TestRangeFrom(t *testing.T) {
	start := 10
	step := 2
	takeLen := 5
	expected := []int{10, 12, 14, 16, 18}

	seq := itertools.RangeFrom(start, step)
	shortSeq := itertools.SliceUntil(seq, takeLen, 1)
	got := slices.Collect(shortSeq)

	require.Equal(t, expected, got)
}

func TestCycle_earlyExit(t *testing.T) {
	baseSeq := itertools.RangeUntil(5, 1)
	takeLen := 3
	expected := []int{0, 1, 2}

	seq := itertools.SliceUntil(itertools.Cycle(baseSeq), takeLen, 1)
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestCycle2_earlyExit(t *testing.T) {
	baseSeq := itertools.Enumerate(itertools.RangeUntil(5, 1), 1)
	takeLen := 3
	expected := [][]int{{1, 0}, {2, 1}, {3, 2}}

	seq := itertools.SliceUntil2(itertools.Cycle2(baseSeq), takeLen, 1)
	got := collectPairs(seq)

	require.Equal(t, expected, got)
}

func TestRepeat_earlyExit(t *testing.T) {
	baseSeq := itertools.Repeat("A", 100)
	takeLen := 5
	expected := []string{"A", "A", "A", "A", "A"}

	seq := itertools.SliceUntil(baseSeq, takeLen, 1)
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestRepeat_infinite(t *testing.T) {
	baseSeq := itertools.Repeat("A", -1)
	takeLen := 5
	expected := []string{"A", "A", "A", "A", "A"}

	seq := itertools.SliceUntil(baseSeq, takeLen, 1)
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestAccumulate_earlyExit(t *testing.T) {
	accumulator := func(x1 int, x2 int) int { return x1 * x2 }
	baseSeq := itertools.Accumulate(itertools.Range(1, 10, 1), accumulator, 1)
	takeLen := 5
	expected := []int{1, 2, 6, 24, 120}

	seq := itertools.SliceUntil(baseSeq, takeLen, 1)
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestCompress_earlyExit(t *testing.T) {
	baseSeq := itertools.RangeUntil(10, 1)
	selectors := slices.Values(
		[]bool{true, false, false, true, true, false, true, true, false, true},
	)
	takeLen := 5
	expected := []int{0, 3, 4, 6, 7}

	seq := itertools.SliceUntil(itertools.Compress(baseSeq, selectors), takeLen, 1)
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestCompress2_earlyExit(t *testing.T) {
	baseSeq := itertools.Enumerate(itertools.RangeUntil(10, 1), 1)
	selectors := slices.Values(
		[]bool{true, false, false, true, true, false, true, true, false, true},
	)
	takeLen := 5
	expected := [][]int{{1, 0}, {4, 3}, {5, 4}, {7, 6}, {8, 7}}

	seq := itertools.SliceUntil2(itertools.Compress2(baseSeq, selectors), takeLen, 1)
	got := collectPairs(seq)

	require.Equal(t, expected, got)
}

func TestCompress2_shortSelector(t *testing.T) {
	baseSeq := itertools.Enumerate(itertools.RangeUntil(10, 1), 1)
	selectors := slices.Values([]bool{true, true, true})
	expected := [][]int{{1, 0}, {2, 1}, {3, 2}}

	seq := itertools.Compress2(baseSeq, selectors)
	got := collectPairs(seq)

	require.Equal(t, expected, got)
}

func TestDropWhile_earlyExit(t *testing.T) {
	baseSeq := itertools.RangeUntil(10, 1)
	predicate := func(x int) bool { return x <= 3 }
	takeLen := 3
	expected := []int{4, 5, 6}

	seq := itertools.SliceUntil(itertools.DropWhile(baseSeq, predicate), takeLen, 1)
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestDropWhile2_earlyExit(t *testing.T) {
	baseSeq := itertools.Enumerate(itertools.Range(1, 10, 1), 0)
	predicate := func(i int, _ int) bool { return i <= 3 }
	takeLen := 3
	expected := [][]int{{4, 5}, {5, 6}, {6, 7}}

	seq := itertools.SliceUntil2(itertools.DropWhile2(baseSeq, predicate), takeLen, 1)
	got := collectPairs(seq)

	require.Equal(t, expected, got)
}

func TestTakeWhile_earlyExit(t *testing.T) {
	baseSeq := itertools.RangeUntil(10, 1)
	predicate := func(x int) bool { return x < 7 }
	takeLen := 3
	expected := []int{0, 1, 2}

	seq := itertools.SliceUntil(itertools.TakeWhile(baseSeq, predicate), takeLen, 1)
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestTakeWhile2_earlyExit(t *testing.T) {
	baseSeq := itertools.Enumerate(itertools.RangeUntil(10, 1), 1)
	predicate := func(x int, _ int) bool { return x < 7 }
	takeLen := 3
	expected := [][]int{{1, 0}, {2, 1}, {3, 2}}

	seq := itertools.SliceUntil2(itertools.TakeWhile2(baseSeq, predicate), takeLen, 1)
	got := collectPairs(seq)

	require.Equal(t, expected, got)
}

func TestIterCtx_earlyExit(t *testing.T) {
	baseSeq := itertools.RangeUntil(10, 1)
	takeLen := 3
	expected := []int{0, 1, 2}

	seq := itertools.SliceUntil(itertools.IterCtx(context.Background(), baseSeq), takeLen, 1)
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestIterCtx_inputExhausted(t *testing.T) {
	baseSeq := itertools.RangeUntil(3, 1)
	expected := []int{0, 1, 2}

	seq := itertools.IterCtx(context.Background(), baseSeq)
	got := slices.Collect(seq)

	require.Equal(t, expected, got)
}

func TestIterCtx2_earlyExit(t *testing.T) {
	baseSeq := itertools.Enumerate(itertools.RangeUntil(10, 1), 1)
	takeLen := 3
	expected := [][]int{{1, 0}, {2, 1}, {3, 2}}

	seq := itertools.SliceUntil2(itertools.IterCtx2(context.Background(), baseSeq), takeLen, 1)
	got := collectPairs(seq)

	require.Equal(t, expected, got)
}

func TestIterCtx2_inputExhausted(t *testing.T) {
	baseSeq := itertools.Enumerate(itertools.RangeUntil(3, 1), 1)
	expected := [][]int{{1, 0}, {2, 1}, {3, 2}}

	seq := itertools.IterCtx2(context.Background(), baseSeq)
	got := collectPairs(seq)

	require.Equal(t, expected, got)
}

func TestPairwise_emptyIfFewerThanTwo(t *testing.T) {
	for _, vals := range [][]int{{}, {1}} {
		t.Run(fmt.Sprintf("%d", len(vals)), func(t *testing.T) {
			iter := itertools.Pairwise(slices.Values(vals))
			require.Empty(t, maps.Collect(iter))
		})
	}
}
