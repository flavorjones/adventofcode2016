package adventofcode2016_test

import (
	"crypto/md5"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
)

type Door struct {
	id string
}

var passwordLen = 8

func md5sum(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

func (d Door) password() string {
	password := make([]byte, passwordLen)
	index := 0
	for j := 0; j < passwordLen; j++ {
		for {
			hash := md5sum(d.id + strconv.Itoa(index))
			index++
			if hash[0:5] == "00000" {
				password[j] = hash[5]
				break
			}
		}
	}
	return string(password)
}

var _ = Describe("Day5", func() {
	Describe("Door", func() {
		Describe("#password", func() {
			It("finds the right password", func() {
				Expect(Door{"abc"}.password()).To(Equal("18f47a30"))
			})
		})
	})

	Describe("star 1", func() {
		It("finds the answer", func() {
			fmt.Println("star 1: ", Door{"ojvtpuvg"}.password())
		})
	})
})
