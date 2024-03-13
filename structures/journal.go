package structures

type Journal struct { //156 bytes
	J_state   int8
	J_command [150]byte
	J_date    int32
	J_content [64]byte
}

//estado int8
//comando 150 byte
//horafecha int32
//contenido 64s
