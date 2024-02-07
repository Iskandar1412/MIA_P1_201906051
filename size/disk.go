package size

import (
	"MIA_P1_201906051/structures"
	"unsafe"
)

func SizeEBR() int64 {
	a01 := unsafe.Sizeof(structures.EBR{}.Part_mount)
	a02 := unsafe.Sizeof(structures.EBR{}.Part_fit)
	a03 := unsafe.Sizeof(structures.EBR{}.Part_start)
	a04 := unsafe.Sizeof(structures.EBR{}.Part_s)
	a05 := unsafe.Sizeof(structures.EBR{}.Part_next)
	a06 := unsafe.Sizeof(structures.EBR{}.Name)
	result := a01 + a02 + a03 + a04 + a05 + a06
	return int64(result)
}

func SizePartition() int64 {
	a01 := unsafe.Sizeof(structures.Partition{}.Part_status)
	a02 := unsafe.Sizeof(structures.Partition{}.Part_type)
	a03 := unsafe.Sizeof(structures.Partition{}.Part_fit)
	a04 := unsafe.Sizeof(structures.Partition{}.Part_start)
	a05 := unsafe.Sizeof(structures.Partition{}.Part_s)
	a06 := unsafe.Sizeof(structures.Partition{}.Part_name)
	a07 := unsafe.Sizeof(structures.Partition{}.Part_correlative)
	a08 := unsafe.Sizeof(structures.Partition{}.Part_id)
	result := a01 + a02 + a03 + a04 + a05 + a06 + a07 + a08
	return int64(result)
}

func SizeMBR() int64 {
	a01 := unsafe.Sizeof(structures.MBR{}.Mbr_tamano)
	a02 := unsafe.Sizeof(structures.MBR{}.Mbr_fecha_creacion)
	a03 := unsafe.Sizeof(structures.MBR{}.Mbr_disk_signature)
	a04 := unsafe.Sizeof(structures.MBR{}.Dsk_fit)
	a05 := SizePartition() * 4
	result := int64(a01+a02+a03+a04) + a05
	return result
}
