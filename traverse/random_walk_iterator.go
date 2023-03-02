package traverse

import (
	"crypto/rand"
	"math/big"

	"github.com/hmdsefi/gograph"
)

// randomWalkIterator implements the Iterator interface to travers
// a graph in a random walk fashion.
//
// Random walk is a stochastic process used to explore a graph, where
// a walker moves through the graph by following random edges. At each
// step, the walker chooses a random neighbor of the current node and
// moves to it, and the process is repeated until a stopping condition
// is met.
type randomWalkIterator[T comparable] struct {
	graph       gograph.Graph[T]   // the graph that being traversed.
	start       *gograph.Vertex[T] // the starting point of the traversal.
	current     *gograph.Vertex[T] // the latest node that has been returned by the iterator.
	steps       int                // the maximum number of steps to be taken during the traversal.
	currentStep int                // the step counter.
}

// NewRandomWalkIterator creates a new instance of randomWalkIterator
// and returns it as the Iterator interface.
func NewRandomWalkIterator[T comparable](graph gograph.Graph[T], start *gograph.Vertex[T], steps int) Iterator[T] {
	return &randomWalkIterator[T]{
		graph:   graph,
		start:   start,
		current: start,
		steps:   steps,
	}
}

// HasNext returns a boolean indicating whether there are more vertices
// to be visited or not.
func (r *randomWalkIterator[T]) HasNext() bool {
	return r.current.OutDegree() > 0 && r.currentStep < r.steps
}

// Next returns the next vertex to be visited in the random walk traversal.
// It chooses one of the neighbors randomly and returns it.
//
// If the HasNext is false, returns nil.
func (r *randomWalkIterator[T]) Next() *gograph.Vertex[T] {
	if !r.HasNext() {
		return nil
	}

	if r.currentStep == 0 {
		r.currentStep++
		return r.current
	}

	neighbors := r.current.Neighbors()

	i, _ := rand.Int(rand.Reader, big.NewInt(int64(len(neighbors))))
	r.current = neighbors[i.Int64()]
	r.currentStep++

	return r.current
}

// Iterate iterates through the vertices in random order and applies
// the given function to each vertex. If the function returns an error,
// the iteration stops and the error is returned.
func (r *randomWalkIterator[T]) Iterate(f func(v *gograph.Vertex[T]) error) error {
	for r.HasNext() {
		if err := f(r.Next()); err != nil {
			return err
		}
	}

	return nil
}

// Reset resets the iterator by setting the initial state of the iterator.
func (r *randomWalkIterator[T]) Reset() {
	r.current = r.start
	r.currentStep = 0
}
