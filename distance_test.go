package main_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	rsa "reapp_students_assignement/pkg"
	mymocks "reapp_students_assignement/pkg/mocks"
)

const twoFcitiesTwoDcities = "{" +
	"\"destination_addresses\" : [ \"Teplice, Czechia\", \"Litvínov, Czechia\" ]," +
	"\"origin_addresses\" : [ \"Most, Czechia\", \"Ústí nad Labem, Czechia\" ]," +
	"\"rows\" : [" +
	"{" +
	"\"elements\" : [" +
	"{" +
	"\"distance\" : {" +
	"\"text\" : \"27.9 km\"," +
	"\"value\" : 27889" +
	"}," +
	"\"duration\" : {" +
	"\"text\" : \"28 mins\"," +
	"\"value\" : 1660" +
	"}," +
	"\"status\" : \"OK\"" +
	"},{" +
	"\"distance\" : {" +
	"\"text\" : \"13.5 km\"," +
	"\"value\" : 13458" +
	"}," +
	"\"duration\" : {" +
	"\"text\" : \"16 mins\"," +
	"\"value\" : 933" +
	"}," +
	"\"status\" : \"OK\"" +
	"}" +
	"]" +
	"},{" +
	"\"elements\" : [" +
	"{" +
	"\"distance\" : {" +
	"\"text\" : \"21.6 km\"," +
	"\"value\" : 21610" +
	"}," +
	"\"duration\" : {" +
	"\"text\" : \"21 mins\"," +
	"\"value\" : 1253" +
	"}," +
	"\"status\" : \"OK\"" +
	"},{" +
	"\"distance\" : {" +
	"\"text\" : \"43.2 km\"," +
	"\"value\" : 43178" +
	"}," +
	"\"duration\" : {" +
	"\"text\" : \"40 mins\"," +
	"\"value\" : 2402" +
	"}," +
	"\"status\" : \"OK\"" +
	"}" +
	"]" +
	"}" +
	"]," +
	"\"status\" : \"OK\"" +
	"}"
const oneFcityTwoTcities = "{" +
	"\"destination_addresses\" : [ \"Teplice, Czechia\", \"Litvínov, Czechia\" ]," +
	"\"origin_addresses\" : [ \"Most, Czechia\" ]," +
	"\"rows\" : [" +
	"{" +
	"\"elements\" : [" +
	"{" +
	"\"distance\" : {" +
	"\"text\" : \"27.9 km\"," +
	"\"value\" : 27889" +
	"}," +
	"\"duration\" : {" +
	"\"text\" : \"28 mins\"," +
	"\"value\" : 1660" +
	"}," +
	"\"status\" : \"OK\"" +
	"},{" +
	"\"distance\" : {" +
	"\"text\" : \"13.5 km\"," +
	"\"value\" : 13458" +
	"}," +
	"\"duration\" : {" +
	"\"text\" : \"16 mins\"," +
	"\"value\" : 933" +
	"}," +
	"\"status\" : \"OK\"" +
	"}" +
	"]" +
	"}" +
	"]," +
	"\"status\" : \"OK\"" +
	"}"
const oneFcityOneTcity = "{" +
	"\"destination_addresses\" : [ \"Litvínov, Czechia\" ]," +
	"\"origin_addresses\" : [ \"Teplice, Czechia\" ]," +
	"\"rows\" : [" +
	"{" +
	"\"elements\" : [" +
	"{" +
	"\"distance\" : {" +
	"\"text\" : \"20.9 km\"," +
	"\"value\" : 20853" +
	"}," +
	"\"duration\" : {" +
	"\"text\" : \"28 mins\"," +
	"\"value\" : 1657" +
	"}," +
	"\"status\" : \"OK\"" +
	"}" +
	"]" +
	"}" +
	"]," +
	"\"status\" : \"OK\"" +
	"}"
const invalidReq = "{" +
	"\"destination_addresses\" : []," +
	"\"origin_addresses\" : []," +
	"\"rows\" : []," +
	"\"status\" : \"INVALID_REQUEST\"" +
	"}"

