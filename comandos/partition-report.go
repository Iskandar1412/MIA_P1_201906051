package comandos

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/fatih/color"
)

func returnstring(s string) string {

	if len(s) <= 0 || !(s != "") {
		return " "
	}
	return s
}

func Report_MBR(name string) {
	path := "MIA/P1/Disks/" + string(name) + ".dsk"
	TempDisk, existe := ObtainMBRDisk(path)
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

	}
	data += "\t\t</TABLE>\n"
	data += "\t>];\n}\n"
	ruta := "MIA/P1/Reports/A.dot"
	err := os.WriteFile(ruta, []byte(data), 0644)
	if err != nil {
		color.Red("[REP]: Error al guardar el reporte del MBR")
		return
	}
	imagepath := "MIA/P1/Reports/A.svg"

	cmd := exec.Command("/usr/bin/dot", "-Tsvg", ruta, "-o", imagepath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error to generate img", err)
		return
	}
	color.Green("Image Generate")
}
