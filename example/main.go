package main

import (
	"fmt"
	"log"

	"github.com/vynazevedo/termx"
)

func main() {
	// Exemplo 1: Formulário simples
	var nome string
	var ambiente string
	var concordou bool

	err := termx.Form().
		Input("Seu nome:", &nome).
		Select("Ambiente:", []string{"dev", "stage", "prod"}, &ambiente).
		Confirm("Você aceita os termos?", &concordou).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nResultados do Formulário:\n")
	fmt.Printf("Nome: %s\n", nome)
	fmt.Printf("Ambiente: %s\n", ambiente)
	fmt.Printf("Concordou: %v\n", concordou)

	// Exemplo 2: Input com validação
	fmt.Println("\nPressione qualquer tecla para o próximo exemplo...")
	fmt.Scanln()

	var email string
	err = termx.Input("Email:", &email).
		WithPlaceholder("usuario@exemplo.com").
		WithValidator(termx.Email()).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nEmail: %s\n", email)

	// Exemplo 3: Input de senha
	fmt.Println("\nPressione qualquer tecla para o próximo exemplo...")
	fmt.Scanln()

	var senha string
	err = termx.Input("Senha:", &senha).
		Password().
		WithValidator(termx.MinLength(8)).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nSenha definida com sucesso!\n")
}