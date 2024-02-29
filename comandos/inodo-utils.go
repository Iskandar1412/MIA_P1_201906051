package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Obtener_Inodo(comando string, path string, inicio int32, numero_inodo int32) (structures.Inode, bool) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return structures.Inode{}, false
	}
	defer file.Close()
	//bytes_con_inodos := inicio_inodos + numero
	posicion_numero_inodo := inicio + (size.SizeInode() * numero_inodo)
	if _, err := file.Seek(int64(posicion_numero_inodo), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return structures.Inode{}, false
	}
	inodo := structures.Inode{}
	if err := binary.Read(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[" + comando + "]: Error en lectura de bloque")
		return structures.Inode{}, false
	}
	return inodo, true
}

func Obtener_Inodo_Disponible(comando string, path string, bm_inode_start int32, inodes_count int32) (int32, bool) {
	for i := 0; i < int(inodes_count); i++ {
		// struct_bitmap := '\x00'
		file, err := os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			color.Red("[" + comando + "]: Error al leer Archivo")
			return 0, false
		}
		defer file.Close()
		//bytes_con_inodos := inicio_inodos + numero
		if _, err := file.Seek(int64(bm_inode_start)+int64(i), 0); err != nil {
			color.Red("[" + comando + "]: Error al mover el puntero")
			return 0, false
		}
		var valor int8
		if err := binary.Read(file, binary.LittleEndian, &valor); err != nil {
			color.Red("[" + comando + "]: Error en lectura de bloque")
			return 0, false
		}
		if valor == 0 { //Inodo libre
			return int32(i), true
		}

	}
	return 0, false
}

func Guardar_Inodo(comando string, path string, inicio int32, inodo structures.Inode, numeroinodo int32) bool {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return false
	}
	posicion_inodo := inicio + (size.SizeInode() * numeroinodo)
	if _, err := file.Seek(int64(posicion_inodo), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return false
	}

	if err := binary.Write(file, binary.LittleEndian, &inodo); err != nil {
		color.Red("[" + comando + "]: Error en escritura de inodo")
		return false
	}
	return true
}

