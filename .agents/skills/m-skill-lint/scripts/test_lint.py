import unittest
import sys
from pathlib import Path

# Add the scripts directory to sys.path so we can import lint
SCRIPT_DIR = Path(__file__).resolve().parent
sys.path.append(str(SCRIPT_DIR))

import importlib.machinery
import importlib.util

try:
    loader = importlib.machinery.SourceFileLoader("lint", str(SCRIPT_DIR / "lint"))
    spec = importlib.util.spec_from_loader("lint", loader)
    lint = importlib.util.module_from_spec(spec)
    loader.exec_module(lint)
    validate_frontmatter = getattr(lint, "validate_frontmatter", None)
except Exception as e:
    print(f"Error loading lint module: {e}")
    validate_frontmatter = None

class TestFrontmatterValidation(unittest.TestCase):
    def setUp(self):
        if validate_frontmatter is None:
            self.fail("validate_frontmatter function is not implemented or cannot be imported from lint.py")

    def test_valid_frontmatter(self):
        content = (
            "---\n"
            "name: my-skill-name\n"
            "description: Use when you need to perform tests.\n"
            "---\n"
            "\n"
            "# My Skill Name\n"
        )
        errors = validate_frontmatter(content)
        self.assertEqual(errors, [])

    def test_valid_multiline_description(self):
        content = (
            "---\n"
            "name: m-skill-lint\n"
            "description: >\n"
            "  Use when validating skill directory structure\n"
            "  and conventions.\n"
            "---\n"
        )
        errors = validate_frontmatter(content)
        self.assertEqual(errors, [])

    def test_missing_frontmatter(self):
        content = "# Simple markdown without frontmatter\n"
        errors = validate_frontmatter(content)
        self.assertTrue(len(errors) > 0, "Expected errors for missing frontmatter")
        self.assertTrue(any("missing" in err.lower() or "invalid" in err.lower() for err in errors))

    def test_invalid_yaml_delimiters(self):
        content = "---\nname: test\ndescription: Use when testing\n"
        errors = validate_frontmatter(content)
        self.assertTrue(len(errors) > 0, "Expected errors for invalid YAML delimiters")

    def test_unsupported_fields(self):
        content = (
            "---\n"
            "name: my-skill\n"
            "description: Use when testing.\n"
            "version: 1.0\n"
            "author: Test Author\n"
            "---\n"
        )
        errors = validate_frontmatter(content)
        self.assertTrue(len(errors) > 0, "Expected errors for unsupported fields")
        self.assertTrue(any("unsupported" in err.lower() or "extra" in err.lower() or "field" in err.lower() for err in errors))

    def test_name_format_invalid(self):
        bad_names = ["my skill", "my(skill)", "my_skill", "MySkill!"]
        for name in bad_names:
            content = (
                "---\n"
                f"name: {name}\n"
                "description: Use when testing.\n"
                "---\n"
            )
            errors = validate_frontmatter(content)
            self.assertTrue(
                len(errors) > 0,
                f"Failed to catch invalid name: {name}"
            )
            self.assertTrue(
                any("name" in err.lower() or "character" in err.lower() or "letters" in err.lower() for err in errors),
                f"Error message for invalid name was unexpected: {errors}"
            )

    def test_description_format_invalid(self):
        bad_descriptions = [
            "This is a description.",
            "Helper for doing things.",
            "Used when we need it.",
            "When we need it.",
        ]
        for desc in bad_descriptions:
            content = (
                "---\n"
                "name: my-skill\n"
                f"description: {desc}\n"
                "---\n"
            )
            errors = validate_frontmatter(content)
            self.assertTrue(
                len(errors) > 0,
                f"Failed to catch invalid description: {desc}"
            )
            self.assertTrue(
                any("description" in err.lower() and "use when" in err.lower() for err in errors),
                f"Error message for invalid description was unexpected: {errors}"
            )

    def test_frontmatter_too_long(self):
        long_desc = "Use when testing. " + ("x" * 1050)
        content = (
            "---\n"
            "name: my-skill\n"
            f"description: {long_desc}\n"
            "---\n"
        )
        errors = validate_frontmatter(content)
        self.assertTrue(len(errors) > 0, "Expected error for too long frontmatter")
        self.assertTrue(any("exceed" in err.lower() or "length" in err.lower() or "1024" in err.lower() for err in errors))

if __name__ == "__main__":
    unittest.main()
