package orgs

import (
	"benefit-calculator-api/core/httpx"
	"benefit-calculator-api/core/shared"
	"net/http"
	"strconv"

	d "benefit-calculator-api/core/decorators"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func getCompanies(c *gin.Context) {
	service := CompanyService{}
	iuser, _ := c.Get("user")
	user, _ := iuser.(shared.User)

	if user != nil && user.IsSuperUser() {
		companies, err := service.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
			return
		}
		c.JSON(http.StatusOK, companies)
		return
	}

	companies, err := service.GetByUserID(user.GetID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, companies)
}

var GetCompanies = d.PermissionRequired([]string{OrgsCompanyRead}, getCompanies)

func getCompany(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service := CompanyService{}
	company, err := service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, company)
}

var GetCompany = d.PermissionRequired([]string{OrgsCompanyRead}, getCompany)

func createCompany(c *gin.Context) {
	payload := CompanyBody{}
	err := payload.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	company := Company{}
	err = company.ApplyPayload(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusCreated, company)
}

var CreateCompany = d.PermissionRequired([]string{OrgsCompanyCreate}, createCompany)

func updateCompany(c *gin.Context) {
	companyID, _ := strconv.Atoi(c.Param("id"))
	service := CompanyService{}
	company, err := service.GetByID(companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Company not found"})
		return
	}
	payload := CompanyBody{}
	err = payload.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	err = company.ApplyPayload(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, company)
}

var UpdateCompany = d.PermissionRequired([]string{OrgsCompanyUpdate}, updateCompany)

func deleteCompany(c *gin.Context) {
	companyID, _ := strconv.Atoi(c.Param("id"))
	service := CompanyService{}
	company, err := service.GetByID(companyID)
	if err != nil {
		zap.S().Errorw("Error while getting company", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Company not found"})
		return
	}
	if err := company.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

var DeleteCompany = d.PermissionRequired([]string{OrgsCompanyDelete}, deleteCompany)
