.PHONY: link-default link-skill unlink-skill

link-default:
	./scripts/link-default.sh

link-skill:
	./scripts/link-skill.sh $(SKILL)

unlink-skill:
	./scripts/unlink-skill.sh $(SKILL)
