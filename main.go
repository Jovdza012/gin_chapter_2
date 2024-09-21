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
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipes []Recipe
var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

// swagger:parameters recipes newRecipe
type Recipe struct {
	//swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions"bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt"bson:"publishedAt"`
}

func init() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	ctx = context.Background()

	// Use the MONGO_URI from .env
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set in .env file or system environment")
	}

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping the database
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	// Get database name from environment variable
	dbName := os.Getenv("MONGO_DATABASE")
	if dbName == "" {
		log.Fatal("MONGO_DATABASE is not set in .env file or system environment")
	}

	// Initialize the database and collection
	db := client.Database(dbName)
	collection = db.Collection("recipes")

	log.Println("Connected to MongoDB")
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

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err = collection.InsertOne(ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, recipe)
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
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)
	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
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
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "name", Value: recipe.Name},
			{Key: "instructions", Value: recipe.Instructions},
			{Key: "ingredients", Value: recipe.Ingredients},
			{Key: "tags", Value: recipe.Tags}}},
		})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipehas been updated"})
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
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filter := bson.M{"_id": objID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe successfully deleted", "deletedCount": result.DeletedCount})
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
	filter := bson.D{{Key: "tags", Value: tag}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	// Iterate over the cursor and decode each document
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the results
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
