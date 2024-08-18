package itertools_test

import (
	"context"
	"fmt"
	"iter"
	"maps"
	"slices"
	"time"

	"github.com/matthewhughes934/go-itertools/itertools"
)

func ExampleChain() {
	seqs := []iter.Seq[int]{
		slices.Values([]int{1, 2, 3}),
		slices.Values([]int{4, 5, 6}),
	}

	for n := range itertools.Chain(seqs...) {
		fmt.Println(n)
	}

	// output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
}

func ExampleChain2() {
	seqs := []iter.Seq2[string, int]{
		maps.All(map[string]int{"foo": 1, "bar": 2}),
		maps.All(map[string]int{"baz": 3, "wat": 4}),
	}

	res := maps.Collect(itertools.Chain2(seqs...))
	for k, v := range res {
		fmt.Println(k, v)
	}

	// unordered output:
	// foo 1
	// bar 2
	// baz 3
	// wat 4
}

func ExampleMap() {
	seq := slices.Values([]int{1, 2, 3})

	res := itertools.Map(func(i int) int { return i * 2 }, seq)

	for n := range res {
		fmt.Println(n)
	}

	// output:
	// 2
	// 4
	// 6
}

func ExampleMap2() {
	seq := maps.All(
		map[string]int{"foo": 1, "bar": 2, "baz": 3},
	)

	res := maps.Collect(
		itertools.Map2(
			func(k string, v int) (string, float64) {
				return k + " halved", float64(v) / 2.0
			},
			seq,
		),
	)

	for k, v := range res {
		fmt.Println(k, v)
	}

	// unordered output:
	// foo halved 0.5
	// bar halved 1
	// baz halved 1.5
}

func ExampleFilter() {
	seq := slices.Values([]int{1, 2, 3, 4, 5})

	res := itertools.Filter(func(i int) bool { return i%2 == 1 }, seq)

	for n := range res {
		fmt.Println(n)
	}

	// output:
	// 1
	// 3
	// 5
}

func ExampleFilter2() {
	seq := maps.All(
		map[string]int{"foo": 1, "bar": 2, "buz": 3, "baz": 4, "bux": 5},
	)

	res := maps.Collect(
		itertools.Filter2(
			func(k string, v int) bool {
				return k[0] == 'b' && v%2 == 0
			},
			seq,
		),
	)

	for k, v := range res {
		fmt.Println(k, v)
	}

	// unordered output:
	// bar 2
	// baz 4
}

func ExampleEnumerate() {
	seq := itertools.Chain(
		slices.Values([]int{10, 11, 12, 13}),
		slices.Values([]int{14, 15, 16, 17}),
	)

	res := itertools.Enumerate(seq, 2)

	for i, n := range res {
		fmt.Println(i, n)
	}

	// output:
	// 2 10
	// 3 11
	// 4 12
	// 5 13
	// 6 14
	// 7 15
	// 8 16
	// 9 17
}

func isEven(i int) bool { return i%2 == 0 }
func isOdd(i int) bool  { return i%2 == 1 }

func ExampleAnyFunc() {
	evenAndOdd := []int{1, 2, 3, 4, 5}
	onlyOdd := []int{1, 3, 5, 7}
	onlyEven := []int{2, 4, 6, 8}

	fmt.Println(itertools.AnyFunc(isEven, slices.Values(evenAndOdd)))
	fmt.Println(itertools.AnyFunc(isEven, slices.Values(onlyOdd)))
	fmt.Println(itertools.AnyFunc(isEven, slices.Values(onlyEven)))

	fmt.Println(itertools.AnyFunc(isOdd, slices.Values(evenAndOdd)))
	fmt.Println(itertools.AnyFunc(isOdd, slices.Values(onlyOdd)))
	fmt.Println(itertools.AnyFunc(isOdd, slices.Values(onlyEven)))

	// output:
	// true
	// false
	// true
	// true
	// true
	// false
}

func ExampleAnyFunc2() {
	m := map[string]int{"foo": 1, "bar": 2, "buz": 3, "bux": 4}

	fmt.Println(itertools.AnyFunc2(func(k string, _ int) bool { return k[0] == 'b' }, maps.All(m)))
	fmt.Println(itertools.AnyFunc2(func(k string, _ int) bool { return k[0] == 'z' }, maps.All(m)))
	fmt.Println(itertools.AnyFunc2(func(_ string, i int) bool { return i > 10 }, maps.All(m)))
	fmt.Println(
		itertools.AnyFunc2(func(k string, i int) bool { return k == "foo" && i == 1 }, maps.All(m)),
	)

	// output:
	// true
	// false
	// false
	// true
}

