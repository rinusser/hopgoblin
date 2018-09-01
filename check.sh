#!/bin/bash
# Copyright 2018 Richard Nusser
# Licensed under GPLv3 (see http://www.gnu.org/licenses/)

shortlog=1
check_test_files_documentation=0

# On Windows systems with Cygwin installed "find" and "sort" without a path might execute Windows's built-ins instead.
# Windows's "find" doesn't work the same, Windows's "sort" messes up newlines.
findcmd=/usr/bin/find
sortcmd=/usr/bin/sort


check_imports_grouping() {
  echo "checking import grouping..."
  invalid=`$findcmd * -name '*.go' ! -name 'doc.go' | xargs grep -l import | xargs grep -L "import ("`
  for file in $invalid; do
    echo "  INVALID: $file"
  done
  echo "...done"
  echo
}

check_imports_order() {
  echo -en "checking imports order...\n  "

  test_regex='^(testing|github\.com/stretchr/testify/assert)$'
  proj_regex='hopgoblin'

  files=`$findcmd * -name '*.go'`
  for file in $files; do
    raw_imports=`grep -zoP "(?<=import \()[^)]*" $file | tr -d '\000' | sed "s/^[ a-z_]*\|\/\/.*$//g" | tr -d '" '`
    imports="$raw_imports"
    imports_test=`echo "$imports" | grep -E "$test_regex"`
    imports=`echo "$imports" | grep -vE "$test_regex"`
    imports_proj=`echo "$imports" | grep -E "$proj_regex"`
    imports=`echo "$imports" | grep -vE "$proj_regex"`
    expected=`(echo "$imports_test" | $sortcmd -r; echo "$imports" | $sortcmd; echo "$imports_proj" | $sortcmd) | xargs`
    actual=`echo "$raw_imports" | xargs`

    if [ "$actual" == "$expected" ]; then
      echo -n .
    else
      echo -e "\n  WRONG ORDER: $file"
      echo "    actual:   $actual"
      echo "    expected: $expected"
    fi
  done
  echo -e "\n...done"
  echo
}


check_init_test() {
  echo "checking for init_test.go files..."
  for dir in `$findcmd . -name '*_test.go' | xargs dirname | $sortcmd | uniq`; do
    [ -e "$dir/init_test.go" ] || echo "  MISSING: $dir/init_test.go"
  done
  echo "...done"
  echo
}

check_package_docs() {
  echo "checking for doc.go files..."
  for dir in `$findcmd . -name '*.go' | xargs grep -L "package main"  | xargs dirname | $sortcmd | uniq`; do
    [ -e "$dir/doc.go" ] || echo "  MISSING: $dir/doc.go"
  done
  echo "...done"
  echo
}


log_ok() {
  if [ $shortlog -eq 0 ]; then
    echo "  OK:      $*"
  else
    [ $prev_status -eq 0 ] && echo -n "  "
    echo -n .
  fi
  prev_status=1
}

log_missing() {
  [ $prev_status -eq 1 ] && echo
  echo "  MISSING: $*"
  prev_status=0
}

prev_status=0
check_missing_docblocks() {
  echo "checking for missing docblocks..."

  export IFS=$'\n'
  files=`$findcmd . -name '*.go'`
  if [ $check_test_files_documentation -eq 0 ]; then
    files=`echo "$files" | grep -vE "_test.go$"`
  fi

  lines=`grep -nE "^(func (\([A-Za-z *]+\) )?|type |var )[A-Z]" $files`

  for line in $lines; do
    file=${line%%:*}
    num=${line#*:}
    num=${num%%:*}
    prevline=`head -n$num $file | tail -n2 | head -n1 | tr -d ' \n\r'`
    logline="$file#$num: ${line##*|}";
    if [[ "$prevline" == "*/" ]]; then
      log_ok "$logline"
    else
      log_missing "$logline"
    fi
  done
  [ $prev_status -eq 1 ] && echo
  echo "...done"
  echo
}


check_file_headers() {
  echo "checking file headers..."
  expected=`echo -ne "// Copyright 2018 Richard Nusser\n// Licensed under GPLv3 (see http://www.gnu.org/licenses/)\n\n"`
  for file in `$findcmd . -name '*.go'`; do
    actual=`head -n3 $file | tr -d '\r'`
    [ "$expected" == "$actual" ] || echo "  WRONG HEADER: $file"
  done
  echo "...done"
  echo
}


[ -f check-private.sh ] && ./check-private.sh

check_imports_grouping
check_imports_order
check_init_test
check_package_docs
check_missing_docblocks
check_file_headers

