package handler

import (
	"errors"
	"net/http"
	"rest-api/internal/resource/model"
	"rest-api/internal/resource/repository"
	"rest-api/utils"
	e "rest-api/utils/errors"
	"rest-api/utils/validator"
)

type UserHandler struct {
	userRepository *repository.UserRepository
	validator      *validator.Validator
	errorResponse  *e.ErrorResponse
}

func NewUserHandler(userRepository *repository.UserRepository, v *validator.Validator, errResp *e.ErrorResponse) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
		validator:      v,
		errorResponse:  errResp,
	}
}

func (a *UserHandler) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Cod   int    `json:"cod"`
		Email string `json:"email"`
	}

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		a.errorResponse.BadRequestResponse(w, r, err)
		return
	}

	if model.ValidateEmail(a.validator, input.Email); !a.validator.Valid() {
		a.errorResponse.FailedValidationResponse(w, r, a.validator.Errors)
		return
	}

	user, err := a.userRepository.GetByCodAndEmail(input.Cod, input.Email)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			a.validator.AddError("code", "invalid validation code or email")
			a.errorResponse.FailedValidationResponse(w, r, a.validator.Errors)
		default:
			a.errorResponse.ServerErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true
	user.Cod = 0

	err = a.userRepository.Update(user)
	if err != nil {
		a.errorResponse.ServerErrorResponse(w, r, err)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user}, nil)
	if err != nil {
		a.errorResponse.ServerErrorResponse(w, r, err)
	}
}

func (a *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var userDTO model.UserSaveDTO
	err := utils.ReadJSON(w, r, &userDTO)
	if err != nil {
		a.errorResponse.BadRequestResponse(w, r, err)
		return
	}

	user, err := userDTO.ToModel()
	if err != nil {
		a.errorResponse.ServerErrorResponse(w, r, err)
		return
	}

	codActivation := utils.GenerateRandomCod()
	user.Cod = codActivation

	v := validator.New()

	if model.ValidateUser(v, user); !v.Valid() {
		a.errorResponse.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.userRepository.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			a.errorResponse.FailedValidationResponse(w, r, v.Errors)
		default:
			a.errorResponse.ServerErrorResponse(w, r, err)
		}

		return
	}

	/*
		MANDA COD VERIFICAÇÃO EMAIL
			app.background(func() {
				data := map[string]interface{}{
					"activationToken": codActivation,
					"userID":          user.ID,
				}

				err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
				if err != nil {
					app.logger.PrintError(err, nil)
				}
			})
	*/

	err = utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user}, nil)
	if err != nil {
		a.errorResponse.ServerErrorResponse(w, r, err)
	}
}
