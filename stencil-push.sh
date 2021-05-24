set -e

project_version=$1
proto_repo=$2
file_path=$3
output_socket=/dev/null

log_info()  { echo "$(date) - [INFO] $@"; }
log_warn()  { echo "$(date) - [WARN] $@"; }
log_error() { echo "$(date) - [ERROR] $@"; }

help() {
    echo ""
    echo "stencil-push.sh <version> <proto-repo> <file-path>"
    echo "   version      - version of the proto descriptor set being uploaded"
    echo "   proto-repo   - repository name where to upload the proto descriptor set"
    echo "   file-path    - path to the proto descriptor set binary file"
    echo "   "
    echo "   stencil-push.sh needs a few ENV variables set - STENCIL_HOSTNAME, STENCIL_USERNAME, STENCIL_PASSWORD"
    echo "   STENCIL_HOSTNAME - Hostname serving the stencil APIs"
    echo "   STENCIL_USERNAME - HTTP Basic Auth username if required"
    echo "   STENCIL_PASSWORD - HTTP Basic Auth password if required"
}

if [[ -z "$STENCIL_HOSTNAME" ]]
then
    log_error "Missing values for ENV variables"
    help
    exit 1
fi
if [ $# -lt 3 ] || [ $# -gt 3 ]
then
    if [ $# -eq 1 ] && [ $1 == "help" ]
    then
        help
        exit 0
    fi
    log_error "Command requires exactly 3 arguments, $# were passed"
    log_error "Passed arguments - " $@
    help
    exit 1
fi

trap 'log_error "Failed: $ACTION"' EXIT

ACTION="Uploading proto descriptor as version $project_version"
log_info $ACTION
curl -sS -w "\n" -u $STENCIL_USERNAME:$STENCIL_PASSWORD -X PUT --fail "https://$STENCIL_HOSTNAME/artifactory/proto-descriptors/$proto_repo/$project_version" -T "$file_path" > /dev/null

ACTION="Uploading proto descriptor as latest"
log_info $ACTION
curl -sS -w "\n" -u $STENCIL_USERNAME:$STENCIL_PASSWORD -X PUT --fail "https://$STENCIL_HOSTNAME/artifactory/proto-descriptors/$proto_repo/latest" -T "$file_path" > /dev/null

ACTION="Updating metadata to set latest version to $project_version"
log_info $ACTION
curl -sS -w "\n" -u $STENCIL_USERNAME:$STENCIL_PASSWORD -X PUT --fail "https://$STENCIL_HOSTNAME/metadata/proto-descriptors/$proto_repo/version" -d value="$project_version" > /dev/null

ACTION="Upload proto descriptor to stencil service as version $project_version"
log_info $ACTION
curl -sS -w "\n" -u $STENCIL_USERNAME:$STENCIL_PASSWORD -X POST --fail "https://$STENCIL_HOSTNAME/v1/descriptors" -F "file=@$file_path" -F "version=$project_version" -F "name=$proto_repo" -F "latest=true" -H "Content-Type: multipart/form-data" > /dev/null

trap - EXIT
