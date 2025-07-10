package models

import (
	"time"
	"gorm.io/gorm"
)

type Autor struct {
	gorm.Model
	Nome           string    `json:"nome"`
	DataNascimento time.Time `json:"data_nascimento"`
	Nacionalidade  string    `json:"nacionalidade"`
	Biografia      string    `json:"biografia"`
	Livros         []Livro   `json:"livros,omitempty"`
}

type Livro struct {
	gorm.Model
	Nome        string    `json:"nome"`
	ISBN        string    `json:"isbn"`
	Ano         int       `json:"ano"`
	Genero      string    `json:"genero"`
	AutorID     uint      `json:"autor_id"`
	Autor       Autor     `json:"autor,omitempty"`
	Colecoes    []Colecao `gorm:"many2many:colecao_livros;" json:"colecoes,omitempty"`
}

type Colecao struct {
	gorm.Model
	Nome      string  `json:"nome"`
	Descricao string  `json:"descricao"`
	Tema      string  `json:"tema"`
	Livros    []Livro `gorm:"many2many:colecao_livros;" json:"livros,omitempty"`
}