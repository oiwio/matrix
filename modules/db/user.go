package db

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	User struct {
		UserId      bson.ObjectId `json:"userId,omitempty" bson:"_id"`
		DeviceToken string        `json:"deviceToken,omitempty" bson:"deviceToken,omitempty"`

		Username       string `json:"username,omitempty" bson:"username,omitempty"`
		HashedPassword []byte `json:"hashedPassword,omitempty" bson:"hashedPassword,omitempty"`

		Email         string `json:"email,omitempty" bson:"email,omitempty"`
		EmailVerified bool   `json:"emailVerfied,omitempty" bson:"emailVerified,omitempty"`
		Phone         string `json:"phone,omitempty" bson:"phone,omitempty"`
		PhoneVerified bool   `json:"phoneVerified,omitempty" bson:"phoneVerified,omitempty"`

		Nickname   string    `json:"nickname,omitempty" bson:"nickname,omitempty"`
		Gender     int       `json:"gender,omitempty" bson:"gender,omitempty"` //女0 男1
		Avatar     string    `json:"avatar,omitempty" bson:"avatar,omitempty"`
		Signature  string    `json:"signature,omitempty" bson:"signature,omitempty"`
		Coordinate []float64 `json:"coordinate,omitempty" bson:"coordinate,omitempty"` //坐标

		Birthday      int64  `json:"birthday,omitempty" bson:"birthday,omitempty"`
		Constellation string `json:"constellation,omitempty" bson:"constellation,omitempty"` //星座
		Age           int    `json:"age,omitempty" bson:"age,omitempty"`

		Following int64 `json:"following,omitempty" bson:"following,omitempty"`
		Follower  int64 `json:"follower,omitempty" bson:"follower,omitempty"`

		Region string `json:"region,omitempty" bson:"region,omitempty"`

		// 黑名单
		BlockList []string `json:"blockList,omitempty" bson:"blockList,omitempty"`

		CreateAt int64 `json:"createAt,omitempty" bson:"createAt,omitempty"`
		UpdateAt int64 `json:"updateAt,omitempty" bson:"updateAt,omitempty"`

		Settings *UserSettings `json:"settings,omitempty" bson:"settings,omitempty"`

		OpenAccounts []*OpenAccount `json:"openIds,omitempty" bson:"openIds,omitempty"`
	}

	UserSettings struct {
		Whisper          bool `json:"whisper" bson:"whisper"`                   // 私聊，true: 允许；false：不允许
		NewMessageNotify bool `json:"newMessageNotify" bson:"newMessageNotify"` // 新消息提醒，true：开启；false：关闭
		SoundNotify      bool `json:"soundNotify" bson:"soundNotify"`           // 声音提醒，true：开启；false：关闭
		ShakeNotify      bool `json:"shakeNotify" bson:"shakeNotify"`           // 震动提示，true：开启；false：关闭
		PushDetail       bool `json:"pushDetail" bson:"pushDetail"`             // 通知详情，true：开启；false：关闭
		DoNotDisturb     bool `json:"doNotDisturb" bson:"doNotDisturb"`         // 免扰，true：开启；false：关闭
		DNDPeriods       bool `json:"dndPeriods" bson:"dndPeriods"`             // 按时段免扰，true：开启；false：关闭
	}

	OpenAccount struct {
		AppName      string `json:"appName" bson:"appName"`
		OpenID       string `json:"openId" bson:"openId"`
		Avatar       string `json:"avatar" bson:"avatar"`
		Signature    string `json:"signature" bson:"signature"`
		AccessToken  string `json:"accessToken,omitempty" bson:"accessToken,omitempty"`
		RefreshToken string `json:"refreshToken,omitempty" bson:"refreshToken,omitempty"`
		ExpiresIn    int64  `json:"expiresIn,omitempty" bson:"expiresIn,omitempty"`
		RemindIn     int64  `json:"remindIn,omitempty" bson:"remindIn,omitempty"`
		CreateAt     int64  `json:"createAt" bson:"createAt"`
	}
)

/**
 * 新建用户
 */
func NewUser(s *mgo.Session, user *User, password string) (*User, error) {
	var (
		err error
	)

	user.UserId = bson.NewObjectId()
	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.CreateAt = time.Now().Unix()
	user.UpdateAt = user.CreateAt

	settings := new(UserSettings)
	settings.Whisper = true
	settings.NewMessageNotify = true
	settings.SoundNotify = true
	settings.ShakeNotify = true
	settings.PushDetail = true
	settings.DoNotDisturb = true
	settings.DNDPeriods = true

	user.Settings = settings
	collection := Collection(s, user)
	err = collection.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, err
}

