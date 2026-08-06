package domain

import "github.com/google/uuid"

var (
	UninitializedUUID      uuid.UUID = uuid.Nil
	UninitializedID        string    = ""
	UninitializedVersion   int64     = -1
	UninitializedRole      UserRole  = RoleRegular
	UninitializedTokenType TokenType = TokenTypeNil
	UninitializedFileTag   FileTag   = FileTagNil
)
