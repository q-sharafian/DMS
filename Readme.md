A user could have multiple job positions (jp).
both a web page and a Windows app's user interface are considered views in the MVC (Model-View-Controller) architectural pattern.
Try write all of logs in the controllers package except some exceptations. For example, you can write Debug and Info logs in other packages too. But try to write error or fatal logs in controllers package.

If we're going to delete a row in database, it has to be a soft delete. That is, the row can't be deleted, just something like a label has to be used to show that it's been deleted.

The hierarchy tree is not connected tree. It means some users/job-positions could have not any parents or creators.

psql -U mohammad -d dms -h localhost -p 5432

jwt has two header field:   
```js
{
  "alg": "RS256",
  "typ": "JWT"
}
```  
and has some payload field   
```js
{
  // Note that times in JWT must be UTC timezone
  "sub": "1234567890", // Subject (e.g. user ID)
  "name": "John Doe", // Currently, we don't use it
  "admin": true, // Currently, we don't use it
  "iat": 1678886400, // Issued at timestamp (Unix timestamp)
  "exp": 1678890000, // Expiration timestamp (Unix timestamp)
  "jp": "345252", // Job position ID
  "login_id": "54"
}
```  
The app has stateful in login/logout. It means we store the details of sessions in database.

To build API documentaion, go to the root dir and then use this command:  
`swag init -g "./cmd/api/main.go" -o "./docs/api"`  
To see the documentation, go to path `/swagger/index.html`.  

All time zones must be UTC.


**How to setup the app:**  
1) Create a RSA public and private key pair to used for JWT.  
To do, run the following command in the project root:  
```
openssl genrsa -out certs/jwt_keypair.pem 2048
openssl rsa -in certs/jwt_keypair.pem -pubout -out certs/jwt_publickey.crt
openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in certs/jwt_keypair.pem -out certs/jwt_pkcs8.key
```

