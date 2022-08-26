package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type websiteListStruct struct {
	SiteList []string
}

type websiteStatusResp struct {
	Site   string
	Status string
}

type StatusChecker interface {
	CheckWebsiteStatus()
}

type statusCheckerStruct struct{}

var siteStatusMap = map[string]string{}
var statusCheckerObject = statusCheckerStruct{}

func (s statusCheckerStruct) CheckWebsiteStatusUtil(site string, m *sync.Mutex, w *sync.WaitGroup) {
	log.Printf("CHECKING SITE: %v", site)
	resp, err := http.Get(site)
	if err != nil || resp.StatusCode != http.StatusOK {
		m.Lock()
		siteStatusMap[site] = "DOWN"
		log.Printf("SITE: %v => DOWN", site)
		m.Unlock()
		w.Done()
		return
	}
	m.Lock()
	siteStatusMap[site] = "UP"
	log.Printf("SITE: %v => UP", site)
	m.Unlock()
	w.Done()
}

func (s statusCheckerStruct) CheckWebsiteStatus() {
	var m sync.Mutex
	var w sync.WaitGroup
	for {
		for curSite := range siteStatusMap {
			w.Add(1)
			go s.CheckWebsiteStatusUtil(curSite, &m, &w)
		}
		w.Wait()
		log.Println("------------------|SLEEPING FOR 1 MINUTE|------------------")
		time.Sleep(60 * time.Second)
	}
}

func handlePostWebsitesList(c *gin.Context) {
	log.Println("(----|POST|----)")
	var newSiteList websiteListStruct

	if err := c.BindJSON(&newSiteList); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "website list not updated",
		})
		return
	}

	for _, s := range newSiteList.SiteList {
		if _, ok := siteStatusMap[s]; !ok {
			siteStatusMap[s] = "WAIT"
		}
	}

	go statusCheckerObject.CheckWebsiteStatus()
	c.JSON(http.StatusCreated, gin.H{
		"message": "website list updated successfully!",
	})

}

func handleQueryParticularWebsiteStatus(c *gin.Context) {
	log.Println("(----|GET #1|----)")
	qSite, ok := c.GetQuery("site")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Query param not found, enter valid url string",
		})
		return
	}
	if _, found := siteStatusMap[qSite]; !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "This website not in the database, please add first",
		})
		return
	}
	resp := websiteStatusResp{Site: qSite, Status: siteStatusMap[qSite]}
	c.IndentedJSON(http.StatusOK, resp)
}

func handleQueryAllWebsitesStatus(c *gin.Context) {
	log.Println("(----|GET #2|----)")
	resp := []websiteStatusResp{}
	for curSite, curStatus := range siteStatusMap {
		resp = append(resp, websiteStatusResp{Site: curSite, Status: curStatus})
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func main() {
	r := gin.Default()
	r.GET("/websites", handleQueryAllWebsitesStatus)
	r.GET("/query", handleQueryParticularWebsiteStatus)
	r.POST("/websites", handlePostWebsitesList)
	r.Run("localhost:8080")
}
