// Package itertools provides tools or working with sequences as defined in
// [iter].
// It is inspired by the [Python itertools package]
//
// [Python itertools package]: https://docs.python.org/3/library/itertools.html
package itertools

import (
	"context"
	"iter"
	"maps"
	"slices"
	"sync"
)

// Chain returns a [iter.Seq] that returns elements from the first sequence
// until it is exhausted, then proceeds to the next sequence, until all
// sequences are exhausted.
func Chain[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// ChainSlices is a convenience method for calling [Chain] given slices.
//
// E.g instead of something like:
//
//	var res []V
//	res = append(res, firstSlice...)
//	res = append(res, secondSlie...)
//	// more slices ...
//
//	for _, v := range res {
//		// do stuff with v ...
//	}
//
// You can call:
//
//	for v := range ChainSlices(firstSlice, secondSlice, /* more slices ... */) {
//		// do stuff with v ...
//	}
func ChainSlices[V any](sls ...[]V) iter.Seq[V] {
	seqs := make([]iter.Seq[V], 0, len(sls))
	for _, slice := range sls {
		seqs = append(seqs, slices.Values(slice))
	}
	return Chain(seqs...)
}

// Chain2 returns a [iter.Seq2] similar to [Chain], returning a [iter.Seq2]
// that runs over all the provider iterators.
func Chain2[K comparable, V any](iters ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, iter := range iters {
			for k, v := range iter {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// ChainMaps is like [Chain] but for maps.
func ChainMaps[K comparable, V any](mps ...map[K]V) iter.Seq2[K, V] {
	seqs := make([]iter.Seq2[K, V], 0, len(mps))
	for _, mp := range mps {
		seqs = append(seqs, maps.All(mp))
	}
	return Chain2(seqs...)
}

// Map returns a [iter.Seq] that applies mapFunc to every item of iterable,
// yielding the results.
func Map[V1 any, V2 any](mapFunc func(V1) V2, seq iter.Seq[V1]) iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for v := range seq {
			if !yield(mapFunc(v)) {
				return
			}
		}
	}
}

// Map2 returns a [iter.Seq2] similar to [Map].
func Map2[K1 comparable, V1 any, K2 comparable, V2 any](
	mapFunc func(K1, V1) (K2, V2),
	seq iter.Seq2[K1, V1],
) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range seq {
			if !yield(mapFunc(k, v)) {
				return
			}
		}
	}
}

