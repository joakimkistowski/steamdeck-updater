# steamdeckupdater

A simple graphical application with gamepad input for updating all flatpaks on the Steam Deck. Is intended to be used in Steam Deck's game mode but should generally work on all linux systems when updating flatpaks in a user installation.

## How to build

Dependencies are Go, libflatpak dev tools and ebitengine dependencies as described [here](https://ebitengine.org/en/documents/install.html).

Assuming you have Go, GCC, and make installed:

1. Install the dependencies:

Fedora:
```shell
sudo dnf install flatpak-devel mesa-libGL-devel mesa-libGLES-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel alsa-lib-devel pkg-config
```
Opensuse:
```shell
sudo zypper install flatpak-devel Mesa-libGL-devel Mesa-libGLESv2-devel Mesa-libGLESv3-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel alsa-devel pkgconf-pkg-config
```
Ubuntu/Debian:
```shell
sudo apt install libflatpak-dev libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
```

2. Build everything:
```shell
go build .
```
