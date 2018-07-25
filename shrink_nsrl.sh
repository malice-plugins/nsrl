#!/bin/sh

# copyright: (c) 2014 by Josh "blacktop" Maine.
# license: MIT

set -ex

RDS_URL=https://s3.amazonaws.com/rds.nsrl.nist.gov/RDS/current/rds_modernm.zip
RDS_SHA_URL=https://s3.amazonaws.com/rds.nsrl.nist.gov/RDS/current/rds_modernm.zip.sha

if ls /nsrl/*.zip 1> /dev/null 2>&1; then
   echo "File '.zip' Exists."
else
    echo "[INFO] Downloading NSRL Reduced Sets..."
    wget --progress=bar:force -P /nsrl/ $RDS_URL
    wget --progress=bar:force -P /nsrl/ $RDS_SHA_URL
    echo " * files downloaded"
    ls -lah /nsrl
    RDS_SHA1=$(cat /nsrl/rds_modernm.zip.sha | grep -o -E -e "[0-9a-f]{40}")
    echo " * checking downloaded ZIPs sha1 hash"
    if [ "$RDS_SHA1" ]; then
      echo "$RDS_SHA1 */nsrl/rds_modernm.zip" | sha1sum -c -; \
    fi
fi

echo "[INFO] Unzip NSRL Database zip to /nsrl/ ..."
# 7za x -o/nsrl/ /nsrl/*.zip
cd /nsrl && unzip *.zip

echo "[INFO] Build bloomfilter from NSRL Database ..."
cd rds_modernm && /bin/nsrl --verbose build

echo "[INFO] Listing created files ..."
ls -lah /nsrl/rds_modernm

echo "[INFO] Saving uncompressed NSRL DB size..."
ls -lah NSRLFile.txt | awk '{print $5}' > /nsrl/DBSZIE

echo "[INFO] Saving bloomfilter size..."
ls -lah nsrl.bloom | awk '{print $5}' > /nsrl/BLOOMSIZE

echo "[INFO] Deleting all unused files ..."
rm -rf /nsrl/rds_modernm
rm -f *.zip *.txt *.sh *.sha
ls -lah /nsrl
