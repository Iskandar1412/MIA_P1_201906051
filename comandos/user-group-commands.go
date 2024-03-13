package comandos

import (
	"MIA_P1_201906051/structures"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Values_LOGIN(instructions []string) (Usuario, bool) {
	var _uss string
	var _pas string
	var _id string
	UsuarioLogeado.Id_Particion = IDParticionByte("")
	UsuarioLogeado.User = NameArchivosByte("")
	UsuarioLogeado.Password = NameArchivosByte("")

	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "id") {
			var value = TieneID("LOGIN", valor)
			_id = value
			continue
		} else if strings.HasPrefix(strings.ToLower(valor), "pass") {
			var value = TienePassword("LOGIN", valor)
			_pas = value
			continue
		} else if strings.HasPrefix(strings.ToLower(valor), "user") {
			var value = TieneUser("LOGIN", valor)
			_uss = value
			continue
		} else {
			color.Yellow("[LOGIN]: Atributo no reconocido")
		}
	}
	if _id == "" || len(_id) == 0 || len(_id) > 4 {
		return UsuarioLogeado, false
	} else if _pas == "" || len(_pas) == 0 || len(_pas) > 10 {
		return UsuarioLogeado, false
	} else if _uss == "" || len(_uss) == 0 || len(_uss) > 10 {
		return UsuarioLogeado, false
	}
	UsuarioLogeado.Id_Particion = IDParticionByte(_id)
	UsuarioLogeado.User = NameArchivosByte(_uss)
	UsuarioLogeado.Password = NameArchivosByte(_pas)
	return UsuarioLogeado, true
}

func LOGIN_EXECUTE(uss string, password string, id_disco string) {
	contenido_users, ecu := Obtener_Contenido_Archivo_Users("LOGIN", id_disco)
	if !ecu {
		return
	}

	if contenido_users == "" || len(contenido_users) == 0 {
		color.Red("No existe archivo USERS.TXT")
		return
	}

	nombre_usuario := ""
	password_usuario := ""
	splita := strings.Split(contenido_users, "\n")
	for _, valor := range splita {
		if strings.Contains(valor, ",U,") {
			splitb := strings.Split(valor, ",")
			nombre_usuario = splitb[3]
			password_usuario = splitb[4]
			if nombre_usuario == uss && password_usuario == password {

				color.Green("[LOGIN]: Usuario correcto")
				UsuarioLogeado.Id_Particion = IDParticionByte(id_disco)
				UsuarioLogeado.User = NameArchivosByte(uss)
				UsuarioLogeado.Password = NameArchivosByte(password)
				// Guardar en Journaling
				numero_journaling_disponible, enjd := Obtener_Journaling_Disponible("LOGIN", id_disco)
				if !enjd {
					return
				}
				if numero_journaling_disponible != -1 {
					comando_journaling := "login -user=" + uss + " -pass=" + password + " -id=" + id_disco
					fecha := ObFechaInt()
					journaling := structures.Journal{}
					journaling.J_state = int8(1)
					journaling.J_command = DevolverContenidoJournal(comando_journaling)
					journaling.J_date = fecha
					journaling.J_content = DevolverContenidoArchivo(uss)
					Modificar_Journaling("LOGIN", id_disco, numero_journaling_disponible, journaling)
					return
				}
			}
		}
	}
	color.Red("[LOGIN]: Usuario incorrecto")
}

