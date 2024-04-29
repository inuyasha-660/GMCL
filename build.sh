#!/usr/bin/env bash

VERSION_BUILD=1.0.0
Line="------------------------------------------------------"

echo "
 ██████╗ ███╗   ███╗ ██████╗██╗     
██╔════╝ ████╗ ████║██╔════╝██║     
██║  ███╗██╔████╔██║██║     ██║     
██║   ██║██║╚██╔╝██║██║     ██║     
╚██████╔╝██║ ╚═╝ ██║╚██████╗███████╗
 ╚═════╝ ╚═╝     ╚═╝ ╚═════╝╚══════╝                 
"

date
echo "Version: "${VERSION_BUILD}

echo ${Line}

echo "==> Build GMCL for linux-amd64"
if 
  mkdir -p bin/linux-amd64
  GOOS=linux 
  GOARCH=amd64 
  go build -o ./bin/linux-amd64/gmcl-${VERSION_BUILD}-linux-amd64 . 
then
  echo "=> Package GMCL Resources"
  cp -r README.md ./Resources ./bin/linux-amd64
  tar -zcvf gmcl-${VERSION_BUILD}-linux-amd64.tar.gz ./bin/linux-amd64
else
  echo "-> Build failed"
fi
