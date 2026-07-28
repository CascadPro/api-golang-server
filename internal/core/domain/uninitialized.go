package domain

import "github.com/google/uuid"

var (
	UninitializedUUID      = uuid.Nil
	UninitializedID        = ""
	UninitializedVersion   = -1
	UninitializedRole      = RoleRegular
	UninitializedTokenType = TokenTypeNil
	UninitializedFileTag   = FileTagNil
)
