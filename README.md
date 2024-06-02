# whisper-transcription

## Anforderungen

### Openai Key erstellen
- Account erstellen bei openai: https://openai.com/
- hier hin gehen: https://platform.openai.com/docs/overview
- Navigieren Dashboard > API Key > create new secret key
- names definieren > create
- den key kopieren (wird für später benötigt)

#### Zahlungsmethode hinterlegen
- Zahnrad Symbol (settings) > Billing > Payment Method > Add payment method

### ffmpeg installieren (optional wenn das video zu audio konvertiert werden muss)
- Kommandozeile öffnen (in windows suche "Eingabeaufforderung" eingeben und öffnen. Ein schwarzes fenster)
- `winget install ffmpeg` eingeben enter drücken
- wenn aufgefordert, mit `Y` Nutzungsbedingungen bestätigen

### .env file erstellen
- ein leeres file erstellen und es umbenennen zu `.env` (im selben ordner wie das transkribier programm)
- `.env` file öffnen
- Diese Zeile einfügen mit dem vorher erstellten openai key: `OPENAI_API_KEY=<kopierten openai key>`

### Ausführen
- (im selben ordner wie das transkribier programm) neuen ordner erstellen mit namen: data
- im data ordner neuen ordner erstellen mit namen: 1-video
- videos in diesen ordner kopieren
- programm ausführen (doppelklick auf das programm)
