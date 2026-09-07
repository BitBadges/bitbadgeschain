#!/usr/bin/env python3
"""Rebuild the Solidity bytecode used by the v35 integration regression."""
import hashlib
import json
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[3]
IMAGE = "ethereum/solc:0.8.24@sha256:e56ef5e376ae846f06b919d7ca4ed0c271f7fb0900daa6c660d53451f5bfd9db"
NAMES = ["RealEstateSecurityToken", "PrivateEquityToken", "TwoFactorSecurityToken", "CarbonCreditToken"]
sources = {}

def visit(path):
    path = path.resolve()
    relative = str(path.relative_to(ROOT))
    if relative in sources:
        return
    data = path.read_bytes()
    sources[relative] = hashlib.sha256(data).hexdigest()
    for dependency in re.findall(r'(?m)^import\s+(?:[^;]*?from\s+)?[\"\']([^\"\']+)[\"\']', data.decode()):
        visit(path.parent / dependency)

files = [f"contracts/examples/{name}.sol" for name in NAMES]
for source in files:
    visit(ROOT / source)
result = subprocess.run([
    "docker", "run", "--rm", "--platform", "linux/amd64",
    "-v", f"{ROOT}:/src:ro", "-w", "/src", IMAGE,
    "--via-ir", "--optimize", "--combined-json", "abi,bin", *files
], check=True, capture_output=True, text=True)
compiled = json.loads(result.stdout)["contracts"]
output = ROOT / "x/tokenization/precompile/test/integration/testdata"
output.mkdir(exist_ok=True)
(output / "v35_examples.json").write_text(json.dumps({
    name: compiled[f"contracts/examples/{name}.sol:{name}"] for name in NAMES
}, indent=2) + "\n")
(output / "v35_examples_sources.json").write_text(json.dumps(sources, sort_keys=True, indent=2) + "\n")
