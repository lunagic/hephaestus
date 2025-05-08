#!/usr/bin/env bash
set -Eeuo pipefail

PROJECTS=(
	"$HOME/Projects/lunagic/hermes"
	"$HOME/Projects/lunagic/hera"
	"$HOME/Projects/lunagic/hephaestus"
	"$HOME/Projects/lunagic/poseidon"
	"$HOME/Projects/lunagic/athena"
	"$HOME/Projects/lunagic/prometheus"
	"$HOME/Projects/lunagic/environment-go"
	"$HOME/Projects/lunagic/typescript-go"
	"$HOME/Projects/aaronellington/dotfiles"
	"$HOME/Projects/aaronellington/remy"
	"$HOME/Projects/aaronellington/ynab-go"
	"$HOME/Projects/aaronellington/clubhouse"
	"$HOME/Projects/aaronellington/aaron.ellington.io"
	"$HOME/Projects/flight1401/flight1401.com"
)

make build > /dev/null
for PROJECT in "${PROJECTS[@]}" ; do
	cd $PROJECT
	echo =========================================================================
	echo =========================================================================
	echo =========================================================================
	pwd
	hephaestus
	git status
done
