package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"encoding/binary"
	"os"
	"reflect"

	"github.com/fatih/color"
)

func Obtener_Journaling(comando string, id string, numero_journaling int) (structures.Journal, bool) {
	inicio_journaling := int32(-1)
	conjunto, path, er := Obtener_Particion_ID(id)
	if !er {
		return structures.Journal{}, false
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
				return structures.Journal{}, false
			}
			inicio_particion := logica.Part_start
			inicio_journaling = inicio_particion + size.SizeSuperBlock()
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
				return structures.Journal{}, false
			}
			inicio_particion := logica.Part_start
			inicio_journaling = inicio_particion + size.SizeSuperBlock()

			if particion.Part_type == 'E' {
				color.Red("[" + comando + "]: No se puede obtener información de particion extendida")
				return structures.Journal{}, false
			}
		}
	}

	if inicio_journaling == superbloque.S_bm_inode_start {
		//ext2
		return structures.Journal{}, false
	}
	//ext3
	ubicabion_journaling := inicio_journaling + (size.SizeJournal() * int32(numero_journaling))
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return structures.Journal{}, false
	}
	defer file.Close()
	if _, err := file.Seek(int64(ubicabion_journaling), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return structures.Journal{}, false
	}
	journal := structures.Journal{}
	if err := binary.Read(file, binary.LittleEndian, &journal); err != nil {
		color.Red("[" + comando + "]: Error en lectura de Journal")
		return structures.Journal{}, false
	}
	return journal, true

}

func Obtener_Journaling_Disponible(comando string, id_particion string) (int32, bool) {
	inicio_journaling := int32(-1)
	conjunto, path, ec := Obtener_Particion_ID(id_particion)
	if !ec {
		return 0, false
	}
	superbloque := structures.SuperBlock{}
	// journaling := structures.Journal{}

	// inicio_particion := 0

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
			inicio_particion := logica.Part_start
			inicio_journaling = inicio_particion + size.SizeSuperBlock()
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
			inicio_particion := particion.Part_start
			inicio_journaling = inicio_particion + size.SizeSuperBlock()

			if particion.Part_type == 'E' {
				color.Red("[" + comando + "]: No se puede obtener información de particion extendida")
				return 0, false
			}
		}
	}

	if inicio_journaling == superbloque.S_bm_inode_start {
		color.Red("[" + comando + "]: Es una particion EXT2")
		return -1, true
	}
	inicio_journaling_actual := inicio_journaling
	contador := 0
	for inicio_journaling_actual < superbloque.S_inode_start {
		journaling_actual, eja := Obtener_Journaling(comando, id_particion, contador)
		if !eja {
			return 0, false
		}
		if journaling_actual.J_state == 0 {
			return int32(contador), true
		}
		contador += 1
	}

	return 1, false
}

func Modificar_Journaling(comando string, id string, numero_journaling int32, journaling structures.Journal) {
	inicio_journaling := int32(-1)
	conjunto, path, ec := Obtener_Particion_ID(id)
	if !ec {
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
			inicio_particion := logica.Part_start
			inicio_journaling = inicio_particion + size.SizeSuperBlock()
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
			inicio_particion := particion.Part_start
			inicio_journaling = inicio_particion + size.SizeSuperBlock()

			if particion.Part_type == 'E' {
				color.Red("[" + comando + "]: No se puede obtener información de particion extendida")
				return
			}
		}
	}

	if inicio_journaling == superbloque.S_bm_inode_start {
		return
	}
	ubicabion_journaling := inicio_journaling + (size.SizeJournal() * numero_journaling)
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		color.Red("[" + comando + "]: Error al leer Archivo")
		return 
	}
	defer file.Close()
	if _, err := file.Seek(int64(ubicabion_journaling), 0); err != nil {
		color.Red("[" + comando + "]: Error al mover el puntero")
		return 
	}

	if err := binary.Write(file, binary.LittleEndian, &journaling); err!= nil {
        color.Red("[" + comando + "]: Error en escritura de Journal")
        return
    }
	color.Green("[" + comando + "]: Journal modificado")
}
