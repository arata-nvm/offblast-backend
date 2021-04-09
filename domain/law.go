package domain

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	"github.com/clbanning/mxj/v2"
	"github.com/labstack/gommon/log"
)

type Law struct {
	Info LawInfo
	Body string
}

type LawInfo struct {
	ID   string
	Name string
}

func GetRandomLaw() (*Law, error) {
	info, err := getRandomLawInfo()
	if err != nil {
		return nil, err
	}
	log.Infof("Getting: %s(%s)\n", info.ID, info.Name)

	rawBody, err := getRawBody(info)
	if err != nil {
		return nil, err
	}
	body, err := formatRawBody(rawBody)
	if err != nil {
		return nil, err
	}

	lawBody := Law{
		Info: *info,
		Body: body,
	}

	return &lawBody, nil
}

func formatRawBody(rawBody []string) (string, error) {
	var buf bytes.Buffer

	for _, raw := range rawBody {
		if !strings.HasSuffix(raw, "。") {
			continue
		}

		raw = strings.ReplaceAll(raw, "「", "")
		raw = strings.ReplaceAll(raw, "」", "")

		re := regexp.MustCompile("（[^（|^）]*）")
		raw = re.ReplaceAllString(raw, "")

		buf.WriteString(raw)
	}
	return buf.String(), nil
}

func getRawBody(info *LawInfo) ([]string, error) {
	url := "https://elaws.e-gov.go.jp/api/1/lawdata/" + info.ID

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	m, err := mxj.NewMapXml(body)
	if err != nil {
		return nil, err
	}

	v, err := m.ValuesForKey("Sentence")
	if err != nil {
		return nil, err
	}

	var rawBody []string
	var raw string
	for _, vv := range v {
		switch vv := vv.(type) {
		case map[string]interface{}:
			if text, ok := vv["#text"]; ok {
				raw = text.(string)
			} else {
				continue
			}
		case string:
			raw = vv
		}
		rawBody = append(rawBody, raw)
	}

	return rawBody, nil
}

func getRandomLawInfo() (*LawInfo, error) {
	body, err := ioutil.ReadFile("./static/lawlists.xml")
	if err != nil {
		return nil, err
	}

	m, err := mxj.NewMapXml(body)
	if err != nil {
		return nil, err
	}

	v, err := m.ValuesForPath("DataRoot.ApplData.LawNameListInfo")
	if err != nil {
		return nil, err
	}

	randomLaw := v[rand.Int31n(int32(len(v)))].(map[string]interface{})

	info := LawInfo{
		ID:   randomLaw["LawId"].(string),
		Name: randomLaw["LawName"].(string),
	}
	return &info, nil
}
