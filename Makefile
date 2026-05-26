## Manage agent skill symlinks (skills/ → ~/.agents/skills/, ~/.claude/skills/, and ~/.gemini/skills/)

help: # show this help
	@echo ""
	@grep "^##" $(MAKEFILE_LIST) | grep -v grep
	@echo ""
	@grep "^[0-9a-zA-Z\-]*: #" $(MAKEFILE_LIST) | grep -v grep
	@echo ""

link-default: # link all skills in scripts/config/default-skills.txt (currently: git-conventionalize)
	./scripts/link-default

link-skill: # link a single skill by name, e.g. make link-skill SKILL=git-conventionalize
	./scripts/link-skill $(SKILL)

unlink-skill: # remove a linked skill symlink, e.g. make unlink-skill SKILL=git-conventionalize
	./scripts/unlink-skill $(SKILL)

