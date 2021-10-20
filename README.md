# Mcli
A command line flac music player.

<img src="./docs/mcli.gif"/>

## Configurations

### 1: All local
If the music files and speaker output all live on the computer:

```console
mcli -m /path/to/flac/files
```

NOTE: see [below](#flac-directory-layout) for the required layout of the flac files

### 2: Local speakers and remote flac files
If the music files live on a remote computer (for example: 192.1.0.22) and 
speaker output is on the local computer:

```console
mcli --host 192.1.0.22:50051
```

NOTE: A mcli server must be running on 192.1.0.22:50051

### 3: Remote speakers and remote flac files
If the music files and speaker output live on a remote computer (for example: 192.1.0.22):

```console
mcli --host 192.1.0.22:50051 --remote
```

### Server
If a remote server is required (client configurations 2 and 3 above):

```console
mcli serve -m /path/to/flac/files --host <dns name or ip address of this machine> --home /path/to/directory/where/mcli/database/lives
```

NOTE:  use --speakers=false if this server doesn't play music on its local machine

## Flac directory layout
The flac files must be orgainized like $MCLI_MUSIC_LOCATION/artist/album/song.flac

So for example, if $MCLI_MUSIC_LOCATION=/mnt/music:

```console
ls -d $PWD/01.Come\ Together.flac
'/mnt/music/The Beatles/Abbey Road (remix)/01.Come Together.flac'
```

## Colors
There are 3 different colors used in this app, and they can be customized by setting
the following environmental variables (the default values are shown here):

```console
export MCLI_C1="252"
export MCLI_C2="2"
export MCLI_C3="11"
```

The colors must be a number from the xterm 256 color palatte.
