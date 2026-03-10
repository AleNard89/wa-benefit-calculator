package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

var constraintMessages = map[string]string{
	"auth_users_email_key":   "a user with this email already exists",
	"auth_roles_name_key":    "a role with this name already exists",
	"orgs_companies_name_key": "a company with this name already exists",
}

var pgCodeMessages = map[string]string{
	"23502": "a required field is missing",
	"23503": "this record is referenced by other data and cannot be modified",
	"23505": "",
	"23514": "the provided value does not meet the required criteria",
	"22001": "the provided value is too long",
	"22003": "the provided number is out of range",
	"22P02": "invalid input format",
	"42501": "insufficient permissions for this operation",
	"40001": "concurrent modification detected, please retry",
	"40P01": "concurrent modification detected, please retry",
}

func GetErrorMessage(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		zap.S().Warnw("Database error",
			"code", pgErr.Code,
			"message", pgErr.Message,
			"detail", pgErr.Detail,
			"constraint", pgErr.ConstraintName,
			"table", pgErr.TableName,
		)

		if pgErr.Code == "23505" {
			if msg, ok := constraintMessages[pgErr.ConstraintName]; ok {
				return msg
			}
			return "a record with this data already exists"
		}

		if msg, ok := pgCodeMessages[pgErr.Code]; ok && msg != "" {
			return msg
		}

		return "an internal error occurred"
	}
	return err.Error()
}

func JSONWithEtag(c *gin.Context, obj any) {
	err := AddEtagHeader(c, obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, obj)
}
