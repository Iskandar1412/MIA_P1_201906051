package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func ExisteExtendida(disk structures.MBR, tipe byte) bool {
	for i := range disk.Mbr_partitions {
		if (disk.Mbr_partitions[i].Part_type == 'E') && (tipe == 'E') {
			//fmt.Println("si")
			return true
		}
		//fmt.Println(i, string(tipe))
	}
	return false
}

// -------------Particiones-------------
func GuardarParticion(path string, estructura structures.MBR) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[FDISK]: No se pudo abrir el archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(0, 0); err != nil {
		color.Red("[FDISK]: No se pudo mover el puntero")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &estructura); err != nil {
		color.Red("[FDISK]: No se pudo escribir en el archivo")
		return
	}
	color.Cyan("[FDISK]: :3")
}

func GuardarParticionV2(path string, particion structures.Partition, numero int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[FDISK]: No se pudo abrir el archivo")
		return
	}
	defer file.Close()

	inicio := int32(-1)
	if numero == 1 {
		inicio = size.SizeMBR_NotPartitions() + 1
	} else if numero == 2 {
		inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + 1
	} else if numero == 3 {
		inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition() + 1
	} else if numero == 4 {
		inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition() + size.SizePartition() + 1
	}
	if inicio == -1 {
		color.Red("Error en particion")
		return
	}

	if _, err := file.Seek(int64(inicio), 0); err != nil {
		color.Red("[FDISK]: No se pudo mover el puntero")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &particion); err != nil {
		color.Red("[FDISK]: No se pudo escribir en el archivo")
		return
	}
	color.Cyan("[FDISK]: :3")
}

// Escribir EBR Particion//-
func Escribir_EBR(path string, ebr_data structures.EBR, start int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[FDISK]: No se pudo abrir el archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(start), 0); err != nil {
		color.Red("[FDISK]: No se pudo mover el puntero")
		return
	}
	//fmt.Println(ebr_data)

	if err := binary.Write(file, binary.LittleEndian, &ebr_data); err != nil {
		color.Red("[FDISK]: No se pudo escribir en el archivo el EBR")
		return
	}
	color.Cyan("[FDISK]: :D")
}

func Escribir_Particion(path string, particion structures.Partition, posicion int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[FDISK]: No se pudo abrir el archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(posicion), 0); err != nil {
		color.Red("[FDISK]: No se pudo mover el puntero")
		return
	}
	//fmt.Println(particion)
	if err := binary.Write(file, binary.LittleEndian, &particion); err != nil {
		color.Red("[FDISK]: No se pudo escribir en el archivo la Particion")
		return
	}
	//color.Cyan("[FDISK]: :D")
} //Escribir EBR Particion//-

// Guardar EBR
func GuardarEBR(path string, ebr structures.EBR, position int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[FDISK]: No se pudo abrir el archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(position), 0); err != nil {
		color.Red("[FDISK]: No se pudo mover el puntero")
		return
	}
	//fmt.Println(ebr)
	if err := binary.Write(file, binary.LittleEndian, &ebr); err != nil {
		color.Red("[FDISK]: No se pudo sobreescribir EBR")
		return
	}
	//color.Cyan("[FDISK]: :D")
}

// ------------VACIAR PARTICION------------
func VaciarParticion(path string, inicio int32, size int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[FDISK]: No se pudo abrir el archivo")
		return
	}
	defer file.Close()
	if _, err := file.Seek(int64(inicio), 0); err != nil {
		color.Red("[FDISK]: No se pudo mover el puntero")
		return
	}
	for i := 0; i < int(size); i++ {
		if _, err := file.Write([]byte{'\x00'}); err != nil {
			color.Red("[FDISK]: No se pudo escribir en el archivo")
			return
		}
	}
	//color.Cyan("[FDISK]: :3")
}

// Obtener EBR
func Obtener_EBR(path string, inicio int32) (structures.EBR, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[EBR]: No se pudo abrir el archivo")
		valor := structures.EBR{}
		return valor, fmt.Errorf("Error")
	}
	defer file.Close()
	if _, err := file.Seek(int64(inicio), 0); err != nil {
		color.Red("[EBR]: No se pudo mover el puntero")
		valor := structures.EBR{}
		return valor, fmt.Errorf("Error")
	}
	var tempEBR structures.EBR
	err = binary.Read(file, binary.LittleEndian, &tempEBR)
	if err != nil {
		return structures.EBR{}, fmt.Errorf("Error")
	}
	return tempEBR, nil
}

