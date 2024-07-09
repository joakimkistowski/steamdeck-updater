# include "flatpak_integration.h"

#include <flatpak.h>
#include <stdio.h>
#include <glib.h>
#include <stdlib.h>
#include <pwd.h>


static FlatpakInstallation* flatpakUserInstallation = NULL;
static GPtrArray* flatpakInstallations = NULL;

static GError *error = NULL;

void clearErrorIfPresent() {
    if (error != NULL) {
        g_clear_error(&error);
        error = NULL;
    }
}

void printErrorIfPresent() {
    if (error != NULL) {
        printf("got an error: %s\n", error->message);
    }
}

void printAndClearErrorIfPresent() {
    printErrorIfPresent();
    clearErrorIfPresent();
}

GPtrArray* getFlatpakInstallations() {
    if (flatpakInstallations == NULL) {
        flatpakInstallations = g_ptr_array_new();
        g_ptr_array_extend_and_steal(flatpakInstallations, flatpak_get_system_installations(NULL, &error));
        printErrorIfPresent();
        if (flatpakUserInstallation == NULL) {
            flatpakUserInstallation = flatpak_installation_new_user(NULL, &error);
            printErrorIfPresent();
            clearErrorIfPresent();
            g_ptr_array_add(flatpakInstallations, flatpakUserInstallation);
        }
        return flatpakInstallations;
    }

    return flatpakInstallations;
}

FlatpakContainer* getFlatpakContainer(FlatpakContainerList *flatpakContainerList, int index) {
    if (flatpakContainerList == NULL || index < 0 || index >= flatpakContainerList->size) {
        return NULL;
    }
    return &(flatpakContainerList->packages[index]);
}

void updateFlatpakRemotesForInstallation(FlatpakInstallation* flatpakInstallation) {
    if (flatpakInstallation == NULL) {
        return;
    }
    GPtrArray* remotes = flatpak_installation_list_remotes(flatpakInstallation, NULL, &error);
    printAndClearErrorIfPresent();
    if (remotes == NULL) {
        return;
    }
    for (int i = 0; i < remotes->len; i++) {
        FlatpakRemote* remote = (FlatpakRemote*) remotes->pdata[i];
        const char* remoteName = flatpak_remote_get_name(remote);
        flatpak_installation_update_remote_sync(flatpakInstallation, remoteName, NULL, &error);
        printAndClearErrorIfPresent();
    }
    g_ptr_array_free(remotes, TRUE);
}

void updateFlatpakRemotes() {
    GPtrArray* installations = getFlatpakInstallations();
    if (installations == NULL || installations->len <= 0) {
        return;
    }
    for (int i = 0; i < installations->len; i++) {
        updateFlatpakRemotesForInstallation((FlatpakInstallation*) installations->pdata[i]);
    }
}

FlatpakContainerList* createMergedFlatpakContainerList(FlatpakContainerList* list1, FlatpakContainerList* list2) {
    if (list1 == NULL && list2 == NULL) {
        return NULL;
    }
    if (list1 == NULL) {
        return list2;
    }
    if (list2 == NULL) {
        return list1;
    }
    FlatpakContainerList* toReturn = malloc(sizeof(FlatpakContainerList));
    toReturn->isErrorGettingPackages = list1->isErrorGettingPackages || list2->isErrorGettingPackages;
    toReturn->size = list1->size + list2->size;
    toReturn->packages = malloc(toReturn->size * sizeof(FlatpakContainer));
    for (int i = 0; i < list1->size; i++) {
        toReturn->packages[i] = list1->packages[i];
    }
    for (int i = 0; i < list2->size; i++) {
        toReturn->packages[i + list1->size] = list2->packages[i];
    }
    return toReturn;
}

