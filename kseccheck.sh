#!/usr/bin/env bash
#
# Kubernetes Resource Security Checker
#
# Andrew Martin, 2017-10-11 21:07:42
# sublimino@gmail.com
#
## Usage: %SCRIPT_NAME% [options] <k8s resource file>
##
## Validate security parameters of a Kubernetes resource
##
## Options:
##   --full             Full console output
##   --json             JSON output
##   -h --help          Display this message
##   --debug            More debug
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
IS_JSON=0
IS_FULL=0
KIND=""

KUBECTL='kubectl'
JQ='jq'
CACHE_DIR=""

POINTS=0
OUTPUT_ADVISE=()
OUTPUT_CRITICAL=()
OUTPUT_POSTITIVE=()
OUTPUT_NEGATIVE=()

main() {
  handle_arguments "$@"

  resolve_kubectl
  resolve_jq

  configure_cache

  if ! JSON=$(read_json_resource_cache "${FILENAME}"); then
    if ! JSON=$(get_json "${FILENAME}") \
      || [[ $(echo "${JSON:0:1}") != '{' ]] \
      || ! get_kind &>/dev/null; then

      error "Invalid input"
    fi

    write_json_resource_cache "${FILENAME}"
  fi

  FULL_JSON="${JSON}"

  check_valid_kind

  RULES=$(get_rules)
  SELECTORS=$(echo "${RULES}" | ${JQ} -r ".selector" | sed 's,\",",g')

  COUNT=0

  # TODO(ajm): /dev/fd/X redirect below fails on Lambda, BASH 4.2 - replaced with tmp file
  local TEMP_FILE=$(mktemp)
  echo "${SELECTORS}" >>"${TEMP_FILE}"

  while read SELECTOR; do
    THIS_POINTS=$(get_points "${SELECTOR}")

    if is_key "${SELECTOR}"; then
      POINTS=$((POINTS + THIS_POINTS))

      THIS_RULE=$(get_rule_by_selector "${SELECTOR}")

      if [[ ${THIS_POINTS} -gt 0 ]]; then
        positive "${THIS_RULE}"
      elif [[ ${THIS_POINTS} -le 10 ]]; then
        critical "${THIS_RULE}"
      else
        negative "${THIS_RULE}"
      fi
    else
      if [[ ${THIS_POINTS} -gt 0 ]]; then
        if get_advise "${SELECTOR}"; then
          THIS_RULE=$(get_rule_by_selector "${SELECTOR}")
          advise "${THIS_RULE}"
        fi
      fi
    fi
    # warning "COUNT is $COUNT : $THIS_POINTS $SELECTOR "
    COUNT=$((COUNT + 1))
  done <"${TEMP_FILE}"

  rule_capabilities
  rule_resources

  print_output

}

print_output() {
  if [[ "${IS_JSON}" == 1 ]]; then
    local JQ_OUT=''
    JQ_OUT=$(${JQ} --null-input \
      ".score |= ${POINTS}")

    for THIS_OUTPUT in "${OUTPUT_CRITICAL[@]:-}"; do
      JQ_OUT=$(echo "${JQ_OUT}" | ${JQ} \
        ".scoring.critical += [${THIS_OUTPUT}]")
    done

    for THIS_OUTPUT in "${OUTPUT_ADVISE[@]:-}"; do
      JQ_OUT=$(echo "${JQ_OUT}" | ${JQ} \
        ".scoring.advise += [${THIS_OUTPUT}]")
    done

    if [[ "${IS_FULL:-}" == 1 ]]; then

      for THIS_OUTPUT in "${OUTPUT_POSITIVE[@]:-}"; do
        JQ_OUT=$(echo "${JQ_OUT}" | ${JQ} \
          ".scoring.positive += [${THIS_OUTPUT}]")
      done

      for THIS_OUTPUT in "${OUTPUT_NEGATIVE[@]:-}"; do
        JQ_OUT=$(echo "${JQ_OUT}" | ${JQ} \
          ".scoring.negative += [${THIS_OUTPUT}]")
      done

    else
      JQ_OUT=$(echo "${JQ_OUT}" | ${JQ} 'del(.scoring[][] | .points)')
    fi

    echo "${JQ_OUT}" | ${JQ} \
      "del(.scoring[][][] | nulls) \
      | .scoring.advise |= \
        (sort_by(.advise) | reverse | .[:5]) \
      | del(.scoring[][] | .advise)"

  else
    if [[ "${IS_FULL:-}" == 1 ]]; then
        output_array "${OUTPUT_CRITICAL[@]:-}"
        output_array "${OUTPUT_ADVISE[@]:-}"
        output_array "${OUTPUT_POSITIVE[@]:-}"
        output_array "${OUTPUT_NEGATIVE[@]:-}"
    fi
    if [[ "${POINTS}" -gt 0 ]]; then
      success "Passed with a score of ${POINTS} points"
    else
      error "Failed with a score of ${POINTS} points"
    fi
  fi
}

