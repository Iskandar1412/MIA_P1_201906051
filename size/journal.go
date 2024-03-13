package size

import (
	"MIA_P1_201906051/structures"
	"unsafe"
)

func SizeJournal() int32 { //156 bytes
	a01 := int32(unsafe.Sizeof(structures.Journal{}.J_state))
	a02 := int32(unsafe.Sizeof(structures.Journal{}.J_command))
	a03 := int32(unsafe.Sizeof(structures.Journal{}.J_date))
	a04 := int32(unsafe.Sizeof(structures.Journal{}.J_content))
	result := a01 + a02 + a03 + a04
	return result
}
