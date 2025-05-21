package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Ecom-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AddAddress() gin.HandlerFunc {

}

func EditAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid serach Index"})
			c.Abort()
			return
		}

		address := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "%Set", Value: bson.D{Key: "address", Value: address}}}

		_ ,err = UserCollection.UpdateOne(ctx,filter,update)
		if err != nil {
			c.IndentedJSON(404,"Wrong commond")
			return
		}
		defer 	cancel()
		ctx.Done()
		c.IndentedJSON(200,"Successfully deleted")
}
