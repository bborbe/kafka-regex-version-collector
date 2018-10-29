// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/bborbe/kafka-regex-version-collector/avro"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

//go:generate counterfeiter -o ../mocks/http_client.go --fake-name HttpClient . HttpClient
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Fetcher struct {
	HttpClient  HttpClient
	Url         string
	Application string
	Regex       string
}

func (f *Fetcher) Fetch(ctx context.Context, versions chan<- avro.ApplicationVersionAvailable) error {
	re, err := regexp.Compile(f.Regex)
	if err != nil {
		return errors.Wrap(err, "parse regex failed")
	}
	req, err := http.NewRequest(http.MethodGet, f.Url, nil)
	if err != nil {
		return errors.Wrap(err, "build request failed")
	}
	glog.V(1).Infof("%s %s", req.Method, req.URL.String())
	resp, err := f.HttpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	if resp.StatusCode/100 != 2 {
		return errors.New("request status code != 2xx")
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read body failed")
	}
	matches := re.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		select {
		case <-ctx.Done():
			glog.Infof("context done => return")
			return nil
		case versions <- avro.ApplicationVersionAvailable{
			App:     f.Application,
			Version: match[1],
		}:
		}
	}
	return nil
}
