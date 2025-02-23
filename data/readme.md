downlad ffmpeg from their offical repo and run `./convert.sh` to convert your video to chromecast appriopriate format

To extract subtitles you would use something similar to commands below.
Keep in mind that spring_original.webm has no subtitles channels embeded into media file.

Get a list of ids embbeded subtitles

```bash
./ffprobe \
		-select_streams s \
		-show_entries stream=index:stream_tags=language,title, \
		-of csv=p=0 \
        ./spring_original.webm
```

than you would pass those ids into map `0:<index>`

```bash
./ffmpeg \
        ./spring_original.webm \
		-i inputFile \
		-map 0:1 \
		subtitles.vtt
```