var _ = Describe("Distance", func() {
	Context("Init DM", func() {
		var (
			mockCtrl       *gomock.Controller
			mockDistMatrix *mymocks.MockDistanceInterface
			dm             *rsa.DistanceMatrix
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockDistMatrix = mymocks.NewMockDistanceInterface(mockCtrl)
			dm = &rsa.DistanceMatrix{DistanceInterface: mockDistMatrix}
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})
		It("GetDistanceFromDM", func() {
			cities := []string{"Most", "Teplice", "Litvínov"}
			mockDistMatrix.EXPECT().GetDistance([]string{"Most, Czechia"}, []string{"Teplice, Czechia", "Litvínov, Czechia"}).Return([]byte(oneFcityTwoTcities))
			mockDistMatrix.EXPECT().GetDistance([]string{"Teplice, Czechia"}, []string{"Litvínov, Czechia"}).Return([]byte(oneFcityOneTcity))
			DM := rsa.NewDistanceMatrix(cities, dm)
			Expect(DM.GetDistanceFromDM("Teplice", "Most")).To(Equal(rsa.Distance{27889, 1660}))
			Expect(DM.GetDistanceFromDM("Litvínov", "Most")).To(Equal(rsa.Distance{13458, 933}))
			Expect(DM.GetDistanceFromDM("Most", "Litvínov")).To(Equal(rsa.Distance{13458, 933}))
		})
		It("Error with one city", func() {
			Expect(rsa.NewDistanceMatrix([]string{"Teplice"}, dm)).To(Equal(rsa.DistanceMatrix{}))
		})
		It("Error with zero city", func() {
			Expect(rsa.NewDistanceMatrix([]string{}, dm)).To(Equal(rsa.DistanceMatrix{}))
		})
	})
	Context("Get distance", func() {
		var (
			mockCtrl       *gomock.Controller
			mockDistMatrix *mymocks.MockDistanceInterface
			dm             *rsa.DistanceMatrix
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockDistMatrix = mymocks.NewMockDistanceInterface(mockCtrl)
			dm = &rsa.DistanceMatrix{DistanceInterface: mockDistMatrix}
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("Run all source and destination cities", func() {
			fcities := []string{"Most", "Ústí nad Labem"}
			tcities := []string{"Teplice", "Litvínov"}
			var expected = make(map[string]rsa.Distance)
			expected["4d6f73747c5465706c696365"] = rsa.Distance{27889, 1660}
			expected["4d6f73747c4c697476c3ad6e6f76"] = rsa.Distance{13458, 933}
			expected["c39a7374c3ad206e6164204c6162656d7c5465706c696365"] = rsa.Distance{21610, 1253}
			expected["c39a7374c3ad206e6164204c6162656d7c4c697476c3ad6e6f76"] = rsa.Distance{43178, 2402}
			mockDistMatrix.EXPECT().GetDistance(fcities, tcities).Return([]byte(twoFcitiesTwoDcities))
			Expect(dm.GetMatrixFor(fcities, tcities)).To(Equal(expected))
		})
		It("Run one source and few destination cities", func() {
			fcities := []string{"Most"}
			tcities := []string{"Teplice", "Litvínov"}
			var expected = make(map[string]rsa.Distance)
			expected["4d6f73747c5465706c696365"] = rsa.Distance{27889, 1660}
			expected["4d6f73747c4c697476c3ad6e6f76"] = rsa.Distance{13458, 933}
			mockDistMatrix.EXPECT().GetDistance(fcities, tcities).Return([]byte(oneFcityTwoTcities))
			Expect(dm.GetMatrixFor(fcities, tcities)).To(Equal(expected))
		})
		It("Run zero source and few destination cities", func() {
			fcities := []string{}
			tcities := []string{"Teplice", "Litvínov"}
			var expected = make(map[string]rsa.Distance)
			mockDistMatrix.EXPECT().GetDistance(fcities, tcities).Return([]byte(invalidReq))
			Expect(dm.GetMatrixFor(fcities, tcities)).To(Equal(expected))
		})
		It("Run few source and zero destination cities", func() {
			fcities := []string{"Most", "Ústí nad Labem"}
			tcities := []string{}
			var expected = make(map[string]rsa.Distance)
			mockDistMatrix.EXPECT().GetDistance(fcities, tcities).Return([]byte(invalidReq))
			Expect(dm.GetMatrixFor(fcities, tcities)).To(Equal(expected))
		})
	})
	Context("Test calculate fitness", func() {
		var (
			mockCtrl       *gomock.Controller
			mockDistMatrix *mymocks.MockDistanceInterface
			dm             *rsa.DistanceMatrix
			attr           rsa.ABCAttributes
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockDistMatrix = mymocks.NewMockDistanceInterface(mockCtrl)
			dm = &rsa.DistanceMatrix{DistanceInterface: mockDistMatrix}
			attr = rsa.ABCAttributes{ReferenceMaxDistance: 50000, PenaltyForTheSameWorkplace: 50}
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})
		It("Fitenss 1", func() {
			cities := []string{"Most", "Teplice", "Litvínov"}
			mockDistMatrix.EXPECT().GetDistance([]string{"Most, Czechia"}, []string{"Teplice, Czechia", "Litvínov, Czechia"}).Return([]byte(oneFcityTwoTcities))
			mockDistMatrix.EXPECT().GetDistance([]string{"Teplice, Czechia"}, []string{"Litvínov, Czechia"}).Return([]byte(oneFcityOneTcity))
			DM := rsa.NewDistanceMatrix(cities, dm)
			Expect(DM.CalculateLogDistance(attr, "Teplice", "Most")).To(Equal(14.767408584948104))
			Expect(DM.CalculateLogDistance(attr, "Teplice", "Zlín")).To(Equal(float64(50000)))
			Expect(DM.CalculateLogDistance(attr, "Teplice", "Teplice")).To(Equal(float64(0)))
		})
	})
})
