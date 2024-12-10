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
    "token": "-",
    "refreshToken": "-",
    "userName": "example name",
    "email": "example@gmail.com"
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
    "token": "-",
    "refreshToken": "-",
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
    "token": "-",
    "refreshToken": "-",
    "userName": "example name",
    "email": "example@gmail.com"
}
```

### login with refreshToken
`/api/v1/user/oauth POST` <br>
`header`<br>
`Authorization` = `refresh token`
 Response
 ```json
{
    "token": "-",
    "refreshToken": "-",
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
        "id": "671652995b41c8ea613df136",
        "value": "headphone"
    }
]
}
```

### Shop
get items as filter <br>
this request will return us max 10 products <br>
`/api/v1/shop/featured?search=&offset= GET`  request 
```json
{
    "category": "671652995b41c8ea613df136", // additional
    "sortBy": "ascending/descending",
    "minPrice": 30,
    "maxPrice": 400
}
```
response
```json
{
    "products" :[
        {
            "id": "6717af4decddea45c1a064de",
            "pictureName": "headphone.jpg",
            "name": "TMA-2 HD Wireless",
            "price": 300
        },
    ]
}
```

### Basket fetch
`/api/v1/basket/ GET` 
response
```json
{
    "basket": [
        {
            "id": "67260023d0ee97422c630d53",
            "product": {
                "id": "6717af4decddea45c1a064de",
                "pictureName": "headphone.jpg",
                "name": "TMA-2 HD Wireless",
                "price": 300
            },
            "count": 1
        }
    ]
}
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
