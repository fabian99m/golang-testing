package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"dbtest/config"
	util "dbtest/util"
	"dbtest/db"
	"dbtest/domain/dto"
	"dbtest/model"
	"dbtest/repository"
	"dbtest/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	heroUseCase model.HeroUseCase
	fileUseCase model.FileUseCase
}

func NewHandler(router *gin.RouterGroup, cfg *config.Config) {
	dbConnection := db.NewDbConnection(&cfg.DataBase)
	heroRepository := repository.NewHeroRespository(dbConnection)

	handler := &Handler{
		heroUseCase: usecase.NewHeroUseCase(heroRepository),
		fileUseCase: usecase.NewFileUseCase(cfg.S3.Bucket, nil),
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
	response := h.fileUseCase.SaveFile(formFile)

	c.JSON(200, response)
}

func (h *Handler) Keys(c *gin.Context) {
	token := c.Query("next")

	keys, error := h.fileUseCase.GetFiles(token);

	if (error != nil) {
		c.Status(http.StatusInternalServerError)
		return;
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

	if (error != nil) {
		c.Status(http.StatusInternalServerError)
		return;
	}

	if output.Data == nil {
		c.Status(http.StatusNotFound)
		return;
	}

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