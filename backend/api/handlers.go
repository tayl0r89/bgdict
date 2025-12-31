package api

import (
	"log"
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

type PostBulkSearchBody struct {
	Queries []string `json:"queries" form:"queries" binding:"required"`
}

type BulkSearchResult struct {
	Results []*WordResult `json:"results"`
}

func BulkSearchHandler(repo WordRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body PostBulkSearchBody
		if err := c.ShouldBind(&body); err != nil {
			c.JSON(400, err.Error())
			return
		}

		results := make([]*WordResult, 0)
		for _, item := range body.Queries {
			found, search_err := repo.SearchWord(item)
			if search_err == nil && len(found) > 0 {
				results = append(results, found...)
			}
		}

		c.JSON(200, BulkSearchResult{Results: results})
	}
}

type searchParams struct {
	Query string `form:"query" json:"query" binding:"required"`
}

func SearchHandler(repo WordRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("HERE")
		var params searchParams = searchParams{}
		err := c.BindQuery(&params)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A suitable query was not provided."})
			return
		}

		res, finderr := repo.SearchWord(params.Query)

		log.Println(res)
		if finderr != nil {
			c.JSON(http.StatusInternalServerError, "Failed finding word")
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

type PostBulkByIdBody struct {
	Queries []int `json:"queries" form:"queries" binding:"required"`
}

type ByIdWordResult struct {
	Word            *Word             `json:"word"`
	DerivativeForms []*DerivativeForm `json:"derivativeForms"`
}

type BulkByIdResult struct {
	Results []*ByIdWordResult `json:"results"`
}

func BulkByIdHandler(repo WordRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body PostBulkByIdBody
		if err := c.ShouldBind(&body); err != nil {
			c.JSON(400, err.Error())
			return
		}

		results := make([]*ByIdWordResult, 0)
		for _, item := range body.Queries {
			found, search_err := repo.GetWordById(item)
			derived_found, derived_err := repo.GetDerivedForms(item)
			if search_err != nil {
				log.Println(search_err.Error())
				c.JSON(400, search_err.Error())
				return
			}
			if derived_err != nil {
				log.Println(derived_err.Error())
				c.JSON(400, derived_err.Error())
				return
			}
			results = append(results, &ByIdWordResult{
				Word:            found,
				DerivativeForms: derived_found,
			})
		}

		c.JSON(200, BulkByIdResult{Results: results})
	}
}
