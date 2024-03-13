package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

func Modificar_Archivo(comando string, path string, inicio_bloques int32, no_bloque int32, content structures.BloqueArchivos) {
	// struct_archivo := structures.BloqueArchivos{}
	ubicacion_bloque := inicio_bloques + (no_bloque * size.SizeBloqueArchivos())

	//grabar contenido
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return
	}
	defer file.Close()
	//bytes_con_inodos := inicio_inodos + numero
	if _, err := file.Seek(int64(ubicacion_bloque), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return
	}
	if err := binary.Write(file, binary.LittleEndian, &content); err != nil {
		color.Red("[" + comando + "]: Error en escritura de bloque")
		return
	}
	color.Green("[" + comando + "]: Archivo modificado")
}

// Creacion de archivo vacíos
func Crear_Archivo_Vacio(comando string, path string, ruta string, superbloque structures.SuperBlock, contenido string, nombre_archivo string, id_user int32, id_group int32) bool {
	numero_inodo_disponible, existe_inodo := Obtener_Inodo_Disponible(comando, path, superbloque.S_bm_inode_start, superbloque.S_inodes_count)
	if !existe_inodo {
		color.Red("[" + comando + "]: Error al buscar inodo")
		return false
	}
	// fmt.Println(numero_inodo_disponible) //se usará luego

	//información padre
	num_inodo_padre, existe_padre := Encontrar_Ruta(comando, path, superbloque.S_inode_start, superbloque.S_block_start, ruta)
	if !existe_padre {
		return false
	}
	//inodo padre
	inodo_padre, existe_ipadre := Obtener_Inodo(comando, path, superbloque.S_inode_start, num_inodo_padre)
	if !existe_ipadre {
		return false
	}

	validacion := Validar_Permisos(comando, inodo_padre.I_uid, inodo_padre.I_gid, id_user, id_group, inodo_padre.I_perm, int32(2))
	//fmt.Println(inodo_padre)
	if !validacion {
		color.Magenta("No se tienen los permisos necesarios")
		return false
	}

	//apuntadores vacios para inodo de archivo
	//lista_apt_bloques
	bloqueApuntadores := structures.BloqueApuntadores{}
	for i := range bloqueApuntadores.B_pointers {
		bloqueApuntadores.B_pointers[i] = int32(-1)
	}

	numember_block_disp := int32(0)

	// numero_bloque_disponible := int32(0)
	if len(contenido) <= 64 {
		// bytes_archivo := []byte(contenido)
		archivo := structures.BloqueArchivos{B_content: DevolverContenidoArchivo(contenido)}
		// fmt.Println(archivo)
		numero_bloque_disponible, existe_bloque := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
		if !existe_bloque {
			color.Red("[" + comando + "]: Erorr al buscar bloque")
			return false
		}
		numember_block_disp = numero_bloque_disponible
		Crear_Bloque_Archivo_Vacio(comando, path, superbloque.S_block_start, numero_bloque_disponible)
		Modificar_Archivo(comando, path, superbloque.S_block_start, numero_bloque_disponible, archivo)

		//modificar bloque apuntadores
		bloqueApuntadores.B_pointers[0] = numero_bloque_disponible
		Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
	} else {
		cadena := ""
		contador := 0
		color.Yellow("[" + comando + "]: Creando archivo")
		for i := 0; i < len(contenido); i++ {
			if len(cadena) == 64 || i == (len(contenido)-1) {
				if i == (len(contenido) - 1) {
					cadena += string(contenido[i])
				}
				//creacion bloque de archivos
				archivo := structures.BloqueArchivos{B_content: DevolverContenidoArchivo(cadena)}
				numero_bloque_disponible, existe := Obtener_Bloque_Disponible(comando, path, superbloque.S_block_start, superbloque.S_blocks_count)
				if !existe {
					color.Red("[" + comando + "]: Erorr al buscar bloque")
					return false
				}
				numember_block_disp = numero_bloque_disponible
				Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
				Crear_Bloque_Archivo_Vacio(comando, path, superbloque.S_block_start, numero_bloque_disponible)
				Modificar_Archivo(comando, path, superbloque.S_block_start, numero_bloque_disponible, archivo)

				//modificar array de apuntadores del inodo
				cadena = string(contenido[i])
				if contador == 13 {
					if bloqueApuntadores.B_pointers[contador] != int32(-1) { //caso que ya tiene un apuntador apuntando a algo
						bloque_apuntador1, esta := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, bloqueApuntadores.B_pointers[contador])
						if !esta {
							return false
						}
						// bloqueApuntadores = bloque_apuntador1
						lista_apt1 := bloque_apuntador1.B_pointers
						for a := range lista_apt1 {
							apt1 := lista_apt1[a]
							if apt1 == -1 {
								//accion
								lista_apt1[a] = numero_bloque_disponible
								lista_a := structures.BloqueApuntadores{B_pointers: lista_apt1}
								Modificar_Apuntador(comando, path, superbloque.S_block_start, bloqueApuntadores.B_pointers[contador], lista_a)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
								if a == 15 {
									contador += 1
								}
								break
							}
						}
					} else {
						//apuntador indirecto simple
						num_bloque_apuntador, eba := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
						if !eba {
							return false
						}
						Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)
						//apuntar a carpeta
						lista_apt_apuntadores := structures.BloqueApuntadores{}
						for i := range lista_apt_apuntadores.B_pointers {
							lista_apt_apuntadores.B_pointers[i] = int32(-1)
						}
						lista_apt_apuntadores.B_pointers[contador] = num_bloque_apuntador
						//guardar apuntador 13 en carpeta
						bloqueApuntadores.B_pointers[contador] = num_bloque_apuntador
						lista_a := structures.BloqueApuntadores{B_pointers: bloqueApuntadores.B_pointers}
						Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador, lista_a)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
					}
					continue
				} else if contador == 14 {
					if bloqueApuntadores.B_pointers[contador] != -1 {
						bloque_apuntador1, eba := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, bloqueApuntadores.B_pointers[contador])
						if !eba {
							return false
						}
						lista_apt1 := bloque_apuntador1.B_pointers
						for a := range lista_apt1 {
							apt1 := lista_apt1[a]
							agrego := false
							apt2lleno := false
							if a == 0 && !apt2lleno {
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, lista_apt1[a])
								if !eba2 {
									return false
								}
								lista_apt2 := bloque_apuntador2.B_pointers
								//lista apt2 = bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if apt2 == -1 {
										lista_apt2[b] = numero_bloque_disponible
										lista_a := structures.BloqueApuntadores{B_pointers: lista_apt2}

										Modificar_Apuntador(comando, path, superbloque.S_block_start, apt1, lista_a)
										Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
										cadena = string(contenido[i])
										agrego = true
										break
									}
								}
							}
							if a == 0 && !agrego {
								apt2lleno = true
								continue
							}
							if apt1 == -1 && !agrego {
								lista_apt1[a] = numero_bloque_disponible
								lista_a := structures.BloqueApuntadores{B_pointers: lista_apt1}
								Modificar_Apuntador(comando, path, superbloque.S_block_start, bloqueApuntadores.B_pointers[contador], lista_a)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
								if a == 15 {
									contador += 1
								}
								cadena = string(contenido[i])
								break
							} else if agrego { //al guardar sale del for
								break
							}
						}
					} else {

						//apuntador 1
						num_bloque_apuntador, nba := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
						if !nba {
							return false
						}
						Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)

						//creacion apuntador 2
						num_bloque_apuntador2, nba2 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
						if !nba2 {
							return false
						}
						Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador2)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador2, 1)

						//apuntar a carpeta
						lista_apt_apuntadores2 := structures.BloqueApuntadores{}
						for j := range lista_apt_apuntadores2.B_pointers {
							lista_apt_apuntadores2.B_pointers[j] = int32(-1)
						}
						lista_apt_apuntadores2.B_pointers[0] = numero_bloque_disponible

						//apuntar carpeta
						lista_apt_apuntadores := structures.BloqueApuntadores{}
						for j := range lista_apt_apuntadores.B_pointers {
							lista_apt_apuntadores.B_pointers[j] = int32(-1)
						}
						lista_apt_apuntadores.B_pointers[0] = num_bloque_apuntador2

						//guardar apuntador 13 del inodo
						bloqueApuntadores.B_pointers[contador] = num_bloque_apuntador

						// lista_a := listaAPTApuntadores2
						Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador2, lista_apt_apuntadores2)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)

						//lista_a := istaaptapuntadores
						Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador, lista_apt_apuntadores)
					}
				} else if contador == 15 {
					if bloqueApuntadores.B_pointers[contador] != -1 {
						//caso de apuntador ya creado
						bloque_apuntador_1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, bloqueApuntadores.B_pointers[contador])
						if !eba1 {
							return false
						}
						lista_apt1 := bloque_apuntador_1.B_pointers
						for a := range lista_apt1 {
							apt1 := lista_apt1[a]
							agrego := false
							apt2lleno := false
							apt3lleno := false
							if a == 0 && !apt2lleno {
								bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
								if !eba2 {
									return false
								}
								lista_apt2 := bloque_apuntador2.B_pointers
								for b := range lista_apt2 {
									apt2 := lista_apt2[b]
									if b == 0 && !apt3lleno {
										bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt2)
										if !eba3 {
											return false
										}
										lista_apt3 := bloque_apuntador3.B_pointers
										for c := range lista_apt3 {
											apt3 := lista_apt3[c]
											if apt3 == -1 {
												lista_apt3[c] = numero_bloque_disponible
												lista_a := structures.BloqueApuntadores{B_pointers: lista_apt3}
												Modificar_Apuntador(comando, path, superbloque.S_block_start, apt2, lista_a)
												Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
												cadena = string(contenido[i])
												agrego = true
												break
											}
										}
									}
									if b == 0 && !agrego {
										apt3lleno = true
										continue
									}
									if apt2 == -1 {
										lista_apt2[b] = numero_bloque_disponible
										lista_a := structures.BloqueApuntadores{B_pointers: lista_apt2}
										Modificar_Apuntador(comando, path, superbloque.S_block_start, apt1, lista_a)
										Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
										cadena = string(contenido[i])
										agrego = true
										break
									} else if agrego {
										break
									}
								}
							}
							if a == 0 && !agrego {
								apt2lleno = true
								break
							}
							if apt1 == 1 && !agrego {
								lista_apt1[a] = numero_bloque_disponible
								list_a := structures.BloqueApuntadores{B_pointers: lista_apt1}
								Modificar_Apuntador(comando, path, superbloque.S_block_start, bloqueApuntadores.B_pointers[contador], list_a)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
								if a == 15 {
									contador += 1
								}
								cadena = string(contenido[i])
								break
							} else if agrego {
								break
							}
						}
					} else {
						//apuntador 1
						num_bloque_apuntador, enba := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
						if !enba {
							return false
						}
						Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)
						//apuntador2

						num_bloque_apuntador2, eba2 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
						if !eba2 {
							return false
						}
						Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador2)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador2, 1)

						//apuntador 3
						num_bloque_apuntador3, eba3 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
						if !eba3 {
							return false
						}
						Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador3)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador3, 1)

						//apuntar a carpeta
						lista_apt_apuntadores3 := structures.BloqueApuntadores{}
						for j := range lista_apt_apuntadores3.B_pointers {
							lista_apt_apuntadores3.B_pointers[j] = int32(-1)
						}
						lista_apt_apuntadores3.B_pointers[0] = numero_bloque_disponible

						lista_apt_apuntadores2 := structures.BloqueApuntadores{}
						for j := range lista_apt_apuntadores2.B_pointers {
							lista_apt_apuntadores2.B_pointers[j] = int32(-1)
						}
						lista_apt_apuntadores2.B_pointers[0] = num_bloque_apuntador3

						lista_apt_apuntadores := structures.BloqueApuntadores{}
						for j := range lista_apt_apuntadores.B_pointers {
							lista_apt_apuntadores.B_pointers[j] = int32(-1)
						}
						lista_apt_apuntadores.B_pointers[0] = num_bloque_apuntador2

						//guardar apuntador 13 del inodo
						lista_apt_apuntadores.B_pointers[contador] = num_bloque_apuntador

						// lista_a := structures.BloqueApuntadores{B_pointers: lista_apt_apuntadores3.B_pointers}
						Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador3, lista_apt_apuntadores3)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)

						Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador2, lista_apt_apuntadores2)

						Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador, lista_apt_apuntadores)

						continue
					}
				} else {
					if contador == 16 {
						color.Red("No hay más espacio en el inodo")
						break
					}
					bloqueApuntadores.B_pointers[contador] = numero_bloque_disponible
					Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
					contador += 1
				}
			} else {
				cadena += string(contenido[i])
			}
		}
	}
	fecha := ObFechaInt()
	size_conenido := len(string(contenido))
	inodo := structures.Inode{}
	inodo.I_uid = id_user
	inodo.I_gid = id_group
	inodo.I_s = int32(size_conenido)
	inodo.I_atime = fecha
	inodo.I_ctime = fecha
	inodo.I_mtime = fecha
	inodo.I_block = bloqueApuntadores.B_pointers
	inodo.I_type = int32(1)
	inodo.I_perm = int32(664)

	if !Guardar_Inodo(comando, path, superbloque.S_inode_start, inodo, numero_inodo_disponible) {
		color.Red("No se pudo guardar el inodo")
		return false
	}

	Modificar_Bitmap(comando, path, superbloque.S_bm_inode_start, numero_inodo_disponible, 1)
	return Agregar_Bloque_Lista_Inodos(comando, path, superbloque.S_block_start, numero_inodo_disponible, superbloque.S_bm_block_start, superbloque.S_blocks_count, nombre_archivo, inodo_padre, superbloque.S_inode_start, num_inodo_padre, superbloque.S_bm_inode_start, numember_block_disp)

}

