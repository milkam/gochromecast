./ffmpeg \
		-y   \
        `# our input` \
		-i spring_original.webm \
        `# try to optimize for speed` \
		-preset fast \
        `# encoding of audio` \
		-c:a aac \
        `# quality of audio` \
		-b:a 192k \
        `# number of audio channels` \
		-ac 2 \
        `# encoder settings` \
		-c:v libx264 \
        `# resolution settings` \
		-b:v 1024k \
         `# quality settings` \
		-profile:v high \
		`# higher number represents better quality while taking longer to rencode` \
        -crf 22 \
        `# pixel type,  (this is not resolution)` \
		-pix_fmt yuv420p \
		`# use all cpus` \
        -threads 0 \
        `# specify m3u8` \
		-f hls \
        `# intal m3u size` \
		-hls_list_size 0 \
        `# target chunks of approx 10 seconds, real will vary` \
		-hls_time 10 \
        `# allow for streaming` \
		-hls_playlist_type event \
        `# tell player that we will add to this list as it plays` \
		-hls_flags independent_segments+append_list \
		`# our output playlist` \
        -master_pl_name playlist.m3u8 \
		`# our output chunks location` \
        ./chunks/*