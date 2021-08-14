package serial

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Serial", func() {
	Describe("Command Parsing", func() {
		It("Should parse a command", func() {
			line := "\nQSE>~DEVICE,GRAFIKEYE,70,3"
			cmd, err := parseQSEMessage(line)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd.Operation).To(Equal(OperationMonitor))
			Expect(cmd.Type).To(Equal(TypeDevice))
			Expect(cmd.IntegrationId).To(Equal(GrafikEye))
			Expect(cmd.CommandFields).To(BeEquivalentTo([]string{"70","3"}))
		})

		It("Should parse a command with double prefix", func() {
			line := "\nQSE>QSE>~DEVICE,GRAFIKEYE,70,3"
			cmd, err := parseQSEMessage(line)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd.Operation).To(Equal(OperationMonitor))
			Expect(cmd.Type).To(Equal(TypeDevice))
			Expect(cmd.IntegrationId).To(Equal(GrafikEye))
			Expect(cmd.CommandFields).To(BeEquivalentTo([]string{"70","3"}))
		})

		It("Should refuse a command with a wrong operation character", func() {
			line := "ABC>~DEVICE,GRAFIKEYE,70,3"
			_, err := parseQSEMessage(line)
			Expect(err).To(HaveOccurred())
		})
	})
})