func Modificar_Contenido_Archivo(comando string, id string, ruta string, contenido_nuevo string) {
	numero_len := len(contenido_nuevo)
	numero_bloque_crear_nuevo := float64(numero_len) / float64(64)
	if numero_bloque_crear_nuevo > 106 {
		color.Red("Contenido exede límite de caracteres")
		return
	}
	conjunto, path, ec := Obtener_Particion_ID(id)
	if !ec {
		return
	}
	superbloque := structures.SuperBlock{}

	ebr := structures.EBR{}
	if temp, ok := conjunto[2].(structures.EBR); ok {
		v := reflect.ValueOf(temp)
		reflect.ValueOf(&ebr).Elem().Set(v)
		// byte_escribir = ebr.Part_start
		var er bool
		superbloque, er = Obtener_Superbloque(comando, path, ToString(ebr.Name[:]))
		if !er {
			return
		}
		conjunto[0] = nil
		conjunto[1] = nil
	}

	particion := structures.Partition{}
	if temp, ok := conjunto[0].(structures.Partition); ok {
		v := reflect.ValueOf(temp)
		reflect.ValueOf(&particion).Elem().Set(v)
		var er bool
		superbloque, er = Obtener_Superbloque(comando, path, ToString(particion.Part_name[:]))
		if !er {
			return
		}
		// byte_escribir = particion.Part_start
	}

	ruta_separada := strings.Split(ruta, "/")
	cantidad_carpetas := len(ruta_separada)
	nombre_archivo := ruta_separada[cantidad_carpetas-1]
	ruta_sin_archivo := strings.ReplaceAll(ruta, "/"+nombre_archivo, "")

	numero_inodo_carpeta, enic := Encontrar_Ruta(comando, path, superbloque.S_inode_start, superbloque.S_block_start, ruta_sin_archivo)
	if !enic {
		return
	}

	inodo_carpeta, eic := Obtener_Inodo(comando, path, superbloque.S_inode_start, numero_inodo_carpeta)
	if !eic {
		return
	}

	lista_apuntadores := inodo_carpeta.I_block
	numero_inodo_archivo := int32(-1)
	for i := range lista_apuntadores {
		// fmt.Println(i, lista_apuntadores)
		apuntador := lista_apuntadores[i]
		if numero_inodo_archivo != -1 {
			break
		}
		if apuntador != -1 {
			if i == 13 {
				bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apuntador)
				if !eba1 {
					return
				}
				lista_apt1 := bloque_apuntador1.B_pointers
				encontro := false
				for a := range lista_apt1 {
					apt1 := lista_apt1[a]
					if apt1 != -1 {
						bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
						if !ebc {
							return
						}
						contenido_carpeta := bloque_carpeta.B_content
						for _, contenido := range contenido_carpeta {
							if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archivo)) {
								numero_inodo_archivo = contenido.B_inodo
								encontro = true
								break
							}
						}
						if encontro {
							break
						}
					}
				}
				continue
			} else if i == 14 {
				bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apuntador)
				if !eba1 {
					return
				}
				lista_apt1 := bloque_apuntador1.B_pointers
				encontro := false
				for a := range lista_apt1 {
					apt1 := lista_apt1[a]
					if a == 0 {
						bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
						if !eba2 {
							return
						}
						lista_apt2 := bloque_apuntador2.B_pointers
						for b := range lista_apt2 {
							apt2 := lista_apt2[b]
							if apt2 != -1 {
								bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
								if !ebc {
									return
								}
								contenido_carpeta := bloque_carpeta.B_content
								for _, contenido := range contenido_carpeta {
									if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archivo)) {
										numero_inodo_archivo = contenido.B_inodo
										encontro = true
										break
									}
								}
								if encontro {
									break
								}
							}
						}
						if encontro {
							break
						}
						continue
					}
					if apt1 != -1 {
						bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
						if !ebc {
							return
						}
						contenido_carpeta := bloque_carpeta.B_content
						for _, contenido := range contenido_carpeta {
							if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archivo)) {
								numero_inodo_archivo = contenido.B_inodo
								encontro = true
								break
							}
						}
						if encontro {
							break
						}
					}
				}
				continue
			} else if i == 15 {
				bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apuntador)
				if !eba1 {
					return
				}
				lista_apt1 := bloque_apuntador1.B_pointers
				encontro := false
				for a := range lista_apt1 {
					apt1 := lista_apt1[a]
					if a == 0 {
						bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
						if !eba2 {
							return
						}
						lista_apt2 := bloque_apuntador2.B_pointers
						for b := range lista_apt2 {
							apt2 := lista_apt2[b]
							if b == 0 {
								bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt2)
								if !eba3 {
									return
								}
								lista_apt3 := bloque_apuntador3.B_pointers
								for c := range lista_apt3 {
									apt3 := lista_apt3[c]
									if apt3 != -1 {
										bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt3)
										if !ebc {
											return
										}
										contenido_carpeta := bloque_carpeta.B_content
										for _, contenido := range contenido_carpeta {
											if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archivo)) {
												numero_inodo_archivo = contenido.B_inodo
												encontro = true
												break
											}
										}
										if encontro {
											break
										}
									}
								}
								if encontro {
									break
								}
								continue
							}
							if apt2 != -1 {
								bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt2)
								if !ebc {
									return
								}
								contenido_carpeta := bloque_carpeta.B_content
								for _, contenido := range contenido_carpeta {
									if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archivo)) {
										numero_inodo_archivo = contenido.B_inodo
										encontro = true
										break
									}
								}
								if encontro {
									break
								}
							}

						}
						continue
					}
					if apt1 != -1 {
						bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apt1)
						if !ebc {
							return
						}
						contenido_carpeta := bloque_carpeta.B_content
						for _, contenido := range contenido_carpeta {
							if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archivo)) {
								numero_inodo_archivo = contenido.B_inodo
								encontro = true
								break
							}
						}
						if encontro {
							break
						}
					}
				}
				continue
			} else {
				bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apuntador)
				if !ebc {
					return
				}
				contenido_carpeta := bloque_carpeta.B_content
				for _, contenido := range contenido_carpeta {
					if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archivo)) {
						numero_inodo_archivo = contenido.B_inodo
						// encontro = true
						break
					}
				}
			}
		}
	}
	if numero_inodo_archivo != -1 {
		//loquesea
		inodo_archivo, eia := Obtener_Inodo(comando, path, superbloque.S_inode_start, numero_inodo_archivo)
		if !eia {
			return
		}
		apuntadores_archivo := inodo_archivo.I_block
		contador := 0
		for i := range apuntadores_archivo {
			parte_archivo := apuntadores_archivo[i]
			if parte_archivo != -1 {
				if i == 13 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, parte_archivo)
					if !eba1 {
						return
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, apt1, 0)
					}
					var lista [16]int32
					for f := range lista {
						lista[f] = int32(-1)
					}
					lista_a := structures.BloqueApuntadores{}
					lista_a.B_pointers = lista
					Modificar_Apuntador(comando, path, superbloque.S_block_start, parte_archivo, lista_a)
					Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, parte_archivo, 0)
					continue
				} else if i == 14 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, parte_archivo)
					if !eba1 {
						return
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if a == 0 {
							bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
							if !eba2 {
								return
							}
							lista_apt2 := bloque_apuntador2.B_pointers
							for b := range lista_apt2 {
								apt2 := lista_apt2[b]
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, apt2, 0)
							}
							var lista [16]int32
							for f := range lista {
								lista[f] = int32(-1)
							}
							lista_a := structures.BloqueApuntadores{}
							lista_a.B_pointers = lista
							Modificar_Apuntador(comando, path, superbloque.S_block_start, apt1, lista_a)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, apt1, 0)
							continue
						}
						lista_apt1[a] = int32(-1)
					}
					lista_a := structures.BloqueApuntadores{}
					lista_a.B_pointers = lista_apt1
					Modificar_Apuntador(comando, path, superbloque.S_block_start, parte_archivo, lista_a)
					Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, parte_archivo, 0)
					continue
				} else if i == 15 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, parte_archivo)
					if !eba1 {
						return
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if a == 0 {
							bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
							if !eba2 {
								return
							}
							lista_apt2 := bloque_apuntador2.B_pointers
							for b := range lista_apt2 {
								apt2 := lista_apt2[b]
								if b == 0 {
									bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt2)
									if !eba3 {
										return
									}
									lista_apt3 := bloque_apuntador3.B_pointers
									for c := range lista_apt3 {
										apt3 := lista_apt3[c]
										Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, apt3, 0)
									}
									var lista [16]int32
									for f := range lista {
										lista[f] = int32(-1)
									}
									lista_a := structures.BloqueApuntadores{}
									lista_a.B_pointers = lista
									Modificar_Apuntador(comando, path, superbloque.S_block_start, apt2, lista_a)
									Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, apt2, 0)
								}
								lista_apt2[b] = int32(-1)
							}
							lista_a := structures.BloqueApuntadores{}
							lista_a.B_pointers = lista_apt2
							Modificar_Apuntador(comando, path, superbloque.S_block_start, apt1, lista_a)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, apt1, 0)
							continue
						}
						lista_apt1[a] = int32(-1)
					}
					lista_a := structures.BloqueApuntadores{}
					lista_a.B_pointers = lista_apt1
					Modificar_Apuntador(comando, path, superbloque.S_block_start, parte_archivo, lista_a)
					Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, parte_archivo, 0)
					continue
				} else {
					Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, parte_archivo, 0)
					inodo_archivo.I_block[contador] = int32(-1)
				}
			}
			contador += 1
		}
		//creacion de los bloques de archivo
		var lista_apt_bloques [16]int32
		for i := range lista_apt_bloques {
			lista_apt_bloques[i] = int32(-1)
		}
		cadena := ""
		if len(contenido_nuevo) <= 64 {
			archivo := structures.BloqueArchivos{B_content: DevolverContenidoArchivo(contenido_nuevo)}
			numero_bloque_disponible, enbd := Obtener_Bloque_Disponible(comando, path, superbloque.S_block_start, superbloque.S_blocks_count)
			if !enbd {
				return
			}
			Crear_Bloque_Archivo_Vacio(comando, path, superbloque.S_block_start, numero_bloque_disponible)
			Modificar_Archivo(comando, path, superbloque.S_block_start, numero_bloque_disponible, archivo)
			lista_apt_bloques[0] = numero_bloque_disponible
			Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
		} else {
			cadena = ""
			contador := 0
			size_cadena_anterior := 0
			agrego_bloque_nuevo_a := false
			for i := 0; i < len(contenido_nuevo); i++ {
				if len(cadena) == 64 || i == len(contenido_nuevo)-1 {
					if (i == (len(contenido_nuevo))-1) && (len(cadena) != 64) {
						cadena += string(contenido_nuevo[i])
						size_cadena_anterior += 1
						agrego_bloque_nuevo_a = true
					}
					archivo := structures.BloqueArchivos{B_content: DevolverContenidoArchivo(cadena)}
					numero_bloque_disponible, enbd := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
					if !enbd {
						return
					}
					Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
					Crear_Bloque_Archivo_Vacio(comando, path, superbloque.S_block_start, numero_bloque_disponible)
					Modificar_Archivo(comando, path, superbloque.S_block_start, numero_bloque_disponible, archivo)
					cadena = string(contenido_nuevo[i])
					if contador == 13 {
						if lista_apt_bloques[contador] != -1 {
							bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador])
							if !eba1 {
								return
							}
							lista_apt1 := bloque_apuntador1.B_pointers
							for a := range lista_apt1 {
								apt1 := lista_apt1[a]
								if apt1 == -1 {
									lista_apt1[a] = numero_bloque_disponible
									lista_a := structures.BloqueApuntadores{B_pointers: lista_apt1}
									Modificar_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador], lista_a)
									Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
									if a == 15 {
										contador += 1
									}
									break
								}
							}
						} else {
							num_bloque_apuntador, enba := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
							if !enba {
								return
							}
							Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)
							var lista_apt_apuntadores [16]int32
							for f := range lista_apt_apuntadores {
								lista_apt_apuntadores[f] = int32(-1)
							}
							lista_apt_apuntadores[0] = numero_bloque_disponible
							lista_apt_bloques[contador] = num_bloque_apuntador
							lista_a := structures.BloqueApuntadores{B_pointers: lista_apt_apuntadores}
							Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador, lista_a)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
						}
					} else if contador == 14 {
						if lista_apt_bloques[contador] != -1 {
							bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador])
							if !eba1 {
								return
							}
							lista_apt1 := bloque_apuntador1.B_pointers
							for a := range lista_apt1 {
								apt1 := lista_apt1[a]
								agrego := false
								apt2lleno := false
								if a == 0 && !apt2lleno {
									bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, lista_apt1[a])
									if !eba2 {
										return
									}
									lista_apt2 := bloque_apuntador2.B_pointers
									for b := range lista_apt2 {
										apt2 := lista_apt2[b]
										if apt2 == -1 {
											lista_apt2[b] = numero_bloque_disponible
											lista_a := structures.BloqueApuntadores{B_pointers: lista_apt2}
											Modificar_Apuntador(comando, path, superbloque.S_block_start, apt1, lista_a)
											Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
											agrego = true
											break
										}
									}
								}
								if a == 0 && !agrego {
									apt2lleno = true
									continue
								}
								if apt1 == -1 && !agrego {
									lista_apt1[a] = numero_bloque_disponible
									lista_a := structures.BloqueApuntadores{B_pointers: lista_apt1}
									Modificar_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador], lista_a)
									Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
									if a == 15 {
										contador += 1
									}
									break
								} else if agrego {
									break
								}
							}
						} else {
							//apuntador 1
							num_bloque_apuntador, enba := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
							if !enba {
								return
							}
							Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)

							//apuntador 2
							num_bloque_apuntador2, enba2 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
							if !enba2 {
								return
							}
							Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador2)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador2, 1)

							//apuntar a carpeta
							var lista_apt_apuntadores2 [16]int32
							for f := range lista_apt_apuntadores2 {
								lista_apt_apuntadores2[f] = int32(-1)
							}
							lista_apt_apuntadores2[0] = numero_bloque_disponible

							//apuntar a carpeta
							var lista_apt_apuntadores [16]int32
							for f := range lista_apt_apuntadores {
								lista_apt_apuntadores[f] = int32(-1)
							}
							lista_apt_apuntadores[0] = num_bloque_apuntador2

							//giardar en el apuntador 13 del inodo
							lista_apt_bloques[contador] = num_bloque_apuntador

							lista_a := structures.BloqueApuntadores{B_pointers: lista_apt_apuntadores2}
							Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador2, lista_a)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)

							lista_a.B_pointers = lista_apt_apuntadores
							Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador, lista_a)
						}
					} else if contador == 15 {
						if lista_apt_bloques[contador] != -1 {
							bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador])
							if !eba1 {
								return
							}
							lista_apt1 := bloque_apuntador1.B_pointers
							for a := range lista_apt1 {
								apt1 := lista_apt1[a]
								agrego := false
								apt2lleno := false
								apt3lleno := false
								if a == 0 && !apt2lleno {
									bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, lista_apt1[a])
									if !eba2 {
										return
									}
									lista_apt2 := bloque_apuntador2.B_pointers
									for b := range lista_apt2 {
										apt2 := lista_apt2[b]
										if b == 0 && !apt3lleno {
											bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt2)
											if !eba3 {
												return
											}
											lista_apt3 := bloque_apuntador3.B_pointers
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												if apt3 == -1 {
													lista_apt3[b] = numero_bloque_disponible
													lista_a := structures.BloqueApuntadores{B_pointers: lista_apt3}
													Modificar_Apuntador(comando, path, superbloque.S_block_start, apt2, lista_a)
													Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
													agrego = true
													break
												}
											}
										}
										if b == 0 && !agrego {
											apt3lleno = true
											continue
										}
										if apt2 == -1 {
											lista_apt2[b] = numero_bloque_disponible
											lista_a := structures.BloqueApuntadores{B_pointers: lista_apt2}
											Modificar_Apuntador(comando, path, superbloque.S_block_start, apt1, lista_a)
											Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
											agrego = true
											break
										} else if agrego {
											break
										}
									}
								}
								if a == 0 && !agrego {
									apt2lleno = true
									continue
								}
								if apt1 == -1 && !agrego {
									lista_apt1[a] = numero_bloque_disponible
									lista_a := structures.BloqueApuntadores{B_pointers: lista_apt1}
									Modificar_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador], lista_a)
									Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
									if a == 15 {
										contador += 1
									}
									break
								} else if agrego {
									break
								}
							}
						} else {
							//apuntador 1
							num_bloque_apuntador, enba := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
							if !enba {
								return
							}
							Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)

							//apuntador 2
							num_bloque_apuntador2, enba2 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
							if !enba2 {
								return
							}
							Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador2)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador2, 1)

							//apuntador 3
							num_bloque_apuntador3, enba3 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
							if !enba3 {
								return
							}
							Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador3)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador3, 1)

							//apuntar a carpeta
							var lista_apt_apuntadores3 [16]int32
							for f := range lista_apt_apuntadores3 {
								lista_apt_apuntadores3[f] = int32(-1)
							}
							lista_apt_apuntadores3[0] = numero_bloque_disponible

							//apuntar a carpeta
							var lista_apt_apuntadores2 [16]int32
							for f := range lista_apt_apuntadores2 {
								lista_apt_apuntadores2[f] = int32(-1)
							}
							lista_apt_apuntadores2[0] = num_bloque_apuntador3

							//apuntar a carpeta
							var lista_apt_apuntadores [16]int32
							for f := range lista_apt_apuntadores {
								lista_apt_apuntadores[f] = int32(-1)
							}
							lista_apt_apuntadores[0] = num_bloque_apuntador2

							//giardar en el apuntador 13 del inodo
							lista_apt_bloques[contador] = num_bloque_apuntador

							lista_a := structures.BloqueApuntadores{B_pointers: lista_apt_apuntadores3}
							Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador3, lista_a)
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)

							lista_a.B_pointers = lista_apt_apuntadores2
							Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador2, lista_a)

							lista_a.B_pointers = lista_apt_apuntadores
							Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador, lista_a)
							continue
						}
					} else {
						if contador == 16 {
							color.Red("[" + comando + "]: No hay más espacio")
						}
						lista_apt_bloques[contador] = numero_bloque_disponible
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)

						contador += 1
					}
					if (i == len(contenido_nuevo)-1) && (size_cadena_anterior == 64) && (!agrego_bloque_nuevo_a) {
						archivo := structures.BloqueArchivos{B_content: DevolverContenidoArchivo(cadena)}
						numero_bloque_disponible, enbd := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
						if !enbd {
							return
						}
						Crear_Bloque_Archivo_Vacio(comando, path, superbloque.S_block_start, numero_bloque_disponible)
						Modificar_Archivo(comando, path, superbloque.S_block_start, numero_bloque_disponible, archivo)

						//moodificar array de apuntadores del inodo
						if contador == 13 {
							if lista_apt_bloques[contador] != -1 {
								bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador])
								if !eba1 {
									return
								}
								lista_apt1 := bloque_apuntador1.B_pointers
								for a := range lista_apt1 {
									apt1 := lista_apt1[a]
									if apt1 == -1 {
										lista_apt1[a] = numero_bloque_disponible
										lista_a := structures.BloqueApuntadores{B_pointers: lista_apt1}
										Modificar_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador], lista_a)
										Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
										if a == 15 {
											contador += 1
										}
										break
									}
								}
							} else {
								num_bloque_apuntador, enba := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
								if !enba {
									return
								}
								Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)
								var lista_apt_apuntadores [16]int32
								for f := range lista_apt_apuntadores {
									lista_apt_apuntadores[f] = int32(-1)
								}
								lista_apt_apuntadores[0] = numero_bloque_disponible
								lista_apt_bloques[contador] = num_bloque_apuntador
								lista_a := structures.BloqueApuntadores{B_pointers: lista_apt_apuntadores}
								Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador, lista_a)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
							}
						} else if contador == 14 {
							if lista_apt_bloques[contador] != -1 {
								bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador])
								if !eba1 {
									return
								}
								lista_apt1 := bloque_apuntador1.B_pointers
								for a := range lista_apt1 {
									apt1 := lista_apt1[a]
									agrego := false
									apt2lleno := false
									if a == 0 && !apt2lleno {
										bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, lista_apt1[a])
										if !eba2 {
											return
										}
										lista_apt2 := bloque_apuntador2.B_pointers
										for b := range lista_apt2 {
											apt2 := lista_apt2[b]
											if apt2 == -1 {
												lista_apt2[b] = numero_bloque_disponible
												lista_a := structures.BloqueApuntadores{B_pointers: lista_apt2}
												Modificar_Apuntador(comando, path, superbloque.S_block_start, apt1, lista_a)
												Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
												agrego = true
												break
											}
										}
									}
									if a == 0 && !agrego {
										apt2lleno = true
										continue
									}
									if apt1 == -1 && !agrego {
										lista_apt1[a] = numero_bloque_disponible
										lista_a := structures.BloqueApuntadores{B_pointers: lista_apt1}
										Modificar_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador], lista_a)
										Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
										if a == 15 {
											contador += 1
										}
										break
									} else if agrego {
										break
									}
								}
							} else {
								//apuntador 1
								num_bloque_apuntador, enba := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
								if !enba {
									return
								}
								Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)

								//apuntador 2
								num_bloque_apuntador2, enba2 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
								if !enba2 {
									return
								}
								Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador2)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador2, 1)

								//apuntar a carpeta
								var lista_apt_apuntadores2 [16]int32
								for f := range lista_apt_apuntadores2 {
									lista_apt_apuntadores2[f] = int32(-1)
								}
								lista_apt_apuntadores2[0] = numero_bloque_disponible

								//apuntar a carpeta
								var lista_apt_apuntadores [16]int32
								for f := range lista_apt_apuntadores {
									lista_apt_apuntadores[f] = int32(-1)
								}
								lista_apt_apuntadores[0] = num_bloque_apuntador2

								//giardar en el apuntador 13 del inodo
								lista_apt_bloques[contador] = num_bloque_apuntador

								lista_a := structures.BloqueApuntadores{B_pointers: lista_apt_apuntadores2}
								Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador2, lista_a)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)

								lista_a.B_pointers = lista_apt_apuntadores
								Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador, lista_a)
							}
						} else if contador == 15 {
							if lista_apt_bloques[contador] != -1 {
								bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador])
								if !eba1 {
									return
								}
								lista_apt1 := bloque_apuntador1.B_pointers
								for a := range lista_apt1 {
									apt1 := lista_apt1[a]
									agrego := false
									apt2lleno := false
									apt3lleno := false
									if a == 0 && !apt2lleno {
										bloque_apuntador2, eba2 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt1)
										if !eba2 {
											return
										}
										lista_apt2 := bloque_apuntador2.B_pointers
										for b := range lista_apt2 {
											apt2 := lista_apt2[b]
											if b == 0 && !apt3lleno {
												bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt2)
												if !eba3 {
													return
												}
												lista_apt3 := bloque_apuntador3.B_pointers
												for c := range lista_apt3 {
													apt3 := lista_apt3[c]
													if apt3 == -1 {
														lista_apt3[b] = numero_bloque_disponible
														lista_a := structures.BloqueApuntadores{B_pointers: lista_apt3}
														Modificar_Apuntador(comando, path, superbloque.S_block_start, apt2, lista_a)
														Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
														agrego = true
														break
													}
												}
											}
											if b == 0 && !agrego {
												apt3lleno = true
												continue
											}
											if apt2 == -1 {
												lista_apt2[b] = numero_bloque_disponible
												lista_a := structures.BloqueApuntadores{B_pointers: lista_apt2}
												Modificar_Apuntador(comando, path, superbloque.S_block_start, apt2, lista_a)
												Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
												agrego = true
												break
											} else if agrego {
												break
											}
										}
									}
									if a == 0 && !agrego {
										apt2lleno = true
										continue
									}
									if apt1 == -1 && !agrego {
										lista_apt1[a] = numero_bloque_disponible
										lista_a := structures.BloqueApuntadores{B_pointers: lista_apt1}
										Modificar_Apuntador(comando, path, superbloque.S_block_start, lista_apt_bloques[contador], lista_a)
										Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
										if a == 15 {
											contador += 1
										}
										break
									} else if agrego {
										break
									}
								}
							} else {
								//apuntador 1
								num_bloque_apuntador, enba := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
								if !enba {
									return
								}
								Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)

								//apuntador 2
								num_bloque_apuntador2, enba2 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
								if !enba2 {
									return
								}
								Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador2)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador2, 1)

								//apuntador 3
								num_bloque_apuntador3, enba3 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
								if !enba3 {
									return
								}
								Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador3)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador3, 1)

								//apuntar a carpeta
								var lista_apt_apuntadores3 [16]int32
								for f := range lista_apt_apuntadores3 {
									lista_apt_apuntadores3[f] = int32(-1)
								}
								lista_apt_apuntadores3[0] = numero_bloque_disponible

								//apuntar a carpeta
								var lista_apt_apuntadores2 [16]int32
								for f := range lista_apt_apuntadores2 {
									lista_apt_apuntadores2[f] = int32(-1)
								}
								lista_apt_apuntadores2[0] = num_bloque_apuntador3

								//apuntar a carpeta
								var lista_apt_apuntadores [16]int32
								for f := range lista_apt_apuntadores {
									lista_apt_apuntadores[f] = int32(-1)
								}
								lista_apt_apuntadores[0] = num_bloque_apuntador2

								//giardar en el apuntador 13 del inodo
								lista_apt_bloques[contador] = num_bloque_apuntador

								lista_a := structures.BloqueApuntadores{B_pointers: lista_apt_apuntadores3}
								Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador3, lista_a)
								Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)

								lista_a.B_pointers = lista_apt_apuntadores2
								Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador2, lista_a)

								lista_a.B_pointers = lista_apt_apuntadores
								Modificar_Apuntador(comando, path, superbloque.S_block_start, num_bloque_apuntador, lista_a)
								continue
							}
						} else {
							if contador == 16 {
								color.Red("[" + comando + "]: No hay espacio")
							}
							lista_apt_bloques[contador] = numero_bloque_disponible
							Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, numero_bloque_disponible, 1)
							contador += 1
							agrego_bloque_nuevo_a = false
						}
					}
					size_cadena_anterior += 1
				} else {
					cadena += string(contenido_nuevo[i])
					size_cadena_anterior += 1
				}
			}
		}
		//Guardar lista de apuntadores en inodo
		inodo_archivo.I_block = lista_apt_bloques
		fecha := ObFechaInt()
		inodo_archivo.I_mtime = fecha
		inodo_archivo.I_s = int32(len(contenido_nuevo))
		igno := Guardar_Inodo(comando, path, superbloque.S_inode_start, inodo_archivo, numero_inodo_archivo)
		if !igno {
			return
		}
		color.Green("Modificando Archivo...")
	} else {
		color.Red("No existe el archivo en la ruta")
		return
	}
}