func LOGOUT_EXECUTE() {
	if (ToString(UsuarioLogeado.User[:]) == "") && (ToString(UsuarioLogeado.Password[:]) == "") && (ToString(UsuarioLogeado.Id_Particion[:]) == "") {
		color.Red("[LOGOUT]: No hay usuario que haya iniciado seción")
		return
	}
	id_usuario := ToString(UsuarioLogeado.Id_Particion[:])
	numero_journaling_disponible, enjd := Obtener_Journaling_Disponible("LOGGOUT", id_usuario)
	if !enjd {
		return
	}
	if numero_journaling_disponible != -1 {
		// fmt.Println("siuuu")
		comando_journaling := "logout"
		fecha := ObFechaInt()
		journaling := structures.Journal{}
		journaling.J_state = int8(1)
		journaling.J_command = DevolverContenidoJournal(comando_journaling)
		journaling.J_date = fecha

		Modificar_Journaling("LOGOUT", id_usuario, numero_journaling_disponible, journaling)

		UsuarioLogeado.Id_Particion = IDParticionByte("")
		UsuarioLogeado.User = NameArchivosByte("")
		UsuarioLogeado.Password = NameArchivosByte("")

	}
	color.Blue("Logout Exitoso")
}

func Values_MKGRP(instructions []string) (string, bool) {
	var _name string
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "name") {
			var value = TieneNombre("MKGRP", valor)
			if len(value) > 10 {
				color.Red("No puede ser nun nombre mayor a 10")
				_name = ""
				break
			} else {
				_name = value
			}
		} else {
			color.Yellow("[MKGRP]: Atributo no reconocido")
		}
	}
	if _name == "" || len(_name) == 0 {
		return "", false
	}
	return _name, true
}

func MKGRP_EXECUTE(_name string) {
	if ToString(UsuarioLogeado.User[:]) != "root" {
		color.Red("Solo root puede usar el comando")
		return
	}
	contenido_users, ecu := Obtener_Contenido_Archivo_Users("MKGRP", ToString(UsuarioLogeado.Id_Particion[:]))
	if !ecu {
		return
	}
	// contenido_users = strings.ReplaceAll(contenido_users, "\t", "")

	ultimo_grupo := ""
	numero := ""
	tipo := ""
	grupo := ""
	usuario := ""
	password := ""
	estado := 0
	for i := 0; i < len(contenido_users); i++ {
		if estado == 0 {
			if contenido_users[i] == ',' {
				estado = 1
			} else {
				numero += string(contenido_users[i])
			}
		} else if estado == 1 {
			if contenido_users[i] == ',' {
				estado = 2
			} else {
				tipo += string(contenido_users[i])
			}
		} else if estado == 2 {
			if contenido_users[i] == ',' {
				estado = 3
			} else if contenido_users[i] == '\n' {
				ultimo_grupo = numero
				if grupo == ToString([]byte(_name)) {
					//grupo ya existente
					color.Yellow("Grupo ya existente")
					return
				}
				numero = ""
				tipo = ""
				grupo = ""
				usuario = ""
				password = ""
				estado = 0
			} else if contenido_users[i] != ' ' {
				grupo += string(contenido_users[i])
			}
		} else if estado == 3 {
			if contenido_users[i] == ',' {
				estado = 4
			} else if contenido_users[i] != ' ' {
				usuario += string(contenido_users[i])
			}
		} else if estado == 4 {
			if contenido_users[i] == '\n' {
				numero = ""
				tipo = ""
				grupo = ""
				usuario = ""
				password = ""
				estado = 0
			} else if contenido_users[i] != ' ' {
				password += string(contenido_users[i])
			}
		}
	}
	color.Blue("Creando grupo nuevo... " + _name)

	nombre_temp, eng := strconv.Atoi(ultimo_grupo)
	if eng != nil {
		return
	}
	nombre_temp = nombre_temp + 1
	contenido_users += strconv.Itoa(nombre_temp)
	contenido_users += ",G,"
	contenido_users += _name
	contenido_users += "\n"

	//Modificar_Contendio_Archivo
	Modificar_Contenido_Archivo("MKGRP", ToString(UsuarioLogeado.Id_Particion[:]), "/users.txt", contenido_users)
	color.Blue("Grupo creado Exitosamente")
	Recalcular_Size_Carpetas("MKGRO", ToString(UsuarioLogeado.Id_Particion[:]))
	numero_journaling_disponible, enjd := Obtener_Journaling_Disponible("MKGRP", ToString(UsuarioLogeado.Id_Particion[:]))
	if !enjd {
		return
	}
	if numero_journaling_disponible != -1 {
		//solo pasa en caos de ser ext3
		comando_journaling := "mkgrp -name=" + _name
		fecha := ObFechaInt()
		journaling := structures.Journal{}
		journaling.J_state = int8(1)
		journaling.J_command = DevolverContenidoJournal(comando_journaling)
		journaling.J_date = fecha

		Modificar_Journaling("MKGRP", ToString(UsuarioLogeado.Id_Particion[:]), numero_journaling_disponible, journaling)
	}
}

