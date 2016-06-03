package vm

import (
	bslcdisk "github.com/cloudfoundry/bosh-softlayer-cpi/softlayer/disk"
	bslcstem "github.com/cloudfoundry/bosh-softlayer-cpi/softlayer/stemcell"

	sldatatypes "github.com/maximilien/softlayer-go/data_types"
)

type VMCloudProperties struct {
	VmNamePrefix             string                               `json:"vmNamePrefix,omitempty"`
	Domain                   string                               `json:"domain,omitempty"`
	StartCpus                int                                  `json:"startCpus,omitempty"`
	MaxMemory                int                                  `json:"maxMemory,omitempty"`
	Datacenter               sldatatypes.Datacenter               `json:"datacenter"`
	BlockDeviceTemplateGroup sldatatypes.BlockDeviceTemplateGroup `json:"blockDeviceTemplateGroup,omitempty"`
	SshKeys                  []sldatatypes.SshKey                 `json:"sshKeys,omitempty"`
	RootDiskSize             int                                  `json:"rootDiskSize,omitempty"`
	EphemeralDiskSize        int                                  `json:"ephemeralDiskSize,omitempty"`

	HourlyBillingFlag              bool                                       `json:"hourlyBillingFlag,omitempty"`
	LocalDiskFlag                  bool                                       `json:"localDiskFlag,omitempty"`
	DedicatedAccountHostOnlyFlag   bool                                       `json:"dedicatedAccountHostOnlyFlag,omitempty"`
	NetworkComponents              []sldatatypes.NetworkComponents            `json:"networkComponents,omitempty"`
	PrivateNetworkOnlyFlag         bool                                       `json:"privateNetworkOnlyFlag,omitempty"`
	PrimaryNetworkComponent        sldatatypes.PrimaryNetworkComponent        `json:"primaryNetworkComponent,omitempty"`
	PrimaryBackendNetworkComponent sldatatypes.PrimaryBackendNetworkComponent `json:"primaryBackendNetworkComponent,omitempty"`
	BlockDevices                   []sldatatypes.BlockDevice                  `json:"blockDevices,omitempty"`
	UserData                       []sldatatypes.UserData                     `json:"userData,omitempty"`
	PostInstallScriptUri           string                                     `json:"postInstallScriptUri,omitempty"`

	BoshIp string `json:"bosh_ip,omitempty"`

	Baremetal             bool   `json:"baremetal,omitempty"`
	BaremetalStemcell     string `json:"bm_stemcell,omitempty"`
	BaremetalNetbootImage string `json:"bm_netboot_image,omitempty"`
}

type AllowedHostCredential struct {
	Iqn      string `json:"iqn"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type VMMetadata map[string]interface{}

type VMCreator interface {
	Create(string, bslcstem.Stemcell, VMCloudProperties, Networks, Environment) (VM, error)
}

type Finder interface {
	Find(int) (VM, bool, error)
}

type VM interface {
	AttachDisk(bslcdisk.Disk) error

	ConfigureNetworks(Networks) error

	DetachDisk(bslcdisk.Disk) error
	Delete(agentId string) error

	GetDataCenterId() int
	GetPrimaryIP() string
	GetPrimaryBackendIP() string
	GetRootPassword() string
	GetFullyQualifiedDomainName() string

	ID() int

	Reboot() error
	ReloadOS(bslcstem.Stemcell) error

	SetMetadata(VMMetadata) error
	SetVcapPassword(string) error
}

type Environment map[string]interface{}

type Mount struct {
	PartitionPath string
	MountPoint    string
}

const etcIscsidConfTemplate = `# Generated by bosh-agent
node.startup = automatic
node.session.auth.authmethod = CHAP
node.session.auth.username = {{.Username}}
node.session.auth.password = {{.Password}}
discovery.sendtargets.auth.authmethod = CHAP
discovery.sendtargets.auth.username = {{.Username}}
discovery.sendtargets.auth.password = {{.Password}}
node.session.timeo.replacement_timeout = 120
node.conn[0].timeo.login_timeout = 15
node.conn[0].timeo.logout_timeout = 15
node.conn[0].timeo.noop_out_interval = 10
node.conn[0].timeo.noop_out_timeout = 15
node.session.iscsi.InitialR2T = No
node.session.iscsi.ImmediateData = Yes
node.session.iscsi.FirstBurstLength = 262144
node.session.iscsi.MaxBurstLength = 16776192
node.conn[0].iscsi.MaxRecvDataSegmentLength = 65536
`

const (
	SOFTLAYER_HARDWARE_LOG_TAG = "SoftLayerHardware"
	SOFTLAYER_VM_FINDER_LOG_TAG = "SoftLayerVMFinder"
	SOFTLAYER_VM_OS_RELOAD_TAG = "OSReload"
	SOFTLAYER_VM_LOG_TAG       = "SoftLayerVM"
	ROOT_USER_NAME             = "root"
)
