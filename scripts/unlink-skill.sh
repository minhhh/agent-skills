#!/bin/bash
set -e

SKILL=$1
if [[ -z "$SKILL" ]]; then
  echo "Usage: ./scripts/unlink-skill.sh <skill-name>"
  exit 1
fi

AGENT_DIRS=".agents/skills .claude/skills"

for dir in $AGENT_DIRS; do
  dest="$HOME/$dir"
  rm -f "$dest/$SKILL"
  echo "Unlinked: $SKILL from $dest/"
done
