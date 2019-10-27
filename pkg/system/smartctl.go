package system

import (
	"encoding/json"
	"github.com/Masterminds/semver"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"strconv"
)

type Smartctl struct {
	path   string
	server Server
	fields data.Fields
}

type TestType string

const TestLong TestType = "long"
const TestShort TestType = "short"

func NewSmartCtl(path string, server Server) (Smartctl, error) {
	if path == "" {
		return Smartctl{}, errs.With("Path cannot be empty")
	}
	return Smartctl{
		path:   path,
		server: server,
		fields: data.WithField("path", path),
	}, nil
}

func SmartctlVersion() (semver.Version, error) {
	cmd := `smartctl -j --version`
	output, err := runner.Local.ExecShellGetStdout(cmd)
	if err != nil {
		return semver.Version{}, errs.WithEF(err, data.WithField("cmd", cmd), "Failed to call smartctl version")
	}

	smartResult := SmartResult{}
	if err = json.Unmarshal([]byte(output), &smartResult); err != nil {
		return semver.Version{}, errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal smartctl response")
	}

	versionString := ""
	for i, v := range smartResult.Smartctl.Version {
		if i > 0 {
			versionString += "."
		}
		versionString += strconv.Itoa(v)
	}

	version, err := semver.NewVersion(versionString)
	if err != nil {
		return semver.Version{}, errs.WithEF(err, data.WithField("versionString", versionString), "Failed to parse smartctl version")
	}
	return *version, nil
}

func (s Smartctl) All() (SmartResult, error) {
	smartResult := SmartResult{}
	output, err := s.server.Exec("smartctl", "--all", "-j", s.path)
	if err != nil {
		return smartResult, errs.WithEF(err, s.fields, "Fail to run smartctl")
	}
	logs.WithField("output", string(output)).Trace("smart output")

	if err = json.Unmarshal([]byte(output), &smartResult); err != nil {
		return smartResult, errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal smartctl result")
	}
	return smartResult, nil
}

func (s Smartctl) RunTest(testType TestType) error {
	output, err := s.server.Exec("sudo smartctl -t  " + s.path + " || true")
	if err != nil {
		return errs.WithEF(err, s.fields, "Fail to run smartctl")
	}
	logs.WithField("output", string(output)).Trace("smart output")
	return nil
}

//////////////////////

