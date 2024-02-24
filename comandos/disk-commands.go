package comandos

import (
	"MIA_P1_201906051/structures"
	"encoding/binary"
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
		}
	}
	return _size, _fit, _unit
}

func MKDISK_Create(_size int32, _fit byte, _unit byte) {
	directorio := "MIA/P1/Disks/"
	for i := 0; i < 26; i++ {
		letra := string('A'+i) + ".dsk"
		archivo := directorio + letra
		if _, err := os.Stat(archivo); os.IsNotExist(err) {
			CreateFile(archivo, _size, _fit, _unit)
			color.Green("[MKDISK]: Disco '" + letra + "' Creado")
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
		color.Red("[RMDISK]: No existe el disco")
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
		color.Red("[MOUNT]: Error en la obtención del disco")
		return
	}

	conjunto, error := BuscarParticion(tempDisk, _name, path)
	if error {
		color.Red("[MOUNT]: Partición no encontrada")
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
				continue
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
	Partitions_Mounted = append(Partitions_Mounted, arr_partition)

	color.Magenta("Montando partición... " + nombre_particion_montada)

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
			return
		}
	}

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

// MKFS

func Values_MKFS(instructions []string) {

}
