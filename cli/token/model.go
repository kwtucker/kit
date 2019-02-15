package main

const (
	getToken      = "https://identity.auth.theplatform.com/idm/web/Authentication/signIn?schema=1.0&form=json"
	getAccount    = "https://mps.theplatform.com/data/Account?schema=1.0&form=cjson"
	getMpxRole    = "http://access.auth.theplatform.com/data/Role?schema=1.0&form=cjson"
	getPermission = "http://access.auth.theplatform.com/data/Permission?schema=1.0&form=cjson"
)

// SignInMPX MPX's response from a auth call.
type signInMPX struct {
	SignInResponse signInDetails `json:"signInResponse"`
}

// SignInDetails Contains detailed info of the logged in user.
type signInDetails struct {
	Token       string `json:"token"`
	Duration    int    `json:"duration"`
	UserName    string `json:"userName"`
	UserID      string `json:"userId"`
	IdleTimeout int    `json:"idleTimeout"`
}

// ErrorMPX Default error response for MPX.
type errorMPX struct {
	ResponseCode  int    `json:"responseCode"`
	IsException   bool   `json:"isException"`
	Description   string `json:"description"`
	Title         string `json:"title"`
	CorrelationID string `json:"correlationId"`
}

type user struct {
	Username, password string
}

// RoleMPX ...
type roleMPX struct {
	ID              string          `json:"id"`
	GUID            string          `json:"guid"`
	Updated         int64           `json:"updated"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Added           int64           `json:"added"`
	OwnerID         string          `json:"ownerId"`
	AddedByUserID   string          `json:"addedByUserId"`
	UpdatedByUserID string          `json:"updatedByUserId"`
	Version         int             `json:"version"`
	Locked          bool            `json:"locked"`
	Disabled        bool            `json:"disabled"`
	RoleType        string          `json:"roleType"`
	AdminTags       []string        `json:"adminTags"`
	AllowOperations []operationsMPX `json:"allowOperations"`
	DenyOperations  []operationsMPX `json:"denyOperations"`
}

// OperationsMPX ...
type operationsMPX struct {
	Service  string `json:"service"`
	Instance string `json:"instance"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

// responseWrapperPerermission is the default values of a response from seattle json altcontent data service.
type responseWrapperPerermission struct {
	StartIndex   int             `json:"startIndex,omitempty"`
	ItemsPerPage int             `json:"itemsPerPage,omitempty"`
	EntryCount   int             `json:"entryCount,omitempty"`
	Entries      []permissionMPX `json:"entries,omitempty"`
}

// responseWrapperRole is the default values of a response from seattle json altcontent data service.
type responseWrapperRole struct {
	StartIndex   int        `json:"startIndex,omitempty"`
	ItemsPerPage int        `json:"itemsPerPage,omitempty"`
	EntryCount   int        `json:"entryCount,omitempty"`
	Entries      []*roleMPX `json:"entries,omitempty"`
}

// PermissionMPX ...
type permissionMPX struct {
	ID              string   `json:"id"`
	GUID            string   `json:"guid"`
	Updated         int64    `json:"updated"`
	Title           string   `json:"title"`
	Added           int64    `json:"added"`
	OwnerID         string   `json:"ownerId"`
	AddedByUserID   string   `json:"addedByUserId"`
	UpdatedByUserID string   `json:"updatedByUserId"`
	Version         int      `json:"version"`
	Locked          bool     `json:"locked"`
	Disabled        bool     `json:"disabled"`
	UserID          string   `json:"userId"`
	RoleIDs         []string `json:"roleIds"`
}
