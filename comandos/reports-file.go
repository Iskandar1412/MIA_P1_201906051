package comandos

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

func Report_FILE(name string, path string, ruta string, id_disco string) {

	conjunto, route, eco := Obtener_Particion_ID(id_disco)
	if !eco && conjunto[0] == nil {
		return
	}

	comando := "REP"
	// usr, grp_temp, eur := Obtener_UID_USER_Nombre_Grupo(comando, id_disco, ToString(UsuarioLogeado.User[:]))
	// if !eur {
	// 	return
	// }
	// grp, egrp := Obtener_UID_Grupo(comando, id_disco, grp_temp)
	// if !egrp {
	// 	return
	// }

	contenido_archivo, eca := Obtener_Contenido_Archivo(id_disco, route, strconv.Itoa(1), strconv.Itoa(1))
	if !eca {
		return
	}

	if contenido_archivo == "" || len(contenido_archivo) == 0 {
		color.Red("[" + comando + "]: Archivo Vacio")
		return
	}

	rep, err := os.Create(path + "/" + name)
	if err != nil {
		fmt.Println("[" + comando + "]: Error al crear archivo de reporte")
		return
	}
	defer rep.Close()

	err = rep.Truncate(0)
	if err != nil {
		return
	}

	_, err = rep.WriteString(fmt.Sprintf("%v\t", contenido_archivo))
	if err != nil {
		fmt.Println("[" + comando + "]: Error al escribir datos del reporte")
		return
	}
	color.Green("Report Generate [file]")
}
