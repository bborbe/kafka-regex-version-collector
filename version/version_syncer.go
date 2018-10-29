// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"context"
	"runtime"
	"time"

	"github.com/bborbe/kafka-regex-version-collector/avro"
	"github.com/bborbe/run"
	"github.com/golang/glog"
)

type fetcher interface {
	Fetch(ctx context.Context, versions chan<- avro.ApplicationVersionAvailable) error
}

type sender interface {
	Send(ctx context.Context, versions <-chan avro.ApplicationVersionAvailable) error
}

type Syncer struct {
	Fetcher fetcher
	Sender  sender
	Timeout time.Duration
}

func (s *Syncer) Sync(ctx context.Context) error {
	glog.V(1).Infof("sync started")
	defer glog.V(1).Infof("sync finished")

	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	versions := make(chan avro.ApplicationVersionAvailable, runtime.NumCPU())
	return run.CancelOnFirstError(
		ctx,
		func(ctx context.Context) error {
			defer close(versions)
			return s.Fetcher.Fetch(ctx, versions)
		},
		func(ctx context.Context) error {
			return s.Sender.Send(ctx, versions)
		},
	)
}
