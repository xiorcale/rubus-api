package models

import (
	"net/http"
	"regexp"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Role is an enum which spcify the role of a `User`
type Role string

// Values for `Role` enum
const (
	EnumRoleAdmin Role = "administrator"
	EnumRoleUser  Role = "user"
)

// User is the `User` model in the database
type User struct {
	ID           int64  `json:"id" pg:",pk"`
	Username     string `json:"username" pg:",unique, notnull"`
	Email        string `json:"email" pg:",unique, notnull"`
	Role         Role   `json:"role"`
	PasswordHash string `json:"-" pg:",notnull"`
}

// NewUser is the model sent to create a new `User`
type NewUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

// PutUser is only use to document the PUT `User` endpoint
type PutUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Bind transforms the given payload into a `User`, with some validations
func (u *User) Bind(c echo.Context, cost int) *JSONError {
	newUser := &NewUser{}
	db := &echo.DefaultBinder{}
	if err := db.Bind(newUser, c); err != nil {
		return NewBadRequestError()
	}

	// fields validation
	if len(newUser.Username) == 0 {
		return &JSONError{
			Status: http.StatusBadRequest,
			Error:  "username is required.",
		}
	}

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(newUser.Email) {
		return &JSONError{
			Status: http.StatusBadRequest,
			Error:  "email address is not valid.",
		}
	}

	if len(newUser.Password) < 8 {
		return &JSONError{
			Status: http.StatusBadRequest,
			Error:  "password should be at least 8 characters.",
		}
	}

	if newUser.Role != EnumRoleAdmin {
		newUser.Role = EnumRoleUser
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), cost)

	// `NewUser` --> `User`
	u.Username = newUser.Username
	u.Email = newUser.Email
	u.PasswordHash = string(bytes)
	u.Role = newUser.Role

	return nil
}

// BindWithEmptyFields transforms the given payload into a `User`, with
// some validations, but does not require any field (they should be either nil
// or valid)
func (u *User) BindWithEmptyFields(c echo.Context, cost int) *JSONError {
	newUser := &NewUser{}
	db := &echo.DefaultBinder{}
	if err := db.Bind(newUser, c); err != nil {
		return NewBadRequestError()
	}

	if newUser.Username != "" {
		u.Username = newUser.Username
	}

	if newUser.Email != "" {
		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if !re.MatchString(newUser.Email) {
			return &JSONError{
				Status: http.StatusBadRequest,
				Error:  "email address is not valid.",
			}
		}
		u.Email = newUser.Email
	}

	if newUser.Password != "" {
		if len(newUser.Password) < 8 {
			return &JSONError{
				Status: http.StatusBadRequest,
				Error:  "password should be at least 8 characters.",
			}
		}
		bytes, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), cost)
		u.PasswordHash = string(bytes)
	}

	return nil
}

// AddUser inserts a new `User` into the database
func AddUser(db *pg.DB, user *User) *JSONError {
	if err := db.Insert(user); err != nil {
		if pgErr, ok := err.(pg.Error); ok && pgErr.IntegrityViolation() {
			return &JSONError{
				Status: http.StatusConflict,
				Error:  "username and/or email already exist(s).",
			}
		}
		return NewInternalServerError()
	}

	return nil
}

// GetUser returns the `User` with the given `uid` from the database
func GetUser(db *pg.DB, uid int64) (*User, *JSONError) {
	user := &User{ID: uid}
	if err := db.Select(user); err != nil {
		if err == pg.ErrNoRows {
			return nil, &JSONError{
				Status: http.StatusNotFound,
				Error:  "user does not exist.",
			}
		}
		return nil, NewInternalServerError()
	}

	return user, nil
}

// GetAllUsers returns all the `User` from the database
func GetAllUsers(db *pg.DB) (users []*User, jsonErr *JSONError) {
	if err := db.Model(users).Select(); err != nil {
		return nil, NewInternalServerError()
	}

	return users, nil
}

// UpdateUser modifies the `User` with the given `uid` in the database, with some validations
func UpdateUser(db *pg.DB, uid int64, uu *User) (u *User, jsonErr *JSONError) {
	u, jsonErr = GetUser(db, uid)
	if jsonErr != nil {
		return nil, jsonErr
	}

	if uu.Username != "" {
		u.Username = uu.Username
	}
	if uu.Email != "" {
		u.Email = uu.Email
	}
	if uu.PasswordHash != "" {
		u.PasswordHash = uu.PasswordHash
	}

	if err := db.Update(u); err != nil {
		if pgErr, ok := err.(pg.Error); ok && pgErr.IntegrityViolation() {
			return nil, &JSONError{
				Status: http.StatusConflict,
				Error:  "username and/or email already exist(s).",
			}
		}
		return nil, NewInternalServerError()
	}

	return u, nil
}

// DeleteUser removes the given Rubus `User` from the database
func DeleteUser(db *pg.DB, uid int64) *JSONError {
	user := &User{ID: uid}
	if err := db.Delete(user); err != nil {
		if err == pg.ErrNoRows {
			return &JSONError{
				Status: http.StatusNotFound,
				Error:  "user does not exist.",
			}
		}
		return NewInternalServerError()
	}

	return nil
}

// Login checks if the given credentials are valid or not
func Login(db *pg.DB, username, password string) (*int64, *Role, bool) {
	user := &User{}

	if err := db.Model(user).Where("username = ?", username).Select(); err != nil {
		if err == pg.ErrNoRows {
			return nil, nil, false
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, false
	}

	return &user.ID, &user.Role, true
}