func Values_RMGRP(instructions []string) (string, bool) {
	var _name string
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "name") {
			var value = TieneNombre("RMGRP", valor)
			if len(value) > 10 {
				color.Red("No puede ser nun nombre mayor a 10")
				_name = ""
				break
			} else {
				_name = value
			}
		} else {
			color.Yellow("[RMGRP]: Atributo no reconocido")
		}
	}
	if _name == "" || len(_name) == 0 {
		return "", false
	}
	return _name, true
}

func RMGRP_EXECUTE(_name string) {
	if ToString(UsuarioLogeado.User[:]) != "root" {
		color.Red("Solo root puede usar el comando")
		return
	}
	contenido_users, ecu := Obtener_Contenido_Archivo_Users("MKGRP", ToString(UsuarioLogeado.Id_Particion[:]))
	if !ecu {
		return
	}

	numero := ""
	tipo := ""
	grupo := ""
	usuario := ""
	password := ""

	estado := 0
	texto := ""
	for i := 0; i < len(contenido_users); i++ {
		if estado == 0 {
			if contenido_users[i] == ',' {
				estado = 1
			} else {
				numero += string(contenido_users[i])
			}
		} else if estado == 1 {
			if contenido_users[i] == ',' {
				estado = 2
			} else {
				tipo += string(contenido_users[i])
			}
		} else if estado == 2 {
			if contenido_users[i] == ',' {
				estado = 3
			} else if contenido_users[i] == '\n' {
				if grupo == _name {
					numero_int, _ := strconv.Atoi(numero)
					if numero_int == 0 {
						color.Red("Grupo no existe.")
						return
					}
					color.Blue("Grupo encontrado, Borrando... " + _name)
					lineas := strings.Split(contenido_users, "\n")
					var newText string
					for _, linea := range lineas {
						if !strings.Contains(linea, _name) {
							if linea == "" {
								continue
							}
							newText += linea + "\n"
						}
					}
					texto = newText
					color.Blue("Borrado exitosamente... " + _name)
					break
				}
				numero = ""
				tipo = ""
				grupo = ""
				usuario = ""
				password = ""
				estado = 0
			} else if contenido_users[i] != ' ' {
				grupo += string(contenido_users[i])
			}
		} else if estado == 3 {
			if contenido_users[i] == ',' {
				estado = 4
			} else if contenido_users[i] != ' ' {
				usuario += string(contenido_users[i])
			}
		} else if estado == 4 {
			if contenido_users[i] == '\n' {
				numero = ""
				tipo = ""
				grupo = ""
				usuario = ""
				password = ""
				estado = 0
			} else if contenido_users[i] != ' ' {
				password += string(contenido_users[i])
			}
		}
	}
	if texto != "" {
		Recalcular_Size_Carpetas("RMGRP", ToString(UsuarioLogeado.Id_Particion[:]))
		Modificar_Contenido_Archivo("RMGRP", ToString(UsuarioLogeado.Id_Particion[:]), "/users.txt", texto)
		color.Magenta("Grupo eliminado")
		numero_journaling_disponible, enjd := Obtener_Journaling_Disponible("RMGRP", ToString(UsuarioLogeado.Id_Particion[:]))
		if !enjd {
			return
		}
		if numero_journaling_disponible != -1 {
			//solo pasa en caos de ser ext3
			comando_journaling := "rmgrp -name=" + _name
			fecha := ObFechaInt()
			journaling := structures.Journal{}
			journaling.J_state = int8(1)
			journaling.J_command = DevolverContenidoJournal(comando_journaling)
			journaling.J_date = fecha

			Modificar_Journaling("RMGRP", ToString(UsuarioLogeado.Id_Particion[:]), numero_journaling_disponible, journaling)
		}
	}
}

