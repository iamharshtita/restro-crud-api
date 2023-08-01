# restro-crud-api
A CRUD application built around Golang.
AIM:  

To create a simple microservice that will implement basic functionalities for a Restaurant Management System. 

Functionalities: 

1-To view the entire menu of the restaurant (Anyone) 

2-To view a single food item (Anyone) 

3-To add items inside the menu (Admin) 

4-To delete a food item from the menu (Admin) 

5-To modify an existing food item (Admin) 

Roles: 

(A) User: 

-A user can only view the entire menu or a particular food item using the REST call 

-Any attempt to add, delete, or modify the food menu will fail. 

-A user cannot use the login api in this current system because only an admin can login by passing his credentials to the request header. 

(B)Admin: 

-An Admin can call all the REST apis. 

-In order to modify the database, admin needs to authenticate himself using the login api. 

-Upon successful authentication, the admin will be granted a JWT. 

-The JWT can be used for authorization purpose to authorize the admin to make any kind of modifications to the database. 
