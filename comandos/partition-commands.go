package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func Values_FDISK(instructions []string) (int32, byte, [16]byte, byte, byte, byte, string, int32) {
	var _size int32
	var _driveletter byte
	var _name [16]byte
	var _unit byte = 'K'
	var _type byte = 'P'
	var _fit byte = 'W'
	var _delete string = "None"
	var _add int32 = 0
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "size") {
			var value = TieneSize("FDISK", valor)
			_size = value
		} else if strings.HasPrefix(strings.ToLower(valor), "driveletter") {
			var value = TieneDriveLetter("FDISK", valor)
			_driveletter = value
		} else if strings.HasPrefix(strings.ToLower(valor), "name") {
			var value = TieneNombre("FDISK", valor)
			if len(value) > 16 {
				color.Red("[FDISK]: El nombre no puede ser mayor a 16 caracteres")
				break
			} else {
				_name = DevolverNombreByte(value)
			}
		} else if strings.HasPrefix(strings.ToLower(valor), "unit") {
			var value = TieneUnit("FDISK", valor)
			_unit = value
		} else if strings.HasPrefix(strings.ToLower(valor), "type") {
			var value = TieneTypeFDISK(valor)
			_type = value
		} else if strings.HasPrefix(strings.ToLower(valor), "fit") {
			var value = TieneFit("FDISK", valor)
			_fit = value
		} else if strings.HasPrefix(strings.ToLower(valor), "delete") {
			var value = TieneDelete(valor)
			_delete = value
		} else if strings.HasPrefix(strings.ToLower(valor), "add") {
			var value = TieneAdd(valor)
			_add = value
		} else {
			color.Yellow("[FDISK]: Atributo no reconocido")
		}
	}
	return _size, _driveletter, _name, _unit, _type, _fit, _delete, _add
}

