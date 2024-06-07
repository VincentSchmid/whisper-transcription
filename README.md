# Whisper-Transkription

## Beschreibung
Dieses Script sendet die Audiodateien an OpenAI Whisper. Whisper transkribiert die Texte mit genauer Zeitangabe. Jedoch sind die Texte dann nicht zwingend auf Schweizerdeutsch. Anschliessend wird das Resultat an ChatGPT gesendet, um es so genau wie möglich ins Schweizerdeutsch zu übersetzen.

## Setup-Guide

### OpenAI Key erstellen
1. Account erstellen bei OpenAI: [OpenAI](https://openai.com/)
2. Navigieren zu: [OpenAI Platform](https://platform.openai.com/docs/overview)
3. Dashboard > API Key > Create new secret key
4. Name definieren > Create
5. Den Key kopieren (dieser wird später benötigt)

#### Zahlungsmethode hinterlegen
1. Zahnrad-Symbol (Settings) > Billing > Payment Method > Add payment method

### Applikation herunterladen
1. Download von: [Whisper Transcription Releases](https://github.com/VincentSchmid/whisper-transcription/releases/latest)
2. Unter Assets: macOS oder Windows Version auswählen und herunterladen (.zip)
3. Ordnerinhalt entpacken

### Konfiguration anpassen
1. `config.env` Datei mit einem Texteditor öffnen
2. Den OpenAI Key in dieser Zeile einfügen: `OPENAI_API_KEY=<kopierter openai key>` (auch die Klammern entfernen). Es sollte so aussehen: `OPENAI_API_KEY=sp-proj-...`
3. Restliche Konfiguration nach Belieben anpassen

### Ausführen
1. Audiodateien in den Ordner `data/audio` kopieren
2. Programm ausführen

#### macOS
1. Swissgerman-Transcription-Darwin ausführen
2. Bei einem Fehler: Einstellungen > Datenschutz & Sicherheit > Trotzdem öffnen klicken

#### Windows
1. Swissgerman-Transcription-Windows ausführen
2. Bei einem blauen Fenster: Weitere Informationen > Trotzdem ausführen

## Konfiguration `config.env`

| Variable                   | Beschreibung                                                                                 | Standardwert        |
|----------------------------|----------------------------------------------------------------------------------------------|---------------------|
| AUDIO_DIR                  | Pfad zum Ordner, wo die Audiodateien liegen. Audiodateien dort hin kopieren.                 | `data/audio`        |
| TRANSCRIPTION_DIR          | Pfad zum Ordner, wo die Zwischenschritte der Transkription gespeichert werden sollen.        | `data/transcription`|
| OUTPUT_DIR                 | Pfad zum Ordner, wo das Endergebnis gespeichert werden soll.                                 | `data/output`       |
| TRANSCRIBE_PROMPT          | Prompt, welcher der Transkriptions-KI mitgegeben wird, um das Resultat zu steuern.           |                     |
| CHAT_GPT_PROMPT            | Anweisung an ChatGPT, um den Text ins Schweizerdeutsch zu übersetzen.                        |                     |
| OPENAI_API_KEY             | Key, um die OpenAI API zu nutzen. Oben beschrieben, wie dieser Key generiert werden kann.    |                     |
| SUBTITLE_TIME_GRANULARITY  | Angabe in Sekunden, wie Granular die Zeitblocke sein müssen.                                 | 30 Sekunden         |
