package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("Can't find the product")
	ErrCantFindProducts   = errors.New("Can't Find the products")
	ErrUserIsNotValid     = errors.New("This user is not valid")
	ErrCantUpdateUser     = errors.New("Can't update the user")
	ErrCantRemoveItemCart = errors.New("Can't Remove this item form the cart")
	ErrCantGetItem        = errors.New("Can't get the item")
	ErrCantBuyCartItem    = errors.New("Can't buy the cart item")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuy() {

}
