#!/usr/bin/env bash
set -euo pipefail

project_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"

output_directory="$project_root/dist"
output_path="$output_directory/PlsRemberDaily"
legacy_windows_resource_path="$project_root/cmd/pls-rember-daily/rsrc.syso"

build_tags=()

for argument in "$@"; do
    case "$argument" in
        -AllowMultipleRuns)
            build_tags+=("allowmultipleruns")
            echo "building with AllowMultipleRuns"
            ;;
        -SkipStartupDelay)
            build_tags+=("skipstartupdelay")
            echo "building with SkipStartupDelay"
            ;;
        *)
            echo "unknown argument: $argument" >&2
            exit 1
            ;;
    esac
done

if ! command -v go >/dev/null 2>&1; then
    echo "go was not found" >&2
    exit 1
fi

if ! command -v pkg-config >/dev/null 2>&1; then
    echo "pkg-config was not found" >&2
    exit 1
fi

if ! pkg-config --atleast-version=4.8 gtk4 gtk4-x11; then
    echo "GTK 4.8 or newer X11 development files were not found" >&2
    exit 1
fi

if [[ "$(go env CGO_ENABLED)" != "1" ]]; then
    echo "cgo must be enabled to build the Linux application" >&2
    exit 1
fi

cd "$project_root"
mkdir -p "$output_directory"
rm -f -- "$legacy_windows_resource_path" "$output_path"

go test ./...
go vet ./...

build_arguments=(-buildvcs=false -o "$output_path")

if ((${#build_tags[@]} > 0)); then
    tags="$(IFS=,; echo "${build_tags[*]}")"
    build_arguments+=(-tags "$tags")
fi

build_arguments+=(./cmd/pls-rember-daily)

go build "${build_arguments[@]}"

echo "built $output_path"
