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
		DeletePartitionFull(path, tempDisk, string(_name))
		return
	}
	//fmt.Println("Salto delete")
	if _add != 0 && _unit != '0' {
		AddInPartition(path, _add, _unit, tempDisk, string(_name))
		return
	}

	//Verificar si existe partición extendida
	if ExisteExtendida(tempDisk, _type) {
		color.Magenta("[FDISK]: Ya hay una partición extendida, " + string(_name) + " no puede ser extendida")
		return
	}

	//PrintarMBR(tempDisk)
	particion_vacia := PartitionVacia()
	bandera := false
	bandera_llena := false
	if _type == 'L' {
		bandera_llena = false
		bandera = true
	} else {
		if !VerifyVoidDisk(tempDisk, particion_vacia) {
			bandera = false
			bandera_llena = true
			color.Magenta("[FDISK]: Todas las particiones estan ocupadas")
			return
		}
	}
	fmt.Println("")
	ObtainDisksPrint(tempDisk)
	fmt.Println("")

	if !bandera && bandera_llena {
		return
	}

	temp_p := PartitionVacia()
	temp_p.Part_status = int8(0)
	temp_p.Part_type = _type
	temp_p.Part_fit = _fit
	temp_p.Part_start = int32(0)
	temp_p.Part_s = Tamano(_size, _unit)
	temp_p.Part_name = DevolverNombreByte(string(_name))

	ebr := structures.EBR{}
	ebr.Part_mount = int8(0)
	ebr.Part_fit = _fit
	ebr.Part_start = int32(0)
	ebr.Part_s = Tamano(_size, _unit)
	ebr.Part_next = int32(-1)
	ebr.Name = DevolverNombreByte(string(_name))
	temp_p.Part_name = DevolverNombreByte(string(_name))
	DiscoCreado := false
	correlative := 0

	// Primaria o Extendida
	if _type == 'E' || _type == 'P' {
		fmt.Println("\t\t\t\t\tPartición primaria o extendida", string(_name))
		//Agregar al primero todos vacios
		if tempDisk.Mbr_partitions[0].Part_status == int8(-1) && tempDisk.Mbr_partitions[1].Part_status == int8(-1) && tempDisk.Mbr_partitions[2].Part_status == int8(-1) && tempDisk.Mbr_partitions[3].Part_status == int8(-1) {
			temp_p.Part_start = int32(size.SizeMBR())
			correlative = 1
			temp_p.Part_correlative = int32(correlative)
			DiscoCreado = true
			tempDisk.Mbr_partitions[0] = temp_p

			//Para el segundo (1 lleno)
		} else if tempDisk.Mbr_partitions[0].Part_status == int8(0) && tempDisk.Mbr_partitions[1].Part_status == int8(-1) && tempDisk.Mbr_partitions[2].Part_status == int8(-1) && tempDisk.Mbr_partitions[3].Part_status == int8(-1) {
			temp_p.Part_start = int32(size.SizeMBR()) + tempDisk.Mbr_partitions[0].Part_s
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
			temp_p.Part_start = int32(size.SizeMBR()) + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s
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
			temp_p.Part_start = int32(size.SizeMBR()) + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + tempDisk.Mbr_partitions[2].Part_s
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
				Escribir_EBR("FDISK", path, ebr, temp_p.Part_start)
				color.Blue("[FDISK]: EBR grabado Exitosamente")
			}
			//fmt.Println(ebr)
			return
		}
		//Caso en el que ya se pudo haber borrado un disco
		//color.Yellow("Ya se ha borrado una particion")
		if string(tempDisk.Mbr_partitions[0].Part_name[:]) == string(_name) || string(tempDisk.Mbr_partitions[1].Part_name[:]) == string(_name) || string(tempDisk.Mbr_partitions[2].Part_name[:]) == string(_name) || string(tempDisk.Mbr_partitions[3].Part_name[:]) == string(_name) {
			color.Red("Nombre de partición ya existente")
			return
		}
		if tempDisk.Dsk_fit == 'B' {
			//se leen las 4 particiones
			//se selecciona el menor tamaño a ajustar
			//Best Fit
			menor_size := 0
			numero_e := 0
			if tempDisk.Mbr_partitions[0].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[0].Part_s > temp_p.Part_s {
					if numero_e != 0 {
						if tempDisk.Mbr_partitions[0].Part_s < int32(menor_size) {
							numero_e = 1
							menor_size = int(tempDisk.Mbr_partitions[0].Part_s)
							tempDisk.Mbr_partitions[0].Part_start = size.SizeMBR()
							temp_p.Part_start = size.SizeMBR()
						}
					} else {
						numero_e = 1
						menor_size = int(tempDisk.Mbr_partitions[0].Part_s)
						tempDisk.Mbr_partitions[0].Part_start = size.SizeMBR()
						temp_p.Part_start = size.SizeMBR()
					}
				}
			}
			if tempDisk.Mbr_partitions[1].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[1].Part_s > temp_p.Part_s {
					if numero_e != 0 {
						if tempDisk.Mbr_partitions[1].Part_s < int32(menor_size) {
							numero_e = 1
							menor_size = int(tempDisk.Mbr_partitions[1].Part_s)
							tempDisk.Mbr_partitions[1].Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s
							temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s
						}
					} else {
						numero_e = 1
						menor_size = int(tempDisk.Mbr_partitions[1].Part_s)
						tempDisk.Mbr_partitions[1].Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s
						temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s
					}
				}
			}
			if tempDisk.Mbr_partitions[2].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[2].Part_s > temp_p.Part_s {
					if numero_e != 0 {
						if tempDisk.Mbr_partitions[2].Part_s < int32(menor_size) {
							numero_e = 1
							menor_size = int(tempDisk.Mbr_partitions[2].Part_s)
							tempDisk.Mbr_partitions[2].Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s
							temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s
						}
					} else {
						numero_e = 1
						menor_size = int(tempDisk.Mbr_partitions[2].Part_s)
						tempDisk.Mbr_partitions[2].Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s
						temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s
					}
				}
			}
			if tempDisk.Mbr_partitions[3].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[3].Part_s > temp_p.Part_s {
					if numero_e != 0 {
						if tempDisk.Mbr_partitions[3].Part_s < int32(menor_size) {
							numero_e = 1
							//menor_size = int(tempDisk.Mbr_partitions[3].Part_s)
							tempDisk.Mbr_partitions[3].Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + tempDisk.Mbr_partitions[2].Part_s
							temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + tempDisk.Mbr_partitions[2].Part_s
						}
					} else {
						numero_e = 1
						//menor_size = int(tempDisk.Mbr_partitions[3].Part_s)
						tempDisk.Mbr_partitions[3].Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + tempDisk.Mbr_partitions[2].Part_s
						temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + tempDisk.Mbr_partitions[2].Part_s
					}
				}
			}
			//guardar particion
			if numero_e != 0 {
				GuardarParticionV2(path, temp_p, int32(numero_e))
				if temp_p.Part_type == 'E' {
					ebr.Part_mount = int8(-1)
					ebr.Part_fit = _fit
					ebr.Part_start = int32(-1)
					ebr.Part_s = int32(-1)
					ebr.Part_next = int32(-1)
					ebr.Name = DevolverNombreByte("-1")
					Escribir_EBR("FDISK", path, ebr, temp_p.Part_start)
				}
			}
		} else if tempDisk.Dsk_fit == 'F' {
			//First Fit
			numero_e := 0
			//menor_size := 0
			if tempDisk.Mbr_partitions[0].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[0].Part_s > temp_p.Part_s {
					if numero_e == 0 {
						numero_e = 1
						//menor_size = int(tempDisk.Mbr_partitions[0].Part_s)
						temp_p.Part_start = size.SizeMBR()
					}
				}
			}
			if tempDisk.Mbr_partitions[1].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[1].Part_s > temp_p.Part_s {
					if numero_e == 0 {
						numero_e = 1
						//menor_size = int(tempDisk.Mbr_partitions[1].Part_s)
						temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s
					}
				}
			}
			if tempDisk.Mbr_partitions[2].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[2].Part_s > temp_p.Part_s {
					if numero_e == 0 {
						numero_e = 1
						//menor_size = int(tempDisk.Mbr_partitions[2].Part_s)
						temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s
					}
				}
			}
			if tempDisk.Mbr_partitions[3].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[3].Part_s > temp_p.Part_s {
					if numero_e == 0 {
						numero_e = 1
						//menor_size = int(tempDisk.Mbr_partitions[3].Part_s)
						temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + tempDisk.Mbr_partitions[2].Part_s
					}
				}
			}
			if numero_e != 0 {
				GuardarParticionV2(path, temp_p, int32(numero_e))
				if temp_p.Part_type == 'E' {
					ebr.Part_mount = int8(-1)
					ebr.Part_fit = _fit
					ebr.Part_start = int32(-1)
					ebr.Part_s = int32(-1)
					ebr.Part_next = int32(-1)
					ebr.Name = DevolverNombreByte("-1")
					Escribir_EBR("FDISK", path, ebr, temp_p.Part_start)
				}
			}
		} else if tempDisk.Dsk_fit == 'W' {
			//Worst Fit
			//Leer 4 particiones
			//Seleccionar la de mayor tamaño
			menor_size := 0
			numero_e := 0
			if tempDisk.Mbr_partitions[0].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[0].Part_s > temp_p.Part_s {
					if numero_e != 0 {
						if tempDisk.Mbr_partitions[0].Part_s > int32(menor_size) {
							numero_e = 1
							menor_size = int(tempDisk.Mbr_partitions[0].Part_s)
							temp_p.Part_start = size.SizeMBR()
						}
					} else {
						numero_e = 1
						menor_size = int(tempDisk.Mbr_partitions[0].Part_s)
						temp_p.Part_start = size.SizeMBR()
					}
				}
			}
			if tempDisk.Mbr_partitions[1].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[1].Part_s > temp_p.Part_s {
					if numero_e != 0 {
						if tempDisk.Mbr_partitions[1].Part_s > int32(menor_size) {
							numero_e = 1
							menor_size = int(tempDisk.Mbr_partitions[1].Part_s)
							temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s
						}
					} else {
						numero_e = 1
						menor_size = int(tempDisk.Mbr_partitions[1].Part_s)
						temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s
					}
				}
			}
			if tempDisk.Mbr_partitions[2].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[2].Part_s > temp_p.Part_s {
					if numero_e != 0 {
						if tempDisk.Mbr_partitions[2].Part_s > int32(menor_size) {
							numero_e = 1
							menor_size = int(tempDisk.Mbr_partitions[2].Part_s)
							temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s
						}
					} else {
						numero_e = 1
						menor_size = int(tempDisk.Mbr_partitions[2].Part_s)
						temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s
					}
				}
			}
			if tempDisk.Mbr_partitions[3].Part_status == int8(-1) {
				if tempDisk.Mbr_partitions[3].Part_s > temp_p.Part_s {
					if numero_e != 0 {
						if tempDisk.Mbr_partitions[3].Part_s > int32(menor_size) {
							numero_e = 1
							//menor_size = int(tempDisk.Mbr_partitions[3].Part_s)
							temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + tempDisk.Mbr_partitions[2].Part_s
						}
					} else {
						numero_e = 1
						//menor_size = int(tempDisk.Mbr_partitions[3].Part_s)
						temp_p.Part_start = size.SizeMBR() + tempDisk.Mbr_partitions[0].Part_s + tempDisk.Mbr_partitions[1].Part_s + tempDisk.Mbr_partitions[2].Part_s
					}
				}
			}
			//Guardar particion
			if numero_e != 0 {
				GuardarParticionV2(path, temp_p, int32(numero_e))
			}
		}
	} else {
		fmt.Println("\t\t\t\t\tPartición Lógica", string(_name))
		//particion logica que se guarda
		if string(tempDisk.Mbr_partitions[0].Part_name[:]) == string(_name) || string(tempDisk.Mbr_partitions[1].Part_name[:]) == string(_name) || string(tempDisk.Mbr_partitions[2].Part_name[:]) == string(_name) || string(tempDisk.Mbr_partitions[3].Part_name[:]) == string(_name) {
			color.Red("Particion existente")
			return
		}
		//println("Seguimos")
		particion_ext := PartitionVacia()
		if tempDisk.Mbr_partitions[0].Part_type == 'E' {
			//fmt.Println("primero")
			particion_ext = tempDisk.Mbr_partitions[0]
		} else if tempDisk.Mbr_partitions[1].Part_type == 'E' {
			particion_ext = tempDisk.Mbr_partitions[1]
			//fmt.Println("segundo")
		} else if tempDisk.Mbr_partitions[2].Part_type == 'E' {
			particion_ext = tempDisk.Mbr_partitions[2]
			//fmt.Println("tercero")
		} else if tempDisk.Mbr_partitions[3].Part_type == 'E' {
			particion_ext = tempDisk.Mbr_partitions[3]
			//fmt.Println("cuarto")
		} else {
			color.Red("[FDISK]: No se puede crear particion logica ya que no existe particion extendida")
			return
		}
		//fmt.Println("siguiendo")
		ebr_inicio := particion_ext.Part_start
		//fmt.Println("inicio", ebr_inicio)
		ebr_anterior := structures.EBR{}
		ebr_anterior.Part_mount = int8(-1)
		ebr_anterior.Part_fit = _fit
		ebr_anterior.Part_start = int32(-1)
		ebr_anterior.Part_s = int32(-1)
		ebr_anterior.Part_next = int32(-1)
		ebr_anterior.Name = DevolverNombreByte("-1")
		inicio_ebr_anterior := int32(0)
		size_particion_actual := int32(0)
		inicio_particion_guardar := int32(0)
		siguiente := -1
		contador := 0
		for {
			if ebr_inicio == int32(-1) {
				break
			}
			//fmt.Println("Obtener")
			ebr_actual, err := Obtener_EBR(path, ebr_inicio)
			//fmt.Println(ebr_actual)
			if err != nil {
				return
			}
			if ebr_actual.Name == ebr.Name {
				color.Yellow("Nombre de EBR '" + ToString(ebr.Name[:]) + "' igual")
				return
			}
			if particion_ext.Part_fit == 'B' {
				//Best Fit
				//fmt.Println("Best")
				if ebr_actual.Part_mount == int8(-1) {
					if temp_p.Part_s < ebr_actual.Part_s {
						if size_particion_actual == 0 {
							size_particion_actual = ebr_actual.Part_s
							inicio_particion_guardar = ebr_inicio
							siguiente = int(ebr_actual.Part_next)
						} else if size_particion_actual < temp_p.Part_s {
							size_particion_actual = ebr_actual.Part_s
							inicio_particion_guardar = ebr_inicio
							siguiente = int(ebr_actual.Part_next)
						}
					}
				}
			} else if particion_ext.Part_fit == 'F' {
				//First
				//fmt.Println("First")
				if ebr_actual.Part_mount == int8(-1) {
					if temp_p.Part_s < ebr_actual.Part_s {
						size_particion_actual = ebr_actual.Part_s
						inicio_particion_guardar = ebr_inicio
						siguiente = int(ebr_actual.Part_next)
					}
				}
			} else if particion_ext.Part_fit == 'W' {
				//fmt.Println("Worst")
				if ebr_actual.Part_mount == int8(-1) {
					if temp_p.Part_s < ebr_actual.Part_s {
						if size_particion_actual == 0 {
							size_particion_actual = ebr_actual.Part_s
							inicio_particion_guardar = ebr_inicio
							siguiente = int(ebr_actual.Part_next)
						} else if size_particion_actual > temp_p.Part_s {
							size_particion_actual = ebr_actual.Part_s
							inicio_particion_guardar = ebr_inicio
							siguiente = int(ebr_actual.Part_next)
						}
					}
				}
			}
			if size_particion_actual == 0 && inicio_particion_guardar == 0 && ebr_actual.Part_next == -1 {
				if contador == 0 && ebr_actual.Part_mount == int8(-1) {
					//fmt.Println("vamos")
					//modificar ebr
					inicio_particion_guardar = ebr_inicio
					if (ebr.Part_s + inicio_particion_guardar) > (particion_ext.Part_start + particion_ext.Part_s) {
						color.Red("Error: No hay espacio en partición extendida")
						return
					}
				} else {
					inicio_particion_guardar = ebr_inicio + ebr_actual.Part_s
					if (ebr.Part_s + inicio_particion_guardar) > (particion_ext.Part_start + particion_ext.Part_s) {
						color.Red("Error: No hay espacio en partición extendida")
						return
					}
					ebr_anterior = ebr_actual
					inicio_ebr_anterior = ebr_inicio
				}
			}
			ebr_inicio = ebr_actual.Part_next
			contador++
		}
		total_size := size.SizeEBR()
		ebr.Part_start = inicio_particion_guardar + total_size
		ebr.Part_next = int32(siguiente)
		GuardarEBR(path, ebr, inicio_particion_guardar)
		color.Green("EBR '" + ToString(ebr.Name[:]) + "' grabado")
		///modificar siguiente anterior
		if inicio_ebr_anterior != 0 {
			ebr_anterior.Part_next = inicio_particion_guardar
			GuardarEBR(path, ebr_anterior, inicio_ebr_anterior)
			color.Green("EBR anterior '" + ToString(ebr_anterior.Name[:]) + "' modificado")
		}
	}
}
