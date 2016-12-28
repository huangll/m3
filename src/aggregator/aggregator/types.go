// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package aggregator

import (
	"github.com/m3db/m3metrics/metric/unaggregated"
	"github.com/m3db/m3x/instrument"
)

// Aggregator aggregates different types of metrics
type Aggregator interface {
	// AddCounterWithPolicies adds a counter with policies for aggregation
	AddCounterWithPolicies(cp unaggregated.CounterWithPolicies)

	// AddBatchTimerWithPolicies adds a batch timer with policies for aggregation
	AddBatchTimerWithPolicies(btp unaggregated.BatchTimerWithPolicies)

	// AddGaugeWithPolicies adds a gauge with policies for aggregation
	AddGaugeWithPolicies(gp unaggregated.GaugeWithPolicies)
}

// SnapshotResult is the snapshot result
type SnapshotResult struct {
	CountersWithPolicies    []unaggregated.CounterWithPolicies
	BatchTimersWithPolicies []unaggregated.BatchTimerWithPolicies
	GaugesWithPolicies      []unaggregated.GaugeWithPolicies
}

// MockAggregator provide an aggregator for testing purposes
type MockAggregator interface {
	Aggregator

	// Snapshot returns a copy of the aggregated data and resets aggregations
	Snapshot() SnapshotResult
}

// Options provide a set of options for the aggregator
type Options interface {
	// SetInstrumentOptions sets the instrument options
	SetInstrumentOptions(value instrument.Options) Options

	// InstrumentOptions returns the instrument options
	InstrumentOptions() instrument.Options
}