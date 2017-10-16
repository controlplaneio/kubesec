#!/usr/bin/env bash
#
# Kubernetes Pod Security Checker
#
# Andrew Martin, 2017-10-11 21:07:42
# sublimino@gmail.com
#
## Usage: %SCRIPT_NAME% [options] <k8s resource>
##
## Validate security parameters of a Kubernetes resource
##
## Options:
##   --debug                     More debug
##
##   -h --help                   Display this message
##

# exit on error or pipe failure
set -eo pipefail
# error on unset variable
set -o nounset
# error on clobber
set -o noclobber

# user defaults
DESCRIPTION="Kubernetes Pod Security Checker"
DEBUG=0

# resolved directory and self
declare -r DIR=$(cd "$(dirname "$0")" && pwd)
declare -r THIS_SCRIPT="${DIR}/$(basename "$0")"

# required defaults
declare -a ARGUMENTS
EXPECTED_NUM_ARGUMENTS=1
ARGUMENTS=()
FILENAME=''
JSON=''
FULL_JSON=''
POINTS=0

KUBECTL='kubectl'
JQ='jq'

KEYS_PLUS_ONE_POINT=(
  'seLinux'
  'supplementalGroups'
  'runAsUser'
  'runAsUser'
  'fsGroup'
  'allowedCapabilities'
  'containers[].securityContext.capabilities.drop'
  'containers[].securityContext.capabilities.drop | index("ALL")'
  'containers[].securityContext.runAsNonRoot == true'
  'containers[].securityContext.runAsUser > 10000'
  'containers[].securityContext.readOnlyRootFilesystem == true'
  'containers[].resources.limits.cpu'
  'containers[].resources.limits.memory'
  'containers[].resources.requests.cpu'
  'containers[].resources.requests.memory'

  # TODO: Root keys, delve
  #  'securityContext'
  #  'securityContext.capabilities'

  # TODO: where is this?
  #  'securityContext.allowPrivilegeEscalation == false'
)

KEYS_MINUS_ONE_POINT=(
  'seLinux'
  'supplementalGroups'
  'runAsUser'
  'runAsUser'
  'fsGroup'
  'allowedCapabilities'
)

KEYS_FAIL=(
  'containers[].securityContext.capabilities.add | index("SYS_ADMIN")'
)

KEYS_STATEFULSET_PLUS_ONE_POINT=(
  '.spec.volumeClaimTemplates[].spec.accessModes | index("ReadWriteOnce")'
  '.spec.volumeClaimTemplates[].spec.resources.requests.storage'
)

#resolve_kubectl() {
#  if ! command -v "${KUBECTL}" &>/dev/null; then
#    KUBECTL='./kubectl'
#
#    if ! command -v "${KUBECTL}" &>/dev/null; then
#      KUBECTL=$(find . -regex '.*/kubectl$' -type f -executable -print -quit)
#
#      if [[ "${KUBECTL:-}" == "" ]]; then
#        KUBECTL=$(find ../ -regex '.*/kubectl$' -type f -executable -print -quit)
#
#        if [[ "${KUBECTL:-}" == "" ]]; then
#          error "kubectl not found"
#        fi
#      fi
#    fi
#  fi
#}

resolve_jq() {
  if ! JQ=$(resolve_binary jq); then
    exit 1
  fi
}

resolve_kubectl() {
  if ! KUBECTL=$(resolve_binary kubectl); then
    exit 1
  fi
}

resolve_binary() {
  local BINARY="${1}"
  local ORIGINAL_BINARY="${BINARY}"

  if ! command -v "${BINARY}" &>/dev/null; then
    BINARY="./${ORIGINAL_BINARY}"

    if ! command -v "${BINARY}" &>/dev/null; then
      BINARY=$(find . -regex ".*/${ORIGINAL_BINARY}$" -type f -executable -print -quit)

      if [[ "${BINARY:-}" == "" ]]; then
        BINARY=$(find ../ -regex ".*/${ORIGINAL_BINARY}$" -type f -executable -print -quit)

        if [[ "${BINARY:-}" == "" ]]; then
          error "${BINARY} not found"
        fi
      fi
    fi
  fi
  echo "${BINARY}"
}

