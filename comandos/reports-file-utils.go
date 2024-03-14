package comandos

import (
	"MIA_P1_201906051/structures"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

func Obtener_Contenido_Archivo(id_disco string, path string, id_user string, grupo_user string) (string, bool) {
	comando := "REP"
	conjunto, route, eco := Obtener_Particion_ID(id_disco)
	if !eco {
		return "", false
	}

	superbloque, esb := ReducirSuperBloqueObtener(route, id_disco, conjunto)
	if !esb {
		return "", false
	}

	ruta_separada := strings.Split(path, "/")
	nombre_archio := ruta_separada[len(ruta_separada)-1]
	ruta_sin_archivo := strings.ReplaceAll(path, "/"+nombre_archio, "")

	numero_inodo_carpeta, enic := Encontrar_Ruta("REP", route, superbloque.S_inode_start, superbloque.S_block_start, ruta_sin_archivo)
	if !enic {
		return "", false
	}

	if numero_inodo_carpeta == -1 {
		return "", false
	}

	inodo_carpeta, ei := Obtener_Inodo("REP", route, superbloque.S_inode_start, numero_inodo_carpeta)
	if !ei {
		return "", false
	}
	lista_apuntadores := inodo_carpeta.I_block
	numero_inodo_archivo := int32(-1)
	for i := range lista_apuntadores {
		apuntador := lista_apuntadores[i]
		if numero_inodo_archivo != -1 {
			break
		}
		if apuntador != -1 {
			if i == 13 {
				bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, apuntador)
				if !eba1 {
					return "", false
				}
				ap_temp := reflect.TypeOf(structures.BloqueApuntadores{})
				if ap_temp == reflect.TypeOf(bloque_apuntador1) {
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							bloque_carpeta, ebc := Obtener_Bloque(comando, route, superbloque.S_block_start, apt1)
							if !ebc {
								return "", false
							}
							contenido_carpeta := bloque_carpeta.B_content
							for _, content := range contenido_carpeta {
								cont_temp := reflect.TypeOf(structures.Content{})
								if cont_temp == reflect.TypeOf(content) {
									//Continuar----
									contenido_carpeta := bloque_carpeta.B_content
									for _, contenido := range contenido_carpeta {
										if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archio)) {
											numero_inodo_archivo = contenido.B_inodo
											break
										}
									}
									//Continuar----
								}
							}
						}
					}
					continue
				}
			} else if i == 14 {
				bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, apuntador)
				if !eba1 {
					return "", false
				}
				ap_temp := reflect.TypeOf(structures.BloqueApuntadores{})
				if ap_temp == reflect.TypeOf(bloque_apuntador1) {
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								//Continuar
								//Continuar
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								ap_temp := reflect.TypeOf(structures.BloqueApuntadores{})
								if ap_temp == reflect.TypeOf(bloque_apuntador2) {
									lista_apt2 := bloque_apuntador2.B_pointers
									for b := range lista_apt2 {
										apt2 := lista_apt2[b]
										if apt2 != -1 {
											bloque_carpeta, ebc := Obtener_Bloque(comando, route, superbloque.S_block_start, apt2)
											if !ebc {
												return "", false
											}
											contenido_carpeta := bloque_carpeta.B_content
											for _, content := range contenido_carpeta {
												cont_temp := reflect.TypeOf(structures.Content{})
												if cont_temp == reflect.TypeOf(content) {
													//Continuar----
													contenido_carpeta := bloque_carpeta.B_content
													for _, contenido := range contenido_carpeta {
														if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archio)) {
															numero_inodo_archivo = contenido.B_inodo
															break
														}
													}
													//Continuar----
												}
											}
										}
									}
									continue
								}
								//Continuar
								//Continuar
							}
							bloque_carpeta, ebc := Obtener_Bloque(comando, route, superbloque.S_block_start, apt1)
							if !ebc {
								return "", false
							}
							contenido_carpeta := bloque_carpeta.B_content
							for _, content := range contenido_carpeta {
								cont_temp := reflect.TypeOf(structures.Content{})
								if cont_temp == reflect.TypeOf(content) {
									//Continuar----
									contenido_carpeta := bloque_carpeta.B_content
									for _, contenido := range contenido_carpeta {
										if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archio)) {
											numero_inodo_archivo = contenido.B_inodo
											break
										}
									}
									//Continuar----
								}
							}
						}
					}
					continue
				}
			} else if i == 15 {
				bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, apuntador)
				if !eba1 {
					return "", false
				}
				ap_temp := reflect.TypeOf(structures.BloqueApuntadores{})
				if ap_temp == reflect.TypeOf(bloque_apuntador1) {
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								//Continuar
								//Continuar
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								ap_temp := reflect.TypeOf(structures.BloqueApuntadores{})
								if ap_temp == reflect.TypeOf(bloque_apuntador2) {
									lista_apt2 := bloque_apuntador2.B_pointers
									for b := range lista_apt2 {
										apt2 := lista_apt2[b]
										if apt2 != -1 {
											if b == 0 {
												//Tercera parte
												//Tercera parte
												bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, apt2)
												if !eba3 {
													return "", false
												}
												ap_temp := reflect.TypeOf(structures.BloqueApuntadores{})
												if ap_temp == reflect.TypeOf(bloque_apuntador3) {
													lista_apt3 := bloque_apuntador3.B_pointers
													for a := range lista_apt3 {
														apt3 := lista_apt3[a]
														if apt3 != -1 {
															bloque_carpeta, ebc := Obtener_Bloque(comando, route, superbloque.S_block_start, apt3)
															if !ebc {
																return "", false
															}
															contenido_carpeta := bloque_carpeta.B_content
															for _, content := range contenido_carpeta {
																cont_temp := reflect.TypeOf(structures.Content{})
																if cont_temp == reflect.TypeOf(content) {
																	//Continuar----
																	contenido_carpeta := bloque_carpeta.B_content
																	for _, contenido := range contenido_carpeta {
																		if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archio)) {
																			numero_inodo_archivo = contenido.B_inodo
																			break
																		}
																	}
																	//Continuar----
																}
															}
														}
													}
													continue
												}
												//Tercera parte
												//Tercera parte
											}
											bloque_carpeta, ebc := Obtener_Bloque(comando, route, superbloque.S_block_start, apt2)
											if !ebc {
												return "", false
											}
											contenido_carpeta := bloque_carpeta.B_content
											for _, content := range contenido_carpeta {
												cont_temp := reflect.TypeOf(structures.Content{})
												if cont_temp == reflect.TypeOf(content) {
													//Continuar----
													contenido_carpeta := bloque_carpeta.B_content
													for _, contenido := range contenido_carpeta {
														if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archio)) {
															numero_inodo_archivo = contenido.B_inodo
															break
														}
													}
													//Continuar----
												}
											}
										}
									}
									continue
								}
								//Continuar
								//Continuar
							}
							bloque_carpeta, ebc := Obtener_Bloque(comando, route, superbloque.S_block_start, apt1)
							if !ebc {
								return "", false
							}
							contenido_carpeta := bloque_carpeta.B_content
							for _, content := range contenido_carpeta {
								cont_temp := reflect.TypeOf(structures.Content{})
								if cont_temp == reflect.TypeOf(content) {
									//Continuar----
									contenido_carpeta := bloque_carpeta.B_content
									for _, contenido := range contenido_carpeta {
										if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archio)) {
											numero_inodo_archivo = contenido.B_inodo
											break
										}
									}
									//Continuar----
								}
							}
						}
					}
					continue
				}
			} else {
				//
				bloque_carpeta, ebc := Obtener_Bloque(comando, route, superbloque.S_block_start, apuntador)
				if !ebc {
					return "", false
				}
				contenido_carpeta := bloque_carpeta.B_content
				for _, content := range contenido_carpeta {
					cont_temp := reflect.TypeOf(structures.Content{})
					if cont_temp == reflect.TypeOf(content) {
						//Continuar----
						contenido_carpeta := bloque_carpeta.B_content
						for _, contenido := range contenido_carpeta {
							if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archio)) {
								numero_inodo_archivo = contenido.B_inodo
								break
							}
						}
						//Continuar----
					}
				}
				//
			}
		}
	}
	if numero_inodo_archivo != -1 {
		inodo_archivo, eia := Obtener_Inodo(comando, route, superbloque.S_inode_start, numero_inodo_archivo)
		if !eia {
			return "", false
		}
		if !Validar_Permisos(comando, inodo_archivo.I_uid, inodo_archivo.I_gid, ToInt(id_user), ToInt(grupo_user), inodo_archivo.I_perm, 4) {
			color.Red("No se tienen premisos necesarios")
			return "", false
		}
		fecha := ObFechaInt()
		inodo_carpeta.I_atime = fecha
		apuntadores_archivo := inodo_archivo.I_block
		contador := 0
		contenido_archivo := ""
		for i := range apuntadores_archivo {
			parte_archivo := apuntadores_archivo[i]
			if parte_archivo != -1 {
				if i == 13 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, parte_archivo)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							bloque_archivo, eba := Obtener_Bloque_Archivo(comando, route, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							arch_temp := reflect.TypeOf(structures.BloqueArchivos{})
							if arch_temp == reflect.TypeOf(bloque_archivo) {
								contenido := ToString(bloque_archivo.B_content[:])
								contenido_archivo += ToString([]byte(contenido))
							}
						}
					}
					continue
				} else if i == 14 { //i == 14
					//
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, parte_archivo)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								//
								//
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, apt1)
								if !eba2 {
									return "", false
								}
								lista_apt2 := bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										bloque_archivo, eba := Obtener_Bloque_Archivo(comando, route, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										arch_temp := reflect.TypeOf(structures.BloqueArchivos{})
										if arch_temp == reflect.TypeOf(bloque_archivo) {
											contenido := ToString(bloque_archivo.B_content[:])
											contenido_archivo += ToString([]byte(contenido))
										}
									}
								}
								continue
								//
								//
							}
							bloque_archivo, eba := Obtener_Bloque_Archivo(comando, route, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							arch_temp := reflect.TypeOf(structures.BloqueArchivos{})
							if arch_temp == reflect.TypeOf(bloque_archivo) {
								contenido := ToString(bloque_archivo.B_content[:])
								contenido_archivo += ToString([]byte(contenido))
							}
						}
					}
					continue
					//
				} else if i == 15 {
					//
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, parte_archivo)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								//
								//
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, apt1)
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
											bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, route, superbloque.S_block_start, apt2)
											if !eba3 {
												return "", false
											}
											lista_apt3 := bloque_apuntador3.B_pointers
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												if apt3 != -1 {
													bloque_archivo, eba := Obtener_Bloque_Archivo(comando, route, superbloque.S_block_start, apt3)
													if !eba {
														return "", false
													}
													arch_temp := reflect.TypeOf(structures.BloqueArchivos{})
													if arch_temp == reflect.TypeOf(bloque_archivo) {
														contenido := ToString(bloque_archivo.B_content[:])
														contenido_archivo += ToString([]byte(contenido))
													}
												}
											}
											continue
											//
											//
										}
										bloque_archivo, eba := Obtener_Bloque_Archivo(comando, route, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										arch_temp := reflect.TypeOf(structures.BloqueArchivos{})
										if arch_temp == reflect.TypeOf(bloque_archivo) {
											contenido := ToString(bloque_archivo.B_content[:])
											contenido_archivo += ToString([]byte(contenido))
										}
									}
								}
								continue
								//
								//
							}
							bloque_archivo, eba := Obtener_Bloque_Archivo(comando, route, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							arch_temp := reflect.TypeOf(structures.BloqueArchivos{})
							if arch_temp == reflect.TypeOf(bloque_archivo) {
								contenido := ToString(bloque_archivo.B_content[:])
								contenido_archivo += ToString([]byte(contenido))
							}
						}
					}
					continue
					//

				} else {
					bloque_archivo, eba := Obtener_Bloque_Archivo(comando, route, superbloque.S_block_start, parte_archivo)
					if !eba {
						return "", false
					}
					contenido := ToString(bloque_archivo.B_content[:])
					contenido_archivo += ToString([]byte(contenido))
				}
			}
			contador += 1
		}
		//retornamos el contenido del archivo
		return contenido_archivo, true
	} else {
		color.Red("[" + comando + "]: Archivo no existe en la ruta")
		return "", false
	}

}
