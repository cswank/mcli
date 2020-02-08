# Mcli
A command line flac music player.

## Configurations

### 1: All local
If the music files and speaker output all live on the computer:

```console
export MCLI_MUSIC_LOCATION="/path/to/flac/files"
mcli
```

NOTE: see [below](#flac-directory-layout) for the required layout of the flac files

### 2: Local speakers and remote flac files
If the music files live on a remote computer (for example: 192.1.0.22) and 
speaker output is on the local computer:

```console
export MCLI_HOST="192.1.0.22:50051"
mcli
```

NOTE: A mcli server must be running on 192.1.0.22:50051

### 3: Remote speakers and remote flac files
If the music files and speaker output live on a remote computer (for example: 192.1.0.22):

```console
export MCLI_HOST="192.1.0.22:50051"
mcli --remote
```

### Server
If a remote server is required (client configurations 2 and 3 above):

```console
export MCLI_MUSIC_LOCATION="/path/to/flac/files"
export MCLI_HOST="hostname or ip address of this machine"
export MCLI_HOME="/path/to/directory/where/mcli/history/will/be/written"
mcli serve
```

## Flac directory layout
The flac files must be orgainized like $MCLI_MUSIC_LOCATION/artist/album/song.flac

So for example, if $MCLI_MUSIC_LOCATION=/mnt/music:

```console
ls -d $PWD/01.Come\ Together.flac
'/mnt/music/The Beatles/Abbey Road (remix)/01.Come Together.flac'
```
