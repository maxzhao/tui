# Fork Homebrew packaging

The inherited formula is disabled: its checksums have not been regenerated from a reviewed fork release. Maxa uses checksum-verified Linux assets from `maxzhao/tui`, not `ZenNotes/homebrew-tap`. `zn update` refuses package-manager/Go/desktop-owned installations rather than invoking an upstream installer.

To maintain an optional fork tap, publish reviewed `X.Y.Z-maxa.N` fork assets first, then run:

```sh
python3 packaging/homebrew/update_formula.py X.Y.Z-maxa.N
python3 -B -m unittest discover -s packaging/homebrew -p 'test_*.py'
```

The updater queries only `maxzhao/tui`, requires all four inherited macOS/Linux archives and their GitHub SHA-256 digests, and removes the disable marker only after validating them. The default fork release workflow publishes Linux archives only; it does not prepare a Homebrew release or push a tap.
