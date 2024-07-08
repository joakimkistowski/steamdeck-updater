 #include <stdio.h>
 //#include <packagekit-glib2/packagekit.h>

 //#include "pk_integration.h"
 #include "../flatpakintegration/flatpak_integration.h"

/*
int main() {
  printf("Refreshing Update Cache\n");
  GenericOperationResult updateCheckResult = checkForUpdates();
  if (updateCheckResult.isReturnedWithError) {
    printf("Error refreshing Cache: %s\n", updateCheckResult.errorMessage);
  }
  printf("Checking for Updates\n");
  PackageList* packageList = getUpdateablePackages();
  if (packageList == NULL || packageList->isErrorGettingPackages) {
    printf("Error getting Updates\n");
    return 1;
  }
  printf("%d updates available\n", packageList->size);
  //for (int i = 0; i < packageList->size; i++ ) {
  //  printf("%s -> %s\n", packageList->packages[i].name, packageList->packages[i].version);
  //}
  printf("Updating everything\n");
  GenericOperationResult updateEverythingResult = updateEverything();
  if(updateEverythingResult.isReturnedWithError) {
    printf("Error updating everything: %s\n", updateEverythingResult.errorMessage);
  }
  if (updateEverythingResult.isRestartRequired) {
    printf("Restart is required. Please restart your system.\n");
  }
  if(!updateEverythingResult.isReturnedWithError) {
    if (updateEverythingResult.errorMessage != NULL && strlen(updateEverythingResult.errorMessage) > 0) {
      printf("Successfully updated everything. System message: %s\n", updateEverythingResult.errorMessage);
    } else {
      printf("Successfully updated everything!\n");
    }
    
  }
  
  clearAndFreePackageIdsForUpdating();
  freePackageList(packageList);
  printf("Done ...\n");
  return 0;
}
*/

int main() {
  printf("Refreshing Update Cache\n");
  updateFlatpakRemotes();
  printf("Checking for Updates\n");
  FlatpakContainerList* containerList =  getUpdatableFlatpaks();
  printf("%d updates available\n", containerList->size);
  int runtimesCount = 0;
  int applicationsCount = 0;
  for (int i = 0; i < containerList->size; i++ ) {
    if (containerList->packages[i].isApplication) {
      applicationsCount++;
    } else {
      runtimesCount++;
    }
  }
  printf("Updating %d applications and %d runtimes\n", applicationsCount, runtimesCount);
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


