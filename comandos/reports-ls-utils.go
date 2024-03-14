package comandos

import (
	"MIA_P1_201906051/structures"
	"fmt"
	"strings"
)

func Obtener_Reporte_LS(comando string, path string, superbloque structures.SuperBlock, num_inodo int32, id_log string) (string, bool) {
	respuesta := ""
	inodo, eino := Obtener_Inodo(comando, path, superbloque.S_inode_start, num_inodo)
	if !eino {
		return "", false
	}
	dato := ""
	if inodo.I_type == 0 {
		for i := range inodo.I_block {
			///continuar luego
			apt := inodo.I_block[i]
			if apt != -1 {
				if i == 13 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
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
										res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !eres {
											return "", false
										}
										respuesta += res
									}
								}
							} else {
								for j := 0; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !eres {
											return "", false
										}
										respuesta += res
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
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								lista_apt2 := bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										info_carpeta := bloque_carpeta.B_content
										if i == 0 {
											for j := 2; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
													if !eres {
														return "", false
													}
													respuesta += res
												}
											}
										} else {
											for j := 0; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
													if !eres {
														return "", false
													}
													respuesta += res
												}
											}
										}
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
										res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !eres {
											return "", false
										}
										respuesta += res
									}
								}
							} else {
								for j := 0; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !eres {
											return "", false
										}
										respuesta += res
									}
								}
							}
						}
					}
					continue
				} else if i == 15 { //caso en el que i == 15
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
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
											lista_apt3 := bloque_apuntador3.B_pointers
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												if apt3 != -1 {
													bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt3)
													if !eba {
														return "", false
													}
													info_carpeta := bloque_carpeta.B_content
													if i == 0 {
														for j := 2; j < 4; j++ {
															info := info_carpeta[j]
															if info.B_inodo != -1 {
																res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
																if !eres {
																	return "", false
																}
																respuesta += res
															}
														}
													} else {
														for j := 0; j < 4; j++ {
															info := info_carpeta[j]
															if info.B_inodo != -1 {
																res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
																if !eres {
																	return "", false
																}
																respuesta += res
															}
														}
													}
												}
											}
											continue
											//
											//
										}
										bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										info_carpeta := bloque_carpeta.B_content
										if i == 0 {
											for j := 2; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
													if !eres {
														return "", false
													}
													respuesta += res
												}
											}
										} else {
											for j := 0; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
													if !eres {
														return "", false
													}
													respuesta += res
												}
											}
										}
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
										res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !eres {
											return "", false
										}
										respuesta += res
									}
								}
							} else {
								for j := 0; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !eres {
											return "", false
										}
										respuesta += res
									}
								}
							}
						}
					}
					continue
				} else {
					//
					bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt)
					if !eba {
						return "", false
					}
					info_carpeta := bloque_carpeta.B_content
					if i == 0 {
						for j := 2; j < 4; j++ {
							info := info_carpeta[j]
							if info.B_inodo != -1 {
								res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
								if !eres {
									return "", false
								}
								respuesta += res
							}
						}
					} else {
						for j := 0; j < 4; j++ {
							info := info_carpeta[j]
							if info.B_inodo != -1 {
								res, eres := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
								if !eres {
									return "", false
								}
								respuesta += res
							}
						}
					}
					//
				}
			}
		}
		dato += respuesta
		return dato, true
	}
	//
	return "", false
}

