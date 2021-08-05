package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestReappStudentsAssignement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReappStudentsAssignement Suite")
}
