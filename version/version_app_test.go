// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version_test

import (
	"time"

	"github.com/bborbe/kafka-regex-version-collector/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version App", func() {
	var app *version.App
	BeforeEach(func() {
		app = &version.App{
			Wait:              time.Hour,
			Port:              1337,
			KafkaBrokers:      "kafka:9092",
			KafkaTopic:        "my-topic",
			SchemaRegistryUrl: "http://localhost:8081",
			Regex:             "(.*)",
			Url:               "http://www.example.com",
			Application:       "myApp",
		}
	})
	It("Validate without error", func() {
		Expect(app.Validate()).NotTo(HaveOccurred())
	})
	It("Validate returns error if port is 0", func() {
		app.Port = 0
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error if wait is 0", func() {
		app.Wait = 0
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error KafkaBrokers is empty", func() {
		app.KafkaBrokers = ""
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error KafkaTopic is empty", func() {
		app.KafkaTopic = ""
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error SchemaRegistryUrl is empty", func() {
		app.SchemaRegistryUrl = ""
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error Url is empty", func() {
		app.Url = ""
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error Regex is empty", func() {
		app.Regex = ""
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error Application is empty", func() {
		app.Application = ""
		Expect(app.Validate()).To(HaveOccurred())
	})
})
