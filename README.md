# UVE - Lightweight UV Environment Manager
UVE provides conda-like workflows for UV virtual environments with minimal overhead.

## Key Features
- 🚀 Create/delete Python environments from any directory
- ⚡ Powered by UV for blazing-fast operations
- 🔄 Activate/deactivate environments globally
- 📦 Auto-installs `uv` and `pip` in new environments

## Quick Start
```bash
# Install
curl -L https://github.com/iamshreeram/uve/releases/latest/download/uve-install.sh | bash

# Usage
uve create newpy311 3.11
uve activate newpy311
uve deactivate
```

## Why UVE?
- No Anaconda bloat
- Cross-platform support
- Centralized environment management
- Seamless UV integration


## New in v0.2.0: Environment Cloning

Quickly duplicate environments:
```bash
uve clone my-env my-env-backup  # Create identical copy
uve clone py311 py311-debug      # Make experimental copy
```
Works exactly like conda create --clone but with UV's speed.

### **Why This Works With Minimal Changes**
1. **Leverages Existing UV Features**:
   - Uses `uv venv --python` to recreate the base environment
   - Copies packages manually (no new dependency resolution needed)

2. **Fits Current Architecture**:
   - Uses same environment layout
   - No database changes required
   - Maintains the "no heavy dependencies" philosophy

3. **Immediate User Benefit**:
   - Solves a common conda workflow pain point
   - More intuitive than recreating environments manually

### **Testing the Feature**
```bash
# Create original environment
uve create myenv 3.11
uve activate myenv
uv pip install numpy pandas

# Clone it
uve clone myenv myenv-copy

# Verify
uve activate myenv-copy
python -c "import numpy; print(numpy.__version__)"
```