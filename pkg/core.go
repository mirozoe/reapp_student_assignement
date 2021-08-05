package reapp_students_assignement

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// ABCrun runs ABC algorithm and coordinates it
func ABCrun(attr ABCAttributes, inputs InputSummary, dm DistanceMatrix) ProposedSolution {
	rand.Seed(time.Now().UnixNano())
	log.Info().Msg("Starting ABC algorithm")
	bestSolution := ProposedSolution{Fitness: -1}
	sbb, wbb := TranslateInputs(inputs)
	paw := PhaseAttributesWrapper{
		AbcAttrs: attr,
		Dm:       dm,
		Sbb:      sbb,
		Wbb:      wbb,
	}
	// init all swarm with random solutions
	swarm := scoutPhase(attr.BeeColonySize, []ProposedSolution{}, paw)
	for i := 1; i <= attr.MaxGenerations; i++ {
		swarm = employeePhase(swarm, paw)
		swarm = scoutPhase(attr.BeeColonySize, swarm, paw)
		swarm = onlookerPhase(swarm, paw)
		bestSolution = findTheBestSolution(bestSolution, swarm)
		log.Debug().Msgf("Best solution is: %+v\n", bestSolution)
	}
	return bestSolution
}

// GenerateAssignement generates random solution for one bee
func GenerateAssignement(sbb []StudentsByBranches, wbb []WorkplacesByBranches) ([]StudentsAssignement, []int) {
	sa := []StudentsAssignement{}
	unassignedStudents := []int{}
	for _, studentsBranch := range sbb {
		index := findWorkplaceByBranches(wbb, studentsBranch.Branch, func(w WorkplacesByBranches, b int) bool {
			return w.Branch == b
		})
		if index == -1 {
			unassignedStudents = append(unassignedStudents, getStudentsId(studentsBranch.Students)...)
			continue
		}

		saindex := findStudentsAssignement(sa, studentsBranch.Branch, func(s StudentsAssignement, b int) bool {
			return s.Branch == b
		})
		if saindex == -1 {
			idList, capList := prepareAllWorkplacesLists(wbb[index].Workplaces)
			new := StudentsAssignement{
				Branch:                 studentsBranch.Branch,
				AssignementToWorkplace: prepareAllWorkplacesMap(wbb[index].Workplaces),
				IdList:                 idList,
				CapacityList:           capList,
			}
			sa = append(sa, new)
			saindex = len(sa) - 1
		}

		for _, student := range studentsBranch.Students {
			windex := selectWorkplace(wbb[index].Workplaces, func(w WorkplaceDetailWithCity) int {
				if w.Branch == studentsBranch.Branch {
					if workpl := generateRandomWorkplace(sa[saindex]); workpl != -1 {
						return workpl
					}
				}
				return -1
			})
			if windex == -1 {
				unassignedStudents = append(unassignedStudents, student.Id)
			} else {
				sa[saindex].AssignementToWorkplace[windex] = append(sa[saindex].AssignementToWorkplace[windex], student.Id)
			}
		}
	}
	return sa, unassignedStudents
}

// TranslateInputs is responsible for prepare better structures from InputSummary object
func TranslateInputs(is InputSummary) ([]StudentsByBranches, []WorkplacesByBranches) {
	var sbb []StudentsByBranches
	for _, s := range is.Students {
		for _, sd := range s.Students {
			sdwc := StudentDetailWithCity{
				Id:                sd.Id,
				Branch:            sd.Branch,
				PreviousWorkplace: sd.PreviousWorkplace,
				City:              s.City,
			}
			translateStudentToSBB(&sbb, sdwc)
		}
	}

	var wbb []WorkplacesByBranches
	for _, w := range is.Workplaces {
		for _, wd := range w.Workplaces {
			wdwc := WorkplaceDetailWithCity{
				Id:       wd.Id,
				Branch:   wd.Branch,
				Capacity: wd.Capacity,
				City:     w.City,
			}
			translateWorkplaceToWBB(&wbb, wdwc)
		}
	}
	return sbb, wbb
}

