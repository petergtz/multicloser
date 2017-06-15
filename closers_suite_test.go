package closers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/petergtz/pegomock"
)

func TestClosers(t *testing.T) {
	RegisterFailHandler(Fail)
	pegomock.RegisterMockFailHandler(Fail)
	RunSpecs(t, "Closers Suite")
}
