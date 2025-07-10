package main

import (
	"github.com/gin-gonic/gin"
	"github.com/LucasRaphael3/biblioteca-api/internal/database"
	"github.com/LucasRaphael3/biblioteca-api/internal/routes"
)

func main() {
	// 1. Conecta ao banco de dados e roda as migrações
	database.ConectarBancoDeDados()

	// 2. Cria o roteador Gin
	r := gin.Default()

	// 3. Configura as rotas da aplicação
	routes.SetupRoutes(r)

	// 4. Inicia o servidor
	r.Run(":8080")
}