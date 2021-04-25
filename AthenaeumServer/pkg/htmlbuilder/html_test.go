package htmlbuilder_test

import (
	"testing"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/htmlbuilder"
)

func TestH1Build(t *testing.T) {
	//given
	expected := "<h1>Heading 1 &lt; 2</h1>"

	//when
	actual := htmlbuilder.H1("Heading 1 < 2")

	//then
	if expected != actual {
		t.Errorf("Did not build HTML heading 1 correctly\nExpected: %s\nWas : %s", expected, actual)
	}
}

func TestH2Build(t *testing.T) {
	//given
	expected := "<h2>Heading 2 &lt; 3</h2>"

	//when
	actual := htmlbuilder.H2("Heading 2 < 3")

	//then
	if expected != actual {
		t.Errorf("Did not build HTML heading 2 correctly\nExpected: %s\nWas : %s", expected, actual)
	}
}

func TestH4Build(t *testing.T) {
	//given
	expected := "<h4>Heading 4 &lt; 5</h4>"

	//when
	actual := htmlbuilder.H4("Heading 4 < 5")

	//then
	if expected != actual {
		t.Errorf("Did not build HTML heading 4 correctly\nExpected: %s\nWas : %s", expected, actual)
	}
}

func TestParagraphBuild(t *testing.T) {
	//given
	expected := "<p>Line1</p><p>This is line 2 &amp; 3 &lt; 4.</p><p>This is a further line.</p>"

	//when
	actual := htmlbuilder.P("Line1\nThis is line 2 & 3 < 4.\nThis is a further line.")

	//then
	if expected != actual {
		t.Errorf("Did not build HTML paragraph correctly\nExpected: %s\nWas : %s", expected, actual)
	}
}
