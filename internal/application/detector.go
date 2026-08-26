// Package implementation for policy-driven content moderation and human review.
package application

import (
	"context"
	"fmt"
	"sync"
)

type DetectorResult struct {
	Name   string
	Score  float64
	Labels []string
}
type Detector interface {
	Name() string
	Detect(context.Context, string) (DetectorResult, error)
}

func RunDetectors(ctx context.Context, text string, detectors []Detector) ([]DetectorResult, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	out := make(chan DetectorResult, len(detectors))
	errs := make(chan error, len(detectors))
	var wg sync.WaitGroup
	for _, d := range detectors {
		wg.Add(1)
		go func(det Detector) {
			defer wg.Done()
			r, err := det.Detect(ctx, text)
			if err != nil {
				errs <- fmt.Errorf("detector %s: %w", det.Name(), err)
				return
			}
			out <- r
		}(d)
	}
	wg.Wait()
	close(out)
	close(errs)
	results := []DetectorResult{}
	for r := range out {
		results = append(results, r)
	}
	if len(results) == 0 {
		for e := range errs {
			return nil, e
		}
	}
	return results, nil
}
