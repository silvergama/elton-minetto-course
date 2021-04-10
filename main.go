package main

import (
	"fmt"
	"strings"
)

func main() {
	// uma função pode retornar mais de um valor
	// a contrução := cria uma var caso ela não exista, já definindo seu tipo
	ok, err := say("Hello World")
	if err != nil {
		panic(err.Error())
	}

	switch ok {
	case true:
		fmt.Println("Deu certo")
	default:
		fmt.Println("Deu errado")
	}
}

//as funções deve declarar o tipo de cada variável que recebe ou que retorna
func say(what string) (bool, error) {
	if strings.TrimSpace(what) == "" {
		return false, fmt.Errorf("Empty string")
	}

	fmt.Println(what)

	return true, nil
}
