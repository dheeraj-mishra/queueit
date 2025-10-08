package models

// Priorities:
const (
	PRIORITY_HIGH   = 1
	PRIORITY_MEDIUM = 2
	PRIORITY_LOW    = 3
)

// Statuses:
const (
	STATUS_PENDING  = 1
	STATUS_WIP      = 2
	STATUS_DONE     = 3
	STATUS_ARCHIVED = 4
)

var ValidStatuses = map[int]bool{
	STATUS_PENDING:  true,
	STATUS_WIP:      true,
	STATUS_DONE:     true,
	STATUS_ARCHIVED: true,
}

var ValidPriorities = map[int]bool{
	PRIORITY_HIGH:   true,
	PRIORITY_MEDIUM: true,
	PRIORITY_LOW:    true,
}
