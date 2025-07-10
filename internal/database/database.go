package database

import (
	"fmt"
	"log"

	"github.com/LucasRaphael3/biblioteca-api/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConectarBancoDeDados() {
	var err error
	// Caminho para o banco de dados dentro do contêiner Docker
	dbPath := "biblioteca.db"
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")

	// Roda as migrações, usando as structs do pacote models
	err = DB.AutoMigrate(&models.Autor{}, &models.Livro{}, &models.Colecao{})
	if err != nil {
		log.Fatalf("Falha ao migrar o banco de dados: %v", err)
	}
}