#!/usr/bin/env bash
set -e

TASKFILE="${1:-Taskfile.yml}"
[ ! -f "$TASKFILE" ] && { echo "Ошибка: $TASKFILE не найден" >&2; exit 1; }

read_vars() {
  local in_vars=0
  while IFS= read -r line; do
    [[ "$line" =~ ^vars: ]] && { in_vars=1; continue; }
    (( in_vars )) || continue
    [[ "$line" =~ ^[a-z] ]] && break   # следующая секция

    # Пропускаем пустые строки и комментарии
    [[ "$line" =~ ^[[:space:]]*$ || "$line" =~ ^[[:space:]]*# ]] && continue

    # Переменная в одинарных кавычках
    if [[ "$line" =~ ^[[:space:]]*([A-Z_][A-Z0-9_]*):[[:space:]]*\'([^\']*)\' ]]; then
      val="${BASH_REMATCH[2]}"
      [[ "$val" != *"{{"* ]] && echo "${BASH_REMATCH[1]}=$val"
    # Переменная в двойных кавычках
    elif [[ "$line" =~ ^[[:space:]]*([A-Z_][A-Z0-9_]*):[[:space:]]*\"([^\"]*)\" ]]; then
      val="${BASH_REMATCH[2]}"
      [[ "$val" != *"{{"* ]] && echo "${BASH_REMATCH[1]}=$val"
    # Переменная без кавычек (до пробела, табуляции или комментария)
    elif [[ "$line" =~ ^[[:space:]]*([A-Z_][A-Z0-9_]*):[[:space:]]*([^[:space:]#][^#]*) ]]; then
      val="${BASH_REMATCH[2]}"
      # Удаляем концевой комментарий (если есть)
      val="${val%%#*}"
      # Убираем пробелы в конце
      val="${val%"${val##*[![:space:]]}"}"
      [[ -n "$val" && "$val" != *"{{"* ]] && echo "${BASH_REMATCH[1]}=$val"
    fi
  done < "$TASKFILE"
}

while IFS= read -r pair; do
  [[ -z "$pair" ]] && continue
  if [[ -n "${GITHUB_ENV:-}" ]]; then echo "$pair" >> "$GITHUB_ENV"; fi
  if [[ -n "${GITHUB_OUTPUT:-}" ]]; then echo "$pair" >> "$GITHUB_OUTPUT"; fi
done < <(read_vars)
