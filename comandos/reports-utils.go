package comandos

import (
	"MIA_P1_201906051/size"
	"MIA_P1_201906051/structures"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func ljust(s string, leng int) string {
	if len(s) >= leng {
		return s
	}
	return s + strings.Repeat(" ", leng-len(s))
}

func TieneNameRep(valor string) (string, bool) {
	if !strings.HasPrefix(strings.ToLower(valor), "name=") {
		color.Red("[REP]: No tiene name o tiene un valor no valido")
		return "", false
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[REP]: No tiene grp Valido")
		return "", false
	}
	if (strings.ToLower(value[1]) == "mbr") || (strings.ToLower(value[1]) == "disk") || (strings.ToLower(value[1]) == "inode") || (strings.ToLower(value[1]) == "journaling") || (strings.ToLower(value[1]) == "block") || (strings.ToLower(value[1]) == "bm_inode") || (strings.ToLower(value[1]) == "bm_block") || (strings.ToLower(value[1]) == "tree") || (strings.ToLower(value[1]) == "sb") || (strings.ToLower(value[1]) == "file") || (strings.ToLower(value[1]) == "ls") {
		return strings.ToLower(value[1]), true
	}
	return "", false
}

func TienePathRep(valor string) (string, bool) {
	if !strings.HasPrefix(strings.ToLower(valor), "path=") {
		color.Red("[REP]: No tiene path o tiene un valor no valido")
		return "", false
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[REP]: No tiene path Valido")
		return "", false
	}
	return value[1], true
}

func TieneIDRep(valor string) (string, bool) {
	comando := "REP"
	if !strings.HasPrefix(strings.ToLower(valor), "id=") {
		color.Red("[" + comando + "]: No tiene id o tiene un valor no valido")
		return "", false
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[" + comando + "]: No tiene id Valido")
		return "", false
	}
	if len(value[1]) > 4 {
		color.Red("[" + comando + "]: No tiene id Valido")
		return "", false
	}
	return value[1], true
}

func TieneRutaRep(valor string) (string, bool) {
	if !strings.HasPrefix(strings.ToLower(valor), "ruta=") {
		color.Red("[REP]: No tiene ruta o tiene un valor no valido")
		return "", false
	}
	value := strings.Split(valor, "=")
	if len(value) < 2 {
		color.Red("[REP]: No tiene ruta Valida")
		return "", false
	}
	return value[1], true
}

func returnstring(s string) string {

	if len(s) <= 0 || !(s != "") {
		return " "
	}
	return s
}

func Obtener_Disco_ID(id string) (string, bool) {
	for _, discos := range Partitions_Mounted {
		if disco, ok := discos.([]string); ok {
			var nombredisco = disco[0]
			if nombredisco == id {
				color.Red("[REP]: Particion Encontrada")
				return disco[3], true
			}
		}
	}
	return "", false
}

func ReducirSuperBloqueObtener(upad string, id_disco string, conjunto []interface{}) (structures.SuperBlock, bool) {
	superbloque := structures.SuperBlock{}

	logica := structures.EBR{}
	if conjunto[2] != nil {
		if temp_log, ok := conjunto[2].(structures.EBR); ok {
			v := reflect.ValueOf(temp_log)
			reflect.ValueOf(&logica).Elem().Set(v)
			var e bool
			superbloque, e = Obtener_Superbloque("REP", upad, ToString(logica.Name[:]))
			if !e {
				return structures.SuperBlock{}, false
			}
		}
		conjunto[0] = nil
		conjunto[1] = nil
	}

	particion := structures.Partition{}
	if conjunto[0] != nil {
		if temp, ok := conjunto[0].(structures.Partition); ok {
			v := reflect.ValueOf(temp)
			reflect.ValueOf(&particion).Elem().Set(v)
			var e bool
			superbloque, e = Obtener_Superbloque("REP", upad, ToString(particion.Part_name[:]))
			if !e {
				return structures.SuperBlock{}, false
			}
		}
	}
	return superbloque, true
}

func Obtener_Superbloque_Journal_Reducido(upad string, id_disco string, conjunto []interface{}) (structures.SuperBlock, int32, bool) {
	inicio_journaling := int32(-1)
	superbloque := structures.SuperBlock{}

	logica := structures.EBR{}
	if conjunto[2] != nil {
		if temp_log, ok := conjunto[2].(structures.EBR); ok {
			v := reflect.ValueOf(temp_log)
			reflect.ValueOf(&logica).Elem().Set(v)
			var e bool
			superbloque, e = Obtener_Superbloque("REP", upad, ToString(logica.Name[:]))
			if !e {
				return structures.SuperBlock{}, -1, false
			}
			inicio_journaling = logica.Part_start + size.SizeSuperBlock()
		}
		conjunto[0] = nil
		conjunto[1] = nil
	}

	particion := structures.Partition{}
	if conjunto[0] != nil {
		if temp, ok := conjunto[0].(structures.Partition); ok {
			v := reflect.ValueOf(temp)
			reflect.ValueOf(&particion).Elem().Set(v)
			var e bool
			superbloque, e = Obtener_Superbloque("REP", upad, ToString(particion.Part_name[:]))
			if !e {
				return structures.SuperBlock{}, -1, false
			}
			inicio_journaling = particion.Part_start + size.SizeSuperBlock()
		}
	}

	if inicio_journaling == superbloque.S_bm_inode_start {
		//ext2
		return structures.SuperBlock{}, -1, false
	}

	TamanioJournal := size.SizeJournal()
	numero_estructuras_journaling := (superbloque.S_bm_inode_start - inicio_journaling) / (TamanioJournal)
	if numero_estructuras_journaling > 0 {
		return superbloque, numero_estructuras_journaling, true
	}
	return structures.SuperBlock{}, -1, false
}

func ToInt(v string) int32 {
	i, _ := strconv.Atoi(v)
	return int32(i)
}
