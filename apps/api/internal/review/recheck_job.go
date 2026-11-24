package review

import (
	"context"
	"log/slog"
	"time"

	"github.com/riverqueue/river"
)

// SampleArgs is the River job that draws a quarter's re-check sample.
type SampleArgs struct{}

// Kind names the job in River.
func (SampleArgs) Kind() string { return "sample_rechecks" }

// SampleWorker draws the current quarter's sample. Drawing is once per
// quarter, so running it again is harmless.
type SampleWorker struct {
	river.WorkerDefaults[SampleArgs]
	svc *Service
	log *slog.Logger
}

// NewSampleWorker returns the re-check sampling worker.
func NewSampleWorker(svc *Service, log *slog.Logger) *SampleWorker {
	return &SampleWorker{svc: svc, log: log}
}

// Work draws the sample for the quarter the job runs in.
func (w *SampleWorker) Work(ctx context.Context, _ *river.Job[SampleArgs]) error {
	quarter := Quarter(time.Now())
	n, err := w.svc.SampleAsSystem(ctx, quarter)
	if err != nil {
		return err
	}
	w.log.InfoContext(ctx, "re-check sample", "quarter", quarter, "items", n)
	return nil
}

// quarterly fires at the start of each calendar quarter, UTC.
type quarterly struct{}

func (quarterly) Next(t time.Time) time.Time { return quarterStart(t).AddDate(0, 3, 0) }

// SamplePeriodically is the periodic job for the worker: at the start of each
// quarter, and once when the worker starts in case it was down then.
func SamplePeriodically() *river.PeriodicJob {
	return river.NewPeriodicJob(quarterly{}, func() (river.JobArgs, *river.InsertOpts) {
		return SampleArgs{}, nil
	}, &river.PeriodicJobOpts{RunOnStart: true})
}
