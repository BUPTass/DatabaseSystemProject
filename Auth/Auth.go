package Auth

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var RegErr = errors.New("username is already in use")
var UserNotFoundErr = errors.New("user not found")
var InvalidUsrErr = errors.New("invalid username or password")

// User represents a user in the database
type User struct {
	ID       int
	Username string
	Password string
}

// RegisterHandler handles the registration of a new user
func RegisterHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	db, err := sql.Open("mysql", "root:1taNWY1vXdTc4_-j@tcp(127.0.0.1:3306)/LTE")
	// Check if username is already in use
	user, err := GetUserByUsername(db, username)
	if err == nil || user.ID != 0 {
		return c.String(http.StatusBadRequest, RegErr.Error())
	}

	// Generate salted hash of password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Insert new user into database
	user = User{Username: username, Password: string(hashedPassword)}
	err = InsertUserReq(db, &user)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "Register Request Received!")
}

// LoginHandler handles the authentication of a user
func LoginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	db, err := sql.Open("mysql", "root:1taNWY1vXdTc4_-j@tcp(127.0.0.1:3306)/LTE")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// Get user from database
	user, err := GetUserByUsername(db, username)
	if err != nil || user.ID == 0 {
		return c.String(http.StatusBadRequest, InvalidUsrErr.Error())
	}
	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.String(http.StatusBadRequest, InvalidUsrErr.Error())
	}
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // One week
		HttpOnly: true,
	}
	sess.Values["username"] = user.Username
	sess.Values["expire"] = time.Now().Unix() + 86400*7 // Expire after one week
	sess.Save(c.Request(), c.Response())
	return c.Redirect(301, "/query")
}

func LogoutHandler(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	sess.Values["username"] = ""
	sess.Values["expire"] = 0
	sess.Save(c.Request(), c.Response())
	return c.Redirect(301, "/login")
}

func GetUserByUsername(db *sql.DB, username string) (User, error) {
	query := `SELECT id, username, password FROM users WHERE username = ?`
	row := db.QueryRow(query, username)

	user := User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, UserNotFoundErr
		}
		log.Println(err)
		return user, err
	}

	return user, nil
}

// InsertUserReq inserts a new user into the database
func InsertUserReq(db *sql.DB, user *User) error {
	// Prepare the SQL statement
	statement := "INSERT INTO users (username, password) VALUES (?, ?)"
	stmt, err := db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(user.Username, user.Password)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func InsertUserByAdmin(db *sql.DB, user *User) error {
	// Prepare the SQL statement
	statement := "INSERT INTO users (username, password, confirmed) VALUES (?, ?, TRUE)"
	stmt, err := db.Prepare(statement)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(user.Username, user.Password)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CheckSession(sess *sessions.Session) bool {
	if _, ok := sess.Values["username"]; ok {
		if _, ok := sess.Values["expire"]; ok {
			if sess.Values["expire"].(int64) < time.Now().Unix() {
				return true
			}
		}
	}
	return false
}
