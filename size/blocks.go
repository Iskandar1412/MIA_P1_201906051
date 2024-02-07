package size

import (
	"MIA_P1_201906051/structures"
	"unsafe"
)

func SizeContent() int64 {
	a01 := unsafe.Sizeof(structures.Content{}.B_name)
	a02 := unsafe.Sizeof(structures.Content{}.B_inodo)
	result := a01 + a02
	return int64(result)
}

func SizeBlockCarpetas() int64 {
	a01 := unsafe.Sizeof(structures.BlockCarpetas{}.B_content[0].B_name) + unsafe.Sizeof(structures.BlockCarpetas{}.B_content[0].B_inodo)
	a02 := unsafe.Sizeof(structures.BlockCarpetas{}.B_content[1].B_name) + unsafe.Sizeof(structures.BlockCarpetas{}.B_content[1].B_inodo)
	a03 := unsafe.Sizeof(structures.BlockCarpetas{}.B_content[2].B_name) + unsafe.Sizeof(structures.BlockCarpetas{}.B_content[2].B_inodo)
	a04 := unsafe.Sizeof(structures.BlockCarpetas{}.B_content[3].B_name) + unsafe.Sizeof(structures.BlockCarpetas{}.B_content[3].B_inodo)
	result := a01 + a02 + a03 + a04
	return int64(result)
}

func SizeBloqueArchivos() int64 {
	a01 := unsafe.Sizeof(structures.BloqueArchivos{}.B_content)
	return int64(a01)
}

func SizeBloqueApuntadores() int64 {
	a01 := unsafe.Sizeof(structures.BloqueApuntadores{}.B_pointers)
	return int64(a01)
}
