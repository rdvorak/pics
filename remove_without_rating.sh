cd $HOME/go/src/gallery
cat without_rating.sql | sqlite3 pics.db | while read file; do FILE="$HOME/Pictures/$file"; test -f "$FILE" && echo "$file" && rm "$FILE"; done;
