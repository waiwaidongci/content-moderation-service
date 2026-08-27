package application

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubDetector struct {
	name string
	res  DetectorResult
	err  error
}

func (s stubDetector) Name() string { return s.name }
func (s stubDetector) Detect(context.Context, string) (DetectorResult, error) {
	return s.res, s.err
}

func TestRunDetectorsReturnsResults(t *testing.T) {
	detectors := []Detector{
		stubDetector{name: "a", res: DetectorResult{Name: "a", Score: 1}},
		stubDetector{name: "b", res: DetectorResult{Name: "b", Score: 2}},
	}
	results, err := RunDetectors(context.Background(), "text", detectors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRunDetectorsReturnsErrorOnFailure(t *testing.T) {
	detectors := []Detector{stubDetector{name: "bad", err: errors.New("boom")}}
	start := make(chan struct{})
	done := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() {
			<-start
			_, err := RunDetectors(context.Background(), "text", detectors)
			done <- err
		}()
	}
	close(start)
	for i := 0; i < 2; i++ {
		select {
		case err := <-done:
			if err == nil {
				t.Fatal("expected error from failing detector")
			}
		case <-time.After(3 * time.Second):
			t.Fatal("RunDetectors hung when a detector returned an error")
		}
	}
}

func TestRunDetectorsReturnsErrorOnPartialFailure(t *testing.T) {
	detectors := []Detector{
		stubDetector{name: "ok", res: DetectorResult{Name: "ok", Score: 1}},
		stubDetector{name: "bad", err: errors.New("boom")},
	}
	start := make(chan struct{})
	done := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() {
			<-start
			_, err := RunDetectors(context.Background(), "text", detectors)
			done <- err
		}()
	}
	close(start)
	for i := 0; i < 2; i++ {
		select {
		case err := <-done:
			if err == nil {
				t.Fatal("expected error when a detector failed among successes")
			}
		case <-time.After(3 * time.Second):
			t.Fatal("RunDetectors hung when one detector failed")
		}
	}
}
