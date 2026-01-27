#!/usr/bin/env bash
set -euo pipefail

# ---- settings ----
BASE_PORT=19999
NUM_FILES=66
OUT_DIR="."
BASE_INSTANCE_DIR="/home/maruyama/workspace/netdata/netdata-dev/instances"
RUN_AS_USER="maruyama"
# ------------------

mkdir -p "${OUT_DIR}"

for i in $(seq 1 "${NUM_FILES}"); do
  port=$((BASE_PORT + i - 1))
  conf="${OUT_DIR}/netdata${i}.conf"
  inst_dir="${BASE_INSTANCE_DIR}/${port}"

  # instance directories (optional but handy)
  mkdir -p "${inst_dir}/"{log,cache,lib,run}

  cat >"${conf}" <<EOF
[global]
  run as user = ${RUN_AS_USER}
  web files owner = ${RUN_AS_USER}

[web]
  bind to = 0.0.0.0:${port}

[directories]
  log = ${inst_dir}/log
  cache = ${inst_dir}/cache
  lib = ${inst_dir}/lib
  run = ${inst_dir}/run

[cloud]
  enabled = no
  conversation log = no

[logs]
  aclk = none
EOF

  echo "generated: ${conf} (port=${port})"
done