// Filter returns a [iter.Seq] from those elements of seq for which filterFunc is true.
func Filter[V any](filterFunc func(V) bool, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if filterFunc(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Filter2 returns a [iter.Seq2] similar to [Filter].
func Filter2[K comparable, V any](filterFunc func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if filterFunc(k, v) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Enumerate returns a [iter.Seq2] which yields a count, starting from start,
// and the values obtained from iterating over seq.
func Enumerate[V any](seq iter.Seq[V], start int) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for v := range seq {
			if !yield(start, v) {
				return
			}
			start++
		}
	}
}

// AnyFunc returns true if checker returns true for any element in seq,
// otherwise it returns false.
func AnyFunc[V any](checker func(V) bool, seq iter.Seq[V]) bool {
	for v := range seq {
		if checker(v) {
			return true
		}
	}
	return false
}

// AnyFunc2 returns true if checker returns true for any element in seq
// otherwise it returns false.
func AnyFunc2[K comparable, V any](checker func(K, V) bool, seq iter.Seq2[K, V]) bool {
	for k, v := range seq {
		if checker(k, v) {
			return true
		}
	}
	return false
}

// FirstFunc returns the first value in seq for which checker returns true and
// 'true' or the zero value for type V and 'false' if no element of seq
// satisfiers checker.
func FirstFunc[V any](checker func(V) bool, seq iter.Seq[V]) (V, bool) { //nolint:ireturn
	for v := range seq {
		if checker(v) {
			return v, true
		}
	}
	var zero V
	return zero, false
}

// FirstFunc2 returns the first pair of values in seq for which checker returns
// true and 'true' or the zero values for types K and V and 'false' if no
// element of seq satisfiers the checker.
func FirstFunc2[K comparable, V any]( //nolint:ireturn
	checker func(K, V) bool,
	seq iter.Seq2[K, V],
) (K, V, bool) {
	for k, v := range seq {
		if checker(k, v) {
			return k, v, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// AllFunc returns true if checker returns true for all values in seq
// otherwise it returns false.
func AllFunc[V any](checker func(V) bool, seq iter.Seq[V]) bool {
	for v := range seq {
		if !checker(v) {
			return false
		}
	}
	return true
}

// AllFunc2 returns true if checker returns true for all values in seq
// otherwise it returns false.
func AllFunc2[K comparable, V any](checker func(K, V) bool, seq iter.Seq2[K, V]) bool {
	for k, v := range seq {
		if !checker(k, v) {
			return false
		}
	}
	return true
}

// Zip returns a [iter.Seq] the iterates over each sequence in seqs in parallel.
// Stops when the shortest sequence is exhausted.
func Zip[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		var nexts []func() (V, bool)
		for _, seq := range seqs {
			next, stop := iter.Pull(seq)
			defer stop()
			nexts = append(nexts, next)
		}

		for {
			for _, next := range nexts {
				v, ok := next()
				if !ok || !yield(v) {
					return
				}
			}
		}
	}
}

// Zip returns a [iter.Seq2] that yields pairs of values from seq1 and seq2.
// Stops when either seq of selectors is exhausted.
func ZipPair[V1 any, V2 any](seq1 iter.Seq[V1], seq2 iter.Seq[V2]) iter.Seq2[V1, V2] {
	return func(yield func(V1, V2) bool) {
		next1, stop1 := iter.Pull(seq1)
		next2, stop2 := iter.Pull(seq2)
		defer stop1()
		defer stop2()

		for {
			v1, ok := next1()
			if !ok {
				return
			}
			v2, ok := next2()
			if !ok {
				return
			}

			if !yield(v1, v2) {
				return
			}
		}
	}
}

// Zip2 returns a [iter.Seq2] like [Zip].
func Zip2[K comparable, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var nexts []func() (K, V, bool)
		for _, seq := range seqs {
			next, stop := iter.Pull2(seq)
			defer stop()
			nexts = append(nexts, next)
		}

		for {
			for _, next := range nexts {
				k, v, ok := next()
				if !ok || !yield(k, v) {
					return
				}
			}
		}
	}
}

// ZipLongest returns a [iter.Seq] like [Zip] but if the sequences are of
// uneven length, missing values are filled-in with fillValue.
func ZipLongest[V any](fillValue V, seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		nexts := make([]func() (V, bool), len(seqs))
		for i, seq := range seqs {
			next, stop := iter.Pull(seq)
			defer stop()
			nexts[i] = next
		}

		liveCount := len(seqs)
		results := make([]V, len(seqs))
		for {
			for i, next := range nexts {
				if next == nil {
					results[i] = fillValue
					continue
				}

				v, ok := next()
				if !ok {
					nexts[i] = nil
					liveCount--
					results[i] = fillValue
					continue
				}
				results[i] = v
			}
			if liveCount == 0 {
				return
			}
			for _, res := range results {
				if !yield(res) {
					return
				}
			}
		}
	}
}

// Range returns a [iter.Seq] that yields values step distance apart from start
// until end, not including end.
//
// Range panics if step is 0.
func Range(start int, end int, step int) iter.Seq[int] {
	if step == 0 {
		panic("step for Range must be non-zero")
	}
	return func(yield func(int) bool) {
		length := getRangeLen(start, end, step)
		for x := start; length > 0; x += step {
			if !yield(x) {
				return
			}
			length--
		}
	}
}

func getRangeLen(start int, end int, step int) int {
	if step > 0 && start < end {
		return 1 + (end-start-1)/step
	} else if step < 0 && start > end {
		return 1 + (start-end-1)/-step
	} else {
		return 0
	}
}

// Range from is like [Range] but has no end.
func RangeFrom(start int, step int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for x := start; ; x += step {
			if !yield(x) {
				return
			}
		}
	}
}

