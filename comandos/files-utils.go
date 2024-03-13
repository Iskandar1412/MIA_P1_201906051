package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Encontrar_Carpeta_En_Inodo(comando string, path string, inicio_inodo int32, inicio_bloques int32, numero_inodo int32, carpeta_archivo string) (int32, bool) {
	inodo, er := Obtener_Inodo(comando, path, inicio_inodo, numero_inodo)
	if !er {
		return 0, false
	}
	for i := range inodo.I_block {
		apuntador := inodo.I_block[i]
		if apuntador != -1 {
			if i == 13 {
				//apuntador indirecto simple
				bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, inicio_bloques, apuntador)
				if !eba1 {
					return 0, false
				}
				lista_apt1 := bloque_apuntador1.B_pointers
				for a := range lista_apt1 {
					apt1 := lista_apt1[a]
					if apt1 != -1 {
						//obtencion de bloque
						// bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, inicio)
						bloque, ebblo := Obtener_Bloque(comando, path, inicio_bloques, apt1)
						if !ebblo {
							return 0, false
						}
						lista_bloque := bloque.B_content
						for datos := range lista_bloque {
							nombre := ToString(lista_bloque[datos].B_name[:])
							if nombre == ToString([]byte(carpeta_archivo)) {
								return lista_bloque[datos].B_inodo, true
							}
						}
					}
				}
				continue
			} else if i == 14 {
				//apuntador indirecto doble
				bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, inicio_bloques, apuntador)
				if !eba1 {
					return 0, false
				}
				lista_apt1 := bloque_apuntador1.B_pointers
				for a := range lista_apt1 {
					apt1 := lista_apt1[a]
					if apt1 != -1 {
						if a == 0 {
							bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, inicio_bloques, apt1)
							if !eba2 {
								return 0, false
							}
							lista_apt2 := bloque_apuntador2.B_pointers
							for b := range lista_apt2 {
								apt2 := lista_apt2[b]
								if apt2 != -1 {
									bloque, eblo := Obtener_Bloque(comando, path, inicio_bloques, apt2)
									if !eblo {
										return 0, false
									}
									lista_bloques := bloque.B_content
									for _, datos := range lista_bloques {
										nombre := ToString(datos.B_name[:])
										if nombre == ToString([]byte(carpeta_archivo)) {
											return datos.B_inodo, true
										}
									}
								}
							}
							continue
						}
					}
					bloque, eblo := Obtener_Bloque(comando, path, inicio_bloques, apt1)
					if !eblo {
						return 0, false
					}
					lista_bloque := bloque.B_content
					for _, datos := range lista_bloque {
						nombre := ToString(datos.B_name[:])
						if nombre == ToString([]byte(carpeta_archivo)) {
							return datos.B_inodo, true
						}
					}
				}
				continue
			} else if i == 15 {
				//apuntador indirecto 13
				bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, inicio_bloques, apuntador)
				if !eba1 {
					return 0, false
				}
				lista_apt1 := bloque_apuntador1.B_pointers
				for a := range lista_apt1 {
					apt1 := lista_apt1[a]
					if apt1 != -1 {
						if a == 0 {
							bloque_apuntador2, eblo2 := Obtener_Bloque_Apuntador(comando, path, inicio_bloques, apt1)
							if !eblo2 {
								return 0, false
							}
							lista_apt2 := bloque_apuntador2.B_pointers
							for b := range lista_apt2 {
								apt2 := lista_apt2[b]
								if apt2 != -1 {
									if b == 0 {
										bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, inicio_bloques, apt2)
										if !eba3 {
											return 0, false
										}
										lista_apt3 := bloque_apuntador3.B_pointers
										for c := range lista_apt3 {
											apt3 := lista_apt3[c]
											if apt3 != -1 {
												bloque, eblo := Obtener_Bloque(comando, path, inicio_bloques, apt3)
												if !eblo {
													return 0, false
												}
												lista_bloque := bloque.B_content
												for _, datos := range lista_bloque {
													nombre := ToString(datos.B_name[:])
													if nombre == ToString([]byte(carpeta_archivo)) {
														return datos.B_inodo, true
													}
												}
											}
										}
										continue
									}
								}
								bloque, eblo := Obtener_Bloque(comando, path, inicio_bloques, apt2)
								if !eblo {
									return 0, false
								}
								lista_bloque := bloque.B_content
								for _, datos := range lista_bloque {
									nombre := ToString(datos.B_name[:])
									if nombre == ToString([]byte(carpeta_archivo)) {
										return datos.B_inodo, true
									}
								}
							}
							continue
						}
						bloque, eblo := Obtener_Bloque(comando, path, inicio_bloques, apt1)
						if !eblo {
							return 0, false
						}
						lista_bloque := bloque.B_content
						for _, datos := range lista_bloque {
							nombre := ToString(datos.B_name[:])
							if nombre == ToString([]byte(carpeta_archivo)) {
								return datos.B_inodo, true
							}
						}
					}
				}
				continue
			} else {
				bloque, eblo := Obtener_Bloque(comando, path, inicio_bloques, apuntador)
				if !eblo {
					return 0, false
				}
				lista_bloque := bloque.B_content
				for _, datos := range lista_bloque {
					nombre := ToString(datos.B_name[:])
					if nombre == ToString([]byte(carpeta_archivo)) {
						return datos.B_inodo, true
					}
				}
			}
		}
	}

	return 0, false
}