func Values_MKUSR(instructions []string) (string, string, string, bool) {
	var _nam string
	var _pas string
	var _grp string
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "user") {
			var value = TieneUser("MKUSR", valor)
			if len(value) > 10 {
				color.Red("No se puede ser un nombre mayor a 10 caracteres")
				_nam = ""
				break
			} else {
				_nam = value
				continue
			}
		} else if strings.HasPrefix(strings.ToLower(valor), "pass") {
			var value = TienePassword("MKUSR", valor)
			if len(value) > 10 {
				color.Red("No se puede ser un password mayor a 10 caracteres")
				_pas = ""
				break
			} else {
				_pas = value
				continue
			}
		} else if strings.HasPrefix(strings.ToLower(valor), "grp") {
			var value = TieneGRP("MKUSR", valor)
			if len(value) > 10 {
				color.Red("No se puede ser un grp mayor a 10 caracteres")
				_grp = ""
				break
			} else {
				_grp = value
				continue
			}
		} else {
			color.Yellow("[MKUSR]: Atrubuto no reconocido")
		}
	}
	if _nam == "" || len(_nam) == 0 || len(_nam) > 10 {
		return "", "", "", false
	} else if _pas == "" || len(_pas) == 0 || len(_pas) > 10 {
		return "", "", "", false
	} else if _grp == "" || len(_grp) == 0 || len(_grp) > 10 {
		return "", "", "", false
	}

	return _nam, _pas, _grp, true
}

func MKUSR_EXECUTE(_us string, _pas string, _grp string) {
	if ToString(UsuarioLogeado.User[:]) != "root" {
		color.Red("Solo root puede usar el comando")
		return
	}
	contenido_users, ecu := Obtener_Contenido_Archivo_Users("MKGRP", ToString(UsuarioLogeado.Id_Particion[:]))
	if !ecu {
		return
	}

	// numero_usuario := ""
	tipo_usuario := ""
	grupo_usuario := ""
	nombre_usuario := ""
	password_usuario := ""
	contador := 1
	existe_grupo := false
	splita := strings.Split(contenido_users, "\n")
	for _, valor := range splita {
		if strings.Contains(valor, ",U,") {
			splitb := strings.Split(valor, ",")
			// numero_usuario = splitb[0]
			tipo_usuario = splitb[1]
			// grupo_usuario = splitb[2]
			nombre_usuario = splitb[3]
			// password_usuario = splitb[4]
			if nombre_usuario == _us {
				// numero_usuario = ""
				tipo_usuario = ""
				// grupo_usuario = ""
				nombre_usuario = ""
				// password_usuario = ""
				color.Red("[MKUSR]: Ya existe un usuario con ese nombre")
				return
			}
			contador += 1
		}
	}

	for _, valor := range splita {
		//verificacion grupos
		if strings.Contains(valor, ",G,") {
			splitb := strings.Split(valor, ",")
			nombre_grupo := splitb[2]
			if nombre_grupo == _grp {
				existe_grupo = true
				break
			}
			existe_grupo = false
		}
	}
	if !existe_grupo {
		color.Red("[MKUSR]: Grupo no existe")
		return
	}

	nombre_usuario = _us
	grupo_usuario = _grp
	password_usuario = _pas

	// numero_usuario_int :=
	contenido_usuario := strconv.Itoa(contador) + "," + tipo_usuario + "," + grupo_usuario + "," + nombre_usuario + "," + password_usuario + "\n"
	contenido_users += contenido_usuario
	Modificar_Contenido_Archivo("MKUSR", ToString(UsuarioLogeado.Id_Particion[:]), "/users.txt", contenido_users)

	// fmt.Println(nombre_usuario, password_usuario, numero_usuario, grupo_usuario, tipo_usuario)
	Recalcular_Size_Carpetas("MKUSR", ToString(UsuarioLogeado.Id_Particion[:]))
	numero_journaling_disponible, enjd := Obtener_Journaling_Disponible("MKUSR", ToString(UsuarioLogeado.Id_Particion[:]))
	if !enjd {
		return
	}
	if numero_journaling_disponible != -1 {
		//solo pasa en caos de ser ext3
		comando_journaling := "mkusr -name=" + nombre_usuario + " -pass=" + password_usuario + " -grp=" + grupo_usuario
		fecha := ObFechaInt()
		journaling := structures.Journal{}
		journaling.J_state = int8(1)
		journaling.J_command = DevolverContenidoJournal(comando_journaling)
		journaling.J_date = fecha

		Modificar_Journaling("MKUSR", ToString(UsuarioLogeado.Id_Particion[:]), numero_journaling_disponible, journaling)
	}
}

