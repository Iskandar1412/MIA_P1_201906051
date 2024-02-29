package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Obtain MBR Disk Obtener disco
func ObtainMBRDisk(path string) (structures.MBR, bool) {
	file, err := os.Open(path)
	if err != nil {
		color.Red("Error al leer Archivo")
		return structures.MBR{}, false
	}
	defer file.Close()
	var tempDisk structures.MBR
	err = binary.Read(file, binary.LittleEndian, &tempDisk)
	if err != nil {
		return structures.MBR{}, false
	}
	return tempDisk, true
}

func VerifyVoidDisk(disk structures.MBR, void_partition structures.Partition) bool {
	bandera := false
	for i := range disk.Mbr_partitions {
		if disk.Mbr_partitions[i] == void_partition {
			// fmt.Print(" :: Particion ", (i + 1), " Vacio :: ")
			bandera = true
		} else {
			// fmt.Print(" :: Particion ", (i + 1), " Lleno :: ")
			bandera = false
		}
	}
	// fmt.Print("\n")
	return bandera
}

func ObtainDisksPrint(disk structures.MBR) {
	particion := PartitionVacia()
	for i := range disk.Mbr_partitions {
		if disk.Mbr_partitions[i] == particion {
			fmt.Print(" :: Particion ", (i + 1), " Vacio :: ")
			// bandera = true
		} else {
			fmt.Print(" :: Particion ", (i + 1), " Lleno :: ")
			// bandera = false
		}
	}
	fmt.Print("\n")
}

func PartitionVacia() structures.Partition {
	var partition structures.Partition
	partition.Part_status = int8(-1)
	partition.Part_type = 'P'
	partition.Part_fit = 'F'
	partition.Part_start = -1
	partition.Part_s = -1
	for i := 0; i < len(partition.Part_name); i++ {
		partition.Part_name[i] = '\x00'
	}
	partition.Part_correlative = -1
	for i := 0; i < len(partition.Part_id); i++ {
		partition.Part_id[i] = '\x00'
	}
	return partition
}

// Disco - Donde inicia la particion
func BuscarParticion(disco structures.MBR, nombre []byte, path string) ([]interface{}, bool) {
	Devolucion := make([]interface{}, 4) //0 (particion ex, prim) 1 (inicio ex, prim) 2 (particion log) 3 (inicio log)
	tempDisk := PartitionVacia()
	inicio := int32(0)
	//inicio_extendida := int32(0)
	error := false
	es_primaria_extendida := false
	// particion := PartitionVacia()
	for i := range disco.Mbr_partitions {
		if (string(disco.Mbr_partitions[i].Part_name[:]) == string(nombre)) && ToString(nombre) != "" {
			tempDisk = disco.Mbr_partitions[i]
			Devolucion[0] = tempDisk
			es_primaria_extendida = true
			if i == 0 {
				inicio = size.SizeMBR_NotPartitions()
				Devolucion[1] = inicio
			} else if i == 1 {
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition()
				Devolucion[1] = inicio
			} else if i == 2 {
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition()
				Devolucion[1] = inicio
			} else if i == 3 {
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition() + size.SizePartition()
				Devolucion[1] = inicio
			} else {
				inicio = 0
				error = true
			}
			break
		} else {
			if ToString(nombre) == "" {
				color.Yellow("[MOUNT]: No hay nombre en la partición que se ha ingresado")
				return Devolucion, true
			}
		}
	}
	if es_primaria_extendida {
		return Devolucion, error
	}
	//puede ser partición logica
	for j := range disco.Mbr_partitions {
		if disco.Mbr_partitions[j].Part_type == 'E' {
			tempDisk = disco.Mbr_partitions[j]
			if j == 0 {
				Devolucion[0] = tempDisk //2 para part log y 3 inicio log
				inicio = size.SizeMBR_NotPartitions()
				Devolucion[1] = inicio
			} else if j == 1 {
				Devolucion[0] = tempDisk //2 para part log y 3 inicio log
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition()
				Devolucion[1] = inicio
			} else if j == 2 {
				Devolucion[0] = tempDisk //2 para part log y 3 inicio log
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition()
				Devolucion[1] = inicio
			} else if j == 3 {
				Devolucion[0] = tempDisk //2 para part log y 3 inicio log
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition() + size.SizePartition()
				Devolucion[1] = inicio
			} else {
				inicio = 0
				error = true
			}
			break
		}
	}
	//existe partición lógica con ello
	siguiente := tempDisk.Part_start
	//ebr_anterior := structures.EBR{Part_mount: int8(-1), Part_fit: 'W', Part_start: int32(-1), Part_s: int32(-1), Part_next: int32(-1), Name: DevolverNombreByte("-1")}
	ebr_anterior_data := structures.EBR{Part_mount: int8(-1), Part_fit: 'W', Part_start: int32(-1), Part_s: int32(-1), Part_next: int32(-1), Name: DevolverNombreByte("-1")}
	for {
		if siguiente == -1 {
			color.Magenta("[MOUNT]: No hay particion lógica con ese nombre")
			Devolucion[0] = nil
			Devolucion[1] = nil
			return Devolucion, true
		}
		ebr_actual, _err := Obtener_EBR(path, siguiente)
		if _err != nil {
			color.Red("[MOUNT]: No se ha encontrado en particiones lógicas")
			Devolucion[0] = nil
			Devolucion[1] = nil
			return Devolucion, true
		}
		if string(ebr_actual.Name[:]) == string(nombre) {
			//existe la particion
			if ebr_anterior_data.Part_start == int32(-1) && ebr_anterior_data.Part_mount == int8(-1) && ebr_anterior_data.Part_s == int32(-1) {
				//caso que sea el primero
				Devolucion[2] = ebr_actual
				Devolucion[3] = tempDisk.Part_start
			} else {
				Devolucion[2] = ebr_actual
				Devolucion[3] = ebr_anterior_data.Part_next
			}
			return Devolucion, error
		}
		//avanzamos al siguiente ebr
		ebr_anterior_data = ebr_actual
		siguiente = ebr_actual.Part_next
	}
}