// CollectAllCities returns only unique cities from InputSummary
func CollectAllCities(is InputSummary) []string {
	rawCities := []string{}
	for _, student := range is.Students {
		if len(rawCities) == 0 {
			rawCities = append(rawCities, student.City)
		} else {
			if !isInSlice(rawCities, student.City) {
				rawCities = append(rawCities, student.City)
			}
		}
	}
	for _, workplace := range is.Workplaces {
		if !isInSlice(rawCities, workplace.City) {
			rawCities = append(rawCities, workplace.City)
		}
	}
	return rawCities
}

func translateStudentToSBB(sbb *[]StudentsByBranches, sdwc StudentDetailWithCity) {
	indx := findStudentByBranches(*sbb, sdwc.Branch, func(s StudentsByBranches, b int) bool {
		return s.Branch == b
	})
	if indx == -1 {
		s := StudentsByBranches{
			Branch:   sdwc.Branch,
			Students: []StudentDetailWithCity{sdwc},
		}
		(*sbb) = append(*sbb, s)
	} else {
		(*sbb)[indx].Students = append((*sbb)[indx].Students, sdwc)
	}
}

func translateWorkplaceToWBB(wbb *[]WorkplacesByBranches, wdwc WorkplaceDetailWithCity) {
	indx := findWorkplaceByBranches(*wbb, wdwc.Branch, func(w WorkplacesByBranches, b int) bool {
		return w.Branch == b
	})
	if indx == -1 {
		w := WorkplacesByBranches{
			Branch:     wdwc.Branch,
			Workplaces: []WorkplaceDetailWithCity{wdwc},
		}
		(*wbb) = append(*wbb, w)
	} else {
		(*wbb)[indx].Workplaces = append((*wbb)[indx].Workplaces, wdwc)
	}
}

func findStudentByBranches(sbb []StudentsByBranches, branch int, fn func(StudentsByBranches, int) bool) int {
	for i, s := range sbb {
		if ret := fn(s, branch); ret {
			return i
		}
	}
	return -1
}

func findWorkplaceByBranches(wbb []WorkplacesByBranches, branch int, fn func(WorkplacesByBranches, int) bool) int {
	for i, w := range wbb {
		if ret := fn(w, branch); ret {
			return i
		}
	}
	return -1
}

func findStudentsAssignement(sa []StudentsAssignement, branch int, fn func(StudentsAssignement, int) bool) int {
	for i, w := range sa {
		if ret := fn(w, branch); ret {
			return i
		}
	}
	return -1
}

func selectWorkplace(wd []WorkplaceDetailWithCity, fn func(w WorkplaceDetailWithCity) int) int {
	for _, w := range wd {
		return fn(w)
	}
	return -1
}

func getStudentsId(sd []StudentDetailWithCity) []int {
	ids := []int{}
	for _, s := range sd {
		ids = append(ids, s.Id)
	}
	return ids
}

func prepareAllWorkplacesLists(workplaces []WorkplaceDetailWithCity) ([]int, []int) {
	idList := []int{}
	capList := []int{}
	for _, w := range workplaces {
		idList = append(idList, w.Id)
		capList = append(capList, w.Capacity)
	}
	return idList, capList
}

func prepareAllWorkplacesMap(workplaces []WorkplaceDetailWithCity) map[int][]int {
	workplMap := make(map[int][]int)
	for _, w := range workplaces {
		workplMap[w.Id] = []int{}
	}
	return workplMap
}

// findTheBestSolution be aware that lower fitness is better
func findTheBestSolution(oldBest ProposedSolution, solutions []ProposedSolution) ProposedSolution {
	best := ProposedSolution{}
	if oldBest.Fitness == -1 {
		best.Fitness = 999999.9
	} else {
		best = oldBest
	}
	for _, solution := range solutions {
		if solution.Fitness < best.Fitness {
			best = solution
		}
	}
	return best
}

// CalculateFitness computes fitness for all students on all workplaces
// Hint: smaller fitness is better
func (sa StudentsAssignement) CalulateFitness(wrapper PhaseAttributesWrapper) float64 {
	var overallFitness float64
	for workplaceId, students := range sa.AssignementToWorkplace {
		cityW, err := findCityByWorkplaceId(wrapper.Wbb, workplaceId)
		if err != nil {
			log.Error().Err(err)
			return -1
		}
		for _, student := range students {
			studentDetail, err := findByStudentId(wrapper.Sbb, student)
			if err != nil {
				log.Error().Err(err)
				continue
			}
			overallFitness += calcStudentFitness(wrapper.AbcAttrs, wrapper.Dm, cityW, studentDetail, workplaceId)
		}
	}
	return overallFitness
}

