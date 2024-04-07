package controller

import (
	"context"
	"near-location/internal/form"
	"near-location/internal/model"
	"near-location/internal/service"
	"near-location/pkg/util"
	"strconv"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	userService service.UserService
}

func NewController(userService service.UserService) *Controller {
	return &Controller{
		userService: userService,
	}
}

func getPagination(c *fiber.Ctx) (int64, int64) {
	pageSizeStr := c.Query("pageSize", "10")
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		log.Warnf("invalid page size input: %s, set default 10", err)
		pageSize = 10
	}
	pageIdxStr := c.Query("pageIdx", "0")
	pageIdx, err := strconv.ParseInt(pageIdxStr, 10, 64)
	if err != nil {
		log.Warnf("invalid page index input: %s, set default 0", err)
		pageIdx = 0
	}
	return pageSize, pageIdx
}

func mapUserLocationToResponse(userLocation model.UserLocation) form.UserLocation {
	return form.UserLocation{
		Id:        userLocation.Id.Hex(),
		UpdatedAt: userLocation.UpdatedAt,
		UserId:    userLocation.UserId,
		Longitude: userLocation.Location.Coordinates[0],
		Latitude:  userLocation.Location.Coordinates[1],
	}
}

func (controller *Controller) GetLocations(c *fiber.Ctx) error {
	ctx := context.Background()
	pageSize, pageIdx := getPagination(c)
	latStr := c.Query("lat")
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return util.ErrBadRequest(err, "invalid latitude value")
	}
	lonStr := c.Query("lon")
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return util.ErrBadRequest(err, "invalid longitude value")
	}
	maxDistanceStr := c.Query("maxDistance")
	var maxDistance int64
	if maxDistanceStr != "" {
		maxDistanceCV, err := strconv.ParseInt(maxDistanceStr, 10, 64)
		if err != nil {
			return util.ErrBadRequest(err, "invalid max distance value")
		}
		maxDistance = maxDistanceCV
	}

	result, err := controller.userService.FindUserLocationsNearDatapoint(ctx, model.Datapoint{
		Longitude: lon,
		Latitude:  lat,
	}, maxDistance, pageSize, pageIdx)
	if err != nil {
		return err
	}
	if int64(len(result)) < pageSize {
		pageSize = int64(len(result))
	}
	var listUserLocationReponse []form.UserLocation
	for _, value := range result {
		res := mapUserLocationToResponse(value)
		listUserLocationReponse = append(listUserLocationReponse, res)
	}
	response := form.GetLocationsResponse{
		Pagination: form.PaginationResponsePartial{
			PageIdx:  pageIdx,
			PageSize: pageSize,
		},
		Data: listUserLocationReponse,
	}
	return c.JSON(response)
}
