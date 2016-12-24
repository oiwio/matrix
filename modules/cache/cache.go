package cache

import (
	// "matrix/modules/db"
	// "matrix/modules/event"
	// "time"
)

// func SetUserRegistry(registry *event.UserRegistry) error {
// 	key := CacheKeyUserRegistry(registry.Phone)
// 	return set(key, registry,"EX","3600")
// }

// func GetUserRegistry(phone string) (*event.UserRegistry, error) {
// 	key := CacheKeyUserRegistry(phone)
// 	registry := new(event.UserRegistry)
// 	err := _cache.Get(key, &registry)
// 	return registry, err
// }

// func DelUserRegistry(registry *event.UserRegistry) error {
// 	key := CacheKeyUserRegistry(registry.DeviceToken)
// 	return _cache.Delete(key)
// }

// func SetVerifyCode(code string, account string) error {
// 	key := CacheKeyVerifyAccount(account)
// 	return _cache.Set(key, code, time.Duration(time.Minute*30))
// }

// func GetVerifyCode(account string) (string, error) {
// 	code := ""
// 	key := CacheKeyVerifyAccount(account)
// 	err := _cache.Get(key, &code)
// 	return code, err
// }

// func DelVerifyCode(account string) error {
// 	key := CacheKeyVerifyAccount(account)
// 	return _cache.Delete(key)
// }

// func SetEmailActive(code string, account string) error {
// 	key := CacheKeyEmailActivate(code)
// 	return _cache.Set(key, account, time.Duration(time.Hour*24*7))
// }

// func GetEmailActive(code string) (string, error) {
// 	account := ""
// 	key := CacheKeyEmailActivate(code)
// 	err := _cache.Get(key, &account)
// 	return account, err
// }

// func DelEmailActive(code string) error {
// 	key := CacheKeyEmailActivate(code)
// 	return _cache.Delete(key)
// }

// func SetEmailVerify(code string, account string) error {
// 	key := CacheKeyEmailVerifyCode(code)
// 	return _cache.Set(key, account, time.Duration(time.Minute*30))
// }

// func GetEmailVerify(code string) (string, error) {
// 	account := ""
// 	key := CacheKeyEmailVerifyCode(code)
// 	err := _cache.Get(key, &account)
// 	return account, err
// }

// func DelEmailVerify(code string) error {
// 	key := CacheKeyEmailVerifyCode(code)
// 	return _cache.Delete(key)
// }

// func SetUserByShortId(user *db.User) error {
//   var err error
//   key := CacheKeyShortId(user.ShortId)
//   err = _cache.Set(key, user, _cache.FOREVER)
//   return err
// }

// func GetUserByShortId(shortId string) (*db.User, error) {
//   var err error
//   key := CacheKeyShortId(shortId)
//   user := new(db.User)
//   err = _cache.Get(key, &user)
//   return user, err
// }

// func SetUserByObjectId(user *db.User) error {
// 	key := CacheKeyUserObjectId(user.UserId.Hex())
// 	return _cache.Set(key, user, _cache.FOREVER)
// }

// func GetUserByObjectId(id string) (*db.User, error) {
// 	key := CacheKeyUserObjectId(id)
// 	user := new(db.User)
// 	err := _cache.Get(key, &user)
// 	return user, err
// }

// func DelUserByObjectId(id string) error {
// 	key := CacheKeyUserObjectId(id)
// 	return _cache.Delete(key)
// }

// func SetUserByAccessToken(user *db.User) error {
// 	key := CacheKeyUserAccessToken(user.AccessToken)
// 	return _cache.Set(key, user, _cache.FOREVER)
// }

// func GetUserByAccessToken(token string) (*db.User, error) {
// 	key := CacheKeyUserAccessToken(token)
// 	user := new(db.User)
// 	err := _cache.Get(key, &user)
// 	return user, err
// }

// func DelUserByAccessToken(token string) error {
// 	key := CacheKeyUserAccessToken(token)
// 	return _cache.Delete(key)
// }

// func SetUserLocation(userId string, location string) error {
// 	key := CacheKeyUserLocation(userId)
// 	return _cache.Set(key, location, _cache.FOREVER)
// }
//
// func GetUserLocation(userId string) (string, error) {
// 	location := ""
// 	key := CacheKeyUserLocation(userId)
// 	err := _cache.Get(key, &location)
// 	return location, err
// }
//
// func MultiGetUserLocation(ids ...string) []string {
//
// 	var err error
//
// 	keys := make([]string, len(ids))
// 	locations := make([]string, len(ids))
// 	for index, id := range ids {
// 		keys[index] = CacheKeyUserLocation(id)
// 	}
//
// 	getter, _ := _cache.GetMulti(keys...)
// 	for index, key := range keys {
// 		err = getter.Get(key, locations[index])
// 		if err != nil {
// 			locations[index] = "0,0"
// 		}
// 	}
//
// 	return locations
// }

