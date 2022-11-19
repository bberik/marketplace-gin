# Marketplace Platform (Multi-vendor shop) API with Gin (Golang web framework) and JWT authentification.

To run the app: 

1. clone this repository to your machine
2. change your directory to marketplace-gin/server
3. ``` go run .```

OR 

1. ```docker run bberik/marketplace-gin```

## Routes:

  - Product Routes:
  
 ![image](https://user-images.githubusercontent.com/85312257/202837394-11b609eb-f5fb-42c9-8570-5b7859616d5a.png)

  - User Routes:
 
 ![image](https://user-images.githubusercontent.com/85312257/202837435-c51eec39-c29e-45fb-b25e-c3520241c600.png)

  - Cart Routes:
  
 ![image](https://user-images.githubusercontent.com/85312257/202837457-a8ff0715-9b4a-4698-a5c4-01c0a7bdb692.png)

  - Order Routes:
  
 ![image](https://user-images.githubusercontent.com/85312257/202837481-476a03a8-b017-4620-bafe-a9041eb1d1c8.png)

  - Shop Routes:
  
 ![image](https://user-images.githubusercontent.com/85312257/202837501-64e6afc4-2134-4474-91c9-6d8d1a63760a.png)


## Examples of request bodies:

  - User object:
  ```
{
    "email": "berik@gmail.com",
    "fullName": "Berik Bazarbayev",
    "username": "berik23",
    "password": "test1234"
}
```

  - To sign in:
```
{
    "email": "berik@gmail.com",
    "password": "test1234"
}
```

  - Shop object:
```
{
    "userId": "6374b5ea3c0201a319270cac",
    "shopName": "Berik's Shop",
    "shopAddress":
        {
            "details": "B",
            "building": "221",
            "street": "Baker",
            "area":
                {
                    "city": "London",
                    "country": "UK"
                }
        }
}
```

  - Address object:
```
{
    "details": "no details",
    "building": "149",
    "street": "Daehak-ro",
    "area":
        {
            "city": "Daejeon",
            "country": "Korea"
        }
}
```

  - Content object:
```
{
    "pID" : "63757a935c8d6c0d92d7aaf6",
    "itemID" : "63757a935c8d6c0d92d7aaf7",
    "quantity" : 5
}
```

  - Order object: 
```
{
    orderContent: 
            {
                "pID" : "63757a935c8d6c0d92d7aaf6",
                "itemID" : "63757a935c8d6c0d92d7aaf8",
                "quantity" : 5
            }
    shippingAddress: 
                {
                    "details": "no details",
                    "building": "149",
                    "street": "Kabanbai batyr",
                    "area":
                        {
                            "city": "Astana",
                            "country": "Kazakhstan"
                        }
                }
}
```

  - Status object:
```
{
    "status": "Order Shipped"
}
```
    
## Access:

- Guest (unauthorized) users:
  - can view al products
  - can view particular product
  - can view products by category
  - can search products
  - can sign in
  - can sign up
- Signed in (authorized) users:
  - all above
  - can add/update/delete products in/to cart
  - can buy products directly
  - can buy products from cart
  - can add/edit/update addresses
  - can view order history
  - can monitor order status
  - can update/delete account
- Sellers:
  - all above
  - can add/update/delete products
  - can view list of orders their shop received
  - can change order status (shipped, canceled, processed and etc.)

## Models used to design the system:

### Product:
  - productID
  - productName
  - description
  - category
  - images (array of URLs)
  - items (array) :
    - itemID
    - color
    - size (can be used as options)
    - instock (amount that can be ordered)
    - price
    - discount

### Content (used to wrap items of products into cart or order):
  - productID
  - itemID
  - quantity

### User:
  - userID
  - username
  - fullname
  - email
  - password
  - addresses (array)
  - cart (array of contents)
  - orders (array of orders)
 
### Shop:
  - userID (owner of shop)
  - shop name
  - shop address
  - products (array of products shop sells)
  - orders (array of orders shop received)
 
 ### Order:
  - orderID
  - order content
  - shipping address
  - order status (array of statuses):
    - status
    - time status created at

### Address:
  - details
  - building
  - street
  - area:
    - city
    - country
