package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Ecom-go/database"
	"github.com/Ecom-go/tokens"
	jwt "github.com/dgrijalwa/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClamins
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email, firstname, lastname, uid string) (SignedToken string, signedrefreshtoken string, err error) {

	claims := &SignedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		Uid:        uid,
		StandardClamins: jwt.StandardClamins{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshclaims := &SignedDetails{
		StandardClamins: jwt.StandardClamins{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshtoken, err
}

func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := tokens.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is not valid"
		return
	}

	if claims.ExpriresAt < time.Now().Local().Unix() {
		msg = "the token is already exprired"
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var updatedobj primitive.D
	updatedobj = append(updatedobj, bson.E{Key: "token", Value: signedtoken})
	updatedobj = append(updatedobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})
	updated_at := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedobj = append(updatedobj, bson.E{Key: "upda ted_at", Value: updated_at})

	upsert := true

	filter := bson.M{"user_id", userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updatedobj},
	}, &opt)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
}
