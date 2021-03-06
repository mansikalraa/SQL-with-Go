// CRUD with Gin Gonic, GORM and MySQL

package main

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Student struct {
	Rno   string `json:"rno" binding:"required"`
	Sname string `json:"sname" binding:"required"`
}

var DB *gorm.DB
var students []Student

func FindStudents(c *gin.Context) {
	var students []Student
	DB.Find(&students)

	c.JSON(http.StatusOK, gin.H{"data": students})
}

func FindStudent(c *gin.Context) {
	var student Student
	if err := DB.Where("rno = ?", c.Param("rno")).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": student})
}

func CreateStudent(c *gin.Context) {
	var input Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student := Student{Rno: input.Rno, Sname: input.Sname}
	DB.Create(&student)

	c.JSON(http.StatusOK, gin.H{"data": "Record created."})
}

func UpdateStudent(c *gin.Context) {
	var student Student
	if err := DB.Where("rno = ?", c.Param("rno")).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var input Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DB.Model(&Student{}).Where("rno = ?", c.Param("rno")).Update(input)
	c.JSON(http.StatusOK, gin.H{"data": "Record updated."})
}

func DeleteStudent(c *gin.Context) {
	var student Student
	if err := DB.Where("rno = ?", c.Param("rno")).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	DB.Where("rno = ?", c.Param("rno")).Delete(&student)

	c.JSON(http.StatusOK, gin.H{"data": "Record deleted."})
}

func main() {
	// Database connection
	db, err := gorm.Open("mysql", "root:belikemee@tcp(127.0.0.1:3306)/try")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	db.AutoMigrate(&Student{})

	DB = db

	r := gin.Default()

	// Routes
	r.GET("/students", FindStudents)
	r.GET("/students/:rno", FindStudent)
	r.POST("/students", CreateStudent)
	r.PUT("/students/:rno", UpdateStudent)
	r.DELETE("/students/:rno", DeleteStudent)

	r.Run(":8000")
}
