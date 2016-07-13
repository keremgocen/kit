// Package teststat provides helpers for testing metrics backends.
package teststat

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/go-kit/kit/metrics2"
)

// TestCounter puts some deltas through the counter, and then calls the value
// func to check that the counter has the correct final value.
func TestCounter(counter metrics.Counter, value func() float64) error {
	a := rand.Perm(100)
	n := rand.Intn(len(a))

	var want float64
	for i := 0; i < n; i++ {
		f := float64(a[i])
		counter.Add(f)
		want += f
	}

	if have := value(); want != have {
		return fmt.Errorf("want %f, have %f", want, have)
	}

	return nil
}

// TestGauge puts some values through the gauge, and then calls the value func
// to check that the gauge has the correct final value.
func TestGauge(gauge metrics.Gauge, value func() float64) error {
	a := rand.Perm(100)
	n := rand.Intn(len(a))

	var want float64
	for i := 0; i < n; i++ {
		f := float64(a[i])
		gauge.Set(f)
		want = f
	}

	if have := value(); want != have {
		return fmt.Errorf("want %f, have %f", want, have)
	}

	return nil
}

// TestHistogram puts some observations through the histogram, and then calls
// the quantiles func to checks that the histogram has computed the correct
// quantiles within some tolerance
func TestHistogram(histogram metrics.Histogram, quantiles func() (p50, p90, p95, p99 float64), tolerance float64) error {
	var (
		seed  = rand.Int()
		mean  = 500
		stdev = 25
	)
	populateNormal(histogram, seed, mean, stdev)

	want50, want90, want95, want99 := normalQuantiles(mean, stdev)
	have50, have90, have95, have99 := quantiles()

	var errs []string
	if want, have := want50, have50; !cmp(want, have, tolerance) {
		errs = append(errs, fmt.Sprintf("p50: want %f, have %f", want, have))
	}
	if want, have := want90, have90; !cmp(want, have, tolerance) {
		errs = append(errs, fmt.Sprintf("p90: want %f, have %f", want, have))
	}
	if want, have := want95, have95; !cmp(want, have, tolerance) {
		errs = append(errs, fmt.Sprintf("p95: want %f, have %f", want, have))
	}
	if want, have := want99, have99; !cmp(want, have, tolerance) {
		errs = append(errs, fmt.Sprintf("p99: want %f, have %f", want, have))
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func cmp(want, have, tol float64) bool {
	if (math.Abs(want-have) / want) > tol {
		return false
	}
	return true
}