package utils_test

import (
	"os"

	resourceutils "github.com/cheqd/cheqd-node/x/resource/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MediaType", func() {
	Describe("GetMediaType", func() {
		DescribeTable("Validate MIME type for different source files",
			func(path string, mt string) {
				data, err := os.ReadFile(path)
				Expect(err).To(BeNil())

				detected := resourceutils.DetectMediaType(data)
				Expect(detected).To(Equal(mt))
			},
			Entry("text file", "testdata/resource.txt", "text/plain; charset=utf-8"),
			Entry("csv file", "testdata/resource.csv", "text/csv"),
			Entry("dat file", "testdata/resource.dat", "application/octet-stream"),
			Entry("json file", "testdata/resource.json", "application/json"),
			Entry("pdf file", "testdata/resource.pdf", "application/pdf"),
			Entry("png file", "testdata/resource.png", "image/png"),
		)
	})
})