FlatpakContainerList* getUpdatableFlatpaksForInstallation(FlatpakInstallation* flatpakInstallation) {
    FlatpakContainerList* toReturn = malloc(sizeof(FlatpakContainerList));
    toReturn->isErrorGettingPackages = TRUE;
    toReturn->size = 0;
    toReturn->packages = NULL;
    if (flatpakInstallation == NULL) {
        return toReturn;
    }
    GPtrArray* updatableFlatpaks = flatpak_installation_list_installed_refs_for_update(flatpakInstallation, NULL, &error);
    printAndClearErrorIfPresent();
    if (updatableFlatpaks == NULL) {
        return toReturn;
    }
    toReturn->isErrorGettingPackages = FALSE;
    if (updatableFlatpaks->len <= 0) {
        g_ptr_array_free(updatableFlatpaks, TRUE);
        return toReturn;
    }
    toReturn->size = updatableFlatpaks->len;
    toReturn->packages = malloc(toReturn->size * sizeof(FlatpakContainer));
    for (int i = 0; i < updatableFlatpaks->len; i++) {
        FlatpakInstalledRef* fref = (FlatpakInstalledRef*) updatableFlatpaks->pdata[i];
        FlatpakContainer flatpakPackage;
        snprintf(flatpakPackage.name, PKG_STR_SIZE, "%s", flatpak_installed_ref_get_appdata_name(fref));
        snprintf(flatpakPackage.origin, PKG_STR_SIZE, "%s", flatpak_installed_ref_get_origin(fref));
        snprintf(flatpakPackage.currentVersion, PKG_STR_SIZE, "%s", flatpak_installed_ref_get_appdata_version(fref));
        flatpakPackage.isApplication = !flatpak_ref_get_kind(FLATPAK_REF(fref));
        toReturn->packages[i] = flatpakPackage;
    }
    g_ptr_array_free(updatableFlatpaks, TRUE);
    return toReturn;
}

FlatpakContainerList* getUpdatableFlatpaks() {
    FlatpakContainerList* toReturn = malloc(sizeof(FlatpakContainerList));
    toReturn->isErrorGettingPackages = FALSE;
    toReturn->size = 0;
    GPtrArray* installations = getFlatpakInstallations();
    if (installations == NULL || installations->len <= 0) {
        return toReturn;
    }
    for (int i = 0; i < installations->len; i++) {
        FlatpakContainerList* updatableFlatpaks = getUpdatableFlatpaksForInstallation((FlatpakInstallation*) installations->pdata[i]);
        FlatpakContainerList* mergedList = createMergedFlatpakContainerList(toReturn, updatableFlatpaks);
        free(toReturn);
        free(updatableFlatpaks);
        toReturn = mergedList;
    }
    return toReturn;
}

GenericFlatpakOperationResult updateEverythingForInstallation(FlatpakInstallation* flatpakInstallation) {
    FlatpakTransaction* flatpakTransaction = flatpak_transaction_new_for_installation(flatpakInstallation, NULL, &error);
    printAndClearErrorIfPresent();
    if (flatpakTransaction == NULL) {
        GenericFlatpakOperationResult res = {TRUE, "Unable to get create Flatpak transaction"};
        return res;
    }
    GPtrArray* updatableFlatpaks = flatpak_installation_list_installed_refs_for_update(flatpakInstallation, NULL, &error);
    printAndClearErrorIfPresent();
    if (updatableFlatpaks == NULL) {
        GenericFlatpakOperationResult res = {TRUE, "Unable to get list of installed and updateable Flatpaks"};
        return res;
    }
    GPtrArray * formattedRefs = g_ptr_array_new();
    for (int i = 0; i < updatableFlatpaks->len; i++) {
        FlatpakInstalledRef* fref = (FlatpakInstalledRef*) updatableFlatpaks->pdata[i];
        FlatpakRef* installedAsRef = FLATPAK_REF(fref);
        const char * formattedRef = flatpak_ref_format_ref(installedAsRef);
        g_ptr_array_add(formattedRefs, GINT_TO_POINTER(formattedRef));
        flatpak_transaction_add_update(flatpakTransaction, formattedRef, NULL, NULL, &error);
    }
    gboolean success = flatpak_transaction_run(flatpakTransaction, NULL, &error);
    g_ptr_array_free(formattedRefs, TRUE);
    
    
    if (!success) {
        GenericFlatpakOperationResult res = {TRUE, "Unknown error running flatpak upgrade"};
        if (error != NULL && error->message != NULL) {
            snprintf(res.errorMessage, PKG_STR_SIZE, "%s", error->message);
        }
        return res;
    }
    printAndClearErrorIfPresent();
    GenericFlatpakOperationResult res = {FALSE, ""};
    return res;
}

GenericFlatpakOperationResult updateEverything() {
    GPtrArray* installations = getFlatpakInstallations();
    if (installations == NULL || installations->len <= 0) {
        GenericFlatpakOperationResult res = {TRUE, "No Flatpak installations found"};
        return res;
    }
    for (int i = 0; i < installations->len; i++) {
        GenericFlatpakOperationResult res = updateEverythingForInstallation((FlatpakInstallation*) installations->pdata[i]);
        if (res.isReturnedWithError) {
            return res;
        }
    }
    GenericFlatpakOperationResult res = {FALSE, ""};
    return res;
}

void freeFlatpakContainerList(FlatpakContainerList* flatpakContainerList) {
    if (flatpakContainerList != NULL && flatpakContainerList->size > 0) {
        free(flatpakContainerList->packages);
    }
    free(flatpakContainerList);
}