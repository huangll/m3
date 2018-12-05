// Copyright (c) 2018 Uber Technologies, Inc.
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

package debug

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/m3db/m3/src/query/api/v1/handler/prometheus"
	"github.com/m3db/m3/src/query/block"
	"github.com/m3db/m3/src/query/models"
	"github.com/m3db/m3/src/query/storage"
	"github.com/m3db/m3/src/query/ts"
)

type debugStorage struct {
	seriesList []*ts.Series
}

// NewStorage creates a new debug storage instance.
func NewStorage(promReadResp prometheus.PromResp) (storage.Storage, error) {
	seriesList, err := PromResultToSeriesList(promReadResp, models.NewTagOptions())
	if err != nil {
		return nil, err
	}

	return &debugStorage{
		seriesList: seriesList,
	}, nil
}

func (s *debugStorage) Fetch(
	ctx context.Context,
	query *storage.FetchQuery,
	options *storage.FetchOptions,
) (*storage.FetchResult, error) {
	return &storage.FetchResult{
		SeriesList: s.seriesList,
	}, nil
}

func (s *debugStorage) FetchBlocks(
	ctx context.Context,
	query *storage.FetchQuery,
	options *storage.FetchOptions,
) (block.Result, error) {
	fetchResult, err := s.Fetch(ctx, query, options)
	if err != nil {
		return block.Result{}, err
	}

	return storage.FetchResultToBlockResult(fetchResult, query)
}

// PromResultToSeriesList converts a prom result to a series list
func PromResultToSeriesList(promReadResp prometheus.PromResp, tagOptions models.TagOptions) ([]*ts.Series, error) {
	results := promReadResp.Data.Result
	seriesList := make([]*ts.Series, len(results))

	for i, result := range results {
		dps := make(ts.Datapoints, len(result.Values))
		for i, dp := range result.Values {
			dps[i].Timestamp = time.Unix(0, int64(dp[0].(float64))*int64(time.Second))
			s, err := strconv.ParseFloat(dp[1].(string), 64)
			if err != nil {
				return nil, err
			}
			dps[i].Value = s
		}

		tagSlice := make([]models.Tag, len(result.Metric))
		tagsIdx := 0
		for name, val := range result.Metric {
			tagSlice[tagsIdx] = models.Tag{
				Name:  []byte(name),
				Value: []byte(val),
			}

			tagsIdx++
		}

		tags := models.Tags{
			Opts: tagOptions,
			Tags: tagSlice,
		}

		name, _ := tags.Name()

		seriesList[i] = ts.NewSeries(
			string(name),
			dps,
			tags,
		)
	}

	return seriesList, nil
}

func (s *debugStorage) Type() storage.Type {
	return storage.TypeDebug
}

func (s *debugStorage) FetchTags(
	ctx context.Context,
	query *storage.FetchQuery,
	_ *storage.FetchOptions,
) (*storage.SearchResults, error) {
	return nil, errors.New("FetchTags not implemented")
}

func (s *debugStorage) CompleteTags(
	ctx context.Context,
	query *storage.CompleteTagsQuery,
	_ *storage.FetchOptions,
) (*storage.CompleteTagsResult, error) {
	return nil, errors.New("CompleteTags not implemented")
}

func (s *debugStorage) Close() error {
	return nil
}

func (s *debugStorage) Write(
	ctx context.Context,
	query *storage.WriteQuery,
) error {
	return errors.New("write not implemented")
}