main() {
  handle_arguments "$@"

  resolve_kubectl
  resolve_jq

  JSON=$(get_json "${FILENAME}")
  FULL_JSON="${JSON}"

  check_valid_kind

  for KEY in "${KEYS_PLUS_ONE_POINT[@]}"; do
    if is_key "${KEY}"; then
      POINTS=$((POINTS + 1))
    else
      advise "Add" "${KEY}"
    fi
  done

  for KEY in "${KEYS_FAIL[@]}"; do
    if is_key "${KEY}"; then
      advise "Remove" "${KEY}"
      POINTS=0
    fi
  done

  rule_capabilities
  rule_resources

  rule_statefulset

  if [[ "${POINTS}" -gt 0 ]]; then
    success "Passed with ${POINTS} points"
  else
    error "Failed with 0 points"
  fi
}

rule_capabilities() {
  # only drop: 5
  # drop + add: 3
  # only add: -5
  :
}

rule_resources() {
  # only limits: 1
  # only requests: 1
  # limits and requests: 3
  :
}

advise() {
  local SPEC_PATH="${3-.spec.}"
  echo "${1}" "${SPEC_PATH}${2}"
}

rule_statefulset() {
  if is_statefulset; then
    for KEY in "${KEYS_STATEFULSET_PLUS_ONE_POINT[@]}"; do
      if is_key_full_json "${KEY}"; then
        POINTS=$((POINTS + 1))
      else
        advise "Add" "${KEY}" ''
      fi
    done
  fi
}

get_json() {
  local FILENAME="${1}"
  ${KUBECTL} convert -o json --local=true --filename="${FILENAME}"
}

get_kind() {
  echo "${FULL_JSON}" | ${JQ} -r '.kind'
}

check_valid_kind() {
  echo "Type: $(get_kind)" >&2
  if ! is_pod; then
    if is_deployment || is_statefulset || is_daemonset; then
      JSON=$(echo "${JSON}" | ${JQ} -r '.spec.template')
    else
      error "Only kinds Pod, Deployment, StatefulSet, DaemonSet accepted"
    fi
  fi

}

is_key() {
  local KEY="${1}"

  #  if ! ${JQ} "select(.spec.${KEY}) | .spec.${KEY}"; then
  #    error "jq error"
  #  fi

  # TODO: if debug read user input

  local RESULT=$(echo "${JSON}" | ${JQ} "select(.spec.${KEY}) | .spec.${KEY}")
  [[ "${RESULT}" != 'null' ]] && [[ "${RESULT}" != '' ]]
}

is_key_full_json() {
  local KEY="${1}"

  #  if ! ${JQ} "select(.spec.${KEY}) | .spec.${KEY}"; then
  #    error "jq error"
  #  fi

  # TODO: if debug read user input

  local RESULT=$(echo "${FULL_JSON}" | ${JQ} "select(${KEY}) | ${KEY}")
  [[ "${RESULT}" != 'null' ]] && [[ "${RESULT}" != '' ]]
}

is_pod() {
  _is_type 'Pod'
}

is_deployment() {
  _is_type 'Deployment'
}

is_daemonset() {
  _is_type 'DaemonSet'
}

is_statefulset() {
  _is_type 'StatefulSet'
}

_is_type() {
  local TYPE="${1:-}"
  [[ $(get_kind) == "${TYPE}" ]]

}

# ---

handle_arguments() {
  [[ $# = 0 && ${EXPECTED_NUM_ARGUMENTS} -gt 0 ]] && usage

  parse_arguments "$@"
  validate_arguments "$@"
}

parse_arguments() {
  while [ $# -gt 0 ]; do
    case $1 in
      -h | --help) usage ;;
      --debug)
        DEBUG=1
        set -xe
        ;;
      --)
        shift
        break
        ;;
      -*) usage "$1: unknown option" ;;
      *) ARGUMENTS+=("$1") ;;
    esac
    shift
  done
}

