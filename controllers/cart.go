package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ecom-go/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct{
	prodCollection *monogo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection,userCollection *mongo.Collection) *Application{
	return &Application{
		prodCollection: prodCollection,
		usercouserCollection: userCollection
	}
}

func (app *Application) AddToCart() gin.Handler {
	return func (c *gin.Context)  {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest,errors.New("product is empty"))
			return
		}

		userQueryId := c.Query("userID")

		if userQueryId == "" {
			log.Println("User id is empty")
			
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("userid is empty"))
			return
		}

		productID,err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil  {
			log.Println(err)

			_ = c.AbortWithError(http.StatusInternalServerError)

			return
		}

		var ctx,cancel := context.WithTimeOut(context.Background(),5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx,app.prodCollection,app.userCollection,productID,userQueryId)

		if err != nil  {
			c.IndentedJSON(http.StatusInternalServerError,err)
		}
		c.IndentedJSON(200,"successfully added the product")
	}
}

func (app *Application) RemoteItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest,errors.New("product is empty"))
			return
		}

		userQueryId := c.Query("userID")

		if userQueryId == "" {
			log.Println("User id is empty")
			
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("userid is empty"))
			return
		}

		productID,err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil  {
			log.Println(err)

			_ = c.AbortWithError(http.StatusInternalServerError)

			return
		}

		var ctx,cancel := context.WithTimeOut(context.Background(),5*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx,app.prodCollection,app.userCollection,productID,userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,err)
			return
		}
		c.IndentedJSON(200,"successfully removed item from the cart")
	}
}

func GetItemFromCart() gin.HandlerFunc {
	
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userQueryId := c.Query("id")

		if userQueryId == "" {
			log.Panicln("user id  is empty")
			c.AbortWithError(http.StatusBadRequest,errors.New("Userid is empty"))
			return
		}

		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		err = database.BuyItemFromCart(ctx,app.userCollection,userQueryId)	

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,err)
		}
		c.IndentedJSON("Successfully placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest,errors.New("product is empty"))
			return
		}

		userQueryId := c.Query("userID")

		if userQueryId == "" {
			log.Println("User id is empty")
			
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("userid is empty"))
			return
		}

		productID,err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil  {
			log.Println(err)

			_ = c.AbortWithError(http.StatusInternalServerError)

			return
		}

		var ctx,cancel := context.WithTimeOut(context.Background(),5*time.Second)
		defer cancel()

		err = database.InstantBuy(ctx,app.prodCollection,app.userCollection,productID,userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,err)
		}

		c.IndentedJSON(200, "successfully placed the order")
	}
}