func Obtener_Contenido_Archivo_Users(comando string, id_usuario string) (string, bool) {
	ruta := "/users.txt"
	conjunto, path, econ := Obtener_Particion_ID(id_usuario)
	if !econ {
		color.Red("[" + comando + "]: Error al obtener particion")
		return "", false
	}
	// size_particion := 0
	superbloque := structures.SuperBlock{}
	// eslogica := false
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
				return "", false
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
				return "", false
			}
			if particion.Part_type == 'E' {
				color.Red("[MKFS]: No se puede obtener información de particion extendida")
				return "", false
			}
		}
	}

	ruta_separada := strings.Split(ruta, "/")
	cantidad_carpetas := len(ruta_separada)
	nombre_archivo := ruta_separada[cantidad_carpetas-1]
	ruta_sin_archivo := strings.ReplaceAll(ruta, "/"+nombre_archivo, "")
	// ruta_separada = strings.Split(ruta_sin_archivo, "/")

	numero_inodo_carpeta, enic := Encontrar_Ruta(comando, path, superbloque.S_inode_start, superbloque.S_block_start, ruta_sin_archivo)
	if !enic {
		return "", false
	}
	inodo_carpeta, eic := Obtener_Inodo(comando, path, superbloque.S_inode_start, numero_inodo_carpeta)
	if !eic {
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
				continue
			} else if i == 14 {
				continue
			} else if i == 15 {
				continue
			} else {
				bloque_carpeta, ebc := Obtener_Bloque(comando, path, superbloque.S_block_start, apuntador)
				if !ebc {
					return "", false
				}
				contenido_carpeta := bloque_carpeta.B_content
				for _, contenido := range contenido_carpeta {
					if ToString(contenido.B_name[:]) == ToString([]byte(nombre_archivo)) {
						numero_inodo_archivo = contenido.B_inodo
						break
					}
				}
			}
		}
	}
	if numero_inodo_archivo != -1 {
		inodo_archivo, eia := Obtener_Inodo(comando, path, superbloque.S_inode_start, numero_inodo_archivo)
		if !eia {
			return "", false
		}
		fecha := ObFechaInt()
		inodo_carpeta.I_atime = fecha
		guardado_inodo := Guardar_Inodo(comando, path, superbloque.S_inode_start, inodo_carpeta, numero_inodo_carpeta)
		if !guardado_inodo {
			return "", false
		}
		apuntadores_archivo := inodo_archivo.I_block
		contador := 0
		contenido_archivo := ""
		for i := range apuntadores_archivo {
			parte_archivo := apuntadores_archivo[i]
			if parte_archivo != -1 {
				if i == 13 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, parte_archivo)
					if !eba1 {
						return "", false
					}
					lista_apt1 := bloque_apuntador1.B_pointers
					for a := range lista_apt1 {
						apt1 := lista_apt1[a]
						if apt1 != -1 {
							bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							contenido := string(bloque_archivo.B_content[:])
							contenido = strings.ReplaceAll(contenido, "\x00", "")
							contenido_archivo += string(contenido)
						}
					}
					continue
				} else if i == 14 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, parte_archivo)
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
										bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										contenido := string(bloque_archivo.B_content[:])
										contenido = strings.ReplaceAll(contenido, "\x00", "")
										contenido_archivo += string(contenido)
									}
								}
								continue
							}
							bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							contenido := string(bloque_archivo.B_content[:])
							contenido = strings.ReplaceAll(contenido, "\x00", "")
							contenido_archivo += string(contenido)
						}
					}
					continue
				} else if i == 15 {
					bloque_apuntador1, eba1 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, parte_archivo)
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
											bloque_apuntador3, eba3 := Obtener_Bloque_Apuntador(comando, path, superbloque.S_block_start, apt2)
											if !eba3 {
												return "", false
											}
											lista_apt3 := bloque_apuntador3.B_pointers
											for c := range lista_apt3 {
												apt3 := lista_apt3[c]
												if apt3 != -1 {
													bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt3)
													if !eba {
														return "", false
													}
													contenido := string(bloque_archivo.B_content[:])
													contenido = strings.ReplaceAll(contenido, "\x00", "")
													contenido_archivo += string(contenido)
												}
											}
											continue
										}
										bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt2)
										if !eba {
											return "", false
										}
										contenido := string(bloque_archivo.B_content[:])
										contenido = strings.ReplaceAll(contenido, "\x00", "")
										contenido_archivo += string(contenido)
									}
								}
								continue
							}
							bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, apt1)
							if !eba {
								return "", false
							}
							contenido := string(bloque_archivo.B_content[:])
							contenido = strings.ReplaceAll(contenido, "\x00", "")
							contenido_archivo += string(contenido)
						}
					}
					continue
				} else {
					bloque_archivo, eba := Obtener_Bloque_Archivo(comando, path, superbloque.S_block_start, parte_archivo)
					if !eba {
						return "", false
					}
					contenido := string(bloque_archivo.B_content[:])
					contenido = strings.ReplaceAll(contenido, "\x00", "")
					contenido_archivo += string(contenido)
				}
			}
			contador += 1
		}
		return contenido_archivo, true
	} else {
		color.Red("El archivo no existe en la ruta")
		return "", false
	}
}
