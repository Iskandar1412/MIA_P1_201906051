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
	estado := 0
	numero := ""
	tipo := ""
	grupo := ""
	usuario := ""
	pass := ""
	// contenido_users = strings.ReplaceAll(contenido_users, "\t", "")
	// fmt.Println(contenido_users)
	if contenido_users == "" || len(contenido_users) == 0 {
		color.Red("No existe archivo USERS.TXT")
		return
	}
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
				numero = ""
				tipo = ""
				grupo = ""
				usuario = ""
				pass = ""
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
				if usuario == uss {
					if password == pass {
						color.Green("[LOGIN]: Usuario correcto")
						UsuarioLogeado.Id_Particion = IDParticionByte(id_disco)
						UsuarioLogeado.User = NameArchivosByte(usuario)
						UsuarioLogeado.Password = NameArchivosByte(password)
						// Guardar en Journaling
						numero_journaling_disponible, enjd := Obtener_Journaling_Disponible("LOGIN", id_disco)
						if !enjd {
							return
						}
						if numero_journaling_disponible != -1 {
							comando_journaling := "login -user=" + usuario + " -pass=" + pass + " -id=" + id_disco
							fecha := ObFechaInt()
							journaling := structures.Journal{}
							journaling.J_state = int8(1)
							journaling.J_command = DevolverContenidoJournal(comando_journaling)
							journaling.J_date = fecha
							journaling.J_content = DevolverContenidoArchivo(usuario)
							Modificar_Journaling("LOGIN", id_disco, numero_journaling_disponible, journaling)
						}
						return
					}
				}
			} else if contenido_users[i] != ' ' {
				pass += string(contenido_users[i])
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
}
