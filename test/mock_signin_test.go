package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"userPage/controller"
	"userPage/database"
	"userPage/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDBSignin(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn: mockDB,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}

	database.SetDB(db)
	mock.ExpectExec("DELETE FROM users").WillReturnResult(sqlmock.NewResult(0, 0))
	database.DB.Exec("DELETE FROM users")
	password, _ := bcrypt.GenerateFromPassword([]byte("user@123"), bcrypt.DefaultCost)
	testUser := models.Users{Name: "user", Email: "user@gmail.com", Password: string(password)}
	if err := database.DB.Create(&testUser).Error; err != nil {
		fmt.Println("-----------", err)
	}

	cleanup := func() {
		mockDB.Close()
	}
	return mock, cleanup
}

func TestSingnIn(t *testing.T) {
	mock, cleanup := SetupTestDBSignin(t)
	defer cleanup()

	gin.SetMode(gin.TestMode)

	t.Run("successful login", func(t *testing.T) {
		password, _ := bcrypt.GenerateFromPassword([]byte("user@123"), bcrypt.DefaultCost)

		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email=\\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
			WithArgs("user@gmail.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(1, "user", "user@gmail.com", password))

		router := gin.Default()
		router.POST("/signin", controller.SignIn)
		loginInput := controller.SignInInput{
			Email:    "user@gmail.com",
			Password: "user@123",
		}
		jsonValue, _ := json.Marshal(loginInput)
		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Login successful")
	})
	t.Run("wrong email", func(t *testing.T) {

		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email=\\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
			WithArgs("wrong@gmail.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}))

		router := gin.Default()
		router.POST("/signin", controller.SignIn)
		loginInput := controller.SignInInput{
			Email:    "wrong@gmail.com",
			Password: "user@123",
		}
		jsonValue, _ := json.Marshal(loginInput)
		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")
	})
	t.Run("wrong password", func(t *testing.T) {
		password, _ := bcrypt.GenerateFromPassword([]byte("user@123"), bcrypt.DefaultCost)
		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email=\\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
			WithArgs("user@gmail.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(1, "user", "user@gmail.com", password))

		router := gin.Default()
		router.POST("/signin", controller.SignIn)
		loginInput := controller.SignInInput{
			Email:    "user@gmail.com",
			Password: "wrong@123",
		}
		jsonValue, _ := json.Marshal(loginInput)
		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")
	})
}
