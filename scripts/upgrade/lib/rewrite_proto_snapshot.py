#!/usr/bin/env python3
"""Rewrite a snapshot of proto/tokenization/*.proto into a versioned package.

    rewrite_proto_snapshot.py DIR vNN

package tokenization            -> package tokenization.vNN
import "tokenization/x.proto"   -> import "tokenization/vNN/x.proto"
option go_package = ".../types" -> ".../types/vNN"
"""
import glob
import os
import sys

MODULE = "github.com/bitbadges/bitbadgeschain"


def main():
    d, v = sys.argv[1], sys.argv[2]
    for p in sorted(glob.glob(os.path.join(d, "*.proto"))):
        with open(p, encoding="utf-8") as f:
            s = f.read()
        s = s.replace("package tokenization;", f"package tokenization.{v};")
        s = s.replace('"tokenization/', f'"tokenization/{v}/')
        s = s.replace(f'option go_package = "{MODULE}/x/tokenization/types";',
                      f'option go_package = "{MODULE}/x/tokenization/types/{v}";')
        with open(p, "w", encoding="utf-8") as f:
            f.write(s)


if __name__ == "__main__":
    main()
