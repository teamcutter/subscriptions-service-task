package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/teamcutter/subscriptions-service-task/internal/handler"
	"github.com/teamcutter/subscriptions-service-task/internal/model"
)

var mockUUID1 uuid.UUID = uuid.New()
var mockUUID2 uuid.UUID = uuid.New()

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(sub *model.Subscription) error {
	args := m.Called(sub)
	return args.Error(0)
}

func (m *MockRepo) GetAll() ([]model.Subscription, error) {
	args := m.Called()
	return args.Get(0).([]model.Subscription), args.Error(1)
}

func (m *MockRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) TotalCost(userID, service, start, end string) (int, error) {
	args := m.Called(userID, service, start, end)
	return args.Int(0), args.Error(1)
}

func setupTest(t *testing.T) (*echo.Echo, *MockRepo, *handler.Handler) {
	e := echo.New()
	mockRepo := new(MockRepo)
	log := slog.Default()
	h := handler.NewHandler(mockRepo, log)
	return e, mockRepo, h
}

func TestCreate(t *testing.T) {
	e, repo, h := setupTest(t)

	sub := model.Subscription{UserID: mockUUID1, ServiceName: "Netflix"}
	repo.On("Create", mock.AnythingOfType("*model.Subscription")).Return(nil)

	body, _ := json.Marshal(sub)
	req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		repo.AssertCalled(t, "Create", mock.AnythingOfType("*model.Subscription"))
	}
}

func TestGetAll(t *testing.T) {
	e, repo, h := setupTest(t)

	expectedSubs := []model.Subscription{
		{ID: 1, UserID: mockUUID1, ServiceName: "Netflix"},
		{ID: 2, UserID: mockUUID2, ServiceName: "Spotify"},
	}
	repo.On("GetAll").Return(expectedSubs, nil)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.GetAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var got []model.Subscription
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got)) {
			assert.Len(t, got, 2)
			assert.Equal(t, "Netflix", got[0].ServiceName)
		}
		repo.AssertCalled(t, "GetAll")
	}
}

func TestDelete(t *testing.T) {
	e, repo, h := setupTest(t)

	repo.On("Delete", 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/subscriptions/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, h.Delete(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		repo.AssertCalled(t, "Delete", 1)
	}
}

func TestTotalCost(t *testing.T) {
	e, repo, h := setupTest(t)

	repo.On("TotalCost", mockUUID1.String(), "Netflix", "01-2024", "12-2024").Return(100, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/subscriptions/total?user=%s&service=Netflix&start=01-2024&end=12-2024", mockUUID1.String()), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.TotalCost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]int
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, 100, resp["total"])

		repo.AssertCalled(t, "TotalCost", mockUUID1.String(), "Netflix", "01-2024", "12-2024")
	}
}
