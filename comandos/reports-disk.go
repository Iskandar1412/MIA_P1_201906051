package comandos

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Report_DISK(name string, path string, ruta string, id_disco string) {

	// ubicacion := ""
	conjunto, path_disco, ec := Obtener_Particion_ID(id_disco)
	if !ec {
		return
	}

	if conjunto == nil {
		return
	}

	// path_disk := "MIA/P1/Disks/" + string(name) + ".dsk"
	TempDisk, existe := ObtainMBRDisk(path_disco)
	if !existe {
		color.Red("[REP]: Error al obtener el disco")
		return
	}

	particiones_logicas := ""

	data := "digraph G {\n\tnode[shape=plaintext fontsize=12];\n\trankdir=LR;\n\t"
	data += "table [label=<\n\t\t"
	data += "<TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"2\">\n"

	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t<TD ALIGN=\"LEFT\" COLSPAN=\"1\" ROWSPAN=\"3\" BGCOLOR=\"white\" ><FONT COLOR=\"black\"><B>MBR</B></FONT></TD>\n"

	porcentaje_acululado := float32(0)
	porcentaje_acululado_extendida := float32(0)
	for _, disco := range TempDisk.Mbr_partitions {
		columnas_extendidas := 0
		if disco.Part_type == 'P' {
			porcentaje := float32(int(disco.Part_s)*100) / float32(TempDisk.Mbr_tamano)
			porcentaje_acululado += porcentaje
			porcentaje_libre := fmt.Sprint(int32(porcentaje)) + "%"
			data += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" ROWSPAN=\"3\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>PRIMARIA</B><BR/> + " + porcentaje_libre + " (Disco)</FONT></TD>\n"

		} else if disco.Part_type == 'E' {
			porcentaje := float32(int(disco.Part_s)*100) / float32(TempDisk.Mbr_tamano)
			porcentaje_acululado += porcentaje
			//logicas
			siguiente := disco.Part_start
			particiones_logicas += "\t\t\t\t<TR>\n"
			porcentaje_extendida := float32(0)
			// porcenaje_e := ""
			for siguiente != -1 {
				ebr_actual, eea := Obtener_EBR(path_disco, siguiente)
				if eea != nil {
					data += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" ROWSPAN=\"2\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>Libre</B><BR/> + " + "0%" + "</FONT></TD>\n"
					columnas_extendidas += 1
					break
				} else {
					porcentaje = float32(ebr_actual.Part_s*100) / float32(TempDisk.Mbr_tamano)
					porcentaje_extendida = float32(ebr_actual.Part_s*100) / float32(TempDisk.Mbr_tamano)
					porcentaje_libre := fmt.Sprint(int32(porcentaje)) + "%"
					porcentaje_acululado_extendida += porcentaje_extendida

					// porcentaje_libre := fmt.Sprint(porcentaje) + "%"
					data += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" ROWSPAN=\"2\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>EBR</B><BR/></FONT></TD>\n"
					data += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" ROWSPAN=\"2\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>Logica</B><BR/> + " + porcentaje_libre + " (Disco)</FONT></TD>\n"
					columnas_extendidas += 2
					siguiente = ebr_actual.Part_next
				}
			}
			if porcentaje_extendida != 100 {
				porcentaje_libre := ((((100 - porcentaje_extendida) / 100) * float32(disco.Part_s)) * 100) / float32(TempDisk.Mbr_tamano)
				data += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" ROWSPAN=\"2\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>Libre</B><BR/> + " + strconv.FormatFloat(float64(porcentaje_libre), 'f', -1, 32) + "% (Disco)</FONT></TD>\n"
				columnas_extendidas += 1
			}
			particiones_logicas += "\t\t\t\t</TR>\n"
			data += "\t\t\t<TD ALIGN=\"CENTER\"  COLSPAN=\"" + strconv.Itoa(columnas_extendidas) + "\" ROWSPAN=\"1\" BGCOLOR=\"WHITE\"><FONT COLOR=\"BLACK\"><B>Extendida</B></FONT></TD>\n"
		} else {
			porcentaje_libre := ""
			if disco.Part_s != -1 {
				porcentaje := float32(disco.Part_s*100) / float32(TempDisk.Mbr_tamano)
				porcentaje_acululado += porcentaje
				porcentaje_libre = fmt.Sprint(int32(porcentaje)) + "%"
			} else {
				if (TempDisk.Mbr_partitions[0].Part_s == -1) && (TempDisk.Mbr_partitions[1].Part_s == -1) && (TempDisk.Mbr_partitions[2].Part_s == -1) && (TempDisk.Mbr_partitions[3].Part_s == -1) && (TempDisk.Mbr_partitions[0] == disco) {
					porcenaje := 25
					porcentaje_acululado += float32(porcenaje)
					porcentaje_libre = fmt.Sprint(int32(porcenaje)) + "%"
				} else if (TempDisk.Mbr_partitions[1].Part_s == -1) && (TempDisk.Mbr_partitions[2].Part_s == -1) && (TempDisk.Mbr_partitions[3].Part_s == -1) && (TempDisk.Mbr_partitions[1] == disco) {
					porcentaje := (100 - porcentaje_acululado) / float32(3)
					porcentaje_acululado += float32(porcentaje)
					porcentaje_libre = fmt.Sprint(int32(porcentaje)) + "%"
				} else if (TempDisk.Mbr_partitions[2].Part_s == -1) && (TempDisk.Mbr_partitions[0] == disco) {
					porcentaje := (100 - porcentaje_acululado) / float32(2)
					porcentaje_acululado += float32(porcentaje)
					porcentaje_libre = fmt.Sprint(int32(porcentaje)) + "%"
				} else {
					porcentaje := (100 - porcentaje_acululado)
					porcentaje_acululado += float32(porcentaje)
					porcentaje_libre = fmt.Sprint(int32(porcentaje)) + "%"
				}
			}
			data += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"1\" ROWSPAN=\"3\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>Libre</B><BR/> + " + porcentaje_libre + "</FONT></TD>\n"
		}
	}

	data += "\t\t\t</TR>\n"

	data += "\t\t</TABLE>\n"
	data += "\t>];\n}\n"

	//conversión
	nombre_sin_extension := strings.Split(name, ".")
	rutaB := path + "/" + nombre_sin_extension[0] + ".dot"
	err := os.WriteFile(rutaB, []byte(data), 0644)
	if err != nil {
		color.Red("[REP]: Error al guardar el reporte del MBR")
		return
	}
	imagepath := path + "/" + name

	cmd := exec.Command("/usr/bin/dot", "-T"+nombre_sin_extension[1], rutaB, "-o", imagepath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error to generate img", err)
		return
	}
	color.Green("Report Generate [Disk]")
}
