package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skantay/hezzl/internal/entity"
)

// Internal Server Error
var ErrISE = errors.New(http.StatusText(http.StatusInternalServerError))

// Bad Request
var ErrBR = errors.New(http.StatusText(http.StatusBadRequest))

func (g ginController) createGoodHandler(c *gin.Context) {
	projectID, err := parseQueryParamAtoi(c, "projectID", -1)
	if err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	if projectID < 0 {
		handleError(c, "project id is invalid", http.StatusBadRequest, ErrBR)

		return
	}

	type Payload struct {
		Name string `json:"name"`
	}

	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	good, err := g.service.Good.Create(c.Request.Context(), projectID, payload.Name)
	if err != nil && !errors.Is(err, entity.ErrProjectNotFound) {
		g.log.Sugar().Errorf("%v", err)
		handleError(c, "", http.StatusInternalServerError, ErrISE)

		return
	} else if errors.Is(err, entity.ErrProjectNotFound) {
		g.log.Sugar().Errorf("%v", err)
		handleError(c, "", http.StatusNotFound, entity.ErrProjectNotFound)

		return
	}

	c.JSON(http.StatusOK, good)
}

func (g ginController) goodsListHandler(c *gin.Context) {
	limit, err := parseQueryParamAtoi(c, "limit", 10)
	if err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	offset, err := parseQueryParamAtoi(c, "offset", 1)
	if err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	if offset < 1 || limit < 1 {
		g.log.Sugar().Errorf("invalid numbers %d and %d", offset, limit)
		handleError(c, "offset and limit are invalid", http.StatusBadRequest, ErrBR)

		return
	}

	type ResponseList struct {
		Meta struct {
			Total   int `json:"total"`
			Removed int `json:"removed"`
			Limit   int `json:"limit"`
			Offset  int `json:"offset"`
		} `json:"meta"`
		Goods []struct {
			entity.Good
		} `json:"goods"`
	}

	var response ResponseList

	goods, err := g.service.Good.List(c, limit, offset)
	if err != nil {
		g.log.Error(err.Error())

		if !errors.Is(err, entity.ErrGoodNotFound) {
			handleError(c, "", http.StatusInternalServerError, ErrISE)

			return
		}
	}

	var removed int

	for _, good := range goods {
		response.Goods = append(response.Goods, struct{ entity.Good }{good})

		if good.Removed {
			removed++
		}
	}

	if len(goods) == 0 {
		response.Goods = []struct{ entity.Good }{}
	}

	response.Meta.Limit = limit

	response.Meta.Offset = offset

	response.Meta.Removed = removed

	response.Meta.Total = len(goods)

	c.JSON(http.StatusOK, response)
}

func (g ginController) reprioritizeGoodHandler(c *gin.Context) {
	id, err := parseQueryParamAtoi(c, "id", -1)
	if err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	projectID, err := parseQueryParamAtoi(c, "projectID", -1)
	if err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	if id < 0 || projectID < 0 {
		g.log.Sugar().Errorf("invalid numbers %d and %d", id, projectID)
		handleError(c, "ID or project ID invalid", http.StatusBadRequest, ErrBR)

		return
	}

	type Payload struct {
		NewPriority int `json:"newPriority"`
	}

	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	goods, err := g.service.Good.Reprioritiize(c.Request.Context(), payload.NewPriority, id, projectID)
	if err != nil {
		if errors.Is(err, entity.ErrGoodNotFound) {

			g.log.Error(err.Error())
			handleError(c, "", http.StatusNotFound, entity.ErrGoodNotFound)

			return
		}

		g.log.Error(err.Error())
		handleError(c, "", http.StatusInternalServerError, ErrISE)

		return
	}

	type response struct {
		ID       int `json:"id"`
		Priority int `json:"priority"`
	}

	responseJSON := make([]response, len(goods))

	for i, good := range goods {
		responseJSON[i].ID = good.ID
		responseJSON[i].Priority = good.Priority
	}

	c.JSON(http.StatusOK, responseJSON)
}

func (g ginController) removeGoodHandler(c *gin.Context) {
	id, err := parseQueryParamAtoi(c, "id", -1)
	if err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	projectID, err := parseQueryParamAtoi(c, "projectID", -1)
	if err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	if id < 0 || projectID < 0 {
		g.log.Sugar().Errorf("invalid numbers %d and %d", id, projectID)
		handleError(c, "ID or project ID invalid", http.StatusBadRequest, ErrBR)

		return
	}

	good, err := g.service.Good.Delete(c.Request.Context(), id, projectID)
	if err != nil {
		if errors.Is(err, entity.ErrGoodNotFound) {

			g.log.Error(err.Error())
			handleError(c, "", http.StatusNotFound, entity.ErrGoodNotFound)

			return
		}
		g.log.Error(err.Error())
		handleError(c, "", http.StatusInternalServerError, ErrISE)

		return
	}

	response := struct {
		Id         int  `json:"id"`
		CampaignID int  `json:"campignID"`
		Removed    bool `json:"removed"`
	}{
		Id:         good.ID,
		CampaignID: good.ProjectID,
		Removed:    good.Removed,
	}

	c.JSON(http.StatusAccepted, response)
}

func (g ginController) updateGoodHandler(c *gin.Context) {
	id, err := parseQueryParamAtoi(c, "id", -1)
	if err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	projectID, err := parseQueryParamAtoi(c, "projectID", -1)
	if err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	if id < 0 || projectID < 0 {
		g.log.Sugar().Errorf("invalid numbers %d and %d", id, projectID)
		handleError(c, "ID or project ID invalid", http.StatusBadRequest, ErrBR)

		return
	}

	type Payload struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		g.log.Error(err.Error())
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	if payload.Name == "" {
		g.log.Error("empty payload name")
		handleError(c, "", http.StatusBadRequest, ErrBR)

		return
	}

	emptyDesc := false

	if payload.Description == nil {
		emptyDesc = true
		word := ""
		payload.Description = &word
	}

	good, err := g.service.Good.Update(c.Request.Context(), id, projectID, payload.Name, *payload.Description, emptyDesc)
	if err != nil {
		if errors.Is(err, entity.ErrGoodNotFound) {

			g.log.Error(err.Error())
			handleError(c, "", http.StatusNotFound, entity.ErrGoodNotFound)

			return
		}

		g.log.Error(err.Error())
		handleError(c, "", http.StatusInternalServerError, ErrISE)

		return
	}

	c.JSON(http.StatusOK, good)
}
