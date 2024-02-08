package comandos

import "fmt"

func GlobalCom(lista []string) {
	for _, comm := range lista {
		fmt.Println(comm)
	}
}
