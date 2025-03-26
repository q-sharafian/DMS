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

**How to create docker image for the app:**
1) Create a docker image for the app:  
```
docker build -t dms .
```

**How to manage setup the project and its dependencies in Kubernetes:**
1) Init kubectl.

2) If you are using GitHub registry, create new token to have ability to pull the images from GitHub registry.  
Create new token with `read:packages` scope. To do that go to `https://github.com/settings/tokens/new?scopes=write:packages` page. After, set registry auth info in `deployment/secret.yml`.   
At the end, apply secret with `kubectl apply -f secret.yml` command.  

3) Run the following command to create a new secret that is used with docker registries.
```sh
kubectl create secret docker-registry registry-secret \
  --docker-server=REGISTRY_URL \
  --docker-username=REGISTRY_USERNAME \      
  --docker-password=REGISTRY_PASS \  
  --docker-email=REGISTRY_EMAIL
```

4) If you want to list images of the registry, Run the following command:
```sh
curl -H "Authorization: Bearer YOUR_PERSONAL_ACCESS_TOKEN" \
     -H "Accept: application/vnd.github.v3+json" \
     https://api.github.com/user/packages?package_type=container
```
Replace `YOUR_PERSONAL_ACCESS_TOKEN` with the PAT you created. (Use Tokens(classic))

5) Create a RSA public and private key pair to used for JWT.    
To do, run the following command in the project root:  
```sh
openssl genrsa -out certs/jwt_keypair.pem 2048
openssl rsa -in certs/jwt_keypair.pem -pubout -out certs/jwt_publickey.crt
openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in certs/jwt_keypair.pem -out certs/jwt_pkcs8.key
```
Then copy their values in `secret.yml`. `JWT_PRIVATE_KEY_FILE_PATH` in the `.env` file represents the contents of `JWT_PRIVATE_KEY` and `JWT_PUBLIC_KEY_FILE_PATH` in the `.env` file represents `JWT_PUBLIC_KEY`.  

6) Apply kubernetes secret and configmap resources:
```sh
kubectl apply -f configmap.yml -f secret.yml
```

7) Apply DMS deployment:
```sh
kubectl apply -f deployment.yml
```

8) Init persistent volume:
```sh
kubectl apply -f deployment/psql/persistent-vol.yaml
```