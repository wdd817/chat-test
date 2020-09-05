package wordfilter

import (
	"testing"
)

func TestIsMatch(t *testing.T) {
	sensitiveList := []string{"china", "china's"}
	input := "i'm from china chengdu"

	util := NewDFAUtil(sensitiveList)
	if util.IsMatch(input) == false {
		t.Errorf("Expected true, but got false")
	}
}

func TestHandleWord(t *testing.T) {
	sensitiveList := []string{"china", "china's"}
	input := "i'm from china chengdu"

	util := NewDFAUtil(sensitiveList)
	newInput := util.HandleWord(input, '*')
	expected := "i'm from ***** chengdu"
	if newInput != expected {
		t.Errorf("Expected %s, but got %s", expected, newInput)
	}
}
