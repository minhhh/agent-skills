#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILLS_DIR="$(dirname "$SCRIPT_DIR")/skills"
AGENT_DIRS=".agents/skills .config/opencode/skills .claude/skills"

while IFS= read -r skill || [[ -n "$skill" ]]; do
  [[ -z "$skill" || "$skill" =~ ^# ]] && continue
  for dir in $AGENT_DIRS; do
    dest="$HOME/$dir"
    mkdir -p "$dest"
    ln -sf "$SKILLS_DIR/$skill" "$dest/$(basename $skill)"
    echo "Linked: $skill -> $dest/"
  done
done < "$SCRIPT_DIR/config/default-skills.txt"

echo "Default skills linked."
