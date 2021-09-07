**Learning GoLang with Graphql Server**

This is a very basic project. Which allows user account creation, login feature and stores data similar to bookmarks. 
It has following mutations and queries:

<details>
<summary><B>NonAuthMutation</B></summary>
<pre>
1. createUser(...): String!    //Create new user account
2. login(...): String!         //Login with existing user account
3. refreshToken(...): String!  //Refresh already logged in users auth token
</pre>
</details>

<details>
<summary><B>AuthorizedMutations</B></summary>
<pre>
1. createLink(...): Link!      //Create a new bookmark in mysql database linked with my account
2. updateLink(...): Link!      //Update a already created bookmark by the same user
3. deleteLink(...): String!    //Delete a bookmark created by the same user
</pre>
</details>

<details>
<summary><B>NonAuthQueries</B></summary>
<pre>
1. links: [Link!]!             //Get all bookmarks in database
</pre>
</details>

<details>
<summary><B>AuthorizedQueries</B></summary>
<pre>
1. users: [User!]!             //Get all users in the database
2. linkById(...): Link         //Get bookmark by bookmark id
3. linksByUserId(...): [Link!]!//Get bookmarks of particular user
4. myLinks: [Link!]!           //Get my bookmarks
</pre>
</details>
<BR>

**Installation**
To run the project locally set up docker:
1. Install mysql image in docker:
```js
docker pull --platform linux/amd64 mysql
```
applicable for mac m1 may change on depending on your machine
<br>
2. Run `mysql` in docker:
```js
docker run -p 3000:3306 --name mysql -e MYSQL_ROOT_PASSWORD=dbpass -e MYSQL_DATABASE=hackernews -d mysql:latest
```
where password for the database is `dbpass` if this is changed then make the changes in `./internal/pkg/db/mysql/mysql.go` at line 17 `db, err := sql.Open("mysql", "root:YOURPASSWORD@tcp(localhost)/mysql")`
3. Assuming you have Go already setup in the machine if not follow the guide over [here](https://golang.org/doc/install)
```js
go run server.go
```
