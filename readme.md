
# HDM

HDM is a software that do Hard Drives Management.

It handle physical and logical disks lifecycle throughout servers.

HDM get: 
- `list` Info (location, serial, size, labels, form factor, ...) about plugged and unplugged disks across servers
- with a `scan` 
- Disks `health`
- `index` files location, size and date 
- `search` for files location, even on unplugged disks
- `verify`, `inspect`, `checkup`
- `index` disks file 
- Check directories are `backupable` 

HDM do:
- New disks `prepare` for usage (partition, encrypt, format, mount)
- Auto `add` plugged-in disks for usage (luksOpen, mount, restart)
- Handle disk `remove` actions (stop, kill, umount, luksClose )
- Auto cleanup removed disks
- Handle disks health `test`
- Pending blocks `repair`
- Trigger `backup`, even for disks of different size
- Run `backup`
- Directories `sync` across servers
- `restore`


HDM will notify for
- Notify for disks in bad state


- `check`
- run an `agent` that

HDM ask for human intervention to:
- Give password to format/open disks
- `destroy` Disk into too bad state to be trustable
- Plug a specific disk for backup, repair or restore 
- handle Directory not backupable due to unmatched source vs target size

## Requirements

- can ssh to servers with ssh agent
- can ssh from server to servers
- can run sudo on servers without password
- lsbk >= 2.33
- smartctl >= 7.0
- hdparm
- rsync

## Install

Same single binary file for agent, controller and cli

## Usage


## TODO

- find non backed-up path



Agent:
- inotify any file change: run backup, re-index
- trigger notification to operator: disk needs to be plugged, size of directory cannot be backuped, disk has failure
- 


sdo -> 325
sdk -> 155


switch 15W
freebox 22W
bbox -> 7W
