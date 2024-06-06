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

### Applikation herunterladen
- auf Release klicken rechts.
- mac oder windows version herunterladen (.zip)
- Ordnerinhalt entpacken

### Konfiguration anpassen
- `config.env` file öffnen
- Den OpenAI Key ersetzen auf dieser Zeile: `OPENAI_API_KEY=<kopierten openai key>` einfügen. Sollte so aussehen:`OPENAI_API_KEY=sp-proj-...`
- Restliche Konfiguration nach belieben anpassen

### Ausführen
- audio files in the data/audio ordner kopieren
- Programm ausführen
