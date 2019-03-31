#!/bin/bash
# This entrypoint is largely based off of https://sosedoff.com/2019/02/12/go-github-actions.html

set -eo pipefail

# Assumes to be set by base docker image
if [[ -z "${GOPATH}" ]]; then
  echo "Environment variable GOPATH must be set."
  exit 1
fi

# see https://developer.github.com/actions/creating-github-actions/accessing-the-runtime-environment/#environment-variables
# Assumes to be set by GitHub Actions
if [[ -z "${GITHUB_WORKSPACE}" ]]; then
  echo "Environment variable GITHUB_WORKSPACE must be set."
  exit 1
fi

# see https://developer.github.com/actions/creating-github-actions/accessing-the-runtime-environment/#environment-variables
# Assumes to be set by GitHub Actions
if [[ -z "${GITHUB_REPOSITORY}" ]]; then
  echo "Environment variable GITHUB_REPOSITORY must be set."
  exit 1
fi

declare workdir="${GOPATH}/src/github.com/${GITHUB_REPOSITORY}"
declare releases="${GITHUB_WORKSPACE}/.releases"
declare targets=${@-"darwin/amd64 darwin/386 linux/amd64 linux/386 windows/amd64 windows/386"}
declare ghproject="$(echo ${GITHUB_REPOSITORY} | cut -d '/' -f2)"

mkdir -p ${workdir}
mkdir -p ${releases}

cp -a ${GITHUB_WORKSPACE}/* ${workdir}/

cd ${workdir}

dep ensure && go test ./...

for target in ${targets}; do
  os="$(echo ${target} | cut -d '/' -f1)"
  arch="$(echo ${target} | cut -d '/' -f2)"
  output="${releases}/${ghproject}_${os}_${arch}"

  echo "Building: $target"
  GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build -o ${output}
  zip -j ${output}.zip ${output} > /dev/null
done
