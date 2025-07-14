#!/bin/bash

echo "Scanning for .env.example files..."

# Tìm tất cả các file .env.example trong thư mục hiện tại và con
find . -type f -name ".env.example" | while read -r SOURCE; do
  TARGET="$(dirname "$SOURCE")/.env"
  echo "Updating $TARGET from $SOURCE"

  # Load existing .env nếu có
  declare -A existing_env
  if [ -f "$TARGET" ]; then
    while IFS= read -r line; do
      [[ "$line" =~ ^[[:space:]]*# ]] && continue
      if [[ "$line" =~ ^[Ee][Xx][Pp][Oo][Rr][Tt][[:space:]]+([A-Za-z_][A-Za-z0-9_]*)=(.*) ]]; then
        key="${BASH_REMATCH[1]}"
        value="${BASH_REMATCH[2]}"
        # Giữ lại dấu ngoặc kép nếu có
        existing_env["$key"]="$value"
      fi
    done < "$TARGET"
  fi

  # Tạo file tạm
  TEMP_FILE=$(mktemp)

  # Duyệt qua từng dòng trong .env.example
  while IFS= read -r line || [[ -n "$line" ]]; do
    # Nếu là comment hoặc dòng trắng => giữ nguyên
    if [[ "$line" =~ ^[[:space:]]*# || -z "$line" ]]; then
      echo "$line" >> "$TEMP_FILE"
    elif [[ "$line" =~ ^[Ee][Xx][Pp][Oo][Rr][Tt][[:space:]]+([A-Za-z_][A-Za-z0-9_]*)= ]]; then
      VAR_NAME="${BASH_REMATCH[1]}"
      VALUE=${existing_env[$VAR_NAME]}
      echo "export $VAR_NAME=${VALUE:-\"\"}" >> "$TEMP_FILE"
    else
      # Không match export thì giữ nguyên dòng
      echo "$line" >> "$TEMP_FILE"
    fi
  done < "$SOURCE"

  mv "$TEMP_FILE" "$TARGET"
  echo "Synced $TARGET"
done

echo "Done."
