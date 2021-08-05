package reapp_students_assignement

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// generateRandomWorkplace select randomly workspace what has enough capacity to
// absorb another student
func generateRandomWorkplace(sa StudentsAssignement) int {
	rand.Seed(time.Now().UnixNano())
	nonAcceptible := []int{}
	for {
		if (len(nonAcceptible) >= len(sa.AssignementToWorkplace)) && len(sa.AssignementToWorkplace) != 0 {
			return -1
		}

		randWorkplaceIndx := rand.Intn(len(sa.IdList))
		if findInt(nonAcceptible, sa.IdList[randWorkplaceIndx]) {
			continue
		}

		if _, ok := sa.AssignementToWorkplace[sa.IdList[randWorkplaceIndx]]; !ok && sa.CapacityList[randWorkplaceIndx] > 0 {
			return sa.IdList[randWorkplaceIndx]
		}
		if len(sa.AssignementToWorkplace[sa.IdList[randWorkplaceIndx]]) < sa.CapacityList[randWorkplaceIndx] {
			return sa.IdList[randWorkplaceIndx]
		}
		nonAcceptible = append(nonAcceptible, sa.IdList[randWorkplaceIndx])
	}
}

func findInt(a []int, i int) bool {
	for _, j := range a {
		if j == i {
			return true
		}
	}
	return false
}

func findIndex(a []int, i int) int {
	for index, j := range a {
		if j == i {
			return index
		}
	}
	return -1
}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func mapStringSlice(ss []string, fn func(string) bool) []string {
	var wrong []string
	for _, s := range ss {
		if !fn(s) {
			wrong = append(wrong, s)
		}
	}
	return wrong
}

func isInSlice(ss []string, str string) bool {
	for _, s := range ss {
		if s == str {
			return true
		}
	}
	return false
}

func replaceWhitespaces(s string) string {
	return strings.ReplaceAll(s, " ", "+")
}

func stripString(s, sep string) string {
	if strings.Contains(s, sep) {
		return s[0:strings.Index(s, sep)]
	}
	return s
}

func stringToHex(s string) string {
	return hex.EncodeToString([]byte(s))
}

func hexToString(s string) (string, error) {
	if decoded, err := hex.DecodeString(s); err != nil {
		return "", err
	} else {
		return string(decoded), nil
	}
}

func prepareMatrixKey(o, d string) string {
	return fmt.Sprintf("%s%s%s", o, DistanceMatrixKeySeparator, d)
}

func appendToDistanceMatrix(orig, new map[string]Distance) map[string]Distance {
	wholeDM := make(map[string]Distance)
	for city2city, val := range orig {
		wholeDM[city2city] = val
	}
	for city2city, val := range new {
		wholeDM[city2city] = val
	}
	return wholeDM
}

func findMinMaxFitness(solutions []ProposedSolution) (float64, float64) {
	var min, max float64
	for _, solution := range solutions {
		if solution.Fitness < min {
			min = solution.Fitness
		} else if solution.Fitness > max {
			max = solution.Fitness
		}
	}
	return min, max
}
