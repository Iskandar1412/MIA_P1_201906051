package comandos

import (
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func ExisteExtendida(disk structures.MBR) bool {
	if disk.Mbr_partitions[0].Part_type == 'E' || disk.Mbr_partitions[1].Part_type == 'E' || disk.Mbr_partitions[2].Part_type == 'E' || disk.Mbr_partitions[3].Part_type == 'E' {
		color.Yellow("[FDISK]: Particion Extendida existente")
		return true
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
	color.Cyan("[FDISK]: :D")
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
	fmt.Println(ebr)
	if err := binary.Write(file, binary.LittleEndian, &ebr); err != nil {
		color.Red("[FDISK]: No se pudo sobreescribir EBR")
		return
	}
	color.Cyan("[FDISK]: :D")
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
	color.Cyan("[FDISK]: :3")
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
func DeletePartitionFull(path string, _add int32, _unit byte, disk structures.MBR, name string) {
	//fmt.Println(disk)
	particion := PartitionVacia()
	var numero int32 = -1
	if string(disk.Mbr_partitions[0].Part_name[:]) == name {
		particion = disk.Mbr_partitions[0]
		numero = 1
	} else if string(disk.Mbr_partitions[1].Part_name[:]) == name {
		particion = disk.Mbr_partitions[1]
		numero = 2
	} else if string(disk.Mbr_partitions[2].Part_name[:]) == name {
		particion = disk.Mbr_partitions[2]
		numero = 3
	} else if string(disk.Mbr_partitions[3].Part_name[:]) == name {
		particion = disk.Mbr_partitions[3]
		numero = 4
	} else {
		// partición lógica
		particion_extendida := PartitionVacia()
		numero_particion_ext := 0
		if disk.Mbr_partitions[0].Part_type == 'E' {
			numero_particion_ext = 1
			particion_extendida = disk.Mbr_partitions[0]
		} else if disk.Mbr_partitions[1].Part_type == 'E' {
			numero_particion_ext = 2
			particion_extendida = disk.Mbr_partitions[1]
		} else if disk.Mbr_partitions[2].Part_type == 'E' {
			numero_particion_ext = 3
			particion_extendida = disk.Mbr_partitions[2]
		} else if disk.Mbr_partitions[3].Part_type == 'E' {
			numero_particion_ext = 4
			particion_extendida = disk.Mbr_partitions[3]
		} else {
			color.Yellow("No existen particiones lógicas para borrar")
			return
		}
		ebr_inicio := particion_extendida.Part_start
		for true {
			if ebr_inicio == -1 {
				break
			}
			ebr_actual, er := Obtener_EBR(path, ebr_inicio)
			if er != nil {
				color.Red("[EBR]: No se pudo obtener el EBR")
				return
			}
			if string(ebr_actual.Name[:]) == name {
				if ebr_actual.Part_next == -1 {
					//como no hay particion logica siguiente
					// mayor a 0 ->  menor a tamaño partición extendida
					if (ebr_actual.Part_s+Tamano(_add, _unit) > 0) && ((ebr_actual.Part_start + ebr_actual.Part_s + Tamano(_add, _unit)) < (particion_extendida.Part_start + particion_extendida.Part_s)) {
						ebr_actual.Part_s = ebr_actual.Part_s + Tamano(_add, _unit)
						GuardarEBR(path, ebr_actual, ebr_inicio)
						color.Green("[FDISK] Espacio lógico modificado")
						return
					}
				} else {
					agregar := false
					//hay particion logica siguiente
				}
			}
		}
	}
}
