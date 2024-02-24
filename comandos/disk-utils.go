package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"fmt"
	"os"
	"reflect"

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

func Obtener_Superbloque(comando string, path string, nombre string) (structures.SuperBlock, bool) {
	superbloque := structures.SuperBlock{}
	byte_escribir := int32(0)
	tempDisk, existe := ObtainMBRDisk(path)
	if !existe {
		color.Red("[" + comando + "]: Error en la obtención del disco")
		return superbloque, false
	}
	value := ToString([]byte(nombre))
	value2 := DevolverNombreByte(value)
	conjunto, _err := BuscarParticion(tempDisk, value2[:], path)
	if _err {
		color.Red("[" + comando + "]: No se ha encontrado en particiones lógicas")
		return superbloque, false
	}

	ebr := structures.EBR{}
	if temp, ok := conjunto[2].(structures.EBR); ok {
		v := reflect.ValueOf(temp)
		reflect.ValueOf(&ebr).Elem().Set(v)
		byte_escribir = ebr.Part_start
		conjunto[0] = nil
		conjunto[1] = nil
	}

	particion := structures.Partition{}
	if temp, ok := conjunto[0].(structures.Partition); ok {
		v := reflect.ValueOf(temp)
		reflect.ValueOf(&particion).Elem().Set(v)
		byte_escribir = particion.Part_start
	}

	if byte_escribir == int32(0) {
		return superbloque, false
	}

	file, err := os.Open(path)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return superbloque, false
	}
	defer file.Close()
	if _, err := file.Seek(int64(byte_escribir), 0); err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return superbloque, false
	}

	err = binary.Read(file, binary.LittleEndian, &superbloque)
	if err != nil {
		color.Red("[" + comando + "]: Error en lectura de superbloque")
		return superbloque, false
	}
	if superbloque.S_mtime == 0 && superbloque.S_umtime == 0 {
		color.Yellow("[" + comando + "]: Super Bloque no creado, Unidad no formateada")
		return superbloque, false
	}
	return superbloque, true
}

func Guardar_Superbloque(comando string, path string, inicio int32, superblock structures.SuperBlock) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return
	}
	defer file.Close()

	if _, err := file.Seek(int64(inicio), 0); err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &superblock); err != nil {
		color.Red("[" + comando + "]: Error en escritura de superbloque")
		return
	}
}