validate_arguments() {
  [[ "${#ARGUMENTS[@]}" != '1' ]] && usage "Single filename required"
  FILENAME="${ARGUMENTS[0]}"

  [[ -f "${FILENAME}" ]] || error "File ${FILENAME} does not exist"

  check_number_of_expected_arguments
}

# helper functions

usage() {
  [ "$*" ] && echo "${THIS_SCRIPT}: ${COLOUR_RED}$*${COLOUR_RESET}" && echo
  sed -n '/^##/,/^$/s/^## \{0,1\}//p' "${THIS_SCRIPT}" | sed "s/%SCRIPT_NAME%/$(basename "${THIS_SCRIPT}")/g"
  exit 2
} 2>/dev/null

success() {
  [ "${*:-}" ] && RESPONSE="$*" || RESPONSE="Unknown Success"
  printf "%s\n" "$(log_message_prefix)${COLOUR_GREEN}${RESPONSE}${COLOUR_RESET}"
} 1>&2

info() {
  [ "${*:-}" ] && INFO="$*" || INFO="Unknown Info"
  printf "%s\n" "$(log_message_prefix)${COLOUR_WHITE}${INFO}${COLOUR_RESET}"
} 1>&2

warning() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Warning"
  printf "%s\n" "$(log_message_prefix)${COLOUR_RED}${ERROR}${COLOUR_RESET}"
} 1>&2

error() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Error"
  printf "%s\n" "$(log_message_prefix)${COLOUR_RED}${ERROR}${COLOUR_RESET}"
  exit 3
} 1>&2

error_env_var() {
  error "${1} environment variable required"
}

log_message_prefix() {
  local TIMESTAMP="[$(date +'%Y-%m-%dT%H:%M:%S%z')]"
  local THIS_SCRIPT_SHORT=${THIS_SCRIPT/$DIR/.}
  tput bold 2>/dev/null
  echo -n "${TIMESTAMP} ${THIS_SCRIPT_SHORT}: "
}

is_empty() {
  [[ -z ${1-} ]] && return 0 || return 1
}

not_empty_or_usage() {
  is_empty "${1-}" && usage "Non-empty value required" || return 0
}

check_number_of_expected_arguments() {
  [[ "${EXPECTED_NUM_ARGUMENTS}" != "${#ARGUMENTS[@]}" ]] && {
    ARGUMENTS_STRING="argument"
    [[ "${EXPECTED_NUM_ARGUMENTS}" -gt 1 ]] && ARGUMENTS_STRING="${ARGUMENTS_STRING}"s
    usage "${EXPECTED_NUM_ARGUMENTS} ${ARGUMENTS_STRING} expected, ${#ARGUMENTS[@]} found"
  }
  return 0
}

hr() {
  printf '=%.0s' $(seq $(tput cols))
  echo
}

wait_safe() {
  local PIDS="${1}"
  for JOB in ${PIDS}; do
    wait "${JOB}"
  done
}

export CLICOLOR=1
export TERM="xterm-color"
export COLOUR_BLACK=$(tput setaf 0 :-"" 2>/dev/null)
export COLOUR_RED=$(tput setaf 1 :-"" 2>/dev/null)
export COLOUR_GREEN=$(tput setaf 2 :-"" 2>/dev/null)
export COLOUR_YELLOW=$(tput setaf 3 :-"" 2>/dev/null)
export COLOUR_BLUE=$(tput setaf 4 :-"" 2>/dev/null)
export COLOUR_MAGENTA=$(tput setaf 5 :-"" 2>/dev/null)
export COLOUR_CYAN=$(tput setaf 6 :-"" 2>/dev/null)
export COLOUR_WHITE=$(tput setaf 7 :-"" 2>/dev/null)
export COLOUR_RESET=$(tput sgr0 :-"" 2>/dev/null)

main "$@"
