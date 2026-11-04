#!/bin/sh
# init.sh - one-shot bootstrap for a new project created from go-pwa-forge.
#
# Renames the Go module + all internal imports, creates the runtime dirs
# the app needs to start, then tidies, regenerates templ, and removes itself.
#
# Usage:
#   ./init.sh <new-module-path>
#
# Example:
#   ./init.sh github.com/jrswab/my-app
#
# Run it from the repo root, once, right after "Use this template".

set -eu

# --- args ---------------------------------------------------------------
NEW_MODULE="${1:-}"

if [ "$#" -ne 1 ] || [ -z "$NEW_MODULE" ]; then
	echo "usage: $0 <new-module-path>" >&2
	echo "example: $0 github.com/jrswab/my-app" >&2
	exit 2
fi

# Must be run from the repo root.
if [ ! -f go.mod ]; then
	echo "error: go.mod not found. Run this from the repo root." >&2
	exit 1
fi

# --- discover the current module path -----------------------------------
# Read it live rather than hardcoding, so this keeps working if the template
# module name changes later.
OLD_MODULE=$(awk '/^module / { print $2; exit }' go.mod)

if [ -z "$OLD_MODULE" ]; then
	echo "error: could not read module path from go.mod" >&2
	exit 1
fi

if [ "$OLD_MODULE" = "$NEW_MODULE" ]; then
	echo "error: new module path is identical to the current one ($OLD_MODULE)" >&2
	exit 1
fi

# --- helpers ------------------------------------------------------------
# Portable in-place edit (GNU sed and BSD sed disagree on -i).
sed_i() {
	_script=$1
	shift
	for _f in "$@"; do
		[ -f "$_f" ] || continue
		_tmp=$(mktemp)
		sed "$_script" "$_f" >"$_tmp" && mv -- "$_tmp" "$_f"
	done
}

# Escape a string for safe use on the LHS of a sed s||| (BRE).
esc_lhs() { printf '%s' "$1" | sed 's/[.[\*^$/|]/\\&/g'; }
# Escape a string for safe use on the RHS of a sed s||| replacement.
esc_rhs() { printf '%s' "$1" | sed 's/[&/|]/\\&/g'; }

echo ">> renaming module"
echo "   $OLD_MODULE  ->  $NEW_MODULE"

# --- 1. module line in go.mod (clean, no sed) ---------------------------
go mod edit -module "$NEW_MODULE"

# --- 2. internal imports across all Go + templ sources ------------------
OLD_ESC=$(esc_lhs "$OLD_MODULE")
NEW_ESC=$(esc_rhs "$NEW_MODULE")

# shellcheck disable=SC2044
for f in $(find . -type f \( -name '*.go' -o -name '*.templ' \) -not -path './vendor/*'); do
	sed_i "s|$OLD_ESC|$NEW_ESC|g" "$f"
done

echo ">> set CORS_ALLOWED_ORIGINS in the runtime environment if you need cross-origin access"

# --- 3. create runtime dirs the app expects on startup ------------------
# main.go opens logs/site.log (O_CREATE makes the file, not the dir) and
# serves media/, so both dirs must exist in a fresh checkout.
echo ">> ensuring runtime dirs exist (logs/, media/)"
mkdir -p logs media
[ -f logs/.gitkeep ] || touch logs/.gitkeep
[ -f media/.gitkeep ] || touch media/.gitkeep

# --- 4. tidy + regenerate (best effort) ---------------------------------
if command -v templ >/dev/null 2>&1; then
	echo ">> templ generate"
	templ generate
else
	echo ">> skipping templ generate (templ not installed)" >&2
fi

echo ">> go mod tidy"
go mod tidy

# --- 5. clean up --------------------------------------------------------
echo ">> removing init script"
rm -f -- "$0"

cat <<EOF

done. next steps:
  - review manifest.json (name / short_name / description are still template defaults)
  - review README.md (module name, deploy paths, secrets)
  - run: go build ./...
EOF
