package cache

import (
	"fmt"
)

const (
	CACHE_BASE           = 0x80000000
	USER_REGISTRY_KEY    = CACHE_BASE + 1
	USER_BASE            = CACHE_BASE + 2
	USER_ACCESS_TOKEN    = CACHE_BASE + 3
	USER_VERIFY          = CACHE_BASE + 4
	USER_EMAIL_ACTIVATE  = CACHE_BASE + 5
	USER_EMAIL_VERIFY    = CACHE_BASE + 6
	USER_DEVICE          = CACHE_BASE + 7
	USER_FOLLOWER        = CACHE_BASE + 8
	USER_FOLLOWING       = CACHE_BASE + 9
	USER_SEED            = CACHE_BASE + 110
	USER_LOCATION        = CACHE_BASE + 10
	FEED_HOT             = CACHE_BASE + 12
	FEED_HOT_FEEDIDS     = CACHE_BASE + 13
	FEED_HOT_USERIDS     = CACHE_BASE + 14
	TOPIC                = CACHE_BASE + 15
	FEED_NEARBY_IDS      = CACHE_BASE + 16
	FEED_NEARBY_DISTANCE = CACHE_BASE + 17
)

func CacheKeyUserRegistry(deviceToken string) string {
	return fmt.Sprintf("%d-%s", USER_REGISTRY_KEY, deviceToken)
}

func CacheKeyUserObjectId(objectId string) string {
	return fmt.Sprintf("%d-%s", USER_BASE, objectId)
}

func CacheKeyUserAccessToken(token string) string {
	return fmt.Sprintf("%v-%v", USER_ACCESS_TOKEN, token)
}

func CacheKeyUserLocation(userId string) string {
	return fmt.Sprintf("%d-%s", USER_LOCATION, userId)
}

func CacheKeySeedUser(gender string) string {
	return fmt.Sprintf("%v-%v", USER_SEED, gender)
}

func CacheKeyUserDeviceToken(deviceToken string) string {
	return fmt.Sprintf("%d-%s", USER_DEVICE, deviceToken)
}

func CacheKeyVerifyAccount(account string) string {
	return fmt.Sprintf("%d-%s", USER_VERIFY, account)
}

func CacheKeyEmailActivate(code string) string {
	return fmt.Sprintf("%d-%s", USER_EMAIL_ACTIVATE, code)
}

func CacheKeyEmailVerifyCode(code string) string {
	return fmt.Sprintf("%d-%s", USER_EMAIL_VERIFY, code)
}

// func CacheKeyShortId(shortId string) string {
//   return fmt.Sprintf("%d-%s", USER_USERNAME, shortId)
// }

func CacheKeyFollowerObjectId(id string) string {
	return fmt.Sprintf("%v-%v", USER_FOLLOWER, id)
}

func CacheKeyFollowingObjectId(id string) string {
	return fmt.Sprintf("%v-%v", USER_FOLLOWING, id)
}

func CacheKeyOfficialTopicId(id string) string {
	return fmt.Sprintf("%d-%s", TOPIC, id)
}

// func CacheKeyBlockObjectId(id string) string {
//   return fmt.Sprintf("%v-%v", USER_BLOCK, id)
// }

// func CacheKeyCloseObjectId(id string) string {
//   return fmt.Sprintf("%v-%v", USER_CLOSE, id)
// }

func CacheKeyUserViewHotFeedCursor(id string) string {
	return fmt.Sprintf("%d-%s", FEED_HOT, id)
}

/**
*   FEED_HOT_FEEDIDS    = CACHE_BASE + 13
 FEED_HOT_USERIDS    = CACHE_BASE + 14
*/
func CacheKeyHotFeedFeedIds(id string) string {
	return fmt.Sprintf("%v-%v", FEED_HOT_FEEDIDS, id)
}

func CacheKeyHotFeedUserIds(id string) string {
	return fmt.Sprintf("%v-%v", FEED_HOT_USERIDS, id)
}

func CacheKeyNearByFeedIds(id string) string {
	return fmt.Sprintf("%v-%v", FEED_NEARBY_IDS, id)
}

func CacheKeyNearbyDistance(id string) string {
	return fmt.Sprintf("%v-%v", FEED_NEARBY_DISTANCE, id)
}
