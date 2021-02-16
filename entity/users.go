package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type contextKey string

var CtxUserSessionKey = contextKey("signed_in_user_session")

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	// FName, LName

	Username string `gorm:"type:varchar(255);not null" json:"username,omitempty" bson:"username,omitempty"`
	Email    string `gorm:"type:varchar(255);not null; unique" json:"email,omitempty" bson:"email,omitempty"`
	Phone    string `gorm:"type:varchar(100);not null; unique" json:"phone,omitempty" bson:"phone,omitempty"`
	Password string `gorm:"type:varchar(255)"`
	Role     string
	Roles []string
	Session []Session


}

type LoginData struct {
	LoginInfo string `json:"info"`
	Password  string
	InfoType  string `json:"info_type"`
}


const (
	Client          = "CLIENT"
	Root             = "ROOT"
	Admin 			= "ADMIN"


)


var RoleDescription = map[string]string{
	Client: "A client is a registered user & can have favorites",
	Admin: "Admins manages the system",
	Root: "Root have all Permissions",
}




//Session represents login user session
type Session struct {
	UUID          string `gorm:"type:varchar(255);not null"`
	LoginDate     time.Time
	//optional
	Expires    int64  `gorm:"type:varchar(255);not null"`
	SigningKey []byte `gorm:"type:varchar(255);not null"`
}

