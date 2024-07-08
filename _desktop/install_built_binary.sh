#!/usr/bin/env sh

if [ -f steamdeckupdater ]; then
    mkdir -p ~/.local/bin
    mkdir -p ~/.local/share/applications
    mkdir -p ~/.local/share/icons/hicolor/scalable/apps
    mkdir -p ~/.local/share/icons/hicolor/48x48/apps
    cp -f steamdeckupdater ~/.local/bin/
    cp -f _desktop/steamdeckupdater.desktop ~/.local/share/applications/
    cp -f _desktop/steamdeckupdater-logo.svg ~/.local/share/icons/hicolor/scalable/apps/steamdeckupdater.svg
    cp -f _desktop/steamdeckupdater-logo.png ~/.local/share/icons/hicolor/48x48/apps/steamdeckupdater.png
    update-desktop-database ~/.local/share/applications
    echo "Success! Have fun updating your Steam Deck!"
else
    echo "steamdeckupdater binary not found. Run this in the steamdeck-updater root directory after building."
fi