// Delete Partition
func DeletePartitionFull(path string, disk structures.MBR, name string) {
	//fmt.Println(disk)
	particion := PartitionVacia()
	valores := 0
	var variableu interface{}
	variableu = 0
	var numero int32 = -1
	if string(disk.Mbr_partitions[0].Part_name[:]) == name {
		particion = disk.Mbr_partitions[0]
		variableu = disk.Mbr_partitions[0]
		numero = 1
		valores = 1
	} else if string(disk.Mbr_partitions[1].Part_name[:]) == name {
		particion = disk.Mbr_partitions[1]
		variableu = disk.Mbr_partitions[1]
		numero = 2
		valores = 1
	} else if string(disk.Mbr_partitions[2].Part_name[:]) == name {
		particion = disk.Mbr_partitions[2]
		variableu = disk.Mbr_partitions[2]
		numero = 3
		valores = 1
	} else if string(disk.Mbr_partitions[3].Part_name[:]) == name {
		particion = disk.Mbr_partitions[3]
		variableu = disk.Mbr_partitions[2]
		numero = 4
		valores = 1
	} else {
		// partición lógica
		particion_extendida := PartitionVacia()
		//numero_particion_ext := 0
		if disk.Mbr_partitions[0].Part_type == 'E' {
			//numero_particion_ext = 1
			particion_extendida = disk.Mbr_partitions[0]
		} else if disk.Mbr_partitions[1].Part_type == 'E' {
			//numero_particion_ext = 2
			particion_extendida = disk.Mbr_partitions[1]
		} else if disk.Mbr_partitions[2].Part_type == 'E' {
			//numero_particion_ext = 3
			particion_extendida = disk.Mbr_partitions[2]
		} else if disk.Mbr_partitions[3].Part_type == 'E' {
			//numero_particion_ext = 4
			particion_extendida = disk.Mbr_partitions[3]
		} else {
			color.Yellow("No existen particiones lógicas para borrar")
			return
		}
		//obtención particiones lógicas
		ebr_inicio := particion_extendida.Part_start
		//recorrer EBR's
		for {
			if ebr_inicio == -1 {
				break
			}
			ebr_actual, err := Obtener_EBR(path, ebr_inicio)
			if err != nil {
				color.Red("Error")
				return
			}
			if string(ebr_actual.Name[:]) == name {
				fmt.Println("Encontrado Lógico")
				ebr_actual.Part_mount = -1
				ebr_actual.Name = DevolverNombreByte("-1")
				GuardarEBR(path, ebr_actual, ebr_inicio)
				VaciarParticion(path, ebr_actual.Part_start, ebr_actual.Part_s-size.SizeEBR())
				color.Green("Particion Eliminada")
				return
			}
			ebr_inicio = ebr_actual.Part_next
			continue
		}
	}
	//En el caso de ser primaraia o extendida
	if variableu != 0 || valores != 0 {
		fmt.Println("Primario o Extendido")
		particion.Part_status = -1
		particion.Part_name = DevolverNombreByte("-1")
		if numero == 1 {
			disk.Mbr_partitions[0] = particion
			GuardarParticion(path, disk)
			VaciarParticion(path, particion.Part_start, particion.Part_s)
			color.Green("Particion 0 Eliminada")
		}
		if numero == 2 {
			disk.Mbr_partitions[1] = particion
			GuardarParticion(path, disk)
			VaciarParticion(path, particion.Part_start, particion.Part_s)
			color.Green("Particion 1 Eliminada")
		}
		if numero == 3 {
			disk.Mbr_partitions[2] = particion
			GuardarParticion(path, disk)
			VaciarParticion(path, particion.Part_start, particion.Part_s)
			color.Green("Particion 2 Eliminada")
		}
		if numero == 4 {
			disk.Mbr_partitions[3] = particion
			GuardarParticion(path, disk)
			VaciarParticion(path, particion.Part_start, particion.Part_s)
			color.Green("Particion 3 Eliminada")
		}
	} else {
		color.Red("No se pudo encontrar la partición")
		return
	}
}

