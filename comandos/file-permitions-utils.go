package comandos

import (
	"MIA_P1_201906051/structures"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Obtener_UID_USER_Nombre_Grupo(comando string, id string, user string) (int32, string, bool) {
	contenido_users, ecu := Obtener_Contenido_Archivo_Users(comando, id)
	if !ecu {
		return -1, "", false
	}

	if contenido_users == "" || len(contenido_users) == 0 {
		color.Red("No existe archivo USERS.TXT")
		return -1, "", false
	}

	nombre_usuario := int32(0)
	// password_usuario := ""
	splita := strings.Split(contenido_users, "\n")
	for _, valor := range splita {
		if strings.Contains(valor, ",U,") {
			splitb := strings.Split(valor, ",")
			numero_us, _ := strconv.Atoi(splitb[0])
			nombre_usuario = int32(numero_us)
			// password_usuario = splitb[4]
			if splitb[3] == user {
				return nombre_usuario, splitb[2], true
			}
		}
	}
	return -1, "", false
}

func Obtener_UID_Grupo(comando string, id string, nombre_grupo string) (int32, bool) {
	contenido_users, ecu := Obtener_Contenido_Archivo_Users(comando, id)
	if !ecu {
		return -1, false
	}

	if contenido_users == "" || len(contenido_users) == 0 {
		color.Red("No existe archivo USERS.TXT")
		return -1, false
	}

	nombre_usuario := int32(0)
	// password_usuario := ""
	splita := strings.Split(contenido_users, "\n")
	for _, valor := range splita {
		if strings.Contains(valor, ",G,") {
			splitb := strings.Split(valor, ",")
			numero_us, _ := strconv.Atoi(splitb[0])
			nombre_usuario = int32(numero_us)
			// password_usuario = splitb[4]
			if splitb[2] == nombre_grupo {
				return nombre_usuario, true
			}
		}
	}
	return -1, false
}

func Encontrar_Ruta_Sin_Path(comando string, id string, ruta string) (int32, bool) {
	conjunto, path, er := Obtener_Particion_ID(id)
	if !er {
		return -1, false
	}
	superbloque := structures.SuperBlock{}

	logica := structures.EBR{}
	if conjunto[2] != nil {
		if temp_log, ok := conjunto[2].(structures.EBR); ok {
			v := reflect.ValueOf(temp_log)
			reflect.ValueOf(&logica).Elem().Set(v)
			conjunto[0] = nil
			conjunto[1] = nil
			var esb bool
			// eslogica = true
			superbloque, esb = Obtener_Superbloque(comando, path, ToString(logica.Name[:]))
			if !esb {
				color.Red("[" + comando + "]: Error al obtener superbloque")
				return -1, false
			}
		}
	}

	// esparticion := false
	particion := structures.Partition{}
	if conjunto[0] != nil {
		if temp, ok := conjunto[0].(structures.Partition); ok {
			v := reflect.ValueOf(temp)
			reflect.ValueOf(&particion).Elem().Set(v)
			var esb bool
			superbloque, esb = Obtener_Superbloque(comando, path, ToString(particion.Part_name[:]))
			if !esb {
				color.Red("[" + comando + "]: Error al obtener superbloque")
				return -1, false
			}

			if particion.Part_type == 'E' {
				color.Red("[" + comando + "]: No se puede obtener información de particion extendida")
				return -1, false
			}
		}
	}

	ruta_separada := strings.Split(ruta, "/")
	inodo_ultima_carpeta := int32(0)
	for _, carpeta := range ruta_separada {
		if carpeta == "" {
			inodo_ultima_carpeta = 0
		} else {
			if inodo_ultima_carpeta == -1 {
				color.Red("No se encontró ruta")
				return -1, false
			}
			var eiuc bool
			inodo_ultima_carpeta, eiuc = Encontrar_Carpeta_En_Inodo(comando, path, superbloque.S_inode_start, superbloque.S_block_start, int32(inodo_ultima_carpeta), carpeta)
			if !eiuc {
				return -1, false
			}
		}
	}
	return inodo_ultima_carpeta, true
	//retornar inodo_ultima_carpeta
}

func Crear_Estructura_Carpeta_Nueva(comando string, path string, inicio_inodo int32, inicio_bloque int32, numero_inodos_total int32, numero_bloques_total int32, ruta string, nombre_carpeta string, inicio_bitmap_inodo int32, inicio_bitmap_bloque int32, id_user int32, id_grupo int32) bool {
	// iuid, _ := strconv.Atoi(id_user)
	// igid, _ := strconv.Atoi(id_grupo)
	numero_inodo_disponible, enid := Obtener_Inodo_Disponible(comando, path, inicio_bitmap_inodo, numero_inodos_total)
	if !enid {
		return false
	}
	numero_bloque_disponible, enbd := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
	if !enbd {
		return false
	}
	contenido := structures.Content{}
	contenido.B_name = NameArchivosByte(nombre_carpeta)
	contenido.B_inodo = numero_inodo_disponible

	numero_inodo_padre, einp := Encontrar_Ruta(comando, path, inicio_inodo, inicio_bloque, ruta)
	if !einp {
		return false
	}

	inodo_padre, einp := Obtener_Inodo(comando, path, inicio_inodo, numero_inodo_padre)
	if !einp {
		return false
	}

	if numero_inodo_padre != 0 {
		if !Validar_Permisos(comando, inodo_padre.I_uid, inodo_padre.I_gid, id_user, id_grupo, inodo_padre.I_perm, 2) {
			color.Red("No se tienen permisos necesarios")
			return false
		}
	}
	// fmt.Println(numero_bloque_disponible)

	bloque_padre, ebp := Obtener_Bloque(comando, path, inicio_bloque, inodo_padre.I_block[0])
	if !ebp {
		return false
	}
	lista_bloque_padre := bloque_padre.B_content
	info_padre := lista_bloque_padre[0]
	nombre_padre := ""

	//
	nombre_padre = ToString(info_padre.B_name[:])
	contenido_padre := structures.Content{}
	contenido_padre.B_inodo = numero_inodo_padre
	contenido_padre.B_name = NameArchivosByte(nombre_padre)
	contenido_vacio := structures.Content{B_name: NameArchivosByte(""), B_inodo: int32(-1)}

	var contenido_carpeta []interface{}
	contenido_carpeta = append(contenido_carpeta, contenido_padre)
	contenido_carpeta = append(contenido_carpeta, contenido_vacio)
	contenido_carpeta = append(contenido_carpeta, contenido_vacio)
	contenido_carpeta = append(contenido_carpeta, contenido_vacio)

	Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, numero_bloque_disponible)
	Modificar_Carpeta(comando, path, inicio_bloque, numero_bloque_disponible, contenido_carpeta)
	Modificar_Bitmap(comando, path, inicio_bitmap_bloque, numero_bloque_disponible, 1)

	//creacion inodo
	var lista_apt_bloques [16]int32
	for i := range lista_apt_bloques {
		lista_apt_bloques[i] = int32(-1)
	}
	lista_apt_bloques[0] = numero_bloque_disponible

	fecha := ObFechaInt()
	inodo := structures.Inode{}
	inodo.I_uid = id_user
	inodo.I_gid = id_grupo
	inodo.I_s = int32(0)
	inodo.I_atime = fecha
	inodo.I_ctime = fecha
	inodo.I_mtime = fecha
	inodo.I_block = lista_apt_bloques
	inodo.I_type = int32(0)
	inodo.I_perm = int32(664)

	if !Guardar_Inodo(comando, path, inicio_inodo, inodo, numero_inodo_disponible) {
		color.Red("Error al guardar el inodo")
		return false
	}

	Modificar_Bitmap(comando, path, inicio_bitmap_inodo, numero_inodo_disponible, 1)
	//modificar inodo padre para apuntar al hijo
	return Agregar_Bloque_Lista_Inodos(comando, path, inicio_bloque, numero_inodo_disponible, inicio_bitmap_bloque, numero_bloques_total, nombre_carpeta, inodo_padre, inicio_inodo, numero_inodo_padre, inicio_bitmap_inodo, numero_bloque_disponible)
}

