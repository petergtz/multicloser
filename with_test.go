package multicloser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/petergtz/pegomock"

	"github.com/petergtz/closers"
)

var _ = Describe("With", func() {
	var closer *MockCloser

	BeforeEach(func() {
		closer = NewMockCloser()
	})

	It("no error", func() {
		e := multicloser.New(closer).CloseAfter(func() error { return nil }, nil)

		closer.VerifyWasCalledOnce().Close()
		Expect(e).NotTo(HaveOccurred())
	})

	It("error during run", func() {
		e := multicloser.New(closer).CloseAfter(func() error { return errors.New("during run") }, nil)

		closer.VerifyWasCalledOnce().Close()
		Expect(e.Error()).To(Equal("during run"))
	})

	It("error during close", func() {
		When(closer.Close()).ThenReturn(errors.New("during close"))

		e := multicloser.New(closer).CloseAfter(func() error { return nil }, nil)

		closer.VerifyWasCalledOnce().Close()
		Expect(e.Error()).To(Equal("during close"))
	})

	It("error during run and close", func() {
		When(closer.Close()).ThenReturn(errors.New("during close"))

		e := multicloser.New(closer).CloseAfter(func() error { return errors.New("during run") }, nil)

		closer.VerifyWasCalledOnce().Close()
		Expect(e.Error()).To(Equal("during run"))
	})

	It("converts a close error correctly", func() {
		When(closer.Close()).ThenReturn(errors.New("during close"))

		e := multicloser.New(closer).CloseAfter(func() error { return nil }, func(e error) error {
			return errors.New("The converted error: " + e.Error())
		})

		closer.VerifyWasCalledOnce().Close()
		Expect(e.Error()).To(Equal("The converted error: during close"))
	})
})
