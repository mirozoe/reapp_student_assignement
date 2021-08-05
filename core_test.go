package main_test

import (
	"math"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "reapp_students_assignement/pkg"
	rsa "reapp_students_assignement/pkg"
	mymocks "reapp_students_assignement/pkg/mocks"
)

var _ = Describe("Core", func() {
	Context("TranslateInputs tests", func() {
		It("All students to all workplaces", func() {
			sd1 := []StudentDetail{
				{1, 0, 20},
				{2, 1, 21},
				{3, 0, 22},
			}
			sd2 := []StudentDetail{
				{4, 1, 20},
				{5, 0, 21},
			}
			si := []StudentInput{
				{"Most", sd1},
				{"Teplice", sd2},
			}
			wd1 := []WorkplaceDetail{
				{1, 0, 10},
				{2, 1, 20},
			}
			wd2 := []WorkplaceDetail{
				{3, 0, 1},
				{4, 1, 2},
			}
			wi := []WorkplaceInput{
				{"Teplice", wd1},
				{"Most", wd2},
			}
			is := InputSummary{si, wi}
			exs1 := []StudentDetailWithCity{
				{1, 0, 20, "Most"},
				{3, 0, 22, "Most"},
				{5, 0, 21, "Teplice"},
			}
			exs2 := []StudentDetailWithCity{
				{2, 1, 21, "Most"},
				{4, 1, 20, "Teplice"},
			}
			exs := []StudentsByBranches{
				{0, exs1},
				{1, exs2},
			}
			exw1 := []WorkplaceDetailWithCity{
				{1, 0, 10, "Teplice"},
				{3, 0, 1, "Most"},
			}
			exw2 := []WorkplaceDetailWithCity{
				{2, 1, 20, "Teplice"},
				{4, 1, 2, "Most"},
			}
			exw := []WorkplacesByBranches{
				{0, exw1},
				{1, exw2},
			}

			rs, rw := TranslateInputs(is)
			Expect(rs).To(Equal(exs))
			Expect(rw).To(Equal(exw))
		})
		It("Some students to some workplaces", func() {
			sd1 := []StudentDetail{
				{1, 0, 20},
				{2, 1, 21},
				{3, 0, 22},
			}
			sd2 := []StudentDetail{
				{4, 1, 20},
				{5, 0, 21},
			}
			si := []StudentInput{
				{"Most", sd1},
				{"Teplice", sd2},
			}
			wd1 := []WorkplaceDetail{
				{1, 0, 10},
				{2, 1, 20},
			}
			wi := []WorkplaceInput{
				{"Teplice", wd1},
			}
			is := InputSummary{si, wi}
			exs1 := []StudentDetailWithCity{
				{1, 0, 20, "Most"},
				{3, 0, 22, "Most"},
				{5, 0, 21, "Teplice"},
			}
			exs2 := []StudentDetailWithCity{
				{2, 1, 21, "Most"},
				{4, 1, 20, "Teplice"},
			}
			exs := []StudentsByBranches{
				{0, exs1},
				{1, exs2},
			}
			exw1 := []WorkplaceDetailWithCity{
				{1, 0, 10, "Teplice"},
			}
			exw2 := []WorkplaceDetailWithCity{
				{2, 1, 20, "Teplice"},
			}
			exw := []WorkplacesByBranches{
				{0, exw1},
				{1, exw2},
			}

			rs, rw := TranslateInputs(is)
			Expect(rs).To(Equal(exs))
			Expect(rw).To(Equal(exw))
		})
		It("no students, no workplaces", func() {
			is := InputSummary{[]StudentInput{}, []WorkplaceInput{}}
			rs, rw := TranslateInputs(is)
			Expect(rs).To(BeEmpty())
			Expect(rw).To(BeEmpty())
		})
	})
	Context("GenerateSolution", func() {
		It("two branches, and all students", func() {
			sbb := []StudentsByBranches{
				{
					Branch: 1,
					Students: []StudentDetailWithCity{
						{10, 1, 100, "Most"},
						{11, 1, 101, "Teplice"},
					},
				}, {
					Branch: 2,
					Students: []StudentDetailWithCity{
						{20, 2, 100, "Most"},
						{21, 2, 101, "Teplice"},
					},
				},
			}
			wbb := []WorkplacesByBranches{
				{
					Branch: 1,
					Workplaces: []WorkplaceDetailWithCity{
						{1, 1, 2, "Teplice"},
						{2, 1, 0, "Most"},
					},
				}, {
					Branch: 2,
					Workplaces: []WorkplaceDetailWithCity{
						{101, 2, 2, "Teplice"},
						{102, 2, 2, "Most"},
					},
				},
			}
			sa, na := GenerateAssignement(sbb, wbb)
			exp := make(map[int][]int)
			exp[1] = []int{10, 11}
			exp[2] = []int{}
			Expect(len(sa)).To(Equal(2))
			Expect(sa[0].AssignementToWorkplace).To(Equal(exp))
			Expect(sa[1].AssignementToWorkplace).To(Not(BeEmpty()))
			Expect(na).To(BeEmpty())
		})
		It("one branch, and all students", func() {
			sbb := []StudentsByBranches{
				{
					Branch: 1,
					Students: []StudentDetailWithCity{
						{10, 1, 100, "Most"},
						{11, 1, 101, "Teplice"},
					},
				}, {
					Branch: 2,
					Students: []StudentDetailWithCity{
						{20, 2, 100, "Most"},
						{21, 2, 101, "Teplice"},
					},
				},
			}
			wbb := []WorkplacesByBranches{
				{
					Branch: 1,
					Workplaces: []WorkplaceDetailWithCity{
						{1, 1, 1, "Teplice"},
						{2, 1, 1, "Most"},
					},
				},
			}
			exp := []int{20, 21}

			sa, na := GenerateAssignement(sbb, wbb)
			Expect(len(sa)).To(Equal(1))
			Expect(len(sa[0].AssignementToWorkplace[1])).To(Equal(1))
			Expect(len(sa[0].AssignementToWorkplace[2])).To(Equal(1))
			Expect(na).To(Equal(exp))
		})
	})
	Context("Test calcFitness", func() {
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
		It("Calc without previous workplace - good one", func() {
			assignement := make(map[int][]int)
			assignement[1] = []int{1, 2}
			assignement[2] = []int{3}
			cities := []string{"Most", "Teplice", "Litvínov"}
			mockDistMatrix.EXPECT().GetDistance([]string{"Most, Czechia"}, []string{"Teplice, Czechia", "Litvínov, Czechia"}).Return([]byte(oneFcityTwoTcities))
			mockDistMatrix.EXPECT().GetDistance([]string{"Teplice, Czechia"}, []string{"Litvínov, Czechia"}).Return([]byte(oneFcityOneTcity))
			DM := rsa.NewDistanceMatrix(cities, dm)
			sa := StudentsAssignement{AssignementToWorkplace: assignement}
			sbb := []rsa.StudentsByBranches{
				{Students: []rsa.StudentDetailWithCity{
					{Id: 1, City: "Teplice"},
					{Id: 2, City: "Most"},
					{Id: 3, City: "Litvínov"},
				}},
			}
			wbb := []WorkplacesByBranches{
				{Workplaces: []rsa.WorkplaceDetailWithCity{
					{Id: 1, City: "Teplice"},
					{Id: 2, City: "Most"},
				}},
			}
			wrapper := PhaseAttributesWrapper{attr, DM, sbb, wbb}
			Expect(sa.CalulateFitness(wrapper)).To(Equal(28.48358499075553))
		})
		It("Calc without previous workplace - worst one", func() {
			assignement := make(map[int][]int)
			assignement[1] = []int{2}
			assignement[2] = []int{1, 3}
			cities := []string{"Most", "Teplice", "Litvínov"}
			mockDistMatrix.EXPECT().GetDistance([]string{"Most, Czechia"}, []string{"Teplice, Czechia", "Litvínov, Czechia"}).Return([]byte(oneFcityTwoTcities))
			mockDistMatrix.EXPECT().GetDistance([]string{"Teplice, Czechia"}, []string{"Litvínov, Czechia"}).Return([]byte(oneFcityOneTcity))
			DM := rsa.NewDistanceMatrix(cities, dm)
			sa := StudentsAssignement{AssignementToWorkplace: assignement}
			sbb := []rsa.StudentsByBranches{
				{Students: []rsa.StudentDetailWithCity{
					{Id: 1, City: "Teplice"},
					{Id: 2, City: "Most"},
					{Id: 3, City: "Litvínov"},
				}},
			}
			wbb := []WorkplacesByBranches{
				{Workplaces: []rsa.WorkplaceDetailWithCity{
					{Id: 1, City: "Teplice"},
					{Id: 2, City: "Most"},
				}},
			}
			wrapper := PhaseAttributesWrapper{attr, DM, sbb, wbb}
			Expect(math.Round(sa.CalulateFitness(wrapper)*10000) / 10000).To(Equal(math.Round(43.25099357570364*10000) / 10000))
		})
		It("Calc with previous workplace", func() {
			assignement := make(map[int][]int)
			assignement[1] = []int{1, 2}
			assignement[2] = []int{3}
			cities := []string{"Most", "Teplice", "Litvínov"}
			mockDistMatrix.EXPECT().GetDistance([]string{"Most, Czechia"}, []string{"Teplice, Czechia", "Litvínov, Czechia"}).Return([]byte(oneFcityTwoTcities))
			mockDistMatrix.EXPECT().GetDistance([]string{"Teplice, Czechia"}, []string{"Litvínov, Czechia"}).Return([]byte(oneFcityOneTcity))
			DM := rsa.NewDistanceMatrix(cities, dm)
			sa := StudentsAssignement{AssignementToWorkplace: assignement}
			sbb := []rsa.StudentsByBranches{
				{Students: []rsa.StudentDetailWithCity{
					{Id: 1, City: "Teplice", PreviousWorkplace: 1},
					{Id: 2, City: "Most"},
					{Id: 3, City: "Litvínov"},
				}},
			}
			wbb := []WorkplacesByBranches{
				{Workplaces: []rsa.WorkplaceDetailWithCity{
					{Id: 1, City: "Teplice"},
					{Id: 2, City: "Most"},
				}},
			}
			wrapper := PhaseAttributesWrapper{attr, DM, sbb, wbb}
			Expect(sa.CalulateFitness(wrapper)).To(Equal(78.48358499075553))
		})
	})
	Context("Test other functions", func() {
		It("Test IsCapacityReached", func() {
			assignement := make(map[int][]int)
			assignement[0] = []int{1, 2}
			assignement[1] = []int{10, 20}
			assignement[2] = []int{100, 200, 300}
			sa := StudentsAssignement{
				Branch:                 1,
				AssignementToWorkplace: assignement,
				IdList:                 []int{0, 1, 2},
				CapacityList:           []int{4, 3, 2},
			}
			Expect(sa.IsCapacityReached(0, 3)).To(BeFalse())
			Expect(sa.IsCapacityReached(1, 3)).To(BeTrue())
			Expect(sa.IsCapacityReached(2, 3)).To(BeTrue())
		})
		It("Test CollectAllCities", func() {
			students := []rsa.StudentInput{
				{
					City: "Most",
					Students: []rsa.StudentDetail{
						{
							Id:                1,
							Branch:            0,
							PreviousWorkplace: -1,
						},
						{
							Id:                2,
							Branch:            0,
							PreviousWorkplace: 1,
						},
					},
				},
				{
					City: "Teplice",
					Students: []rsa.StudentDetail{
						{
							Id:                2,
							Branch:            0,
							PreviousWorkplace: -1,
						},
						{
							Id:                3,
							Branch:            0,
							PreviousWorkplace: 3,
						},
					},
				},
				{
					City: "Ústí nad Labem",
					Students: []rsa.StudentDetail{
						{
							Id:                4,
							Branch:            0,
							PreviousWorkplace: 4,
						},
					},
				},
			}
			workplaces := []rsa.WorkplaceInput{
				{
					City: "Most",
					Workplaces: []rsa.WorkplaceDetail{
						{
							Id:       1,
							Branch:   0,
							Capacity: 2,
						},
						{
							Id:       2,
							Branch:   0,
							Capacity: 1,
						},
					},
				},
				{
					City: "Teplice",
					Workplaces: []rsa.WorkplaceDetail{
						{
							Id:       3,
							Branch:   0,
							Capacity: 2,
						},
					},
				},
				{
					City: "Litvínov",
					Workplaces: []rsa.WorkplaceDetail{
						{
							Id:       4,
							Branch:   0,
							Capacity: 2,
						},
					},
				},
			}
			inputs := rsa.InputSummary{
				Students:   students,
				Workplaces: workplaces,
			}
			expected := rsa.CollectAllCities(inputs)
			Expect(expected).To(ConsistOf("Litvínov", "Teplice", "Most", "Ústí nad Labem"))
			Expect(len(expected)).To(Equal(4))
		})
		It("Test GetStudentIndexFromNonEmptyWorkplace", func() {
			assignement := make(map[int][]int)
			assignement[0] = []int{1, 2}
			assignement[1] = []int{10, 20}
			assignement[2] = []int{}
			sa := StudentsAssignement{
				Branch:                 1,
				AssignementToWorkplace: assignement,
				IdList:                 []int{0, 1, 2},
				CapacityList:           []int{4, 3, 2},
			}
			wrkpl, std := sa.GetStudentIndexFromNonEmptyWorkplace()
			Expect(wrkpl).NotTo(Equal(-1))
			Expect(wrkpl).To(BeElementOf(0, 1))
			Expect(std).NotTo(Equal(-1))
			Expect(std).To(BeElementOf(0, 1))

			delete(assignement, 0)
			delete(assignement, 1)
			sa.AssignementToWorkplace = assignement
			wrkpl, std = sa.GetStudentIndexFromNonEmptyWorkplace()
			Expect(wrkpl).To(Equal(-1))
			Expect(std).To(Equal(-1))
		})
	})
})
