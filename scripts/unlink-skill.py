#!/usr/bin/env python3
import os
import sys
from pathlib import Path

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 unlink-skill.py <skill-name>", file=sys.stderr)
        sys.exit(1)
        
    skill_name = sys.argv[1]
    agent_dirs = [".agents/skills", ".claude/skills"]
    
    home = Path.home()
    for agent_dir in agent_dirs:
        dest_path = home / agent_dir / skill_name
        if dest_path.exists() or dest_path.is_symlink():
            dest_path.unlink()
            print(f"Unlinked: {skill_name} from {home / agent_dir}/")
        else:
            print(f"Not linked: {skill_name} in {home / agent_dir}/")

if __name__ == "__main__":
    main()
