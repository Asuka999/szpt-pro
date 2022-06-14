package szpt

import (
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"szpt/utils"
)

type User struct {
	Account  string
	Password string
}

func SzptLogin(user User) *http.Client {
	Jar, _ := cookiejar.New(nil)
	//uri, _ := url.Parse("http://127.0.0.1:9090")
	Client := http.Client{
		Transport: &http.Transport{
			// set proxyman proxy
			//Proxy: http.ProxyURL(uri),
		},
		Jar: Jar,
	}
	resp, _ := Client.Get("https://authserver.szpt.edu.cn/authserver/login")
	Lt, pwdEncryptSalt := utils.GetEncry(resp)
	EncryptedPwd := utils.EncryPasswd(user.Password, pwdEncryptSalt)
	requestForm := strings.NewReader(url.Values{"username": {user.Account}, "password": {EncryptedPwd}, "lt": {Lt}, "dllt": {"userNamePasswordLogin"}, "execution": {"e1s1"}, "_eventId": {"submit"}, "rmShown": {"1"}}.Encode())
	req, _ := http.NewRequest("POST", "https://authserver.szpt.edu.cn/authserver/login", requestForm)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	Client.Do(req)
	Client.Get("https://ehall.szpt.edu.cn/publicappinternet/sys/szptpubxsjkxxbs/*default/index.do#/")
	menuinfoForm := "data=%7B%22APPID%22%3A%225812981499622390%22%2C%22APPNAME%22%3A%22szptpubxsjkxxbs%22%7D"
	menuinfoForms := strings.NewReader(menuinfoForm)
	Client.Post("https://ehall.szpt.edu.cn/publicappinternet/sys/itpub/MobileCommon/getMenuInfo.do", "application/x-www-form-urlencoded", menuinfoForms)
	return &Client
}

type ReportInfoJson struct {
	Code      int                    `json:"code"`
	Datas     map[string]interface{} `json:"datas"`
	DaySchool bool                   `json:"daySchool"`
}

func GetSaveReportInfo(client *http.Client) (reportInfoDatas map[string]interface{}) {
	url := "https://ehall.szpt.edu.cn/publicappinternet/sys/szptpubxsjkxxbs/mrxxbs/getSaveReportInfo.do"
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	var ReportInfoMap ReportInfoJson
	jsoniter.Unmarshal(body, &ReportInfoMap)
	reportInfoDatas = ReportInfoMap.Datas
	return reportInfoDatas
}

func SaveReportInfo(client *http.Client, eportInfoDatas map[string]interface{}) {
	pletReportMap := replenishReport(eportInfoDatas)
	pletReportJson, _ := jsoniter.Marshal(pletReportMap)
	body := url.Values{}
	body.Add("formData", string(pletReportJson))
	_body := body.Encode()
	bodyReader := strings.NewReader(_body)
	//file, _ := os.OpenFile("json.json", os.O_RDWR|os.O_CREATE, 0666)
	//file.Write([]byte(_body))
	url_ := "https://ehall.szpt.edu.cn/publicappinternet/sys/szptpubxsjkxxbs/mrxxbs/saveReportInfo.do"
	req, _ := http.NewRequest("POST", url_, bodyReader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client.Do(req)

}

func replenishReport(eportInfoDatas map[string]interface{}) (pletReportMap map[string]interface{}) {
	ReplenishJson := `
		{"WID":"",
		"ZRCXFHJZXCLX":"",
		"JRFXXCLX":"",
		"ZSDZ":"",
		"SXFS":"",
		"SFZZSXDWSS":"",
		"QYTZWTW":"",
		"QYTWSTW":"",
		"DTZSTW":"",
		"FHTJGJ":"",
		"QTXYSMDJWQK":"",
		"SSSQ":"",
		"XSQBDSJ":"",
		"JSJJGCJTSJ":"",
		"JSJTGCJTSJ":"",
		"JSJJJTGCYY":"",
		"STYCZK":"",
		"STYXZK":"",
		"HSJCBG":"",
		"XGYMJZJJ":"",
		"SFYYYXGYMJZ":"",
		"WJZXGYMYY":"",
		"YJZXGYMZJS":"",
		"HSJCJG":"1"
		}
	`

	var ReplenishJson_ map[string]interface{}
	jsoniter.UnmarshalFromString(ReplenishJson, &ReplenishJson_)
	pletReportMap = utils.MergeJSONMaps(eportInfoDatas, ReplenishJson_)
	return pletReportMap
}

func Report(user User) {
	Client := SzptLogin(user)
	GetSaveReportInfo(Client)
	eportInfoDatas := GetSaveReportInfo(Client)
	SaveReportInfo(Client, eportInfoDatas)
}
