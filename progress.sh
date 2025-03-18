# get argment name
name=$1
timestamp=$2

# if 5 file, remove 1 oldest
ls -1t | tail -n +5 | xargs rm -f

# save progress
sudo cp  data backup/backup-${name}-${timestamp}