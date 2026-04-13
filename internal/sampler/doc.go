// Package sampler implements log-line sampling strategies for logslice.
//
// Two complementary strategies are provided:
//
//   - Rate sampling: retain every Nth line (e.g. WithRate(10) keeps lines
//     10, 20, 30, …). Useful for evenly-spaced thinning of high-volume logs.
//
//   - Random sampling: retain each line independently with probability p
//     (e.g. WithRandom(0.1) keeps ~10 % of lines). Useful when a
//     statistically representative subset is required.
//
// Both strategies can be combined; a line must pass both filters to be kept.
//
// Example usage:
//
//	s := sampler.New(time.Now().UnixNano(),
//		sampler.WithRate(5),
//		sampler.WithRandom(0.8),
//	)
//	for _, line := range lines {
//		if s.Keep() {
//			// process line
//		}
//	}
package sampler
