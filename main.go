package main

import (
	"fmt"
	rsa "reapp_students_assignement/pkg"
)

// main is not copied to GCP FaaS and so you can place here whatever you want to test localy
func main() {
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
					Id:                3,
					Branch:            0,
					PreviousWorkplace: -1,
				},
				{
					Id:                4,
					Branch:            0,
					PreviousWorkplace: 3,
				},
			},
		},
		{
			City: "Ústí nad Labem",
			Students: []rsa.StudentDetail{
				{
					Id:                5,
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
	attributes := rsa.ABCAttributes{ReferenceMaxDistance: 60000, PenaltyForTheSameWorkplace: 50, BeeColonySize: 1000, MaxGenerations: 100}
	allCities := rsa.CollectAllCities(inputs)
	distanceMatrix := rsa.NewDistanceMatrix(allCities, rsa.GMapGetter{})
	fmt.Print(rsa.ABCrun(attributes, inputs, distanceMatrix))
}
