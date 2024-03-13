package comandos

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Report_MBR(name string, path string, ruta string, id_disco string) {

	// ubicacion := ""
	path_disco, edisco := Obtener_Disco_ID(id_disco)
	if !edisco {
		color.Red("[REP]: Disco no encontrado")
		return
	}

	// path_disk := "MIA/P1/Disks/" + string(name) + ".dsk"
	TempDisk, existe := ObtainMBRDisk(path_disco)
	if !existe {
		color.Red("[REP]: Error al obtener el disco")
		return
	}

	data := "digraph G {\n\tnode[shape=plaintext fontsize=12];\n\trankdir=LR;\n\t"
	data += "table [label=<\n\t\t"
	data += "<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"

	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t<TD ALIGN=\"LEFT\" COLSPAN=\"2\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>REPORTE MBR</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"

	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>mbr_tamano</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(TempDisk.Mbr_tamano)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"

	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>mbr_fecha_creacion</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>" + IntFechaToStr(TempDisk.Mbr_fecha_creacion) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"

	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>mbr_disk_signature</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(TempDisk.Mbr_disk_signature)) + "</B></FONT></TD>\n"
	data += "\t\t\t</TR>\n"

	for i := range TempDisk.Mbr_partitions {
		data += "\t\t\t<TR>\n"
		data += "\t\t\t\t<TD ALIGN=\"LEFT\" COLSPAN=\"2\" BGCOLOR=\"#392250\" ><FONT COLOR=\"WHITE\"><B>Partition</B></FONT></TD>\n"
		data += "\t\t\t</TR>\n"

		data += "\t\t\t<TR>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>part_status</B></FONT></TD>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(TempDisk.Mbr_partitions[i].Part_status)) + "</B></FONT></TD>\n"
		data += "\t\t\t</TR>\n"

		data += "\t\t\t<TR>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>part_type</B></FONT></TD>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>" + string(TempDisk.Mbr_partitions[i].Part_type) + "</B></FONT></TD>\n"
		data += "\t\t\t</TR>\n"

		data += "\t\t\t<TR>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>part_fit</B></FONT></TD>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>" + string(TempDisk.Mbr_partitions[i].Part_fit) + "</B></FONT></TD>\n"
		data += "\t\t\t</TR>\n"

		data += "\t\t\t<TR>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>part_start</B></FONT></TD>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(TempDisk.Mbr_partitions[i].Part_start)) + "</B></FONT></TD>\n"
		data += "\t\t\t</TR>\n"

		data += "\t\t\t<TR>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>part_size</B></FONT></TD>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>" + strconv.Itoa(int(TempDisk.Mbr_partitions[i].Part_s)) + "</B></FONT></TD>\n"
		data += "\t\t\t</TR>\n"

		data += "\t\t\t<TR>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>part_name</B></FONT></TD>\n"
		data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"WHITE\" ><FONT COLOR=\"BLACK\"><B>" + returnstring(ToString(TempDisk.Mbr_partitions[i].Part_name[:])) + "</B></FONT></TD>\n"
		data += "\t\t\t</TR>\n"
		if TempDisk.Mbr_partitions[i].Part_type == 'E' {
			//particion extendida
			siguiente := TempDisk.Mbr_partitions[i].Part_start
			for siguiente != -1 {
				ebr_actual, ebr := Obtener_EBR(path_disco, siguiente)
				if ebr != nil {
					return
				}
				if ebr_actual.Part_mount != -1 {
					data += "\t\t\t\t<TR>\n"
					data += "\t\t\t\t\t<TD ALIGN=\"LEFT\" COLSPAN=\"2\" BGCOLOR=\"#F58C8F\"><FONT COLOR=\"white\"><B>Particion Logica</B></FONT></TD>"
					data += "\t\t\t\t</TR>\n"
					data += "\t\t\t\t<TR>\n"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>Part_Status</B></FONT></TD>"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>" + strconv.Itoa(int(ebr_actual.Part_mount)) + "</B></FONT></TD>"
					data += "\t\t\t\t</TR>\n"
					data += "\t\t\t\t<TR>\n"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>Part_Next</B></FONT></TD>"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>" + strconv.Itoa(int(ebr_actual.Part_next)) + "</B></FONT></TD>"
					data += "\t\t\t\t</TR>\n"
					data += "\t\t\t\t<TR>\n"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>Part_Fit</B></FONT></TD>"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>" + string(ebr_actual.Part_fit) + "</B></FONT></TD>"
					data += "\t\t\t\t</TR>\n"
					data += "\t\t\t\t<TR>\n"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>Part_Start</B></FONT></TD>"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>" + strconv.Itoa(int(ebr_actual.Part_start)) + "</B></FONT></TD>"
					data += "\t\t\t\t</TR>\n"
					data += "\t\t\t\t<TR>\n"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>Part_S</B></FONT></TD>"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>" + strconv.Itoa(int(ebr_actual.Part_s)) + "</B></FONT></TD>"
					data += "\t\t\t\t</TR>\n"
					data += "\t\t\t\t<TR>\n"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>Part_Name</B></FONT></TD>"
					data += "\t\t\t\t\t<TD ALIGN=\"CENTER\" BGCOLOR=\"white\"><FONT COLOR=\"black\"><B>" + returnstring(ToString(ebr_actual.Name[:])) + "</B></FONT></TD>"
					data += "\t\t\t\t</TR>\n"
				}
				siguiente = ebr_actual.Part_next
			}
		}

	}
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
	color.Green("Report Generate [MBR]")
}
