package comandos

import (
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func ReportCommandProps(name string, instructions []string) {
	path := "MIA/P1/Disks/"
	/*
		var _name string
		var _path string
		var _id string
		var _ruta string
	*/
	if strings.ToUpper(name) == "REP" {
		fmt.Println("GENERANDO REPORTE")
		file, err := os.Open(path)
		if err != nil {
			color.Red("Error al leer Archivo")
			return
		}
		defer file.Close()
		var tempDisk structures.MBR
		err = binary.Read(file, binary.LittleEndian, &tempDisk)
		if err != nil {
			return
		}
		fmt.Printf("[MBR]\n\tTamaño: %d\n\tFechaCreacion: %s\n\tDiskSignature: %d\n\tDiskFit: %s", tempDisk.Mbr_tamano, IntFechaToStr(tempDisk.Mbr_fecha_creacion), tempDisk.Mbr_disk_signature, string(tempDisk.Dsk_fit))
		for i := 0; i < len(tempDisk.Mbr_partitions); i++ {
			fmt.Printf("\n\tPartition%d: [Status: %s, Type: %s, Fit: %s, Start: %d, Size: %d, Name: %s, Correlative: %d, ID: %s]", i+1, string(tempDisk.Mbr_partitions[i].Part_status), string(tempDisk.Mbr_partitions[i].Part_type), string(tempDisk.Mbr_partitions[i].Part_fit), tempDisk.Mbr_partitions[i].Part_start, tempDisk.Mbr_partitions[i].Part_s, ToString(tempDisk.Mbr_partitions[i].Part_name[:]), tempDisk.Mbr_partitions[i].Part_correlative, ToString(tempDisk.Mbr_partitions[i].Part_id[:]))
		}
		fmt.Println()
	} else {
		color.Red("[ReportCommandProps]: Internal Error")
	}
}
