# GoChromecast

GoChromecast is a simple tool to stream your videos to your Android TV or google Chromecast device.

## Requirements

Since chromecast supports only certain formats I have prepared a small open source video [spring_original](<https://en.wikipedia.org/wiki/Spring_(2019_film)>) you can cast found in [data_dir](./data)

## Casting to tv

To use it run `go run ./cmd/mdns/mdns`

It will try to discover devices for 5s and output results.

In last log lines it will output device names that it finds on network.

```
2025/02/23 20:34:00 got devices on your network names:'[]string{"myAndroidName", "Iskon\\"}' url:'192.168.0.149:8009'
2025/02/23 20:34:00 got devices on your network names:'[]string{"myAndroidName2\\ tv", "TPM192E"}' url:'192.168.0.141:8009'
```

After that pass that device name to cast command

```go
    go run ./cmd/cast/cast.go -device myAndroidName
```

You should now see your movie spring_original beeing played on your tv together with subtitles.

## Future development

I might improve the chromecast thing adding features to it, however i don't plan on adding ffmpeg support at this time.

## Learning more about chromecast

I have written a blog series about chromecast diving deep into it.
You can find it on [my blog](https://vjerci.com/writings/chromecast/basics/)
