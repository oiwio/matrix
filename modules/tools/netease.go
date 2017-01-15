package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"matrix/modules/db"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

const (
	gNeteaseRetOk             = 200
	gNeteaseProvider          = "http://music.163.com/"
	gNeteaseAPIUrlBase        = "http://music.163.com/api"
	gNeteaseAlbumUrl          = "/album/"
	gNeteaseSongListUrl       = "/song/detail?ids=[%s]"
	gNeteasePlayListUrl       = "/playlist/detail?id="
	gNeteaseEIDCacheKeyPrefix = "163eid:" // encrypted dfsId
	gNeteaseMusicCDNUrlF      = "http://p1.music.126.net/%s/%s.mp3"
)

func GetUrl(client *http.Client, url string) []byte {

	// cache missed, do http request
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("error get url %s: %s", url, err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error getting response body from url %s: %s", url, err)
		return nil
	}
	// fmt.Printf(string(body))
	return body
}

func PostUrl(client *http.Client, songurl string, songname string, limit string, offset string) []byte {
	songparam := "s=" + songname + "&limit=" + limit + "&type=1&offset=" + offset
	resp, err := client.Post(songurl, "application/x-www-form-urlencoded", strings.NewReader(songparam))
	if err != nil {
		log.Printf("error get url %s: %s", songurl, err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error getting response body from url %s: %s", songurl, err)
		return nil
	}
	// fmt.Printf(string(body))
	return body
}

var (
	gNeteaseClient      = &http.Client{}
	gNeteaseEIDReplacer = strings.NewReplacer(
		"/", "_",
		"+", "-",
	)
	// bypass cross-domain problem
	gNeteaseUrlReplacer = strings.NewReplacer(
		"http://m", "http://p",
	)
)

func init() {
	// init netease http client
	cookies, err := cookiejar.New(nil)
	if nil != err {
		fmt.Printf("failed to init netease httpclient cookiejar: %s", err)
	}
	apiUrl, err := url.Parse(gNeteaseAPIUrlBase)
	if nil != err {
		fmt.Printf("failed to parse netease api url %s: %s", gNeteaseAPIUrlBase, err)
	}
	// netease api requires some cookies to work
	cookies.SetCookies(apiUrl, []*http.Cookie{
		&http.Cookie{Name: "appver", Value: "1.4.1.62460"},
		&http.Cookie{Name: "os", Value: "pc"},
		&http.Cookie{Name: "osver", Value: "Microsoft-Windows-7-Ultimate-Edition-build-7600-64bit"},
	})
	gNeteaseClient.Jar = cookies
}

type NeteaseRetStatus struct {
	StatusCode int    `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
}

type NeteaseAlbum struct {
	Songs    []NeteaseSong `json:"songs,omitempty"`
	CoverUrl string        `json:"picUrl,omitempty"`
}

type NeteaseSongListRet struct {
	NeteaseRetStatus
	Songs []NeteaseSong `json:"songs"`
}

type SearchSong struct {
	Id       int             `json:"id"`
	Name     string          `json:"name"`
	Artists  []NeteaseArtist `json:"artists,omitempty"`
	Duration int             `json:"duration"`
}

type SearchResult struct {
	Songs []SearchSong `json:"songs"`
}

type SearchResultRet struct {
	// NeteaseRetStatus
	Result SearchResult `json:"result"`
}

type NeteaseSong struct {
	Id       int             `json:"id,omitempty"`
	Artists  []NeteaseArtist `json:"artists,omitempty"`
	Album    NeteaseAlbum    `json:"album,omitempty"`
	Duration int             `json:"duration,omitempty"`
	Name     string          `json:"name,omitempty"`
	Url      string          `json:"mp3Url,omitempty"`

	HighQualityMusic   NeteaseMusicDetail `json:"hMusic,omitempty"`
	MediumQualityMusic NeteaseMusicDetail `json:"mMusic,omitempty"`
	LowQualityMusic    NeteaseMusicDetail `json:"lMusic,omitempty"`
}

func (song *NeteaseSong) UpdateUrl(quality string) *NeteaseSong {
	if "" == quality || quality == "medium" {
		song.Url = gNeteaseUrlReplacer.Replace(song.Url)
		return song
	}
	musicDetail := &song.HighQualityMusic
	if quality == "low" {
		musicDetail = &song.LowQualityMusic
	}
	song.Url = musicDetail.MakeUrl()
	return song
}

type NeteaseMusicDetail struct {
	Bitrate int `json:"bitrate"`
	DfsID   int `json:"dfsId"`
}

func (md *NeteaseMusicDetail) MakeUrl() string {
	strDfsID := strconv.Itoa(md.DfsID)
	fmt.Print(md.DfsID)
	// load eid from cache first
	// eidKey := gCacheKeyPrefix + gNeteaseEIDCacheKeyPrefix + strDfsID
	// eid := GetCache(eidKey, false)
	// if nil == eid {
	// build encrypted dfsId, see https://github.com/yanunon/NeteaseCloudMusic/wiki/网易云音乐API分析#歌曲id加密代码
	byte1 := []byte("3go8&$8*3*3h0k(2)2")
	byte2 := []byte(strDfsID)
	byte1Len := len(byte1)
	for i := range byte2 {
		byte2[i] = byte2[i] ^ byte1[i%byte1Len]
	}
	sum := md5.Sum(byte2)
	var buff bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buff)
	_, err := enc.Write(sum[:])
	if nil != err {
		log.Printf("error encoding(base64) netease dfsId %s:%s", strDfsID, err)
		return ""
	}
	enc.Close()
	eid := []byte(gNeteaseEIDReplacer.Replace(buff.String()))
	// }
	// update cache, no expiration, no compression
	// SetCache(eidKey, eid, 0, false)
	return fmt.Sprintf(gNeteaseMusicCDNUrlF, eid, strDfsID)
}

func (ns *NeteaseSong) ArtistsString() string {
	arts := ""
	for i, _ := range ns.Artists {
		arts += (ns.Artists[i].Name + ",")
	}
	return strings.TrimRight(arts, ",")
}

func (ns *SearchSong) ArtistsString() string {
	arts := ""
	for i, _ := range ns.Artists {
		arts += (ns.Artists[i].Name + ",")
	}
	return strings.TrimRight(arts, ",")
}

func (ns *NeteaseSong) CoverString() string {
	return ns.Album.CoverUrl
}

type NeteaseArtist struct {
	Name string `json:"name,omitempty"`
}

func GetSearchList(songname string, limit string, offset string) ([]*db.Music, error) {
	songurl := "http://music.163.com/api/search/get/"
	ret := PostUrl(gNeteaseClient, songurl, songname, limit, offset)
	// fmt.Printf(string(ret))
	if nil == ret {
		return nil, errors.New("error accessing url")
	}

	var songlistRet SearchResultRet
	err := json.Unmarshal(ret, &songlistRet)

	if nil != err {
		return nil, errors.New("error parsing songlist return data from url")
	}
	// if gNeteaseRetOk != songlistRet.StatusCode {
	// 	return sl.SetAndLogErrorf("error getting url %s: %s", songurl, songlistRet.Message)
	// }
	var musics []*db.Music
	for i, _ := range songlistRet.Result.Songs {
		song := &songlistRet.Result.Songs[i]
		musics = append(musics, &db.Music{
			Name:     song.Name,
			Id:       song.Id,
			Artists:  song.ArtistsString(),
			Duration: song.Duration,
		})
	}
	return musics, nil
}

func GetNeteaseSongList(songId string) (*db.Music, error) {
	url := fmt.Sprintf(gNeteaseAPIUrlBase+gNeteaseSongListUrl, songId)

	ret := GetUrl(gNeteaseClient, url)
	if nil == ret {
		return nil, errors.New("error accessing url")
	}

	var songlistRet NeteaseSongListRet
	err := json.Unmarshal(ret, &songlistRet)
	if nil != err {
		return nil, errors.New("error parsing songlist return data from url")
	}
	if gNeteaseRetOk != songlistRet.StatusCode {
		return nil, errors.New("error getting data")
	}
	var musics []*db.Music
	for i, _ := range songlistRet.Songs {
		song := (&songlistRet.Songs[i]).UpdateUrl("low")
		musics = append(musics, &db.Music{
			Id:       song.Id,
			Name:     song.Name,
			Url:      song.Url,
			Artists:  song.ArtistsString(),
			CoverUrl: song.CoverString(),
			Duration: song.Duration,
		})
	}

	if len(musics) != 0 {
		return musics[0], nil
	}
	return nil, nil
}
