package test

import (
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
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func UserListMock() (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic("failed to open sqlmock database connection")
	}
	dialecter := postgres.New(postgres.Config{
		Conn: mockDB,
	})
	gormDB, err := gorm.Open(dialecter, &gorm.Config{})
	if err != nil {
		panic("failed to open gorm db connection")
	}
	database.SetDB(gormDB)
	users := []models.Users{
		{Name: "user1", Email: "user1@gmail.com", Password: "user1@123"},
		{Name: "user2", Email: "user2@gmail.com", Password: "user2@123"},
	}
	mock.ExpectExec("DELETE FROM users").WillReturnResult(sqlmock.NewResult(0, 0))
	database.DB.Exec("DELETE FROM users")
	if err := database.DB.Create(&users).Error; err != nil {
		fmt.Println("-----------", err)
	}

	cleanup := func() {
		mockDB.Close()
	}
	return mock, cleanup
}
func TestUserList(t *testing.T) {
	mock, cleanup := UserListMock()
	defer cleanup()
	gin.SetMode(gin.TestMode)
	t.Run("Successfull fetch all users", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL").WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(1, "user1", "user1@gmail.com", "user1@123").
				AddRow(2, "user2", "user2@gmail.com", "user2@123"),
		)
		router := gin.Default()
		router.GET("/userlist", controller.UserList)
		req, _ := http.NewRequest("GET", "/userlist", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string][]models.Users
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response["users"], 2)
		assert.Equal(t, "user1", response["users"][0].Name)
		assert.Equal(t, "user1@gmail.com", response["users"][0].Email)
		assert.Equal(t, "user2", response["users"][1].Name)
		assert.Equal(t, "user2@gmail.com", response["users"][1].Email)
	})
	t.Run("Failure to fetch users", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL").
			WillReturnError(gorm.ErrInvalidTransaction)

		router := gin.Default()
		router.GET("/userlist", controller.UserList)

		req, _ := http.NewRequest("GET", "/userlist", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Failed to fetch users", response["error"])

	})

}
