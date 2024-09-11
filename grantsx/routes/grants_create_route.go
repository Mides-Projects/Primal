package routes

import (
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/grantsx/model"
	"github.com/holypvp/primal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GrantsCreateRoute(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return common.HTTPError(http.StatusBadRequest, "No name found for the player")
	}

	id := service.Account().FetchAccountId(name)
	if id == "" {
		return common.HTTPError(http.StatusNotFound, "Player not found")
	}

	var body map[string]interface{}
	if err := c.Bind(&body); err != nil {
		return common.HTTPError(http.StatusBadRequest, "Failed to bind request body: "+err.Error())
	}

	g := &model.Grant{}
	if err := g.Unmarshal(body); err != nil {
		return common.HTTPError(http.StatusBadRequest, "Failed to unmarshal grant: "+err.Error())
	}

	ga, err := service.Grants().Lookup(id)
	if err != nil {
		return common.HTTPError(http.StatusInternalServerError, "Failed to lookup grant: "+err.Error())
	}

	if ga == nil {
		return common.HTTPError(http.StatusNotFound, "Player not found")
	}

	if ga.Account().Id() != id {
		return common.HTTPError(http.StatusConflict, "Player ID mismatch")
	}

	gaAdder := service.Grants().LookupAtCache(g.AddedBy())
	if gaAdder == nil {
		return common.HTTPError(http.StatusNotFound, "Grants for source adder not found (Not cached)")
	}

	if gaAdder.Account().Id() != g.AddedBy() {
		return common.HTTPError(http.StatusConflict, "Source adder ID mismatch")
	}

	hgAdder := service.Grants().HighestGroupBy(gaAdder)
	if hgAdder == nil {
		return common.HTTPError(http.StatusNotFound, "Highest group not found for who added the grant")
	}

	hg := service.Grants().HighestGroupBy(ga)
	if hg != nil && hg.Weight() > hgAdder.Weight() {
		return common.HTTPError(http.StatusUnauthorized, "You cannot grant a rank to someone with a higher rank than you")
	}

	go func() {
		if err := service.Grants().Save(id, g); err != nil {
			common.Log.Errorf("Failed to save grant: %v", err)
		}
	}()

	return common.HTTPError(http.StatusOK, "Grant saved")
}
