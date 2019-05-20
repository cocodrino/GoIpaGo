package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"strings"

	"github.com/gorilla/websocket"
	"net/http"

	"github.com/bregydoc/gtranslate"
)

var arphabetToIPA = map[string]string{
	"AA": "ɑ",
	"AE": "æ",
	"AH": "ʌ",
	"AO": "ɔ",
	"AW": "aʊ",
	"AY": "aɪ",
	"B":  "b",
	"CH": "tʃ",
	"D":  "d",
	"DH": "ð",
	"EH": "ɛ",
	"ER": "ɝ",
	"EY": "eɪ",
	"F":  "ɾ",
	"G":  "ɡ",
	"HH": "h",
	"IH": "ɪ",
	"IY": "i",
	"JH": "dʒ",
	"K":  "k",
	"L":  "l",
	"M":  "m",
	"N":  "n",
	"NG": "ŋ",
	"OW": "oʊ",
	"OY": "ɔɪ",
	"P":  "p",
	"R":  "ɹ",
	"S":  "s",
	"SH": "ʃ",
	"T":  "t",
	"TH": "θ",
	"UW": "u",
	"UH": "ʊ",
	"V":  "v",
	"W":  "w",
	"Y":  "j",
	"Z":  "z",
	"ZH": "ʒ",
}

var simplifySounds = map[string]string{
	"AA": "A",
	"AE": "a",
	"AH": "ʌ",
	"AO": "ao",
	"AW": "Au",
	"AY": "ai",
	"B":  "b",
	"CH": "ch",
	"D":  "d",
	"DH": "D",
	"EH": "E",
	"ER": "or",
	"EY": "ei",
	"F":  "f",
	"G":  "ɡ",
	"HH": "J",
	"IH": "i",
	"IY": "ii",
	"JH": "y",
	"K":  "k",
	"L":  "l",
	"M":  "m",
	"N":  "n",
	"NG": "n",
	"OW": "Ou",
	"OY": "Oy",
	"P":  "p",
	"R":  "r",
	"S":  "s",
	"SH": "ʃ",
	"T":  "t",
	"TH": "tt",
	"UW": "U",
	"UH": "Uu",
	"V":  "v",
	"W":  "w",
	"Y":  "y",
	"Z":  "z",
	"ZH": "zz",
}

type WebsocketResponse struct {
	Original   string `json:"original"`
	Ipa        string `json:"ipa"`
	Simplified string `json:"simplified"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
// en dictionary se van a cargar todas las palabras que contiene el archivo cmudict
var dictionary = make(map[string]string)

type Format int

const (
	Simplified Format = iota
	Ipa
)

var parser *ParsePronunciator



func init() {
	parser, _ = NewParsePronunciator()
	file, err := os.Open("cmudict-0.7b.txt")
	if err != nil {
		log.Fatal("error opening file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	regxValidLine := regexp.MustCompile(`^[A-Z]+\s`)
	for scanner.Scan() {
		text := scanner.Text()
		if !regxValidLine.MatchString(text) {
			continue
		}

		words := strings.Split(text, " ")
		if len(words) > 1 {
			//for instance "A" cmu sound is AH0
			//so we take the A as key for our dictionary and join the rest of words AH0 and other sounds
			// hello would ve dictionary["hello"] "HH AH0 L OW1"
			dictionary[strings.ToLower(words[0])] = strings.TrimSpace(strings.Join(words[1:], " "))
		}
	}

	//fmt.Println("imprimiendo hola")
	//fmt.Println(dictionary["hello"])

	if err := scanner.Err(); err != nil {
		log.Fatal("error with scanner ", err)
	}
}

//run go run $(ls -1 *.go | grep -v _test.go)
func main() {
	//dictionary := make(map[string]string)

	//fmt.Println(Pronounce(Ipa, "what! are you doing today let me know now."))
	http.HandleFunc("/translate", func(w http.ResponseWriter, request *http.Request) {
		conn, _ := upgrader.Upgrade(w, request, nil)

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			fmt.Println(string(msg))
			//resp := Pronounce(Ipa, string(msg))
			translatedText,err := gtranslate.TranslateWithFromTo(string(msg),gtranslate.FromTo{From:"auto",To:"en"})
			if err != nil{
				fmt.Println("error translating ",err)
			}

			jsonObj := &WebsocketResponse{Original: translatedText, Ipa: Pronounce(Ipa, translatedText), Simplified: Pronounce(Simplified, translatedText)}
			fmt.Println(jsonObj)
			resp, err := json.Marshal(jsonObj)
			fmt.Println("json", string(resp))

			if err != nil {
				fmt.Println("error json ", err)
			}

			if err = conn.WriteMessage(msgType, resp); err != nil {

			}
		}
	})

	fmt.Println("server running in port :8080, you can start to translating using the endpoint /translate")
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println("error running the server >\n", err)
	}

}