// IsCapacityReached checks if workplace capacity is already reached or not
func (sa StudentsAssignement) IsCapacityReached(workplace, desired int) (bool, error) {
	if i := findIndex(sa.IdList, workplace); i > -1 {
		if sa.CapacityList[i] > desired {
			return false, nil
		}
		return true, nil
	}
	return false, ErrorFactory(WorkplaceNotFound, strconv.Itoa(workplace))
}

// GetStudentIndexFromNonEmptyWorkplace returns random workplace and student position in actual assignement
func (sa StudentsAssignement) GetStudentIndexFromNonEmptyWorkplace() (int, int) {
	for i := 0; i < len(sa.AssignementToWorkplace); i++ {
		workplacePosition := selectRandomFromMap(sa.AssignementToWorkplace)
		if len(sa.AssignementToWorkplace[workplacePosition]) != 0 {
			return workplacePosition, rand.Intn(len(sa.AssignementToWorkplace[workplacePosition]))
		}
	}
	return -1, -1
}

func calcStudentFitness(attr ABCAttributes, dm DistanceMatrix, cityW string, student StudentDetailWithCity, newWorkplaceId int) float64 {
	distanceFitness := dm.CalculateLogDistance(attr, cityW, student.City)
	if newWorkplaceId == student.PreviousWorkplace {
		distanceFitness += float64(attr.PenaltyForTheSameWorkplace)
	}
	return distanceFitness
}

func calculateAssignementsFitness(assignements []StudentsAssignement, wrapper PhaseAttributesWrapper) float64 {
	fitness := float64(0)
	for _, assignement := range assignements {
		fitness += assignement.CalulateFitness(wrapper)
	}
	return fitness
}

func (ps ProposedSolution) Copy() ProposedSolution {
	a := make([]StudentsAssignement, len(ps.Assignements))
	copy(a, ps.Assignements)
	us := make([]int, len(ps.UnassignedStudents))
	copy(us, ps.UnassignedStudents)
	return ProposedSolution{
		Assignements:            a,
		UnassignedStudents:      us,
		Fitness:                 ps.Fitness,
		CycleWithoutImprovement: 0,
	}
}

func findByStudentId(sbbs []StudentsByBranches, studentId int) (StudentDetailWithCity, error) {
	for _, sbb := range sbbs {
		for _, student := range sbb.Students {
			if student.Id == studentId {
				return student, nil
			}
		}
	}
	return StudentDetailWithCity{}, ErrorFactory(CityDidntFound, strconv.Itoa(studentId))
}

func findCityByWorkplaceId(wbbs []WorkplacesByBranches, workplaceId int) (string, error) {
	for _, wbb := range wbbs {
		for _, workplace := range wbb.Workplaces {
			if workplace.Id == workplaceId {
				return workplace.City, nil
			}
		}
	}
	return "", ErrorFactory(CityDidntFound, strconv.Itoa(workplaceId))
}

func selectRandomFromMap(m map[int][]int) int {
	var key int
	r := rand.Intn(len(m))
	i := 1
	for k, _ := range m {
		if i == r {
			return k
		}
		key = k
		i++
	}
	return key
}

func isAtrracted(fitness, min, max float64) bool {
	var percent float64
	if max-min < 0 {
		return false
	}
	if fitness == min {
		percent = 90.0
	} else if fitness == max {
		percent = 10.0
	} else {
		percent = 0.90 - (((fitness - min) * 0.80) / (max - min))
	}
	if rand.Float64() < percent {
		return true
	}
	return false
}