func FDISK_Create(_size int32, _driveletter byte, _name []byte, _unit byte, _type byte, _fit byte, _delete string, _add int32) {
	//fmt.Println(_name)
	path := "MIA/P1/Disks/" + string(_driveletter) + ".dsk"
	if !ExisteArchivo("FDISK", path) {
		color.Yellow("[FDISK] Cancel the operation because not yet a file")
		return
	}

	//obtener disco
	tempDisk, existe := ObtainMBRDisk(path)
	if !existe {
		color.Red("[FDISK]: Error al obtener el disco")
		return
	}

	//fmt.Println(_name)
	if _delete == "FULL" {
		DeletePartitionFull(path, _add, _unit, tempDisk, string(_name))
		return
	}
	//fmt.Println("Salto delete")
	if _add != 0 && _unit != '0' {
		fmt.Println("Dentro del add")
		return
	}

	//Verificar si existe partición extendida
	if ExisteExtendida(tempDisk) {
		color.Magenta("[FDISK]: Ya hay una partición extendida")
		return
	}

	//PrintarMBR(tempDisk)
	particion_vacia := PartitionVacia()
	if !VerifyVoidDisk(tempDisk, particion_vacia) {
		color.Magenta("[FDISK]: Todas las particiones estan ocupadas")
		return
	}
	temp_p := PartitionVacia()
	temp_p.Part_status = int8(0)
	temp_p.Part_type = _type
	temp_p.Part_fit = _fit
	temp_p.Part_start = int32(0)
	temp_p.Part_s = Tamano(_size, _unit)
	copy(temp_p.Part_name[:], _name)

	ebr := structures.EBR{}
	ebr.Part_mount = '0'
	ebr.Part_fit = _fit
	ebr.Part_start = int32(0)
	ebr.Part_s = Tamano(_size, _unit)
	ebr.Part_next = int32(-1)
	copy(ebr.Name[:], _name)
	DiscoCreado := false
	correlative := 0

	// Primaria o Extendida
	if _type == 'E' || _type == 'P' {
		//Agregar al primero todos vacios
		if tempDisk.Mbr_partitions[0].Part_status == int8(-1) && tempDisk.Mbr_partitions[1].Part_status == int8(-1) && tempDisk.Mbr_partitions[2].Part_status == int8(-1) && tempDisk.Mbr_partitions[3].Part_status == int8(-1) {
			temp_p.Part_start = int32(size.SizeMBR())
			correlative = 1
			temp_p.Part_correlative = int32(correlative)
			DiscoCreado = true
			tempDisk.Mbr_partitions[0] = temp_p

			//Para el segundo (1 lleno)
		} else if tempDisk.Mbr_partitions[0].Part_status == int8(0) && tempDisk.Mbr_partitions[1].Part_status == int8(-1) && tempDisk.Mbr_partitions[2].Part_status == int8(-1) && tempDisk.Mbr_partitions[3].Part_status == int8(-1) {
			temp_p.Part_start = int32(size.SizeMBR()) + tempDisk.Mbr_partitions[0].Part_s + 1
			correlative = 2
			temp_p.Part_correlative = int32(correlative)
			DiscoCreado = true
			tempDisk.Mbr_partitions[1] = temp_p
			if tempDisk.Mbr_partitions[0].Part_name == [16]byte(_name) {
				color.Magenta("[FDISK]: Nombre de disco ya existente")
				return
			}

			//Para el tercero (1, 2 llenos)
		} else if tempDisk.Mbr_partitions[0].Part_status == int8(0) && tempDisk.Mbr_partitions[1].Part_status == int8(0) && tempDisk.Mbr_partitions[2].Part_status == int8(-1) && tempDisk.Mbr_partitions[3].Part_status == int8(-1) {
			temp_p.Part_start = int32(size.SizeMBR()) + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + 1
			correlative = 3
			temp_p.Part_correlative = int32(correlative)
			DiscoCreado = true
			tempDisk.Mbr_partitions[2] = temp_p
			if tempDisk.Mbr_partitions[0].Part_name == [16]byte(_name) || tempDisk.Mbr_partitions[1].Part_name == [16]byte(_name) {
				color.Magenta("[FDISK]: Nombre de disco ya existente")
				return
			}

			//Para el cuarto (1, 2, 3, llenos)
		} else if tempDisk.Mbr_partitions[0].Part_status == int8(0) && tempDisk.Mbr_partitions[1].Part_status == int8(0) && tempDisk.Mbr_partitions[2].Part_status == int8(0) && tempDisk.Mbr_partitions[3].Part_status == int8(-1) {
			temp_p.Part_start = int32(size.SizeMBR()) + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + tempDisk.Mbr_partitions[2].Part_s + 1
			correlative = 4
			temp_p.Part_correlative = int32(correlative)
			DiscoCreado = true
			tempDisk.Mbr_partitions[3] = temp_p
			if tempDisk.Mbr_partitions[0].Part_name == [16]byte(_name) || tempDisk.Mbr_partitions[1].Part_name == [16]byte(_name) || tempDisk.Mbr_partitions[2].Part_name == [16]byte(_name) {
				color.Magenta("[FDISK]: Nombre de disco ya existente")
				return
			}

			//Para todos llenos (1,2,3,4 llenos)
		} else if tempDisk.Mbr_partitions[0].Part_status == int8(0) && tempDisk.Mbr_partitions[1].Part_status == int8(0) && tempDisk.Mbr_partitions[2].Part_status == int8(0) && tempDisk.Mbr_partitions[3].Part_status == int8(0) {
			temp_p.Part_start = int32(size.SizeMBR()) + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + tempDisk.Mbr_partitions[2].Part_s
			color.Magenta("[FDISK]: Error al crear particiones, Todas llenas")
			return
		}

		//crear disco
		if DiscoCreado {
			if tempDisk.Mbr_tamano < (temp_p.Part_s + temp_p.Part_start) {
				color.Magenta("[FDISK]: Error al crear la partición " + string(_name) + ", Disco insuficiente")
				return
			}

			GuardarParticion(path, tempDisk)
			color.Green("Partición " + string(_name) + " Creada Exitosamente")
			if temp_p.Part_type == 'E' {
				ebr.Part_mount = int8(-1)
				ebr.Part_fit = _fit
				ebr.Part_start = int32(-1)
				ebr.Part_s = int32(-1)
				ebr.Part_next = int32(-1)
				ebr.Name = DevolverNombreByte("-1")
				Escribir_EBR(path, ebr, temp_p.Part_start)
				color.Blue("[FDISK]: EBR grabado Exitosamente")
			}
			//fmt.Println(ebr)
			return
		}
		//Caso en el que ya se pudo haber borrado un disco

	}
}
