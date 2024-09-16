// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/BuildingDistributed-Applications-in-Gin.
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Jovan Milanovic <milanovic97jovan@gmail.com>
// GitHub: https://github.com/milanovic97
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal(file, &recipes)
}

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

// swagger:operation POST /recipes recipes createRecipe
// ---
// summary: Creates a new recipe
// produces:
// - application/json
// parameters:
//   - name: recipe
//     in: body
//     description: Recipe to be created
//     required: true
//     schema:
//     $ref: '#/definitions/Recipe'
//
// responses:
//
//	'200':
//	  description: Recipe created successfully
//	  schema:
//	    $ref: '#/definitions/Recipe'
//	'400':
//	  description: Invalid payload
//	'500':
//	  description: Internal server error
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// ---
// summary: Returns a list of recipes
// produces:
// - application/json
// responses:
//
//	'200':
//	  description: List of recipes
//	  schema:
//	    type: array
//	    items:
//	      $ref: '#/definitions/Recipe'
//	'500':
//	  description: Internal server error
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// ---
// summary: Updates an existing recipe
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe to update
//     required: true
//     type: string
//   - name: recipe
//     in: body
//     description: Updated recipe information
//     required: true
//     schema:
//     $ref: '#/definitions/Recipe'
//
// produces:
// - application/json
// responses:
//
//	'200':
//	  description: Recipe updated successfully
//	  schema:
//	    $ref: '#/definitions/Recipe'
//	'400':
//	  description: Invalid input
//	'404':
//	  description: Recipe not found
//	'500':
//	  description: Internal server error
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, r := range recipes {
		if r.ID == id {
			recipes[i] = recipe
			c.JSON(http.StatusOK, recipe)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// ---
// summary: Deletes an existing recipe
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe to delete
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: Recipe deleted successfully
//	'404':
//	  description: Recipe not found
//	'500':
//	  description: Internal server error
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	for i, r := range recipes {
		if r.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "recipe deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
}

// swagger:operation GET /recipes/search recipes searchRecipes
// ---
// summary: Searches for recipes by tag
// parameters:
//   - name: tag
//     in: query
//     description: Tag to search for
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//	  description: List of recipes matching the tag
//	  schema:
//	    type: array
//	    items:
//	      $ref: '#/definitions/Recipe'
//	'500':
//	  description: Internal server error
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	var results []Recipe
	for _, r := range recipes {
		for _, t := range r.Tags {
			if t == tag {
				results = append(results, r)
				break
			}
		}
	}
	c.JSON(http.StatusOK, results)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	router.Run()
}