func searchNeighborSolution(solution ProposedSolution, wrapper PhaseAttributesWrapper) ProposedSolution {
	newSolution := solution.Copy()
	exchangeNo := rand.Intn(2)
	for i := 0; i <= exchangeNo; i++ {
		branchPosition := rand.Intn(len(newSolution.Assignements))
		specBranch := newSolution.Assignements[branchPosition]
		workplacePosition, studentPosition := specBranch.GetStudentIndexFromNonEmptyWorkplace()
		if workplacePosition == -1 || studentPosition == -1 {
			continue
		}
		toWorkplacePosition := selectRandomFromMap(specBranch.AssignementToWorkplace)
		reached, err := specBranch.IsCapacityReached(toWorkplacePosition, len(specBranch.AssignementToWorkplace[toWorkplacePosition])+1)
		if err != nil {
			panic(err)
		}

		if reached {
			var swapStudentPosition int
			if len(specBranch.AssignementToWorkplace[toWorkplacePosition]) > 0 {
				swapStudentPosition = rand.Intn(len(specBranch.AssignementToWorkplace[toWorkplacePosition]))
				studentToSwap := specBranch.AssignementToWorkplace[toWorkplacePosition][swapStudentPosition]
				specBranch.AssignementToWorkplace[toWorkplacePosition][swapStudentPosition] = specBranch.AssignementToWorkplace[workplacePosition][studentPosition]
				specBranch.AssignementToWorkplace[workplacePosition][studentPosition] = studentToSwap
			}
		} else {
			specBranch.AssignementToWorkplace[toWorkplacePosition] = append(specBranch.AssignementToWorkplace[toWorkplacePosition], specBranch.AssignementToWorkplace[workplacePosition][studentPosition])
			remove(specBranch.AssignementToWorkplace[workplacePosition], studentPosition)
		}
		newSolution.Fitness = calculateAssignementsFitness(newSolution.Assignements, wrapper)
		newSolution.CycleWithoutImprovement = 0
	}
	return newSolution
}

func solutionGreedySelection(solution ProposedSolution, wrapper PhaseAttributesWrapper) ProposedSolution {
	neighborSolution := searchNeighborSolution(solution, wrapper)
	if neighborSolution.Fitness < solution.Fitness {
		return neighborSolution
	} else {
		solution.CycleWithoutImprovement += 1
		return solution
	}
}

func scoutPhase(colonySize int, solutions []ProposedSolution, wrapper PhaseAttributesWrapper) []ProposedSolution {
	if len(solutions) == 0 {
		var newSolutions []ProposedSolution
		for i := 0; i < colonySize; i++ {
			assignements, unassignedStudents := GenerateAssignement(wrapper.Sbb, wrapper.Wbb)
			fitness := calculateAssignementsFitness(assignements, wrapper)
			newSolutions = append(newSolutions, ProposedSolution{
				Assignements:            assignements,
				UnassignedStudents:      unassignedStudents,
				Fitness:                 fitness,
				CycleWithoutImprovement: 0,
			})
		}
		return newSolutions
	}
	newSolutions := make([]ProposedSolution, colonySize)
	copy(newSolutions, solutions)
	for index, solution := range solutions {
		if solution.CycleWithoutImprovement >= EmployeeBeeMaxCycleNo {
			assignements, unassignedStudents := GenerateAssignement(wrapper.Sbb, wrapper.Wbb)
			fitness := calculateAssignementsFitness(assignements, wrapper)
			newSolutions[index] = ProposedSolution{
				Assignements:            assignements,
				UnassignedStudents:      unassignedStudents,
				Fitness:                 fitness,
				CycleWithoutImprovement: 0,
			}
		}
	}
	return newSolutions
}

func employeePhase(solutions []ProposedSolution, wrapper PhaseAttributesWrapper) []ProposedSolution {
	newSolutions := make([]ProposedSolution, len(solutions))
	copy(newSolutions, solutions)
	for index, solution := range newSolutions {
		newSolutions[index] = solutionGreedySelection(solution, wrapper)
	}
	return newSolutions
}

func onlookerPhase(solutions []ProposedSolution, wrapper PhaseAttributesWrapper) []ProposedSolution {
	newSolutions := make([]ProposedSolution, len(solutions))
	copy(newSolutions, solutions)
	min, max := findMinMaxFitness(solutions)
	for onlookerBeer := 0; onlookerBeer < len(solutions); onlookerBeer++ {
		for index, solution := range newSolutions {
			if isAtrracted(solution.Fitness, min, max) {
				newSolutions[index] = solutionGreedySelection(solution, wrapper)
				break
			}
		}
	}
	return newSolutions
}
