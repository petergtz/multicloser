package multicloser_test

import (
	"errors"
	"io"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/petergtz/pegomock"

	"github.com/petergtz/closers"
)

func Example_CopyFile(dstName, srcName string) (int64, error) {
	multiCloser := multicloser.New()

	src, e := os.Open(srcName)
	if e != nil {
		multiCloser.Close() // no-op. Not strictly necessary, but helps to create some kind of pattern
		return 0, errors.New("Error while opening file for reading. Caused by: " + e.Error())
	}
	multiCloser.Add(src)

	dst, e := os.Create(dstName)
	if e != nil {
		multiCloser.Close()
		return 0, errors.New("Error while opening file for writing. Caused by: " + e.Error())
	}
	multiCloser.Add(dst)

	var numBytesWritten int64
	e = multiCloser.CloseAfter(func() error {
		var e error
		numBytesWritten, e = io.Copy(dst, src)
		if e != nil {
			return errors.New("Error while copying. Caused by: " + e.Error())
		}
		return nil
	}, func(e error) error {
		return errors.New("Error while closing. Caused by: " + e.Error())
	})
	return numBytesWritten, e
}

var _ = Describe("Closers", func() {
	var closer1, closer2, closer3 *MockCloser

	BeforeEach(func() {
		closer1 = NewMockCloser()
		closer2 = NewMockCloser()
		closer3 = NewMockCloser()
	})
	Context("1 closer", func() {
		It("returns no errors", func() {
			Expect(multicloser.New(closer1).Close()).To(Succeed())

			closer1.VerifyWasCalledOnce().Close()
		})

		It("returns the error the closer returns", func() {
			When(closer1.Close()).ThenReturn(errors.New("Close failed"))

			e := multicloser.New(closer1).Close()

			Expect(e.Error()).To(Equal("Close failed"))
			closer1.VerifyWasCalledOnce().Close()
		})
	})

	Context("2 closers", func() {
		It("returns no errors", func() {
			Expect(multicloser.New(closer1, closer2).Close()).To(Succeed())

			inOrderContext := new(InOrderContext)
			closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
		})

		Context("closer 1 returns error", func() {
			It("closes all closers, returns error from closer 1", func() {
				When(closer1.Close()).ThenReturn(errors.New("Close failed"))

				e := multicloser.New(closer1, closer2).Close()

				Expect(e.Error()).To(Equal("Close failed"))
				inOrderContext := new(InOrderContext)
				closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			})
		})

		Context("closer 2 returns error", func() {
			It("closes all closers, returns error from closer 2", func() {
				When(closer2.Close()).ThenReturn(errors.New("Close failed"))

				e := multicloser.New(closer1, closer2).Close()

				Expect(e.Error()).To(Equal("Close failed"))
				inOrderContext := new(InOrderContext)
				closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			})
		})

		Context("closer 1 and 2 return an error", func() {
			It("closes all closers, returns error from closer 2", func() {
				When(closer2.Close()).ThenReturn(errors.New("Close failed"))

				e := multicloser.New(closer1, closer2).Close()

				Expect(e.Error()).To(Equal("Close failed"))
				inOrderContext := new(InOrderContext)
				closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			})
		})
	})

	Context("3 closers", func() {
		It("returns no errors", func() {
			multicloser.New(closer1, closer2, closer3).Close()

			inOrderContext := new(InOrderContext)
			closer3.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
		})
		Context("closer 1 returns error", func() {
			It("closes all closers, returns error from closer 1", func() {
				When(closer1.Close()).ThenReturn(errors.New("Close failed"))

				e := multicloser.New(closer1, closer2, closer3).Close()

				Expect(e.Error()).To(Equal("Close failed"))
				inOrderContext := new(InOrderContext)
				closer3.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			})
		})

		Context("closer 2 returns error", func() {
			It("closes all closers, returns error from closer 2", func() {
				When(closer2.Close()).ThenReturn(errors.New("Close failed"))

				e := multicloser.New(closer1, closer2, closer3).Close()

				Expect(e.Error()).To(Equal("Close failed"))
				inOrderContext := new(InOrderContext)
				closer3.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			})
		})

		Context("closer 3 returns error", func() {
			It("closes all closers, returns error from closer 2", func() {
				When(closer3.Close()).ThenReturn(errors.New("Close failed"))

				e := multicloser.New(closer1, closer2, closer3).Close()

				Expect(e.Error()).To(Equal("Close failed"))
				inOrderContext := new(InOrderContext)
				closer3.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			})
		})

		Context("closer 2 and 3 return an error", func() {
			It("closes all closers, returns error from closer 3", func() {
				When(closer3.Close()).ThenReturn(errors.New("Close failed"))

				e := multicloser.New(closer1, closer2, closer3).Close()

				Expect(e.Error()).To(Equal("Close failed"))
				inOrderContext := new(InOrderContext)
				closer3.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			})
		})

		Context("closer 1 and 3 return an error", func() {
			It("closes all closers, returns error from closer 3", func() {
				When(closer3.Close()).ThenReturn(errors.New("Close failed"))

				e := multicloser.New(closer1, closer2, closer3).Close()

				Expect(e.Error()).To(Equal("Close failed"))
				inOrderContext := new(InOrderContext)
				closer3.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			})
		})

		Context("closer 1 and 2 return an error", func() {
			It("closes all closers, returns error from closer 2", func() {
				When(closer2.Close()).ThenReturn(errors.New("Close failed"))

				e := multicloser.New(closer1, closer2, closer3).Close()

				Expect(e.Error()).To(Equal("Close failed"))
				inOrderContext := new(InOrderContext)
				closer3.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer2.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
				closer1.VerifyWasCalledInOrder(Once(), inOrderContext).Close()
			})
		})
	})
})
