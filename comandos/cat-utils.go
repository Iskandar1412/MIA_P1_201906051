package comandos

import (
	"fmt"
	"regexp"

	"github.com/fatih/color"
)

func Values_CAT(instructions []string) ([]string, bool) {
	var stringVacio []string
	var string_enviar []string
	// var er bool
	re := regexp.MustCompile(`file\d{1,2}=.*`)

	for _, value := range instructions {
		if re.MatchString(value) {
			val := TieneFile(value)
			string_enviar = append(string_enviar, val)
		} else {
			color.Yellow("[CAT]: Valor no reconocido")
			return stringVacio, false
		}
	}
	return string_enviar, true
}

func CAT_EXECUTE(values []string) {
	if ToString(UsuarioLogeado.User[:]) == "" {
		color.Red("No hay usuario con seción iniciada")
		return
	}
	comando := "CAT"
	id_usuario, grupo, eco := Obtener_UID_USER_Nombre_Grupo(comando, ToString(UsuarioLogeado.Id_Particion[:]), ToString(UsuarioLogeado.User[:]))
	if !eco {
		return
	}

	id_grupo, eig := Obtener_UID_Grupo(comando, ToString(UsuarioLogeado.Id_Particion[:]), grupo)
	if !eig {
		return
	}

	unido := ""
	for _, val := range values {
		contenido, econ := Obtener_Contenido_Archivo_Con_Permisos(comando, ToString(UsuarioLogeado.Id_Particion[:]), val, id_usuario, id_grupo)
		if !econ {
			continue
		}
		unido += "\n"
		unido += contenido
	}

	fmt.Println(unido)
}
