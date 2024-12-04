# About API
This is a **fake** REST API specifically built for my mobile application. It includes all the essential features required for a basic e-commerce application.
This API uses mongoDB as database.

## End points
### Login 
`/api/v1/user/login POST` request 
```json 
{
    "email": "example@gmail.com",
    "password":"examplePass"
}
``` 
Response
 ```json 
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOjE3MzA1NDQ0MDgsInVzZXJJRCI6IjY3MTY0ZmRmMGYzOTFhMWZiMGZkNjE3ZSJ9.aZ5P1NTltnjs1F-CKLXfsSOBg6nHvUcMzm_h6uZ7ss4",
    "userName": "mal yay",
    "email": "wiro@gmail.com"
}
```

### Register
`/api/v1/user/register POST` 
```json
{
    "userName": "example name",
    "email": "example@gmail.com",
    "password":"examplePass"
}
```
 Response
 ```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOjE3MzA1NDQ0MDgsInVzZXJJRCI6IjY3MTY0ZmRmMGYzOTFhMWZiMGZkNjE3ZSJ9.aZ5P1NTltnjs1F-CKLXfsSOBg6nHvUcMzm_h6uZ7ss4",
    "userName": "example name",
    "email": "example@gmail.com"
}
```

### OAuth
`/api/v1/user/oauth POST` 
```json
{
    "userName": "Berk OAuth",
    "email": "wiroOAuth@gmail.com",
    "oAuthId": "12434563563"
}
```
 Response
 ```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOjE3MzA1NDQ0MDgsInVzZXJJRCI6IjY3MTY0ZmRmMGYzOTFhMWZiMGZkNjE3ZSJ9.aZ5P1NTltnjs1F-CKLXfsSOBg6nHvUcMzm_h6uZ7ss4",
    "userName": "example name",
    "email": "example@gmail.com"
}
```

### Categories
Get all categories 
`/api/v1//categories/ GET`  <br> Response
 ```json
{
    "categories": [
    {
        "_id": "671652995b41c8ea613df136",
        "value": "headphone"
    }
]
}
```

### Shop
get items as filter <br>
this request will return us max 20 products <br>
so if we want to access more/diffrent products we have to set `startIndex` as much as we fetch products
`/api/v1/shop/featured/ GET`  request 
```json
{
    "category": "671652995b41c8ea613df136",
    "search": "",
    "startIndex": 0,
    "sortBy": "ascending/descending",
    "minPrice": 30,
    "maxPrice": 400
}
```
response
```json
[
    {
        "_id": "6717af4decddea45c1a064de",
        "pictureName": "headphone.jpg",
        "name": "TMA-2 HD Wireless",
        "price": 300
    },
]
```

### Basket fetch
`/api/v1/basket/ GET` 
response
```json
[
    {
        "_id": "67260023d0ee97422c630d53",
        "product": {
            "_id": "6717af4decddea45c1a064de",
            "pictureName": "headphone.jpg",
            "name": "TMA-2 HD Wireless",
            "price": 300
        },
        "count": 1
    }
]
```

### Basket Update
`/api/v1/basket/update POST` 
request
```json
{
    "basketId": "671fda1b69e5fdfac40ae895",
    "count": -1
}
```

### Basket Add To Basket
`/api/v1/basket/add POST` 
request
```json
{
    "productId": "6717af4decddea45c1a064de"
}
```
response <br>
id: id of the basket
```json
{
    "id": "6717af4decddea45c1a064de"
}
```

### Remove From basket
`/api/v1/basket/removeOne DELETE` 
request
```json
{
    "basketId": "671fda1b69e5fdfac40ae895"
}
```

### Clear Basket
`/api/v1/basket/removeAll DELETE` 

### IMAGE
```/api/v1/image/{ImageFile} GET```

## How do I use
download the project <br>
create a `.env` file and fill the variables 
```env
MONGO_URI="PASTE your mongoURI"
DB_NAME="DB name"
JWT_SECRET="SUPERSECRET"
```
create all the collections in ur DB and import jsons the correct collections <br>
the dumy datas can be found in: `/collectionDumData/` <br>
run `go run main.go` command inside the project directory
