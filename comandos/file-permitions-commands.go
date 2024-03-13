package comandos

import (
	"MIA_P1_201906051/structures"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// --------------------------------------- path r size cont error
func Values_MKFILE(instructions []string) (string, bool, int32, string, bool) {
	//  path r size cont error
	_r := false
	_size := int32(0)
	_cont := ""
	_path := ""
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "path") {
			var value = TienePathV2("MKFILE", valor)
			_path = value
		} else if strings.HasPrefix(strings.ToLower(valor), "r") {
			var value = TieneR("MKFILE", valor)
			_r = value
		} else if strings.HasPrefix(strings.ToLower(valor), "cont") {
			var value = TieneCont("MKFILE", valor)
			_cont = value
		} else if strings.HasPrefix(strings.ToLower(valor), "size") {
			var value = TieneSizeV2("MKFILE", valor)
			_size = value
		} else {
			color.Yellow("[MKFILE]: Atributo no reconocido")
			_path = ""
			break
		}
	}
	if _path == "" {
		color.Red("[MKFILE]: No hay path")
		return "", false, 0, "", false
	}
	if _size == -1 {
		return "", false, 0, "", false
		// returnstring()
	}
	return _path, _r, _size, _cont, true
}

func MKFILE_EXECUTE(path string, r bool, size int32, cont string) {
	comando := "MKFILE"
	comando_journaling := "mkfile -path=\"" + path + "\" -size=" + strconv.Itoa(int(size))
	if r {
		comando_journaling += " -r"
	}
	if cont != "" {
		comando_journaling += " -cont=\"" + cont + "\""
	}

	conjunto, path, er := Obtener_Particion_ID(ToString(UsuarioLogeado.Id_Particion[:]))
	if !er {
		return
	}

	if conjunto[0] == nil {
		return
	}

	usr, grp_temp, eur := Obtener_UID_USER_Nombre_Grupo(comando, ToString(UsuarioLogeado.Id_Particion[:]), ToString(UsuarioLogeado.User[:]))
	if !eur {
		return
	}
	grp, egrp := Obtener_UID_Grupo(comando, ToString(UsuarioLogeado.Id_Particion[:]), grp_temp)
	if !egrp {
		return
	}

	contenido_archivo := ""
	if size >= 0 {
		contador := 0
		for i := 0; i < int(size); i++ {
			contenido_archivo += fmt.Sprint(contador)
			contador++
			if contador == 10 {
				contador = 0
			}
		}
	}

	if cont != "" {
		contenidoArchivo, ecA := os.ReadFile(cont)
		if ecA != nil {
			color.Red("[" + comando + "]: Error al leer archivo")
			return
		}
		fmt.Println("Contenido: -------- ")
		contenido_archivo = ToString(contenidoArchivo[:])
	}

	ruta_separada := strings.Split(path, "/")
	cantidad_carpetas := len(ruta_separada)
	nombre_archivo := ruta_separada[cantidad_carpetas-1]
	ruta_sin_archivo := strings.ReplaceAll(path, "/"+nombre_archivo, "")

	//verificar que existe ruta
	numero_inodo_ultima_carpeta, _ := Encontrar_Ruta_Sin_Path(comando, ToString(UsuarioLogeado.Id_Particion[:]), ruta_sin_archivo)
	if numero_inodo_ultima_carpeta == -1 {
		//como no existe la carpeta ni ruta
		if r {
			//crear carpetas en caso que no existan
			fmt.Println("Creando rutas (r activo)")
			if !Crear_Ruta_Carpetas(comando, ToString(UsuarioLogeado.Id_Particion[:]), ruta_sin_archivo, strconv.Itoa(int(usr)), strconv.Itoa(int(grp))) {
				color.Red("Error al crear carpetas")
				return
			}
			fmt.Println("Creando Arhivo.... " + nombre_archivo)
			//crear archivo
			if !Crear_Archivo_Vacio_Sin_Path(comando, ToString(UsuarioLogeado.Id_Particion[:]), ruta_sin_archivo, contenido_archivo, nombre_archivo, usr, grp) {
				color.Red("Error al crear archivo")
			} else {
				color.Green("Archivo Creado")
				numero_journaling_disponible, enjd := Obtener_Journaling_Disponible(comando, ToString(UsuarioLogeado.Id_Particion[:]))
				if !enjd {
					return
				}
				if numero_journaling_disponible != -1 {
					//ext3
					fecha := ObFechaInt()
					journaling := structures.Journal{}
					journaling.J_state = int8(1)
					journaling.J_command = DevolverContenidoJournal(comando_journaling)
					journaling.J_date = fecha
					journaling.J_content = DevolverContenidoArchivo(contenido_archivo)
					Modificar_Journaling(comando, ToString(UsuarioLogeado.Id_Particion[:]), numero_journaling_disponible, journaling)
				}
			}
		} else {
			color.Red("[" + comando + "]: Ruta no encontrada")
		}
	} else {
		fmt.Println("Creando archivo... " + nombre_archivo)
		if !Crear_Archivo_Vacio_Sin_Path(comando, ToString(UsuarioLogeado.Id_Particion[:]), ruta_sin_archivo, contenido_archivo, nombre_archivo, usr, grp) {
			color.Red("No se pudo crear el archivo")
		} else {
			color.Green("Archivo Creado")
			numero_journaling_disponible, enjd := Obtener_Journaling_Disponible(comando, ToString(UsuarioLogeado.Id_Particion[:]))
			if !enjd {
				return
			}
			if numero_journaling_disponible != -1 {
				//ext3
				fecha := ObFechaInt()
				journaling := structures.Journal{}
				journaling.J_state = int8(1)
				journaling.J_command = DevolverContenidoJournal(comando_journaling)
				journaling.J_date = fecha
				journaling.J_content = DevolverContenidoArchivo(contenido_archivo)
				Modificar_Journaling(comando, ToString(UsuarioLogeado.Id_Particion[:]), numero_journaling_disponible, journaling)
			}
		}
	}

}