func Agregar_Bloque_Lista_Inodos(comando string, path string, inicio_bloque int32, numero_inodo_disponible int32, inicio_bitmap_bloque int32, numero_bloques_total int32, nombre_carpeta string, inodo_padre structures.Inode, inicio_inodo int32, numero_inodo_padre int32, inicio_bitmap_inodo int32, numero_bloque_disponible int32) {
	apuntadores_padre := inodo_padre.I_block
	contenido_vacio := structures.Content{B_name: NameArchivosByte(""), B_inodo: int32(-1)}
	contador := 0
	for _, apuntadores := range apuntadores_padre { //recorrer apuntadores padre
		if apuntadores != -1 {
			if contador == 13 {
				bloque_apt, eba := Obtener_Bloque_Apuntador(comando, path, inicio_bloque, apuntadores)
				if !eba {
					return
				}
				lista_apt := bloque_apt.B_pointers
				for i := range lista_apt {
					apt := lista_apt[i]
					if apt != -1 {
						bloque_padre, ebp := Obtener_Bloque(comando, path, inicio_bloque, apt)
						if !ebp {
							return
						}
						lista_bloque_padre := bloque_padre.B_content
						for _, apt_bloque_padre := range lista_bloque_padre {
							if apt_bloque_padre.B_inodo == -1 {
								lista_bloque_padre[i].B_inodo = numero_inodo_disponible
								lista_bloque_padre[i].B_name = NameArchivosByte(nombre_carpeta)
								// apt_bloque_padre.B_inodo = numero_inodo_disponible
								// apt_bloque_padre.B_name = NameArchivosByte(nombre_carpeta)
								var lista_enviar []interface{}
								lista_enviar = append(lista_enviar, lista_bloque_padre[0])
								lista_enviar = append(lista_enviar, lista_bloque_padre[1])
								lista_enviar = append(lista_enviar, lista_bloque_padre[2])
								lista_enviar = append(lista_enviar, lista_bloque_padre[3])

								Modificar_Carpeta(comando, path, inicio_bloque, apt, lista_enviar)
								return
							}
						}
					} else {
						var contenido_carpeta []interface{}
						carpeta_nueva := structures.Content{B_name: NameArchivosByte(nombre_carpeta), B_inodo: numero_inodo_disponible}
						contenido_carpeta = append(contenido_carpeta, carpeta_nueva)
						contenido_carpeta = append(contenido_carpeta, contenido_vacio)
						contenido_carpeta = append(contenido_carpeta, contenido_vacio)
						contenido_carpeta = append(contenido_carpeta, contenido_vacio)
						nuevo_bloque_padre, enbp := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
						if !enbp {
							return
						}
						Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, nuevo_bloque_padre)
						Modificar_Carpeta(comando, path, inicio_bloque, nuevo_bloque_padre, contenido_carpeta)
						Modificar_Bitmap(comando, path, inicio_bitmap_bloque, nuevo_bloque_padre, 1)

						//modificar apuntador inodo
						lista_apt[i] = nuevo_bloque_padre
						lista_a := structures.BloqueApuntadores{B_pointers: lista_apt}
						Modificar_Apuntador(comando, path, inicio_bloque, apuntadores, lista_a)
						return
					}
				}
				contador += 1
				continue
			} else if contador == 14 {
				bloque_apt, eba := Obtener_Bloque_Apuntador(comando, path, inicio_bloque, apuntadores)
				if !eba {
					return
				}
				lista_apt := bloque_apt.B_pointers
				for i := range lista_apt {
					apt := lista_apt[i]
					if apt != -1 {
						if i == 0 {
							bloque_apt2, eba2 := Obtener_Bloque_Apuntador(comando, path, inicio_bloque, apt)
							if !eba2 {
								return
							}
							lista_apt2 := bloque_apt2.B_pointers
							for j := range lista_apt2 {
								apt2 := lista_apt2[j]
								if apt2 != -1 {
									bloque_padre, ebp := Obtener_Bloque(comando, path, inicio_bloque, apt2)
									if !ebp {
										return
									}
									lista_bloque_padre := bloque_padre.B_content
									for _, apt_bloque_padre := range lista_bloque_padre {
										if apt_bloque_padre.B_inodo == -1 {
											lista_bloque_padre[i].B_inodo = numero_inodo_disponible
											lista_bloque_padre[i].B_name = NameArchivosByte(nombre_carpeta)
											// apt_bloque_padre.B_inodo = numero_inodo_disponible
											// apt_bloque_padre.B_name = NameArchivosByte(nombre_carpeta)
											var lista_enviar []interface{}
											lista_enviar = append(lista_enviar, lista_bloque_padre[0])
											lista_enviar = append(lista_enviar, lista_bloque_padre[1])
											lista_enviar = append(lista_enviar, lista_bloque_padre[2])
											lista_enviar = append(lista_enviar, lista_bloque_padre[3])

											Modificar_Carpeta(comando, path, inicio_bloque, apt2, lista_enviar)
											return
										}
									}
								} else {
									var contenido_carpeta []interface{}
									carpeta_nueva := structures.Content{B_name: NameArchivosByte(nombre_carpeta), B_inodo: numero_inodo_disponible}
									contenido_carpeta = append(contenido_carpeta, carpeta_nueva)
									contenido_carpeta = append(contenido_carpeta, contenido_vacio)
									contenido_carpeta = append(contenido_carpeta, contenido_vacio)
									contenido_carpeta = append(contenido_carpeta, contenido_vacio)
									nuevo_bloque_padre, enbp := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
									if !enbp {
										return
									}
									Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, nuevo_bloque_padre)
									Modificar_Carpeta(comando, path, inicio_bloque, nuevo_bloque_padre, contenido_carpeta)
									Modificar_Bitmap(comando, path, inicio_bitmap_bloque, nuevo_bloque_padre, 1)

									//modificar apuntador inodo
									lista_apt2[j] = nuevo_bloque_padre
									lista_a := structures.BloqueApuntadores{B_pointers: lista_apt2}
									Modificar_Apuntador(comando, path, inicio_bloque, apuntadores, lista_a)
									return
								}
							}
							continue
						}
						bloque_padre, ebp := Obtener_Bloque(comando, path, inicio_bloque, apt)
						if !ebp {
							return
						}
						lista_bloque_padre := bloque_padre.B_content
						for _, apt_bloque_padre := range lista_bloque_padre {
							if apt_bloque_padre.B_inodo == -1 {
								lista_bloque_padre[i].B_inodo = numero_inodo_disponible
								lista_bloque_padre[i].B_name = NameArchivosByte(nombre_carpeta)
								// apt_bloque_padre.B_inodo = numero_inodo_disponible
								// apt_bloque_padre.B_name = NameArchivosByte(nombre_carpeta)
								var lista_enviar []interface{}
								lista_enviar = append(lista_enviar, lista_bloque_padre[0])
								lista_enviar = append(lista_enviar, lista_bloque_padre[1])
								lista_enviar = append(lista_enviar, lista_bloque_padre[2])
								lista_enviar = append(lista_enviar, lista_bloque_padre[3])

								Modificar_Carpeta(comando, path, inicio_bloque, apt, lista_enviar)
								return
							}
						}
					} else {
						var contenido_carpeta []interface{}
						carpeta_nueva := structures.Content{B_name: NameArchivosByte(nombre_carpeta), B_inodo: numero_inodo_disponible}
						contenido_carpeta = append(contenido_carpeta, carpeta_nueva)
						contenido_carpeta = append(contenido_carpeta, contenido_vacio)
						contenido_carpeta = append(contenido_carpeta, contenido_vacio)
						contenido_carpeta = append(contenido_carpeta, contenido_vacio)
						nuevo_bloque_padre, enbp := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
						if !enbp {
							return
						}
						Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, nuevo_bloque_padre)
						Modificar_Carpeta(comando, path, inicio_bloque, nuevo_bloque_padre, contenido_carpeta)
						Modificar_Bitmap(comando, path, inicio_bitmap_bloque, nuevo_bloque_padre, 1)

						//modificar apuntador inodo
						lista_apt[i] = nuevo_bloque_padre
						lista_a := structures.BloqueApuntadores{B_pointers: lista_apt}
						Modificar_Apuntador(comando, path, inicio_bloque, apuntadores, lista_a)
						return
					}
				}
				contador += 1
				continue
			} else if contador == 15 {
				//apuntador triple
				bloque_apt, eba := Obtener_Bloque_Apuntador(comando, path, inicio_bloque, apuntadores)
				if !eba {
					return
				}
				lista_apt := bloque_apt.B_pointers
				for i := range lista_apt {
					apt := lista_apt[i]
					if apt != -1 {
						if i == 0 {
							bloque_apt2, eba2 := Obtener_Bloque_Apuntador(comando, path, inicio_bloque, apt)
							if !eba2 {
								return
							}
							lista_apt2 := bloque_apt2.B_pointers
							for j := range lista_apt2 {
								apt2 := lista_apt2[j]
								if apt2 != -1 {
									if j == 0 {
										bloque_apt3, eba3 := Obtener_Bloque_Apuntador(comando, path, inicio_bloque, apt2)
										if !eba3 {
											return
										}
										lista_apt3 := bloque_apt3.B_pointers
										for k := range lista_apt3 {
											apt3 := lista_apt3[k]
											if apt3 != -1 {
												bloque_padre, ebp3 := Obtener_Bloque(comando, path, inicio_bloque, apt3)
												if !ebp3 {
													return
												}
												lista_bloque_padre := bloque_padre.B_content
												for _, apt_bloque_padre := range lista_bloque_padre {
													if apt_bloque_padre.B_inodo == -1 {
														lista_bloque_padre[i].B_inodo = numero_inodo_disponible
														lista_bloque_padre[i].B_name = NameArchivosByte(nombre_carpeta)
														// apt_bloque_padre.B_inodo = numero_inodo_disponible
														// apt_bloque_padre.B_name = NameArchivosByte(nombre_carpeta)
														var lista_enviar []interface{}
														lista_enviar = append(lista_enviar, lista_bloque_padre[0])
														lista_enviar = append(lista_enviar, lista_bloque_padre[1])
														lista_enviar = append(lista_enviar, lista_bloque_padre[2])
														lista_enviar = append(lista_enviar, lista_bloque_padre[3])

														Modificar_Carpeta(comando, path, inicio_bloque, apt3, lista_enviar)
														return
													}
												}
											} else {
												var contenido_carpeta []interface{}
												carpeta_nueva := structures.Content{B_name: NameArchivosByte(nombre_carpeta), B_inodo: numero_inodo_disponible}
												contenido_carpeta = append(contenido_carpeta, carpeta_nueva)
												contenido_carpeta = append(contenido_carpeta, contenido_vacio)
												contenido_carpeta = append(contenido_carpeta, contenido_vacio)
												contenido_carpeta = append(contenido_carpeta, contenido_vacio)
												nuevo_bloque_padre, enbp := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
												if !enbp {
													return
												}
												Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, nuevo_bloque_padre)
												Modificar_Carpeta(comando, path, inicio_bloque, nuevo_bloque_padre, contenido_carpeta)
												Modificar_Bitmap(comando, path, inicio_bitmap_bloque, nuevo_bloque_padre, 1)

												//modificar apuntador inodo
												lista_apt3[k] = nuevo_bloque_padre
												lista_a := structures.BloqueApuntadores{B_pointers: lista_apt3}
												Modificar_Apuntador(comando, path, inicio_bloque, apuntadores, lista_a)
												return
											}
										}
									}
									bloque_padre, ebp := Obtener_Bloque(comando, path, inicio_bloque, apt2)
									if !ebp {
										return
									}
									lista_bloque_padre := bloque_padre.B_content
									for _, apt_bloque_padre := range lista_bloque_padre {
										if apt_bloque_padre.B_inodo == -1 {
											lista_bloque_padre[i].B_inodo = numero_inodo_disponible
											lista_bloque_padre[i].B_name = NameArchivosByte(nombre_carpeta)
											// apt_bloque_padre.B_inodo = numero_inodo_disponible
											// apt_bloque_padre.B_name = NameArchivosByte(nombre_carpeta)
											var lista_enviar []interface{}
											lista_enviar = append(lista_enviar, lista_bloque_padre[0])
											lista_enviar = append(lista_enviar, lista_bloque_padre[1])
											lista_enviar = append(lista_enviar, lista_bloque_padre[2])
											lista_enviar = append(lista_enviar, lista_bloque_padre[3])

											Modificar_Carpeta(comando, path, inicio_bloque, apt2, lista_enviar)
											return
										}
									}
								} else {
									var contenido_carpeta []interface{}
									carpeta_nueva := structures.Content{B_name: NameArchivosByte(nombre_carpeta), B_inodo: numero_inodo_disponible}
									contenido_carpeta = append(contenido_carpeta, carpeta_nueva)
									contenido_carpeta = append(contenido_carpeta, contenido_vacio)
									contenido_carpeta = append(contenido_carpeta, contenido_vacio)
									contenido_carpeta = append(contenido_carpeta, contenido_vacio)
									nuevo_bloque_padre, enbp := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
									if !enbp {
										return
									}
									Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, nuevo_bloque_padre)
									Modificar_Carpeta(comando, path, inicio_bloque, nuevo_bloque_padre, contenido_carpeta)
									Modificar_Bitmap(comando, path, inicio_bitmap_bloque, nuevo_bloque_padre, 1)

									//modificar apuntador inodo
									lista_apt2[j] = nuevo_bloque_padre
									lista_a := structures.BloqueApuntadores{B_pointers: lista_apt2}
									Modificar_Apuntador(comando, path, inicio_bloque, apuntadores, lista_a)
									return
								}
							}
							continue
						}
						bloque_padre, ebp := Obtener_Bloque(comando, path, inicio_bloque, apt)
						if !ebp {
							return
						}
						lista_bloque_padre := bloque_padre.B_content
						for _, apt_bloque_padre := range lista_bloque_padre {
							if apt_bloque_padre.B_inodo == -1 {
								lista_bloque_padre[i].B_inodo = numero_inodo_disponible
								lista_bloque_padre[i].B_name = NameArchivosByte(nombre_carpeta)
								// apt_bloque_padre.B_inodo = numero_inodo_disponible
								// apt_bloque_padre.B_name = NameArchivosByte(nombre_carpeta)
								var lista_enviar []interface{}
								lista_enviar = append(lista_enviar, lista_bloque_padre[0])
								lista_enviar = append(lista_enviar, lista_bloque_padre[1])
								lista_enviar = append(lista_enviar, lista_bloque_padre[2])
								lista_enviar = append(lista_enviar, lista_bloque_padre[3])

								Modificar_Carpeta(comando, path, inicio_bloque, apt, lista_enviar)
								return
							}
						}
					} else {
						var contenido_carpeta []interface{}
						carpeta_nueva := structures.Content{B_name: NameArchivosByte(nombre_carpeta), B_inodo: numero_inodo_disponible}
						contenido_carpeta = append(contenido_carpeta, carpeta_nueva)
						contenido_carpeta = append(contenido_carpeta, contenido_vacio)
						contenido_carpeta = append(contenido_carpeta, contenido_vacio)
						contenido_carpeta = append(contenido_carpeta, contenido_vacio)
						nuevo_bloque_padre, enbp := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
						if !enbp {
							return
						}
						Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, nuevo_bloque_padre)
						Modificar_Carpeta(comando, path, inicio_bloque, nuevo_bloque_padre, contenido_carpeta)
						Modificar_Bitmap(comando, path, inicio_bitmap_bloque, nuevo_bloque_padre, 1)

						//modificar apuntador inodo
						lista_apt[i] = nuevo_bloque_padre
						lista_a := structures.BloqueApuntadores{B_pointers: lista_apt}
						Modificar_Apuntador(comando, path, inicio_bloque, apuntadores, lista_a)
						return
					}
				}
				contador += 1
				if contador == 16 {
					color.Yellow("No hay más espacio en inodo")
					return
				}
			} else {
				bloque_padre, ebp := Obtener_Bloque(comando, path, inicio_bloque, apuntadores)
				if !ebp {
					return
				}
				lista_bloque_padre := bloque_padre.B_content

				for i, apt_bloque_padre := range lista_bloque_padre {
					if apt_bloque_padre.B_inodo == int32(-1) {
						lista_bloque_padre[i].B_inodo = numero_inodo_disponible
						lista_bloque_padre[i].B_name = NameArchivosByte(nombre_carpeta)
						// apt_bloque_padre.B_inodo = numero_inodo_disponible
						// apt_bloque_padre.B_name = NameArchivosByte(nombre_carpeta)
						var lista_enviar []interface{}
						lista_enviar = append(lista_enviar, lista_bloque_padre[0])
						lista_enviar = append(lista_enviar, lista_bloque_padre[1])
						lista_enviar = append(lista_enviar, lista_bloque_padre[2])
						lista_enviar = append(lista_enviar, lista_bloque_padre[3])
						Modificar_Carpeta(comando, path, inicio_bloque, apuntadores, lista_enviar)
						return
					}
				}
			}
		} else {
			if contador == 13 {
				num_bloque_apuntador, enba := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
				if !enba {
					return
				}
				Crear_Bloque_Apuntador_Vacio(comando, path, inicio_bloque, num_bloque_apuntador)
				Modificar_Bitmap(comando, path, inicio_bitmap_bloque, num_bloque_apuntador, 1)

				var contenido_carpeta []interface{}
				contenido := structures.Content{B_name: NameArchivosByte(nombre_carpeta), B_inodo: numero_inodo_disponible}
				contenido_carpeta = append(contenido_carpeta, contenido)
				contenido_carpeta = append(contenido_carpeta, contenido_vacio)
				contenido_carpeta = append(contenido_carpeta, contenido_vacio)
				contenido_carpeta = append(contenido_carpeta, contenido_vacio)

				nuevo_bloque_padre, enbp := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
				if !enbp {
					return
				}
				Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, nuevo_bloque_padre)
				Modificar_Carpeta(comando, path, inicio_bloque, nuevo_bloque_padre, contenido_carpeta)

				//modificar bitmap
				Modificar_Bitmap(comando, path, inicio_bitmap_bloque, nuevo_bloque_padre, 1)

				//guardar en apuntador 13 inodo
				apuntadores_padre[contador] = num_bloque_apuntador
				inodo_padre.I_block = apuntadores_padre
				fecha := ObFechaInt()
				inodo_padre.I_mtime = fecha

				err_in_pad := Guardar_Inodo(comando, path, inicio_inodo, inodo_padre, numero_inodo_padre)
				if !err_in_pad {
					return
				}

				//apuntar a carpeta
				lista_apt_apuntadores := structures.BloqueApuntadores{}
				for i := range lista_apt_apuntadores.B_pointers {
					lista_apt_apuntadores.B_pointers[i] = int32(-1)
				}
				lista_apt_apuntadores.B_pointers[0] = nuevo_bloque_padre

				// lista_a := apuntadores
				Modificar_Apuntador(comando, path, inicio_bloque, num_bloque_apuntador, lista_apt_apuntadores)
				return
			} else if contador == 14 {
				//indirecto doble
				num_bloque_apuntador, enba := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
				if !enba {
					return
				}
				Crear_Bloque_Apuntador_Vacio(comando, path, inicio_bloque, num_bloque_apuntador)
				Modificar_Bitmap(comando, path, inicio_bitmap_bloque, num_bloque_apuntador, 1)

				//apuntador 2
				num_bloque_apuntador2, enba2 := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
				if !enba2 {
					return
				}
				Crear_Bloque_Apuntador_Vacio(comando, path, inicio_bloque, num_bloque_apuntador2)
				Modificar_Bitmap(comando, path, inicio_bitmap_bloque, num_bloque_apuntador2, 1)

				lista_apt_apuntadores := structures.BloqueApuntadores{}
				for i := range lista_apt_apuntadores.B_pointers {
					lista_apt_apuntadores.B_pointers[i] = int32(-1)
				}
				lista_apt_apuntadores.B_pointers[0] = num_bloque_apuntador2

				//modificar apuntadores 1
				//lista_a = lista_apt_apuntadores
				Modificar_Apuntador(comando, path, inicio_bloque, num_bloque_apuntador, lista_apt_apuntadores)

				//estructura carpeta
				var contenido_carpeta []interface{}
				carpeta_nueva := structures.Content{B_name: NameArchivosByte(nombre_carpeta), B_inodo: numero_inodo_disponible}
				contenido_carpeta = append(contenido_carpeta, carpeta_nueva)
				contenido_carpeta = append(contenido_carpeta, contenido_vacio)
				contenido_carpeta = append(contenido_carpeta, contenido_vacio)
				contenido_carpeta = append(contenido_carpeta, contenido_vacio)
				nuevo_bloque_padre, enbp := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
				if !enbp {
					return
				}
				Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, nuevo_bloque_padre)
				Modificar_Carpeta(comando, path, inicio_bloque, nuevo_bloque_padre, contenido_carpeta)

				Modificar_Bitmap(comando, path, inicio_bitmap_bloque, nuevo_bloque_padre, 1)

				//guardar en el apuntador 14 del inodo
				apuntadores_padre[contador] = num_bloque_apuntador
				inodo_padre.I_block = apuntadores_padre
				fecha := ObFechaInt()
				inodo_padre.I_mtime = fecha
				err_in_pa := Guardar_Inodo(comando, path, inicio_inodo, inodo_padre, numero_inodo_padre)
				if !err_in_pa {
					return
				}

				//apuntar a a carpeta
				lista_apt_apuntadores = structures.BloqueApuntadores{}
				for i := range lista_apt_apuntadores.B_pointers {
					lista_apt_apuntadores.B_pointers[i] = int32(-1)
				}
				lista_apt_apuntadores.B_pointers[0] = nuevo_bloque_padre

				Modificar_Apuntador(comando, path, inicio_bloque, num_bloque_apuntador2, lista_apt_apuntadores)
				return
			} else if contador == 15 {
				//apuntador 1
				num_bloque_apuntador, enbc := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
				if !enbc {
					return
				}
				Crear_Bloque_Apuntador_Vacio(comando, path, inicio_bloque, num_bloque_apuntador)
				Modificar_Bitmap(comando, path, inicio_bitmap_bloque, num_bloque_apuntador, 1)

				//apuntador 2
				num_bloque_apuntador2, enbc2 := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
				if !enbc2 {
					return
				}
				Crear_Bloque_Apuntador_Vacio(comando, path, inicio_bloque, num_bloque_apuntador2)
				Modificar_Bitmap(comando, path, inicio_bitmap_bloque, num_bloque_apuntador2, 1)

				//apuntador 3
				num_bloque_apuntador3, enbc3 := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
				if !enbc3 {
					return
				}
				Crear_Bloque_Apuntador_Vacio(comando, path, inicio_bloque, num_bloque_apuntador3)
				Modificar_Bitmap(comando, path, inicio_bitmap_bloque, num_bloque_apuntador3, 1)

				//se apunta al apuntador 3 desde el 2
				lista_apt_apuntadores := structures.BloqueApuntadores{}
				for i := range lista_apt_apuntadores.B_pointers {
					lista_apt_apuntadores.B_pointers[i] = int32(-1)
				}
				lista_apt_apuntadores.B_pointers[0] = num_bloque_apuntador3

				//modificar apuntadores 2
				Modificar_Apuntador(comando, path, inicio_bloque, num_bloque_apuntador2, lista_apt_apuntadores)

				//apuntar del 2 desde 1
				lista_apt_apuntadores = structures.BloqueApuntadores{}
				for i := range lista_apt_apuntadores.B_pointers {
					lista_apt_apuntadores.B_pointers[i] = int32(-1)
				}
				lista_apt_apuntadores.B_pointers[0] = num_bloque_apuntador2

				//modificar apuntadores 1
				Modificar_Apuntador(comando, path, inicio_bloque, num_bloque_apuntador, lista_apt_apuntadores)

				//struct carpeta
				var contenido_carpeta []interface{}
				carpeta_nueva := structures.Content{B_name: NameArchivosByte(nombre_carpeta), B_inodo: numero_inodo_disponible}
				contenido_carpeta = append(contenido_carpeta, carpeta_nueva)
				contenido_carpeta = append(contenido_carpeta, contenido_vacio)
				contenido_carpeta = append(contenido_carpeta, contenido_vacio)
				contenido_carpeta = append(contenido_carpeta, contenido_vacio)

				nuevo_bloque_padre, enbp := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
				if !enbp {
					return
				}
				Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, nuevo_bloque_padre)
				Modificar_Carpeta(comando, path, inicio_bloque, nuevo_bloque_padre, contenido_carpeta)

				//modificar bitmap
				Modificar_Bitmap(comando, path, inicio_bitmap_bloque, nuevo_bloque_padre, 1)

				//guardar en apuntador 14 inodo
				apuntadores_padre[contador] = num_bloque_apuntador
				inodo_padre.I_block = apuntadores_padre
				fecha := ObFechaInt()
				inodo_padre.I_mtime = fecha
				Guardar_Inodo(comando, path, inicio_inodo, inodo_padre, numero_inodo_padre)

				//apuntar a carpeta
				lista_apt_apuntadores = structures.BloqueApuntadores{}
				for i := range lista_apt_apuntadores.B_pointers {
					lista_apt_apuntadores.B_pointers[i] = int32(-1)
				}
				lista_apt_apuntadores.B_pointers[0] = nuevo_bloque_padre

				Modificar_Apuntador(comando, path, inicio_bloque, num_bloque_apuntador3, lista_apt_apuntadores)

				return
			} else {
				//crear un nuevo espacio
				//no hay más apuntadores llenos y hay que crear un bloque nuevo para almacenar la carpeta
				var contenido_capetas []interface{}
				carpeta_nueva := structures.Content{B_name: NameArchivosByte(nombre_carpeta), B_inodo: numero_inodo_disponible}
				contenido_capetas = append(contenido_capetas, carpeta_nueva)
				contenido_capetas = append(contenido_capetas, contenido_vacio)
				contenido_capetas = append(contenido_capetas, contenido_vacio)
				contenido_capetas = append(contenido_capetas, contenido_vacio)

				nuevo_bloque_padre, enbp := Obtener_Bloque_Disponible(comando, path, inicio_bitmap_bloque, numero_bloques_total)
				if !enbp {
					return
				}
				Crear_Bloque_Carpeta_Vacio(comando, path, inicio_bloque, nuevo_bloque_padre)
				Modificar_Carpeta(comando, path, inicio_bloque, nuevo_bloque_padre, contenido_capetas)

				//modificar apuntador del inodo
				apuntadores_padre[contador] = nuevo_bloque_padre
				inodo_padre.I_block = apuntadores_padre
				fecha := ObFechaInt()
				inodo_padre.I_mtime = fecha
				Guardar_Inodo(comando, path, inicio_inodo, inodo_padre, numero_inodo_padre)
				return
			}
		}
		contador += 1
	}
	fmt.Println("[" + comando + "]: No hay espacio")
	Modificar_Bitmap(comando, path, inicio_bitmap_bloque, numero_bloque_disponible, 0)
	Modificar_Bitmap(comando, path, inicio_bitmap_inodo, numero_inodo_disponible, 0)
}
