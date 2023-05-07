package main

import (
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	type Params struct {
		Params string
	}

	r.POST("/checkParams", func(c *gin.Context) {
		jsonData, err := c.GetRawData()

		var params Params

		var secret string = "qIzWg28UmfSacBGRRvEZ" /*Защищенный ключ из настроек приложения*/

		var sign string

		json.Unmarshal([]byte(jsonData), &params)

		m, _ := url.ParseQuery(params.Params)

		vk_ := make(map[string]string)
		vk_Sorted := make(map[string]string)

		for k, v := range m {
			/*перебираем массив параметров по значению ключа*/
			if strings.HasPrefix(k, "vk_") {
				vk_[k] = v[0]
			}
			if strings.HasPrefix(k, "sign") {
				sign = v[0]
			}
		}

		keys := make([]string, 0, len(vk_))

		for k := range vk_ {
			keys = append(keys, k)
		}
		/*сортируем ключи с vk_ по алфавиту*/
		sort.Strings(keys)

		for _, k := range keys {
			for a, b := range vk_ {
				if a == k {
					/*добавляем ключи и их значения в отсортированный массив*/
					vk_Sorted[k] = b
				}
			}

		}

		sortedQuery := url.Values{}

		/*добавляем ключ=значение в URL query*/
		for k, v := range vk_Sorted {
			sortedQuery.Add(k, v)
		}
		/*кодируем URL query в строку*/
		encodedSortedQuery := sortedQuery.Encode()

		/*вычисляем SHA256 хеш*/
		h := hmac.New(sha256.New, []byte(secret))
		h.Write([]byte(encodedSortedQuery))
		sEnc := b64.StdEncoding.EncodeToString([]byte(h.Sum(nil)))

		/*заменяем в итоговой строке символы*/
		encodedReplacedPlus := strings.Replace(sEnc, "+", "-", 10)
		encodedReplacedSlash := strings.Replace(encodedReplacedPlus, "/", "_", 10)
		encodedReplacedSim := strings.Replace(encodedReplacedSlash, "=", "", 10)

		isSignValid := encodedReplacedSim == sign

		/*подготавливаем интерфейс для JSON ответа*/
		respMap := map[string]interface{}{
			"isSignValid": isSignValid,
		}

		/*конвертируем в JSON*/
		jsonDataResponse, err := json.Marshal(respMap)

		if err != nil {
			//Тут можно положить обработку ошибок
		}

		c.Data(http.StatusOK, "text/plain; charset=utf-8", jsonDataResponse)

	})
	r.Run()
}
