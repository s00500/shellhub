package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/labstack/echo"
	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/api/user"
)

const (
	UpdateUserURL = "/user"
)

func UpdateUser(c apicontext.Context) error {
	var req struct {
		Username        string `json:"username"`
		Email           string `json:"email"`
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json: "newPassword"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}

	ID := ""
	if v := c.ID(); v != nil {
		ID = v.ID
	}
	if req.CurrentPassword != "" {
		sum := sha256.Sum256([]byte(req.CurrentPassword))
		sumByte := sum[:]
		req.CurrentPassword = hex.EncodeToString(sumByte)
	}
	if req.NewPassword != "" {
		sum := sha256.Sum256([]byte(req.NewPassword))
		sumByte := sum[:]
		req.NewPassword = hex.EncodeToString(sumByte)
	}

	svc := user.NewService(c.Store())

	if invalidFields, err := svc.UpdateDataUser(c.Ctx(), req.Username, req.Email, req.CurrentPassword, req.NewPassword, ID); err != nil {
		switch {
		case err == user.ErrUnauthorized:
			return echo.ErrUnauthorized
		case err == user.ErrConflict:
			return c.JSON(http.StatusConflict, invalidFields)
		default:
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
