package prng

// PRNG interface for deterministic random number generation.
// TODO: Implement xoshiro256** or PCG64 for better quality
// For now, using standard library rand is sufficient for the stub.

type PRNG interface {
	Uint64() uint64
	Intn(n int) int
	Float64() float64
}