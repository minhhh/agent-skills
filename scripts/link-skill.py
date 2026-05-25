#!/usr/bin/env python3
import os
import sys
from pathlib import Path

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 link-skill.py <skill-name>", file=sys.stderr)
        sys.exit(1)
    
    skill_name = sys.argv[1]
    
    repo_root = Path(__file__).resolve().parent.parent
    skills_dir = repo_root / "skills"
    private_skills_dir = repo_root / "private-agent-skills" / "skills"
    agent_dirs = [".agents/skills", ".claude/skills"]
    
    # Find skill
    src_path = None
    if (skills_dir / skill_name).is_dir():
        src_path = skills_dir / skill_name
    elif (private_skills_dir / skill_name).is_dir():
        src_path = private_skills_dir / skill_name
        
    if not src_path:
        print(f"Error: Skill '{skill_name}' not found in skills/ or private-agent-skills/skills/", file=sys.stderr)
        sys.exit(1)
        
    home = Path.home()
    for agent_dir in agent_dirs:
        dest_dir = home / agent_dir
        dest_dir.mkdir(parents=True, exist_ok=True)
        dest_path = dest_dir / skill_name
        
        # Remove existing symlink or file if present
        if dest_path.exists() or dest_path.is_symlink():
            dest_path.unlink()
            
        dest_path.symlink_to(src_path)
        print(f"Linked: {skill_name} -> {dest_dir}/")

if __name__ == "__main__":
    main()
