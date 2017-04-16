for YEAR
do
	find $HOME/Pictures/$YEAR -type f -name "*.js" | while read f; do curl -POST --data-binary @"$f" http://localhost:8081/gallery/submit/metadata;
done;
done
