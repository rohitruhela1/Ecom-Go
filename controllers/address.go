package controllers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os/user"
	"time"

	"github.com/Ecom-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AddAddress() gin.HandlerFunc {
	user_id = c.Query("id")
	return func(ctx *gin.Context) {
		if user_id =="" {
			c.Header("Content-Type","application/json")
			c.JSON(http.StatusNotFound,gin.H{"error","invalid user id"})
			c.Abort()
			return
		}
		
		address,err := ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500,"internal server error")
		}

		var addresses models.Address

		addresses.Address_ID = primitive.NewObjectID(user)
		if err := c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable,err.Error())	
		}

		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second )

		match_filter := bson.D{Key:"$match", Value: bson.D{primitive.E{Key: "_id",Value: address}} }
		unwind := bson.D{{Key: "$unwind",Value: bson.D{primitive.E{Key: "$path",Value: "$address"}}}}
		group := bson.D{{Key: "$group",Value: bson.D{primitive.E{Key: "_id",Value: "$address_id"},{Key: "count",Value: bson.D{primitive.E{Key: "$sum",Value: 1}}}}}}
		
		pointcursor ,err := UserCollection.Aggregate(match_filter,unwind,group)

		if err != nil {
			c.IndentedJSON(500,"Internal server error")
		}

		var addressinfo []bson.M
		if err = pointcursor.All(ctx,&addressinfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addreaddressinfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size <2 {
			filter := bson.D{primitive.E{Key: "_id",Value: address}}
			update := bson.D{{Key: "$push",Value: bson.D{primitive.E{Key: "address",Value: addresses}}}}
			_,err := UserCollection.UpdateOne(ctx,filter,update)
			if err != nil {
				fmt.Println(err)
			}
		}else{
			c.IndentedJSON(400,"Not Allowed")
		}
		defer cancel()
		ctx.Done()
	}
}

func EditAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid user id"})
			c.Abort()
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest,err.Error())
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set",Value: bson.D{primitive.E{Key: "address.0.house_name",Value: editaddress.House},{Key: "address.0.street_name",Value: editaddress.Street},{Key: "address.0.city_name",Value: editaddress.City},{Key: "address.0.pin_code",Value: editaddress.Pincode}}}}
		_,err = UserCollection.UpdateOne(ctx,filter,update)
		if err != nil {
			c.IndentedJSON(500,"Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200,"successfully updated the home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid user id"})
			c.Abort()
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest,err.Error())
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set",Value: bson.D{primitive.E{Key: "address.1.house_name",Value: editaddress.House},{Key: "address.1.street_name",Value: editaddress.Street},{Key: "address.1.city_name",Value: editaddress.City},{Key: "address.1.pin_code",Value: editaddress.Pincode}}}}
		_,err = UserCollection.UpdateOne(ctx,filter,update)

		if err != nil {
			c.IndentedJSON(500,"Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200,"successfully updated the work address")
	}
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
