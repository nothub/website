---
title: "guestfish"
description: "Manipulate qcow2 images with the guest filesystem shell"
date: "2024-06-02"
tags: ["qemu", "shell"]
draft: false
---

I recently found this handy tool called [guestfish](https://libguestfs.org/guestfish.1.html).
It is great for manipulating disk images in a precise manner without much workflow overhead.

```sh
sudo guestfish --rw --add "my-cool-image.qcow2" -i <<EOF
   write /root/.ssh/authorized_keys  "$(cat "~/.ssh/id_ed25519.pub")"
   upload local_file.txt             /etc/foo/bar.txt
EOF
```
