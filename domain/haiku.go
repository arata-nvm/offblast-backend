package domain

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type Word struct {
	class       tokenizer.TokenClass
	surface     string
	typ         string
	subTyp1     string
	subTyp2     string
	subTyp3     string
	orginalForm string
	pron        string

	length int
}

func listWords(text string) ([]Word, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}

	tokens := t.Tokenize(text)

	words := []Word{}
	for _, token := range tokens {
		stat := token.Class
		surface := token.Surface
		features := token.Features()
		if len(features) != 9 {
			// fmt.Println(node)
			// TODO
			continue
		}

		typ := features[0]
		subTyp1 := features[1]
		subTyp2 := features[2]
		subTyp3 := features[3]
		orginalForm := features[6]
		pron := features[8]
		re := regexp.MustCompile("[ァィゥェォャュョ]")
		pron = re.ReplaceAllString(pron, "")
		length := len([]rune(pron))

		if stat == tokenizer.UNKNOWN {
			length = 0
		}

		if typ == "記号" {
			length = 0
		}

		words = append(words, Word{
			stat,
			surface,
			typ,
			subTyp1,
			subTyp2,
			subTyp3,
			orginalForm,
			pron,
			length,
		})
	}

	return words, nil
}

func canBeFirstWord(word Word) bool {
	cond1 := word.typ != "助詞" && word.typ != "助動詞" && word.typ != "記号"
	cond2 := word.subTyp1 != "接尾" && word.subTyp1 != "非自立"
	cond3 := word.subTyp1 != "自立" || (word.orginalForm != "する" && word.orginalForm != "できる")
	return cond1 && cond2 && cond3
}

func canBePart(word Word) bool {
	cond1 := word.class != tokenizer.UNKNOWN
	cond2 := !strings.ContainsAny(word.surface, "（）「」。…")
	cond3 := word.subTyp1 != "アルファベット"
	cond4 := word.orginalForm != "第" || word.subTyp1 != "数接続"
	return cond1 && cond2 && cond3 && cond4
}

func canBeLastWord(word Word) bool {
	cond1 := word.typ != "連体詞"
	re := regexp.MustCompile("(名詞接続|格助詞|係助詞|連体化|接続助詞|並立助詞|副詞化|数接続)")
	cond2 := !re.MatchString(word.subTyp1)
	cond3 := word.typ != "助動詞" || word.orginalForm != "だ"
	cond4 := word.typ != "名詞" || word.subTyp1 != "数"
	return cond1 && cond2 && cond3 && cond4
}

var PartsLen = [...]int{5, 7, 5}

func FindHaikus(text string) ([]string, error) {
	haikus := []string{}

	words, err := listWords(text)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(words); i++ {
		var parts bytes.Buffer
		partIndex := 0
		sumLength := 0
		for j := i; j < len(words); j++ {
			word := words[j]
			if sumLength == 0 && !canBeFirstWord(word) {
				break
			}
			if !canBePart(word) {
				break
			}

			parts.WriteString(word.surface)
			sumLength += word.length

			if sumLength > PartsLen[partIndex] {
				break
			}
			if sumLength == PartsLen[partIndex] {
				partIndex++
				if partIndex == len(PartsLen) {
					if !canBeLastWord(word) {
						break
					}
					haikus = append(haikus, parts.String())
					break
				}
				parts.WriteByte(' ')
				sumLength = 0
			}
		}
	}

	return unique(haikus), nil
}

func unique(haikus []string) []string {
	newHaikus := []string{}
	m := map[string]bool{}

	for _, haiku := range haikus {
		if !m[haiku] {
			m[haiku] = true
			newHaikus = append(newHaikus, haiku)
		}
	}

	return newHaikus
}
