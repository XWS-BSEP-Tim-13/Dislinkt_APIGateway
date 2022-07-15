package enum

type EmploymentType int

const (
	FULL_TIME EmploymentType = iota
	PART_TIME
	INTERNSHIP
	NONE
)
