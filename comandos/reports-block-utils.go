package comandos

import (
	"MIA_P1_201906051/structures"
	"fmt"
	"reflect"
	"strings"
)

func Obtener_Codigo_Bloques(comando string, path string, superbloque structures.SuperBlock, num_inodo int32, espaciado string) (string, bool) {
	respuesta := ""
	inodo, ei := Obtener_Inodo("REP", path, superbloque.S_inode_start, num_inodo)
	if !ei {
		return "", false
	}
	dato := "\n"
	if inodo.I_type == 0 {
		for i := range inodo.I_block {
			apt := inodo.I_block[i]
			if apt != -1 {
				if i == 13 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers

					dato += "\tbloque" + fmt.Sprint(apt) + "[label=<\n\n"
					dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
					dato += "\t\t\t<TR>\n"
					dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt) + "</B></TD>\n"
					dato += "\t\t\t</TR>\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(a) + "</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt1) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"
					}
					dato += "\t\t</TABLE>\n"
					dato += "\t>];\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							info_carpeta := bloque_carpeta.B_content
							if i == 0 {
								for j := 2; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
										if !efaul {
											return "", false
										}
										respuesta += espaciado + res
									}
								}
							} else {
								for j := range info_carpeta {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
										if !efaul {
											return "", false
										}
										respuesta += espaciado + res
									}
								}
							}
							contenido1 := info_carpeta[0]
							contenido2 := info_carpeta[1]
							contenido3 := info_carpeta[2]
							contenido4 := info_carpeta[3]
							tipo := reflect.TypeOf(structures.Content{})
							if (tipo == reflect.TypeOf(contenido1)) && (tipo == reflect.TypeOf(contenido2)) && (tipo == reflect.TypeOf(contenido3)) && (tipo == reflect.TypeOf(contenido4)) {
								dato += "\tbloque" + fmt.Sprint(apt1) + "[label=<\n\n"
								dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt1) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"

								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_name</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_inodo</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido1.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido1.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido2.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido2.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido3.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido3.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido4.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido4.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"

								dato += "\t\t</TABLE>\n"
								dato += "\t>];\n"
							}

						}
					}
					continue
				} else if i == 14 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers

					dato += "\tbloque" + fmt.Sprint(apt) + "[label=<\n\n"
					dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
					dato += "\t\t\t<TR>\n"
					dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt) + "</B></TD>\n"
					dato += "\t\t\t</TR>\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(a) + "</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt1) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"
					}
					dato += "\t\t</TABLE>\n"
					dato += "\t>];\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								lista_apt2 := bloque_apuntador2.B_pointers
								dato += "\tbloque" + fmt.Sprint(apt1) + "[label=<\n\n"
								dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt1) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									dato += "\t\t\t<TR>\n"
									dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(b) + "</B></TD>\n"
									dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt2) + "</B></TD>\n"
									dato += "\t\t\t</TR>\n"
								}
								dato += "\t\t</TABLE>\n"
								dato += "\t>];\n"
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										//
										bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										info_carpeta := bloque_carpeta.B_content
										if i == 0 {
											for j := 2; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
													if !efaul {
														return "", false
													}
													respuesta += espaciado + res
												}
											}
										} else {
											for j := range info_carpeta {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
													if !efaul {
														return "", false
													}
													respuesta += espaciado + res
												}
											}
										}
										contenido1 := info_carpeta[0]
										contenido2 := info_carpeta[1]
										contenido3 := info_carpeta[2]
										contenido4 := info_carpeta[3]
										tipo := reflect.TypeOf(structures.Content{})
										if (tipo == reflect.TypeOf(contenido1)) && (tipo == reflect.TypeOf(contenido2)) && (tipo == reflect.TypeOf(contenido3)) && (tipo == reflect.TypeOf(contenido4)) {
											dato += "\tbloque" + fmt.Sprint(apt2) + "[label=<\n\n"
											dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt2) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"

											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_name</B></TD>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_inodo</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido1.B_name[:])) + "</B></TD>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido1.B_inodo) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido2.B_name[:])) + "</B></TD>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido2.B_inodo) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido3.B_name[:])) + "</B></TD>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido3.B_inodo) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido4.B_name[:])) + "</B></TD>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido4.B_inodo) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											dato += "\t\t</TABLE>\n"
											dato += "\t>];\n"
										}
										//
									}
								}
								continue
							}
							bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							info_carpeta := bloque_carpeta.B_content
							if i == 0 {
								for j := 2; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
										if !efaul {
											return "", false
										}
										respuesta += espaciado + res
									}
								}
							} else {
								for j := range info_carpeta {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
										if !efaul {
											return "", false
										}
										respuesta += espaciado + res
									}
								}
							}
							contenido1 := info_carpeta[0]
							contenido2 := info_carpeta[1]
							contenido3 := info_carpeta[2]
							contenido4 := info_carpeta[3]
							tipo := reflect.TypeOf(structures.Content{})
							if (tipo == reflect.TypeOf(contenido1)) && (tipo == reflect.TypeOf(contenido2)) && (tipo == reflect.TypeOf(contenido3)) && (tipo == reflect.TypeOf(contenido4)) {
								dato += "\tbloque" + fmt.Sprint(apt) + "[label=<\n\n"
								dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"

								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_name</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_inodo</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido1.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido1.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido2.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido2.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido3.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido3.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido4.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido4.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"

								dato += "\t\t</TABLE>\n"
								dato += "\t>];\n"
							}

						}
					}
					continue
				} else if i == 15 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers

					dato += "\tbloque" + fmt.Sprint(apt) + "[label=<\n\n"
					dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
					dato += "\t\t\t<TR>\n"
					dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt) + "</B></TD>\n"
					dato += "\t\t\t</TR>\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(a) + "</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt1) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"
					}
					dato += "\t\t</TABLE>\n"
					dato += "\t>];\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								lista_apt2 := bloque_apuntador2.B_pointers
								dato += "\tbloque" + fmt.Sprint(apt1) + "[label=<\n\n"
								dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt1) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									dato += "\t\t\t<TR>\n"
									dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(b) + "</B></TD>\n"
									dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt2) + "</B></TD>\n"
									dato += "\t\t\t</TR>\n"
								}
								dato += "\t\t</TABLE>\n"
								dato += "\t>];\n"
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
											lista_apt3 := bloque_apuntador3.B_pointers
											dato += "\tbloque" + fmt.Sprint(apt2) + "[label=<\n\n"
											dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt2) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												dato += "\t\t\t<TR>\n"
												dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(c) + "</B></TD>\n"
												dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt3) + "</B></TD>\n"
												dato += "\t\t\t</TR>\n"
											}
											dato += "\t\t</TABLE>\n"
											dato += "\t>];\n"
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												if apt3 != -1 {
													//
													bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt3)
													if !eba {
														return "", false
													}
													info_carpeta := bloque_carpeta.B_content
													if i == 0 {
														for j := 2; j < 4; j++ {
															info := info_carpeta[j]
															if info.B_inodo != -1 {
																res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
																if !efaul {
																	return "", false
																}
																respuesta += espaciado + res
															}
														}
													} else {
														for j := range info_carpeta {
															info := info_carpeta[j]
															if info.B_inodo != -1 {
																res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
																if !efaul {
																	return "", false
																}
																respuesta += espaciado + res
															}
														}
													}
													contenido1 := info_carpeta[0]
													contenido2 := info_carpeta[1]
													contenido3 := info_carpeta[2]
													contenido4 := info_carpeta[3]
													tipo := reflect.TypeOf(structures.Content{})
													if (tipo == reflect.TypeOf(contenido1)) && (tipo == reflect.TypeOf(contenido2)) && (tipo == reflect.TypeOf(contenido3)) && (tipo == reflect.TypeOf(contenido4)) {
														dato += "\tbloque" + fmt.Sprint(apt3) + "[label=<\n\n"
														dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
														dato += "\t\t\t<TR>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt3) + "</B></TD>\n"
														dato += "\t\t\t</TR>\n"

														dato += "\t\t\t<TR>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_name</B></TD>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_inodo</B></TD>\n"
														dato += "\t\t\t</TR>\n"
														dato += "\t\t\t<TR>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido1.B_name[:])) + "</B></TD>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido1.B_inodo) + "</B></TD>\n"
														dato += "\t\t\t</TR>\n"
														dato += "\t\t\t<TR>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido2.B_name[:])) + "</B></TD>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido2.B_inodo) + "</B></TD>\n"
														dato += "\t\t\t</TR>\n"
														dato += "\t\t\t<TR>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido3.B_name[:])) + "</B></TD>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido3.B_inodo) + "</B></TD>\n"
														dato += "\t\t\t</TR>\n"
														dato += "\t\t\t<TR>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido4.B_name[:])) + "</B></TD>\n"
														dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido4.B_inodo) + "</B></TD>\n"
														dato += "\t\t\t</TR>\n"
														dato += "\t\t</TABLE>\n"
														dato += "\t>];\n"
													}
													//
												}
											}
											continue

											//
											//
										}
										//
										bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										info_carpeta := bloque_carpeta.B_content
										if i == 0 {
											for j := 2; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
													if !efaul {
														return "", false
													}
													respuesta += espaciado + res
												}
											}
										} else {
											for j := range info_carpeta {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
													if !efaul {
														return "", false
													}
													respuesta += espaciado + res
												}
											}
										}
										contenido1 := info_carpeta[0]
										contenido2 := info_carpeta[1]
										contenido3 := info_carpeta[2]
										contenido4 := info_carpeta[3]
										tipo := reflect.TypeOf(structures.Content{})
										if (tipo == reflect.TypeOf(contenido1)) && (tipo == reflect.TypeOf(contenido2)) && (tipo == reflect.TypeOf(contenido3)) && (tipo == reflect.TypeOf(contenido4)) {
											dato += "\tbloque" + fmt.Sprint(apt2) + "[label=<\n\n"
											dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt2) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"

											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_name</B></TD>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_inodo</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido1.B_name[:])) + "</B></TD>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido1.B_inodo) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido2.B_name[:])) + "</B></TD>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido2.B_inodo) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido3.B_name[:])) + "</B></TD>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido3.B_inodo) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido4.B_name[:])) + "</B></TD>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido4.B_inodo) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											dato += "\t\t</TABLE>\n"
											dato += "\t>];\n"
										}
										//
									}
								}
								continue
							}
							bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							info_carpeta := bloque_carpeta.B_content
							if i == 0 {
								for j := 2; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
										if !efaul {
											return "", false
										}
										respuesta += espaciado + res
									}
								}
							} else {
								for j := range info_carpeta {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
										if !efaul {
											return "", false
										}
										respuesta += espaciado + res
									}
								}
							}
							contenido1 := info_carpeta[0]
							contenido2 := info_carpeta[1]
							contenido3 := info_carpeta[2]
							contenido4 := info_carpeta[3]
							tipo := reflect.TypeOf(structures.Content{})
							if (tipo == reflect.TypeOf(contenido1)) && (tipo == reflect.TypeOf(contenido2)) && (tipo == reflect.TypeOf(contenido3)) && (tipo == reflect.TypeOf(contenido4)) {
								dato += "\tbloque" + fmt.Sprint(apt) + "[label=<\n\n"
								dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"

								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_name</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_inodo</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido1.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido1.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido2.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido2.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido3.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido3.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido4.B_name[:])) + "</B></TD>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido4.B_inodo) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"

								dato += "\t\t</TABLE>\n"
								dato += "\t>];\n"
							}

						}
					}
					continue
				} else {
					bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt)
					if !eba {
						return "", false
					}
					info_carpeta := bloque_carpeta.B_content
					if i == 0 {
						for j := 2; j < 4; j++ {
							info := info_carpeta[j]
							if info.B_inodo != -1 {
								res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
								if !efaul {
									return "", false
								}
								respuesta += espaciado + res
							}
						}
					} else {
						for j := range info_carpeta {
							info := info_carpeta[j]
							if info.B_inodo != -1 {
								res, efaul := Obtener_Codigo_Bloques(comando, path, superbloque, info.B_inodo, espaciado+"   ")
								if !efaul {
									return "", false
								}
								respuesta += espaciado + res
							}
						}
					}
					contenido1 := info_carpeta[0]
					contenido2 := info_carpeta[1]
					contenido3 := info_carpeta[2]
					contenido4 := info_carpeta[3]
					tipo := reflect.TypeOf(structures.Content{})
					if (tipo == reflect.TypeOf(contenido1)) && (tipo == reflect.TypeOf(contenido2)) && (tipo == reflect.TypeOf(contenido3)) && (tipo == reflect.TypeOf(contenido4)) {
						dato += "\tbloque" + fmt.Sprint(apt) + "[label=<\n\n"
						dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"

						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_name</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>b_inodo</B></TD>\n"
						dato += "\t\t\t</TR>\n"
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido1.B_name[:])) + "</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido1.B_inodo) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido2.B_name[:])) + "</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido2.B_inodo) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido3.B_name[:])) + "</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido3.B_inodo) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + returnstring(ToString(contenido4.B_name[:])) + "</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(contenido4.B_inodo) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"

						dato += "\t\t</TABLE>\n"
						dato += "\t>];\n"
					}
				}
			}
		}
		dato += respuesta
		return dato, true
	} else {
		//recorrerr cada bloque del inodo
		//
		for i := range inodo.I_block {
			apt := inodo.I_block[i]
			if apt != -1 {
				if i == 13 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					dato += "\tbloque" + fmt.Sprint(apt) + "[label=<\n\n"
					dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
					dato += "\t\t\t<TR>\n"
					dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt) + "</B></TD>\n"
					dato += "\t\t\t</TR>\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(a) + "</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt1) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"
					}
					dato += "\t\t</TABLE>\n"
					dato += "\t>];\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							bloque, eb := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt1)
							if !eb {
								return "", false
							}
							// a:= ToString(bloque.B_content[:])
							contenido := ToString(bloque.B_content[:])
							contenido = strings.ReplaceAll(contenido, "\n", "<BR/>")
							dato += "\tbloque" + fmt.Sprint(apt1) + "[label=<\n\n"
							dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
							dato += "\t\t\t<TR>\n"
							dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>Bloque Archivo " + fmt.Sprint(apt1) + "</B></TD>\n"
							dato += "\t\t\t</TR>\n"
							//
							dato += "\t\t\t<TR>\n"
							dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>" + fmt.Sprint(contenido) + "</B></TD>\n"
							dato += "\t\t\t</TR>\n"
							//
							dato += "\t\t</TABLE>\n"
							dato += "\t>];\n"
						}
					}
					continue
				} else if i == 14 { //para 14
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					dato += "\tbloque" + fmt.Sprint(apt) + "[label=<\n\n"
					dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
					dato += "\t\t\t<TR>\n"
					dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt) + "</B></TD>\n"
					dato += "\t\t\t</TR>\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(a) + "</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt1) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"
					}
					dato += "\t\t</TABLE>\n"
					dato += "\t>];\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								//
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								lista_apt2 := bloque_apuntador2.B_pointers
								dato += "\tbloque" + fmt.Sprint(apt1) + "[label=<\n\n"
								dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt1) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									dato += "\t\t\t<TR>\n"
									dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(b) + "</B></TD>\n"
									dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt2) + "</B></TD>\n"
									dato += "\t\t\t</TR>\n"
								}
								dato += "\t\t</TABLE>\n"
								dato += "\t>];\n"
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										bloque, eb := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt2)
										if !eb {
											return "", false
										}
										contenido := ToString(bloque.B_content[:])
										contenido = strings.ReplaceAll(contenido, "\n", "<BR/>")
										dato += "\tbloque" + fmt.Sprint(apt2) + "[label=<\n\n"
										dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
										dato += "\t\t\t<TR>\n"
										dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>Bloque Archivo " + fmt.Sprint(apt2) + "</B></TD>\n"
										dato += "\t\t\t</TR>\n"
										//
										dato += "\t\t\t<TR>\n"
										dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>" + fmt.Sprint(contenido) + "</B></TD>\n"
										dato += "\t\t\t</TR>\n"
										//
										dato += "\t\t</TABLE>\n"
										dato += "\t>];\n"
									}
								}
								continue
								//
							}
							bloque, eb := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt1)
							if !eb {
								return "", false
							}
							contenido := ToString(bloque.B_content[:])
							contenido = strings.ReplaceAll(contenido, "\n", "<BR/>")
							dato += "\tbloque" + fmt.Sprint(apt1) + "[label=<\n\n"
							dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
							dato += "\t\t\t<TR>\n"
							dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>Bloque Archivo " + fmt.Sprint(apt1) + "</B></TD>\n"
							dato += "\t\t\t</TR>\n"
							//
							dato += "\t\t\t<TR>\n"
							dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>" + fmt.Sprint(contenido) + "</B></TD>\n"
							dato += "\t\t\t</TR>\n"
							//
							dato += "\t\t</TABLE>\n"
							dato += "\t>];\n"
						}
					}
					continue
				} else if i == 15 { //para 15
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					dato += "\tbloque" + fmt.Sprint(apt) + "[label=<\n\n"
					dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
					dato += "\t\t\t<TR>\n"
					dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt) + "</B></TD>\n"
					dato += "\t\t\t</TR>\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						dato += "\t\t\t<TR>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(a) + "</B></TD>\n"
						dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt1) + "</B></TD>\n"
						dato += "\t\t\t</TR>\n"
					}
					dato += "\t\t</TABLE>\n"
					dato += "\t>];\n"
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								//
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								lista_apt2 := bloque_apuntador2.B_pointers
								dato += "\tbloque" + fmt.Sprint(apt1) + "[label=<\n\n"
								dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
								dato += "\t\t\t<TR>\n"
								dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt1) + "</B></TD>\n"
								dato += "\t\t\t</TR>\n"
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									dato += "\t\t\t<TR>\n"
									dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(b) + "</B></TD>\n"
									dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt2) + "</B></TD>\n"
									dato += "\t\t\t</TR>\n"
								}
								dato += "\t\t</TABLE>\n"
								dato += "\t>];\n"
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										if b == 0 {
											//
											//
											//

											bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt2)
											if !eba3 {
												return "", false
											}
											lista_apt3 := bloque_apuntador3.B_pointers
											dato += "\tbloque" + fmt.Sprint(apt2) + "[label=<\n\n"
											dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
											dato += "\t\t\t<TR>\n"
											dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\" WIDTH=\"150\"><B>Bloque Apuntador " + fmt.Sprint(apt2) + "</B></TD>\n"
											dato += "\t\t\t</TR>\n"
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												dato += "\t\t\t<TR>\n"
												dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>apt " + fmt.Sprint(a) + "</B></TD>\n"
												dato += "\t\t\t\t<TD ALIGN=\"CENTER\"><B>" + fmt.Sprint(apt3) + "</B></TD>\n"
												dato += "\t\t\t</TR>\n"
											}
											dato += "\t\t</TABLE>\n"
											dato += "\t>];\n"
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												if apt3 != -1 {
													bloque, eb := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt3)
													if !eb {
														return "", false
													}
													contenido := ToString(bloque.B_content[:])
													contenido = strings.ReplaceAll(contenido, "\n", "<BR/>")
													dato += "\tbloque" + fmt.Sprint(apt3) + "[label=<\n\n"
													dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
													dato += "\t\t\t<TR>\n"
													dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>Bloque Archivo " + fmt.Sprint(apt3) + "</B></TD>\n"
													dato += "\t\t\t</TR>\n"
													//
													dato += "\t\t\t<TR>\n"
													dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>" + fmt.Sprint(contenido) + "</B></TD>\n"
													dato += "\t\t\t</TR>\n"
													//
													dato += "\t\t</TABLE>\n"
													dato += "\t>];\n"
												}
											}
											continue

											//
											//
											//
										}
										bloque, eb := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt2)
										if !eb {
											return "", false
										}
										contenido := ToString(bloque.B_content[:])
										contenido = strings.ReplaceAll(contenido, "\n", "<BR/>")
										dato += "\tbloque" + fmt.Sprint(apt2) + "[label=<\n\n"
										dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
										dato += "\t\t\t<TR>\n"
										dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>Bloque Archivo " + fmt.Sprint(apt2) + "</B></TD>\n"
										dato += "\t\t\t</TR>\n"
										//
										dato += "\t\t\t<TR>\n"
										dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>" + fmt.Sprint(contenido) + "</B></TD>\n"
										dato += "\t\t\t</TR>\n"
										//
										dato += "\t\t</TABLE>\n"
										dato += "\t>];\n"
									}
								}
								continue
								//
							}
							bloque, eb := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt1)
							if !eb {
								return "", false
							}
							contenido := ToString(bloque.B_content[:])
							contenido = strings.ReplaceAll(contenido, "\n", "<BR/>")
							dato += "\tbloque" + fmt.Sprint(apt1) + "[label=<\n\n"
							dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
							dato += "\t\t\t<TR>\n"
							dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>Bloque Archivo " + fmt.Sprint(apt1) + "</B></TD>\n"
							dato += "\t\t\t</TR>\n"
							//
							dato += "\t\t\t<TR>\n"
							dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>" + fmt.Sprint(contenido) + "</B></TD>\n"
							dato += "\t\t\t</TR>\n"
							//
							dato += "\t\t</TABLE>\n"
							dato += "\t>];\n"
						}
					}
					continue
				} else { //ninguno de los anteriores
					bloque, eb := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt)
					if !eb {
						return "", false
					}
					contenido := ToString(bloque.B_content[:])
					contenido = strings.ReplaceAll(contenido, "\n", "<BR/>")
					dato += "\tbloque" + fmt.Sprint(apt) + "[label=<\n\n"
					dato += "\t\t<TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\">\n"
					dato += "\t\t\t<TR>\n"
					dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>Bloque Archivo " + fmt.Sprint(apt) + "</B></TD>\n"
					dato += "\t\t\t</TR>\n"
					//
					dato += "\t\t\t<TR>\n"
					dato += "\t\t\t\t<TD ALIGN=\"CENTER\" COLSPAN=\"2\"><B>" + fmt.Sprint(contenido) + "</B></TD>\n"
					dato += "\t\t\t</TR>\n"
					//
					dato += "\t\t</TABLE>\n"
					dato += "\t>];\n"
				}
			}
		}
		return dato, true
		//
	}
	// return dato, true
	// return "", false
}
