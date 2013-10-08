package utils

import (
	"testing"
)

type ExpectedResult struct {
	title string
	slug  string
}

var expectedResults = []ExpectedResult{
	ExpectedResult{
		title: "Fix low wifi speed on Linux (Ubuntu) with chip Atheros AR9285",
		slug:  "fix-low-wifi-speed-on-linux-ubuntu-with-chip-atheros-ar9285",
	},
	ExpectedResult{
		title: "Python and Scala smoke the peace pipe",
		slug:  "python-and-scala-smoke-the-peace-pipe",
	},
	ExpectedResult{
		title: "Graphite, Carbon and Diamond",
		slug:  "graphite-carbon-and-diamond",
	},
	ExpectedResult{
		title: "Here I go PyGrunn'13!",
		slug:  "here-i-go-pygrunn13",
	},
	ExpectedResult{
		title: "How-to install GNOME 3 instead of Unity on Ubuntu 11.04",
		slug:  "how-to-install-gnome-3-instead-of-unity-on-ubuntu-1104",
	},
	ExpectedResult{
		title: "This - is - just - a - test",
		slug:  "this-is-just-a-test",
	},
}

func TestSlugify(t *testing.T) {
	for _, testExpected := range expectedResults {
		title := testExpected.title
		expectedSlug := testExpected.slug

		if result := Slugify(title); result != expectedSlug {
			t.Errorf("%s\n+++ %s\n--- %s", title, result, expectedSlug)
		}
	}
}