// RangeUntil is equivalent to
//
//	Range(0, end, step)
func RangeUntil(end int, step int) iter.Seq[int] {
	return Range(0, end, step)
}

// Cycle returns a [iter.Seq] that returns elements from the iterable and saves a copy of each.
// When the iterable is exhausted, elements from the saved copy are returned.
// Repeats indefinitely.
func Cycle[V any](iterable iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		var saved []V
		for v := range iterable {
			if !yield(v) {
				return
			}
			saved = append(saved, v)
		}

		for {
			for _, v := range saved {
				if !yield(v) {
					return
				}
			}
		}
	}
}

type seq2Store[K comparable, V any] struct {
	k K
	v V
}

func Cycle2[K comparable, V any](iterable iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var saved []seq2Store[K, V]
		for k, v := range iterable {
			if !yield(k, v) {
				return
			}
			saved = append(saved, seq2Store[K, V]{k, v})
		}

		for {
			for _, s := range saved {
				if !yield(s.k, s.v) {
					return
				}
			}
		}
	}
}

// Repeat returns a [iter.Seq] that returns value over and over again.
// Runs indefinitely if times is negative, otherwise runs that many times.
func Repeat[V any](value V, times int) iter.Seq[V] {
	return func(yield func(V) bool) {
		if times < 0 {
			for {
				if !yield(value) {
					return
				}
			}
		}
		for range times {
			if !yield(value) {
				return
			}
		}
	}
}

// Accumulate returns a [iter.Seq] that returns accumulated results from
// function.
// The function should accept two arguments, an accumulated total and a value
// from the input sequence.
func Accumulate[V1 any, V2 any](
	seq iter.Seq[V1],
	function func(acc V2, val V1) V2,
	initial V2,
) iter.Seq[V2] {
	current := initial
	return func(yield func(V2) bool) {
		for v := range seq {
			current = function(current, v)
			if !yield(current) {
				return
			}
		}
	}
}

