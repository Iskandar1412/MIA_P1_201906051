package comandos

import (
	"MIA_P1_201906051/structures"
	"reflect"

	"github.com/fatih/color"
)

func Recalcular_Size_Carpetas_Recursivo(comando string, id string, num_inodo int32) (int32, bool) {
	conjunto, path, er := Obtener_Particion_ID(id)
	if !er {
		return 0, false
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
				return 0, false
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
			// esparticion = true
			superbloque, esb = Obtener_Superbloque(comando, path, ToString(particion.Part_name[:]))
			if !esb {
				color.Red("[" + comando + "]: Error al obtener superbloque")
				return 0, false
			}
			if particion.Part_type == 'E' {
				color.Red("[" + comando + "]: No se puede obtener información de particion extendida")
				return 0, false
			}
		}
	}

	inodo, ei := Obtener_Inodo(comando, path, superbloque.S_inode_start, num_inodo)
	if !ei {
		return 0, false
	}
	// dato := ""
	respuesta := int32(0)
	if inodo.I_type == 0 {
		for i := range inodo.I_block {
			apt := inodo.I_block[i]
			if apt != -1 {
				if i == 13 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return 0, false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return 0, false
							}
							info_carpeta := bloque_carpeta.B_content
							if i == 0 {
								for j := 2; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
										respuesta += size_total
									}
								}
							} else {
								for j := range info_carpeta {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
										respuesta += size_total
									}
								}
							}
						}
					}
					continue
				} else if i == 14 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return 0, false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return 0, false
								}
								lista_apt2 := bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return 0, false
										}
										info_carpeta := bloque_carpeta.B_content
										if i == 0 {
											for j := 2; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
													respuesta += size_total
												}
											}
										} else {
											for j := range info_carpeta {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
													respuesta += size_total
												}
											}
										}
									}
								}
								continue
							}
							bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return 0, false
							}
							info_carpeta := bloque_carpeta.B_content
							if i == 0 {
								for j := 2; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
										respuesta += size_total
									}
								}
							} else {
								for j := range info_carpeta {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
										respuesta += size_total
									}
								}
							}
						}
					}
					continue
				} else if i == 15 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt)
					if !eba1 {
						return 0, false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							if a == 0 {
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return 0, false
								}
								lista_apt2 := bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 != -1 {
										if b == 0 {
											bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt2)
											if !eba3 {
												return 0, false
											}
											lista_apt3 := bloque_apuntador3.B_pointers
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												if apt3 != -1 {
													bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt3)
													if !eba {
														return 0, false
													}
													info_carpeta := bloque_carpeta.B_content
													if i == 0 {
														for j := 2; j < 4; j++ {
															info := info_carpeta[j]
															if info.B_inodo != -1 {
																size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
																respuesta += size_total
															}
														}
													} else {
														for j := range info_carpeta {
															info := info_carpeta[j]
															if info.B_inodo != -1 {
																size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
																respuesta += size_total
															}
														}
													}
												}
											}
											continue
										}
										bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return 0, false
										}
										info_carpeta := bloque_carpeta.B_content
										if i == 0 {
											for j := 2; j < 4; j++ {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
													respuesta += size_total
												}
											}
										} else {
											for j := range info_carpeta {
												info := info_carpeta[j]
												if info.B_inodo != -1 {
													size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
													respuesta += size_total
												}
											}
										}
									}
								}
								continue
							}
							bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return 0, false
							}
							info_carpeta := bloque_carpeta.B_content
							if i == 0 {
								for j := 2; j < 4; j++ {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
										respuesta += size_total
									}
								}
							} else {
								for j := range info_carpeta {
									info := info_carpeta[j]
									if info.B_inodo != -1 {
										size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
										respuesta += size_total
									}
								}
							}
						}
					}
					continue
				} else {
					bloque_carpeta, eba := Obtener_Bloque(comando, path, superbloque.S_block_start, apt)
					if !eba {
						return 0, false
					}
					info_carpeta := bloque_carpeta.B_content
					if i == 0 {
						for j := 2; j < 4; j++ {
							info := info_carpeta[j]
							if info.B_inodo != -1 {
								size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
								respuesta += size_total
							}
						}
					} else {
						for j := range info_carpeta {
							info := info_carpeta[j]
							if info.B_inodo != -1 {
								size_total, _ := Recalcular_Size_Carpetas_Recursivo(comando, id, info.B_inodo)
								respuesta += size_total
							}
						}
					}
				}
			}
		}
		//asignar valor a carpeta
		inodo.I_s = respuesta
		Guardar_Inodo(comando, path, superbloque.S_inode_start, inodo, num_inodo)
		return respuesta, true
	} else {
		//spñp se retorna el tamaño para sumar a carpeta
		return inodo.I_s, true
	}
}
