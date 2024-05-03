package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"dbtest/app"
	"dbtest/domain/dto"
	"dbtest/model"

	util "dbtest/util"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	heroUseCase model.HeroUseCase
	fileUseCase model.FileUseCase
}

func NewHandler(router *gin.RouterGroup, cfg *app.AppConfig) {
	handler := &Handler{
		heroUseCase: cfg.HeroUseCase,
		fileUseCase: cfg.FileUseCase,
	}

	router.GET("/heros", handler.GetAll)
	router.GET("/heros/:id", handler.GetHeroById)
	router.POST("/heros", handler.SaveHero)
	
	router.POST("/upload", handler.UploadFile)
	router.GET("/download", handler.Download)
	router.GET("/keys", handler.Keys)
}

func (h *Handler) UploadFile(c *gin.Context) {
	formFile, formFileError := c.FormFile("file")
	if formFileError != nil {
		log.Println("error abriendo multipart", formFileError)
		c.Status(500)
		return
	}

	if (formFile.Size == 0) {
		log.Println("File size incorrecta", formFileError)
		c.Status(500)
		return
	}

	multipartFile, multipartFileError := formFile.Open()
	if multipartFileError != nil {
		log.Println("Error abriendo multipartFile...", multipartFileError)
		c.Status(500)
		return
	}
	
	defer multipartFile.Close()
	
	fileContent, fileContentError := io.ReadAll(multipartFile)
	if fileContentError != nil {
		log.Println("Error leyendo todos los bytes de multipartFile", fileContentError)
		c.Status(500)
		return
	}

	response, er := h.fileUseCase.SaveFile(&fileContent, &formFile.Filename)

	if er != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(200, response)
}

func (h *Handler) Keys(c *gin.Context) {
	token := c.Query("next")

	keys, error := h.fileUseCase.GetFiles(token)

	if error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(200, keys)
}

func (h *Handler) Download(c *gin.Context) {
	key := c.Query("key")

	if key == "" {
		c.AbortWithError(401, ErrBadRequest)
		return
	}

	output, error := h.fileUseCase.GetFile(key)

	if error != nil {
		c.Status(http.StatusInternalServerError)
		return;
	}

	if output.Data == nil {
		c.Status(http.StatusNotFound)
		return;
	}
	
	defer output.Data.Close()

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`"attachment; filename="%s"`, key),
	}

	c.DataFromReader(200, output.ContentLength, output.ContentType, output.Data, extraHeaders)
}

func (h *Handler) GetAll(c *gin.Context) {
	response := h.heroUseCase.GetAllHeros()
	c.JSON(response.Status, response)
}

func (h *Handler) GetHeroById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.AbortWithError(401, ErrBadRequest)
		return
	}

	response := h.heroUseCase.GetHeroById(id)

	c.JSON(response.Status, response)
}

func (h *Handler) SaveHero(c *gin.Context) {
	var newHero dto.HeroDto

	if err := c.BindJSON(&newHero); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	valid, fields := util.Validate(newHero)

	if !valid {
		dto.InvalidData.Data = fields
		c.Error(dto.InvalidData)
		return
	}

	response := h.heroUseCase.SaveHero(newHero)
	c.JSON(response.Status, response)
}