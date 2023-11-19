package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/bookstore")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

type KategoriBuku struct {
	ID           int    `json:"id"`
	NamaKategori string `json:"nama_kategori"`
}

type Buku struct {
	ID          int    `json:"id"`
	Judul       string `json:"judul"`
	Penulis     string `json:"penulis"`
	KategoriID  int    `json:"kategori_id"`
}

func main() {
	r := gin.Default()

	r.GET("/kategori", getKategoriBuku)
	r.POST("/kategori", createKategoriBuku)
	r.GET("/buku", getBuku)
	r.POST("/buku", createBuku)
	r.PUT("/buku/:id", updateBuku)
	r.DELETE("/buku/:id", deleteBuku)
	r.GET("/buku/search", searchBuku)


	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func getKategoriBuku(c *gin.Context) {
	rows, err := db.Query("SELECT id, nama_kategori FROM kategori_buku")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	kategoris := []KategoriBuku{}
	for rows.Next() {
		var kategori KategoriBuku
		err := rows.Scan(&kategori.ID, &kategori.NamaKategori)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		kategoris = append(kategoris, kategori)
	}

	c.JSON(http.StatusOK, kategoris)
}

func createKategoriBuku(c *gin.Context) {
	var kategori KategoriBuku
	if err := c.ShouldBindJSON(&kategori); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO kategori_buku (nama_kategori) VALUES (?)", kategori.NamaKategori)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	kategori.ID = int(id)

	c.JSON(http.StatusCreated, kategori)
}

func getBuku(c *gin.Context) {
	rows, err := db.Query("SELECT id, judul, penulis, kategori_id FROM buku")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	bukus := []Buku{}
	for rows.Next() {
		var buku Buku
		err := rows.Scan(&buku.ID, &buku.Judul, &buku.Penulis, &buku.KategoriID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bukus = append(bukus, buku)
	}

	c.JSON(http.StatusOK, bukus)
}

func createBuku(c *gin.Context) {
	var buku Buku
	if err := c.ShouldBindJSON(&buku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO buku (judul, penulis, kategori_id) VALUES (?, ?, ?)", buku.Judul, buku.Penulis, buku.KategoriID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	buku.ID = int(id)

	c.JSON(http.StatusCreated, buku)
}

func updateBuku(c *gin.Context) {
	id := c.Param("id")

	var buku Buku
	if err := c.ShouldBindJSON(&buku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE buku SET judul=?, penulis=?, kategori_id=? WHERE id=?", buku.Judul, buku.Penulis, buku.KategoriID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteBuku(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM buku WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func searchBuku(c *gin.Context) {
	// Ambil nilai pencarian dari query string
	query := c.Query("q")
 
	// Query ke database untuk mencari buku berdasarkan judul atau penulis
	rows, err := db.Query("SELECT id, judul, penulis, kategori_id FROM buku WHERE judul LIKE ? OR penulis LIKE ?", "%"+query+"%", "%"+query+"%")
	if err != nil {
	   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	   return
	}
	defer rows.Close()
 
	bukus := []Buku{}
	for rows.Next() {
	   var buku Buku
	   err := rows.Scan(&buku.ID, &buku.Judul, &buku.Penulis, &buku.KategoriID)
	   if err != nil {
		  c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		  return
	   }
	   bukus = append(bukus, buku)
	}
 
	c.JSON(http.StatusOK, bukus)
 }
 