func Values_RMUSR(instructions []string) (string, bool) {
	var _name string
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "user") {
			var value = TieneUser("RMUSR", valor)
			if len(value) > 10 {
				color.Red("No puede ser nun nombre mayor a 10")
				_name = ""
				break
			} else {
				_name = value
			}
		} else {
			color.Yellow("[RMUSR]: Atributo no reconocido")
		}
	}
	if _name == "" || len(_name) == 0 {
		return "", false
	}
	return _name, true
}

func RMUSR_EXECUTE(_user string) {
	if ToString(UsuarioLogeado.User[:]) != "root" {
		color.Red("Solo root puede usar el comando")
		return
	}
	contenido_users, ecu := Obtener_Contenido_Archivo_Users("MKGRP", ToString(UsuarioLogeado.Id_Particion[:]))
	if !ecu {
		return
	}

	// tipo_usuario := ""
	// grupo_usuario := ""
	nombre_usuario := ""
	// password_usuario := ""
	// contador := 1
	// existe_grupo := false
	existe_usuario := false
	splita := strings.Split(contenido_users, "\n")
	for _, valor := range splita {
		if strings.Contains(valor, ",U,") {
			splitb := strings.Split(valor, ",")
			nombre_usuario = splitb[3]
			if nombre_usuario == _user {
				nombre_usuario = ""
				color.Magenta("[RMUSR]: Usuario Encontrado")
				existe_usuario = true
				break
			}
			existe_usuario = false
			// contador += 1
		}
	}
	if !existe_usuario {
		color.Red("[RMUSR]: Usuario no encontrado")
		return
	}

	color.Blue("[RMUSR]: Eliminado Usuario... " + _user)
	texto := ""
	for _, valor := range splita {
		if !strings.Contains(valor, _user) {
			if valor == "" {
				continue
			}
			texto += valor + "\n"
		}
	}

	Modificar_Contenido_Archivo("RMUSR", ToString(UsuarioLogeado.Id_Particion[:]), "/users.txt", texto)
	Recalcular_Size_Carpetas("RMUSR", ToString(UsuarioLogeado.Id_Particion[:]))
	numero_journaling_disponible, enjd := Obtener_Journaling_Disponible("RMUSR", ToString(UsuarioLogeado.Id_Particion[:]))
	if !enjd {
		return
	}
	if numero_journaling_disponible != -1 {
		//solo para ext3
		comando_journaling := "rmusr -name=" + _user
		fecha := ObFechaInt()
		journaling := structures.Journal{}
		journaling.J_state = int8(1)
		journaling.J_command = DevolverContenidoJournal(comando_journaling)
		journaling.J_date = fecha

		Modificar_Journaling("RMUSR", ToString(UsuarioLogeado.Id_Particion[:]), numero_journaling_disponible, journaling)
	}
	color.Green("[RMUSR]: Usuario Eliminado exitosamente")
}
