package comandos

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func ReportCommandProps(name string, instructions []string) {
	/*
		var _name string
		var _path string
		var _id string
		var _ruta string
	*/
	if strings.ToUpper(name) == "REP" {
		fmt.Println("GENERANDO REPORTE")
	} else {
		color.Red("[ReportCommandProps]: Internal Error")
	}
}
