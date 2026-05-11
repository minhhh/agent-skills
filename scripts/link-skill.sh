#!/bin/bash
set -e

SKILL=$1
if [[ -z "$SKILL" ]]; then
  echo "Usage: ./scripts/link-skill.sh <skill-name>"
  exit 1
fi

SKILLS_DIR="$(cd "$(dirname "$0")/../skills" && pwd)"
AGENT_DIRS=".agents/skills .claude/skills"

for dir in $AGENT_DIRS; do
  dest="$HOME/$dir"
  mkdir -p "$dest"
  rm -f "$dest/$SKILL"
  ln -s "$SKILLS_DIR/$SKILL" "$dest/$SKILL"
  echo "Linked: $SKILL -> $dest/"
done