output_array() {
  local OUTPUT_ARRAY=("${@:-}")
  for THIS_OUTPUT in "${OUTPUT_ARRAY[@]}"; do
    echo "${THIS_OUTPUT}"
  done
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

rule_seccomp_all_pods() {
  # container.seccomp.security.alpha.kubernetes.io/${container_name}
  :
}

critical() {
  OUTPUT_CRITICAL+=("${1}")
}

advise() {
  OUTPUT_ADVISE+=("${1}")
}

positive() {
  OUTPUT_POSITIVE+=("${1}")
}

negative() {
  OUTPUT_NEGATIVE+=("${1}")
}

get_json() {
  local FILENAME="${1}"
  local COMMAND="${KUBECTL} convert -o json --local=true --filename=\"${FILENAME}\""
  if is_unprivileged_userns_clone; then
    unshare --net \
      --map-root-user \
      ${COMMAND} 2>&1
  else
    ${COMMAND} 2>&1
  fi
}

get_kind() {
  if [[ "${KIND:-}" == "" ]]; then
    KIND=$(echo "${FULL_JSON}" | ${JQ} -r '.kind')
  fi
  echo "${KIND}"
}

check_valid_kind() {
  if ! is_pod; then
    if is_deployment || is_statefulset || is_daemonset; then
      JSON=$(echo "${JSON}" | ${JQ} -r '.spec.template')
    else
      error "Only kinds Pod, Deployment, StatefulSet, DaemonSet accepted"
    fi
  fi

}

get_rules() {
  cat "${DIR}/"k8s-rules.json | ${JQ} ".rules[] | select(.kind == \"$(get_kind)\" or .kind == null)"
}

get_rule_by_selector() {
  local SELECTOR="${1//\"/\\\"}"
  local CACHE_FILE="${CACHE_DIR}/$(echo "${SELECTOR}" | sed 's,[^a-zA-Z],-,g')"
  if [[ -f "${CACHE_FILE}" ]]; then
    cat "${CACHE_FILE}"
  else
    echo "${RULES}" | ${JQ} \
      ". | select(.selector == \"${SELECTOR}\")" \
      | tee "${CACHE_FILE}"
  fi
}

escape_json() {
  :
}

get_points() {
  local SELECTOR="${1}"
  get_rule_by_selector "${SELECTOR}" | ${JQ} '.points'
}

get_advise() {
  local SELECTOR="${1}"
  get_rule_by_selector "${SELECTOR}" | ${JQ} --exit-status '.advise > 0' >/dev/null
}

is_key() {
  local KEY="${1}"
  local TEST_JSON="${JSON}"

  if [[ "${KEY:0:12}" == 'containers[]' ]]; then
    KEY=".spec.${KEY}"
  else
    TEST_JSON="${FULL_JSON}"
  fi

  # TODO: if debug read user input

  local RESULT=$(echo "${TEST_JSON}" | ${JQ} "select(${KEY}) | ${KEY}" 2>&1)
  if [[ "${RESULT}" =~ ^jq:[[:space:]]error ]]; then
    warning "${RESULT}"
    warning "${TEST_JSON}"
    error "${JQ} \"select(${KEY}) | ${KEY}\""
  fi

  [[ "${RESULT}" != 'null' ]] && [[ "${RESULT}" != '' ]]
}

# ---

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

is_unprivileged_userns_clone() {
  command -v sysctl &>/dev/null \
    && sysctl kernel.unprivileged_userns_clone | grep -q ' = 1' \
      && unshare --net \
        --map-root-user \
        touch /dev/null &>/dev/null
}

resolve_jq() {
  if ! JQ=$(resolve_binary ${JQ}); then
    exit 1
  fi
}

resolve_kubectl() {
  if ! KUBECTL=$(resolve_binary ${KUBECTL}); then
    exit 1
  fi
}

resolve_binary() {
  local BINARY="${1}"
  local ORIGINAL_BINARY="${BINARY}"

  if ! BINARY=$(command -v "${BINARY}" 2>/dev/null); then
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

# ---

configure_cache() {
  if [[ "${CACHE_DIR:-}" == "" ]]; then
    CACHE_DIR="/dev/shm/$(echo "${THIS_SCRIPT}" | base64_fs_sanitise)"
    mkdir -p "${CACHE_DIR}"
  fi
}

base64_fs_sanitise() {
  base64 -w0 | tr '/' '-' | sed -r 's,(.{200}),\1/,g'
}

read_json_resource_cache() {
  local FILENAME="${1:-}"
  local CACHE_KEY=$(get_resource_cache_key "${FILENAME}")
  if [[ "${CACHE_KEY}" != "" ]]; then
    if [[ -f "${CACHE_DIR}/${CACHE_KEY}/data" ]]; then
      mkdir -p "${CACHE_DIR}/${CACHE_KEY}"
      cat "${CACHE_DIR}/${CACHE_KEY}/data"
      return 0
    fi
  fi

  return 1
}

get_resource_cache_key() {
  local FILENAME="${1:-}"
  local CACHE_KEY=$(cat "${FILENAME}" | base64_fs_sanitise)
  echo "${CACHE_KEY}"
}

write_json_resource_cache() {
  local FILENAME="${1:-}"
  local CACHE_KEY=$(get_resource_cache_key "${FILENAME}")
  if [[ "${CACHE_KEY}" != "" ]]; then
    mkdir -p  "${CACHE_DIR}/${CACHE_KEY}"
    echo "${JSON}" | tee "${CACHE_DIR}/${CACHE_KEY}/data" >/dev/null
  fi
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
      --json) IS_JSON=1 ;;
      --full) IS_FULL=1 ;;
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
  if [[ "${IS_JSON:-0}" == 1 ]]; then
    return
  fi
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Warning"
  printf "%s\n" "$(log_message_prefix)${COLOUR_RED}${ERROR}${COLOUR_RESET}"
} 1>&2

error() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Error"
  if [[ "${IS_JSON:-0}" == 1 ]]; then
    json_error "${ERROR}"
  else
    printf "%s\n" "$(log_message_prefix)${COLOUR_RED}${ERROR}${COLOUR_RESET}"  1>&2
    exit 3
  fi
}

json_error() {
  ${JQ} --null-input ".error |= \"${*//\"/\\\"}\"" 2>&1
  exit 3
}

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
