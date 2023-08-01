package controllers

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func DetermineBreed(pet models.Pet) string {
	if(strings.ToLower(pet.Animal) == "dog"){
		return pet.Breed
	} else {
		return "unknown"
	}
}

func CreatePet() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var pet models.Pet
        defer cancel()

        //validate the request body
        if err := c.BindJSON(&pet); err != nil {
            c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&pet); validationErr != nil {
            c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

				newSize := models.PetSize{
					Height: pet.Size.Height,
					Weight: pet.Size.Weight,
				}

        newUser := models.Pet{
            Id:       primitive.NewObjectID(),
            Name:    	pet.Name,
          	DOB: 			pet.DOB,
            Owner:    pet.Owner,
						Animal:		pet.Animal,
						Breed:		DetermineBreed(pet),
						Size:			newSize,
						Toy:			pet.Toy,
        }

        result, err := userCollection.InsertOne(ctx, newUser)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
    }
}

func GetPet() gin.HandlerFunc {
	return func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			userId := c.Param("userId")
			var pet models.Pet
			defer cancel()

			objId, _ := primitive.ObjectIDFromHex(userId)

			err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&pet)
			if err != nil {
					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					return
			}

			c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": pet}})
	}
}

func EditAPet() gin.HandlerFunc {
	return func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			userId := c.Param("userId")
			var pet models.Pet
			defer cancel()
			objId, _ := primitive.ObjectIDFromHex(userId)

			//validate the request body
			if err := c.BindJSON(&pet); err != nil {
					c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					return
			}

			//use the validator library to validate required fields
			if validationErr := validate.Struct(&pet); validationErr != nil {
					c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
					return
			}

			newSize := models.PetSize{
				Height: pet.Size.Height,
				Weight: pet.Size.Weight,
			}

			update := bson.M{"name": pet.Name, "dob": pet.DOB, "owner": pet.Owner, "animal": pet.Animal, "breed":DetermineBreed(pet), "size":newSize, "toy": pet.Toy}
			result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
			if err != nil {
					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					return
			}

			//get updated user details
			var updatedPet models.Pet
			if result.MatchedCount == 1 {
					err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedPet)
					if err != nil {
							c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
							return
					}
			}

			c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedPet}})
	}
}

func DeleteAPet() gin.HandlerFunc {
	return func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			userId := c.Param("userId")
			defer cancel()

			objId, _ := primitive.ObjectIDFromHex(userId)

			result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
			if err != nil {
					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					return
			}

			if result.DeletedCount < 1 {
					c.JSON(http.StatusNotFound,
							responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}},
					)
					return
			}

			c.JSON(http.StatusOK,
					responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}},
			)
	}
}

func GetAllPets() gin.HandlerFunc {
	return func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			var users []models.Pet
			defer cancel()

			results, err := userCollection.Find(ctx, bson.M{})

			if err != nil {
					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					return
			}

			//reading from the db in an optimal way
			defer results.Close(ctx)
			for results.Next(ctx) {
					var singlePet models.Pet
					if err = results.Decode(&singlePet); err != nil {
							c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					}

					users = append(users, singlePet)
			}

			c.JSON(http.StatusOK,
					responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}},
			)
	}
}