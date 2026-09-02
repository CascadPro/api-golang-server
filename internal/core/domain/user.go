package domain

import (
	"fmt"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleForeman        = UserRole("foreman")
	RoleProjectManager = UserRole("project_manager")
	RoleClerk          = UserRole("clerk")
	RoleEngineer       = UserRole("engineer")
	RoleDirector       = UserRole("director")
	RoleRegular        = UserRole("regular")
	RoleAdmin          = UserRole("admin")
	RoleNil            = UserRole("")
)

var (
	Roles = []UserRole{RoleAdmin, RoleClerk, RoleDirector, RoleEngineer, RoleForeman, RoleProjectManager, RoleRegular}
)

type User struct {
	ID        uuid.UUID
	Version   int64
	Activated bool

	Email        *string
	PasswordHash *string
	Role         UserRole

	Name     string
	Surname  string
	LastName *string

	AvatarFileID *string

	LastActiveAt time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser() User {
	return User{}
}

func NewRegisterUser(name string, surname string, lastName *string) User {
	return User{
		ID:        UninitializedUUID,
		Version:   UninitializedVersion,
		Activated: false,
		Role:      UninitializedRole,
		Name:      name,
		Surname:   surname,
		LastName:  lastName,
	}
}

func NewUserLogin(email string, passwordHash string) User {
	return User{
		ID:           UninitializedUUID,
		Version:      UninitializedVersion,
		Activated:    true,
		Role:         UninitializedRole,
		Email:        &email,
		PasswordHash: &passwordHash,
	}
}

func (u *User) Validate() error {
	minimum, maximum := core_validation.NameMinLen, core_validation.NameMaxLen

	if err := core_validation.ValidateStringLength(&u.Name, "Name", minimum, maximum); err != nil {
		return err
	}

	if err := core_validation.ValidateStringLength(&u.Surname, "Surname", minimum, maximum); err != nil {
		return err
	}

	if u.LastName != nil {
		if err := core_validation.ValidateStringLength(u.LastName, "LastName", minimum, maximum); err != nil {
			return err
		}
	}

	if u.Role == RoleNil {
		return fmt.Errorf("`Role` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	if u.Activated {
		if _, err := core_validation.ValidateStringEmail(u.Email); err != nil {
			return err
		}

		if u.PasswordHash == nil {
			return fmt.Errorf("`PasswordHash` can't be NULL if `Activated` is true: %w", core_errors.ErrInvalidArgument)
		}
	}

	if u.CreatedAt.After(u.LastActiveAt) {
		return fmt.Errorf("`LastActiveAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	if u.CreatedAt.After(u.UpdatedAt) {
		return fmt.Errorf("`UpdatedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (u *User) GetFullName() string {
	if u.LastName == nil {
		return fmt.Sprintf("%s %s", u.Name, u.Surname)
	} else {
		return fmt.Sprintf("%s %s %s", u.Name, u.Surname, *u.LastName)
	}
}

type UserPatch struct {
	Activated    Nullable[bool]
	Email        Nullable[string]
	PasswordHash Nullable[string]
	Role         Nullable[UserRole]

	Name     Nullable[string]
	Surname  Nullable[string]
	LastName Nullable[string]

	AvatarFileID Nullable[string]
}

func NewUserRegisterPatch(email Nullable[string], passwordHash Nullable[string]) UserPatch {
	return UserPatch{
		Activated:    NewNullable(true),
		Email:        email,
		PasswordHash: passwordHash,
	}
}

func NewAvatarUserPatch(fileID Nullable[string]) UserPatch {
	return UserPatch{
		AvatarFileID: fileID,
	}
}

func (p *UserPatch) Validate() error {
	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf("'Name' can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.Surname.Set && p.Surname.Value == nil {
		return fmt.Errorf("'Surname' can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.Role.Set && p.Role.Value == nil {
		return fmt.Errorf("'Role' can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.Activated.Set && *p.Activated.Value {
		if p.Email.Set && p.Email.Value == nil {
			return fmt.Errorf("'Email' can't be patched to NULL if `Activated` is true: %w", core_errors.ErrInvalidArgument)
		}

		if p.PasswordHash.Set && p.PasswordHash.Value == nil {
			return fmt.Errorf("'PasswordHash' can't be patched to NULL if `Activated` is true: %w", core_errors.ErrInvalidArgument)
		}
	}

	if p.AvatarFileID.Set && p.AvatarFileID.Value != nil {
		if err := core_validation.ValidateID(*p.AvatarFileID.Value, FileIDByteLength); err != nil {
			return fmt.Errorf("validate avatar file id: %w", err)
		}
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) (User, error) {
	if err := patch.Validate(); err != nil {
		return User{}, fmt.Errorf("validate user patch: %w", err)
	}

	var patched User
	patched.Activated = u.Activated
	patched.Version = u.Version

	if patch.Email.Set {
		patched.Email = patch.Email.Value
	}

	if patch.PasswordHash.Set {
		patched.PasswordHash = patch.PasswordHash.Value
	}

	if patch.Role.Set {
		patched.Role = *patch.Role.Value
	}

	if patch.Name.Set {
		patched.Name = *patch.Name.Value
	}

	if patch.Surname.Set {
		patched.Surname = *patch.Surname.Value
	}

	if patch.LastName.Set {
		patched.LastName = patch.LastName.Value
	}

	if patch.Activated.Set {
		patched.Activated = *patch.Activated.Value
	}

	if patch.AvatarFileID.Set {
		patched.AvatarFileID = patch.AvatarFileID.Value
	}

	return patched, nil
}
