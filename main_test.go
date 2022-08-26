package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func Test_handleQueryAllWebsiteStatusEmptyList(t *testing.T) {
	mockResponse := `[]`
	r := SetUpRouter()
	r.GET("/websites", handleQueryAllWebsitesStatus)
	req, _ := http.NewRequest("GET", "/websites", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_handlePostWebsitesList(t *testing.T) {
	testSitesList := websiteListStruct{
		SiteList: []string{"https://xyzw.zyx", "https://youtube.com", "https://google.com"},
	}
	r := SetUpRouter()
	r.POST("/websites", handlePostWebsitesList)
	jsonValue, _ := json.Marshal(testSitesList)
	req, _ := http.NewRequest("POST", "/websites", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	testSiteList2 := `[]string{"https://xyzw.zyx", "https://youtube.com", "https://google.com"}`
	jsonValue, _ = json.Marshal(testSiteList2)
	req, _ = http.NewRequest("POST", "/websites", bytes.NewBuffer(jsonValue))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	respStr := string(responseData)
	mockResp1 := `{"message":"website list not updated"}`
	assert.Equal(t, mockResp1, respStr)

}

func Test_handleQueryAllWebsitesStatus(t *testing.T) {
	testSitesList := websiteListStruct{
		SiteList: []string{"https://xyzw.zyx", "https://google.com"},
	}
	r := SetUpRouter()
	r.POST("/websites", handlePostWebsitesList)
	jsonValue, _ := json.Marshal(testSitesList)
	req, _ := http.NewRequest("POST", "/websites", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	r.GET("/websites", handleQueryAllWebsitesStatus)
	req, _ = http.NewRequest("GET", "/websites", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_handleQueryParticularWebsitesStatus(t *testing.T) {
	testSitesList := websiteListStruct{
		SiteList: []string{"https://xyzw.zyx", "https://google.com"},
	}
	r := SetUpRouter()
	r.POST("/websites", handlePostWebsitesList)
	jsonValue, _ := json.Marshal(testSitesList)
	req, _ := http.NewRequest("POST", "/websites", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	r.GET("/query", handleQueryParticularWebsiteStatus)

	req, _ = http.NewRequest("GET", "/query?site=https://xyzw.zyx", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	responseData, _ := ioutil.ReadAll(w.Body)
	respStr := string(responseData)
	if !strings.Contains(respStr, "UP") && !strings.Contains(respStr, "DOWN") && !strings.Contains(respStr, "WAIT") {
		t.Fail()
	}

	req, _ = http.NewRequest("GET", "/query?site=abc.cba", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	responseData, _ = ioutil.ReadAll(w.Body)
	respStr = string(responseData)
	mockResp1 := `{"message":"This website not in the database, please add first"}`
	assert.Equal(t, mockResp1, respStr)

	req, _ = http.NewRequest("GET", "/query", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	responseData, _ = ioutil.ReadAll(w.Body)
	respStr = string(responseData)
	mockResp2 := `{"message":"Query param not found, enter valid url string"}`
	assert.Equal(t, mockResp2, respStr)
}
