package credential

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mrexmelle/connect-authx/internal/config"
	"github.com/mrexmelle/connect-authx/internal/dto/dtorespwithoutdata"
	"github.com/mrexmelle/connect-authx/internal/localerror"
)

type Controller struct {
	ConfigService     *config.Service
	CredentialService *Service
	LocalErrorService *localerror.Service
}

func NewController(cfg *config.Service, svc *Service, les *localerror.Service) *Controller {
	return &Controller{
		ConfigService:     cfg,
		CredentialService: svc,
		LocalErrorService: les,
	}
}

// Post Credentials : HTTP endpoint to post new credentials
// @Tags Credentials
// @Description Post a new credential
// @Accept json
// @Produce json
// @Param data body PostRequestDto true "Credential Request"
// @Success 200 {object} PostResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /credentials [POST]
func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	var requestBody PostRequestDto
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		dtorespwithoutdata.New(
			localerror.ErrBadJson.Error(),
			err.Error(),
		).RenderTo(w, http.StatusBadRequest)
		return
	}

	err = c.CredentialService.Create(requestBody)
	info := c.LocalErrorService.Map(err)
	dtorespwithoutdata.New(
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}

// Delete Credentials : HTTP endpoint to delete credentials
// @Tags Credentials
// @Description Delete a credential
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Success 200 {object} DeleteResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /credentials/{employee_id} [DELETE]
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	err := c.CredentialService.DeleteByEmployeeId(chi.URLParam(r, "employee_id"))
	info := c.LocalErrorService.Map(err)
	dtorespwithoutdata.New(
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}

// Patch Password : HTTP endpoint to patch password
// @Tags Credentials
// @Description Patch a credential's password
// @Accept json
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Param data body PatchPasswordRequestDto true "Password Patch Request"
// @Success 200 {object} PatchPasswordResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /credentials/{employee_id}/password [PATCH]
func (c *Controller) PatchPassword(w http.ResponseWriter, r *http.Request) {
	var requestBody PatchPasswordRequestDto
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		dtorespwithoutdata.New(
			localerror.ErrBadJson.Error(),
			err.Error(),
		).RenderTo(w, http.StatusBadRequest)
		return
	}
	err = c.CredentialService.UpdatePasswordByEmployeeId(
		chi.URLParam(r, "employee_id"),
		requestBody)
	info := c.LocalErrorService.Map(err)
	dtorespwithoutdata.New(
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}

// Reset Password : HTTP endpoint to reset password
// @Tags Credentials
// @Description Reset a credential's password
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Success 200 {object} PatchPasswordResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /credentials/{employee_id}/password [DELETE]
func (c *Controller) DeletePassword(w http.ResponseWriter, r *http.Request) {
	err := c.CredentialService.ResetPasswordByEmployeeId(chi.URLParam(r, "employee_id"))
	info := c.LocalErrorService.Map(err)
	dtorespwithoutdata.New(
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}
