#!/bin/sh
# Copyright (c) 2026 Parable, LLC
# Both the shell and child ignore TERM, exercising group KILL escalation.
trap '' TERM
/bin/sleep 60 &
wait
