package user

import (
	"net/http"

	"github.com/M15t/ghoul/internal/model"
	"github.com/M15t/ghoul/pkg/rbac"
	"github.com/M15t/ghoul/pkg/server"
	dbutil "github.com/M15t/ghoul/pkg/util/db"
	structutil "github.com/M15t/ghoul/pkg/util/struct"
)

// Custom errors
var (
	ErrIncorrectPassword = server.NewHTTPError(http.StatusBadRequest, "INCORRECT_PASSWORD", "Incorrect old password")
	ErrUserNotFound      = server.NewHTTPError(http.StatusBadRequest, "USER_NOTFOUND", "User not found")
	ErrUsernameExisted   = server.NewHTTPValidationError("Username already existed")
)

// Create creates a new user account
func (s *User) Create(authUsr *model.AuthUser, data CreationData) (*model.User, error) {
	if err := s.enforce(authUsr, model.ActionCreateAll); err != nil {
		return nil, err
	}

	if existed, err := s.udb.Exist(s.db, map[string]interface{}{"username": data.Username}); err != nil || existed {
		return nil, ErrUsernameExisted.SetInternal(err)
	}

	rec := &model.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Mobile:    data.Mobile,
		Username:  data.Username,
		Password:  s.cr.HashPassword(data.Password),
		Blocked:   data.Blocked,
		Role:      data.Role,
	}

	if err := s.udb.Create(s.db, rec); err != nil {
		return nil, server.NewHTTPInternalError("Error creating user").SetInternal(err)
	}

	return rec, nil
}

// View returns single user
func (s *User) View(authUsr *model.AuthUser, id int) (*model.User, error) {
	if err := s.enforce(authUsr, model.ActionViewAll); err != nil {
		return nil, err
	}

	rec := new(model.User)
	if err := s.udb.View(s.db, rec, id); err != nil {
		return nil, ErrUserNotFound.SetInternal(err)
	}

	return rec, nil
}

// List returns list of users
func (s *User) List(authUsr *model.AuthUser, lq *dbutil.ListQueryCondition, count *int) ([]*model.User, error) {
	if err := s.enforce(authUsr, model.ActionViewAll); err != nil {
		return nil, err
	}

	var data []*model.User
	if err := s.udb.List(s.db, &data, lq, count); err != nil {
		return nil, server.NewHTTPInternalError("Error listing user").SetInternal(err)
	}

	return data, nil
}

// Update updates user information
func (s *User) Update(authUsr *model.AuthUser, id int, data UpdateData) (*model.User, error) {
	if err := s.enforce(authUsr, model.ActionUpdateAll); err != nil {
		return nil, err
	}

	// optimistic update
	updates := structutil.ToMap(data)
	if err := s.udb.Update(s.db, updates, id); err != nil {
		return nil, server.NewHTTPInternalError("Error updating user").SetInternal(err)
	}

	rec := new(model.User)
	if err := s.udb.View(s.db, rec, id); err != nil {
		return nil, ErrUserNotFound.SetInternal(err)
	}

	return rec, nil
}

// Delete deletes a user
func (s *User) Delete(authUsr *model.AuthUser, id int) error {
	if err := s.enforce(authUsr, model.ActionDeleteAll); err != nil {
		return err
	}

	if existed, err := s.udb.Exist(s.db, id); err != nil || !existed {
		return ErrUserNotFound.SetInternal(err)
	}

	if err := s.udb.Delete(s.db, id); err != nil {
		return server.NewHTTPInternalError("Error deleting user").SetInternal(err)
	}

	return nil
}

// Me returns authenticated user
func (s *User) Me(authUsr *model.AuthUser) (*model.User, error) {
	rec := new(model.User)
	if err := s.udb.View(s.db, rec, authUsr.ID); err != nil {
		return nil, ErrUserNotFound.SetInternal(err)
	}
	return rec, nil
}

// ChangePassword changes authenticated user password
func (s *User) ChangePassword(authUsr *model.AuthUser, data PasswordChangeData) error {
	rec, err := s.Me(authUsr)
	if err != nil {
		return err
	}

	if !s.cr.CompareHashAndPassword(rec.Password, data.OldPassword) {
		return ErrIncorrectPassword
	}

	hashedPwd := s.cr.HashPassword(data.NewPassword)
	if err = s.udb.Update(s.db, map[string]interface{}{"password": hashedPwd}, rec.ID); err != nil {
		return server.NewHTTPInternalError("Error changing password").SetInternal(err)
	}

	return nil
}

// enforce checks user permission to perform the action
func (s *User) enforce(authUsr *model.AuthUser, action string) error {
	if !s.rbac.Enforce(authUsr.Role, model.ObjectUser, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
