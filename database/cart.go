package database

import (
	"context"
	"errors"
	"log"

	"github.com/Ecom-go/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrCantFindProduct    = errors.New("Can't find the product")
	ErrCantFindProducts   = errors.New("Can't Find the products")
	ErrUserIsNotValid     = errors.New("This user is not valid")
	ErrCantUpdateUser     = errors.New("Can't update the user")
	ErrCantRemoveItemCart = errors.New("Can't Remove this item form the cart")
	ErrCantGetItem        = errors.New("Can't get the item")
	ErrCantBuyCartItem    = errors.New("Can't buy the cart item")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb,err := prodCollection.Find(ctx,bson.M{"_id": productID})

	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser
	err = searchfromdb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return err
	}

	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIsNotValid
	}

	filtered := bson.D{primitive.E{Key: "_id",Value: id}}
	update 	 := bson.D{{Key: "$push",Value: bson.D{primitive.E{Key: "usercart",Value: bson.D{{Key: "$each",Value: productCart}}}}}}	

	_,err = userCollection.UpdateOne(ctx,filtered,update)

	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuy() {

}
