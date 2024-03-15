package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// MKDISK

func Values_MKDISK(instructions []string) (int32, byte, byte) {
	var _size int32
	var _fit byte = 'F'
	var _unit byte = 'M'
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "size") {
			var value = TieneSize("MKDISK", valor)
			_size = value
		} else if strings.HasPrefix(strings.ToLower(valor), "fit") {
			var value = TieneFit("MKDISK", valor)
			_fit = value
		} else if strings.HasPrefix(strings.ToLower(valor), "unit") {
			var value = TieneUnit("mkdisk", valor)
			_unit = value
		} else {
			color.Yellow("[MKDISK]: Atributo no reconocido")
			return -1, '0', '0'
		}
	}
	return _size, _fit, _unit
}

func MKDISK_Create(_size int32, _fit byte, _unit byte) {
	directorio := "MIA/P1/Disks/"
	for i := 0; i < 26; i++ {
		letra := fmt.Sprintf("%c.dsk", 'A'+i)
		archivo := directorio + letra
		if _, err := os.Stat(archivo); os.IsNotExist(err) {
			CreateFile(archivo, _size, _fit, _unit)
			color.Green("[MKDISK]: Disco '" + letra + "' Creado -> S(" + fmt.Sprint(_size) + ") U(" + string(_unit) + ")")
			break
		} else {
			continue
			//color.Yellow("[MKDISK]: Disco Existente")
		}
	}
}

func CreateFile(archivo string, _size int32, _fit byte, _unit byte) {
	file, err := os.Create(archivo)
	if err != nil {
		color.Red("Error al crear el archivo")
		return
	}
	defer file.Close()
	//Escribir datos en archivo
	var estructura structures.MBR
	tamanio := Tamano(_size, _unit)
	estructura.Mbr_tamano = tamanio
	estructura.Mbr_fecha_creacion = ObFechaInt()
	estructura.Mbr_disk_signature = ObDiskSignature()
	estructura.Dsk_fit = _fit
	for i := 0; i < len(estructura.Mbr_partitions); i++ {
		estructura.Mbr_partitions[i] = PartitionVacia()
	}
	//Llenar el archivo
	bytes_llenar := make([]byte, int(tamanio))
	if _, err := file.Write(bytes_llenar); err != nil {
		color.Red("Error al escribir bytes en el archivo")
		return
	}

	//mover de posicion
	if _, err := file.Seek(0, 0); err != nil {
		color.Red("Error al mover puntero del archivo")
		return
	}

	if err := binary.Write(file, binary.LittleEndian, &estructura); err != nil {
		color.Red("Error al escribir datos del MBR")
		return
	}
}

// RMDISK

func Values_RMDISK(instructions []string) (byte, bool) {
	var _driveletter byte = '0'
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "driveletter") {
			var value = TieneDriveLetter("RMDISK", valor)
			_driveletter = value
			break
		} else {
			color.Yellow("[RMDISK]: Atributo no reconocido")
			_driveletter = '0'
			break
		}
	}
	if _driveletter == '0' {
		return '0', false
	} else {
		return _driveletter, true
	}
}

func RMDISK_EXECUTE(_driveletter byte) {
	PATH := "MIA/P1/Disks/" + string(_driveletter) + ".dsk"
	if _, err := os.Stat(PATH); os.IsNotExist(err) {
		color.Red("[RMDISK]: No existe el disco -> " + string(_driveletter))
		return
	}
	err := os.Remove(PATH)
	if err != nil {
		color.Red("[RMDISK]: Error al borrar el disco")
		return
	}
	color.Green("[RMDISK]: Disco '" + string(_driveletter) + ".dsk' Borrado")
}

// MOUNT

func Values_Mount(instructions []string) (byte, [16]byte, bool) {
	var _driveletter byte
	var _name [16]byte
	var _error = false
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "driveletter") {
			var value = TieneDriveLetter("MOUNT", valor)
			_driveletter = value
		} else if strings.HasPrefix(strings.ToLower(valor), "name") {
			var value = TieneNombre("MOUNT", valor)
			if len(value) > 16 {
				color.Red("[MOUNT]: El nombre no puede ser mayor a 16 caracteres")
				_error = true
				break
			} else {
				_name = DevolverNombreByte(value)
			}
		} else {
			color.Yellow("[MOUNT]: Atributo no reconocido")
		}
	}
	return _driveletter, _name, _error
}

