package flatpakintegration

// #cgo CFLAGS: -I /usr/include/flatpak/ -I /usr/include/glib-2.0/ -I /usr/lib64/glib-2.0/include/ -I /usr/lib/x86_64-linux-gnu/glib-2.0/include/
// #cgo LDFLAGS: -l flatpak -l glib-2.0 -l gobject-2.0
// #include "flatpak_integration.h"
import "C"

const PkgStrSize = (C.int)(C.PKG_STR_SIZE)

type SDUFlatpakContainer struct {
	Name           string
	CurrentVersion string
	Origin         string
}

type SDUFlatpakContainerList struct {
	IsErrorGettingPackages bool
	RuntimePackages        []SDUFlatpakContainer
	ApplicationPackages    []SDUFlatpakContainer
}

type SDUFlatpakUpgradeResult struct {
	IsErrorUpgrading bool
	Message          string
}

func UpdateFlatpakRemotes() {
	C.updateFlatpakRemotes()
}

func GetUpdateableFlatpaks() *SDUFlatpakContainerList {
	cPackageList := C.getUpdatableFlatpaks()
	flatpakContainerList := &SDUFlatpakContainerList{
		IsErrorGettingPackages: int(cPackageList.isErrorGettingPackages) != 0,
		RuntimePackages:        make([]SDUFlatpakContainer, 0),
		ApplicationPackages:    make([]SDUFlatpakContainer, 0),
	}
	for i := 0; i < int(cPackageList.size); i++ {
		cCurrentFlatpakContainer := C.getFlatpakContainer(cPackageList, (C.int)(i))
		flatpakContainer := SDUFlatpakContainer{
			Name:           C.GoString(&cCurrentFlatpakContainer.name[0]),
			CurrentVersion: C.GoString(&cCurrentFlatpakContainer.currentVersion[0]),
			Origin:         C.GoString(&cCurrentFlatpakContainer.origin[0]),
		}

		if cCurrentFlatpakContainer.isApplication != 0 {
			flatpakContainerList.ApplicationPackages = append(flatpakContainerList.ApplicationPackages, flatpakContainer)
		} else {
			flatpakContainerList.RuntimePackages = append(flatpakContainerList.RuntimePackages, flatpakContainer)
		}
	}
	C.freeFlatpakContainerList(cPackageList)
	return flatpakContainerList
}

func UpgradeAllFlatpaks() *SDUFlatpakUpgradeResult {
	cUpgradeResult := C.updateEverything()
	upgradeResult := &SDUFlatpakUpgradeResult{
		IsErrorUpgrading: int(cUpgradeResult.isReturnedWithError) != 0,
		Message:          C.GoString(&cUpgradeResult.errorMessage[0]),
	}
	return upgradeResult
}
