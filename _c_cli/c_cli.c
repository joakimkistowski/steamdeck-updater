 #include <stdio.h>

 #include "../flatpakintegration/flatpak_integration.h"


int main() {
  printf("Refreshing Update Cache\n");
  updateFlatpakRemotes();
  printf("Checking for Updates\n");
  FlatpakContainerList* containerList =  getUpdatableFlatpaks();
  int runtimesCount = 0;
  int applicationsCount = 0;
  for (int i = 0; i < containerList->size; i++ ) {
    if (containerList->packages[i].isApplication) {
      applicationsCount++;
    } else {
      runtimesCount++;
    }
  }
  printf("Updating %d applications and %d runtimes (%d total)\n", applicationsCount, runtimesCount, containerList->size);
  if (applicationsCount > 0) {
    printf("Applications to be updated:\n");
  }
  for (int i = 0; i < containerList->size; i++ ) {
    if (containerList->packages[i].isApplication) {
      printf("    %s @%s\n", containerList->packages[i].name, containerList->packages[i].origin);  
    }
  }
  printf("Running updates ...\n");
  updateEverything();
  printf("Done\n");
  freeFlatpakContainerList(containerList);
  return 0;
}


