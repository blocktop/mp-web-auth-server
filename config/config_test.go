package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

const (
	webAuthEndpoint = "https://localhost:3000/token"
)

var _ = Describe("config tests", func() {

	Describe("#Parse", func() {

		It("parses from embedded struct", func() {
			os.Setenv("MP_WEB_AUTH_ENDPOINT", webAuthEndpoint)
			c := GetConfig()
			Expect(c.WebAuthEndpoint).To(Equal(webAuthEndpoint))
		})
	})
})