type SmartResult struct {
	JSONFormatVersion []int `json:"json_format_version"`
	Smartctl          struct {
		Version      []int    `json:"version"`
		SvnRevision  string   `json:"svn_revision"`
		PlatformInfo string   `json:"platform_info"`
		BuildInfo    string   `json:"build_info"`
		Argv         []string `json:"argv"`
		ExitStatus   int      `json:"exit_status"`
	} `json:"smartctl"`
	Device struct {
		Name     string `json:"name"`
		InfoName string `json:"info_name"`
		Type     string `json:"type"`
		Protocol string `json:"protocol"`
	} `json:"device"`
	ModelFamily  string `json:"model_family"`
	ModelName    string `json:"model_name"`
	SerialNumber string `json:"serial_number"`
	Wwn          struct {
		Naa int `json:"naa"`
		Oui int `json:"oui"`
		ID  int `json:"id"`
	} `json:"wwn"`
	FirmwareVersion string `json:"firmware_version"`
	UserCapacity    struct {
		Blocks int   `json:"blocks"`
		Bytes  int64 `json:"bytes"`
	} `json:"user_capacity"`
	LogicalBlockSize  int `json:"logical_block_size"`
	PhysicalBlockSize int `json:"physical_block_size"`
	RotationRate      int `json:"rotation_rate"`
	FormFactor        struct {
		AtaValue int    `json:"ata_value"`
		Name     string `json:"name"`
	} `json:"form_factor"`
	InSmartctlDatabase bool `json:"in_smartctl_database"`
	AtaVersion         struct {
		String     string `json:"string"`
		MajorValue int    `json:"major_value"`
		MinorValue int    `json:"minor_value"`
	} `json:"ata_version"`
	SataVersion struct {
		String string `json:"string"`
		Value  int    `json:"value"`
	} `json:"sata_version"`
	InterfaceSpeed struct {
		Max struct {
			SataValue      int    `json:"sata_value"`
			String         string `json:"string"`
			UnitsPerSecond int    `json:"units_per_second"`
			BitsPerUnit    int    `json:"bits_per_unit"`
		} `json:"max"`
		Current struct {
			SataValue      int    `json:"sata_value"`
			String         string `json:"string"`
			UnitsPerSecond int    `json:"units_per_second"`
			BitsPerUnit    int    `json:"bits_per_unit"`
		} `json:"current"`
	} `json:"interface_speed"`
	LocalTime struct {
		TimeT   int    `json:"time_t"`
		Asctime string `json:"asctime"`
	} `json:"local_time"`
	ReadLookahead struct {
		Enabled bool `json:"enabled"`
	} `json:"read_lookahead"`
	WriteCache struct {
		Enabled bool `json:"enabled"`
	} `json:"write_cache"`
	AtaSecurity struct {
		State   int    `json:"state"`
		String  string `json:"string"`
		Enabled bool   `json:"enabled"`
		Frozen  bool   `json:"frozen"`
	} `json:"ata_security"`
	SmartStatus struct {
		Passed bool `json:"passed"`
	} `json:"smart_status"`
	AtaSmartData struct {
		OfflineDataCollection struct {
			Status struct {
				Value  int    `json:"value"`
				String string `json:"string"`
			} `json:"status"`
			CompletionSeconds int `json:"completion_seconds"`
		} `json:"offline_data_collection"`
		SelfTest struct {
			Status struct {
				Value  int    `json:"value"`
				String string `json:"string"`
				Passed bool   `json:"passed"`
			} `json:"status"`
			PollingMinutes struct {
				Short      int `json:"short"`
				Extended   int `json:"extended"`
				Conveyance int `json:"conveyance"`
			} `json:"polling_minutes"`
		} `json:"self_test"`
		Capabilities struct {
			Values                        []int `json:"values"`
			ExecOfflineImmediateSupported bool  `json:"exec_offline_immediate_supported"`
			OfflineIsAbortedUponNewCmd    bool  `json:"offline_is_aborted_upon_new_cmd"`
			OfflineSurfaceScanSupported   bool  `json:"offline_surface_scan_supported"`
			SelfTestsSupported            bool  `json:"self_tests_supported"`
			ConveyanceSelfTestSupported   bool  `json:"conveyance_self_test_supported"`
			SelectiveSelfTestSupported    bool  `json:"selective_self_test_supported"`
			AttributeAutosaveEnabled      bool  `json:"attribute_autosave_enabled"`
			ErrorLoggingSupported         bool  `json:"error_logging_supported"`
			GpLoggingSupported            bool  `json:"gp_logging_supported"`
		} `json:"capabilities"`
	} `json:"ata_smart_data"`
	AtaSctCapabilities struct {
		Value                         int  `json:"value"`
		ErrorRecoveryControlSupported bool `json:"error_recovery_control_supported"`
		FeatureControlSupported       bool `json:"feature_control_supported"`
		DataTableSupported            bool `json:"data_table_supported"`
	} `json:"ata_sct_capabilities"`
	AtaSmartAttributes struct {
		Revision int `json:"revision"`
		Table    []struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			Value      int    `json:"value"`
			Worst      int    `json:"worst"`
			Thresh     int    `json:"thresh"`
			WhenFailed string `json:"when_failed"`
			Flags      struct {
				Value         int    `json:"value"`
				String        string `json:"string"`
				Prefailure    bool   `json:"prefailure"`
				UpdatedOnline bool   `json:"updated_online"`
				Performance   bool   `json:"performance"`
				ErrorRate     bool   `json:"error_rate"`
				EventCount    bool   `json:"event_count"`
				AutoKeep      bool   `json:"auto_keep"`
			} `json:"flags"`
			Raw struct {
				Value  int    `json:"value"`
				String string `json:"string"`
			} `json:"raw"`
		} `json:"table"`
	} `json:"ata_smart_attributes"`
	PowerOnTime struct {
		Hours int `json:"hours"`
	} `json:"power_on_time"`
	PowerCycleCount int `json:"power_cycle_count"`
	Temperature     struct {
		Current                   int `json:"current"`
		PowerCycleMin             int `json:"power_cycle_min"`
		PowerCycleMax             int `json:"power_cycle_max"`
		LifetimeMin               int `json:"lifetime_min"`
		LifetimeMax               int `json:"lifetime_max"`
		OpLimitMax                int `json:"op_limit_max"`
		OpLimitMin                int `json:"op_limit_min"`
		LimitMin                  int `json:"limit_min"`
		LimitMax                  int `json:"limit_max"`
		LifetimeOverLimitMinutes  int `json:"lifetime_over_limit_minutes"`
		LifetimeUnderLimitMinutes int `json:"lifetime_under_limit_minutes"`
	} `json:"temperature"`
	AtaLogDirectory struct {
		GpDirVersion        int  `json:"gp_dir_version"`
		SmartDirVersion     int  `json:"smart_dir_version"`
		SmartDirMultiSector bool `json:"smart_dir_multi_sector"`
		Table               []struct {
			Address      int    `json:"address"`
			Name         string `json:"name"`
			Read         bool   `json:"read,omitempty"`
			Write        bool   `json:"write,omitempty"`
			GpSectors    int    `json:"gp_sectors"`
			SmartSectors int    `json:"smart_sectors,omitempty"`
		} `json:"table"`
	} `json:"ata_log_directory"`
	AtaSmartErrorLog struct {
		Extended struct {
			Revision int `json:"revision"`
			Sectors  int `json:"sectors"`
			Count    int `json:"count"`
		} `json:"extended"`
	} `json:"ata_smart_error_log"`
	AtaSmartSelfTestLog struct {
		Extended struct {
			Revision int `json:"revision"`
			Sectors  int `json:"sectors"`
			Count    int `json:"count"`
		} `json:"extended"`
	} `json:"ata_smart_self_test_log"`
	AtaSmartSelectiveSelfTestLog struct {
		Revision int `json:"revision"`
		Table    []struct {
			LbaMin int `json:"lba_min"`
			LbaMax int `json:"lba_max"`
			Status struct {
				Value  int    `json:"value"`
				String string `json:"string"`
			} `json:"status"`
		} `json:"table"`
		Flags struct {
			Value                int  `json:"value"`
			RemainderScanEnabled bool `json:"remainder_scan_enabled"`
		} `json:"flags"`
		PowerUpScanResumeMinutes int `json:"power_up_scan_resume_minutes"`
	} `json:"ata_smart_selective_self_test_log"`
	AtaSctStatus struct {
		FormatVersion int `json:"format_version"`
		SctVersion    int `json:"sct_version"`
		DeviceState   struct {
			Value  int    `json:"value"`
			String string `json:"string"`
		} `json:"device_state"`
		Temperature struct {
			Current         int `json:"current"`
			PowerCycleMin   int `json:"power_cycle_min"`
			PowerCycleMax   int `json:"power_cycle_max"`
			LifetimeMin     int `json:"lifetime_min"`
			LifetimeMax     int `json:"lifetime_max"`
			OpLimitMax      int `json:"op_limit_max"`
			UnderLimitCount int `json:"under_limit_count"`
			OverLimitCount  int `json:"over_limit_count"`
		} `json:"temperature"`
		VendorSpecific []int `json:"vendor_specific"`
	} `json:"ata_sct_status"`
	AtaSctTemperatureHistory struct {
		Version                int `json:"version"`
		SamplingPeriodMinutes  int `json:"sampling_period_minutes"`
		LoggingIntervalMinutes int `json:"logging_interval_minutes"`
		Temperature            struct {
			OpLimitMin int `json:"op_limit_min"`
			OpLimitMax int `json:"op_limit_max"`
			LimitMin   int `json:"limit_min"`
			LimitMax   int `json:"limit_max"`
		} `json:"temperature"`
		Size  int   `json:"size"`
		Index int   `json:"index"`
		Table []int `json:"table"`
	} `json:"ata_sct_temperature_history"`
	AtaSctErc struct {
		Read struct {
			Enabled bool `json:"enabled"`
		} `json:"read"`
		Write struct {
			Enabled bool `json:"enabled"`
		} `json:"write"`
	} `json:"ata_sct_erc"`
	AtaDeviceStatistics struct {
		Pages []struct {
			Number   int    `json:"number"`
			Name     string `json:"name"`
			Revision int    `json:"revision"`
			Table    []struct {
				Offset int    `json:"offset"`
				Name   string `json:"name"`
				Size   int    `json:"size"`
				Value  int    `json:"value"`
				Flags  struct {
					Value                 int    `json:"value"`
					String                string `json:"string"`
					Valid                 bool   `json:"valid"`
					Normalized            bool   `json:"normalized"`
					SupportsDsn           bool   `json:"supports_dsn"`
					MonitoredConditionMet bool   `json:"monitored_condition_met"`
				} `json:"flags"`
			} `json:"table"`
		} `json:"pages"`
	} `json:"ata_device_statistics"`
	SataPhyEventCounters struct {
		Table []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Size     int    `json:"size"`
			Value    int    `json:"value"`
			Overflow bool   `json:"overflow"`
		} `json:"table"`
		Reset bool `json:"reset"`
	} `json:"sata_phy_event_counters"`
}
