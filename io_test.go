package moist

import (
	"log"
	"testing"
)

func TestIO_01(t *testing.T) {

	log.Println("IO Test")

	log.Println("current directory: " + CurrentDir())

	log.Println("files of current directory:")
	f := ListFile(CurrentDir())
	log.Println(f)
}
