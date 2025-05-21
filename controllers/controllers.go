package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go/token"
	"log"
	"net/http"
	"time"

	"github.com/Ecom-go/database"
	"github.com/Ecom-go/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/tools/go/analysis/passes/defers"
)

var (
	UserCollection *mongo.Collection = database.UserData(database.Client, "Users" )
	ProductCollection *mongo.Collection = database.ProductData(databse.Client,"Products")
	Validate = validator.New()
)

func HashPassword(password string) string {
	bytes,err := bcrypt.GenerateFromPassword([]byte(password),14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword),[]byte(userPassword))
	
	if err != nil {
		return false,"invalid credentials"
	}
	return true,""
}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON{http.StatusBadRequest,gin.H{"error": err.Error()}}
		}

		validateErr := Validate.Struct(user	)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest,gin.H("error": validateErr))
			return
		}

		count , err := UserCollection.CountDocments(ctx,bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error": err})
			return 
		}
		if count>0 {
			c.JSON(http.StatusBadRequest,gin.H{"error": "user already exists"})
		}

		count ,err = UserCollection.CountDocments(ctx,bson.M{"phone": user.Phone	})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error": err})
			return
		}

		if count >0 {
			c.JSON(http.StatusBadRequest,gin.H{"error ": "phone already exists"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At,_  	= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.Updated_At,_	= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.ID				= primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token,refreshtoken,_ := generate.TokenGenerator(*user.Email,*user.First_Name,*user.Last_Name,*&user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart	= make([]models.ProductUser,0)
		user.Address_Details = make([]models.Address,0)
		user.Order_Status	= make([]models.Order,0)
		_ , inserterror := UserCollection.InsertOne(ctx,user)
		if inserterror != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error" : "the user did not get created"})
			return 
		}
		defer cancel()
		
		c.JSON(http.StatusCreated,"successfully created user")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		var user  models.User
		if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": err})
		return
		}
		err := UserCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Invalid email or password"})
			return
		}

		PasswordIsValid,msg := VerifyPassword(*user.Password,*founduser.Password)
		defer cancel()
		
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			fmt.Println(msg)
			return
		}

		token,refreshToken,_ := generate.TokenGenerator(*founduser.Email,*founduser.First_Name,*founduser.Last_Name, *founduser.User_ID )
		defer cancel()

		generate.UpdateAllTokens(token ,refreshToken ,founduser.User_ID)
		c.JSON(http.StatusFound,founduser)
	}
}

func ProductViewAdmin() gin.HandlerFunc {
	
}

func SearchProdcut() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var productList []models.Product 
		var ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		cursor ,err := ProductCollection.Find(ctx,bson.D{{}})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,"Something went wrong,plaese try again after some time")
			return
		}

		err = cursor.All(ctx,&productList )
		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}
		defer cursor.Close()

		if err:= cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400,"Invalid")
			return
		}

		defer cancel()
		c.IndentedJSON(200,productList)
 	}
}

func SearchProdcutByQuery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var searchProdcuts []models.Product
		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type","application/json")
			c.JSON(http.StatusNotFound,gin.H{"Error","invalid search index"})
			c.Abort()
			return
		}

		var ctx,cancel := context.WithTimeOut(context.Background(),100*time.Second)
		defer cancel()
		
		searchquerydb,err := ProductCollection.Find(ctx,bson.M{"product_name":bson.M{"$regex",queryParam}})

		if err != nil {
			c.IndentedJSON(404,"something went wrong while fetching the data")
			return
		}

		err = searchquerydb.All(ctx,&searchProdcuts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400,"Invalid")
			return
		}

		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err != nil {
			log.Printf(err)
			c.IndentedJSON(400,"invalide request")
			return
		}
		defer cancel()
		c.IndentedJSON(200,searchProdcuts)
	}
}
