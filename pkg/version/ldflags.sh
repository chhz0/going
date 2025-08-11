#!/bin/bash

VERSION_PKG="github.com/chhz0/going/pkg/version"

get_git_info() {
    # --always: Show the commit hash even if there is no tag.
    # --dirty: If there are modifications in the working directory, append the -dirty suffix.
    # --tags: Use all tags.
    # --abbrev=7: Shorten the commit hash to 7 characters.
    VERSION_INFO=$(git describe --always --dirty=-dev --tags --abbrev=7 2>/dev/null || echo "v0.0.0-dev")

    GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "")
    GIT_COMMIT_STAMP=$(git show -s --format=%ct 2>/dev/null || echo "")
    GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")

    GIT_STATE="clean"
    if ! git diff-index --quiet HEAD -- 2>/dev/null; then
        GIT_STATE="dirty"
    fi

    echo "$VERSION_INFO $GIT_COMMIT $GIT_COMMIT_STAMP $GIT_BRANCH $GIT_STATE"
}

get_build_info() {
    BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    echo "$BUILD_DATE"
}

gen_ldflags() {
    local version_info git_commit git_commit_stamp git_branch git_state
    read -r version_info git_commit git_commit_stamp git_branch git_state <<EOF
$(get_git_info)
EOF

    local build_date
    read -r build_date <<EOF
$(get_build_info)
EOF

    local flags=""
    add_flag() {
        flags="$flags -X '$VERSION_PKG.$1=$2'"
    }

    add_flag "version" "$version_info"
    add_flag "gitCommit" "$git_commit"
    add_flag "gitCommitStamp" "$git_commit_stamp"
    add_flag "gitBranch" "$git_branch"
    add_flag "gitState" "$git_state"
    add_flag "buildDate" "$build_date"

    echo "$flags"
}

main() {
    gen_ldflags
}

main "$@"