package model

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"

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
	Role     int64  `json:"role"`                                      // user's role
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
	Permission uint64 `json:"permission"`
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

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	var res []webauthn.Credential
	err := json.Unmarshal([]byte(u.Authn), &res)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func (u *User) SetPassword(pwd string) *User {
	u.Salt = utils.RandString(16)
	u.PwdHash = TwoHashPwd(pwd, u.Salt)
	u.PwdTS = time.Now().Unix()
	return u
}

func (u *User) CanSeeHides() bool {
	return u.Permission&1 == 1
}

func (u *User) CanAccessWithoutPassword() bool {
	return (u.Permission>>1)&1 == 1
}

func (u *User) CanAddOfflineDownloadTasks() bool {
	return (u.Permission>>2)&1 == 1
}

func (u *User) CanWrite() bool {
	return (u.Permission>>3)&1 == 1
}

func (u *User) CanRename() bool {
	return (u.Permission>>4)&1 == 1
}

func (u *User) CanMove() bool {
	return (u.Permission>>5)&1 == 1
}

func (u *User) CanCopy() bool {
	return (u.Permission>>6)&1 == 1
}

func (u *User) CanRemove() bool {
	return (u.Permission>>7)&1 == 1
}

func (u *User) CanWebdavRead() bool {
	return (u.Permission>>8)&1 == 1
}

func (u *User) CanWebdavManage() bool {
	return (u.Permission>>9)&1 == 1
}

func (u *User) CanFTPAccess() bool {
	return (u.Permission>>10)&1 == 1
}

func (u *User) CanFTPManage() bool {
	return (u.Permission>>11)&1 == 1
}
