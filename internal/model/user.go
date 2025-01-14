package model

import (
	"encoding/binary"
	"fmt"
	"time"

	"ndm/internal/utils"
)

const (
	GENERAL = iota
	GUEST   // only one exists
	ADMIN
)

const StaticHashSalt = "https://github.com/midoks/ndm"

type User struct {
	ID       int64  `json:"id" gorm:"primaryKey"`                      // unique key
	Username string `json:"username" gorm:"unique" binding:"required"` // username
	PwdHash  string `json:"-"`                                         // password hash
	PwdTS    int64  `json:"-"`                                         // password timestamp
	Salt     string `json:"salt"`                                      // unique salt
	Password string `json:"password"`                                  // password
	BasePath string `json:"base_path"`                                 // base path
	Role     int    `json:"role"`                                      // user's role
	Disabled bool   `json:"disabled"`
	// Determine permissions by bit
	//   0:  can see hidden files
	//   1:  can access without password
	//   2:  can add offline download tasks
	//   3:  can mkdir and upload
	//   4:  can rename
	//   5:  can move
	//   6:  can copy
	//   7:  can remove
	//   8:  webdav read
	//   9:  webdav write
	//   10: ftp/sftp login and read
	//   11: ftp/sftp write
	Permission int64  `json:"permission"`
	OtpSecret  string `json:"-"`
	SsoID      string `json:"sso_id"` // unique by sso platform
	Authn      string `gorm:"type:text" json:"-"`
}

func (u *User) IsGuest() bool {
	return u.Role == GUEST
}

func (u *User) IsAdmin() bool {
	return u.Role == ADMIN
}

func StaticHash(password string) string {
	return utils.HashData(utils.SHA256, []byte(fmt.Sprintf("%s-%s", password, StaticHashSalt)))
}

func HashPwd(static string, salt string) string {
	return utils.HashData(utils.SHA256, []byte(fmt.Sprintf("%s-%s", static, salt)))
}

func TwoHashPwd(password string, salt string) string {
	return HashPwd(StaticHash(password), salt)
}

func (u *User) WebAuthnID() []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(u.ID))
	return bs
}

func (u *User) WebAuthnName() string {
	return u.Username
}

func (u *User) WebAuthnDisplayName() string {
	return u.Username
}

func (u *User) SetPassword(pwd string) *User {
	u.Salt = utils.RandString(16)
	u.PwdHash = TwoHashPwd(pwd, u.Salt)
	u.PwdTS = time.Now().Unix()
	return u
}
