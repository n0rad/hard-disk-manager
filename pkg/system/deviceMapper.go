package system


//for file in $(ls -la /dev/mapper/* | grep "\->" | grep -oP "\-> .+" | grep -oP " .+"); do echo "MAPPER:"$(F=$(echo $file | grep -oP "[a-z0-9-]+");echo $F":"$(ls "/sys/block/${F}/slaves/");)":"$(df -h "/dev/mapper/${file}" | sed 1d); done;
// dmsetup table
// dmsetup remove /dev/mapper/yopla
func findFromBlockDevice(name string) {

}