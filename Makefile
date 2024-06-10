#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

CORELIB_PKG := go-corelibs/rxp
VERSION_TAGS += MAIN
MAIN_MK_SUMMARY := ${CORELIB_PKG}
MAIN_MK_VERSION := v0.8.0

DEPS += golang.org/x/perf/cmd/benchstat

STATS_BENCH       := testdata/bench
STATS_FILE        := ${STATS_BENCH}/${MAIN_MK_VERSION}
STATS_PATH        := ${STATS_BENCH}/${MAIN_MK_VERSION}-d
STATS_FILE_OUTPUT := ${STATS_BENCH}/${MAIN_MK_VERSION}-d/output
STATS_FILE_REGEXP := ${STATS_BENCH}/${MAIN_MK_VERSION}-d/regexp
STATS_FILE_RXP    := ${STATS_BENCH}/${MAIN_MK_VERSION}-d/rxp

STATS_FILES += ${STATS_FILE}
STATS_FILES += ${STATS_FILE_OUTPUT}
STATS_FILES += ${STATS_FILE_REGEXP}
STATS_FILES += ${STATS_FILE_RXP}

.PHONY += benchmark
.PHONY += benchstats-history
.PHONY += benchstats-regexp

define _perl_regexp_rxp
if (m!_Regexp!) { \
  s/_Regexp//; \
  print STDERR "$$_"; \
} elsif (m!_Rxp!) { \
	s/_Rxp//; \
  print STDOUT "$$_"; \
} else { \
  print STDOUT "$$_"; \
  print STDERR "$$_"; \
}; $$_="";
endef

include CoreLibs.mk

benchmark: export BENCH_COUNT=50
benchmark:
	@rm -fv    "${STATS_FILE}" || true
	@rm -rfv   "${STATS_PATH}" || true
	@mkdir -vp "${STATS_PATH}"
	@$(MAKE) bench | egrep -v '^make' > "${STATS_FILE_OUTPUT}"
	@cat "${STATS_FILE_OUTPUT}" \
			| grep -v "_Regexp" \
			> "${STATS_FILE}"
	@cat "${STATS_FILE_OUTPUT}" \
			| perl -pe '$(call _perl_regexp_rxp)' \
			> "${STATS_FILE_RXP}" \
			2> "${STATS_FILE_REGEXP}"
	@shasum ${STATS_FILES}

benchstats-history:
	@pushd ${STATS_BENCH} > /dev/null \
		&& ${CMD} benchstat \
			`ls | egrep -v '\-d$$' | sort -V` \
		&& popd > /dev/null

benchstats-regexp:
	@pushd ${STATS_PATH} > /dev/null \
		&& ${CMD} benchstat \
			`basename ${STATS_FILE_REGEXP}` \
			`basename ${STATS_FILE_RXP}` \
		&& popd > /dev/null
