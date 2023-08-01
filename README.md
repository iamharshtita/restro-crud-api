# restro-crud-api
A CRUD application built around Golang.
Implemented:
a) Authentication System for login purposes.
b) An Authorization middleware to protect enpoints.
c) REST apis to communicate with the server and handled routes using the Go's MUX router.
d) Postgresql to store data.


# AIM:  

To create a simple microservice that will implement basic functionalities for a Restaurant Management System. 

# Functionalities: 

1-To view the entire menu of the restaurant (Anyone) 

2-To view a single food item (Anyone) 

3-To add items inside the menu (Admin) 

4-To delete a food item from the menu (Admin) 

5-To modify an existing food item (Admin) 

# Roles: 

(A) User: 

-A user can only view the entire menu or a particular food item using the REST call 

-Any attempt to add, delete, or modify the food menu will fail. 

-A user cannot use the login api in this current system because only an admin can login by passing his credentials to the request header. 

(B) Admin: 

-An Admin can call all the REST apis. 

-In order to modify the database, admin needs to authenticate himself using the login api. 

-Upon successful authentication, the admin will be granted a JWT. 

-The JWT can be used for authorization purpose to authorize the admin to make any kind of modifications to the database.

# Endpoints: 

1) /api/login: This endpoint will generate the jwt token for the admin. 

Request Payload: None 

Method: GET 

Request Headers: Username, Password 

Response: Status Code along with JWT token. 

 

2) /api/view: This endpoint will display the entire list of food items present in the database. 

Request Payload: None 

Method: GET 

Request Headers: None 

Response: Status Code along with the list of food items. 

 

3) /api/view/{name}: This endpoint will display the food item of the corresponding {name} passed in the URL and present in the database. 

Request Payload: None 

Method: GET 

Request Headers: None 

Response: Status Code along with the description of the food item. 

 

4) /api/add: This endpoint will allow only the 'admin' to add a food item into the database. 

Request Payload: Food Item in JSON format 

Method: POST 

Request Headers: JWT Token 

Response: Status Code along with an output message 

 

5) /api/delete: This endpoint allows the admin to delete a food item from the database. 

Request Payload: Food Item in JSON format 

Method: DELETE 

Request Headers: JWT Token 

Response: Status Code along with an output message 

 

 

 

6) /api/update: This endpoint will allow the admin to update an existing food item. 

Request Payload: Food Item in JSON format 

Method: PUT 

Request Headers: JWT Token 

Response: Status Code along with an output message 
