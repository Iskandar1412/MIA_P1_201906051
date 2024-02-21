package comandos

import (
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func ljust(s string, leng int) string {
	if len(s) >= leng {
		return s
	}
	return s + strings.Repeat(" ", leng-len(s))
}

func MBR_Report(tempDisk structures.MBR) {
	fmt.Printf("[MBR] => Tamaño: %d :: FechaCreacion: %s :: DiskSignature: %d :: DiskFit: %s", tempDisk.Mbr_tamano, IntFechaToStr(tempDisk.Mbr_fecha_creacion), tempDisk.Mbr_disk_signature, string(tempDisk.Dsk_fit))
	for i := 0; i < len(tempDisk.Mbr_partitions); i++ {
		fmt.Printf("\n\tPartition%d: [Status: %d, Type: %s, Fit: %s, Start: %d, Size: %d, Name: %s, Correlative: %d, ID: %s]", i+1, tempDisk.Mbr_partitions[i].Part_status, string(tempDisk.Mbr_partitions[i].Part_type), string(tempDisk.Mbr_partitions[i].Part_fit), tempDisk.Mbr_partitions[i].Part_start, tempDisk.Mbr_partitions[i].Part_s, ljust(ToString(tempDisk.Mbr_partitions[i].Part_name[:]), 16), tempDisk.Mbr_partitions[i].Part_correlative, ljust(ToString(tempDisk.Mbr_partitions[i].Part_id[:]), 4))
	}
	fmt.Println()
}

func ReportCommandProps(name string, instructions []string) {
	path := "MIA/P1/Disks/A.dsk"
	/*
		var _name string
			var _path string
			var _id string
			var _ruta string
	*/
	if strings.ToUpper(name) == "REP" {
		//fmt.Println("GENERANDO REPORTE")
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
		MBR_Report(tempDisk)
		Report_MBR("A")
	} else {
		color.Red("[ReportCommandProps]: Internal Error")
	}
}

func PrintarMBR(disk structures.MBR) {
	MBR_Report(disk)
}
