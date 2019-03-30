package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"strings"
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
	"L":  "l̩",
	"M":  "m̩",
	"N":  "n̩",
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
	"L":  "l̩",
	"M":  "m̩",
	"N":  "n̩",
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

var dictionary = make(map[string]string)

type Format int

const (
	Simplified Format = iota
	Ipa
)

var parser *ParsePronunciator

func Pronounce(format Format, text string) string {
	var str strings.Builder
	elements := parser.getElements(text)

	for _, element := range elements {
		switch element.elementType {
		case Symbol:
			str.WriteString(element.value)
		case Space:
			str.WriteString(" ")
		case World:
			arphabetSound, wasFound := getArphabetPhonetic(strings.ToLower(element.value))
			if !wasFound {
				str.WriteString(element.value)
				continue
			}
			word := arphabetTo(format, arphabetSound)
			str.WriteString(word)
		}
	}
	return str.String()
}


//receive format: ipa or simplified and some cmu string extracted from the dictionary, return the same world but
//displaying the ipa or simplified pronunciation
func arphabetTo(format Format, soundCMU string) string {
	requiredMap := simplifySounds
	if format == Ipa {
		requiredMap = arphabetToIPA
	}
	soundsCMU := strings.Split(soundCMU, " ")
	soundWords := make([]string, len(soundsCMU))

	notDigitRg := regexp.MustCompile(`\d`)
	for _, sound := range soundsCMU {
		sound = notDigitRg.ReplaceAllString(sound, "")

		if soundWord, ok := requiredMap[sound]; ok {
			soundWords = append(soundWords, soundWord)
		} else {
			soundWords = append(soundWords, sound)
		}

	}

	return strings.Join(soundWords, "")
}

// check the cmu dictionary wor the world and returns the cmu pronunciation
func getArphabetPhonetic(word string) (string, bool) {
	if val, ok := dictionary[word]; ok {
		return val, true
	}
	return word, false
}

func init(){
	parser,_ = NewParsePronunciator()
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
			dictionary[strings.ToLower(words[0])] = strings.TrimSpace(strings.Join(words[1:], " "))
		}
	}

	fmt.Println("imprimiendo hola")
	fmt.Println(dictionary["HELLO"])

	if err := scanner.Err(); err != nil {
		log.Fatal("error with scanner ", err)
	}
}

func main() {
	//dictionary := make(map[string]string)
	
	fmt.Println(Pronounce(Ipa, "what! are you doing today let me know now."))
}
