package main

import (
	"argotech/ml"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func calc_loss(deltaT, deltaH, t float64, dFilm, p string) float64 {
	type_of_apple := make(map[string]float64)
	type_of_apple["golden"] = 1.0
	type_of_apple["gala"] = 1.4
	type_of_apple["fuji"] = 1.2
	type_of_apple["aidared"] = 0.8
	type_of_apple["jonagold"] = 1.6

	type_of_box := make(map[string]float64)
	type_of_box["wood"] = 0.55
	type_of_box["plastic"] = 0.4
	type_of_box["cardboard"] = 0.7

	var A = 1.2
	var k1 = 0.1
	var k2 = 0.05
	var k3 = 0.2
	var C = 2.0
	expPart := math.Exp(k1*deltaT + k2*deltaH + k3*type_of_apple[dFilm])
	return A * expPart * type_of_box[p] * t + C
}

func parse(num string) float64 {
	floatVal, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return -999.0
	}
	return floatVal
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./template/*")
	r.Static("/static", "./static")
  r.GET("/", func(c *gin.Context) {
  	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
  })
  r.POST("/calculate", func(c *gin.Context) {
  	type Req struct {
   		Variety     string 	`json:"variety"`
      Temperature string  `json:"temperature"`
      Humidity    string  `json:"humidity"`
      Package     string 	`json:"package"`
      Duration 		string  `json:"duration"`
    }
    var req Req
    if err := c.BindJSON(&req); err != nil {
    	c.JSON(http.StatusBadRequest, gin.H{})
	    return
    }

    temp := parse(req.Temperature)
    if temp == -999.0 {
    	c.JSON(http.StatusBadRequest, gin.H{
     	  "Message": "Uncorrect temperature",
      })
    }
    humi := parse(req.Humidity)
    if humi == -999.0 {
    	c.JSON(http.StatusBadRequest, gin.H{
     	  "Message": "Uncorrect humidity",
      })
    }
    dura := parse(req.Duration)
    if dura == -999.0 {
    	c.JSON(http.StatusBadRequest, gin.H{
     	  "Message": "Uncorrect duration",
      })
    }
    recommendations := ml.ReqML(fmt.Sprintf(
      `Give recommendations for better storage of apples based on the following details:
      Apple type: %s
      Target shelf life: %f days
      Packaging: %s
      Humidity: %f procent
      Temperature: %f °C
      Provide a coherent, informative paragraph (around 150 words) explaining whether the current storage conditions are appropriate, what potential risks exist, and how they can be improved to meet or exceed the desired shelf life. Use clear, practical language for farmers or distributors.`,
      req.Variety,
      dura,
      req.Package,
      humi,
      temp,
    ))
    loss := calc_loss(temp, humi, dura / 30, req.Variety, req.Package)
    var res float64
    if loss > 100.0 {
    	res = 100.0
    } else {
    	res = loss
    }
    c.JSON(http.StatusOK, gin.H{
    	"loss": res,
      "recommendations": recommendations,
    })
  })
  r.Run(":80")
}