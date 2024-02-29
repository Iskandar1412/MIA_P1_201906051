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

func Obtener_Superbloque(comando string, path string, nombre string) (structures.SuperBlock, bool) {
	superbloque := structures.SuperBlock{}
	byte_escribir := int32(0)
	tempDisk, existe := ObtainMBRDisk(path)
	if !existe {
		color.Red("[" + comando + "]: Error en la obtención del disco")
		return superbloque, false
	}
	value := ToString([]byte(nombre))
	value2 := DevolverNombreByte(value)
	conjunto, _err := BuscarParticion(tempDisk, value2[:], path)
	if _err {
		color.Red("[" + comando + "]: No se ha encontrado en particiones lógicas")
		return superbloque, false
	}

	ebr := structures.EBR{}
	if temp, ok := conjunto[2].(structures.EBR); ok {
		v := reflect.ValueOf(temp)
		reflect.ValueOf(&ebr).Elem().Set(v)
		byte_escribir = ebr.Part_start
		conjunto[0] = nil
		conjunto[1] = nil
	}

	particion := structures.Partition{}
	if temp, ok := conjunto[0].(structures.Partition); ok {
		v := reflect.ValueOf(temp)
		reflect.ValueOf(&particion).Elem().Set(v)
		byte_escribir = particion.Part_start
	}

	if byte_escribir == int32(0) {
		return superbloque, false
	}

	file, err := os.Open(path)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return superbloque, false
	}
	defer file.Close()
	if _, err := file.Seek(int64(byte_escribir), 0); err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return superbloque, false
	}

	err = binary.Read(file, binary.LittleEndian, &superbloque)
	if err != nil {
		color.Red("[" + comando + "]: Error en lectura de superbloque")
		return superbloque, false
	}
	if superbloque.S_mtime == 0 && superbloque.S_umtime == 0 && superbloque.S_magic == 0 {
		color.Yellow("[" + comando + "]: Super Bloque no creado, Unidad no formateada")
		return superbloque, false
	}
	return superbloque, true
}

func Guardar_Superbloque(comando string, path string, inicio int32, superblock structures.SuperBlock) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return
	}
	defer file.Close()

	if _, err := file.Seek(int64(inicio), 0); err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &superblock); err != nil {
		color.Red("[" + comando + "]: Error en escritura de superbloque")
		return
	}
	color.Green("[" + comando + "]: Superbloque grabado correctamente")
}

func Crear_Bloque_Carpeta_Vacio(comando string, path string, inicio_bloque int32, num_bloque int32) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return
	}
	defer file.Close()

	carpet := structures.BlockCarpetas{}
	ubicacion_bloque := inicio_bloque + (num_bloque * size.SizeBlockCarpetas())
	not_used_name := DevolverNombreByte("")
	nombre_bytes := [10]byte(not_used_name[:])
	for i := range carpet.B_content {
		carpet.B_content[i].B_name = nombre_bytes
		carpet.B_content[i].B_inodo = int32(-1)
	}
	if _, err := file.Seek(int64(ubicacion_bloque), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &carpet); err != nil {
		color.Red("[" + comando + "]: Error en escritura de bloque")
		return
	}

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

func Obtener_Bloque(comando string, path string, inicio_bloques int32, no_bloque int32) (structures.BlockCarpetas, bool) {
	carpet := structures.BlockCarpetas{}
	ubicacion_bloque := inicio_bloques + (no_bloque * size.SizeBlockCarpetas())
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return structures.BlockCarpetas{}, false
	}
	defer file.Close()
	if _, err := file.Seek(int64(ubicacion_bloque), 0); err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return structures.BlockCarpetas{}, false
	}
	err = binary.Read(file, binary.LittleEndian, &carpet)
	if err != nil {
		color.Red("[" + comando + "]: Error en obtención de bloque")
		return structures.BlockCarpetas{}, false
	}
	return carpet, true
}

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

