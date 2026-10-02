package main

import "testing"

func TestUrl(t *testing.T) {
	uris := []string{
		"https://example.com",
		"https://www.google.com",
		"https://github.com",
		"https://go.dev",
		"https://this-domain-definitely-does-not-exist-12345.com",
		"https://example.com/does-not-exist",
		"not-a-uri",
		"https://",
		"https://example .com",
	}

	expected := []bool{
		true,
		true,
		true,
		true,
		false,
		false,
		false,
		false,
		false,
	}

	for i := range uris {
		outcome := isUrlValid(uris[i])
		if outcome != expected[i] {
			t.Errorf("URI %s \n got %t want %t", uris[i], outcome, expected[i])
		}
	}
}
