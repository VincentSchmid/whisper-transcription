package helpers

import (
	"fmt"
	"testing"

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
Ja, wenn mer vielleicht grad noch schnell auf das zurückkos, ich has noch echt eindrücklich gefunden, wenn sie so die Krebsarten vergleichet, an was stirbt mer hützetschad, an Krebs in de Schwiiz.

00:00:22.000 --> 00:00:32.000
Und es isch einfach noch eindrücklich, wenn sie zäge, mer hätt irgendwie da den Lungenkrebs und da hämmer die drei andere Arte und da chömme mer dann nachher noch uf Präventionsschwäche.

00:00:32.000 --> 00:00:40.000
Ja, sie händs vorig grad eklärt, also Lungenkrebs, das isch die weitaus dödlichste Krebsart und Hand in de Schwiiz.

00:00:40.000 --> 00:00:51.000
Das isch korrekt. Mer händ ja vier Krebsarten, wo über gleich hüufig sind in de Schwiiz. Das isch de Dickdarmkrebs, das isch die Brustkrebs bide Männer, die Brustkrebs bide Frauen und eben da Lungenkrebs.

00:00:51.000 --> 00:00:58.000
Das sind plus minus zwischen vier bis sechstausend Fälle pro Jahr die Verteilung von den Krebsarten.

00:00:58.000 --> 00:01:08.000
Wenn mer dann aber die Todesrate angeschaut, dann gseht mer, dass am Lungenkrebs mehr Leute sterben als an den drei Krebsarten.

00:01:08.000 --> 00:01:19.000
Und das zeigt ja, dass da ein gewisses Missverhältnis isch und wenn mer sich überlegt, warum stirbt mer vielleicht weniger an den anderen drei Krebsarten,

00:01:19.000 --> 00:01:24.000
dann händ sicher damit zu due, dass es da gewisse Früherkennungsmöglichkeiten gibt.

00:01:24.000 --> 00:01:37.000
Dickdarmkrebs kann mer eine Dickdarmspiegelung mache, Prostata kann mer beim Hausarzt untersuche, kann am Blutwert mache und beim Brustkrebs kann mer eine Ultraschall-MRI oder allenfalls eine Mammographie mache.

00:01:37.000 --> 00:01:43.000
Und ausgerechnet beim Lungenkrebs, wo am meisten Leute sterben, wird nichts gemacht.

00:01:43.000 --> 00:01:47.000
Und darum händ mer uns das auf die Fahne geschrieben und gseht, das kann nicht sein.

00:01:47.000 --> 00:01:55.000
Auf der einen Seite gibt es einen gewissen Zusammenhang, dass mer sagt, beim Lungenkrebs gibt es ja eigentlich eine Erklärung, warum mer das überchumt.

00:01:55.000 --> 00:02:00.000
Weil 90% der Lungenkrebspatienten sind Raucher oder X-Raucher.

00:02:00.000 --> 00:02:07.000
Und bei den anderen Krebsarten gibt es den direkten Zusammenhang mit der Lebensform nicht direkt.

00:02:07.000 --> 00:02:11.000
Und deshalb ist das ein bisschen vernachlässigt worden.

00:02:11.000 --> 00:02:13.000
Wir händ aber vor 20 Jahren...

00:02:13.000 --> 00:02:17.000
Darf ich schnell nachhaken? Sie sagen mit anderen Worten, weil die einfach selber schuld sind.

00:02:17.000 --> 00:02:18.000
Die sind selber schuld?

00:02:18.000 --> 00:02:20.000
Das sind die Raucher, denen schaue mer jetzt nicht.

00:02:20.000 --> 00:02:22.000
Genau, genau.

00:02:22.000 --> 00:02:25.000
Und das ist der Grund, wieso es dort kein...

00:02:25.000 --> 00:02:31.000
Nein, das ist nicht der Grund allein, aber das ist einer der Gründe, warum mer dort nicht so aktiv wurde.

00:02:31.000 --> 00:02:39.000
Weil mer gesagt hat, die Leute gehen ja bewusst, steht ja auf jedem Päckli, Rauchen tötet, bringt dich um.

00:02:39.000 --> 00:02:44.000
Die Leute gehen das Risiko ein, also wieso soll mer sich dort engagieren?

00:02:44.000 --> 00:02:49.000
Und wir haben dann 1999 angefangen mit der Früherkennung für Lungenkrebs.

00:02:49.000 --> 00:02:56.000
Haben damals eine Pilotstudie gemacht, sind aber eingebettet in eine grosse internationale Vereinigung.

00:02:56.000 --> 00:03:00.000
Die Early Detection for Lung Cancer heisst das.

00:03:00.000 --> 00:03:07.000
Das ist ursprünglich von der Cornell University, jetzt vom Mount Sinai Hospital in New York geleitet.

00:03:07.000 --> 00:03:13.000
Da sind unterdessen weltweit über 100'000 Patienten diagnostiziert worden.

00:03:13.000 --> 00:03:17.000
Wir haben ungefähr 1'500 Krebs insgesamt gefunden.

00:03:17.000 --> 00:03:23.000
Und haben jetzt diese Patienten über bis zu 20 Jahre verfolgen können, was passiert.

00:03:23.000 --> 00:03:32.000
Nach 10 Jahren nach der Früherkennung leben noch 90% der Lungenkrebspatienten, die früher erkannt wurden.

00:03:32.000 --> 00:03:40.000
Wenn wir nach dem herkömmlichen Prinzip, wie wir heute arbeiten, warten, bis der Patient kommt und irgendwelche Symptome hat,

00:03:40.000 --> 00:03:50.000
z.B. hat Blut im Auswurf oder hat irgendwelche Brustschmerzen, dann leben nach 5 Jahren noch höchstens 20% der Leute.

00:03:50.000 --> 00:03:58.000
Das zeigt, man kann bei diesem Krebs nicht warten, bis der Patient Symptome hat, bis man eine Diagnose stellt.

00:03:58.000 --> 00:04:04.000
Sondern man muss die Diagnose stellen zu einem Zeitpunkt, wo der Patient noch nichts davon merkt.

00:04:04.000 --> 00:04:08.000
Das Problem bei der Lunge ist, dass sie keine sensiblen Fasern hat.

00:04:08.000 --> 00:04:12.000
Es tut nicht weh. Die Lunge hat keine Schmerzfasern.

00:04:12.000 --> 00:04:15.000
Es tut nicht weh, wenn etwas wächst in der Lunge.

00:04:15.000 --> 00:04:18.000
Darum kommen die Patienten häufig viel zu spät.

00:04:18.000 --> 00:04:21.000
Zweitens sind viele dieser Patienten Raucher.

00:04:21.000 --> 00:04:25.000
Sie haben sowieso ein bisschen Husten und Auswurf.

00:04:25.000 --> 00:04:28.000
Man gewöhnt sich an das.

00:04:28.000 --> 00:04:36.000
Das ist nicht das primäre Symptom, bei dem der Patient sagt, dass er Husten und Auswurf hat und seine Untersuchungen lassen muss.

00:04:36.000 --> 00:04:40.000
Darum haben wir in der Früherkennung definierte Einschlusskriterien.

00:04:40.000 --> 00:04:47.000
Man hat das untersucht, ab welchem Alter, ab welcher Menge von Zigaretten das Risiko plötzlich ansteigt.

00:04:47.000 --> 00:04:55.000
Dann hat man gesehen, dass ab Alter 50, wenn man mindestens 20 Jahre 20 Zigaretten geraucht hat und nicht länger als 10 Jahre aufgehört hat,

00:04:55.000 --> 00:04:59.000
dann qualifiziert man für eine Früherkennung für Lungenkrebs.

00:04:59.000 --> 00:05:03.000
Herr Klinger, ich glaube, das war auch noch etwas umstritten.

00:05:03.000 --> 00:05:04.000
Nicht mehr.

00:05:04.000 --> 00:05:08.000
Ich glaube, in der Zwischenzeit ist aber etwas gegangen.

00:05:08.000 --> 00:05:16.000
Es war die Vereinigung der verschiedenen Ärzte, das Gremium, die eine definitive Empfehlung ausgesprochen hat.

00:05:16.000 --> 00:05:21.000
Wir haben uns immer an die internationalen Richtlinien gehalten.

00:05:21.000 --> 00:05:28.000
Wir haben das nicht selbst erfunden, sondern wir haben das in der grossen Gemeinschaft, die das erarbeitet hat,

00:05:28.000 --> 00:05:32.000
waren wir ein kleines Rädchen an diesem ZARA und haben dort mitgearbeitet.

00:05:32.000 --> 00:05:36.000
Das hat sich schon vor Jahren herauskristallisiert.

00:05:36.000 --> 00:05:43.000
Dass jetzt plötzlich die Universitäten das Gefühl haben, sie müssten es auch machen, hat natürlich mitgeteilt.

00:05:43.000 --> 00:05:45.000
Sie haben das Bild für uns.

00:05:45.000 --> 00:05:47.000
Das haben wir noch nicht gesehen.

00:05:47.000 --> 00:05:49.000
Entschuldigung.

00:05:49.000 --> 00:05:52.000
Können wir das noch schnell abtrennen?

00:05:52.000 --> 00:05:58.000
Wenn man das so eingibt, das Screening und so, dann sehen wir plötzlich die grossen Studien vom Unispital.

00:05:58.000 --> 00:06:01.000
Was ist das für eine Entwicklung?

00:06:01.000 --> 00:06:12.000
Ich sage, dass sie die richtige Literatur gelesen und aufmerksam verfolgt haben, dass das etwas Sinnvolles ist.

00:06:12.000 --> 00:06:15.000
Darum springen sie jetzt auch auf den Zug auf.

00:06:15.000 --> 00:06:18.000
Das ist nicht verkehrt.

00:06:18.000 --> 00:06:22.000
Man forscht natürlich weiter.

00:06:22.000 --> 00:06:24.000
Ich weiss gar nicht, in welche Richtung das geht.

00:06:24.000 --> 00:06:26.000
Aber wahrscheinlich muss man das noch bestätigen.

00:06:26.000 --> 00:06:28.000
Nein, nein.

00:06:28.000 --> 00:06:30.000
Was geht da eigentlich?

00:06:30.000 --> 00:06:34.000
Das ist unzählige Mal bestätigt worden, dass das sinnvoll ist.

00:06:34.000 --> 00:06:39.000
Aber es ist halt immer so, dass man lieber alles selber nochmals überprüft.

00:06:39.000 --> 00:06:43.000
Dass man dann zum Entscheid kommt, man sollte das eigentlich auch machen.

00:06:43.000 --> 00:06:50.000
In Tat und Wahrheit, das müssen Sie jetzt nicht aufzeichnen, ist es so, weil es von uns gekommen ist, wurde es abgelehnt.

00:06:50.000 --> 00:06:56.000
Weil es quasi aus dem Hirsland gekommen ist, haben die anderen gesagt, das taugt nichts.

00:06:56.000 --> 00:06:58.000
Das ist die Wahrheit.

00:06:58.000 --> 00:07:01.000
Was heisst abgelehnt worden? Einfach von diesen Fachleuten?

00:07:01.000 --> 00:07:03.000
Ja.

00:07:03.000 --> 00:07:09.000
Die grössere Frage in diesem Zusammenhang ist ja auch noch, was kostet, was bringt.

00:07:09.000 --> 00:07:17.000
Und dort ist die Situation immer noch so, dass es nicht wie bei den anderen Krebsarten, die Sie erwähnten, einfach gezahlt wird.

00:07:17.000 --> 00:07:27.000
Gut, aber da gibt es natürlich in den USA, wo ja die Kassen das zum Teil übernehmen, gibt es Kosten-Nutzen-Analysen.

00:07:27.000 --> 00:07:40.000
Wenn Sie in der Schweiz schauen, kostet die Behandlung eines Lungenkrebses heute ungefähr 40'000 Franken, wenn der Patient geheilt wird.

00:07:40.000 --> 00:07:48.000
Ein Patient, der ein vorgeschrittenes Tumorstadium hat, wie wir sie ja heute in 80% der Fälle erleben,

00:07:48.000 --> 00:07:52.000
dieser Patient bekommt nachher eine sogenannte multimodale Therapie,

00:07:52.000 --> 00:07:58.000
dann kommt vielleicht eine Chemotherapie über, kombiniert mit einer Immuntherapie, vielleicht sogar noch eine Bestrahlungstherapie.

00:07:58.000 --> 00:08:07.000
Diese Patienten geben 90% ihrer gesamten Gesundheitskosten, die sie im ganzen Leben haben, im letzten Lebensjahr aus.

00:08:07.000 --> 00:08:17.000
Und das sind ja dann nur Therapien, die den Krankheitsverlauf verzögern können, die aber keine Heilung mehr bringen.

00:08:17.000 --> 00:08:24.000
Darum ist es doch viel, viel gescheiter, man gibt das Geld am Anfang aus, tut die Leute frühzeitig trinken.
`

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

	assert.Equal(t, 88, len(condensedTranscript.Captions))
	assert.Equal(t, 1, len(condensedTranscript.Captions[10].Text))
}
