package trash



// state
cmd.AddCommand(
//command("list", []string{"ls"}, List,
//	"List known disks (even unplugged)"),
commandWithDiskSelector("index", []string{}, Index,
"Index files from disks"),

// cycle
cmd.AddCommand(
commandWithDiskSelector("add", []string{}, Add,
"AddBlockDevice disks as usable (mdadm,crypt,mount,restart)"),
commandWithRequiredDiskSelector("remove", []string{}, Remove,
"Remove or cleanup removed disk (kill,umount,restart,mdadm,crypt)"),
commandWithRequiredServerDiskAndLabel("erase", []string{}, Erase,
"securely erase disk"),
)

// heal

// backup
cmd.AddCommand(
commandWithRequiredDiskSelector("backupable", []string{}, Backupable,
"Find backup configs and run backups"),
commandWithRequiredDiskSelector("backup", []string{}, BackupCmd,
"Find backup configs and run backups"),
command("backups", []string{}, Backups,
"Find backup configs and run backups"),
)