func Obtener_Bloque_Disponible(comando string, path string, inicio_bloques int32, cantidad_bloques int32) (int32, bool) {
	for i := 0; i < int(cantidad_bloques); i++ {
		// struct_bitmap := '\x00'
		file, err := os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			color.Red("[" + comando + "]: Error al leer Archivo")
			return 0, false
		}
		defer file.Close()
		//bytes_con_inodos := inicio_inodos + numero
		if _, err := file.Seek(int64(inicio_bloques)+int64(i), 0); err != nil {
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

func Crear_Bloque_Archivo_Vacio(comando string, path string, inicio_bloques int32, no_bloques int32) {
	struct_archivo := structures.BloqueArchivos{}
	size_archivo := size.SizeBloqueArchivos()
	ubicacion_bloque := inicio_bloques + (no_bloques * size_archivo)
	contenido_bytes := DevolverContenidoArchivo("")
	struct_archivo.B_content = contenido_bytes

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
	if err := binary.Write(file, binary.LittleEndian, &struct_archivo); err != nil {
		color.Red("[" + comando + "]: Error en escritura de bloque")
		return
	}
	color.Green("[" + comando + "]: Archivo creado")
}

func Obtener_Bloque_Apuntador(comando string, path string, inicio_bloques int32, no_bloque int32) (structures.BloqueApuntadores, bool) {
	ubicacion_bloque := inicio_bloques + (no_bloque * size.SizeBloqueArchivos())

	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return structures.BloqueApuntadores{}, false
	}
	defer file.Close()
	//bytes_con_inodos := inicio_inodos + numero
	if _, err := file.Seek(int64(ubicacion_bloque), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return structures.BloqueApuntadores{}, false
	}

	//Lectura archivo
	valor := structures.BloqueApuntadores{}
	if err := binary.Read(file, binary.LittleEndian, &valor); err != nil {
		color.Red("[" + comando + "]: Error en lectura de bloque")
		return structures.BloqueApuntadores{}, false
	}

	return valor, true
}

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

func Modificar_Apuntador(comando string, path string, inicio_bloques int32, no_bloque int32, content structures.BloqueApuntadores) {
	ubicacion_bloque := inicio_bloques + (no_bloque * size.SizeBloqueApuntadores())
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
	color.Green("[" + comando + "]: Apuntador modificado")
}

func Crear_Bloque_Apuntador_Vacio(comando string, path string, inicio_bloques int32, no_bloque int32) {
	bloque_apt := structures.BloqueApuntadores{}
	for i := range bloque_apt.B_pointers {
		bloque_apt.B_pointers[i] = int32(-1)
	}
	ubicacion_bloque := inicio_bloques + (no_bloque * size.SizeBloqueArchivos())
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

	if err := binary.Write(file, binary.LittleEndian, &bloque_apt); err != nil {
		color.Red("[" + comando + "]: Error en escritura de bloque")
		return
	}
	color.Green("[" + comando + "]: Apuntador creado")
}

// Creacion de archivo vacíos
func Crear_Archivo_Vacio(comando string, path string, ruta string, superbloque structures.SuperBlock, contenido string, nombre_archivo string, id_user int32, id_group int32) {
	numero_inodo_disponible, existe_inodo := Obtener_Inodo_Disponible(comando, path, superbloque.S_bm_inode_start, superbloque.S_inodes_count)
	if !existe_inodo {
		color.Red("[" + comando + "]: Error al buscar inodo")
		return
	}
	// fmt.Println(numero_inodo_disponible) //se usará luego

	//información padre
	num_inodo_padre, existe_padre := Encontrar_Ruta(comando, path, superbloque.S_inode_start, superbloque.S_block_start, ruta)
	if !existe_padre {
		return
	}
	//inodo padre
	inodo_padre, existe_ipadre := Obtener_Inodo(comando, path, superbloque.S_inode_start, num_inodo_padre)
	if !existe_ipadre {
		return
	}

	validacion := Validar_Permisos(comando, inodo_padre.I_uid, inodo_padre.I_gid, id_user, id_group, inodo_padre.I_perm, int32(2))
	//fmt.Println(inodo_padre)
	if !validacion {
		color.Magenta("No se tienen los permisos necesarios")
		return
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
			return
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
					return
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
							return
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
							return
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
							return
						}
						Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)

						//creacion apuntador 2
						num_bloque_apuntador2, nba2 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
						if !nba2 {
							return
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
							return
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
							return
						}
						Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador, 1)
						//apuntador2

						num_bloque_apuntador2, eba2 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
						if !eba2 {
							return
						}
						Crear_Bloque_Apuntador_Vacio(comando, path, superbloque.S_block_start, num_bloque_apuntador2)
						Modificar_Bitmap(comando, path, superbloque.S_bm_block_start, num_bloque_apuntador2, 1)

						//apuntador 3
						num_bloque_apuntador3, eba3 := Obtener_Bloque_Disponible(comando, path, superbloque.S_bm_block_start, superbloque.S_blocks_count)
						if !eba3 {
							return
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

	err_g_inodo := Guardar_Inodo(comando, path, superbloque.S_inode_start, inodo, numero_inodo_disponible)
	if !err_g_inodo {
		color.Red("No se pudo guardar el inodo")
		return
	}

	Modificar_Bitmap(comando, path, superbloque.S_bm_inode_start, numero_inodo_disponible, 1)
	Agregar_Bloque_Lista_Inodos(comando, path, superbloque.S_block_start, numero_inodo_disponible, superbloque.S_bm_block_start, superbloque.S_blocks_count, nombre_archivo, inodo_padre, superbloque.S_inode_start, num_inodo_padre, superbloque.S_bm_inode_start, numember_block_disp)
}
