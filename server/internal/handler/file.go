package handler

import (
	"github.com/Garagator3000/cft-api/server"
	"github.com/Garagator3000/cft-api/server/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (handler *Handler) GetAll(ctx *gin.Context) {
	response := make([]service.File, 0)

	files, err := handler.service.GetList()
	if err != nil {
		server.Trace(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	for _, file := range files {
		response = append(response, service.File{
			Name:     file.Name,
			Checksum: file.Checksum,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (handler *Handler) Get(ctx *gin.Context) {
	name := ctx.Param("name")
	file, err := handler.service.GetFile(name)
	if err != nil {
		server.Trace(err)
		ctx.JSON(http.StatusNotFound, "File '"+name+"' not found")
	} else {
		ctx.Header("filename", name)
		ctx.File(file.Name())
	}
}

func (handler *Handler) Post(ctx *gin.Context) {
	file, fileHeader, err := ctx.Request.FormFile("filename")
	if err != nil {
		server.Trace(err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	defer file.Close()

	content := make([]byte, fileHeader.Size)
	_, err = file.Read(content)
	if err != nil {
		server.Trace(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	err = handler.service.SaveFile(fileHeader.Filename, content)
	if err != nil {
		server.Trace(err)
		ctx.JSON(http.StatusConflict, err)
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

func (handler *Handler) Put(ctx *gin.Context) {
	file, fileHeader, err := ctx.Request.FormFile("filename")
	if err != nil {
		server.Trace(err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	defer file.Close()

	content := make([]byte, fileHeader.Size)
	_, err = file.Read(content)
	if err != nil {
		server.Trace(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	err = handler.service.UpdateFile(fileHeader.Filename, content)
	if err != nil {
		server.Trace(err)
		ctx.JSON(http.StatusConflict, err)
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

func (handler *Handler) Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	err := handler.service.DeleteFile(name)
	if err != nil {
		if err.Error() == "file does not exist" {
			server.Trace(err)
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		server.Trace(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}
