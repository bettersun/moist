package moist

import (
	"log"
	"regexp"
	"testing"
)

func TestRegExp(t *testing.T) {

	commentMark := `//`
	sRegExp := `^[\s]*` + commentMark + `.*`

	regexpCode := regexp.MustCompile(sRegExp)

	v := `//aa`
	// v := `/*aa`
	// v := `aa*/`

	result := regexpCode.MatchString(v)
	log.Println(result)

	v = `   //aa`
	result = regexpCode.MatchString(v)
	log.Println(result)

	v = `   int main(void)`
	result = regexpCode.MatchString(v)
	log.Println(result)
}
