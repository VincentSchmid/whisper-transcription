package helpers

import (
	"fmt"
	"testing"
	"time"

	subtitles "github.com/martinlindhe/subtitles"
	"github.com/stretchr/testify/assert"
)

var shortFile = `WEBVTT

00:00:00.000 --> 00:00:08.000
Es war die beste Zeit, es war die schlimmste Zeit, es war das Zeitalter der Weisheit, es war das Zeitalter der Torheit, es war die Epoche des Glaubens, es war die Epoche des Unglaubens

00:00:08.000 --> 00:00:20.000
es war die Zeit des Lichts, es war die Zeit der Finsternis, es war der Frühling der Hoffnung, es war der Winter der Verzweiflung, wir hatten alles vor uns, wir hatten nichts vor uns.`

var longFile = `WEBVTT

00:00:00.000 --> 00:00:22.000
Es war die beste Zeit, es war die schlimmste Zeit, es war das Zeitalter der Weisheit, es war das Zeitalter der Torheit, es war die Epoche des Glaubens, es war die Epoche des Unglaubens

00:00:22.000 --> 00:00:32.000
es war die Zeit des Lichts, es war die Zeit der Finsternis, es war der Frühling der Hoffnung, es war der Winter der Verzweiflung, wir hatten alles vor uns, wir hatten nichts vor uns.

00:00:32.000 --> 00:00:40.000
Es war der Moment der Klarheit, es war der Moment der Verwirrung, es war die Ära des Fortschritts, es war die Ära des Rückschritts, es war die Zeit des Friedens, es war die Zeit des Krieges.`

func TestShorterThanConcatRange(t *testing.T) {
	subs, err := subtitles.NewFromVTT(shortFile)
	if err != nil {
		fmt.Println(err)
	}

	condensedTranscript := ConcatSubs(subs, 30)

	assert.Equal(t, 1, len(condensedTranscript.Captions))
	assert.Equal(t, 2, len(condensedTranscript.Captions[0].Text))
}

func TestExactConcatRange(t *testing.T) {
	subs, err := subtitles.NewFromVTT(shortFile)
	if err != nil {
		fmt.Println(err)
	}

	condensedTranscript := ConcatSubs(subs, 20)

	assert.Equal(t, 1, len(condensedTranscript.Captions))
	assert.Equal(t, 2, len(condensedTranscript.Captions[0].Text))
}

func TestGreaterThanConcatRange(t *testing.T) {
	subs, err := subtitles.NewFromVTT(shortFile)
	if err != nil {
		fmt.Println(err)
	}

	condensedTranscript := ConcatSubs(subs, 1)

	assert.Equal(t, 2, len(condensedTranscript.Captions))
	assert.Equal(t, 1, len(condensedTranscript.Captions[0].Text))
}

func TestLongFile(t *testing.T) {
	subs, err := subtitles.NewFromVTT(longFile)
	if err != nil {
		fmt.Println(err)
	}

	condensedTranscript := ConcatSubs(subs, 30)

	time32Seconds := condensedTranscript.Captions[0].Start.Add(32 * time.Second)

	assert.Equal(t, 2, len(condensedTranscript.Captions))
	assert.Equal(t, 2, len(condensedTranscript.Captions[0].Text))
	assert.Equal(t, 1, len(condensedTranscript.Captions[1].Text))
	assert.Equal(t, time32Seconds, condensedTranscript.Captions[1].Start)
}
