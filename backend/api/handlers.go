package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWord(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Word{Id: 1, Name: "test"})
}

type findWordParams struct {
	Query string `form:"query" json:"query" binding:"required"`
}

func FindWordHandler(repo WordRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params findWordParams = findWordParams{}
		err := c.BindQuery(&params)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A suitable query was not provided."})
			return
		}

		res, finderr := repo.FindWords(params.Query)

		if finderr != nil {
			c.JSON(404, gin.H{"error": finderr.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	}
}

type getWordByIdParams struct {
	Id int `form:"id" json:"id" binding:"required"`
}

func GetWordByIdHandler(repo WordRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params getWordByIdParams = getWordByIdParams{}
		err := c.BindQuery(params)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No suitable id was provided."})
			return
		}

		res, getErr := repo.GetWordById(params.Id)

		if getErr != nil {
			c.JSON(404, "Failed to find word for id")
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	}
}
