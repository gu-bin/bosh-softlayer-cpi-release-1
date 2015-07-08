package disk

import (
	"fmt"
	"strconv"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	sl "github.com/maximilien/softlayer-go/softlayer"
)

const softLayerCreatorLogTag = "SoftLayerCreator"

type SoftLayerCreator struct {
	softLayerClient sl.Client
	logger          boshlog.Logger
}

func NewSoftLayerDiskCreator(client sl.Client, logger boshlog.Logger) SoftLayerCreator {
	return SoftLayerCreator{
		softLayerClient: client,
		logger:          logger,
	}
}

func (c SoftLayerCreator) Create(size int, cloudProps DiskCloudProperties, virtualGuestId int) (Disk, error) {
	c.logger.Debug(softLayerCreatorLogTag, "Creating disk of size '%d'", size)

	vmService, err := c.softLayerClient.GetSoftLayer_Virtual_Guest_Service()
	if err != nil {
		return SoftLayerDisk{}, bosherr.WrapError(err, "Create SoftLayer Virtual Guest Service error.")
	}

	vm, err := vmService.GetObject(virtualGuestId)
	if err != nil || vm.Id == 0 {
		return SoftLayerDisk{}, bosherr.WrapError(err, fmt.Sprintf("Can not retrieve vitual guest with id: %d.", virtualGuestId))
	}

	storageService, err := c.softLayerClient.GetSoftLayer_Network_Storage_Service()
	if err != nil {
		return SoftLayerDisk{}, bosherr.WrapError(err, "Create SoftLayer Network Storage Service error.")
	}

	disk, err := storageService.CreateIscsiVolume(c.getSoftLayerDiskSize(size), strconv.Itoa(vm.Datacenter.Id))
	if err != nil {
		return SoftLayerDisk{}, bosherr.WrapError(err, "Create SoftLayer iSCSI disk error.")
	}

	return NewSoftLayerDisk(disk.Id, c.softLayerClient, c.logger), nil
}

func (c SoftLayerCreator) getSoftLayerDiskSize (size int) int {
	sizeArray := []int{20, 40, 80, 100, 250, 500, 1000, 2000, 4000, 8000, 12000}

	for i := range sizeArray {
		if ret := size / 1024; ret <= sizeArray[i] {
			return sizeArray[i]
		}
	}
	return 12000
}
