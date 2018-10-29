// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/bborbe/kafka-regex-version-collector/avro"
	"github.com/bborbe/kafka-regex-version-collector/mocks"
	"github.com/bborbe/kafka-regex-version-collector/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version Fetcher", func() {
	var httpClient *mocks.HttpClient
	var fetcher *version.Fetcher
	BeforeEach(func() {
		httpClient = &mocks.HttpClient{}
		fetcher = &version.Fetcher{
			HttpClient:  httpClient,
			Regex:       "<a>(.*)</a>",
			Application: "MyApp",
		}
	})
	It("work with empty", func() {
		httpClient.DoReturnsOnCall(0, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
		}, nil)
		versions := make(chan avro.ApplicationVersionAvailable)
		var list []avro.ApplicationVersionAvailable
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for version := range versions {
				list = append(list, version)
			}
		}()
		err := fetcher.Fetch(context.Background(), versions)
		Expect(err).NotTo(HaveOccurred())
		close(versions)
		wg.Wait()
		Expect(list).To(HaveLen(0))
	})
	It("work with empty", func() {
		httpClient.DoReturnsOnCall(0, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`<a>v1</a>`)),
		}, nil)
		versions := make(chan avro.ApplicationVersionAvailable)
		var list []avro.ApplicationVersionAvailable
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for version := range versions {
				list = append(list, version)
			}
		}()
		err := fetcher.Fetch(context.Background(), versions)
		close(versions)
		Expect(err).NotTo(HaveOccurred())
		wg.Wait()
		Expect(list).To(HaveLen(1))
		Expect(list[0].App).To(Equal("MyApp"))
		Expect(list[0].Version).To(Equal("v1"))
	})
})