func ExampleZip() {
	seqs := []iter.Seq[int]{
		slices.Values([]int{1, 2, 3}),
		slices.Values([]int{11, 12, 13}),
	}

	res := itertools.Zip(seqs...)

	for n := range res {
		fmt.Println(n)
	}

	// output:
	// 1
	// 11
	// 2
	// 12
	// 3
	// 13
}

func ExampleZipPair() {
	names := slices.Values([]string{"one", "two", "three"})
	values := slices.Values([]int{1, 2, 3})

	for n, v := range itertools.ZipPair(names, values) {
		fmt.Println(n, v)
	}

	// output:
	// one 1
	// two 2
	// three 3
}

func ExampleZip2() {
	seqs := []iter.Seq2[int, int]{
		slices.All([]int{1, 2, 3}),
		slices.All([]int{11, 12, 13}),
	}

	res := itertools.Zip2(seqs...)

	for i, n := range res {
		fmt.Println(i, n)
	}

	// output:
	// 0 1
	// 0 11
	// 1 2
	// 1 12
	// 2 3
	// 2 13
}

func ExampleZipLongest() {
	seqs := []iter.Seq[int]{
		slices.Values([]int{1, 2, 3, 4, 5, 6}),
		slices.Values([]int{11, 12, 13}),
	}

	res := itertools.ZipLongest(-1, seqs...)

	for n := range res {
		fmt.Println(n)
	}

	// output:
	// 1
	// 11
	// 2
	// 12
	// 3
	// 13
	// 4
	// -1
	// 5
	// -1
	// 6
	// -1
}

func ExampleCycle() {
	res := itertools.Cycle(slices.Values([]string{"A", "B", "C", "D"}))

	for i, s := range itertools.Enumerate(res, 1) {
		fmt.Println(s)
		if i == 12 {
			break
		}
	}

	// output:
	// A
	// B
	// C
	// D
	// A
	// B
	// C
	// D
	// A
	// B
	// C
	// D
}

func ExampleRepeat() {
	res := itertools.Repeat("A", 7)

	for s := range res {
		fmt.Println(s)
	}

	// output:
	// A
	// A
	// A
	// A
	// A
	// A
	// A
}

func ExampleAllFunc() {
	seq := slices.Values([]int{1, 2, 3, 4, 5})

	fmt.Println(itertools.AllFunc(func(i int) bool { return i > 0 }, seq))
	fmt.Println(itertools.AllFunc(func(i int) bool { return i < 0 }, seq))
	fmt.Println(itertools.AllFunc(func(i int) bool { return i < 4 }, seq))

	// output:
	// true
	// false
	// false
}

func ExampleAllFunc2() {
	seq := maps.All(map[string]int{"a": 1, "b": 2, "c": 3})

	fmt.Println(itertools.AllFunc2(func(k string, v int) bool { return k < "d" && v < 4 }, seq))
	fmt.Println(itertools.AllFunc2(func(k string, v int) bool { return k > "d" || v < 4 }, seq))
	fmt.Println(itertools.AllFunc2(func(k string, v int) bool { return k > "d" && v > 4 }, seq))

	// output:
	// true
	// true
	// false
}

func ExampleAccumulate() {
	seq := slices.Values([]int{1, 2, 3, 4, 5})
	add := func(x int, y int) int { return x + y }

	for v := range itertools.Accumulate(seq, add, 0) {
		fmt.Println(v)
	}

	// output:
	// 1
	// 3
	// 6
	// 10
	// 15
}

func ExampleCompress() {
	seq := slices.Values([]string{"A", "B", "C", "D", "E"})
	selectors := slices.Values([]bool{true, true, false, true, false})

	for s := range itertools.Compress(seq, selectors) {
		fmt.Println(s)
	}

	// output:
	// A
	// B
	// D
}

func ExampleDropWhile() {
	seq := slices.Values([]int{1, 4, 6, 3, 8})
	predicate := func(i int) bool { return i < 5 }

	for n := range itertools.DropWhile(seq, predicate) {
		fmt.Println(n)
	}

	// output:
	// 6
	// 3
	// 8
}