func Crear_Archivo_Vacio_Sin_Path(comando string, id string, ruta string, contenido string, nombre_archivo string, id_user int32, id_grupo int32) bool {
	conjunto, path, er := Obtener_Particion_ID(id)
	if !er {
		return false
	}
	superbloque := structures.SuperBlock{}

	logica := structures.EBR{}
	if conjunto[2] != nil {
		if temp_log, ok := conjunto[2].(structures.EBR); ok {
			v := reflect.ValueOf(temp_log)
			reflect.ValueOf(&logica).Elem().Set(v)
			conjunto[0] = nil
			conjunto[1] = nil
			var esb bool
			// eslogica = true
			superbloque, esb = Obtener_Superbloque(comando, path, ToString(logica.Name[:]))
			if !esb {
				color.Red("[" + comando + "]: Error al obtener superbloque")
				return false
			}
		}
	}

	particion := structures.Partition{}
	if conjunto[0] != nil {
		if temp, ok := conjunto[0].(structures.Partition); ok {
			v := reflect.ValueOf(temp)
			reflect.ValueOf(&particion).Elem().Set(v)
			var esb bool
			// esparticion = true
			superbloque, esb = Obtener_Superbloque(comando, path, ToString(particion.Part_name[:]))
			if !esb {
				color.Red("[" + comando + "]: Error al obtener superbloque")
				return false
			}
			if particion.Part_type == 'E' {
				color.Red("[MKFS]: No se puede obtener información de particion extendida")
				return false
			}
		}
	}

	return Crear_Archivo_Vacio(comando, path, ruta, superbloque, contenido, nombre_archivo, id_user, id_grupo)
}

