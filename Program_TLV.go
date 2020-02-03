package main
	
import (
	"fmt"
	"strconv"
	"io/ioutil"
	"os"
	"bufio"
	"strings"
)

func main(){	
	//Lee Archivo externo con datos TLV.
	dat := readFile("TLV.txt")
	
	//Procesa array byte con TLV.
	mapResult, errores  := processFile(dat)
	
	//Muestra resultados por consola.
	exportResult(mapResult,errores)	
}

//-------------------------Funciones.------------------------------

//Retorna Array de Byte con TLV.
func readFile(path string) []byte {
	dat, err := ioutil.ReadFile(path)
	check(err,("Error al abrir archivo: "+ path + ". Verifique que se encuentre en carpeta."))
	return dat
}

//Procesa TLV.
func processFile(dat []byte) (map[int]map[string]string, string){

	mapResult := map[int]map[string]string{}
	errores := ""
	
	tlv := string(dat)
	tlv_len := len(tlv)
	i := 0
	count := 1
	
	for{
		//Cierra iteración de string.
		if(i >= tlv_len){
			break
		}
		//Obtener largo de string valor tlv
		large, err := strconv.Atoi(string(tlv[i:i+2]))
		if err != nil {
			errores = "Cadena de no cumple con formato TLV."
			break
		}
		if(i+large+5>tlv_len){
			errores = "Cadena de no cumple con formato TLV."
			break
		}

		tlv_aux := string(tlv[i:i+large+5])

		//Función que genera valores.
		values := getValues(tlv_aux)
		
		//Validar tipo de dato.
		if validaType(values["Tipo"][:1],values["Valor"]){
			errores = "Valor no corresponde al tipo de dato Numerico. TLV:" + strconv.Itoa(count)
			break
		}				
		
		//Llena varible Result.
		mapResult[count] = values

		//Variables de iteración.
		i += large + 5
		count++
	}
	return mapResult, errores
}

func exportResult(mapResult map[int]map[string]string,errores string){
	if strings.TrimSpace(errores) == ""{
		for key, tlv := range mapResult {
			fmt.Println("TLV:", key)
			for _key, value := range tlv {
				fmt.Println("\t",_key, "\t:", value)
			}
		}
		fmt.Print("\n\nPresione cualquier tecla para continuar")
  		bufio.NewReader(os.Stdin).ReadBytes('\n') 
	}else{
		fmt.Print(errores)
		fmt.Print("\n\nPresione cualquier tecla para continuar")
  		bufio.NewReader(os.Stdin).ReadBytes('\n') 
	}	
}

func check(e error, mensaje string) {
    if e != nil {
		fmt.Println(mensaje)
    }
}

func getValues(tlv string) map[string]string{
	values := make(map[string]string)
	values["Largo"] = string(tlv[:2])
	values["Tipo"] = valueType(string(tlv[2:5]))
	values["Valor"] = string(tlv[5:])
	return values
}


func valueType(tipo string) string{	
	switch string(tipo[:1]) {
	case "A":
		return "Alfanumerico de largo:" + string(tipo[1:])
	case "N":
		return "Numerico de largo:" + string(tipo[1:])
	default:
		return "No identificado"
	}
}

func validaType(tipo string, value string) bool{	
	error := false
	switch tipo {
	case "A":
		error = false
	case "N":
		if _, err := strconv.Atoi(value); err != nil {
			error = true
		}
	}
	return error
}
