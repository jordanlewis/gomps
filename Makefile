# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
#include $(GOROOT)/src/Make.$(GOARCH)
#
#TARG=gomps
#GOFILES=\
#	program.go\
#
#include $(GOROOT)/src/Make.pkg
TARGETS = debug inst token scanner parser eval
all:
	for d in $(TARGETS); do (cd $$d; $(MAKE) && $(MAKE) install); done
	6g gomps.go && 6l -o gomps gomps.6


clean:
	for d in $(TARGETS); do (cd $$d; $(MAKE) clean); done
	rm gomps.6 gomps
