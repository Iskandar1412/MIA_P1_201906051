package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"os"
	"reflect"

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

func Obtener_Bloque_Archivo(comando string, path string, inicio_bloques int32, no_bloque int32) (structures.BloqueArchivos, bool) {
	struct_archivo := structures.BloqueArchivos{}
	ubicacion_bloque := inicio_bloques + (no_bloque * size.SizeBloqueArchivos())
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return structures.BloqueArchivos{}, false
	}
	defer file.Close()
	if _, err := file.Seek(int64(ubicacion_bloque), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return structures.BloqueArchivos{}, false
	}
	if err := binary.Read(file, binary.LittleEndian, &struct_archivo); err != nil {
		color.Red("[" + comando + "]: Error en lectura de bloque")
		return structures.BloqueArchivos{}, false
	}
	return struct_archivo, true
}
