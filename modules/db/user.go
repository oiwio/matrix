package db

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

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

		Nickname    string     `json:"nickname,omitempty" bson:"nickname,omitempty"`
		Gender      int        `json:"gender,omitempty" bson:"gender,omitempty"`
		Orientation int        `json:"orientation,omitempty" bson:"orientation,omitempty"` //0男 1女 2双性
		Avatar      string     `json:"avatar,omitempty" bson:"avatar,omitempty"`
		Tags        []*UserTag `json:"tags,omitempty" bson:"tags,omitempty"`
		Signature   string     `json:"signature,omitempty" bson:"signature,omitempty"`
		Coordinate  []float64  `json:"coordinate,omitempty" bson:"coordinate,omitempty"` //坐标

		Birthday      int64  `json:"birthday,omitempty" bson:"birthday,omitempty"`
		Constellation string `json:"constellation,omitempty" bson:"constellation,omitempty"` //星座
		Age           int    `json:"age,omitempty" bson:"age,omitempty"`
		Height        int    `json:"height,omitempty" bson:"height,omitempty"`
		Weight        int    `json:"weight,omitempty" bson:"weight,omitempty"`

		Region string `json:"region,omitempty" bson:"region,omitempty"`

		// 黑名单
		BlockList []string `json:"blockList,omitempty" bson:"blockList,omitempty"`

		CreateAt int64 `json:"createAt,omitempty" bson:"createAt,omitempty"`
		UpdateAt int64 `json:"updateAt,omitempty" bson:"updateAt,omitempty"`

		Status   *UserStatus   `json:"status,omitempty" bson:"status,omitempty"`
		Settings *UserSettings `json:"settings,omitempty" bson:"settings,omitempty"`

		OpenAccounts []*OpenAccount `json:"openIds,omitempty" bson:"openIds,omitempty"`

		//用户授权
		AccessToken string `json:"accessToken,omitempty" bson:"accessToken,omitempty"`
		RefreshToken string `json:"refreshToken" bson:"refreshToken"`
		ExpiresIn    int64  `json:"expiresIn" bson:"expiresIn"`
	}

	UserTag struct {
		TagId   bson.ObjectId `json:"tagId,omitempty" bson:"tagId"`
		TagName string        `json:"feedName,omitempty" bson:"feedName,omitempty"`
	}

	UserStatus struct {
		Following int64 `json:"following,omitempty" bson:"following,omitempty"`
		Follower  int64 `json:"follower,omitempty" bson:"follower,omitempty"`
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
func NewUser(s *mgo.Session, user *User) (*User, error) {
	var (
		err error
	)

	user.UserId = bson.NewObjectId()

	//user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
	user.Nickname = GenerateNickname(1)
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
	u.Height = user.Height
	u.Weight = user.Weight
	u.Region = user.Region
	u.Tags = user.Tags

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

func GetUserByAccessToken(s *mgo.Session, token string) (*User, error) {

	var (
		err error
	)

	user := new(User)
	err = Collection(s, user).Find(bson.M{"accessToken": token}).One(user)
	if err == nil {
		return user, err
	}
	return nil, errors.New(fmt.Sprintf("Can not found user with phone no. %v | Reason:%v", token, err.Error()))
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

//根据语言生成昵称,0 英文，1 中文，2 俄文，3 火星文
func GenerateNickname(langurage int) string {
	cnFirstName := [...]string{"可爱的", "变态的", "害羞的", "性感的", "迷人的", "火辣的", "无聊的"}
	cnLastName := [...]string{"佐助", "鸣人", "路飞", "白胡子", "小樱", "柯南", "一休", "佐为", "卡卡西", "艾斯"}
	enFirstName := [...]string{"lucky", "fucking", "angry"}
	enLastName := [...]string{"dog", "cat", "girl"}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	switch langurage {
	case 0:
		return enFirstName[r.Intn(len(enFirstName))] + enLastName[r.Intn(len(enLastName))]
	case 1:
		return cnFirstName[r.Intn(len(cnFirstName))] + cnLastName[r.Intn(len(cnLastName))]
	}
	return cnFirstName[r.Intn(len(cnFirstName))] + cnLastName[r.Intn(len(cnLastName))]
}
