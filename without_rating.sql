select substr(name,1,7)||'/web'||substr(replace(name,'.JPG','.jpg'),8) from pictures where not exists( select picture_name from picture_tags where name = picture_name and meta='Rating' and tag like '*%');
