package structures

type Journal struct { //172bytes
	J_operation   [10]byte //operacion
	J_content     [60]byte //contenido
	J_permissions int64    //permisos
	J_name        [60]byte //nombre
	J_owner       [10]byte //propietario
	J_date        int64    //fecha
	J_type        int64    //type
	J_size        int64    //tamano
}