func BuscarParticionV2(nombre []byte, path string) ([]interface{}, bool) {
	Devolucion := make([]interface{}, 4) //0 (particion ex, prim) 1 (inicio ex, prim) 2 (particion log) 3 (inicio log)
	disco, existe := ObtainMBRDisk(path)
	if !existe {
		color.Red("Error en la obtención del disco")
		return Devolucion, false
	}
	tempDisk := PartitionVacia()
	inicio := int32(0)
	//inicio_extendida := int32(0)
	error := false
	es_primaria_extendida := false
	// particion := PartitionVacia()
	for i := range disco.Mbr_partitions {
		if (string(disco.Mbr_partitions[i].Part_name[:]) == string(nombre)) && ToString(nombre) != "" {
			tempDisk = disco.Mbr_partitions[i]
			Devolucion[0] = tempDisk
			es_primaria_extendida = true
			if i == 0 {
				inicio = size.SizeMBR_NotPartitions()
				Devolucion[1] = inicio
			} else if i == 1 {
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition()
				Devolucion[1] = inicio
			} else if i == 2 {
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition()
				Devolucion[1] = inicio
			} else if i == 3 {
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition() + size.SizePartition()
				Devolucion[1] = inicio
			} else {
				inicio = 0
				error = true
			}
			break
		} else {
			if ToString(nombre) == "" {
				color.Yellow("[MOUNT]: No hay nombre en la partición que se ha ingresado")
				return Devolucion, true
			}
		}
	}
	if es_primaria_extendida {
		return Devolucion, error
	}
	//puede ser partición logica
	for j := range disco.Mbr_partitions {
		if disco.Mbr_partitions[j].Part_type == 'E' {
			tempDisk = disco.Mbr_partitions[j]
			if j == 0 {
				Devolucion[0] = tempDisk //2 para part log y 3 inicio log
				inicio = size.SizeMBR_NotPartitions()
				Devolucion[1] = inicio
			} else if j == 1 {
				Devolucion[0] = tempDisk //2 para part log y 3 inicio log
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition()
				Devolucion[1] = inicio
			} else if j == 2 {
				Devolucion[0] = tempDisk //2 para part log y 3 inicio log
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition()
				Devolucion[1] = inicio
			} else if j == 3 {
				Devolucion[0] = tempDisk //2 para part log y 3 inicio log
				inicio = size.SizeMBR_NotPartitions() + size.SizePartition() + size.SizePartition() + size.SizePartition()
				Devolucion[1] = inicio
			} else {
				inicio = 0
				error = true
			}
			break
		}
	}
	//existe partición lógica con ello
	siguiente := tempDisk.Part_start
	//ebr_anterior := structures.EBR{Part_mount: int8(-1), Part_fit: 'W', Part_start: int32(-1), Part_s: int32(-1), Part_next: int32(-1), Name: DevolverNombreByte("-1")}
	ebr_anterior_data := structures.EBR{Part_mount: int8(-1), Part_fit: 'W', Part_start: int32(-1), Part_s: int32(-1), Part_next: int32(-1), Name: DevolverNombreByte("-1")}
	for {
		if siguiente == -1 {
			color.Magenta("[MOUNT]: No hay particion lógica con ese nombre")
			Devolucion[0] = nil
			Devolucion[1] = nil
			return Devolucion, true
		}
		ebr_actual, _err := Obtener_EBR(path, siguiente)
		if _err != nil {
			color.Red("[MOUNT]: No se ha encontrado en particiones lógicas")
			Devolucion[0] = nil
			Devolucion[1] = nil
			return Devolucion, true
		}
		if string(ebr_actual.Name[:]) == string(nombre) {
			//existe la particion
			if ebr_anterior_data.Part_start == int32(-1) && ebr_anterior_data.Part_mount == int8(-1) && ebr_anterior_data.Part_s == int32(-1) {
				//caso que sea el primero
				Devolucion[2] = ebr_actual
				Devolucion[3] = tempDisk.Part_start
			} else {
				Devolucion[2] = ebr_actual
				Devolucion[3] = ebr_anterior_data.Part_next
			}
			return Devolucion, error
		}
		//avanzamos al siguiente ebr
		ebr_anterior_data = ebr_actual
		siguiente = ebr_actual.Part_next
	}
}

func EliminarParticionMount(comando string, particion string) {
	if particion == "" {
		color.Red("[" + comando + "]: No se ha ingresado el nombre de la partición")
		return
	}
	var Partis []interface{}
	for _, discos := range Partitions_Mounted {
		if disco, ok := discos.([]string); ok {
			if disco[0] != particion {
				Partis = append(Partis, disco)
			}
		}
	}
	Partitions_Mounted = Partis
	//fmt.Println(Partis)
	color.Cyan("[" + comando + "]: Desmontada exitosamente partición... '" + particion + "'")
}

func Obtener_Particion_ID(id string) ([]interface{}, string, bool) {
	var conecta []interface{}
	for _, particion := range Partitions_Mounted {
		if particion, ok := particion.([]string); ok {
			if particion[0] == id {
				nombre := DevolverNombreByte(particion[1])
				conjunto, error := BuscarParticionV2(nombre[:], particion[3])
				if error {
					color.Red("Particion no encontrada")
					return conecta, "", false
				}
				return conjunto, particion[3], true
			}
		}
	}
	return conecta, "", false
}
