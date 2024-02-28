package ginV1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skantay/hezzl/internal/domain/good/model"
)

func (g ginController) goodsListHandler(c *gin.Context) {
	limit, err := parseQueryParamAtoi(c, "limit", 10)
	if err != nil {
		g.log.Error(err.Error())
		badRequest(c, http.StatusText(http.StatusBadRequest))

		return
	}

	offset, err := parseQueryParamAtoi(c, "offset", 1)
	if err != nil {
		g.log.Error(err.Error())
		badRequest(c, http.StatusText(http.StatusBadRequest))

		return
	}

	if offset < 1 || limit < 1 {
		g.log.Sugar().Errorf("invalid numbers %d and %d", offset, limit)
		badRequest(c, http.StatusText(http.StatusBadRequest))

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
			model.Good
		} `json:"goods"`
	}

	var response ResponseList

	goods, err := g.service.GoodService.GetGoods(c, limit, offset)
	if err != nil {
		g.log.Error(err.Error())

		if !errors.Is(err, model.ErrGoodNotFound) {
			internalServerError(c)

			return
		}
	}

	var removed int

	for _, good := range goods {
		response.Goods = append(response.Goods, struct{ model.Good }{good})

		if good.Removed {
			removed++
		}
	}

	if len(goods) == 0 {
		response.Goods = []struct{ model.Good }{}
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
		badRequest(c, err.Error())

		return
	}

	projectID, err := parseQueryParamAtoi(c, "projectID", -1)
	if err != nil {
		g.log.Error(err.Error())
		badRequest(c, err.Error())

		return
	}

	if id < 0 || projectID < 0 {
		g.log.Sugar().Errorf("invalid numbers %d and %d", id, projectID)
		badRequest(c, http.StatusText(http.StatusBadRequest))

		return
	}

	type Payload struct {
		NewPriority int `json:"newPriority"`
	}

	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		g.log.Error(err.Error())
		badRequest(c, http.StatusText(http.StatusBadRequest))

		return
	}

	goods, err := g.service.GoodService.Reprioritiize(c.Request.Context(), payload.NewPriority, id, projectID)
	if err != nil {
		if errors.Is(err, model.ErrGoodNotFound) {

			g.log.Error(err.Error())
			handleError(c, err)

			return
		}

		g.log.Error(err.Error())
		internalServerError(c)

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
		badRequest(c, err.Error())

		return
	}

	projectID, err := parseQueryParamAtoi(c, "projectID", -1)
	if err != nil {
		g.log.Error(err.Error())
		badRequest(c, err.Error())

		return
	}

	if id < 0 || projectID < 0 {
		g.log.Sugar().Errorf("invalid numbers %d and %d", id, projectID)
		badRequest(c, http.StatusText(http.StatusBadRequest))

		return
	}

	good, err := g.service.GoodService.DeleteGood(c.Request.Context(), id, projectID)
	if err != nil {
		if errors.Is(err, model.ErrGoodNotFound) {

			g.log.Error(err.Error())
			handleError(c, err)

			return
		}
		g.log.Error(err.Error())
		internalServerError(c)

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
		badRequest(c, err.Error())

		return
	}

	projectID, err := parseQueryParamAtoi(c, "projectID", -1)
	if err != nil {
		g.log.Error(err.Error())
		badRequest(c, err.Error())

		return
	}

	if id < 0 || projectID < 0 {
		g.log.Sugar().Errorf("invalid numbers %d and %d", id, projectID)
		badRequest(c, http.StatusText(http.StatusBadRequest))

		return
	}

	type Payload struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		g.log.Error(err.Error())
		badRequest(c, http.StatusText(http.StatusBadRequest))

		return
	}

	if payload.Name == "" {
		g.log.Error("empty payload name")
		badRequest(c, http.StatusText(http.StatusBadRequest))

		return
	}

	emptyDesc := false

	if payload.Description == nil {
		emptyDesc = true
		word := ""
		payload.Description = &word
	}

	good, err := g.service.GoodService.UpdateGood(c.Request.Context(), id, projectID, payload.Name, *payload.Description, emptyDesc)
	if err != nil {
		if errors.Is(err, model.ErrGoodNotFound) {

			g.log.Error(err.Error())
			handleError(c, err)

			return
		}

		g.log.Error(err.Error())
		internalServerError(c)

		return
	}

	c.JSON(http.StatusOK, good)
}

func (g ginController) createGoodHandler(c *gin.Context) {
	projectID, err := parseQueryParamAtoi(c, "projectID", -1)
	if err != nil {
		g.log.Error(err.Error())
		badRequest(c, err.Error())

		return
	}

	if projectID < 0 {
		badRequest(c, http.StatusText(http.StatusBadRequest))

		return
	}

	type Payload struct {
		Name string `json:"name"`
	}

	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		g.log.Error(err.Error())
		badRequest(c, http.StatusText(http.StatusBadRequest))

		return
	}

	good, err := g.service.GoodService.CreateGood(c.Request.Context(), projectID, payload.Name)
	if err != nil && !errors.Is(err, model.ErrGoodNotFound) {
		g.log.Sugar().Errorf("%v", err)
		internalServerError(c)

		return
	} else if errors.Is(err, model.ErrGoodNotFound) {
		g.log.Sugar().Errorf("%v", err)
		badRequest(c, "project id not found")

		return
	}

	c.JSON(http.StatusOK, good)
}
