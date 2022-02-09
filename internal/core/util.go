package core

import (
	"encoding/json"
	"github.com/anaskhan96/soup"
	"github.com/jaeha-choi/DFF/internal/datatype"
	"time"
)

func (client *DFFClient) getFromJson(url string) (r *datatype.OPGGResponse, ok bool) {
	var err error

	soup.Cookie("customLocale", client.Language)

	var resp string
	retryCnt := 3
	for i := 0; i < retryCnt; i++ {
		resp, err = soup.Get(url)
		if err == nil {
			break
		} else if i == retryCnt-1 {
			client.Log.Debug(err)
			client.Log.Error("Couldn't connect to the given url", url)
			return nil, false
		}
		client.Log.Debug(err)
		time.Sleep(500 * time.Millisecond)
		client.Log.Debug("Retrying..")
	}

	doc := soup.HTMLParse(resp)
	doc = doc.Find("script", "id", "__NEXT_DATA__")

	if err = json.Unmarshal([]byte(doc.Text()), &r); err != nil {
		client.Log.Debug(err)
		return nil, ok
	}

	//if err = json.Unmarshal([]byte(doc.Text()), &jsonRoot); err != nil {
	//	client.Log.Debug(err)
	//	return nil, ok
	//}
	//
	//if out, ok = jsonRoot["props"]; !ok {
	//	client.Log.Debug("value 'props' does not exist")
	//	return nil, ok
	//}
	//if jsonRoot, ok = out.(map[string]interface{}); !ok {
	//	client.Log.Debug("value 'props' is not map[string]interface{}")
	//	return nil, ok
	//}
	//
	//if out, ok = jsonRoot["pageProps"]; !ok {
	//	client.Log.Debug("value 'pageProps' does not exist")
	//	return nil, ok
	//}
	//if jsonRoot, ok = out.(map[string]interface{}); !ok {
	//	client.Log.Debug("value 'pageProps' is not map[string]interface{}")
	//	return nil, ok
	//}
	//
	//if out, ok = jsonRoot[dataKey]; !ok {
	//	client.Log.Debug(dataKey, "does not exist")
	//	return nil, ok
	//}

	return r, true
}

func min(i int, j int) int {
	if i < j {
		return i
	}
	return j
}