// Encontrar ruta
func Encontrar_Ruta(comando string, path string, inicio_inodos int32, inicio_bloques int32, ruta string) (int32, bool) {
	ruta_separada := strings.Split(ruta, "/")
	inodo_ultima_carpeta := int32(0)
	for carpeta := range ruta_separada {
		if string(ruta_separada[carpeta]) == "" {
			inodo_ultima_carpeta = 0
		} else {
			if inodo_ultima_carpeta == -1 {
				color.Red("[" + comando + "]: No se encontro la ruta")
				return -1, false
			}
			var er bool
			inodo_ultima_carpeta, er = Encontrar_Carpeta_En_Inodo(comando, path, inicio_inodos, inicio_bloques, inodo_ultima_carpeta, ruta_separada[carpeta])
			if !er {
				return -1, false
			}
		}
		// return
	}
	return inodo_ultima_carpeta, true
}

func Modificar_Carpeta(comando string, path string, inicio_bloque int32, num_bloque int32, contenido []interface{}) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return
	}
	defer file.Close()
	carpet := structures.BlockCarpetas{}
	size_carpet := size.SizeBlockCarpetas()
	ubicacion_bloque := inicio_bloque + (num_bloque * size_carpet)
	contador := 0
	for _, carpetas := range contenido {
		if carpeta, ok := carpetas.(structures.Content); ok {
			carpet.B_content[contador].B_inodo = carpeta.B_inodo
			carpet.B_content[contador].B_name = carpeta.B_name
		}
		contador += 1
	}

	if _, err := file.Seek(int64(ubicacion_bloque), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &carpet); err != nil {
		color.Red("[" + comando + "]: Error en escritura de bloque")
		return
	}
	color.Green("[" + comando + "]: Carpeta modificada exitosamente")
}

func Recalcular_Size_Carpetas(comando string, id string) {
	conjunto, path, er := Obtener_Particion_ID(id)
	if !er {
		return
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
				return
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
				return
			}
			if particion.Part_type == 'E' {
				color.Red("[MKFS]: No se puede obtener información de particion extendida")
				return
			}
		}
	}

	inodo, einodo := Obtener_Inodo(comando, path, superbloque.S_inode_start, 0)
	if !einodo {
		return
	}
	size_total, est := Recalcular_Size_Carpetas_Recursivo(comando, id, 0)
	if !est {
		return
	}
	inodo.I_s = size_total
	Guardar_Inodo(comando, path, superbloque.S_inode_start, inodo, 0)
}

func Crear_Ruta_Carpetas(comando string, id string, ruta string, uid_user string, uid_grupo string) bool {
	conjunto, path, ec := Obtener_Particion_ID(id)
	if !ec {
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

	//separar ruta
	ruta_separada := strings.Split(ruta, "/")
	ruta_anterior := "/"
	nueva_ruta := ""
	for _, carpeta := range ruta_separada {
		if carpeta == "" {
			continue
		}
		ruta_anterior = nueva_ruta
		if ruta_anterior == "" {
			ruta_anterior = "/"
		}
		nueva_ruta += "/" + carpeta
		ultimo_inodo_padre, _ := Encontrar_Ruta(comando, path, superbloque.S_inode_start, superbloque.S_block_start, nueva_ruta)
		if ultimo_inodo_padre != -1 {
			color.Red("[" + comando + "]: La ruta ya existe")
		} else {
			color.Blue("No se encontró carpeta... " + carpeta + " [Creando...] ")
			iuid, _ := strconv.Atoi(uid_user)
			iuig, _ := strconv.Atoi(uid_grupo)
			//crear_estructura_carpeta_nueva
			if !Crear_Estructura_Carpeta_Nueva(comando, path, superbloque.S_inode_start, superbloque.S_block_start, superbloque.S_inodes_count, superbloque.S_blocks_count, ruta_anterior, carpeta, superbloque.S_bm_inode_start, superbloque.S_bm_block_start, int32(iuid), int32(iuig)) {
				// color.Red("Error al crear carpeta")
				return false
			}
			// color.Blue("Creando Archivo")
		}
	}
	return true
}
