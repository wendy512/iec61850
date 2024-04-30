package iec61850

type DataModel struct {
	LDs []LD
}

type LD struct {
	Data string
	LNs  []LN
}

type LN struct {
	Data      string
	DOs       []DO
	DSs       []DS
	URReports []URReport
	BRReports []BRReport
}

type URReport struct {
	Data string
}

type BRReport struct {
	Data string
}

type DS struct {
	Data   string
	DSRefs []DSRef
}

type DSRef struct {
	Data string
}

type DO struct {
	Data string
	DAs  []DA
}

type DA struct {
	Data string
	DAs  []DA
}
