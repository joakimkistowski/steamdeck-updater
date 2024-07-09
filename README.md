# steamdeckupdater

A simple graphical application with gamepad input for updating all flatpaks on the Steam Deck. Is intended to be used in Steam Deck's game mode but should generally work on all linux systems when updating flatpaks, as long as the user running the update has the correct permissions.

## How to install

In Desktop mode: Download `steamdeckupdater.tar.gz` the latest release build [from our releases page](https://github.com/joakimkistowski/steamdeck-updater/releases). Extract it and then run `install.sh` either by double clicking and clicking `Execute` or by running `./install.sh` from your command line.

The application is now installed. You can add it to Steam on the Steam Deck by right clicking on it in the application launcher (start menu) and clicking `Add to Steam`.

### Help: The Application was not added to Steam correctly

For some reason adding the application to Steam does not always work and it will show up garbled. If this happens:

1. Right click the garbled application name (such as `steamdeckupdater.desktop`) in your Steam library (Desktop Mode)
2. Go to `Properties...`
3. Under *Shortcut*: Enter a nicer name, such as "Steam Deck Updater"
4. As *Target*: Enter `/home/deck/.local/bin/steamdeckupdater`
5. [**Optional**] For *Start In*: Enter `/home/deck`
6. [**Optional**] A nice logo can be found here: `/home/deck/.local/share/icons/hicolor/48x48/apps/steamdeckupdater.png` (or in the `desktop` directory within your extracted download)
7. [**Optional**] Images for your Steam Deck Grid can be found in the `desktop` directory within your extracted download


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
