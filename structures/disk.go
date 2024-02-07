package structures

// ESTRUCTURAS DISCOS

type EBR struct { //42bytes
	Part_mount byte     //indica si esta montada o no
	Part_fit   byte     //tipo de ajuste B (best) F (first) W (worst)
	Part_start int64    //indica en que byte del disco inicia particion
	Part_s     int64    //contiene tamaño total de particion en bytes
	Part_next  int64    //byte en el que esta proximo EBR -1 si no hay siguiente
	Name       [16]byte //nombre de particion
}

type Partition struct { //47bytes
	Part_status      byte     //indica si particion esta montada o no
	Part_type        byte     //indica el tipo de particion P (primaria) E (extendida)
	Part_fit         byte     //tipo de ajuste de particion B (best) F (first) W (worst)
	Part_start       int64    //indica en que byte del disco inicia la particion
	Part_s           int64    //contiene el tamaño total de la particion en bytes
	Part_name        [16]byte //nombre de particion
	Part_correlative int64    //correlativo de particion
	Part_id          [4]byte  //indica el id de particion generada al montar la misma
}

type MBR struct { //213bytes
	Mbr_tamano         int64        //tamaño total del disco en bytes
	Mbr_fecha_creacion int64        //fecha y hora de creacion del disco (time)
	Mbr_disk_signature int64        //numero random que identifica de forma unica cada disco
	Dsk_fit            byte         //tipo de ajuste de particion B (best) F (first) W (worst)
	Mbr_partitions     [4]Partition //estructura con informacion de las 4 particiones
}
