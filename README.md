# whisper-transcription

## Beschreibung
Dieses Script sendet die Audio files an OpenAI Whisper. Whisper transkibiert die Texte mit genauer Zeitangabe. Jedoch sind die Texte dann nicht zwingend auf Schweizerdeutsch.
Dann wird das Resultat an ChatGPT gesendet, wo es dann so genau wie möglich ins Schweizerdeutsch übersetzt wird.

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

## Konfiguration config.env
### AUDIO_DIR
Pfad zum Ordner wo die Audio Dateien Liegen. Per default `data/audio` Audio Files dort hin kopieren.

### TRANSCRIPTION_DIR
Pfad zum Ordner wo die zwischenschritte der Transkription gespeichert werden sollen.

### OUTPUT_DIR
Pfad zum Ordner wo das Endresultat gespeichert werden sollte.

### TRANSCRIBE_PROMPT
Prompt welcher der Transkriptions AI Mitgegeben wird um das Resultat etwas zu steuern.

### CHAT_GPT_PROMPT
Anweisung an Chat GPT um die Text ins Schweizerdeutsch zu übersetzen

### OPENAI_API_KEY
Key um die OpenAI API zu nutzen. Oben beschrieben wie dieser Key generiert werden kann.
