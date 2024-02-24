package iec61850

// #include <iec61850_client.h>
import "C"

type FC int

// fc types
const (
	// ST Status information
	ST FC = iota
	// MX Measurands - analogue values
	MX
	// SP Setpoint
	SP
	// SV Substitution
	SV
	// CF Configuration
	CF
	// DC Description
	DC
	// SG Setting group
	SG
	// SE Setting group editable
	SE
	// SR service response / service tracking
	SR
	// OR Operate received
	OR
	// BL Blocking
	BL
	// EX Extended definition
	EX
	// CO Control, deprecated but kept here for backward compatibility
	CO
	// RP Unbuffered Reporting
	RP
	// BR Buffered Reporting
	BR
	// ALL All FCs - wildcard value
	ALL  FC = 99
	NONE FC = -1
)