func MOUNT_EXECUTE(_driveletter byte, _name []byte) {
	path := "MIA/P1/Disks/" + string(_driveletter) + ".dsk"
	if !ExisteArchivo("MOUNT", path) {
		color.Yellow("[MOUNT]: Cancel the operation because not yet a file")
		return
	}

	//obtener disco
	tempDisk, existe := ObtainMBRDisk(path)
	if !existe {
		color.Red("[MOUNT]: Error en la obtención del disco -> " + path)
		return
	}

	conjunto, error := BuscarParticion(tempDisk, _name, path)
	if error {
		color.Red("[MOUNT]: Partición no encontrada -> " + path)
		return
	}

	//fmt.Println(string(_driveletter))
	//fmt.Println(conjunto)

	if conjunto[0] == nil {
		color.Red("[MOUNT]: No hay particion para usar")
		return
	}

	particion := structures.Partition{}
	if temp, ok := conjunto[0].(structures.Partition); ok {
		v := reflect.ValueOf(temp)
		reflect.ValueOf(&particion).Elem().Set(v)
	}
	//inicio_particion := conjunto[1]
	//logica_interfaz := conjunto[2]
	logica := structures.EBR{}
	//inicio_logica := conjunto[3]
	numero := 1
	//bandera := false
	contador := 1
	//revisar anteriores montadas en el caso que hayan
	for _, discos := range Partitions_Mounted {
		//fmt.Println(discos, " - ")
		if disco, ok := discos.([]string); ok {
			var nombredisco = ""
			// fmt.Println(disco[2])
			nombredisco = disco[2]
			if nombredisco == string(_driveletter) {
				//color.Red("[MOUNT]: El disco ya está montado")
				if disco[1] == ToString(_name) {
					color.Red("[MOUNT]: La partición ya está montada")
					// bandera = true
					return
				}
				contador += 1

			}
		}
		numero = contador
		//fmt.Println(nombredisco)
		//validacion := false
	}

	//crear particion
	nombre_particion_montada := string(_driveletter) + strconv.Itoa(numero) + "51"
	nombre_bytes := DevolverNombreByte(nombre_particion_montada)
	var arr_partition []string
	arr_partition = append(arr_partition, nombre_particion_montada, ToString(_name), string(_driveletter), path)

	color.Magenta("Montando partición... " + nombre_particion_montada + ", de la particion: " + ToString(_name))

	if reflect.TypeOf(conjunto[2]) == reflect.TypeOf(structures.EBR{}) {
		if temp_log, ok := conjunto[2].(structures.EBR); ok {
			v := reflect.ValueOf(temp_log)
			reflect.ValueOf(&logica).Elem().Set(v)
		}
		logica.Part_mount = int8(1)
		conjunto[0] = nil
		conjunto[1] = nil
		Escribir_EBR("MOUNT", path, logica, conjunto[3].(int32))
	}

	if reflect.TypeOf(conjunto[0]) == reflect.TypeOf(structures.Partition{}) {
		particion.Part_status = int8(1)
		particion.Part_id = [4]byte(nombre_bytes[:])
		Escribir_Particion("MOUNT", path, particion, conjunto[1].(int32))
		if particion.Part_type == 'E' {
			color.Red("[MOUNT]: No se puede montar una partición extendida")
			return
		}
	}

	Partitions_Mounted = append(Partitions_Mounted, arr_partition)

	inicio := int32(0)
	if conjunto[0] != nil {
		inicio = particion.Part_start
	}
	if conjunto[2] != nil {
		inicio = logica.Part_start
	}

	superblock, err := Obtener_Superbloque("MOUNT", path, ToString(_name))
	if err {
		fecha := ObFechaInt()
		superblock.S_mtime = fecha
		val := superblock.S_mnt_count
		val += 1
		superblock.S_mnt_count = val
		Guardar_Superbloque("MOUNT", path, inicio, superblock)
	}
	//fmt.Println(superblock)
	//Particiones_Montadas = append(Particiones_Montadas, "asdf")
	//fmt.Println(particion, inicio_particion, logica, inicio_logica)
}

// UNMOUNT

func Values_Unmount(instructions []string) (string, bool) {
	var _id = ""
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "id") {
			var value = TieneID("UNMOUNT", valor)
			_id = value
			break
		} else {
			color.Yellow("[UNMOUNT]: Atributo no reconocido")
			_id = ""
			break
		}
	}
	if _id == "" || len(_id) == 0 || len(_id) > 4 {
		return "", false
	} else {
		return _id, true
	}
}

