#!/bin/sh 
#parse command line arguments 
while [[ $# -gt 0 ]]; do 
    case "$1" in 
        -t|--accessToken)
            ACCESS_TOKEN="$2"
            shift 2
            ;;
        -u|--username)
            USERNAME="$2"
            shift 2
            ;;
        -r|--repo)
            REPO="$2"
            shift 2
            ;;
        -b|--branch)
            BRANCH="$2"
            shift 2
            ;;
        -p|--projectFolder)
            PROJECT_FOLDER="$2"
            shift 2
            ;;
        -s|--sha)
            PARENT_SHA="$2"
            shift 2
            ;;
        *)
            echo "Unrecognized argument $1"
            exit 1
            ;;
    esac
done

if [ -z "$ACCESS_TOKEN" ] || [ -z "$USERNAME" ] || [ -z "$REPO" ] || [ -z "$BRANCH" ] || [ -z "$PROJECT_FOLDER" ]; then 
    echo "Access token not provided"
    exit 1
fi
createTree() {
    local dir="$1"
    local tree_content=""
    while IFS= read -r entry; do 
        local entry_type="blob" 
        local entry_path="$dir/$entry"
        if [ -d "$entry_path" ]; then 
            entry_type="tree"
            local tree_sha=$(createTree "$entry_path")
            tree_content+="{\"path\":\"$entry\",\"mode\":\"040000\",\"type\":\"tree\",\"sha\":\"$tree_sha\"},"
        else 
            #file_content=$(base64 -w 0 <  "$entry")  
            local file_content=$(cat "$entry_path" | jq -Rs .)
            local entry_sha=$(curl -X POST -H "Authorization: Bearer $ACCESS_TOKEN" \
            -d "{\"content\":$file_content, \"encoding\":\"utf-8\"}" \
            "https://api.github.com/repos/$USERNAME/$REPO/git/blobs" | jq -r '.sha') 
            tree_content+="{\"path\":\"$entry\",\"mode\":\"100644\",\"type\":\"$entry_type\",\"sha\":\"$entry_sha\"},"
        fi
    
    done < <(ls -A "$dir")

    tree_content="[${tree_content%,}]"
    echo "$tree_content"
    local tree_sha=$(curl -L \
    -X POST \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "{\"tree\":$tree_content}" \
    "https://api.github.com/repos/$USERNAME/$REPO/git/trees" | jq -r '.sha')
    
    echo "$tree_sha"
}

final_tree_sha=$(createTree "$PROJECT_FOLDER")
echo "$final_tree_sha"
commit_sha=$(curl -L \
    -X POST \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    https://api.github.com/repos/$USERNAME/$REPO/git/commits \
    -d "{\"message\":\"Initial commit\",\"tree\":\"$final_tree_sha\",\"parents\":[\"$PARENT_SHA\"]}" \
    | jq -r '.sha')
echo "commit_sha: $commit_sha"
curl -X PATCH -H "Authorization: token $ACCESS_TOKEN" \
    -d "{\"sha\":\"$commit_sha\",\"force\":true}" \
    "https://api.github.com/repos/$USERNAME/$REPO/git/refs/heads/$BRANCH"