package comandos

import (
	"MIA_P1_201906051/structures"
	"strings"

	"github.com/fatih/color"
)

func Obtener_Numero_Estructura_Journaling(id_log string) (int32, bool) {

	// inicio_journaling:= int32(-1)
	conjunto, path, ec := Obtener_Particion_ID(id_log)
	if !ec {
		return 0, false
	}

	superbloque, inicio_journaling, esb := Obtener_Superbloque_Journal_Reducido(path, id_log, conjunto)
	if !esb && (superbloque == structures.SuperBlock{}) {
		color.Red("No existe Journal")
		return -1, false
	}

	return inicio_journaling, true

	// fmt.Println(superbloque, inicio_journaling)

	// return -1, false
}

func Obtener_Nodo_Journaling(comando string, contenido string, fecha int32) (string, bool) {
	operacion := ""
	parametro_actual := ""
	path := ""
	destino := ""
	ruta := ""
	cadena := ""
	estado := 0
	for i := 0; i < len(comando); i++ {
		if estado == 0 {
			if comando[i] == ' ' {
				cadena = strings.ToLower(cadena)
				cadena = strings.ReplaceAll(cadena, " ", "")
				if cadena == "" {
					estado = 0
					continue
				}
				if cadena == "login" || cadena == "mkgrp" || cadena == "rmgrp" || cadena == "mkusr" || cadena == "rmusr" || cadena == "mkfile" || cadena == "cat" || cadena == "remove" || cadena == "edit" || cadena == "rename" || cadena == "mkdir" || cadena == "copy" || cadena == "move" || cadena == "find" || cadena == "chown" || cadena == "chgrp" || cadena == "chmod" || cadena == "pause" || cadena == "recovery" || cadena == "loss" || cadena == "execute" || cadena == "rep" {
					operacion = cadena
					if cadena == "" {
						estado = 0
						continue
					}
					cadena = ""
					estado = 1
				} else {
					estado = 0
				}
				cadena = ""
			} else {
				cadena += string(comando[i])
				cadena = strings.ToLower(cadena)
				if cadena == "logout" {
					operacion = cadena
					cadena = ""
					estado = 1
				}
			}
		} else if estado == 1 {
			if comando[i] == '=' {
				cadena = strings.ToLower(cadena)
				cadena = strings.ReplaceAll(cadena, " ", "")
				if cadena == "-size" || cadena == "-fit" || cadena == "-path" || cadena == "-unit" || cadena == "-name" || cadena == "-type" || cadena == "-delete" || cadena == "-add" || cadena == "-id" || cadena == "-fs" || cadena == "-user" || cadena == "-pass" || cadena == "-grp" || cadena == "-r" || cadena == "-cont" || cadena == "-destino" || cadena == "-ruta" || cadena == "-ugo" {
					parametro_actual = cadena
					parametro_actual = strings.ReplaceAll(parametro_actual, ">", "")
					estado = 2
					cadena = ""
				} else if parametro_actual == "-file" {
					estado = 2
					cadena = ""
				}
			} else {
				if comando[i] != ' ' {
					cadena += string(comando[i])
					if cadena == "-r" {
						if i+1 < len(comando) {
							if comando[i+1] == '=' {
								color.Red("Error con parámetro -r")
								return "", false
							} else if comando[i+1] != ' ' {
								cadena = "-r"
								continue
							}
						}
						// r := 1
						cadena = ""
						estado = 1
					}
				} else if comando[i] == ' ' {
					if cadena == "-r" {
						if i+1 < len(comando) {
							if comando[i+1] == '=' {
								color.Red("Error con parámetro -r")
								return "", false
							}
						}
						// r := 1
						cadena = ""
						estado = 1
					}
				}
			}
		} else if estado == 2 {
			if comando[i] == ' ' {
				if parametro_actual == "-path" {
					path = cadena
				} else if parametro_actual == "-destino" {
					destino = strings.ToLower(cadena)
				} else if parametro_actual == "-ruta" {
					ruta = cadena
				}
				estado = 1
				cadena = ""
			} else if i == (len(comando) - 1) {
				cadena += string(comando[i])
				if parametro_actual == "-path" {
					path = cadena
				} else if parametro_actual == "-destino" {
					destino = strings.ToLower(cadena)
				} else if parametro_actual == "-ruta" {
					ruta = cadena
				}
				cadena = ""
			} else if comando[i] == '"' {
				estado = 3
			} else {
				cadena += string(comando[i])
			}
		} else if estado == 3 {
			if comando[i] == '"' {
				if parametro_actual == "-path" {
					path = cadena
				} else if parametro_actual == "-destino" {
					destino = strings.ToLower(cadena)
				} else if parametro_actual == "-ruta" {
					ruta = cadena
				}
				estado = 1
				cadena = ""
			} else {
				cadena += string(comando[i])
			}
		}
	}

	nodo := "\t<TR>\n"
	nodo += "\t\t<TD ALIGN=\"CENTER\">" + returnstring(operacion) + "</TD>\n"
	dato_path := ""
	if path != "" {
		dato_path = returnstring(path)
	}
	if ruta != "" {
		dato_path = returnstring(ruta)
	}
	if destino != "" {
		dato_path = returnstring(destino)
	}

	nodo += "\t\t<TD ALIGN=\"CENTER\">" + returnstring(ToString([]byte(dato_path))) + "</TD>\n"
	nodo += "\t\t<TD ALIGN=\"CENTER\">" + returnstring(ToString([]byte(contenido))) + "</TD>\n"
	nodo += "\t\t<TD ALIGN=\"CENTER\">" + IntFechaToStr(fecha) + "</TD>\n"
	nodo += "\t</TR>\n"

	return nodo, true
}
