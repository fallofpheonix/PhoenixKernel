#!/usr/bin/env sh
set -eu

if [ "$#" -ne 1 ]; then
	printf 'usage: %s project-name\n' "$0" >&2
	exit 2
fi

name="$1"
case "$name" in
	*[!a-z0-9-]* | "" | -* | *-)
		printf 'invalid project name: use lowercase kebab-case\n' >&2
		exit 2
		;;
esac

root="${ENGINEERING_ROOT:-$HOME/engineering}"
src="$root/infrastructure/templates/project-base"
dst="$root/workspace/active/$name"
note="$root/brain/05_PROJECTS/ACTIVE/$name"

if [ -e "$dst" ]; then
	printf 'project path exists: %s\n' "$dst" >&2
	exit 1
fi

mkdir -p "$dst" "$note"
cp -R "$src"/. "$dst"/
cp "$dst/.env.example" "$dst/.env"

cat > "$note/Project.md" <<EOF
# Project: $name

## One-Liner
TBD.

## Status
PLANNING

## Repo
\`~/engineering/workspace/active/$name\`

## Ports
- API: localhost:TBD
- DB: localhost:TBD

## Database
TBD

## Run Command
\`cd ~/engineering/workspace/active/$name && docker compose up -d\`

## Dependencies On Other Projects
None

## What I Deliver To Others
None

## Links
- [[Architecture]]
- [[Decisions]]
- [[Mistakes]]

## Current Blockers
None

## Last Worked On
$(date +%F)
EOF

touch "$note/Architecture.md" "$note/Decisions.md" "$note/Mistakes.md"

printf 'created project: %s\ncreated note: %s\n' "$dst" "$note"

