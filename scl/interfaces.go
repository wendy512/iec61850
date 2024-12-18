package scl

type SclType interface {
	GetId() string
	GetDesc() string
	GetUsed() bool
	SetUsed(used bool)
}

type DataModelNode interface {
	GetName() string
	GetSclType() SclType
	GetChildByName(childName string) DataModelNode
}
