#!/usr/bin/env python3
"""Idempotent wiring of a new upgrade version into the app.

    wire_upgrade.py --repo DIR --version vNN [--dry-run]

Edits, each only when needed:
  Makefile                     VERSION := vNN
  app/upgrades/vNN/upgrades.go scaffolded from templates/upgrades.go.tmpl
  app/upgrades.go              import + SetUpgradeHandler + StoreUpgrades case

If vNN is already registered in app/upgrades.go the file is left alone
entirely: a hand-wired version has already decided whether it needs a
StoreUpgrades case. Prints one "changed:"/"unchanged:" line per file.
"""
import argparse
import os
import re
import sys

MODULE = "github.com/bitbadges/bitbadgeschain"


def read(p):
    with open(p, encoding="utf-8") as f:
        return f.read()


def write(p, s, dry):
    if not dry:
        os.makedirs(os.path.dirname(p), exist_ok=True)
        with open(p, "w", encoding="utf-8") as f:
            f.write(s)


def report(path, changed):
    print(("changed:   " if changed else "unchanged: ") + path)
    return changed


def bump_makefile(repo, v, dry):
    p = os.path.join(repo, "Makefile")
    s = read(p)
    new, n = re.subn(r"(?m)^VERSION := v\d+[ \t]*$", f"VERSION := {v}", s, count=1)
    if n == 0:
        sys.exit("Makefile has no 'VERSION := vNN' line")
    if new != s:
        write(p, new, dry)
    return report("Makefile", new != s)


def scaffold_handler(repo, v, dry, template):
    p = os.path.join(repo, "app", "upgrades", v, "upgrades.go")
    if os.path.exists(p):
        return report(os.path.relpath(p, repo), False)
    write(p, read(template).replace("__VERSION__", v), dry)
    return report(os.path.relpath(p, repo), True)


REGISTRATION = """\tapp.UpgradeKeeper.SetUpgradeHandler(
\t\t{v}.UpgradeName,
\t\t{v}.CreateUpgradeHandler(
\t\t\tapp.ModuleManager,
\t\t\tapp.Configurator(),
\t\t),
\t)
"""

STORE_CASE = """\tcase {v}.UpgradeName:
\t\tstoreUpgrades = &storetypes.StoreUpgrades{{
\t\t\tRenamed: []storetypes.StoreRename{{}},
\t\t\tDeleted: []string{{}},
\t\t\tAdded:   []string{{}},
\t\t}}
"""


def wire_app(repo, v, dry):
    p = os.path.join(repo, "app", "upgrades.go")
    s = read(p)
    rel = os.path.relpath(p, repo)
    if re.search(rf"\b{v}\.UpgradeName\b", s):
        return report(rel, False)

    imp = f'\t{v} "{MODULE}/app/upgrades/{v}"\n'
    imports = list(re.finditer(rf'(?m)^\tv\d+ "{re.escape(MODULE)}/app/upgrades/v\d+"\n', s))
    if not imports:
        sys.exit("app/upgrades.go: no existing app/upgrades/vNN import to anchor on")
    s = s[: imports[-1].end()] + imp + s[imports[-1].end():]

    # Register after the last SetUpgradeHandler(...) call, before the
    # ReadUpgradeInfoFromDisk block.
    calls = list(re.finditer(r"(?m)^\tapp\.UpgradeKeeper\.SetUpgradeHandler\(\n(?:.*\n)*?\t\)\n", s))
    if not calls:
        sys.exit("app/upgrades.go: no SetUpgradeHandler call to anchor on")
    end = calls[-1].end()
    s = s[:end] + REGISTRATION.format(v=v) + s[end:]

    m = re.search(r"(?m)^\tswitch upgradeInfo\.Name \{\n", s)
    if not m:
        sys.exit("app/upgrades.go: no 'switch upgradeInfo.Name {' block")
    close = re.search(r"(?m)^\t\}\n", s[m.end():])
    if not close:
        sys.exit("app/upgrades.go: switch block never closes")
    at = m.end() + close.start()
    s = s[:at] + STORE_CASE.format(v=v) + s[at:]

    write(p, s, dry)
    return report(rel, True)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--repo", required=True)
    ap.add_argument("--version", required=True)
    ap.add_argument("--dry-run", action="store_true")
    ap.add_argument("--template", default=os.path.join(os.path.dirname(__file__), "..", "templates", "upgrades.go.tmpl"))
    a = ap.parse_args()
    if not re.fullmatch(r"v\d+", a.version):
        sys.exit(f"version must look like v35, got {a.version!r}")
    bump_makefile(a.repo, a.version, a.dry_run)
    scaffold_handler(a.repo, a.version, a.dry_run, a.template)
    wire_app(a.repo, a.version, a.dry_run)


if __name__ == "__main__":
    main()
