#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/root-x86"

cleanup() {
  sudo umount tmp/mount 2>/dev/null || true
}
trap cleanup EXIT

mkdir -p tmp
cat blk*.bin > tmp/rootfs.img

mkdir -p tmp/mount
sudo mount -o loop tmp/rootfs.img tmp/mount

sudo cp -a includes/. tmp/mount

sudo cp -a includes/. tmp/mount/
sudo cp -a ../../data tmp/mount/root/
sudo cp -a ../../posts tmp/mount/root/

sync
sudo umount tmp/mount
sudo rm -rf tmp/mount

rm -f blk.txt blk*.bin
split -b 256K -d -a 9 tmp/rootfs.img blk
rm -rf tmp

i=0
for f in blk*; do
  mv "${f}" "blk$(printf '%09d' ${i}).bin"
  i=$((i+1))
done

for f in blk*.bin; do
  truncate -s $((256*1024)) "${f}"
done

n=$(ls blk*.bin | wc -l)
cat > blk.txt <<EOF
{
  block_size: 256,
  n_block: $n,
}
EOF