func TieneFile(valor string) string {
	value := strings.Split(valor, "=")
	return value[1]
}

func Obtener_Contenido_Archivo_Con_Permisos(comando string, id string, ruta string, id_user int32, grupo_user int32) (string, bool) {

	conjunto, path, eco := Obtener_Particion_ID(id)
	if !eco {
		return "", false
	}

	superbloque, esb := ReducirSuperBloqueObtener(path, id, conjunto)
	if !esb {
		return "", false
	}

	ruta_separada := strings.Split(ruta, "/")
	nombre_archivo := ruta_separada[len(ruta_separada)-1]
	ruta_sin_archivo := strings.ReplaceAll(ruta, "/"+nombre_archivo, "")

	numero_inodo_carpeta, enic := Encontrar_Ruta(comando, path, superbloque.S_inode_start, superbloque.S_block_start, ruta_sin_archivo)
	if !enic {
		return "", false
	}

	inodo_carpeta, eic := Obtener_Inodo(comando, path, superbloque.S_inode_start, numero_inodo_carpeta)
	if !eic {
		return "", false
	}

	if !Validar_Permisos(comando, inodo_carpeta.I_uid, inodo_carpeta.I_gid, id_user, grupo_user, inodo_carpeta.I_perm, 4) {
		color.Red("No se tienen permisos necesarios")
		return "", false
	}
	fecha := ObFechaInt()
	inodo_carpeta.I_atime = fecha
	if !Guardar_Inodo(comando, path, superbloque.S_inode_start, inodo_carpeta, numero_inodo_carpeta) {
		return "", false
	}
	content, econ := Obtener_Contenido_Archivo(id, ruta, strconv.Itoa(1), strconv.Itoa(1))
	if !econ {
		return "", false
	}
	return content, true

	// return "", false
}
