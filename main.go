package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "FOOD_ITEMS"
)

type Food struct {
	ID       int    `json:"id"`
	Name     string `json:"foodName"`
	CATEGORY string `json:"category"`
	PRICE    int    `json:"price"`
}

func main() {

	db := dbConnect()

	defer db.Close()

	//Create Table if doesn't exist
	res, err := db.Exec("CREATE TABLE IF NOT EXISTS CATALOGUE (id SERIAL PRIMARY KEY, foodName TEXT, category TEXT, price INTEGER )")

	numofRow, _ := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	} else if numofRow != 0 {
		fmt.Println("Table Created")
	}

	router := mux.NewRouter()

	router.HandleFunc("/api", greet()).Methods(http.MethodGet)
	router.HandleFunc("/api/login", login()).Methods(http.MethodGet)
	router.HandleFunc("/api/viewMenu/{category}", viewMenuByCategory(db)).Methods(http.MethodGet)
	router.HandleFunc("/api/viewMenu", viewMenu(db)).Methods(http.MethodGet)
	router.HandleFunc("/api/add", AuthMiddleware(addMenu(db))).Methods(http.MethodPost)
	router.HandleFunc("/api/update", AuthMiddleware(updateFood(db))).Methods(http.MethodPut)
	router.HandleFunc("/api/delete", AuthMiddleware(deleteFood(db))).Methods(http.MethodDelete)

	http.ListenAndServe(":8080", jsonContentTypeMiddleware(router))

}

var (
	secretKey = []byte(goDotEnvVariable("SecretKey"))
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Harsh Tita",
			ExpiresAt: jwt.NewTime(float64(time.Now().Add(time.Minute * 10).Unix())),
		},
	}

	//Createthe token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AuthMiddleware is a middleware function to verify the JWT token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Missing authorization token")
			return
		}

		// Verify the token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return secretKey, nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Invalid authorization token")
			return
		}

		// Token is valid, call the next handler
		next.ServeHTTP(w, r)
	}
}


func viewMenu(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM catalogue")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		foods := []Food{}
		for rows.Next() {
			var f Food
			if err := rows.Scan(&f.ID, &f.Name, &f.CATEGORY, &f.PRICE); err != nil {
				log.Fatal(err)
			}
			foods = append(foods, f)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(foods)
	}
}

func deleteFood(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var f Food
		json.NewDecoder(r.Body).Decode(&f)

		idFromBody := f.ID

		rows, err := db.Query("Select * from catalogue where id=$1", idFromBody)

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()
		if !rows.Next() {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "No record found with given id")
			return
		}

		_, errrr := db.Exec("Delete from catalogue where id= $1", idFromBody)

		if errrr != nil {
			log.Fatal(errrr)
		} else {
			fmt.Fprintf(w, "Deleted food %s with id: %d, successfully", f.Name, idFromBody)
		}

		// Construct the ALTER SEQUENCE statement
		query := fmt.Sprintf("ALTER SEQUENCE catalogue_id_seq RESTART WITH %d", idFromBody)
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("ALTER statement executed successfully.")

	}

}

func updateFood(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var f Food
		json.NewDecoder(r.Body).Decode(&f)

		idFromBody := f.ID

		rows, err := db.Query("Select * from catalogue where id=$1", idFromBody)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		if !rows.Next() {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("No record found with given id")
			return
		}

		_, errrr := db.Exec("Update catalogue Set price=$1 where id= $2", f.PRICE, idFromBody)

		if errrr != nil {
			log.Fatal(errrr)
		} else {
			fmt.Printf("Updated food with id: %d, successfully", idFromBody)
		}

	}
}

func addMenu(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var f Food
		json.NewDecoder(r.Body).Decode(&f)

		rows, err := db.Query("Select * from catalogue where foodname=$1", f.Name)

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()
		if rows.Next() {

			var f Food
			if err := rows.Scan(&f.ID, &f.Name, &f.CATEGORY, &f.PRICE); err != nil {
				log.Fatal(err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "The food %s already exists with the id:%d ", f.Name, f.ID)
			return
		}

		res, err := db.Exec("INSERT INTO catalogue (foodname, category, price) VALUES ($1, $2, $3) RETURNING id", f.Name, f.CATEGORY, f.PRICE)
		if err != nil {
			log.Fatal(err)
		}

		if num, _ := res.RowsAffected(); num == 1 {
			fmt.Fprintf(w, "Food '%s' added in the database succesfully!", f.Name)
		} else if num == 0 {
			fmt.Fprintf(w, "Food '%s' could not be added in the database!", f.Name)
		}

		//json.NewEncoder(w).Encode(f)
	}
}

func viewMenuByCategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		category := vars["category"]
		rows, err := db.Query("Select * from catalogue where category=$1", category)

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		var numRows int = 0
		foods := []Food{}
		for rows.Next() {
			numRows++
			var f Food
			if err := rows.Scan(&f.ID, &f.Name, &f.CATEGORY, &f.PRICE); err != nil {
				log.Fatal(err)
			}
			foods = append(foods, f)
		}

		if numRows == 0 {
			fmt.Fprintf(w, "No record found for the given category: %s", category)
			return
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(foods)

	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func dbConnect() *sql.DB {

	password := goDotEnvVariable("password")
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")

	return db
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	val, isPresent:=os.LookupEnv(key)

	if !isPresent{
		log.Fatal("The key is not present in the .env file")
	} 

	return val
}