func Obtener_Codigo_LS(comando string, path string, superbloque structures.SuperBlock, num_inodo int32, id_log string, nombre string) (string, bool) {
	respuesta := ""
	inodo, einodo := Obtener_Inodo(comando, path, superbloque.S_inode_start, num_inodo)
	if !einodo {
		return "", false
	}
	dato := ""
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
										res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !ores {
											return "", false
										}
										respuesta += res
									}
								}
							} else {
								for j := 0; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !ores {
											return "", false
										}
										respuesta += res
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
								lista_apt2 := bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										info_carpeta := bloque_carpeta.B_content
										if i == 0 {
											for j := 2; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
													if !ores {
														return "", false
													}
													respuesta += res
												}
											}
										} else {
											for j := 0; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
													if !ores {
														return "", false
													}
													respuesta += res
												}
											}
										}
									}
								}
								continue
								//
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
										res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !ores {
											return "", false
										}
										respuesta += res
									}
								}
							} else {
								for j := 0; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !ores {
											return "", false
										}
										respuesta += res
									}
								}
							}
						}
					}
					continue
				} else if i == 15 { //i == 15
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return "", false
					}
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
											lista_apt3 := bloque_apuntador3.B_pointers
											for c := range lista_apt1 {
												apt3 := lista_apt3[c]
												if apt3 != -1 {
													bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt3)
													if !eba {
														return "", false
													}
													info_carpeta := bloque_carpeta.B_content
													if i == 0 {
														for j := 2; j < 4; j++ {
															info := info_carpeta[j]
															if info.B_inodo != -1 {
																res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
																if !ores {
																	return "", false
																}
																respuesta += res
															}
														}
													} else {
														for j := 0; j < 4; j++ {
															info := info_carpeta[j]
															if info.B_inodo != -1 {
																res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
																if !ores {
																	return "", false
																}
																respuesta += res
															}
														}
													}
												}
											}
											continue
											//
											//
										}
										bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										info_carpeta := bloque_carpeta.B_content
										if i == 0 {
											for j := 2; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
													if !ores {
														return "", false
													}
													respuesta += res
												}
											}
										} else {
											for j := 0; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
													if !ores {
														return "", false
													}
													respuesta += res
												}
											}
										}
									}
								}
								continue
								//
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
										res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !ores {
											return "", false
										}
										respuesta += res
									}
								}
							} else {
								for j := 0; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
										if !ores {
											return "", false
										}
										respuesta += res
									}
								}
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
								res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
								if !ores {
									return "", false
								}
								respuesta += res
							}
						}
					} else {
						for j := 0; j < 4; j++ {
							info := info_carpeta[j]
							if info.B_inodo != -1 {
								res, ores := Obtener_Codigo_LS(comando, path, superbloque, info.B_inodo, id_log, ToString(info.B_name[:]))
								if !ores {
									return "", false
								}
								respuesta += res
							}
						}
					}
				}
			}
		}
		dato += "<TR>\n"
		dato += "<TD ALIGN=\"CENTER\">" + fmt.Sprint(inodo.I_perm) + "</TD>"
		//usuario grupo nombre
		us_name, group, eus := Obtener_Usuario_Grupo_Nombre(id_log)
		if !eus {
			return "", false
		}
		dato += "<TD ALIGN=\"CENTER\">" + fmt.Sprint(us_name) + "</TD>"
		dato += "<TD ALIGN=\"CENTER\">" + fmt.Sprint(group) + "</TD>"
		dato += "<TD ALIGN=\"CENTER\">" + fmt.Sprint(inodo.I_s) + "</TD>"
		dato += "<TD ALIGN=\"CENTER\">" + IntFechaToStr(inodo.I_atime) + "</TD>"
		dato += "<TD ALIGN=\"CENTER\">" + IntFechaToStr(inodo.I_mtime) + "</TD>"
		dato += "<TD ALIGN=\"CENTER\">CARPETA</TD>"
		dato += "<TD ALIGN=\"CENTER\">" + nombre + "</TD>"
		dato += "</TR>\n"

		dato += respuesta
		return dato, true
	} else {
		dato += "<TR>\n"
		dato += "<TD ALIGN=\"CENTER\">" + fmt.Sprint(inodo.I_perm) + "</TD>"
		//usuario grupo nombre
		us_name, group, eus := Obtener_Usuario_Grupo_Nombre(id_log)
		if !eus {
			return "", false
		}

		//-------------------
		//-------------------
		//-------------------
		dato += "<TD ALIGN=\"CENTER\">" + fmt.Sprint(us_name) + "</TD>"
		dato += "<TD ALIGN=\"CENTER\">" + fmt.Sprint(group) + "</TD>"
		dato += "<TD ALIGN=\"CENTER\">" + fmt.Sprint(inodo.I_s) + "</TD>"
		dato += "<TD ALIGN=\"CENTER\">" + IntFechaToStr(inodo.I_atime) + "</TD>"
		dato += "<TD ALIGN=\"CENTER\">" + IntFechaToStr(inodo.I_mtime) + "</TD>"
		dato += "<TD ALIGN=\"CENTER\">ARCHIVO</TD>"
		dato += "<TD ALIGN=\"CENTER\">" + nombre + "</TD>"
		dato += "</TR>\n"
		return dato, true

	}
	// return "", false
}

func Obtener_Usuario_Grupo_Nombre(id_disco string) (string, string, bool) {
	contenido_usuarios, eus := Obtener_Contenido_Archivo_Users("REP", id_disco)
	if !eus {
		return "", "", false
	}

	// id_usuario := int32(-1)
	splita := strings.Split(contenido_usuarios, "\n")
	for _, valor := range splita {
		if strings.Contains(valor, ",U,") {
			splitb := strings.Split(valor, ",")
			// numero_us, _ := strconv.Atoi(splitb[0])
			// id_usuario = int32(numero_us)
			if splitb[3] == ToString(UsuarioLogeado.User[:]) {
				return splitb[3], splitb[2], true
			}
		}
	}
	return "", "", false
}