// Add
func AddInPartition(path string, _add int32, _unit byte, disk structures.MBR, name string) {
	particion_siguiente := PartitionVacia()
	numero := -1
	valores := 0
	particion := PartitionVacia()
	inicio_siguiente := int32(0)
	if string(disk.Mbr_partitions[0].Part_name[:]) == name {
		particion = disk.Mbr_partitions[0]
		particion_siguiente = disk.Mbr_partitions[1]
		numero = 1
		valores = 1
		inicio_siguiente = particion.Part_start
	} else if string(disk.Mbr_partitions[1].Part_name[:]) == name {
		particion = disk.Mbr_partitions[1]
		particion_siguiente = disk.Mbr_partitions[2]
		numero = 2
		valores = 1
		inicio_siguiente = particion.Part_start
	} else if string(disk.Mbr_partitions[2].Part_name[:]) == name {
		particion = disk.Mbr_partitions[2]
		particion_siguiente = disk.Mbr_partitions[3]
		numero = 3
		valores = 1
		inicio_siguiente = particion.Part_start
	} else if string(disk.Mbr_partitions[3].Part_name[:]) == name {
		particion = disk.Mbr_partitions[3]
		particion_siguiente = PartitionVacia()
		numero = 4
		valores = 1
		inicio_siguiente = disk.Mbr_tamano
	} else {
		fmt.Println("No es Particion primaria o extendida")
		//numero_particion_ext := 0
		particion_extendida := PartitionVacia()
		if disk.Mbr_partitions[0].Part_type == 'E' {
			particion_extendida = disk.Mbr_partitions[0]
		} else if disk.Mbr_partitions[1].Part_type == 'E' {
			particion_extendida = disk.Mbr_partitions[1]
		} else if disk.Mbr_partitions[2].Part_type == 'E' {
			particion_extendida = disk.Mbr_partitions[2]
		} else if disk.Mbr_partitions[3].Part_type == 'E' {
			particion_extendida = disk.Mbr_partitions[3]
		} else {
			fmt.Println("No ha particiones logicas")
			return
		}

		//obtención de partiiones logicas
		ebr_inicio := particion_extendida.Part_start
		contador := 0
		for {
			if ebr_inicio == -1 {
				color.Red("Fin Lectura EBR")
				break
			}
			ebr_actual, err := Obtener_EBR(path, ebr_inicio)
			if err != nil {
				color.Red("Error ADD")
				return
			}
			if string(ebr_actual.Name[:]) == name {
				fmt.Println("Logico encontrado")
				if ebr_actual.Part_next == -1 {
					//no hay siguiente partición logica, mayor a 0, menor a tamaño particion ex
					if (ebr_actual.Part_s+Tamano(_add, _unit) > 0) && ((ebr_actual.Part_start + ebr_actual.Part_s + Tamano(_add, _unit)) < (particion_extendida.Part_start + particion_extendida.Part_s)) {
						ebr_actual.Part_s = ebr_actual.Part_s + Tamano(_add, _unit)
						GuardarEBR(path, ebr_actual, ebr_inicio)
						color.Green("Espacio Logico Modificado")
						return
					}
				} else {
					agrego := false
					//particion logica siguiente
					if ebr_actual.Part_next != -1 {
						ebr_siguiente, err := Obtener_EBR(path, ebr_actual.Part_next)
						if err != nil {
							color.Red("Error ADD")
							return
						}
						if ebr_siguiente.Part_mount == -1 {
							//tienen siguiente particion logica
							//mayor a 0, menor al tamaño de la particion extendida
							//menora l final de la siguiente particion logica
							if (ebr_actual.Part_s+Tamano(_add, _unit) > 0) && ((ebr_actual.Part_s + ebr_actual.Part_start + Tamano(_add, _unit)) < (particion_extendida.Part_start + particion_extendida.Part_s)) && ((ebr_actual.Part_s + ebr_actual.Part_start + Tamano(_add, _unit)) < (ebr_siguiente.Part_start + ebr_siguiente.Part_s)) {
								//recorrer ebr de la siguiente particion
								//inicio = inicio siguiente + add
								GuardarEBR(path, ebr_actual, (ebr_actual.Part_next + Tamano(_add, _unit)))
								ebr_actual.Part_s = ebr_actual.Part_s + Tamano(_add, _unit)
								GuardarEBR(path, ebr_actual, ebr_inicio)
								color.Green("Espacio Logico Modificado")
							}
						}
					} else {
						//no hay siguiente particion logica
						//mayor a 0, menor al tamaño de la particion extendida
						if (ebr_actual.Part_s+Tamano(_add, _unit) > 0) && ((ebr_actual.Part_s + ebr_actual.Part_start + Tamano(_add, _unit)) < (particion_extendida.Part_start + particion_extendida.Part_s)) {
							ebr_actual.Part_s = ebr_actual.Part_s + Tamano(_add, _unit)
							GuardarEBR(path, ebr_actual, ebr_inicio)
							color.Green("Espacio Logico Modificado")
						}
					}
					if !agrego {
						//no importa si hay particion siguiente
						//no pasa del tamaño inicial de la particion logica (solo quita y añade otra vez)
						//mayor a 0
						if (ebr_actual.Part_s+Tamano(_add, _unit) > 0) && (ebr_actual.Part_start+ebr_actual.Part_s+Tamano(_add, _unit) < ebr_actual.Part_next) {
							ebr_actual.Part_s = ebr_actual.Part_s + Tamano(_add, _unit)
							GuardarEBR(path, ebr_actual, ebr_inicio)
							color.Green("Espacio Logico Modificado")
							return
						}
					}
					return
				}
			}
			ebr_inicio = ebr_actual.Part_next
			contador += 1
		}
		return
	}
	//aqui se borra la particion primaria o extendida
	agrego := false
	if valores != 0 {
		//si se creo un aparticion luego de esta
		//esta inactivo, tamaño actual mayor a 0
		if particion_siguiente.Part_status == -1 && particion_siguiente.Part_s > 0 { //-1 no usado, 0 desmontado, 1 montado
			//mayor a 0, menor al final de la siguietne particion
			if particion.Part_s+Tamano(_add, _unit) > 0 && ((particion.Part_s + Tamano(_add, _unit)) < (particion_siguiente.Part_start + particion_siguiente.Part_s)) {
				particion_siguiente.Part_start = particion_siguiente.Part_start + Tamano(_add, _unit)
				GuardarParticionV2(path, particion, int32(numero+1))
				particion.Part_s = particion.Part_s + Tamano(_add, _unit)
				GuardarParticionV2(path, particion, int32(numero))
				agrego = true
				color.Green("Espacio Modificado")
			}
			//no se creo particion des puesd de esta
		} else if particion_siguiente.Part_status == -1 && particion_siguiente.Part_s == 0 {
			//no se han creado mas particiones
			//mayora  0, menor a tamaño de disco
			if particion.Part_s+Tamano(_add, _unit) > 0 && ((particion.Part_s + Tamano(_add, _unit) + particion.Part_start) < disk.Mbr_tamano) {
				particion.Part_s = particion.Part_s + Tamano(_add, _unit)
				GuardarParticionV2(path, particion, int32(numero))
				agrego = true
				color.Green("Espacio Modificado")
			}
		}
	} else {
		//si particion siguiente es igual a 0 solo es la particion 4 la que se modificara
		//mayor a 0, menor a tamaño de particion extendida
		if particion.Part_s+Tamano(_add, _unit) > 0 && (particion.Part_s+Tamano(_add, _unit)+particion.Part_start) < disk.Mbr_tamano {
			particion.Part_s = particion.Part_s + Tamano(_add, _unit)
			GuardarParticionV2(path, particion, int32(numero))
			agrego = true
			color.Green("Espacio Modificado")
		}
	}
	if !agrego {
		if (particion.Part_s+Tamano(_add, _unit) > 0) && (particion.Part_s+particion.Part_start+Tamano(_add, _unit) < inicio_siguiente) {
			particion.Part_s = particion.Part_s + Tamano(_add, _unit)
			GuardarParticionV2(path, particion, int32(numero))
			color.Green("Espacio Modificado")
			return
		} else {
			color.Red("No se pudo modificar el espacio")
			return
		}
	}
}
