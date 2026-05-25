## Manage agent skill symlinks (skills/ → ~/.agents/skills/ and ~/.claude/skills/)

help: # show this help
	@echo ""
	@grep "^##" $(MAKEFILE_LIST) | grep -v grep
	@echo ""
	@grep "^[0-9a-zA-Z\-]*: #" $(MAKEFILE_LIST) | grep -v grep
	@echo ""

link-default: # link all skills in scripts/config/default-skills.txt (currently: git-conventionalize)
	python3 ./scripts/link-default.py

link-skill: # link a single skill by name, e.g. make link-skill SKILL=git-conventionalize
	python3 ./scripts/link-skill.py $(SKILL)

unlink-skill: # remove a linked skill symlink, e.g. make unlink-skill SKILL=git-conventionalize
	python3 ./scripts/unlink-skill.py $(SKILL)

