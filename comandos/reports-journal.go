package comandos

import (
	"MIA_P1_201906051/structures"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func Report_Journal(name string, path string, ruta string, id_disco string) {

	conjunto, route, eco := Obtener_Particion_ID(id_disco)
	if !eco {
		return
	}

	superbloque, esb := ReducirSuperBloqueObtener(route, id_disco, conjunto)
	if !esb && (superbloque == structures.SuperBlock{}) {
		return
	}

	data := "digraph G {"
	data += "\tnode [shape=plaintext]\n"
	data += "\trankdir=LR\n\n"
	data += "\ttabla0 [label=<\n"
	data += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"2\">\n"
	data += "\t\t\t<TR>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"150\" ><B>Operacion</B></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"500\" ><B>Path</B></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"500\" ><B>Contenido</B></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"150\" ><B>Fecha</B></TD>\n"
	data += "\t\t\t</TR>\n"
	//-------

	//-obtener numero del journaling
	contador, econt := Obtener_Numero_Estructura_Journaling(id_disco)
	if !econt {
		return
	}

	for i := 0; i < int(contador); i++ {
		journaling_actual, eja := Obtener_Journaling("REP", id_disco, i)
		if !eja {
			return
		}
		if journaling_actual.J_state == 1 {
			//obtener nodos del journaling
			codigo_nodo, ecn := Obtener_Nodo_Journaling(ToString(journaling_actual.J_command[:]), ToString(journaling_actual.J_content[:]), journaling_actual.J_date)
			if !ecn {
				return
			}
			data += codigo_nodo
		}
	}
	// fmt.Println(contador, superbloque)

	//-------
	data += "\t\t</TABLE>\n"
	data += "\t>]\n"
	data += "\n}\n"

	//
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
	color.Green("Report Generate [Journal]")
}
