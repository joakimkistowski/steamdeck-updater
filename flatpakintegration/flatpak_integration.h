#ifndef __FLATPAK_INTEGRATION__
#define __FLATPAK_INTEGRATION__


#ifndef PKG_STR_SIZE
#define PKG_STR_SIZE 128
#endif

typedef struct s_FlatpakContainer {
    char name[PKG_STR_SIZE];
    char origin[PKG_STR_SIZE];
    char currentVersion[PKG_STR_SIZE];
    short isApplication;
} FlatpakContainer;

typedef struct s_FlatpakContainerList {
    short isErrorGettingPackages;
    int size;
    FlatpakContainer* packages;
} FlatpakContainerList;

typedef struct s_GenericFlatpakOperationResult {
    short isReturnedWithError;
    char errorMessage[PKG_STR_SIZE];
} GenericFlatpakOperationResult;

FlatpakContainer* getFlatpakContainer(FlatpakContainerList* flatpakContainerList, int index);

void updateFlatpakRemotes();
FlatpakContainerList* getUpdatableFlatpaks();
GenericFlatpakOperationResult updateEverything();

void freeFlatpakContainerList(FlatpakContainerList* flatpakContainerList);

#endif