package size

import (
	"MIA_P1_201906051/structures"
	"unsafe"
)

func SizeJournal() int64 {
	a01 := int64(unsafe.Sizeof(structures.Journal{}.J_operation))
	a02 := int64(unsafe.Sizeof(structures.Journal{}.J_content))
	a03 := int64(unsafe.Sizeof(structures.Journal{}.J_permissions))
	a04 := int64(unsafe.Sizeof(structures.Journal{}.J_name))
	a05 := int64(unsafe.Sizeof(structures.Journal{}.J_owner))
	a06 := int64(unsafe.Sizeof(structures.Journal{}.J_date))
	a07 := int64(unsafe.Sizeof(structures.Journal{}.J_type))
	a08 := int64(unsafe.Sizeof(structures.Journal{}.J_size))
	result := a01 + a02 + a03 + a04 + a05 + a06 + a07 + a08
	return result
}
