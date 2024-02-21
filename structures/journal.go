package structures

type Journal struct { //156 bytes
	J_operation   [10]byte //operacion
	J_content     [60]byte //contenido
	J_permissions int32    //permisos
	J_name        [60]byte //nombre
	J_owner       [10]byte //propietario
	J_date        int32    //fecha
	J_type        int32    //type
	J_size        int32    //tamano
}
