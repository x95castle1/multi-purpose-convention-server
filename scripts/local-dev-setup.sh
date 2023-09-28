#!/bin/bash

check_dependencies() {
  R=$1 && REQUIRED_DEPENDENCIES=${R[*]}
  MISSING_DEPENDENCIES=()

  for DEPENDENCY in ${REQUIRED_DEPENDENCIES}; do
    command -v "$DEPENDENCY" >/dev/null 2>&1 || { MISSING_DEPENDENCIES+=("$DEPENDENCY"); }
  done

  # Instead of exiting when we see a missing command, let's be nice and give the user a list.
  if [ ${#MISSING_DEPENDENCIES[@]} -ne 0 ]; then
    printf '> Missing %s, please install it!\n' "${MISSING_DEPENDENCIES[@]}"
    printf 'Exiting.\n'
    exit 1
  fi
}

install_utility_deps() {
  brew install jq gsed
}

install_dev_deps() {
  brew install go
}

install_required_deps() {
  install_utility_deps
  install_dev_deps
  install_carvel_tools
  install_and_configure_tanzu_cli
  install_and_configure_build_packs_cli
}

install_carvel_tools() {
  brew tap vmware-tanzu/carvel
  brew install ytt kbld kapp imgpkg kwt vendir kctrl
}

install_and_configure_tanzu_cli() {
  brew install vmware-tanzu/tanzu/tanzu-cli
  tanzu plugin install --group vmware-tap/default:v1.6.3
}

install_and_configure_build_packs_cli() {
  brew install go buildpacks/tap/pack
  pack config default-builder paketobuildpacks/builder-jammy-tiny
}

main() {
  check_dependencies "docker git"
  install_required_deps
}

main
