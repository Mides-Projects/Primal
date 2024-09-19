package grants

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/grantsx/model"
	"github.com/holypvp/primal/service"
	"net/http"
)

func GrantsCreateRoute(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return common.HTTPError(c, http.StatusBadRequest, "No name found for the player")
	}

	var body map[string]interface{}
	if err := c.Bind().Body(&body); err != nil {
		return common.HTTPError(c, http.StatusBadRequest, "Failed to bind request body: "+err.Error())
	}

	g := &model.Grant{}
	if err := g.Unmarshal(body); err != nil {
		return common.HTTPError(c, http.StatusBadRequest, "Failed to unmarshal grant: "+err.Error())
	}

	gaAdder := service.Grants().Lookup(g.AddedBy())
	if gaAdder == nil {
		return common.HTTPError(c, http.StatusNotFound, "Grants for account adder not found (Not cached)")
	}

	if gaAdder.Account().Id() != g.AddedBy() {
		return common.HTTPError(c, http.StatusConflict, "Source adder ID mismatch")
	}

	hgAdder := service.Grants().HighestGroupBy(gaAdder)
	if hgAdder == nil {
		return common.HTTPError(c, http.StatusNotFound, "Highest group not found for who added the grant")
	}

	// Retrieve the account of the player from our redis cache
	// but if they are online, we can fetch it from the RAM Cache
	acc, err := service.Account().UnsafeLookupByName(name)
	if err != nil {
		return common.HTTPError(c, http.StatusInternalServerError, "Failed to lookup account: "+err.Error())
	} else if acc == nil {
		return common.HTTPError(c, http.StatusNotFound, "Player not found")
	}

	ga, err := service.Grants().UnsafeLookup(acc.Id())
	if err != nil {
		return common.HTTPError(c, http.StatusInternalServerError, "Failed to lookup grant: "+err.Error())
	}

	if ga == nil {
		return common.HTTPError(c, http.StatusNotFound, "Player not found")
	}

	if ga.Account().Id() != acc.Id() {
		return common.HTTPError(c, http.StatusConflict, "Player ID mismatch")
	}

	hg := service.Grants().HighestGroupBy(ga)
	if hg != nil && hg.Weight() > hgAdder.Weight() {
		return common.HTTPError(c, http.StatusUnauthorized, "You cannot grant a rank to someone with a higher rank than you")
	}

	go func() {
		if err := service.Grants().Save(acc.Id(), g); err != nil {
			common.Log.Fatalf("Failed to save grant: %v", err)
		}
	}()

	return common.HTTPError(c, http.StatusOK, "Grant saved")
}
