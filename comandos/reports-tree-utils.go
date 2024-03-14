package comandos

import (
	"MIA_P1_201906051/structures"
	"fmt"
	"reflect"
	"strings"
)

func Obtener_Codigo_Tree(comando string, path string, superbloque structures.SuperBlock, num_inodo int32, id_log string) (string, bool) {
	respuesta := ""
	inodo, eino := Obtener_Inodo(comando, path, superbloque.S_inode_start, num_inodo)
	if !eino {
		return "", false
	}
	data := ""
	conexiones := ""

	data += "node [shape=plaintext]; \n"
	data += "\tinodo" + fmt.Sprint(num_inodo) + " [label=<\n"
	data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
	//
	data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\" COLSPAN=\"2\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(num_inodo) + "</B></FONT></TD></TR>\n"
	data += "\t\t\t<TR>"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\"><B>I_UID</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(inodo.I_uid) + "</FONT></TD>\n"
	data += "\t\t\t</TR>"
	data += "\t\t\t<TR>"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\"><B>I_GID</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(inodo.I_gid) + "</FONT></TD>\n"
	data += "\t\t\t</TR>"
	data += "\t\t\t<TR>"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\"><B>I_s</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(inodo.I_s) + "</FONT></TD>\n"
	data += "\t\t\t</TR>"
	data += "\t\t\t<TR>"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\"><B>I_atime</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\">" + IntFechaToStr(inodo.I_atime) + "</FONT></TD>\n"
	data += "\t\t\t</TR>"
	data += "\t\t\t<TR>"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\"><B>I_ctime</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\">" + IntFechaToStr(inodo.I_ctime) + "</FONT></TD>\n"
	data += "\t\t\t</TR>"
	data += "\t\t\t<TR>"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\"><B>I_mtime</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\">" + IntFechaToStr(inodo.I_mtime) + "</FONT></TD>\n"
	data += "\t\t\t</TR>"
	for i := range inodo.I_block {
		apt := inodo.I_block[i]
		data += "\t\t\t<TR>"
		data += "\t\t\t\t<TD port=\"apte" + fmt.Sprint(i) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\"><B>apt" + fmt.Sprint(i) + "</B></FONT></TD>\n"
		data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(i) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt) + "</FONT></TD>\n"
		data += "\t\t\t</TR>"
	}
	data += "\t\t\t<TR>"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\"><B>I_type</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(inodo.I_type) + "</FONT></TD>\n"
	data += "\t\t\t</TR>"
	data += "\t\t\t<TR>"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\"><B>I_perm</B></FONT></TD>\n"
	data += "\t\t\t\t<TD ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#71B4D6\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(inodo.I_perm) + "</FONT></TD>\n"
	data += "\t\t\t</TR>"
	//
	data += "\t\t</table>\n"
	data += "\t>];\n\n"
	if inodo.I_type == 0 {
		for i := range inodo.I_block {
			apt := inodo.I_block[i]
			if apt != -1 {
				if i == 13 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					apuntadores := bloque_apuntador1.B_pointers
					data += "node [shape=plaintext]; \n"
					data += "\tbloque" + fmt.Sprint(apt) + " [label=<\n"
					data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
					data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt) + "</B></FONT></TD></TR>\n"
					//
					for j := range apuntadores {
						apt1 := apuntadores[j]
						data += "\t\t\t<TR>"
						data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt1) + "</FONT></TD>\n"
						data += "\t\t\t</TR>"
					}
					//
					data += "\t\t</table>\n"
					data += "\t>];\n\n"
					data += "inodo" + fmt.Sprint(num_inodo) + ":apts" + fmt.Sprint(i) + " -> bloque" + fmt.Sprint(apt) + ":titulo;\n"

					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
							if !ebc {
								return "", false
							}
							info_carpeta := bloque_carpeta.B_content
							contenido1 := info_carpeta[0]
							contenido2 := info_carpeta[1]
							contenido3 := info_carpeta[2]
							contenido4 := info_carpeta[3]
							temp := reflect.TypeOf(structures.Content{})
							if (temp == reflect.TypeOf(contenido1)) && (temp == reflect.TypeOf(contenido2)) && (temp == reflect.TypeOf(contenido3)) && (temp == reflect.TypeOf(contenido4)) {
								data += "node [shape=plaintext]; \n"
								data += "\tbloque" + fmt.Sprint(apt1) + " [label=<\n"
								data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
								data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt1) + "</B></FONT></TD></TR>\n"
								//
								for z := range info_carpeta {
									dat := info_carpeta[z]
									data += "\t\t\t<TR>"
									data += "\t\t\t\t<TD port=\"apte" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + returnstring(ToString(dat.B_name[:])) + "</B></FONT></TD>\n"
									data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(dat.B_inodo) + "</FONT></TD>\n"
									data += "\t\t\t</TR>"
								}
								//
								data += "\t\t</table>\n"
								data += "\t>];\n\n"
								data += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(a) + " -> bloque" + fmt.Sprint(apt1) + ":titulo;\n"
							}
							if i == 0 {
								for j := 2; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
										if !eres {
											return "", false
										}
										respuesta += resp
										conexiones += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
									}
								}
							} else {
								for j := 0; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
										if !eres {
											return "", false
										}
										respuesta += resp
										conexiones += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
									}
								}
							}
						}
					}
					continue
				} else if i == 14 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					apuntadores := bloque_apuntador1.B_pointers
					data += "node [shape=plaintext]; \n"
					data += "\tbloque" + fmt.Sprint(apt) + " [label=<\n"
					data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
					data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt) + "</B></FONT></TD></TR>\n"
					//
					for j := range apuntadores {
						apt1 := apuntadores[j]
						data += "\t\t\t<TR>"
						data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt1) + "</FONT></TD>\n"
						data += "\t\t\t</TR>"
					}
					//
					data += "\t\t</table>\n"
					data += "\t>];\n\n"
					data += "inodo" + fmt.Sprint(num_inodo) + ":apts" + fmt.Sprint(i) + " -> bloque" + fmt.Sprint(apt) + ":titulo;\n"

					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								//
								//
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								apuntadores2 := bloque_apuntador2.B_pointers
								data += "node [shape=plaintext]; \n"
								data += "\tbloque" + fmt.Sprint(apt1) + " [label=<\n"
								data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
								data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt1) + "</B></FONT></TD></TR>\n"
								//
								for j := range apuntadores2 {
									apt2 := apuntadores2[j]
									data += "\t\t\t<TR>"
									data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt2) + "</FONT></TD>\n"
									data += "\t\t\t</TR>"
								}
								//
								data += "\t\t</table>\n"
								data += "\t>];\n\n"
								data += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(a) + " -> bloque" + fmt.Sprint(apt1) + ":titulo;\n"

								lista_apt2 := bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
										if !ebc {
											return "", false
										}
										info_carpeta := bloque_carpeta.B_content
										contenido1 := info_carpeta[0]
										contenido2 := info_carpeta[1]
										contenido3 := info_carpeta[2]
										contenido4 := info_carpeta[3]
										temp := reflect.TypeOf(structures.Content{})
										if (temp == reflect.TypeOf(contenido1)) && (temp == reflect.TypeOf(contenido2)) && (temp == reflect.TypeOf(contenido3)) && (temp == reflect.TypeOf(contenido4)) {
											data += "node [shape=plaintext]; \n"
											data += "\tbloque" + fmt.Sprint(apt2) + " [label=<\n"
											data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
											data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt2) + "</B></FONT></TD></TR>\n"
											//
											for z := range info_carpeta {
												dat := info_carpeta[z]
												data += "\t\t\t<TR>"
												data += "\t\t\t\t<TD port=\"apte" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + returnstring(ToString(dat.B_name[:])) + "</B></FONT></TD>\n"
												data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(dat.B_inodo) + "</FONT></TD>\n"
												data += "\t\t\t</TR>"
											}
											//
											data += "\t\t</table>\n"
											data += "\t>];\n\n"
											data += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(b) + " -> bloque" + fmt.Sprint(apt2) + ":titulo;\n"
										}
										if i == 0 {
											for j := 2; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
													if !eres {
														return "", false
													}
													respuesta += resp
													conexiones += "bloque" + fmt.Sprint(apt2) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
												}
											}
										} else {
											for j := 0; j < 4; j++ {
												info := info_carpeta[j]
												// info_temp := reflect.TypeOf(structures.Content{})
												// if info_temp == reflect.TypeOf(info)
												if info.B_inodo != -1 {
													resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
													if !eres {
														return "", false
													}
													respuesta += resp
													conexiones += "bloque" + fmt.Sprint(apt2) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
												}
											}
										}
									}
								}
								continue
								//
								//
							}
							bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
							if !ebc {
								return "", false
							}
							info_carpeta := bloque_carpeta.B_content
							contenido1 := info_carpeta[0]
							contenido2 := info_carpeta[1]
							contenido3 := info_carpeta[2]
							contenido4 := info_carpeta[3]
							temp := reflect.TypeOf(structures.Content{})
							if (temp == reflect.TypeOf(contenido1)) && (temp == reflect.TypeOf(contenido2)) && (temp == reflect.TypeOf(contenido3)) && (temp == reflect.TypeOf(contenido4)) {
								data += "node [shape=plaintext]; \n"
								data += "\tbloque" + fmt.Sprint(apt1) + " [label=<\n"
								data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
								data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt1) + "</B></FONT></TD></TR>\n"
								//
								for z := range info_carpeta {
									dat := info_carpeta[z]
									data += "\t\t\t<TR>"
									data += "\t\t\t\t<TD port=\"apte" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + returnstring(ToString(dat.B_name[:])) + "</B></FONT></TD>\n"
									data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(dat.B_inodo) + "</FONT></TD>\n"
									data += "\t\t\t</TR>"
								}
								//
								data += "\t\t</table>\n"
								data += "\t>];\n\n"
								data += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(a) + " -> bloque" + fmt.Sprint(apt1) + ":titulo;\n"
							}
							if i == 0 {
								for j := 2; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
										if !eres {
											return "", false
										}
										respuesta += resp
										conexiones += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
									}
								}
							} else {
								for j := 0; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
										if !eres {
											return "", false
										}
										respuesta += resp
										conexiones += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
									}
								}
							}
						}
					}
					// continue
				} else if i == 15 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					apuntadores := bloque_apuntador1.B_pointers
					data += "node [shape=plaintext]; \n"
					data += "\tbloque" + fmt.Sprint(apt) + " [label=<\n"
					data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
					data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt) + "</B></FONT></TD></TR>\n"
					//
					for j := range apuntadores {
						apt1 := apuntadores[j]
						data += "\t\t\t<TR>"
						data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt1) + "</FONT></TD>\n"
						data += "\t\t\t</TR>"
					}
					//
					data += "\t\t</table>\n"
					data += "\t>];\n\n"
					data += "inodo" + fmt.Sprint(num_inodo) + ":apts" + fmt.Sprint(i) + " -> bloque" + fmt.Sprint(apt) + ":titulo;\n"

					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								//
								//
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								apuntadores2 := bloque_apuntador2.B_pointers
								data += "node [shape=plaintext]; \n"
								data += "\tbloque" + fmt.Sprint(apt1) + " [label=<\n"
								data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
								data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt1) + "</B></FONT></TD></TR>\n"
								//
								for j := range apuntadores2 {
									apt2 := apuntadores2[j]
									data += "\t\t\t<TR>"
									data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt2) + "</FONT></TD>\n"
									data += "\t\t\t</TR>"
								}
								//
								data += "\t\t</table>\n"
								data += "\t>];\n\n"
								data += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(a) + " -> bloque" + fmt.Sprint(apt1) + ":titulo;\n"

								lista_apt2 := bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										if b == 0 {
											///
											///
											///
											bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt2)
											if !eba3 {
												return "", false
											}
											apuntadores3 := bloque_apuntador3.B_pointers
											data += "node [shape=plaintext]; \n"
											data += "\tbloque" + fmt.Sprint(apt2) + " [label=<\n"
											data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
											data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt2) + "</B></FONT></TD></TR>\n"
											//
											for j := range apuntadores3 {
												apt3 := apuntadores3[j]
												data += "\t\t\t<TR>"
												data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt3) + "</FONT></TD>\n"
												data += "\t\t\t</TR>"
											}
											//
											data += "\t\t</table>\n"
											data += "\t>];\n\n"
											data += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(b) + " -> bloque" + fmt.Sprint(apt2) + ":titulo;\n"

											lista_apt3 := bloque_apuntador3.B_pointers
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												if apt3 != -1 {
													bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt3)
													if !ebc {
														return "", false
													}
													info_carpeta := bloque_carpeta.B_content
													contenido1 := info_carpeta[0]
													contenido2 := info_carpeta[1]
													contenido3 := info_carpeta[2]
													contenido4 := info_carpeta[3]
													temp := reflect.TypeOf(structures.Content{})
													if (temp == reflect.TypeOf(contenido1)) && (temp == reflect.TypeOf(contenido2)) && (temp == reflect.TypeOf(contenido3)) && (temp == reflect.TypeOf(contenido4)) {
														data += "node [shape=plaintext]; \n"
														data += "\tbloque" + fmt.Sprint(apt3) + " [label=<\n"
														data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
														data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt3) + "</B></FONT></TD></TR>\n"
														//
														for z := range info_carpeta {
															dat := info_carpeta[z]
															data += "\t\t\t<TR>"
															data += "\t\t\t\t<TD port=\"apte" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + returnstring(ToString(dat.B_name[:])) + "</B></FONT></TD>\n"
															data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(dat.B_inodo) + "</FONT></TD>\n"
															data += "\t\t\t</TR>"
														}
														//
														data += "\t\t</table>\n"
														data += "\t>];\n\n"
														data += "bloque" + fmt.Sprint(apt2) + ":apts" + fmt.Sprint(c) + " -> bloque" + fmt.Sprint(apt3) + ":titulo;\n"
													}
													if i == 0 {
														for j := 2; j < 4; j++ {
															info := info_carpeta[j]
															if info.B_inodo != -1 {
																resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
																if !eres {
																	return "", false
																}
																respuesta += resp
																conexiones += "bloque" + fmt.Sprint(apt3) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
															}
														}
													} else {
														for j := 0; j < 4; j++ {
															info := info_carpeta[j]
															if info.B_inodo != -1 {
																resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
																if !eres {
																	return "", false
																}
																respuesta += resp
																conexiones += "bloque" + fmt.Sprint(apt3) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
															}
														}
													}
												}
											}
											continue
											///
											///
											///
										}
										bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
										if !ebc {
											return "", false
										}
										info_carpeta := bloque_carpeta.B_content
										contenido1 := info_carpeta[0]
										contenido2 := info_carpeta[1]
										contenido3 := info_carpeta[2]
										contenido4 := info_carpeta[3]
										temp := reflect.TypeOf(structures.Content{})
										if (temp == reflect.TypeOf(contenido1)) && (temp == reflect.TypeOf(contenido2)) && (temp == reflect.TypeOf(contenido3)) && (temp == reflect.TypeOf(contenido4)) {
											data += "node [shape=plaintext]; \n"
											data += "\tbloque" + fmt.Sprint(apt2) + " [label=<\n"
											data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
											data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt2) + "</B></FONT></TD></TR>\n"
											//
											for z := range info_carpeta {
												dat := info_carpeta[z]
												data += "\t\t\t<TR>"
												data += "\t\t\t\t<TD port=\"apte" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + returnstring(ToString(dat.B_name[:])) + "</B></FONT></TD>\n"
												data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(dat.B_inodo) + "</FONT></TD>\n"
												data += "\t\t\t</TR>"
											}
											//
											data += "\t\t</table>\n"
											data += "\t>];\n\n"
											data += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(b) + " -> bloque" + fmt.Sprint(apt2) + ":titulo;\n"
										}
										if i == 0 {
											for j := 2; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
													if !eres {
														return "", false
													}
													respuesta += resp
													conexiones += "bloque" + fmt.Sprint(apt2) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
												}
											}
										} else {
											for j := 0; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
													if !eres {
														return "", false
													}
													respuesta += resp
													conexiones += "bloque" + fmt.Sprint(apt2) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
												}
											}
										}
									}
								}
								continue
								//
								//
							}
							bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
							if !ebc {
								return "", false
							}
							info_carpeta := bloque_carpeta.B_content
							contenido1 := info_carpeta[0]
							contenido2 := info_carpeta[1]
							contenido3 := info_carpeta[2]
							contenido4 := info_carpeta[3]
							temp := reflect.TypeOf(structures.Content{})
							if (temp == reflect.TypeOf(contenido1)) && (temp == reflect.TypeOf(contenido2)) && (temp == reflect.TypeOf(contenido3)) && (temp == reflect.TypeOf(contenido4)) {
								data += "node [shape=plaintext]; \n"
								data += "\tbloque" + fmt.Sprint(apt1) + " [label=<\n"
								data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
								data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt1) + "</B></FONT></TD></TR>\n"
								//
								for z := range info_carpeta {
									dat := info_carpeta[z]
									data += "\t\t\t<TR>"
									data += "\t\t\t\t<TD port=\"apte" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + returnstring(ToString(dat.B_name[:])) + "</B></FONT></TD>\n"
									data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(dat.B_inodo) + "</FONT></TD>\n"
									data += "\t\t\t</TR>"
								}
								//
								data += "\t\t</table>\n"
								data += "\t>];\n\n"
								data += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(a) + " -> bloque" + fmt.Sprint(apt1) + ":titulo;\n"
							}
							if i == 0 {
								for j := 2; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
										if !eres {
											return "", false
										}
										respuesta += resp
										conexiones += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
									}
								}
							} else {
								for j := 0; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
										if !eres {
											return "", false
										}
										respuesta += resp
										conexiones += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
									}
								}
							}
						}
					}
					// continue
				} else {
					bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt)
					if !ebc {
						return "", false
					}
					info_carpeta := bloque_carpeta.B_content
					contenido1 := info_carpeta[0]
					contenido2 := info_carpeta[1]
					contenido3 := info_carpeta[2]
					contenido4 := info_carpeta[3]
					temp := reflect.TypeOf(structures.Content{})
					if (temp == reflect.TypeOf(contenido1)) && (temp == reflect.TypeOf(contenido2)) && (temp == reflect.TypeOf(contenido3)) && (temp == reflect.TypeOf(contenido4)) {
						data += "node [shape=plaintext]; \n"
						data += "\tbloque" + fmt.Sprint(apt) + " [label=<\n"
						data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
						data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt) + "</B></FONT></TD></TR>\n"
						//
						for z := range info_carpeta {
							dat := info_carpeta[z]
							data += "\t\t\t<TR>"
							data += "\t\t\t\t<TD port=\"apte" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + returnstring(ToString(dat.B_name[:])) + "</B></FONT></TD>\n"
							data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(z+1) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(dat.B_inodo) + "</FONT></TD>\n"
							data += "\t\t\t</TR>"
						}
						//
						data += "\t\t</table>\n"
						data += "\t>];\n\n"
						data += "inodo" + fmt.Sprint(num_inodo) + ":apts" + fmt.Sprint(i) + " -> bloque" + fmt.Sprint(apt) + ":titulo;\n"
					}
					if i == 0 {
						for j := 2; j < 4; j++ {
							info := info_carpeta[j]
							if info.B_inodo != -1 {
								resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
								if !eres {
									// fmt.Println("ASDFADFASDFADSFAFASDFs")
									return "", false
								}
								respuesta += resp
								conexiones += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
							}
						}
					} else {
						for j := 0; j < 4; j++ {
							info := info_carpeta[j]
							if info.B_inodo != -1 {
								resp, eres := Obtener_Codigo_Tree(comando, path, superbloque, info.B_inodo, id_log)
								if !eres {
									return "", false
								}
								respuesta += resp
								conexiones += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(j+1) + " -> " + "inodo" + fmt.Sprint(info.B_inodo) + ":titulo;\n"
							}
						}
					}
				}
			}
		}
		data += respuesta
		data += conexiones
		return data, true
	} else { //negacion
		//creacion inodo archivo
		for i := range inodo.I_block {
			apt := inodo.I_block[i]
			if apt != -1 {
				if i == 13 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					apuntadores := bloque_apuntador1.B_pointers
					data += "node [shape=plaintext]; \n"
					data += "\tbloque" + fmt.Sprint(apt) + " [label=<\n"
					data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
					data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt) + "</B></FONT></TD></TR>\n"
					//
					for j := range apuntadores {
						apt1 := apuntadores[j]
						data += "\t\t\t<TR>"
						data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt1) + "</FONT></TD>\n"
						data += "\t\t\t</TR>"
					}
					//
					data += "\t\t</table>\n"
					data += "\t>];\n\n"
					data += "inodo" + fmt.Sprint(num_inodo) + ":apts" + fmt.Sprint(i) + " -> bloque" + fmt.Sprint(apt) + ":titulo;\n"

					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							contenido := strings.ReplaceAll(returnstring(ToString(bloque_archivo.B_content[:])), "\n", "<BR/>")
							data += "node [shape=plaintext]; \n"
							data += "\tbloque" + fmt.Sprint(apt1) + " [label=<\n"
							data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
							data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt1) + "</B></FONT></TD></TR>\n"
							data += "\t\t\t<TR>"
							data += "\t\t\t\t<TD port=\"apts1\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + contenido + "</B></FONT></TD>\n"
							data += "\t\t\t</TR>"
							data += "\t\t</table>\n"
							data += "\t>];\n\n"
							data += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(a) + " -> bloque" + fmt.Sprint(apt1) + ":titulo;\n"
						}
					}
					continue
				} else if i == 14 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					apuntadores := bloque_apuntador1.B_pointers
					data += "node [shape=plaintext]; \n"
					data += "\tbloque" + fmt.Sprint(apt) + " [label=<\n"
					data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
					data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt) + "</B></FONT></TD></TR>\n"
					//
					for j := range apuntadores {
						apt1 := apuntadores[j]
						data += "\t\t\t<TR>"
						data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt1) + "</FONT></TD>\n"
						data += "\t\t\t</TR>"
					}
					//
					data += "\t\t</table>\n"
					data += "\t>];\n\n"
					data += "inodo" + fmt.Sprint(num_inodo) + ":apts" + fmt.Sprint(i) + " -> bloque" + fmt.Sprint(apt) + ":titulo;\n"

					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								//
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								apuntadores2 := bloque_apuntador2.B_pointers
								data += "node [shape=plaintext]; \n"
								data += "\tbloque" + fmt.Sprint(apt1) + " [label=<\n"
								data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
								data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt1) + "</B></FONT></TD></TR>\n"
								//
								for j := range apuntadores2 {
									apt2 := apuntadores2[j]
									data += "\t\t\t<TR>"
									data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt2) + "</FONT></TD>\n"
									data += "\t\t\t</TR>"
								}
								//
								data += "\t\t</table>\n"
								data += "\t>];\n\n"
								data += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(a) + " -> bloque" + fmt.Sprint(apt1) + ":titulo;\n"

								lista_apt2 := bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										contenido := strings.ReplaceAll(returnstring(ToString(bloque_archivo.B_content[:])), "\n", "<BR/>")
										data += "node [shape=plaintext]; \n"
										data += "\tbloque" + fmt.Sprint(apt2) + " [label=<\n"
										data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
										data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt2) + "</B></FONT></TD></TR>\n"
										data += "\t\t\t<TR>"
										data += "\t\t\t\t<TD port=\"apts1\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + contenido + "</B></FONT></TD>\n"
										data += "\t\t\t</TR>"
										data += "\t\t</table>\n"
										data += "\t>];\n\n"
										data += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(b) + " -> bloque" + fmt.Sprint(apt2) + ":titulo;\n"
									}
								}
								continue
								//
							}
							bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							contenido := strings.ReplaceAll(returnstring(ToString(bloque_archivo.B_content[:])), "\n", "<BR/>")
							data += "node [shape=plaintext]; \n"
							data += "\tbloque" + fmt.Sprint(apt1) + " [label=<\n"
							data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
							data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt1) + "</B></FONT></TD></TR>\n"
							data += "\t\t\t<TR>"
							data += "\t\t\t\t<TD port=\"apts1\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + contenido + "</B></FONT></TD>\n"
							data += "\t\t\t</TR>"
							data += "\t\t</table>\n"
							data += "\t>];\n\n"
							data += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(a) + " -> bloque" + fmt.Sprint(apt1) + ":titulo;\n"
						}
					}

				} else if i == 15 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					apuntadores := bloque_apuntador1.B_pointers
					data += "node [shape=plaintext]; \n"
					data += "\tbloque" + fmt.Sprint(apt) + " [label=<\n"
					data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
					data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt) + "</B></FONT></TD></TR>\n"
					//
					for j := range apuntadores {
						apt1 := apuntadores[j]
						data += "\t\t\t<TR>"
						data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt1) + "</FONT></TD>\n"
						data += "\t\t\t</TR>"
					}
					//
					data += "\t\t</table>\n"
					data += "\t>];\n\n"
					data += "inodo" + fmt.Sprint(num_inodo) + ":apts" + fmt.Sprint(i) + " -> bloque" + fmt.Sprint(apt) + ":titulo;\n"

					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								//
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								apuntadores2 := bloque_apuntador2.B_pointers
								data += "node [shape=plaintext]; \n"
								data += "\tbloque" + fmt.Sprint(apt1) + " [label=<\n"
								data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
								data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt1) + "</B></FONT></TD></TR>\n"
								//
								for j := range apuntadores2 {
									apt2 := apuntadores2[j]
									data += "\t\t\t<TR>"
									data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt2) + "</FONT></TD>\n"
									data += "\t\t\t</TR>"
								}
								//
								data += "\t\t</table>\n"
								data += "\t>];\n\n"
								data += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(a) + " -> bloque" + fmt.Sprint(apt1) + ":titulo;\n"

								lista_apt2 := bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										if b == 0 {
											//
											//
											bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt2)
											if !eba3 {
												return "", false
											}
											apuntadores2 := bloque_apuntador3.B_pointers
											data += "node [shape=plaintext]; \n"
											data += "\tbloque" + fmt.Sprint(apt2) + " [label=<\n"
											data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
											data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\"><B>Bloque Apuntador" + fmt.Sprint(apt2) + "</B></FONT></TD></TR>\n"
											//
											for j := range apuntadores2 {
												apt3 := apuntadores2[j]
												data += "\t\t\t<TR>"
												data += "\t\t\t\t<TD port=\"apts" + fmt.Sprint(j) + "\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D8B3EE\"><FONT COLOR=\"BLACK\">" + fmt.Sprint(apt3) + "</FONT></TD>\n"
												data += "\t\t\t</TR>"
											}
											//
											data += "\t\t</table>\n"
											data += "\t>];\n\n"
											data += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(b) + " -> bloque" + fmt.Sprint(apt2) + ":titulo;\n"

											lista_apt3 := bloque_apuntador3.B_pointers
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												if apt3 != -1 {
													bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt3)
													if !eba {
														return "", false
													}
													contenido := strings.ReplaceAll(returnstring(ToString(bloque_archivo.B_content[:])), "\n", "<BR/>")
													data += "node [shape=plaintext]; \n"
													data += "\tbloque" + fmt.Sprint(apt3) + " [label=<\n"
													data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
													data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt3) + "</B></FONT></TD></TR>\n"
													data += "\t\t\t<TR>"
													data += "\t\t\t\t<TD port=\"apts1\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + contenido + "</B></FONT></TD>\n"
													data += "\t\t\t</TR>"
													data += "\t\t</table>\n"
													data += "\t>];\n\n"
													data += "bloque" + fmt.Sprint(apt2) + ":apts" + fmt.Sprint(c) + " -> bloque" + fmt.Sprint(apt3) + ":titulo;\n"
												}
											}
											continue
											//
											//
										}
										bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										contenido := strings.ReplaceAll(returnstring(ToString(bloque_archivo.B_content[:])), "\n", "<BR/>")
										data += "node [shape=plaintext]; \n"
										data += "\tbloque" + fmt.Sprint(apt2) + " [label=<\n"
										data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
										data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt2) + "</B></FONT></TD></TR>\n"
										data += "\t\t\t<TR>"
										data += "\t\t\t\t<TD port=\"apts1\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + contenido + "</B></FONT></TD>\n"
										data += "\t\t\t</TR>"
										data += "\t\t</table>\n"
										data += "\t>];\n\n"
										data += "bloque" + fmt.Sprint(apt1) + ":apts" + fmt.Sprint(b) + " -> bloque" + fmt.Sprint(apt2) + ":titulo;\n"
									}
								}
								continue
								//
							}
							bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							contenido := strings.ReplaceAll(returnstring(ToString(bloque_archivo.B_content[:])), "\n", "<BR/>")
							data += "node [shape=plaintext]; \n"
							data += "\tbloque" + fmt.Sprint(apt1) + " [label=<\n"
							data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
							data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt1) + "</B></FONT></TD></TR>\n"
							data += "\t\t\t<TR>"
							data += "\t\t\t\t<TD port=\"apts1\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + contenido + "</B></FONT></TD>\n"
							data += "\t\t\t</TR>"
							data += "\t\t</table>\n"
							data += "\t>];\n\n"
							data += "bloque" + fmt.Sprint(apt) + ":apts" + fmt.Sprint(a) + " -> bloque" + fmt.Sprint(apt1) + ":titulo;\n"
						}
					}

				} else {
					bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt)
					if !eba {
						return "", false
					}
					contenido := strings.ReplaceAll(returnstring(ToString(bloque_archivo.B_content[:])), "\n", "<BR/>")
					data += "node [shape=plaintext]; \n"
					data += "\tbloque" + fmt.Sprint(apt) + " [label=<\n"
					data += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
					data += "\t\t\t<TR><TD port=\"titulo\" ALIGN=\"CENTER\"  BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>Inodo" + fmt.Sprint(apt) + "</B></FONT></TD></TR>\n"
					data += "\t\t\t<TR>"
					data += "\t\t\t\t<TD port=\"apts1\" ALIGN=\"CENTER\" WIDTH=\"200\" BGCOLOR=\"#D68671\"><FONT COLOR=\"BLACK\"><B>" + contenido + "</B></FONT></TD>\n"
					data += "\t\t\t</TR>"
					data += "\t\t</table>\n"
					data += "\t>];\n\n"
					data += "inodo" + fmt.Sprint(num_inodo) + ":apts" + fmt.Sprint(i) + " -> bloque" + fmt.Sprint(apt) + ":titulo;\n"
				}
			}
		}
		return data, true
	}

	//
	// return "", false
}
