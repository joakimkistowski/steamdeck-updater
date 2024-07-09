INCLUDES := -I /usr/include/flatpak/ -I /usr/include/glib-2.0/ -I /usr/lib64/glib-2.0/include/ -I /usr/lib/x86_64-linux-gnu/glib-2.0/include/
LIBRARIES := -l flatpak -l glib-2.0 -l gobject-2.0
OUTNAME := steamdeck-updater-cli
SOURCES := c_cli.c ../flatpakintegration/flatpak_integration.c

.PHONY: all
all: $(OUTNAME)

.PHONY: clean
clean:
	rm $(OUTNAME)

$(OUTNAME): $(SOURCES)
	$(CC) $(SOURCES) $(INCLUDES) -o $(OUTNAME)  $(LIBRARIES)