func ExampleTakeWhile() {
	seq := slices.Values([]int{1, 4, 6, 3, 8})
	predicate := func(i int) bool { return i < 5 }

	for n := range itertools.TakeWhile(seq, predicate) {
		fmt.Println(n)
	}

	// output:
	// 1
	// 4
}

func ExampleCollectIntoSlice() {
	data := []int{1, 2, 3, 4}
	mapper := func(x int) int { return 2 * x }
	dest := make([]int, len(data))

	itertools.CollectIntoSlice(itertools.Map(mapper, slices.Values(data)), dest)

	for _, n := range dest {
		fmt.Println(n)
	}

	// output:
	// 2
	// 4
	// 6
	// 8
}

func ExampleCollectIntoMap() {
	data := map[string]int{"foo": 1, "bar": 2}
	mapper := func(k string, x int) (string, int) { return k + "_new", x + 2 }
	dest := make(map[string]int, len(data))

	itertools.CollectIntoMap(itertools.Map2(mapper, maps.All(data)), dest)

	for k, v := range dest {
		fmt.Println(k, v)
	}

	// unordered output:
	// foo_new 3
	// bar_new 4
}

func ExampleCollectIntoMap_concatMaps() {
	map1 := map[string]int{"A": 1, "B": 2, "C": 3}
	map2 := map[string]int{"D": 4, "E": 5, "F": 6}
	dest := make(map[string]int, len(map1)+len(map2))

	itertools.CollectIntoMap(itertools.Chain2(maps.All(map1), maps.All(map2)), dest)

	for k, v := range dest {
		fmt.Println(k, v)
	}

	// unordered output:
	// A 1
	// B 2
	// C 3
	// D 4
	// E 5
	// F 6
}

func ExampleIterCtx() {
	tickIter := func(yield func(int) bool) {
		ticker := time.NewTicker(time.Millisecond * 10)
		defer ticker.Stop()

		count := 1
		for {
			<-ticker.C
			if !yield(count) {
				return
			}
			count++
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*55)
	defer cancel()

	for c := range itertools.IterCtx(ctx, tickIter) {
		fmt.Println(c)
	}

	// output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleIterCtx2() {
	tickIter := func(yield func(int) bool) {
		ticker := time.NewTicker(time.Millisecond * 10)
		defer ticker.Stop()

		count := 1
		for {
			<-ticker.C
			if !yield(count) {
				return
			}
			count++
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*55)
	defer cancel()

	for c := range itertools.IterCtx(ctx, tickIter) {
		fmt.Println(c)
	}

	// output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleSlice() {
	seq := slices.Values([]string{"A", "B", "C", "D", "E", "F", "G", "H"})

	for s := range itertools.Slice(seq, 1, 9, 2) {
		fmt.Println(s)
	}

	// output:
	// B
	// D
	// F
	// H
}

func ExampleSliceUntil() {
	seq := slices.Values([]string{"A", "B", "C", "D", "E", "F", "G", "H"})

	for s := range itertools.SliceUntil(seq, 5, 2) {
		fmt.Println(s)
	}

	// output:
	// A
	// C
	// E
}

func ExampleSlice2() {
	seq := itertools.ZipPair(
		slices.Values([]int{1, 2, 3, 4, 5}),
		slices.Values([]int{11, 12, 13, 14, 15}),
	)

	for x1, x2 := range itertools.Slice2(seq, 1, 5, 2) {
		fmt.Println(x1, x2)
	}

	// output:
	// 2 12
	// 4 14
}

func ExampleSliceUntil2() {
	seq := itertools.ZipPair(
		slices.Values([]int{1, 2, 3, 4, 5}),
		slices.Values([]int{11, 12, 13, 14, 15}),
	)

	for x1, x2 := range itertools.SliceUntil2(seq, 4, 2) {
		fmt.Println(x1, x2)
	}

	// output:
	// 1 11
	// 3 13
}

func ExampleRange() {
	for n := range itertools.Range(2, 10, 2) {
		fmt.Println(n)
	}

	// output:
	// 2
	// 4
	// 6
	// 8
}

func ExampleRangeUntil() {
	for n := range itertools.RangeUntil(10, 3) {
		fmt.Println(n)
	}

	// output:
	// 0
	// 3
	// 6
	// 9
}