func UNMOUNT_EXECUTE(comando string, nombre string) {
	//fmt.Println(comando, nombre)
	//tempDisk := structures.MBR{}
	bandera := false
	for _, discos := range Partitions_Mounted {
		if disco, ok := discos.([]string); ok {
			temp, existe := ObtainMBRDisk(disco[3])
			if !existe {
				color.Red("[UNMOUNT]: Error en la obtención del disco")
				return
			}

			// tempDisk = temp
			if disco[2] != string(nombre[0]) {
				//color.Red("Nombre Disco '" + string(nombre[0]) + "' no existe")
				bandera = false
				continue
			}
			if disco[0] != string(nombre) {
				bandera = false
				continue
			}

			bandera = true

			//if DevolverNombreByte(disco[1]) ==
			nombre2 := DevolverNombreByte(disco[1])
			conjunto, error := BuscarParticion(temp, nombre2[:], disco[3])
			if error {
				color.Red("[UNMOUNT]: Partición no encontrada")
				return
			}

			logica := structures.EBR{}
			if conjunto[2] != nil {
				if temp_log, ok := conjunto[2].(structures.EBR); ok {
					v := reflect.ValueOf(temp_log)
					reflect.ValueOf(&logica).Elem().Set(v)
				}
				conjunto[0] = nil
				conjunto[1] = nil
			}

			particion := structures.Partition{}
			if conjunto[0] != nil {
				if temp, ok := conjunto[0].(structures.Partition); ok {
					v := reflect.ValueOf(temp)
					reflect.ValueOf(&particion).Elem().Set(v)
				}
			}

			//caso de ser partición lógica
			inicio := int32(0)
			if conjunto[0] != nil {
				inicio = particion.Part_start
				particion.Part_id = [4]byte{'\x00', '\x00', '\x00', '\x00'}
				particion.Part_status = int8(0)
				Escribir_Particion("UNMOUNT", disco[3], particion, conjunto[1].(int32))
				color.Green("[UNMOUNT]: Particion '" + ToString(particion.Part_name[:]) + "' desmontada, ID: " + nombre)
			}
			if conjunto[2] != nil {
				inicio = logica.Part_start
				logica.Part_mount = int8(0)
				Escribir_EBR("UNMOUNT", disco[3], logica, conjunto[3].(int32))
				color.Green("[UNMOUNT]: Particion '" + ToString(logica.Name[:]) + "' desmontada, ID: " + nombre)
			}

			superblock, er := Obtener_Superbloque("UNMOUNT", disco[3], disco[1])
			if er {
				fecha := ObFechaInt()
				superblock.S_umtime = fecha
				Guardar_Superbloque("UNMOUNT", disco[3], inicio, superblock)
			}

			//Eliminar particion de la lista
			EliminarParticionMount("UMOUNT", disco[0])
			return
			//fmt.Println(conjunto)
		}
	}
	if !bandera {
		color.Red("[UNMOUNT]: ID de partición '" + nombre + "' no encontrada")
		return
	}
	//fmt.Println(tempDisk)
}

// MKFS

func Values_MKFS(instructions []string) (string, string, string, bool) {
	var _id string
	var _type = "FULL"
	var _fs = "2fs"
	for _, valor := range instructions {
		if strings.HasPrefix(strings.ToLower(valor), "id") {
			var value = TieneID("MKFS", valor)
			_id = value
		} else if strings.HasPrefix(strings.ToLower(valor), "type") {
			var value = TieneTypeMKFS(strings.ToLower(valor))
			_type = value
		} else if strings.HasPrefix(strings.ToLower(valor), "fs") {
			var value = TieneFS(strings.ToLower(valor))
			_fs = value
		} else {
			color.Yellow("[MKFS]: Atributo no reconocido")
			_id = ""
			break
		}
	}
	if (_id == "" || len(_id) == 0 || len(_id) > 4) || (_type == "" || _fs == "") {
		return "", "", "", false
	} else {
		return _id, _type, _fs, true
	}
}