// 更新用户表
// func UpdateUser(client *elastic.Client, user *User) error {
//
// 	return err
// }

// 删除用户
func DeleteUser(s *mgo.Session, userId string) error {
	var (
		err  error
		user *User
	)
	err = Collection(s, user).Remove(bson.M{"_id": userId})
	return err
}

// 检查用户是否存在
func IsUserExist(s *mgo.Session, userId bson.ObjectId) bool {

	var (
		err error
	)

	// exist, err := client.Exists().Index(gender).Type("User").Id(UserId).Do()
	user := new(User)
	err = Collection(s, user).FindId(userId).One(user)
	if err != nil {
		return false
	}
	return true

}

// 更新用户密码
func UpdateUserPassword(s *mgo.Session, user *User) error {

	var (
		err error
	)

	u := new(User)
	collection := Collection(s, user)
	err = collection.FindId(user.UserId).One(u)
	if err != nil {
		return err
	}

	u.HashedPassword = user.HashedPassword
	err = collection.Update(bson.M{"_id": user.UserId}, u)
	if err != nil {
		return err
	}

	return err
}

// 更新用户资料
func UpdateUserProfile(s *mgo.Session, user *User) error {

	var err error

	u := new(User)
	collection := Collection(s, u)

	err = collection.FindId(user.UserId).One(u)
	if err != nil {
		return err
	}

	if len(user.Nickname) > 0 {
		u.Nickname = user.Nickname
	}

	u.Birthday = user.Birthday
	u.Constellation = user.Constellation
	u.Age = user.Age

	err = collection.Update(bson.M{"_id": user.UserId}, u)
	if err != nil {
		return err
	}

	return err
}

// 更新用户设置
func UpdateUserSettings(s *mgo.Session, user *User) error {

	var (
		err error
	)

	u := new(User)
	collection := Collection(s, user)
	err = collection.FindId(user.UserId).One(u)
	if err != nil {
		return err
	}

	// log.Println(u.Settings, user.Settings)
	u.Settings = user.Settings
	err = collection.Update(bson.M{"_id": user.UserId}, u)
	if err != nil {
		return err
	}

	return err
}

/**
 * 更新用户签名
 */
func UpdateUserSignature(s *mgo.Session, user *User) error {

	var (
		err error
	)

	u := new(User)
	collection := Collection(s, user)
	err = collection.FindId(user.UserId).One(u)
	if err != nil {
		return err
	}

	u.Signature = user.Signature
	err = collection.Update(bson.M{"_id": user.UserId}, u)
	if err != nil {
		return err
	}

	return err
}

/**
 * 更新设备 Token
 */
func UpdateDeviceToken(s *mgo.Session, user *User) error {

	var (
		err        error
		collection *mgo.Collection
		u          *User
	)

	u = new(User)
	collection = Collection(s, user)

	err = collection.FindId(user.UserId).One(u)
	if err != nil {
		return err
	}

	u.DeviceToken = user.DeviceToken
	u.UpdateAt = time.Now().Unix()
	err = collection.Update(bson.M{"_id": user.UserId}, u)
	if err != nil {
		return err
	}

	return err
}

// 根据Id获取用户
func GetUserById(s *mgo.Session, userId bson.ObjectId) (*User, error) {
	var (
		err error
	)

	user := new(User)
	err = Collection(s, user).FindId(userId).One(user)
	fmt.Println(err)
	return user, err
}

func GetUserByDeviceToken(s *mgo.Session, deviceToken string) (*User, error) {

	var (
		err error
	)

	user := new(User)
	err = Collection(s, user).Find(bson.M{"deviceToken": deviceToken}).One(user)
	if err == nil {
		return user, err
	}
	return nil, errors.New(fmt.Sprintf("Can not found user with phone no. %v | Reason:%v", deviceToken, err.Error()))
}

//根据手机号搜索用户
func GetUserByPhone(s *mgo.Session, phone string) (*User, error) {
	var (
		err error
	)

	user := new(User)
	err = Collection(s, user).Find(bson.M{"phone": phone}).One(user)
	if err == nil {
		return user, err
	}
	return nil, errors.New(fmt.Sprintf("Can not found user with phone no. %v", phone))
}

//根据邮箱搜索用户
func GetUserByEmail(s *mgo.Session, email string) (*User, error) {
	var (
		err error
	)

	user := new(User)
	err = Collection(s, user).Find(bson.M{"email": email}).One(user)
	if err == nil {
		return user, err
	}

	return nil, errors.New(fmt.Sprintf("Can not found user with email. %v", email))
}
