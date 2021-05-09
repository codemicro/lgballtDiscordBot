# Patches

Files in this directory are patches to be applied to vendored libraries.

New patches should have the directory they should be run in *relative to the repository root* defined in `patches.json`.

Before committing a patch, ensure it is applied using `mage ApplyPatches`. Patches are not automatically applied when building.

## Generating a patch

```bash
git clone https://example.com/blah.git # clone library repository
code . # make your edits, don't commit them
git diff > libraryname.patch
```