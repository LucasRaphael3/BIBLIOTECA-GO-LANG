package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/LucasRaphael3/biblioteca-api/internal/handlers"
)

func SetupRoutes(router *gin.Engine) {
	// Agrupamento de rotas para melhor organização
	api := router.Group("/api")
		// Rotas obrigatórias
		api.GET("/autores", handlers.ListarAutores)
		api.GET("/autores/:id/livros", handlers.ListarLivrosDeAutor)
		api.POST("/autores", handlers.CadastrarAutor)
		api.POST("/livros", handlers.CadastrarLivro)
		api.POST("/colecoes", handlers.CriarColecao)
		api.POST("/colecoes/:id/adicionar-livro", handlers.AdicionarLivroAColecao)
		api.GET("/colecoes/:id/livros", handlers.ListarLivrosDeColecao)

		extras := api.Group("/extras")
		{
			extras.GET("/livros/por-nacionalidade/:nacionalidade", handlers.ListarLivrosPorNacionalidade)
			extras.GET("/colecoes/por-tema/:tema", handlers.ListarColecoesPorTema)
			extras.GET("/autores/com-mais-de/:quantidade", handlers.ListarAutoresComMaisDeXLivros)
			extras.GET("/livros/publicados-em/:ano", handlers.ListarLivrosPorAno)
			extras.GET("/generos", handlers.ListarGenerosDistintos)
			extras.GET("/colecoes/:id/detalhes", handlers.DetalhesColecao)
			extras.GET("/estatisticas/nacionalidade", handlers.EstatisticasLivrosPorNacionalidade)
	}
}
