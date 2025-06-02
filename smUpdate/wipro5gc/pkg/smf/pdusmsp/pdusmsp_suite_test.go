package pdusmsp_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPdusmsp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pdusmsp Suite")
}
