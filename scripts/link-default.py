#!/usr/bin/env python3
import os
import sys
from pathlib import Path

def main():
    repo_root = Path(__file__).resolve().parent.parent
    config_file = repo_root / "scripts" / "config" / "default-skills.txt"
    
    if not config_file.exists():
        print(f"Error: Config file not found at {config_file}", file=sys.stderr)
        sys.exit(1)
        
    skills_dir = repo_root / "skills"
    private_skills_dir = repo_root / "private-agent-skills" / "skills"
    agent_dirs = [".agents/skills", ".claude/skills"]
    home = Path.home()
    
    with open(config_file, "r", encoding="utf-8") as f:
        skills = [line.strip() for line in f if line.strip() and not line.strip().startswith("#")]
        
    for skill_name in skills:
        src_path = None
        if (skills_dir / skill_name).is_dir():
            src_path = skills_dir / skill_name
        elif (private_skills_dir / skill_name).is_dir():
            src_path = private_skills_dir / skill_name
            
        if not src_path:
            print(f"Warning: Skill '{skill_name}' not found in skills/ or private-agent-skills/skills/", file=sys.stderr)
            continue
            
        for agent_dir in agent_dirs:
            dest_dir = home / agent_dir
            dest_dir.mkdir(parents=True, exist_ok=True)
            dest_path = dest_dir / skill_name
            
            if dest_path.exists() or dest_path.is_symlink():
                dest_path.unlink()
                
            dest_path.symlink_to(src_path)
            print(f"Linked: {skill_name} -> {dest_dir}/")
            
    print("Default skills linked.")

if __name__ == "__main__":
    main()