// Compress returns a [iter.Seq] that returns elements from seq where the corresponding
// element in selectors is true. Stops when either the data or selectors iterables have been exhausted.
func Compress[V any](seq iter.Seq[V], selectors iter.Seq[bool]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v, sel := range ZipPair(seq, selectors) {
			if sel {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Compress2 is like [Compress] but for [iter.Seq2].
func Compress2[K comparable, V any](seq iter.Seq2[K, V], selectors iter.Seq[bool]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seqNext, seqStop := iter.Pull2(seq)
		selNext, selStop := iter.Pull(selectors)
		defer seqStop()
		defer selStop()

		for {
			k, v, ok := seqNext()
			if !ok {
				return
			}

			selekt, ok := selNext()
			if !ok {
				return
			}

			if selekt {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// DropWhile returns a [iter.Seq] that drops elements from seq while the
// predicate is true and afterwards returns every element.
func DropWhile[V any](seq iter.Seq[V], predicate func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		start := false
		for v := range seq {
			if !start && !predicate(v) {
				start = true
			}

			if start {
				if !yield(v) {
					return
				}
				continue
			}
		}
	}
}

// DropWhile2 is like [DropWhile] but for [iter.Seq2].
func DropWhile2[K comparable, V any](
	seq iter.Seq2[K, V],
	predicate func(K, V) bool,
) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		start := false
		for k, v := range seq {
			if !start && !predicate(k, v) {
				start = true
			}

			if start {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// TakeWhile returns a [iter.Seq] that returns elements from seq while the
// predicate is true.
func TakeWhile[V any](seq iter.Seq[V], predicate func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !predicate(v) {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

// TakeWhile2 is like [TakeWhile2] but for [iter.Seq2].
func TakeWhile2[K comparable, V any](
	seq iter.Seq2[K, V],
	predicate func(K, V) bool,
) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if !predicate(k, v) {
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// CollectIntoSlice is like [slices.Collect] but accepts a pre-allocated slice
// to collect values into.
func CollectIntoSlice[V any](seq iter.Seq[V], dest []V) {
	for i, v := range Enumerate(seq, 0) {
		dest[i] = v
	}
}

// CollectIntoMap is like [maps.Collect] but accepts a pre-allocated map to
// collect values into.
func CollectIntoMap[K comparable, V any](seq iter.Seq2[K, V], dest map[K]V) {
	for k, v := range seq {
		dest[k] = v
	}
}

// Keys returns a [iter.Seq] over the keys of seq.
func Keys[K comparable, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns a [iter.Seq] over the values of seq.
func Values[K comparable, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

// IterCtx returns a [iter.Seq] that yields values from seq until either
// seq is exhausted or ctx is cancelled, whichever comes first.
func IterCtx[V any](ctx context.Context, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		res := make(chan V)
		next, stop := iter.Pull(seq)

		// 'next' and 'stop' must not be called from multiple gorountines
		// simultaneously
		var pullMutex sync.Mutex
		defer func() {
			pullMutex.Lock()
			defer pullMutex.Unlock()
			stop()
		}()

		for {
			var ok bool
			var v V
			go func() {
				pullMutex.Lock()
				defer pullMutex.Unlock()
				v, ok = next()
				res <- v
			}()

			select {
			case v := <-res:
				if !ok || !yield(v) {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

// IterCtx2 is like [IterCtx] but for [iter.Seq2] sequences.
func IterCtx2[K comparable, V any](ctx context.Context, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		res := make(chan seq2Store[K, V])
		next, stop := iter.Pull2(seq)

		// 'next' and 'stop' must not be called from multiple gorountines
		// simultaneously
		var pullMutex sync.Mutex
		defer func() {
			pullMutex.Lock()
			defer pullMutex.Unlock()
			stop()
		}()

		for {
			var ok bool
			var v V
			var k K
			go func() {
				pullMutex.Lock()
				defer pullMutex.Unlock()
				k, v, ok = next()
				res <- seq2Store[K, V]{k, v}
			}()

			select {
			case s := <-res:
				if !ok || !yield(s.k, s.v) {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

// Slice returns a [iter.Seq] that slices up the provided sequence: returning
// elements step distance apart from start until end (excluding end).
//
// Slice will panic if step is not a positive integer.
func Slice[V any](seq iter.Seq[V], start int, end int, step int) iter.Seq[V] {
	if step <= 0 {
		panic("step for Slice must be a positive integer")
	}
	return func(yield func(V) bool) {
		next, stop := iter.Pull(seq)
		defer stop()

		for range RangeUntil(start, 1) {
			if _, ok := next(); !ok {
				return
			}
		}

		for i := start; i < end; i++ {
			v, ok := next()
			if !ok {
				return
			}

			if (i-start)%step == 0 {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// SliceUntil is a equivalent to
//
//	Slice(seq, 0, end, step)
func SliceUntil[V any](seq iter.Seq[V], end int, step int) iter.Seq[V] {
	return Slice(seq, 0, end, step)
}

// Slice2 is like [Slice] but for [iter.Seq2].
//
// Like [Slice] it will panic if step is not a positive integer.
func Slice2[K comparable, V any](
	seq iter.Seq2[K, V],
	start int,
	end int,
	step int,
) iter.Seq2[K, V] {
	if step <= 0 {
		panic("step for Slice2 must be a positive integer")
	}
	return func(yield func(K, V) bool) {
		next, stop := iter.Pull2(seq)
		defer stop()

		for range RangeUntil(start, 1) {
			if _, _, ok := next(); !ok {
				return
			}
		}

		for i := start; i < end; i++ {
			k, v, ok := next()
			if !ok {
				return
			}

			if (i-start)%step == 0 {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// SliceUntil2 is like [SliceUntil] but for [iter.Seq2].
func SliceUntil2[K comparable, V any](seq iter.Seq2[K, V], end int, step int) iter.Seq2[K, V] {
	return Slice2(seq, 0, end, step)
}

// Flatten returns a sequence that iterates across all keys and then all
// values of seq.
func Flatten[K comparable](seq iter.Seq2[K, K]) iter.Seq[K] {
	return Zip(Keys(seq), Values(seq))
}

// FlattenMap is a convenience wrapper of [Flatten].
func FlattenMap[K comparable](m map[K]K) iter.Seq[K] {
	return Flatten(maps.All(m))
}
