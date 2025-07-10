package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/LucasRaphael3/biblioteca-api/internal/database" // Importa a variável database.DB
	"github.com/LucasRaphael3/biblioteca-api/internal/models"   // Importa as structs
)

func ListarAutores(c *gin.Context) {
	var autores []models.Autor
	database.DB.Find(&autores)
	c.JSON(http.StatusOK, autores)
}

func ListarLivrosDeAutor(c *gin.Context) {
	autorID := c.Param("id")

	var livros []models.Livro
	if err := database.DB.Where("autor_id = ?", autorID).Find(&livros).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Livros não encontrados para este autor"})
		return
	}

	c.JSON(http.StatusOK, livros)
}

func CadastrarAutor(c *gin.Context) {
	var autor models.Autor
	// 1. Pega os dados JSON do corpo da requisição e coloca na variável 'autor'
	if err := c.ShouldBindJSON(&autor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Cria o registro do autor no banco de dados usando GORM
	if err := database.DB.Create(&autor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao cadastrar autor"})
		return
	}

	// 3. Retorna uma resposta de sucesso (201 Created) com os dados do autor criado
	c.JSON(http.StatusCreated, autor)
}

func CadastrarLivro(c *gin.Context) {
	var livro models.Livro
	if err := c.ShouldBindJSON(&livro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var autor models.Autor
	if err := database.DB.First(&autor, livro.AutorID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Autor com ID especificado não existe!"})
		return
	}

	if err := database.DB.Create(&livro).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar livro"})
		return
	}
	c.JSON(http.StatusCreated, livro)
}

func CriarColecao(c *gin.Context) {
	var colecao models.Colecao
	if err := c.ShouldBindJSON(&colecao); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&colecao).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar coleção"})
		return
	}
	c.JSON(http.StatusCreated, colecao)
}

func AdicionarLivroAColecao(c *gin.Context) {
	colecaoID := c.Param("id")

	var body struct {
		LivroID uint `json:"livro_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da requisição inválido, esperado {\"livro_id\": ...}"})
		return
	}

	var colecao models.Colecao
	if err := database.DB.First(&colecao, colecaoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coleção não encontrada"})
		return
	}

	var livro models.Livro
	if err := database.DB.First(&livro, body.LivroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
		return
	}

	database.DB.Model(&colecao).Association("Livros").Append(&livro)
	c.JSON(http.StatusOK, gin.H{"message": "Livro adicionado à coleção com sucesso"})
}

func ListarLivrosDeColecao(c *gin.Context) {
	colecaoID := c.Param("id")
	var colecao models.Colecao

	// Usamos Preload para carregar os livros associados a esta coleção
	if err := database.DB.Preload("Livros").First(&colecao, colecaoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coleção não encontrada"})
		return
	}

	c.JSON(http.StatusOK, colecao.Livros)
}

func ListarLivrosPorNacionalidade(c *gin.Context) {
	nacionalidade := c.Param("nacionalidade")
	var livros []models.Livro

	// Usamos Joins para filtrar livros baseados em um campo da tabela relacionada (autores)
	database.DB.Joins("Autor").Where("autores.nacionalidade = ?", nacionalidade).Find(&livros)

	c.JSON(http.StatusOK, livros)
}

// 2. GET /colecoes/por-tema/:tema
func ListarColecoesPorTema(c *gin.Context) {
	tema := c.Param("tema")
	var colecoes []models.Colecao
	database.DB.Where("tema = ?", tema).Find(&colecoes)
	c.JSON(http.StatusOK, colecoes)
}

// 3. GET /autores/com-mais-de/:quantidade
func ListarAutoresComMaisDeXLivros(c *gin.Context) {
    quantidadeStr := c.Param("quantidade")
    quantidade, err := strconv.Atoi(quantidadeStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Quantidade inválida"})
        return
    }

    var autores []models.Autor
    // Subconsulta para encontrar autor_id's que aparecem mais de 'quantidade' vezes na tabela de livros
    subQuery := database.DB.Model(&models.Livro{}).Select("autor_id").Group("autor_id").Having("COUNT(id) > ?", quantidade)
    // Busca autores cujos IDs estão na lista da subconsulta
    database.DB.Where("id IN (?)", subQuery).Find(&autores)

    c.JSON(http.StatusOK, autores)
}

// 4. GET /livros/publicados-em/:ano
func ListarLivrosPorAno(c *gin.Context) {
	anoStr := c.Param("ano")
	ano, err := strconv.Atoi(anoStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ano inválido"})
		return
	}
	var livros []models.Livro
	database.DB.Where("ano = ?", ano).Find(&livros)
	c.JSON(http.StatusOK, livros)
}

func ListarGenerosDistintos(c *gin.Context) {
	var generos []string
	// Usa Distinct e Pluck para pegar uma lista de valores únicos de uma coluna
	database.DB.Model(&models.Livro{}).Distinct().Pluck("genero", &generos)
	c.JSON(http.StatusOK, generos)
}


func DetalhesColecao(c *gin.Context) {
	colecaoID := c.Param("id")
	var colecao models.Colecao

	// Preload aninhado: Carrega os Livros da coleção e, para cada Livro, carrega seu Autor.
	if err := database.DB.Preload("Livros.Autor").First(&colecao, colecaoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coleção não encontrada"})
		return
	}
	c.JSON(http.StatusOK, colecao)
}

func EstatisticasLivrosPorNacionalidade(c *gin.Context) {
	// Struct anônimo para receber o resultado da consulta SQL
	type Resultado struct {
		Nacionalidade string
		TotalLivros   int
	}
	var resultados []Resultado

	// Query customizada com Joins, Group By e Count
	database.DB.Model(&models.Livro{}).
		Select("autores.nacionalidade, count(livros.id) as total_livros").
		Joins("join autores on autores.id = livros.autor_id").
		Group("autores.nacionalidade").
		Scan(&resultados)

	// Converte a lista de resultados em um mapa para o formato de resposta desejado
	mapaDeResultados := make(map[string]int)
	for _, res := range resultados {
		mapaDeResultados[res.Nacionalidade] = res.TotalLivros
	}

	c.JSON(http.StatusOK, mapaDeResultados)
}