// func BatchGetUsersByShortId(ids ...string) []*db.User {
//   keys := make([]string, len(ids))
//   for index, id := range ids {
//     keys[index] = CacheKeyShortId(id)
//   }
//   getter, _ := _cache.GetMulti(keys...)
//   users := make([]*db.User, len(ids))
//   for index, key := range keys {
//     _ = getter.Get(key, users[index])
//   }
//   return users
// }

// func stringChannelIterator(strArray) <-chan string {
//   ch := make(chan string)
//   go func() {
//     for _, val := range strArray {
//       ch <- val
//     }
//     close(ch)
//   }()
//   return ch
// }

// func SetSeedUserIds(gender string, seedIds []string) error {
// 	return _cache.Set(CacheKeySeedUser(gender), seedIds, _cache.FOREVER)
// }

// func GetSeedUserIds(gender string) ([]string, error) {
// 	seedIds := []string{}
// 	err := _cache.Get(CacheKeySeedUser(gender), &seedIds)
// 	return seedIds, err
// }

// func SetFollowerByObjectId(id string, followers []string) error {
// 	err := _cache.Set(CacheKeyFollowerObjectId(id), followers, _cache.FOREVER)
// 	return err
// }

// func GetFollowerByObjectId(id string) ([]string, error) {
// 	followers := []string{}
// 	err := _cache.Get(CacheKeyFollowerObjectId(id), &followers)
// 	return followers, err
// }

// func SetFollowingByObjectId(id string, followings []string) error {
// 	err := _cache.Set(CacheKeyFollowingObjectId(id), followings, _cache.FOREVER)
// 	return err
// }

// func GetFollowingByObjectId(id string) ([]string, error) {
// 	followings := []string{}
// 	err := _cache.Get(CacheKeyFollowingObjectId(id), &followings)
// 	return followings, err
// }

// 用于缓存官方topic中
// func SetUsersForOfficialTopic(topicId string, joinedIds []string) error {
// 	return _cache.Set(CacheKeyOfficialTopicId(topicId), joinedIds, _cache.FOREVER)
// }

// func GetUsersFromOfficialTopic(topicId string) ([]string, error) {
// 	joined := []string{}
// 	err := _cache.Get(CacheKeyOfficialTopicId(topicId), &joined)
// 	return joined, err
// }

// func SetUserViewHotFeedCursor(userId string, feeds []*db.Feed, period time.Duration) error {
// 	return _cache.Set(CacheKeyUserViewHotFeedCursor(userId), feeds, period)
// }

// func GetUserViewHotFeedCursor(userId string) ([]*db.Feed, error) {
// 	feeds := []*db.Feed{}
// 	err := _cache.Get(CacheKeyUserViewHotFeedCursor(userId), &feeds)
// 	return feeds, err
// }

// func ReplaceUserViewHotFeedCursor(userId string, feeds []*db.Feed, period time.Duration) error {
// 	return _cache.Replace(CacheKeyUserViewHotFeedCursor(userId), feeds, period)
// }

// func DeleteUserViewHotFeedCursor(userId string) error {
// 	return _cache.Delete(CacheKeyUserViewHotFeedCursor(userId))
// }

// func SetHotFeedIds(userId string, feedIds []string, period time.Duration) error {
// 	return _cache.Set(CacheKeyHotFeedFeedIds(userId), feedIds, period)
// }

// func GetHotFeedIds(userId string) ([]string, error) {
// 	feedIds := []string{}
// 	err := _cache.Get(CacheKeyHotFeedFeedIds(userId), &feedIds)
// 	return feedIds, err
// }

// func ReplaceHotFeedIds(userId string, feedIds []string, period time.Duration) error {
// 	return _cache.Replace(CacheKeyHotFeedFeedIds(userId), feedIds, period)
// }

// func SetHotUserIds(userId string, userIds []string, period time.Duration) error {
// 	return _cache.Set(CacheKeyHotFeedUserIds(userId), userIds, period)
// }

// func GetHotUserIds(userId string) ([]string, error) {
// 	userIds := []string{}
// 	err := _cache.Get(CacheKeyHotFeedUserIds(userId), &userIds)
// 	return userIds, err
// }

// func ReplaceHotUserIds(userId string, userIds []string, period time.Duration) error {
// 	return _cache.Replace(CacheKeyHotFeedUserIds(userId), userIds, period)
// }

// func SetNearbyFeedIdsByUserId(userId string, Ids []string, period time.Duration) error {
// 	return _cache.Set(CacheKeyNearByFeedIds(userId), Ids, period)
// }

// func GetNearbyFeedIdsByUserId(userId string) ([]string, error) {
// 	Ids := []string{}
// 	err := _cache.Get(CacheKeyNearByFeedIds(userId), &Ids)
// 	return Ids, err
// }

// func SetNearbyDistanceByUserId(userId string, distance float64, period time.Duration) error {
// 	return _cache.Set(CacheKeyNearbyDistance(userId), distance, period)
// }

// func GetNearbyDistanceByUserId(userId string) (float64, error) {
// 	var distance float64
// 	err := _cache.Get(CacheKeyNearbyDistance(userId), &distance)
// 	return distance, err
// }
