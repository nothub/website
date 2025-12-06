---
title: "OCI by hand"
description: "How do these container images even work?"
date: "2025-12-05"
tags: [ "shell" ]
draft: true
---

Did you ever ask yourself how these container images even work?

Yes? No? Never mind, let's build one by hand anyway.

---

First, we create a filesystem and populate it with our stuff:

```
image_dir="image"
rootfs="${image_dir}/rootfs"
app='f0VMRgEBAfCfkYsKAAAAAAIAAwABAAAAgIAECDQAAAAAuAQAAADNgOtYIAACACgABQAEAAEAAAAAAAAAAIAECACABAiiAAAAogAAAAUAAAAAEAAAAQAAAKQAAACkkAQIpJAECAkAAAAJAAAAugkAAAC5B5AECLsBAAAA66QAAADr6rsAAAAAuAEAAADNgA=='
mkdir -p "${rootfs}/bin"
echo "${app}" | base64 -d > "${rootfs}/bin/hi"
chmod +x "${rootfs}/bin/hi"
```

Then wrap it as a tarball:

```
cd "${rootfs}"
tar -cf ../layer.tar .
cd ..
```

And grab some meta-infos:

```
layer_digest="$(sha256sum "layer.tar" | awk '{print $1}')"
layer_size="$(wc -c < "layer.tar")"
```

To comply with the OCI image spec, we also need to create a [config.json](https://github.com/opencontainers/image-spec/blob/main/config.md):

```
timestamp="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
cat > "config.json" << EOF
{
  "created": "${timestamp}",
  "architecture": "amd64",
  "os": "linux",
  "config": {
    "User": "1000:1000",
    "Env": [
      "PATH=/bin"
    ],
    "Cmd": [
      "/bin/hi"
    ],
    "WorkingDir": "/"
  },
  "rootfs": {
    "type": "layers",
    "diff_ids": [
      "sha256:${layer_digest}"
    ]
  }
}
EOF

config_digest="$(sha256sum "config.json" | awk '{print $1}')"
config_size="$(wc -c < "config.json")"
```

And a [manifest.json](https://github.com/opencontainers/image-spec/blob/main/manifest.md):

```
cat > "manifest.json" << EOF
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "digest": "sha256:${config_digest}",
    "size": ${config_size}
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar",
      "digest": "sha256:${layer_digest}",
      "size": ${layer_size}
    }
  ]
}
EOF

manifest_digest="$(sha256sum "manifest.json" | awk '{print $1}')"
manifest_size="$(wc -c < "manifest.json")"
```

---

Then we assemble the [OCI image layout](https://specs.opencontainers.org/image-spec/image-layout/).

We rename some files to be content-addressable [blobs](https://github.com/opencontainers/image-spec/blob/main/image-layout.md#blobs):

```
blobs="blobs/sha256"
mkdir -p "${blobs}"
mv "config.json" "${blobs}/${config_digest}"
mv "layer.tar" "${blobs}/${layer_digest}"
mv "manifest.json" "${blobs}/${manifest_digest}"
```

Create an [oci-layout](https://specs.opencontainers.org/image-spec/image-layout/#oci-layout-file) file:

```
cat > "oci-layout" << EOF
{
  "imageLayoutVersion": "1.0.0"
}
EOF
```

And an [index.json](https://github.com/opencontainers/image-spec/blob/main/image-index.md):

```
cat > "index.json" << EOF
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": ${manifest_size},
      "digest": "sha256:${manifest_digest}"
    }
  ]
}
EOF
```

And just like that, everything is ready.

That was not very hard, just a pile of JSON files. 🤔

---

Let's bundle and import the image with the containerd cli:

```
cd ..
tar -cf image.tar -C image .
sudo ctr images import --base-name oci-by-hand:latest --digests --all-platforms image.tar
```

And run it:

```
sudo ctr run --rm "$(sudo ctr images list | grep -F "oci-by-hand@" | awk '{print $1}')" hi
```

Done 😎
