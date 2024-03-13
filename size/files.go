package size

import (
	"MIA_P1_201906051/structures"
	"unsafe"
)

func SizeSuperBlock() int32 { //68 bytes
	a01 := unsafe.Sizeof(structures.SuperBlock{}.S_filesistem_type)
	a02 := unsafe.Sizeof(structures.SuperBlock{}.S_inodes_count)
	a03 := unsafe.Sizeof(structures.SuperBlock{}.S_blocks_count)
	a04 := unsafe.Sizeof(structures.SuperBlock{}.S_free_blocks_count)
	a05 := unsafe.Sizeof(structures.SuperBlock{}.S_free_inodes_count)
	a06 := unsafe.Sizeof(structures.SuperBlock{}.S_mtime)
	a07 := unsafe.Sizeof(structures.SuperBlock{}.S_umtime)
	a08 := unsafe.Sizeof(structures.SuperBlock{}.S_mnt_count)
	a09 := unsafe.Sizeof(structures.SuperBlock{}.S_magic)
	a10 := unsafe.Sizeof(structures.SuperBlock{}.S_inode_s)
	a11 := unsafe.Sizeof(structures.SuperBlock{}.S_block_s)
	a12 := unsafe.Sizeof(structures.SuperBlock{}.S_first_ino)
	a13 := unsafe.Sizeof(structures.SuperBlock{}.S_first_blo)
	a14 := unsafe.Sizeof(structures.SuperBlock{}.S_bm_inode_start)
	a15 := unsafe.Sizeof(structures.SuperBlock{}.S_bm_block_start)
	a16 := unsafe.Sizeof(structures.SuperBlock{}.S_inode_start)
	a17 := unsafe.Sizeof(structures.SuperBlock{}.S_block_start)
	result := a01 + a02 + a03 + a04 + a05 + a06 + a07 + a08 + a09 + a10 + a11 + a12 + a13 + a14 + a15 + a16 + a17
	return int32(result)
}

func SizeInode() int32 { //92 bytes
	a01 := unsafe.Sizeof(structures.Inode{}.I_uid)
	a02 := unsafe.Sizeof(structures.Inode{}.I_gid)
	a03 := unsafe.Sizeof(structures.Inode{}.I_s)
	a04 := unsafe.Sizeof(structures.Inode{}.I_atime)
	a05 := unsafe.Sizeof(structures.Inode{}.I_ctime)
	a06 := unsafe.Sizeof(structures.Inode{}.I_mtime)
	a07 := unsafe.Sizeof(structures.Inode{}.I_block)
	a08 := unsafe.Sizeof(structures.Inode{}.I_type)
	a09 := unsafe.Sizeof(structures.Inode{}.I_perm)
	result := a01 + a02 + a03 + a04 + a05 + a06 + a07 + a08 + a09
	return int32(result)
}