func MKFS_EXECUTE(nombre string, tipo string, fs string) {
	conjunto, path, funciona := Obtener_Particion_ID(nombre)
	if !funciona {
		color.Red("[MKFS]: Partición no encontrada '" + nombre + "' en las particiones montadas")
		return
	}
	// fmt.Println(conjunto)
	///fmt.Println(conjunto)

	filesystemtype := 0
	if strings.ToUpper(fs) == "2FS" {
		filesystemtype = 2
	}
	if strings.ToUpper(fs) == "3FS" {
		filesystemtype = 3
	}
	color.Cyan("Se encontro la partición... " + nombre)
	// bytes_inicio := -1
	size_particion := 0

	logica := structures.EBR{}
	if conjunto[2] != nil {
		if temp_log, ok := conjunto[2].(structures.EBR); ok {
			v := reflect.ValueOf(temp_log)
			reflect.ValueOf(&logica).Elem().Set(v)
			// bytes_inicio = int(logica.Part_start)
			size_particion = int(logica.Part_s)
			conjunto[0] = nil
			conjunto[1] = nil
		}
	}

	particion := structures.Partition{}
	if conjunto[0] != nil {
		if temp, ok := conjunto[0].(structures.Partition); ok {
			v := reflect.ValueOf(temp)
			reflect.ValueOf(&particion).Elem().Set(v)
			// bytes_inicio = int(particion.Part_start)
			size_particion = int(particion.Part_s)
			if particion.Part_type == 'E' {
				color.Red("[MKFS]: No se puede formatear partición extendida")
				return
			}
		}
	}

	//tamaños
	tamanio_superblock := size.SizeSuperBlock()
	tamanio_journaling := size.SizeJournal()
	tamanio_inodo := size.SizeInode()
	tamanio_bloque_archivo := size.SizeBloqueArchivos()

	//valores
	NumEstructuras := 0
	inicio_journal := int32(0)
	inicio_bitmap_inodos := int32(0)
	inicio_bitmap_bloques := int32(0)
	inicio_inodos := int32(0)
	inicio_bloques := int32(0)
	fecha := int32(0)
	SupahBlock := structures.SuperBlock{}

	if filesystemtype == 2 {
		Estructuras := float64(size_particion-int(tamanio_superblock)) / float64((4+int(tamanio_inodo))+(3*int(tamanio_bloque_archivo)))
		NumEstructuras = int(math.Ceil(Estructuras))
		// numeroEntero := int(numeroEstructuras)
		// fmt.Println(numeroEntero, tamanio_journaling-int32(bytes_inicio))

		inicio_bitmap_inodos = particion.Part_start + tamanio_superblock
		inicio_bitmap_bloques = inicio_bitmap_inodos + int32(NumEstructuras)
		inicio_inodos = inicio_bitmap_bloques + (int32(NumEstructuras) * 3)
		inicio_bloques = inicio_inodos + (tamanio_inodo * int32(NumEstructuras))

		fecha = ObFechaInt()

		superBloque := structures.SuperBlock{}
		superBloque.S_filesistem_type = int32(filesystemtype)
		superBloque.S_inodes_count = int32(NumEstructuras)
		superBloque.S_blocks_count = int32(NumEstructuras) * 3
		superBloque.S_free_blocks_count = int32(NumEstructuras) * 3
		superBloque.S_free_inodes_count = int32(NumEstructuras)
		superBloque.S_mtime = fecha
		superBloque.S_umtime = fecha
		superBloque.S_mnt_count = 1
		superBloque.S_magic = int32(0xEF53)
		superBloque.S_inode_s = tamanio_inodo
		superBloque.S_block_s = tamanio_bloque_archivo
		superBloque.S_first_ino = inicio_inodos
		superBloque.S_first_blo = inicio_bloques
		superBloque.S_bm_inode_start = inicio_bitmap_inodos
		superBloque.S_bm_block_start = inicio_bitmap_bloques
		superBloque.S_inode_start = inicio_inodos
		superBloque.S_block_start = inicio_bloques
		SupahBlock = superBloque

		//fmt.Println(tamanio_journaling, bytes_inicio)
		Guardar_Superbloque("MKFS", path, particion.Part_start, superBloque)
		color.Green("[MKFS]: Disco Formateado correctamente")
	}
	if filesystemtype == 3 {
		Estructuras := float64(size_particion-int(tamanio_superblock)) / float64((4+int(tamanio_inodo))+(3*int(tamanio_bloque_archivo)+int(tamanio_journaling)))
		NumEstructuras = int(math.Ceil(Estructuras))
		inicio_journal = particion.Part_start + (tamanio_superblock)
		inicio_bitmap_inodos = inicio_journal + (int32(NumEstructuras) * tamanio_journaling)
		inicio_bitmap_bloques = inicio_bitmap_inodos + int32(NumEstructuras)
		inicio_inodos = inicio_bitmap_bloques + (int32(NumEstructuras) + 3)
		inicio_bloques = inicio_inodos + (tamanio_inodo * int32(NumEstructuras))

		fecha = ObFechaInt()

		superBloque := structures.SuperBlock{}
		superBloque.S_filesistem_type = int32(filesystemtype)
		superBloque.S_inodes_count = int32(NumEstructuras)
		superBloque.S_blocks_count = int32(NumEstructuras) * 3
		superBloque.S_free_blocks_count = int32(NumEstructuras) * 3
		superBloque.S_free_inodes_count = int32(NumEstructuras)
		superBloque.S_mtime = fecha
		superBloque.S_umtime = fecha
		superBloque.S_mnt_count = 1
		superBloque.S_magic = int32(0xEF53)
		superBloque.S_inode_s = tamanio_inodo
		superBloque.S_block_s = tamanio_bloque_archivo
		superBloque.S_first_ino = inicio_inodos
		superBloque.S_first_blo = inicio_bloques
		superBloque.S_bm_inode_start = inicio_bitmap_inodos
		superBloque.S_bm_block_start = inicio_bitmap_bloques
		superBloque.S_inode_start = inicio_inodos
		superBloque.S_block_start = inicio_bloques
		SupahBlock = superBloque

		Guardar_Superbloque("MKFS", path, particion.Part_start, superBloque)
		color.Green("[MKFS]: Disco Formateado correctamente")
	}
	lista_inodo := make([]int32, 16)
	for i := range lista_inodo {
		lista_inodo[i] = int32(-1)
	}
	inodo := structures.Inode{}
	//fmt.Println(inodo)
	inodo.I_uid = int32(-1)
	inodo.I_gid = int32(-1)
	inodo.I_s = int32(-1)
	inodo.I_atime = fecha
	inodo.I_ctime = fecha
	inodo.I_mtime = fecha
	inodo.I_block = [16]int32(lista_inodo)
	inodo.I_type = int32(-1)
	inodo.I_perm = int32(-1)

	for i := 0; i < NumEstructuras; i++ {
		res := Guardar_Inodo("MKFS", path, inicio_inodos, inodo, int32(i))
		if !res {
			color.Red("[MKFS]: No se guardo el Inodo")
		}
	}
	Vaciar_Conjunto_Bitmap("MKFS", path, inicio_bitmap_inodos, int32(NumEstructuras))
	// bitmap = '\0' int8
	Vaciar_Conjunto_Bitmap("MKFS", path, inicio_bitmap_bloques, (int32(NumEstructuras) * 3))

	color.Cyan("Creando carpeta Raíz")

	var contenido_carpeta []interface{}
	nombre_carpeta := DevolverNombreByte("raiz")
	nombre_content := DevolverNombreByte("")
	contenido := structures.Content{B_name: [10]byte(nombre_carpeta[:])}
	contenido_vacio := structures.Content{B_name: [10]byte(nombre_content[:]), B_inodo: int32(-1)}

	contenido_carpeta = append(contenido_carpeta, contenido)
	contenido_carpeta = append(contenido_carpeta, contenido)
	contenido_carpeta = append(contenido_carpeta, contenido_vacio)
	contenido_carpeta = append(contenido_carpeta, contenido_vacio)

	//crear bloque vacio carpeta
	Crear_Bloque_Carpeta_Vacio("MKFS", path, inicio_bloques, int32(0))
	//modificar carpeta
	Modificar_Carpeta("MKFS", path, inicio_bloques, int32(0), contenido_carpeta)

	//CREACION INODO RAIZ
	//CREACION INODO RAIZ
	lista_inodo[0] = (int32(0)) //apuntamos al inicio del bloque (raiz)
	inodo_raiz := structures.Inode{}
	inodo_raiz.I_uid = int32(0)
	inodo_raiz.I_gid = int32(0)
	inodo_raiz.I_s = int32(0)
	inodo_raiz.I_atime = fecha
	inodo_raiz.I_ctime = fecha
	inodo_raiz.I_mtime = fecha
	inodo_raiz.I_block = [16]int32(lista_inodo)
	inodo_raiz.I_type = int32(0)
	inodo_raiz.I_perm = int32(777)

	res := Guardar_Inodo("MKFS", path, inicio_inodos, inodo_raiz, 0)
	if !res {
		color.Red("[MKFS]: No se guardo el Inodo")
		return
	}
	Modificar_Bitmap("MKFS", path, inicio_bitmap_inodos, int32(0), int32(1))
	Modificar_Bitmap("MKFS", path, inicio_bitmap_bloques, int32(0), int32(1))

	//Creacion archivo users.txt
	color.Yellow("Creando archivo Users.TXT")

	contenido_usuario := "1,G,root		\n"
	contenido_usuario += "1,U,root		,root	,123	\n"
	contenido_usuario = strings.ReplaceAll(contenido_usuario, "\t", "")

	Crear_Archivo_Vacio("MKFS", path, "/", SupahBlock, contenido_usuario, "users.txt", 1, 1)
}
