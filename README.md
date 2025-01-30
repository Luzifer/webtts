# Luzifer / webtts

This project is a simple wrapper around the [Google Cloud Text-To-Speech](https://cloud.google.com/text-to-speech) and [Azure Text-To-Speech](https://azure.microsoft.com/en-us/services/cognitive-services/speech-service/) API to output OGG Vorbis Audio to be used with OBS overlays.

## Usage

### Google Cloud Text-To-Speech

- Create a project in the [Google Cloud Console](https://console.cloud.google.com/).
- Enable the [Text-to-Speech API](https://cloud.google.com/text-to-speech/docs/apis).
- Create credentials (Service Account Key) and download it as JSON.
- Set the environment variable `GOOGLE_APPLICATION_CREDENTIALS` to the path of the downloaded JSON file.

### Azure Text-To-Speech

- Create a Text-To-Speech resource in the [Azure Portal](https://portal.azure.com/).
- Navigate to the resource and find the "Keys" section. Copy one of the keys.
- Search for a voice in the [Speech Studio](https://speech.microsoft.com/portal)
- Set the environment variable `AZURE_SPEECH_RESOURCE_KEY` to the copied key and `AZURE_SPEECH_REGION` to the region where your resource is located.

### Request

```
GET /tts.ogg
  ?provider=google|azure
  &lang=en-US
  &text=The%20text%20to%20convert%20to%20speech
  &valid-to=<RFC3339-Timestamp>
  &voice=<name-of-the-voice>
  &signature=<HMAC-SHA256-Signature>
```

- The `signature` is an HMAC-SHA256 signature of the request parameters, using the secret key
  - It contains all parameters except for the signature itself
  - The parameters are sorted by name
  - The HMAC is generated over the parameters concatinated with `\n`: `param1=value1\nparam2=value2\n...`
  - The signature is a lower-case hex encoding of the HMAC
- The `valid-to` timestamp is a RFC3339 timestamp indicating when the request is valid

So for example these would be valid URLs for the key `topsecret`:

```
http://localhost:3000/tts.ogg
  ?lang=en-EN
  &provider=google
  &text=Hello%20there%2C%20general%20Kenobi%21
  &valid-to=2025-01-31T01%3A22%3A17.405263Z
  &voice=de-DE-Standard-G
  &signature=afb82dc41b444f9573d585094cf4a22a517853b307d98031c9a324f294db026e

http://localhost:3000/tts.ogg
  ?lang=en-EN
  &provider=azure
  &text=Hello%20there%2C%20general%20Kenobi%21
  &valid-to=2025-01-31T01%3A23%3A57.905524Z
  &voice=en-US-AvaMultilingualNeural
  &signature=ad4b15b78acd7d59a9d659b6b7a67dce2eee070a49634f9be66be3727eb1f5fc
```
