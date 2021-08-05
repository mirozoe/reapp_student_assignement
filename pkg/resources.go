package reapp_students_assignement

// InputSummary is input data structure what is then used
// for calculating combination model
type InputSummary struct {
	Students   []StudentInput
	Workplaces []WorkplaceInput
}

type StudentInput struct {
	City     string
	Students []StudentDetail
}

type WorkplaceInput struct {
	City       string
	Workplaces []WorkplaceDetail
}

type StudentDetail struct {
	Id                int
	Branch            int
	PreviousWorkplace int
}

type WorkplaceDetail struct {
	Id       int
	Branch   int
	Capacity int
}

// StudentsByBranches translated structure representing students by branch
type StudentsByBranches struct {
	Branch   int
	Students []StudentDetailWithCity
}

// WorkplacesByBranches translated structure representing workplaces
// by branch
type WorkplacesByBranches struct {
	Branch     int
	Workplaces []WorkplaceDetailWithCity
}

type StudentDetailWithCity struct {
	Id                int
	Branch            int
	PreviousWorkplace int
	City              string
}

type WorkplaceDetailWithCity struct {
	Id       int
	Branch   int
	Capacity int
	City     string
}

// ABCAttributes
type ABCAttributes struct {
	// ReferenceMaxDistance is maximal awailable distance for student to commute
	// to workplace (100 km)
	ReferenceMaxDistance int64

	// PenaltyForTheSameWorkplace is constant what is added to fitness
	// when student is assigned to workplace where was previously
	PenaltyForTheSameWorkplace int

	// BeeColonySize nuber of bees
	BeeColonySize int

	// MaxGenerations max nuber of colony generations
	MaxGenerations int
}

// StudentsAssignement is output of ABC for specific branch
type StudentsAssignement struct {
	// branch
	Branch int
	// key is workplace ID and value is slice of student IDs assigned to it
	AssignementToWorkplace map[int][]int
	// slice of workplace IDs with specific branch
	IdList []int
	// corresponding capacity to workplace
	CapacityList []int
}

// ProposedSolution is container for assigned and unassigned students
type ProposedSolution struct {
	Assignements            []StudentsAssignement
	UnassignedStudents      []int
	Fitness                 float64
	CycleWithoutImprovement int
}

// Structs and interfaces for distance

type DistanceMatrix struct {
	DistanceInterface
	distanceMatrix map[string]Distance
	cities         []string
}

// DistanceInterface generic interface what provides connection to distance data source
type DistanceInterface interface {
	GetDistance([]string, []string) []byte
}

// DistanceGetter is wrapper arround DistanceInterface and provides mockable resource
type DistanceGetter struct {
	DistanceInterface
}

// GMapGetter is specific DistanceInterface implementation
type GMapGetter struct{}

// GMapResponse describes JSON struct response from GMap
type GMapResponse struct {
	DestinationAddresses []string      `json:"destination_addresses"`
	OriginAddreses       []string      `json:"origin_addresses"`
	Rows                 []GMapRespRow `json:"rows"`
}

// GMapRespRow is row in GMap response JSON
type GMapRespRow struct {
	Elements []GMapRespElement `json:"elements"`
}

// GMapRespElement is an element in GMap response JSON
type GMapRespElement struct {
	Distance GMapRespElementStruct `json:"distance"`
	Duration GMapRespElementStruct `json:"duration"`
	Status   string                `json:"status"`
}

// GMapRespElementStruct is detail from element in GMap response JSON
type GMapRespElementStruct struct {
	Text  string `json:"text"`
	Value int64  `json:"value"`
}

// Distance is distance between two cities
type Distance struct {
	// DistanceM distance in meters
	DistanceM int64
	// Duration is time in seconds to drive between two cities
	DurationS int64
}

// PhaseAttributesWrapper wraps all attributes what are needed in scout, employee and onlooker phases
type PhaseAttributesWrapper struct {
	// AbcAttrs are BeeColony settings
	AbcAttrs ABCAttributes
	// Dm is not Dungeon Master but distance matrix
	Dm DistanceMatrix
	// Sbb students sorted by branches
	Sbb []StudentsByBranches
	// Wbb workplaces sorted by branches
	Wbb []WorkplacesByBranches
}
