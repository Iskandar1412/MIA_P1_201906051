package comandos

import (
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
			fmt.Print(" :: Particion ", (i + 1), " Vacio :: ")
			bandera = true
		} else {
			fmt.Print(" :: Particion ", (i + 1), " Lleno :: ")
			bandera = false
		}
	}
	fmt.Print("\n")
	return bandera
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
