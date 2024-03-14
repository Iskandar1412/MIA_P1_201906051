package comandos

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func Report_LS(name string, path string, ruta string, id_disco string) {
	conjunto, route, econ := Obtener_Particion_ID(id_disco)
	if !econ {
		return
	}

	superbloque, esb := ReducirSuperBloqueObtener(route, id_disco, conjunto)
	if !esb {
		return
	}

	data := "digraph G{\n"
	data += "\tnode [shape=plaintext];\n"
	data += "\trankdir=LR;\n\n"
	data += "\ttabla0 [label=<\n"
	data += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>Permisos</B></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>Owner</B></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>Grupo</B></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>Size</B></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>Fecha</B></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>Hora</B></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>Tipo</B></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>Name</B></TD>\n"
	data += "\t\t\t</TR>\n"

	//
	num_inodo := 0
	info_carpetas, eic := Obtener_Reporte_LS("REP", route, superbloque, int32(num_inodo), id_disco)
	if !eic {
		return
	}
	data += info_carpetas
	//

	data += "\t\t</TABLE>\n"
	data += "\t>]\n"
	data += "\n}\n"

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
	color.Green("Report Generate [LS]")
}
