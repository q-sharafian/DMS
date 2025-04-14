A user could have multiple job positions (jp).  
both a web page and a Windows app's user interface are considered views in the MVC (Model-View-Controller) architectural pattern.
Try write all of logs in the controllers package except some exceptations. For example, you can write Debug and Info logs in other packages too. But try to write error or fatal logs in controllers package.

If we're going to delete a row in database, it has to be a soft delete. That is, the row can't be deleted, just something like a label has to be used to show that it's been deleted.

The hierarchy tree is not connected tree. It means some job-positions could have not any parents or creators.
job-positions that have not any parents, could access all things. (docs, events, and etc)

psql -U mohammad -d dms -h localhost -p 5432

Download/Upload policy is simple; If the job position have access to an event, then he could upload/download files.

**Sample of upload/downlaod file request:**
`auth-token` structure is as this: `event-id:jwt:job-position-id`. Then it must be encoded with `base64`.  

- Upload:  
```json
{
  "auth-token": "your-token",
  "object-types":{
    "jpg": 10,
    "pdf": 1
  }
}
```

TODO: Set redis memory cleaning policy

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

**Notes on running the app in local:**
1) Run the following command in the project root to have an `.env` file and set its values:  
```sh
mv .env-template .env
```  
2) Init redis and store its details in the `.env` file.
```sh
docker run --name some-redis -d redis
```  
3) Run go wtih below command in the root directory:
```sh
go run cmd/api/main.go
```  


**How to create docker image for the app:**
1) Create a docker image for the app:  
```
docker build -t dms .
```

**How to manage setup the project and its dependencies in Kubernetes:**
1) Init kubectl.

2) If you are using GitHub registry, create new token to have ability to pull the images from GitHub registry.  
Create new token with `read:packages` scope. To do that go to `https://github.com/settings/tokens/new?scopes=write:packages` page. After, set registry auth info in `deployment/secret.yml`.   
At the end, apply secret with `kubectl apply -f deployment/secret.yml` command.  

3) Run the following command to create a new secret that is used with docker registries.
```sh
kubectl create secret docker-registry registry-secret \
  --docker-server=REGISTRY_URL \
  --docker-username=REGISTRY_USERNAME \      
  --docker-password=REGISTRY_PASS \  
  --docker-email=REGISTRY_EMAIL
```
4) If you want to list images of the GitHub Docker registry, Run the following command:
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
Then copy their values in `deployment/secret.yml`. `JWT_PRIVATE_KEY_FILE_PATH` in the `.env` file represents the contents of `JWT_PRIVATE_KEY` and `JWT_PUBLIC_KEY_FILE_PATH` in the `.env` file represents `JWT_PUBLIC_KEY`. In production, the `JWT_PUBLIC_KEY_FILE_PATH` and `JWT_PRIVATE_KEY_FILE_PATH` are useless.  

6) Apply kubernetes secret and configmap resources:
```sh
kubectl apply -f deployment/configmap.yml -f deployment/secret.yml
```  
7) Init persistent volume claim to use in PSQL deployment:
```sh
kubectl apply -f deployment/psql/psql-volume-claim.yml
```  
9) deploy PostgreSQL:
```sh
kubectl apply -f deployment/psql/psql-deployment.yml -f deployment/psql/psql-service.yml 
```
If you need to connect a URL to your Postgres service, run the following code. But note taht edit the file.  
```sh
kubectl apply -f deployment/psql/psql-ingress.yml
```  
10) Deploy Redis:
```sh
kubectl apply -f deployment/redis/redis-deployment.yml -f deployment/redis/redis-service.yml
```  
11) Apply DMS deployment:
```sh
kubectl apply -f deployment/dms/dms-deployment.yml -f deployment/dms/dms-service.yml -f deployment/dms/dms-ingress.yml
```
12) Init `file-transfer` container. First set env variables inside `deployment/file-transfer/file-transfer-configmap.yml` and `deployment/file-transfer/file-transfer-secret.yml` file. Then apply resources:
```sh
kubectl apply \
-f deployment/file-transfer/file-transfer-configmap.yml \
-f deployment/file-transfer/file-transfer-secret.yml \
-f deployment/file-transfer/file-transfer-deployment.yml \
-f deployment/file-transfer/file-transfer-service.yml \
-f deployment/file-transfer/file-transfer-ingress.yml
```
13) If you want to send HTTP request to the `file-transfer`, send the request to the url `http://host-addr:80/downloadOrUpload`.

get list of pods related to a deployment resource:
kubectl get pods --selector=app=<app-name>
kubectl get pods -l app=<app-name>
deleting a pod:
kubectl delete pods <pod-name>
Deleting an ingress:
kubectl delete ingress <ingress-name>