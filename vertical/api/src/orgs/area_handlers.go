package orgs

import (
	"orbita-api/core/httpx"
	d "orbita-api/core/decorators"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getAreas(c *gin.Context) {
	companyID := c.GetInt("companyID")
	if companyID == 0 {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Missing company id"})
		return
	}
	service := AreaService{}
	areas, err := service.GetByCompany(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, areas)
}

var GetAreas = d.AuthenticationRequired(getAreas)

func createArea(c *gin.Context) {
	companyID := c.GetInt("companyID")
	if companyID == 0 {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Missing company id"})
		return
	}
	body := AreaBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	area := Area{CompanyID: companyID, Name: body.Name}
	service := AreaService{}
	if err := service.Insert(&area); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusCreated, area)
}

var CreateArea = d.CompanyRequiredAndPermissionRequired([]string{OrgsCompanyCreate}, createArea)

func updateArea(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service := AreaService{}
	area, err := service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Area not found"})
		return
	}
	body := AreaBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	area.Name = body.Name
	if err := service.Update(area); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, area)
}

var UpdateArea = d.CompanyRequiredAndPermissionRequired([]string{OrgsCompanyUpdate}, updateArea)

func deleteArea(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service := AreaService{}
	area, err := service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Area not found"})
		return
	}
	if err := service.Delete(area.ID); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

var DeleteArea = d.CompanyRequiredAndPermissionRequired([]string{OrgsCompanyDelete}, deleteArea)

func getUserAreas(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userId"))
	companyID := c.GetInt("companyID")
	service := AreaService{}
	areas, err := service.GetUserAreas(userID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, areas)
}

var GetUserAreas = d.CompanyRequiredAndPermissionRequired([]string{OrgsCompanyRead}, getUserAreas)

type SetUserAreasBody struct {
	AreaIDs []int `json:"areaIds" binding:"required"`
}

func setUserAreas(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userId"))
	body := SetUserAreasBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	service := AreaService{}
	if err := service.SetUserAreas(userID, body.AreaIDs); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

var SetUserAreas = d.CompanyRequiredAndPermissionRequired([]string{OrgsCompanyUpdate}, setUserAreas)

// Superuser endpoints for managing areas by company ID

func getCompanyAreas(c *gin.Context) {
	companyID, _ := strconv.Atoi(c.Param("id"))
	service := AreaService{}
	areas, err := service.GetByCompany(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, areas)
}

var GetCompanyAreas = d.SuperuserRequired(getCompanyAreas)

func createCompanyArea(c *gin.Context) {
	companyID, _ := strconv.Atoi(c.Param("id"))
	body := AreaBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	area := Area{CompanyID: companyID, Name: body.Name}
	service := AreaService{}
	if err := service.Insert(&area); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusCreated, area)
}

var CreateCompanyArea = d.SuperuserRequired(createCompanyArea)
