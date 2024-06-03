package helpers

import "github.com/martinlindhe/subtitles"

// Hier wird die Transkription mit der Zeitangabe zusammengefasst, sodass pro Zeitangabe weniger genau angezeigt wird
// parameter:
// subs: das Untertiel Objekt welche die Untertitel enthält
// lengthInSec: die Länge in Sekunden wie lange die Untertitel in einen Zeitblock zusammengefasst werden sollen
func ConcatSubs(subs subtitles.Subtitle, lengthInSec int) subtitles.Subtitle {
	var tmpCaption subtitles.Caption

	newSubs := subtitles.Subtitle{}
	startTime := subs.Captions[0].Start
	captionText := make([]string, 0)

	for i, caption := range subs.Captions {
		if caption.End.Sub(startTime).Seconds() < float64(lengthInSec) {
			captionText = append(captionText, caption.Text...)

		} else {
			tmpCaption = subtitles.Caption{
				Start: startTime,
				End:   caption.End,
				Text:  captionText,
			}

			newSubs.Captions = append(newSubs.Captions, tmpCaption)

			if i+1 < len(subs.Captions) {
				startTime = subs.Captions[i+1].Start
				captionText = make([]string, 0)
			}
			
		}
	}
	return newSubs
}
